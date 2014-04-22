package subchal

import (
    "testing"
    "time"
)

func Test_Stoptime_DayLater(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    initialTime, _ := time.Parse("15:04:05", "11:32:19")
    st := Stoptime{
        "A20131215WKD_WOLF_NS_01",
        initialTime,
        initialTime,
        nil,
        1,
        "",
    }
    newSt := st.DayLater()

    switch false {
    case newSt.TripID == st.TripID:
        t.Log("Stoptime.DayLater() didn't preserve TripID")
        t.FailNow()
    case 24 * time.Hour == newSt.ArrivalTime.Sub(st.ArrivalTime):
        t.Log("Stoptime.DayLater() didn't transform ArrivalTime")
        t.FailNow()
    case 24 * time.Hour == newSt.DepartureTime.Sub(st.DepartureTime):
        t.Log("Stoptime.DayLater() didn't transform DepartureTime")
        t.FailNow()
    }
}

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
        firstSt := stSlice[0]
        lastSt := stSlice[len(stSlice) - 1]
        if lastSt.ArrivalTime.Sub(firstSt.ArrivalTime) < 24 * time.Hour {
            t.Log("Stoptimes for stop", stopID, "were not copied to the next day")
            t.FailNow()
        }

        for _, st := range stSlice {
            if st.Stop.StopID != stopID {
                t.Log("Stoptime at stop", st.Stop.StopID, "doesn't match stoptimes key", stopID)
                t.FailNow()
            }
            if st.Stop.IsStation() {
                t.Log("Stoptime at stop", st.Stop.StopID, "is a station rather than a stop")
                t.FailNow()
            }
        }
    }
}
