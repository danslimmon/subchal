package subchal

import (
    "os"
    "time"
    "encoding/csv"
)

// A transfer from one route to another, as included in a Walk.
type Routeswitch struct {
    FromStop *Stop
    ToStop *Stop
    FromRoute *Route
    ToRoute *Route
}

// A complete walk through the subway, touching every station at least once.
type Walk struct {
    StartStop *Stop
    EndStop *Stop
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


// Determines the next stoptime from the given Stop after the given time.
//
// Also returns an integer indicating the number of times we passed midnight.
func NextStoptime(s *Stop, t time.Time) (*Stoptime, int, error) {
    for _, st := range s.Stoptimes {
        if st.DepartureTime.After(t) {
            return st, 0, nil
        }
    }

    // We had to go past midnight
    return s.Stoptimes[0], 1, nil
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
func LoadWalk(csvPath string, stops map[string]*Stop, routes map[string]*Route) (*Walk, error) {
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
                sw.FromStop = stops[rec[i]]
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
                    sw.ToStop = stops[rec[i]]
                }
            case "from_route_id":
                sw.FromRoute = routes[rec[i]]
            case "to_route_id":
                sw.ToRoute = routes[rec[i]]
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
