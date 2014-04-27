package subchal

import (
    "testing"
    "time"
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
                t.FailNow()
            }
            if st.Stop.IsStation() {
                t.Log("Stoptime at stop", st.Stop.StopID, "is a station rather than a stop")
                t.FailNow()
            }
        }
    }
}

func Test_ParseTime(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    rslt, err := ParseTime("00:00:00")
    switch false {
    case err == nil:
        t.Log("Got error from ParseTime:", err)
        t.FailNow()
    case VerifyReferenceDay(rslt) && rslt.Hour() == 0 && rslt.Minute() == 0 && rslt.Second() == 0:
        t.Log("ParseTime produced an incorrect time:", rslt.Format("2006-01-02 15:04:05 MST"))
        t.FailNow()
    }

    rslt, err = ParseTime("13:59:06")
    switch false {
    case err == nil:
        t.Log("Got error from ParseTime:", err)
        t.FailNow()
    case VerifyReferenceDay(rslt) && rslt.Hour() == 13 && rslt.Minute() == 59 && rslt.Second() == 6:
        t.Log("ParseTime produced an incorrect time:", rslt.Format("2006-01-02 15:04:05 MST"))
        t.FailNow()
    }

    rslt, err = ParseTime("25:00:14")
    switch false {
    case err == nil:
        t.Log("Got error from ParseTime:", err)
        t.FailNow()
    case VerifyReferenceDay(rslt) && rslt.Hour() == 1 && rslt.Minute() == 0 && rslt.Second() == 14:
        t.Log("ParseTime produced an incorrect time:", rslt.Format("2006-01-02 15:04:05 MST"))
        t.FailNow()
    }
}

// Verifies that the given time.Time occurs on the reference day (2006-01-02)
func VerifyReferenceDay(t time.Time) bool {
    return (t.Year() == 2006 && t.YearDay() == 2)
}
