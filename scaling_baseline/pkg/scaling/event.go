package scaling

type Event struct {
	EventType       string
	TargetDomainID  DomainID
	TargetServiceID ServiceID
}
