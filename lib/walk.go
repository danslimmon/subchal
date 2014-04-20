package subchal

// A complete walk through the subway, touching every station at least once.
type Walk struct {
    StartStop *Stop
    EndStop *Stop
    Transfers []*Transfer
}
