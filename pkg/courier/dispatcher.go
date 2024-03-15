package courier

type Dispatcher interface {
	// ProcessDispatchResponse stores responses
	ProcessDispatchResponse(response DispatchResponse)
	// GenerateCouriers produces a minimal collection of couriers (1:1 with recipient)
	GenerateCouriers() (couriers []Courier)
	// SendCouriersAndWait dispatches couriers and blocks until all couriers have returned
	SendCouriersAndWait(couriers []Courier) (response []DispatchResponse)
}

type DispatcherImpl struct {
}
