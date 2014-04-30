package subchal

import (
    "testing"

    "time"
    "io/ioutil"
    "strings"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
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

func Test_TimeToTransfer(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    db, err := sql.Open("sqlite3", "../test-data/subchal.sqlite")
    if err != nil {
        t.Log("Error opening SQLite database:", err)
        t.FailNow()
    }

    // Test a normal transfer
    m, err := TimeToTransfer(db, "WFELL_S", "WFELL_S", "WOLF", confidentParseTime("09:40:00"))
    if err != nil {
        t.Log("Error calculating transfer time:", err)
        t.FailNow()
    }
    if m != 30 * 60 + 90 {
        t.Log("Got incorrect transfer time:", m)
        t.FailNow()
    }

    // Test a transfer that cycles through midnight
    m, err = TimeToTransfer(db, "TWINS_N", "TWINS_E", "TROUT", confidentParseTime("15:12:00"))
    if err != nil {
        t.Log("Error calculating transfer time:", err)
        t.FailNow()
    }
    if m != 21 * 1800 + 90 {
        t.Log("Got incorrect transfer time:", m)
        t.FailNow()
    }

    // Test a transfer that doesn't exist
    m, err = TimeToTransfer(db, "EYRIE_W", "DFORT_N", "WOLF", confidentParseTime("12:25:00"))
    if err == nil {
        t.Log("Didn't get an error from TimeToTransfer when we should've")
        t.FailNow()
    }

}


func Test_LoadWalk(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    wkLineCount, _ := CountFileLines("../test-data/walk.txt")
    wk, err := LoadWalk("../test-data/walk.txt")

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
