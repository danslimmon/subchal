package subchal

type Stop struct {
    StopID string
    StopCode string
    StopName string
    StopDesc string
    StopLat float64
    StopLon float64
    ZoneID string
    StopURL string
    LocationType int
    ParentStation *Stop

    ParentStationID string
    Stoptimes []*Stoptime
    Transfers []*Transfer
}

func (s *Stop) IsStation() (bool) {
    return (s.ParentStation == nil)
}


// Associates the array of matching Stoptimes with each Stop.
func AddStoptimesToStops(stops map[string]*Stop, stoptimes map[string][]*Stoptime) error {
    for stopID, stop := range stops {
        stop.Stoptimes = stoptimes[stopID]
    }
    return nil
}


// Associates the array of matching Transfers with each station
//
// After running this function, each station's Transfers property will contain a list of
// all *Transfers with a FromStation matching that station.
func AddTransfersToStations(stations map[string]*Stop, transfers map[string][]*Transfer) error {
    for stopID, station := range stations {
        station.Transfers = transfers[stopID]
    }
    return nil
}
