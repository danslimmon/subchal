package subchal

import (
    "os"
    "fmt"
    "time"
    "io/ioutil"
    "database/sql"

    "launchpad.net/goyaml"
)

// A transfer from one route to another, or a turnaround, as included in a Walk.
type Step struct {
    FromStation string `from_station`
    ToStation string `to_station`
    ToRoute string `to_route`
}

// A complete walk through the subway, touching every station at least once.
type Walk struct {
    StartStation string `start_station`
    StartRoute string `start_route`
    // StartTime is tagged with a nonexistent variable name so that
    // goyaml.Marshal doesn't try to populate it
    StartTime time.Time `nonexistent`
    EndStation string `end_station`
    Steps []Step `steps`

    // Gets converted to a time.Time and placed in StartTime
    StartTimeStr string `start_time`
    // Counts visits to stations so we can make sure we hit them all
    StationVisits map[string]int
}


type SimulationError struct { s string }
func (e SimulationError) Error() string { return e.s }


// Simulates an execution of the Walk through the subway system.
//
// Returns the number of seconds that the Walk would take, from the beginning
// of the first trip to the end of the last.
func (wk *Walk) RunSim(db *sql.DB) (int, error) {
    dur := 0
    t := wk.StartTime
    dt, err := TimeToTransfer(db, t,
                              wk.StartStation, wk.StartStation, wk.StartRoute,
                              wk.Steps[0].FromStation)
    if err != nil { return 0, err }
    t.Add(time.Duration(dt) * time.Second)
    dur += int(t.Sub(wk.StartTime).Seconds())
    return dur, nil
    return 0, nil
}


// Finds the stop_id of all stations and stops with the given name
func StopIDsFromName(db *sql.DB, stopName string) ([]string, error) {
    stopIDs := make([]string, 0)
    stopRows, err := db.Query(`
        SELECT s.stop_id
        FROM stops s
        WHERE s.stop_name = ?;
        `, stopName)
    if err != nil { return stopIDs, err }
    for stopRows.Next() {
        var stopID string
        err = stopRows.Scan(&stopID)
        if err != nil { return stopIDs, SimulationError{"Error getting stop ID from database"} }
        stopIDs = append(stopIDs, stopID)
    }
    if len(stopIDs) < 1 {
        return stopIDs, SimulationError{fmt.Sprintf("No station with name '%s'", stopName)}
    }
    return stopIDs, nil
}

// Finds the stop_id the station in fromStops with service to a station in toStops
//
// Returns only the first result found, if there's more than one result
func StopGoingToward(db *sql.DB, []string fromStops, []string toStops) ([]string, error) {
    stopIDs := make([]string, 0)
    tripRows, err := db.Query(`
        SELECT t.trip_id,
        FROM stop_times st
            JOIN trips t ON st.trip_id = t.trip_id
        WHERE st.stop_id IN ?
           OR st.stop_id IN ?
        LIMIT 1;
        `, fromStops, toStops)
    if err != nil { return stopIDs, err }
    for stopRows.Next() {
        var stopID string
        err = stopRows.Scan(&stopID)
        if err != nil { return stopIDs, SimulationError{"Error getting stop ID from database"} }
        stopIDs = append(stopIDs, stopID)
    }
    if len(stopIDs) < 1 {
        return stopIDs, SimulationError{fmt.Sprintf("No station with name '%s'", stopName)}
    }
    return stopIDs, nil
}



// Determines the number of seconds it will take to transfer to the
// given route from the given stop at the given time.
func TimeToTransfer(db *sql.DB, t time.Time,
                    fromStation string, toStation string, toRoute string,
                    towardStationName string) (int, error) {

    towardStation, err := StopIDsFromName(db, towardStationName)
    if err != nil { return 0, err }

//    toStopID, err := StopGoingToward(towardStation)
//    if err != nil { return 0, err }

    // First we need to find the stops that are contained in the given
    // station.
    stopRows, err := db.Query(`
        SELECT s.stop_id
        FROM stops s
        WHERE s.parent_station = ?;
        `, toStation)
    if err != nil { return 0, err }
    toStops := make([]string, 0)
    for stopRows.Next() {
        var stopID string
        stopRows.Scan(&stopID)
        toStops = append(toStops, stopID)
    }
    if err != nil {
        return 0, SimulationError{fmt.Sprintf("No stops with parent station '%s'", towardStation)}
    }

    timeStringRows, err := db.Query(`
        SELECT st.departure_time
        FROM stop_times st
            JOIN trips t ON st.trip_id = t.trip_id
            JOIN routes r ON t.route_id = r.route_id
        WHERE st.stop_id = ?
          AND r.route_id = ?
        ORDER BY departure_time ASC;
    `, toStation, toRoute)
    if err != nil { return 0, err }
    var timeStrings []string
    for timeStringRows.Next() {
        var timeString string
        timeStringRows.Scan(&timeString)
        timeStrings = append(timeStrings, timeString)
    }

    times := make([]time.Time, 0)
    for _, timeStr := range timeStrings {
        t, err := ParseTime(timeStr)
        if err != nil { return 0, err }
        times = append(times, t)
    }
    if len(times) < 1 {
        return 0, SimulationError{fmt.Sprintf("No stoptimes found for stop %s and route %s", toStation, toRoute)}
    }

    // Find the next 2 stop times
    nearbyTimes := make([]time.Time, 0)
    for _, stoptime := range times {
        if stoptime.After(t) {
            nearbyTimes = append(nearbyTimes, stoptime)
        }
        if len(nearbyTimes) == 2 {
            break
        }
    }
    if len(nearbyTimes) < 2 {
        // We had to go past midnight.
        nearbyTimes = append(nearbyTimes, times[0].Add(24 * time.Hour))
        if len(nearbyTimes) < 2 {
            nearbyTimes = append(nearbyTimes, times[1].Add(24 * time.Hour))
        }
    }
    interval := int(nearbyTimes[1].Sub(nearbyTimes[0]).Seconds())

    // Okay, now that we have the interval between trains, we need to add the time
    // it takes to run from platform to platform
    transferTimeRow := db.QueryRow(`
        SELECT t.min_transfer_time
        FROM transfers t
        WHERE t.from_stop_id = (
            SELECT s.parent_station
            FROM stops s
            WHERE s.stop_id = ?
            )
          AND t.to_stop_id = (
            SELECT s.parent_station
            FROM stops s
            WHERE s.stop_id = ?
            )
        LIMIT 1;
    `, fromStation, toStation)
    var transferTime int
    err = transferTimeRow.Scan(&transferTime)
    if err != nil {
        return 0, SimulationError{fmt.Sprintf("No transfers possible from %s to %s", fromStation, toStation)}
    }

    return (transferTime + interval)/2, nil
}


// Finds the number of seconds it will take to travel the given segment near the given time.
func TimeToTravel(db *sql.DB, fromStop string, toStop string, route string, t time.Time) (int, error) {
    formattedTime := t.Format("15:04:05")
    dayCycled := 0

    // Find the next trip on route `route` that leaves from `fromStop` in the direction of `toStop`
    tripIDRow := db.QueryRow(`
        SELECT t.trip_id
        FROM trips t
            JOIN stop_times st ON t.trip_id = st.trip_id
        WHERE st.stop_id = ?
          AND t.route_id = ?
          AND st.departure_time > ?
        ORDER BY st.departure_time ASC
        LIMIT 1;
    `, fromStop, route, formattedTime)
    var tripID string
    err := tripIDRow.Scan(&tripID)
    if err != nil {
        // We couldn't find any trips after `t`, so we'll cycle over to
        // the next day
        dayCycled = 1
        tripIDRow := db.QueryRow(`
            SELECT t.trip_id
            FROM trips t
                JOIN stop_times st ON t.trip_id = st.trip_id
            WHERE st.stop_id = ?
              AND t.route_id = ?
            ORDER BY st.departure_time ASC
            LIMIT 1;
        `, fromStop, route)
        err = tripIDRow.Scan(&tripID)
        if err != nil {
            return 0, SimulationError{
                fmt.Sprintf("No trips from %s to %s after %s", fromStop, toStop, formattedTime),
            }
        }
    }

    departureTime, err := TimeFromQuery(db, `
        SELECT st.departure_time
        FROM stop_times st
        WHERE st.stop_id = ?
          AND st.trip_id = ?;
    `, fromStop, tripID)
    if err != nil  {
        return 0, err
    }
    if dayCycled == 1 {
        departureTime = departureTime.Add(24 * time.Hour)
    }

    arrivalTime, err := TimeFromQuery(db, `
        SELECT st.arrival_time
        FROM stop_times st
        WHERE st.stop_id = ?
          AND st.trip_id = ?;
    `, toStop, tripID)
    if err != nil {
        return 0, err
    }
    if arrivalTime.Before(departureTime) {
        arrivalTime = arrivalTime.Add(24 * time.Hour)
    }

    return int(arrivalTime.Sub(departureTime).Seconds()), nil
}


// Parses a time.Time out of a single DB row returned by the given query.
//
// Takes the same arguments as QueryRow: a query followed by zero or more
// strings to interpolate into that query.
func TimeFromQuery(db *sql.DB, query string, params ...interface{}) (time.Time, error) {
    var tStr string
    tRow := db.QueryRow(query, params...)
    err := tRow.Scan(&tStr)
    if err != nil {
        return time.Time{}, err
    }

    t, err := ParseTime(tStr)
    if err != nil {
        return time.Time{}, SimulationError{
            fmt.Sprintf("Malformatted time: %s", tStr),
        }
    }

    return t, nil
}

// Loads an initial Walk from the specified YAML
func LoadWalk(yamlPath string) (*Walk, error) {
    f, err := os.Open(yamlPath)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    yamlBytes, err := ioutil.ReadAll(f)
    if err != nil {
        return nil, err
    }

    wk := new(Walk)
    goyaml.Unmarshal(yamlBytes, wk)
    wk.StartTime, err = time.Parse("15:04:05", wk.StartTimeStr)
    if err != nil {
        return nil, err
    }

    return wk, nil
}
