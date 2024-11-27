package scaling

// type time int

type Event struct {
	EventTime       int
	EventType       string
	TargetDomainID  DomainID
	TargetServiceID ServiceID
}
