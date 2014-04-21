package subchal

import (
    "testing"
)

func Test_LoadStoptimes(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    stops, _, _ := LoadStops("../test-data/stops.txt")
    stoptimes, err := LoadStoptimes("../test-data/stop_times.txt", stops)

    switch false {
    case err == nil:
        t.Log("Got error from LoadStoptimes():", err)
        t.FailNow()
    case len(stoptimes) == len(stops):
        t.Log("Didn't load stoptimes for enough stops")
        t.FailNow()
    }

    for stopID, stSlice := range stoptimes {
        for _, st := range stSlice {
            if st.Stop.StopID != stopID {
                t.Log("Stoptime at stop", st.Stop.StopID, "doesn't match stoptimes key", stopID)
            }
            if st.Stop.IsStation() {
                t.Log("Stoptime at stop", st.Stop.StopID, "is a station rather than a stop")
            }
        }
    }
}
