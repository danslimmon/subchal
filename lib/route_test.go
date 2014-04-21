package subchal

import (
    "testing"
)

func Test_LoadRoutes(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    routes, err := LoadRoutes("../test-data/routes.txt")
    switch false {
    case err == nil:
        t.Log("Got error from LoadRoutes():", err)
        t.FailNow()
    case len(routes) > 1:
        t.Log("Didn't load enough routes")
        t.FailNow()
    }
}
