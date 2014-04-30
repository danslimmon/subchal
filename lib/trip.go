package subchal

type Trip struct {
    Route *Route
    ServiceID string
    TripID string
    TripHeadsign string
    DirectionID int64

    Stoptimes []*Stoptime
}
