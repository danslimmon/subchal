package subchal

import (
    "os"
    "sort"
    "time"
    "strconv"
    "encoding/csv"
)

type Stoptime struct {
    TripID string
    ArrivalTime time.Time
    DepartureTime time.Time
    Stop *Stop
    StopSequence int64
    StopHeadsign string
}

// ByArrivalTime implements sort.Interface for []*Stoptime based on the ArrivalTime field.
type ByArrivalTime []*Stoptime
func (bat ByArrivalTime) Len() int { return len(bat) }
func (bat ByArrivalTime) Swap(i, j int) { bat[i], bat[j] = bat[j], bat[i] }
func (bat ByArrivalTime) Less(i, j int) bool { return bat[i].ArrivalTime.Before(bat[j].ArrivalTime) }

// Reads the CSV stoptimes.txt at the given path
//
// 'stops' should be a map as returned by LoadStops().
//
// Returns a map[StopID string][]*Stoptime. The []*Stoptime will be ordered by ArrivalTime.
// We will cycle back through the Stoptimes to create an extra virtual day of them.
func LoadStoptimes(csvPath string, stops map[string]*Stop) (map[string][]*Stoptime, error) {
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

    stoptimes := make(map[string][]*Stoptime)
    for _, rec := range records {
        st := new(Stoptime)
        for i, colName := range colNames {
            switch colName {
            case "trip_id":
                st.TripID = rec[i]
            case "arrival_time":
                st.ArrivalTime, err = ParseTime(rec[i])
            case "departure_time":
                st.DepartureTime, err = ParseTime(rec[i])
            case "stop_id":
                st.Stop = stops[rec[i]]
            case "stop_sequence":
                st.StopSequence, err = strconv.ParseInt(rec[i], 10, 8)
            case "stop_headsign":
                st.StopHeadsign = rec[i]
            }

            if err != nil {
                return nil, err
            }
        }

        stoptimes[st.Stop.StopID] = append(stoptimes[st.Stop.StopID], st)
    }

    // Order each slice of *Stoptimes by ArrivalTime.
    for _, stSlice := range stoptimes {
        sort.Sort(ByArrivalTime(stSlice))
    }

    return stoptimes, nil
}

// Converts a schedule time (e.g. "11:36:14") to a time.Time object.
//
// The time zone will be UTC. The date will be the zero date
// (January 1, 0001)
func ParseTime(timeStr string) (t time.Time, err error) {
    t, err = time.Parse("15:04:05", timeStr)
    return
}
