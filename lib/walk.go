package subchal

import (
    "os"
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
    Routeswitches []Routeswitch

    StationVisits map[string]int
}


// Loads an initial Walk from the specified CSV file.
//
// The following example travels from the Eyrie west to the Iron
// Islands, then back to the Twins to transfer to the Wolf line. From
// there it goes up to Winterfell and then turns around and goes all
// the way south to King's Landing.
//
//    from_stop_id,to_stop_id,from_route_id,to_route_id
//    EYRIE_W,,TROUT,
//    FEISL_W,FEISL_E,TROUT,TROUT
//    TWINS_E,TWINS_N,TROUT,WOLF
//    WFELL_N,WFELL_S,WOLF,WOLF
//    KLAND_S,,WOLF,
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
