package subchal

import (
    "time"
)

type Stoptime struct {
    TripID string
    ArrivalTime time.Time
    DepartureTime time.Time
    Stop *Stop
    StopSequence int64
    StopHeadsign string
}

// ByArrivalTime implements sort.Interface for []*Stoptime based on the ArrivalTime field.
type ByArrivalTime []*Stoptime
func (bat ByArrivalTime) Len() int { return len(bat) }
func (bat ByArrivalTime) Swap(i, j int) { bat[i], bat[j] = bat[j], bat[i] }
func (bat ByArrivalTime) Less(i, j int) bool { return bat[i].ArrivalTime.Before(bat[j].ArrivalTime) }
