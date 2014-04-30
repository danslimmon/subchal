package subchal

import (
    "os"
    "fmt"
    "time"
    "encoding/csv"
    "database/sql"
)

// A transfer from one route to another, as included in a Walk.
type Routeswitch struct {
    FromStop string
    ToStop string
    FromRoute string
    ToRoute string
}

// A complete walk through the subway, touching every station at least once.
type Walk struct {
    StartStop string
    EndStop string
    StartTime time.Time
    Routeswitches []Routeswitch

    StationVisits map[string]int
}


type SimulationError struct { s string }
func (e SimulationError) Error() string { return e.s }


// Simulates an execution of the Walk through the subway system.
//
// Returns the number of seconds that the Walk would take, from the beginning
// of the first trip to the end of the last.
//func (wk *Walk) RunSim() (int, error) {
//    currentStoptime, err := FirstStoptime(wk.StartStop, wk.StartTime)
//    for _, sw := range wk.Routeswitches {
//        currentStoptime, err := NextStoptime(sw, currentStoptime)
//    }
//    return 0, nil
//}


// Determines the number of seconds it will take to transfer to the
// given route from the given stop at the given time.
func TimeToTransfer(db *sql.DB, fromStop string, toStop string, route string, t time.Time) (int, error) {

    timeStringRows, err := db.Query(`
        SELECT st.departure_time
        FROM stop_times st
            JOIN trips t ON st.trip_id = t.trip_id
            JOIN routes r ON t.route_id = r.route_id
        WHERE st.stop_id = ?
          AND r.route_id = ?
        ORDER BY departure_time ASC;
    `, toStop, route)
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
        return 0, SimulationError{fmt.Sprintf("No stoptimes found for stop %s and route %s", toStop, route)}
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
    `, fromStop, toStop)
    if err != nil { return 0, err }
    var transferTime int;
    err = transferTimeRow.Scan(&transferTime)
    if err != nil {
        return 0, SimulationError{fmt.Sprintf("No transfers possible from %s to %s", fromStop, toStop)}
    }

    return (transferTime + interval)/2, nil
}


// Loads an initial Walk from the specified CSV file.
//
// The following example travels from the Eyrie west to the Iron
// Islands, then back to the Twins to transfer to the Wolf line. From
// there it goes up to Winterfell and then turns around and goes all
// the way south to King's Landing.
//
//    from_stop_id,to_stop_id,from_route_id,to_route_id,start_time
//    EYRIE_W,,TROUT,10:20:07
//    FEISL_W,FEISL_E,TROUT,TROUT,
//    TWINS_E,TWINS_N,TROUT,WOLF,
//    WFELL_N,WFELL_S,WOLF,WOLF,
//    KLAND_S,,WOLF,,
func LoadWalk(csvPath string) (*Walk, error) {
    f, err := os.Open(csvPath)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    colNames, err := csvReader.Read()
    if err != nil {
        return nil, err
    }

    records, err := csvReader.ReadAll()
    if err != nil {
        return nil, err
    }

    wk := new(Walk)
    routeswitches := make([]Routeswitch, 0)
    for _, rec := range records {
        sw := new(Routeswitch)
        for i, colName := range colNames {
            if sw == nil {
                // This happens if we just processed a starting line
                continue
            }

            switch colName {
            case "from_stop_id":
                sw.FromStop = rec[i]
            case "to_stop_id":
                if rec[i] == "" {
                    // This is either the beginning or ending of the walk.
                    if len(routeswitches) == 0 {
                        wk.StartStop = sw.FromStop
                        sw = nil
                        break
                    } else {
                        wk.EndStop = sw.FromStop
                        sw = nil
                        break
                    }
                } else {
                    // Normal (non-beginning-or-ending) routeswitch
                    sw.ToStop = rec[i]
                }
            case "from_route_id":
                sw.FromRoute = rec[i]
            case "to_route_id":
                sw.ToRoute = rec[i]
            case "start_time":
                if rec[i] != "" {
                    wk.StartTime, err = time.Parse("15:04:05", rec[i])
                }
            }

            if err != nil {
                return nil, err
            }
        }

        if sw != nil {
            // sw will be nil if we just processed a starting or ending line
            routeswitches = append(routeswitches, *sw)
        }
    }

    wk.Routeswitches = routeswitches
    return wk, nil
}
