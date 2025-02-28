package events

import (
	"os"

	orc "github.com/nasim-samimi/scaling-simulator/pkg/orchestrator"
	"github.com/sirupsen/logrus"
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

type EventID string

// Global logger instance
var log = logrus.New()

// Automatically runs when the package is imported
func init() {
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)       // Log to console
	log.SetLevel(logrus.InfoLevel) // Set log level
}
