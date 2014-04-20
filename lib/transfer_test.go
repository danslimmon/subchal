package subchal

import (
    "testing"
)

func Test_LoadTransfers(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    _, stations, _ := LoadStops("../test-data/stops.txt")
    transfers, err := LoadTransfers("../test-data/transfers.txt", stations)

    switch false {
    case err == nil:
        t.Log("Got error from LoadTransfers():", err)
        t.FailNow()
    case len(transfers) > 0:
        t.Log("Didn't load enough transfers")
        t.FailNow()
    }

    for stopID, transfers := range transfers {
        for _, xfer := range transfers {
            if xfer.FromStation.StopID != stopID {
                t.Log("Transfer from station", xfer.FromStation.StopID, "doesn't match transfers key", stopID)
            }
            if ! (xfer.FromStation.IsStation() && xfer.ToStation.IsStation()) {
                t.Log("Transfer from", xfer.FromStation.StopID, "to", xfer.ToStation.StopID, "includes a stop")
            }
        }
    }
}
