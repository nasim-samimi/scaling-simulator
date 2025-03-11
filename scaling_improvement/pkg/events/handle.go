package events

import (
	"math"
	"sort"
	"time"

	cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"
	src "github.com/nasim-samimi/scaling-simulator/pkg/orchestrator"
)

type EventResult struct {
	Results *cnfg.ResultContext
}

func NewEventResult() *EventResult {
	qosPerCost := make([]float64, 0)
	qos := make([]float64, 0)
	cost := make([]float64, 0)
	durations := make([]float64, 0)
	eventTime := make([]float64, 0)
	return &EventResult{
		Results: &cnfg.ResultContext{
			QosPerCost: qosPerCost,
			Qos:        qos,
			Cost:       cost,
			Durations:  durations,
			EventTime:  eventTime,
		},
	}
}

func ProcessEvents(events []Event, orchestrator *src.Orchestrator) (*cnfg.ResultContext, error) {

	qosPerCost := make([]float64, 0)
	qos := make([]float64, 0)
	cost := make([]float64, 0)
	durations := make([]float64, 0)
	eventTime := make([]float64, 0)
	// rejected := make([]float64, 0)
	test := 0
	for _, event := range events {
		eventID := event.EventID
		log.Info("event:", event)
		log.Info("service:", orchestrator.AllServices[event.TargetServiceID])
		if event.EventType == "allocate" {
			startTime := time.Now()
			allocated, err := orchestrator.Allocate(event.TargetDomainID, event.TargetServiceID, eventID)
			duration := time.Since(startTime)
			log.Info("Allocate:", allocated, orchestrator.QoS, orchestrator.Cost)
			log.Info("Time to allocate:", duration)
			if err != nil {
				log.Info(err)
			}
			if allocated {
				log.Info("/////////////////////")
				log.Info("service is allocated ")
				log.Info("/////////////////////")
			} else {
				log.Info("/////////////////////")
				log.Info("service is rejected ")
				// delete(orchestrator.RunningServices, eventID)
				log.Info("/////////////////////")

			}
			log.Info("//////////////////bw of the domain nodes////////////////////")
			for _, node := range orchestrator.Domains[event.TargetDomainID].ActiveNodes {
				log.Info("node:", node.NodeName, "consumed bandwidth:", node.TotalConsumedBandwidth)
			}
			log.Info("//////////////////bw of the cloud nodes////////////////////")
			for _, node := range orchestrator.Cloud.ActiveNodes {
				log.Info("node:", node.NodeName, "consumed bandwidth:", node.TotalConsumedBandwidth)
			}
			qosPerCost = append(qosPerCost, math.Round(float64(orchestrator.QoS)*1000/float64(orchestrator.Cost))/1000)
			qos = append(qos, float64(orchestrator.QoS))
			cost = append(cost, float64(orchestrator.Cost))
			durations = append(durations, float64(duration.Microseconds())/1000)
			eventTime = append(eventTime, float64(event.EventTime))
			test++
			// if test == 50 {
			// 	break
			// }
		}
		if event.EventType == "deallocate" {
			log.Info("/////////////////////")
			log.Info("Deallocate")
			upgrade := false
			if _, ok := orchestrator.RunningServices[eventID]; ok {
				svc := *orchestrator.RunningServices[eventID]
				if svc.AllocationMode == "Standard" {
					upgrade = true
				}

				orchestrator.Deallocate(event.TargetDomainID, event.TargetServiceID, eventID)
				if upgrade {
					orchestrator.UpgradeServiceIfEnabled(orchestrator.Config.UpgradeHeuristic, svc, event.TargetDomainID) // change this to only one domain.
				}
				orchestrator.BasicNodeReclaim(event.TargetDomainID)
				orchestrator.NodeReclaimIfEnabled(event.TargetDomainID)
			} else {
				log.Info("Service does not exist. rejected?")
			}
			qosPerCost = append(qosPerCost, math.Round(float64(orchestrator.QoS)*1000/float64(orchestrator.Cost))/1000)
			qos = append(qos, float64(orchestrator.QoS))
			cost = append(cost, float64(orchestrator.Cost))
			eventTime = append(eventTime, float64(event.EventTime))
			log.Info("/////////////////////")
			log.Info("//////////////////bw of the domain nodes////////////////////")
			for _, node := range orchestrator.Domains[event.TargetDomainID].ActiveNodes {
				log.Info("node:", node.NodeName, "consumed bandwidth:", node.TotalConsumedBandwidth, "allocated services:", node.AllocatedServices)
			}
			log.Info("//////////////////bw of the cloud nodes////////////////////")
			for _, node := range orchestrator.Cloud.ActiveNodes {
				log.Info("node:", node.NodeName, "consumed bandwidth:", node.TotalConsumedBandwidth, "allocated services:", node.AllocatedServices)
			}
		}
	}
	log.Info("QoS per Cost: ", qosPerCost)

	log.Info("Durations: ", durations)
	return &cnfg.ResultContext{
		QosPerCost: qosPerCost,
		Qos:        qos,
		Cost:       cost,
		Durations:  durations,
		EventTime:  eventTime,
	}, nil
}

func BufferEvents(events []Event, interval float64, orchestrator *src.Orchestrator) (*cnfg.ResultContext, error) {
	initTime := 0.0
	endTime := interval + initTime
	eventsBuffer := EventsBuffer{
		DeallocEvents: make(DeallocEvents, 0),
		AllocEvents:   make(AllocEvents, 0),
	}
	// processedIDs := make(map[src.ServiceID]bool)
	// qosPerCost, qos, cost, durations, eventTime := []float64{}, []float64{}, []float64{}, []float64{}, []float64{}
	results := NewEventResult()
	unprocessedDeallocs := []Event{}
	rejectedEvents := []Event{}

	log.Info("Total events:", len(events), "Interval:", interval)
	log.Info("interval:", interval)
	log.Info("initTime:", initTime)
	log.Info("endTime:", endTime)
	eventIndex := 0
	processedEvents := 0
	processedEventsa := 0
	processedEventsd := 0
	bufferedEventsa := 0
	bufferedEventsd := 0
	for eventIndex < len(events) {
		event := events[eventIndex]
		if event.EventTime >= initTime && event.EventTime <= endTime {
			switch event.EventType {
			case "allocate":
				eventsBuffer.AllocEvents = append(eventsBuffer.AllocEvents, event)
				bufferedEventsa++
			case "deallocate":
				eventsBuffer.DeallocEvents = append(eventsBuffer.DeallocEvents, event)
				bufferedEventsd++
			}
			eventIndex++
			continue
		}
		log.Info("events are buffered")
		if event.EventTime > endTime {
			unprocessedDeallocs = make([]Event, 0)
			deallocateBufferedEvents(orchestrator, eventsBuffer.DeallocEvents, results, endTime, unprocessedDeallocs)
			allocateBufferedEvents(orchestrator, eventsBuffer, results, endTime)
			deallocateBufferedEvents(orchestrator, unprocessedDeallocs, results, endTime, rejectedEvents)
			// unprocessedDeallocs = make([]Event, 0)

			initTime = endTime
			endTime = interval + initTime
			eventsBuffer = EventsBuffer{
				DeallocEvents: make(DeallocEvents, 0),
				AllocEvents:   make(AllocEvents, 0),
			}
			log.Info("initTime:", initTime)
			log.Info("endTime:", endTime)

		}
	}

	if len(eventsBuffer.DeallocEvents) > 0 {
		unprocessedDeallocs = make([]Event, 0)
		deallocateBufferedEvents(orchestrator, eventsBuffer.DeallocEvents, results, endTime, unprocessedDeallocs)
	}
	if len(eventsBuffer.AllocEvents) > 0 {
		allocateBufferedEvents(orchestrator, eventsBuffer, results, endTime)
	}
	if len(unprocessedDeallocs) > 0 {
		deallocateBufferedEvents(orchestrator, unprocessedDeallocs, results, endTime, rejectedEvents)
	}

	log.Info("Total processed events:", processedEvents)
	log.Info("number of events:", len(events))
	log.Info("number of buffered deallocate events:", bufferedEventsd)
	log.Info("number of buffered allocate events:", bufferedEventsa)
	log.Info("number of processed deallocate events:", processedEventsd)
	log.Info("number of processed allocate events:", processedEventsa)
	log.Info("number of services in running services:", len(orchestrator.RunningServices))
	log.Info("remaining services in the running services:", orchestrator.RunningServices)
	return results.Results, nil
}

func sortEvents(allocEvents []Event, allServices src.Services) ([]src.ServiceID, map[src.ServiceID]Event) {
	sortedEventIDs := make([]src.ServiceID, 0, len(allocEvents))
	sortedEvents := make(map[src.ServiceID]Event)
	for _, event := range allocEvents {
		sortedEventIDs = append(sortedEventIDs, event.EventID)
		sortedEvents[event.EventID] = event
	}
	sort.Slice(sortedEventIDs, func(i, j int) bool {
		serviceimpi := allServices[sortedEvents[sortedEventIDs[i]].TargetServiceID].ImportanceFactor
		serviceimpj := allServices[sortedEvents[sortedEventIDs[j]].TargetServiceID].ImportanceFactor
		return float64(sortedEvents[sortedEventIDs[i]].TotalUtil)*serviceimpi > float64(sortedEvents[sortedEventIDs[j]].TotalUtil)*serviceimpj
	})
	return sortedEventIDs, sortedEvents

}

func computeQoSCost(qos float64, cost float64) float64 {
	return math.Round(qos*1000/cost) / 1000
}

func (e *EventResult) appendMetrics(orchestrator *src.Orchestrator, endTime float64, duration float64) {
	// *qosPerCost = append(*qosPerCost, computeQoSCost(float64(orchestrator.QoS), float64(orchestrator.Cost)))
	// *qos = append(*qos, float64(orchestrator.QoS))
	// *cost = append(*cost, float64(orchestrator.Cost))
	// if duration != 0 {
	// 	*durations = append(*durations, duration)
	// }
	// *eventTime = append(*eventTime, endTime)
	e.Results.QosPerCost = append(e.Results.QosPerCost, computeQoSCost(float64(orchestrator.QoS), float64(orchestrator.Cost)))
	e.Results.Qos = append(e.Results.Qos, float64(orchestrator.QoS))
	e.Results.Cost = append(e.Results.Cost, float64(orchestrator.Cost))
	if duration != 0 {
		e.Results.Durations = append(e.Results.Durations, duration)
	}
	e.Results.EventTime = append(e.Results.EventTime, endTime)
}

func allocateBufferedEvents(orchestrator *src.Orchestrator, eventsBuffer EventsBuffer, results *EventResult, endTime float64) {
	sortedEventIDs, sortedEvents := sortEvents(eventsBuffer.AllocEvents, orchestrator.AllServices)
	if len(sortedEventIDs) != len(eventsBuffer.AllocEvents) {
		log.Info("sorted events do not match the buffer")
	}
	log.Info("sorted event ids for allocation:", sortedEventIDs)
	for _, eventID := range sortedEventIDs {
		aEvent := sortedEvents[eventID]
		log.Info("allocate event:", aEvent)
		eventID := aEvent.EventID
		startTime := time.Now()
		allocated, err := orchestrator.Allocate(aEvent.TargetDomainID, aEvent.TargetServiceID, eventID)
		duration := time.Since(startTime)
		log.Info("Allocate:", allocated, orchestrator.QoS, orchestrator.Cost)
		log.Info("Time to allocate:", duration)
		if err != nil {
			log.Info(err)
		}
		if allocated {
			log.Info("/////////////////////")
			log.Info("service is allocated ")
			log.Info("/////////////////////")
		} else {
			log.Info("/////////////////////")
			log.Info("service is rejected ")
			// delete(orchestrator.RunningServices, eventID)
			log.Info("/////////////////////")

		}
		results.appendMetrics(orchestrator, endTime, float64(duration.Microseconds()/1000))
	}
}

func deallocateBufferedEvents(orchestrator *src.Orchestrator, eventsBuffer []Event, results *EventResult, endTime float64, unprocessedDeallocs []Event) {
	log.Info("attempt to process deallocate events")
	log.Info("deallocation events:", eventsBuffer)
	for _, dEvent := range eventsBuffer {
		// if _, ok := processedIDs[dEvent.EventID]; !ok {
		log.Info("deallocate event:", dEvent)
		// }
		eventID := dEvent.EventID
		if _, ok := orchestrator.RunningServices[eventID]; ok {
			orchestrator.Deallocate(dEvent.TargetDomainID, dEvent.TargetServiceID, eventID)
			results.appendMetrics(orchestrator, endTime, 0)
			log.Info("service exists in the first round")
			log.Info("/////////////////////")
		} else {
			unprocessedDeallocs = append(unprocessedDeallocs, dEvent)
			log.Info("service does not exist in first round. rejected?")
		}

	}
}
