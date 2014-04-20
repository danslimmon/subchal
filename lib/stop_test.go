package subchal

import (
    "testing"
)

func Test_LoadStops(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    stops, stations, err := LoadStops("../test-data/stops.txt")
    switch false {
    case err == nil:
        t.Log("Got error from LoadStops():", err)
        t.FailNow()
    case len(stops) > 1:
        t.Log("Didn't load enough stops")
        t.FailNow()
    case len(stations) > 1:
        t.Log("Didn't load enough stations")
        t.FailNow()
    }

    for _, stop := range stops {
        if stop.IsStation() {
            t.Log("Found a station in the stops list:", stop.StopID)
            t.FailNow()
        }
        if stations[stop.ParentStation.StopID] == nil {
            t.Log("Station's ParentStation is not in the stations array")
            t.FailNow()
        }
    }

    for _, station := range stations {
        if station.IsStation() == false {
            t.Log("Found a stop in the stations list:", station.StopID)
        }
    }
}
