package subchal

import (
    "os"
    "strconv"
    "encoding/csv"
)

type Transfer struct {
    FromStation *Stop
    ToStation *Stop
    TransferType int64
    MinTransferTime int64

    ToRoute *Route
}

// Reads the CSV transfers.txt at the given path
//
// Returns a map[FromStopID string][]*Transfer
func LoadTransfers(csvPath string, stations map[string]*Stop) (map[string][]*Transfer, error) {
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
    
    transfers := make(map[string][]*Transfer)
    for _, rec := range records {
        t := new(Transfer)
        for i, colName := range colNames {
            switch colName {
            case "from_stop_id":
                t.FromStation = stations[rec[i]]
            case "to_stop_id":
                t.ToStation = stations[rec[i]]
            case "transfer_type":
                t.TransferType, err = strconv.ParseInt(rec[i], 10, 8)
            case "min_transfer_time":
                t.MinTransferTime, err = strconv.ParseInt(rec[i], 10, 16)
            }

            if err != nil {
                return nil, err
            }
        }

        transfers[t.FromStation.StopID] = append(transfers[t.FromStation.StopID], t)
    }

    return transfers, nil
}
