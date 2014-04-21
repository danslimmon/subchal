package subchal

import (
    "testing"
)

func Test_LoadTrips(t *testing.T) {
    t.Parallel()
    SetTestLogger(t)

    routes, _ := LoadRoutes("../test-data/routes.txt")
    stops, _, _ := LoadStops("../test-data/stops.txt")
    stoptimes, _ := LoadStoptimes("../test-data/stoptimes.txt", stops)
    trips, err := LoadTrips("../test-data/trips.txt", routes, stoptimes)

    switch false {
    case err == nil:
        t.Log("Got error from LoadTrips():", err)
        t.FailNow()
    case len(trips) > 1:
        t.Log("Didn't load enough trips")
        t.FailNow()
    }

    for stopID, stSlice := range stoptimes {
        for _, st := range stSlice {
            found := false
            trip, ok := trips[st.TripID]
            if ! ok {
                t.Log("Trip with ID", st.TripID, "not found in 'trips'")
                t.FailNow()
            }
            
            for _, tripST := range trip.Stoptimes {
                if tripST == st {
                    found = true
                    break
                }
            }

            if ! found {
                t.Log("Stoptime at stop", stopID, "missing from trips[", trip.TripID, "]")
                t.FailNow()
            }
        }
    }
}
