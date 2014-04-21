package subchal

import (
    "os"
    "strconv"
    "encoding/csv"
)

type Route struct {
    RouteID string
    RouteShortName string
    RouteLongName string
    RouteDesc string
    RouteType int64
    RouteURL string
}

// Reads the CSV routes.txt at the given path
func LoadRoutes(csvPath string) (map[string]*Route, error) {
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

    routes := make(map[string]*Route)
    for _, rec := range records {
        rt := new(Route)
        for i, colName := range colNames {
            switch colName {
            case "route_id":
                rt.RouteID = rec[i]
            case "route_short_name":
                rt.RouteShortName = rec[i]
            case "route_long_name":
                rt.RouteLongName = rec[i]
            case "route_desc":
                rt.RouteDesc = rec[i]
            case "route_type":
                rt.RouteType, err = strconv.ParseInt(rec[i], 10, 64)
            case "route_url":
                rt.RouteURL = rec[i]
            }

            if err != nil {
                return nil, err
            }
        }
        routes[rt.RouteID] = rt
    }

    return routes, nil
}
