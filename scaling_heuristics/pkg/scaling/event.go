package scaling

// type time int

type Event struct {
	EventTime       int
	EventID         ServiceID
	EventType       string
	TargetDomainID  DomainID
	TargetServiceID ServiceID
}
