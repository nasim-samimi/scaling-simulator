package scaling

type Event struct {
	EventType       string
	TargetDomainID  string
	TargetServiceID TaskID
}
