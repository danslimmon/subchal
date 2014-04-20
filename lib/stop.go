package subchal

import (
    "os"
    "strconv"
    "encoding/csv"
)

type Stop struct {
    StopID string
    StopCode string
    StopName string
    StopDesc string
    StopLat float64
    StopLon float64
    ZoneID string
    StopURL string
    LocationType int
    ParentStation *Stop

    ParentStationID string
}

func (s *Stop) IsStation() (bool) {
    return (s.ParentStation == nil)
}


// Reads the CSV stops.txt at the given path
//
// Return values are:
//
//   stopsByID: A map[StopID string]*Stop describing all stops with a parent station
//   stationsByID: A map[StopID string]*Stop describing all stops that are parents of
//        other stations
//   err
func LoadStops(csvPath string) (map[string]*Stop, map[string]*Stop, error) {
    f, err := os.Open(csvPath)
    if err != nil {
        return nil, nil, err
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    colNames, err := csvReader.Read()
    if err != nil {
        return nil, nil, err
    }

    records, err := csvReader.ReadAll()
    if err != nil {
        return nil, nil, err
    }
    
    stops := make(map[string]*Stop)
    stations := make(map[string]*Stop)
    for _, rec := range records {
        s := new(Stop)
        for i, colName := range colNames {
            switch colName {
            case "stop_id":
                s.StopID = rec[i]
            case "stop_code":
                s.StopCode = rec[i]
            case "stop_name":
                s.StopName = rec[i]
            case "stop_desc":
                s.StopDesc = rec[i]
            case "stop_lat":
                s.StopLat, err = strconv.ParseFloat(rec[i], 64)
            case "stop_lon":
                s.StopLon, err = strconv.ParseFloat(rec[i], 64)
            case "zone_id":
                s.ZoneID = rec[i]
            case "stop_url":
                s.StopURL = rec[i]
            case "parent_station":
                s.ParentStationID = rec[i]
            }

            if err != nil {
                return nil, nil, err
            }
        }

        if s.ParentStationID != "" {
            stops[s.StopID] = s
        } else {
            stations[s.StopID] = s
        }
    }

    // Now that we've loaded all the stops, we can populate the ParentStation
    // attribute for all stops with a parent station
    for _, s := range stops {
        if s.ParentStationID != "" {
            s.ParentStation = stations[s.ParentStationID]
        }
    }
    return stops, stations, nil
}
