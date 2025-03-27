package config

type Heuristic string

type Cost float64

type ResultContext struct {
	QosPerCost []float64
	Qos        []float64
	Cost       []float64
	Durations  []float64
	EventTime  []float64
	Acceptance []float64
}
