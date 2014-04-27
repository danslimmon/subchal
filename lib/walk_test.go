package subchal

import (
    "testing"

    "time"
    "io/ioutil"
    "strings"
)

func CountFileLines(path string) (int, error) {
    contentsBytes, err := ioutil.ReadFile(path)
    contents := string(contentsBytes)
    if err != nil {
        return 0, err
    }
    contents = strings.TrimRight(contents, "\n")
    return strings.Count(contents, "\n") + 1, nil
}

func confidentParseTime(timeStr string) time.Time {
    t, _ := ParseTime(timeStr)
    return t
}

func Test_NextStoptime(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    s := new(Stop)
    s.Stoptimes = []*Stoptime{
        &Stoptime{
            ArrivalTime: confidentParseTime("25:17:06"),
            DepartureTime: confidentParseTime("25:17:06"),
            Stop: s,
        },
        &Stoptime{
            ArrivalTime: confidentParseTime("06:19:55"),
            DepartureTime: confidentParseTime("06:19:55"),
            Stop: s,
        },
        &Stoptime{
            ArrivalTime: confidentParseTime("12:47:00"),
            DepartureTime: confidentParseTime("12:47:00"),
            Stop: s,
        },
        &Stoptime{
            ArrivalTime: confidentParseTime("18:30:31"),
            DepartureTime: confidentParseTime("18:30:31"),
            Stop: s,
        },
    }

    st, days, err := NextStoptime(s, confidentParseTime("10:00:00"))
    switch false {
    case err == nil:
        t.Log("Got error from NextStoptime:", err)
        t.FailNow()
    case days == 0:
        t.Log("Went forward a day when we shouldn't've")
        t.FailNow()
    case st.ArrivalTime.Equal(confidentParseTime("12:47:00")):
        t.Log("Got wrong ArrivalTime from NextStoptime:", st.ArrivalTime.Format("15:04:05"))
        t.FailNow()
    }

    st, days, err = NextStoptime(s, confidentParseTime("18:30:32"))
    switch false {
    case err == nil:
        t.Log("Got error from NextStoptime:", err)
        t.FailNow()
    case days == 1:
        t.Log("Didn't go forward a day when we should've", st)
        t.FailNow()
    case st.ArrivalTime.Equal(confidentParseTime("01:17:06")):
        t.Log("Got wrong ArrivalTime from NextStoptime:", st.ArrivalTime.Format("15:04:05"))
        t.FailNow()
    }
}


func Test_LoadWalk(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    wkLineCount, _ := CountFileLines("../test-data/walk.txt")
    stops, _, _ := LoadStops("../test-data/stops.txt")
    routes, _ := LoadRoutes("../test-data/routes.txt")
    wk, err := LoadWalk("../test-data/walk.txt", stops, routes)

    switch false {
    case err == nil:
        t.Log("Got error from LoadWalk():", err)
        t.FailNow()
    case len(wk.Routeswitches) == wkLineCount - 3:
        t.Log("Didn't load the right number of Routeswitches; loaded",
              len(wk.Routeswitches), "instead of", wkLineCount - 3)
        t.FailNow()
    }
}
