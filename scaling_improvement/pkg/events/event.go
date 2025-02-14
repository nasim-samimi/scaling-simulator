package events

import (
	orc "github.com/nasim-samimi/scaling-simulator/pkg/orchestrator"
)

type Event struct {
	EventTime       float64
	EventID         orc.ServiceID
	EventType       string
	TargetDomainID  orc.DomainID
	TargetServiceID orc.ServiceID
	TotalUtil       int
}

type DeallocEvents []Event

type AllocEvents []Event

type EventsBuffer struct {
	DeallocEvents DeallocEvents
	AllocEvents   AllocEvents
}
