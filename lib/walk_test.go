package subchal

import (
    "io/ioutil"
    "testing"
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
