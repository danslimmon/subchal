package subchal

type Transfer struct {
    FromStation *Stop
    ToStation *Stop
    TransferType int64
    MinTransferTime int64

    ToRoute *Route
}
