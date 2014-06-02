package subchal

import (
    "testing"

    "time"
    "strings"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func confidentParseTime(timeStr string) time.Time {
    t, _ := ParseTime(timeStr)
    return t
}

func Test_StopIDFromName(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    db, err := sql.Open("sqlite3", "../test-data/subchal.sqlite")
    if err != nil {
        t.Log("Error opening SQLite database:", err)
        t.FailNow()
    }

    stopIDs, err := StopIDsFromName(db, "Winterfell")
    if err != nil {
        t.Log("Error determining station ID from name:", err)
        t.FailNow()
    }
    for _, stopID := range stopIDs {
        if ! strings.HasPrefix(stopID, "WFELL") {
            t.Log("Got incorrect stop ID:", stopID)
            t.FailNow()
        }
    }

    stopIDs, err = StopIDsFromName(db, "Bloopty Spoop")
    if err == nil {
        t.Log("Expected error when getting stop ID for nonexistent station", err)
        t.FailNow()
    }
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
    s, err := TimeToTransfer(db, confidentParseTime("09:40:00"), "Winterfell", "Winterfell", "Wolf", "King's Landing")
    if err != nil {
        t.Log("Error calculating transfer time:", err)
        t.FailNow()
    }
    if s != 30 * 60 + 90 {
        t.Log("Got incorrect transfer time:", s)
        t.FailNow()
    }

    // Test a transfer that cycles through midnight
    s, err = TimeToTransfer(db, confidentParseTime("15:12:00"), "Twins", "Twins", "Trout", "Eyrie")
    if err != nil {
        t.Log("Error calculating transfer time:", err)
        t.FailNow()
    }
    if s != 21 * 1800 + 90 {
        t.Log("Got incorrect transfer time:", s)
        t.FailNow()
    }

    // Test a transfer that doesn't exist
    s, err = TimeToTransfer(db, confidentParseTime("12:25:00"), "Eyrie", "Dreadfort", "Wolf", "Winterfell")
    if err == nil {
        t.Log("Didn't get an error from TimeToTransfer when we should've")
        t.FailNow()
    }

}


func Test_TimeToTravel(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    db, err := sql.Open("sqlite3", "../test-data/subchal.sqlite")
    if err != nil {
        t.Log("Error opening SQLite database:", err)
        t.FailNow()
    }

    // Test a normal leg of a walk
    s, err := TimeToTravel(db, "EYRIE_W", "FEISL_W", "TROUT", confidentParseTime("11:58:00"))
    if err != nil {
        t.Log("Error calculating travel time:", err)
        t.FailNow()
    }
    if s != 2 * 3600 {
        t.Log("Got incorrect travel time:", s)
        t.FailNow()
    }

    // Test a leg that goes through midnight
    s, err = TimeToTravel(db, "DFORT_S", "KLAND_S", "WOLF", confidentParseTime("22:46:00"))
    if err != nil {
        t.Log("Error calculating travel time:", err)
        t.FailNow()
    }
    if s != 2 * 3600 {
        t.Log("Got incorrect travel time:", s)
        t.FailNow()
    }

    // Test a leg where the next trip doesn't start until the next day
    s, err = TimeToTravel(db, "TWINS_E", "EYRIE_E", "TROUT", confidentParseTime("22:46:00"))
    if err != nil {
        t.Log("Error calculating travel time:", err)
        t.FailNow()
    }
    if s != 3600 {
        t.Log("Got incorrect travel time:", s)
        t.FailNow()
    }

    // Test an impossible leg
    s, err = TimeToTravel(db, "WFELL_S", "EYRIE_E", "WOLF", confidentParseTime("10:00:00"))
    if err == nil {
        t.Log("Didn't get an error from TimeToTravel() when we should've")
        t.FailNow()
    }
}


func Test_LoadWalk(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    wk, err := LoadWalk("../test-data/walk.yaml")

    switch false {
    case err == nil:
        t.Log("Got error from LoadWalk():", err)
        t.FailNow()
    case ! (wk.StartTime.Equal(time.Time{})):
        t.Log("Time not populated for Walk", wk.StartTime)
        t.FailNow()
    case len(wk.Steps) > 2:
        t.Log("Not enough steps loaded for walk")
        t.FailNow()
    }
}
