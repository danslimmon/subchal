package subchal

import (
    "os"
    "strconv"
    "encoding/csv"
)

type Trip struct {
    Route *Route
    ServiceID string
    TripID string
    TripHeadsign string
    DirectionID int64

    Stoptimes []*Stoptime
}

// Reads the CSV trips.txt at the given path
//
// 'routes' should be a map as returned by LoadRoutes().
// 'stoptimes' should be a map as returned by LoadStoptimes().
//
// Returns a map[TripID string]*Trip.
func LoadTrips(csvPath string, routes map[string]*Route, stoptimes map[string][]*Stoptime) (map[string]*Trip, error) {
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

    trips := make(map[string]*Trip)
    for _, rec := range records {
        tr := new(Trip)
        for i, colName := range colNames {
            switch colName {
            case "route_id":
                tr.Route = routes[rec[i]]
            case "service_id":
                tr.ServiceID = rec[i]
            case "trip_id":
                tr.TripID = rec[i]
            case "trip_headsign":
                tr.TripHeadsign = rec[i]
            case "direction_id":
                tr.DirectionID, err = strconv.ParseInt(rec[i], 10, 8)
            }

            if err != nil {
                return nil, err
            }
        }

        trips[tr.TripID] = tr
    }

    // Populate the Stoptimes attribute
    for _, stSlice := range stoptimes {
        for _, st := range stSlice {
            trips[st.TripID].Stoptimes = append(trips[st.TripID].Stoptimes, st)
        }
    }

    return trips, nil
}
