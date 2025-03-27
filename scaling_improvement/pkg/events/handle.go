package events

/*
#include <pthread.h>
#include <time.h>
#include <stdio.h>

static unsigned long long int getProcessTime() {
    struct timespec t;
    if (clock_gettime(CLOCK_PROCESS_CPUTIME_ID, &t)) {
        perror("clock_gettime");
        return 0;
    }
    return t.tv_sec * 1000000000LL + t.tv_nsec;
}
*/
import "C"
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
	acceptance := make([]float64, 0)
	// rejected := make([]float64, 0)
	test := 0
	e := 0
	for _, event := range events {
		eventID := event.EventID
		log.Info("event:", event)
		log.Info("service:", orchestrator.AllServices[event.TargetServiceID])
		if event.EventType == "allocate" {
			// startTime := time.Now()
			startTime := C.getProcessTime()
			allocated, err := orchestrator.Allocate(event.TargetDomainID, event.TargetServiceID, eventID)
			endTime := C.getProcessTime()
			cpuTime := endTime - startTime
			duration := float64(cpuTime) / 1000000
			// duration := time.Since(startTime)
			log.Info("Allocate:", allocated, orchestrator.QoS, orchestrator.Cost)
			log.Info("Time to allocate:", duration)
			if err != nil {
				log.Info(err)
			}
			if allocated {
				log.Info("/////////////////////")
				log.Info("service is allocated ")
				e++
				log.Info("/////////////////////")
			} else {
				log.Info("/////////////////////")
				log.Info("service is rejected ")
				// delete(orchestrator.RunningServices, eventID)
				log.Info("/////////////////////")

			}

			qosPerCost = append(qosPerCost, math.Round(float64(orchestrator.QoS)*1000/float64(orchestrator.Cost))/1000)
			qos = append(qos, float64(orchestrator.QoS))
			cost = append(cost, float64(orchestrator.Cost))
			durations = append(durations, float64(duration))
			eventTime = append(eventTime, float64(event.EventTime))

			test++
			// if test == 50 {
			// 	break
			// }
		}
		if event.EventType == "deallocate" {
			log.Info("/////////////////////")
			log.Info("Deallocate")

			if _, ok := orchestrator.RunningServices[eventID]; ok {
				// svc := *orchestrator.RunningServices[eventID]

				orchestrator.Deallocate(event.TargetDomainID, event.TargetServiceID, eventID)

				// orchestrator.UpgradeServiceIfEnabled(orchestrator.Config.UpgradeHeuristic, svc, event.TargetDomainID) // change this to only one domain.

				orchestrator.BasicNodeReclaim(event.TargetDomainID)
				// orchestrator.NodeReclaimIfEnabled(event.TargetDomainID)
			} else {
				log.Info("Service does not exist. rejected?")
			}
			qosPerCost = append(qosPerCost, math.Round(float64(orchestrator.QoS)*1000/float64(orchestrator.Cost))/1000)
			qos = append(qos, float64(orchestrator.QoS))
			cost = append(cost, float64(orchestrator.Cost))
			eventTime = append(eventTime, float64(event.EventTime))
		}
	}
	log.Info("QoS per Cost: ", qosPerCost)

	log.Info("Durations: ", durations)
	acceptance = append(acceptance, float64(e)/float64(len(events)))
	return &cnfg.ResultContext{
		QosPerCost: qosPerCost,
		Qos:        qos,
		Cost:       cost,
		Durations:  durations,
		EventTime:  eventTime,
		Acceptance: acceptance,
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
	unprocessedDeallocs := map[src.ServiceID]Event{}

	log.Info("Total events:", len(events), "Interval:", interval)
	log.Info("interval:", interval)
	log.Info("initTime:", initTime)
	log.Info("endTime:", endTime)
	eventIndex := 0
	// processedEvents := 0
	// processedEventsa := 0
	// processedEventsd := 0
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
			unprocessedDeallocs = map[src.ServiceID]Event{}
			deallocateBufferedEvents(orchestrator, eventsBuffer.DeallocEvents, results, endTime, unprocessedDeallocs)
			allocateBufferedEvents(orchestrator, eventsBuffer, results, endTime, unprocessedDeallocs)
			// unprocessedDeallocs = make([]Event, 0)
			for id, _ := range orchestrator.Domains {
				// if upgrade {
				// orchestrator.UpgradeServiceIfEnabledIntervalBased(id) // change this to only one domain.
				// }
				orchestrator.BasicNodeReclaim(id)
			}

			initTime = math.Floor(event.EventTime/interval) * interval
			endTime = interval + initTime
			eventsBuffer = EventsBuffer{
				DeallocEvents: make(DeallocEvents, 0),
				AllocEvents:   make(AllocEvents, 0),
			}

		}
	}

	if len(eventsBuffer.DeallocEvents) > 0 {
		unprocessedDeallocs = map[src.ServiceID]Event{}
		deallocateBufferedEvents(orchestrator, eventsBuffer.DeallocEvents, results, initTime, unprocessedDeallocs)
	}
	if len(eventsBuffer.AllocEvents) > 0 {
		allocateBufferedEvents(orchestrator, eventsBuffer, results, initTime, unprocessedDeallocs)
	}

	// log.Info("Total processed events:", processedEvents)
	// log.Info("number of events:", len(events))
	// log.Info("number of buffered deallocate events:", bufferedEventsd)
	// log.Info("number of buffered allocate events:", bufferedEventsa)
	// log.Info("number of processed deallocate events:", processedEventsd)
	// log.Info("number of processed allocate events:", processedEventsa)
	// log.Info("number of services in running services:", len(orchestrator.RunningServices))
	// log.Info("remaining services in the running services:", orchestrator.RunningServices)
	log.Info("does it reach here?")
	return results.Results, nil
}

func sortEvents(allocEvents []Event, allServices src.Services, unprocessedDealloc map[src.ServiceID]Event) ([]src.ServiceID, map[src.ServiceID]Event) {
	sortedEventIDs := make([]src.ServiceID, 0, len(allocEvents))
	sortedEvents := make(map[src.ServiceID]Event)
	for _, event := range allocEvents {
		if _, ok := unprocessedDealloc[event.EventID]; ok {
			continue
		}
		sortedEventIDs = append(sortedEventIDs, event.EventID)
		sortedEvents[event.EventID] = event
	}
	// sort.Slice(sortedEventIDs, func(i, j int) bool {
	// 	serviceimpi := allServices[sortedEvents[sortedEventIDs[i]].TargetServiceID].ImportanceFactor
	// 	serviceimpj := allServices[sortedEvents[sortedEventIDs[j]].TargetServiceID].ImportanceFactor
	// 	return float64(serviceimpi)*float64(sortedEvents[sortedEventIDs[i]].TotalUtil) > float64(serviceimpj)*float64(sortedEvents[sortedEventIDs[j]].TotalUtil)
	// })
	sort.Slice(sortedEventIDs, func(i, j int) bool {
		serviceimpi := allServices[sortedEvents[sortedEventIDs[i]].TargetServiceID].ImportanceFactor
		serviceimpj := allServices[sortedEvents[sortedEventIDs[j]].TargetServiceID].ImportanceFactor
		servicebwi := allServices[sortedEvents[sortedEventIDs[i]].TargetServiceID].StandardMode.BandwidthEdge
		servicebwj := allServices[sortedEvents[sortedEventIDs[j]].TargetServiceID].StandardMode.BandwidthEdge
		serviceci := allServices[sortedEvents[sortedEventIDs[i]].TargetServiceID].StandardMode.CpusEdge
		servicecj := allServices[sortedEvents[sortedEventIDs[j]].TargetServiceID].StandardMode.CpusEdge
		return float64(serviceimpi)/((servicebwi)*float64(serviceci)) > float64(serviceimpj)/((servicebwj)*float64(servicecj)) //+servicecj
	})
	return sortedEventIDs, sortedEvents

}

func computeQoSCost(qos float64, cost float64) float64 {
	return math.Round(qos / cost)
}

func (e *EventResult) appendMetrics(orchestrator *src.Orchestrator, time float64, duration float64) {

	e.Results.Qos = append(e.Results.Qos, float64(orchestrator.QoS))
	e.Results.Cost = append(e.Results.Cost, float64(orchestrator.Cost))
	if duration != 0 {
		e.Results.Durations = append(e.Results.Durations, duration)
	}
	e.Results.EventTime = append(e.Results.EventTime, time)
	e.Results.QosPerCost = append(e.Results.QosPerCost, float64(orchestrator.QoS)/float64(orchestrator.Cost))
}

func allocateBufferedEvents(orchestrator *src.Orchestrator, eventsBuffer EventsBuffer, results *EventResult, timed float64, unprocessedDealloc map[src.ServiceID]Event) {

	sortedEventIDs, sortedEvents := sortEvents(eventsBuffer.AllocEvents, orchestrator.AllServices, unprocessedDealloc)
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
		results.appendMetrics(orchestrator, aEvent.EventTime, float64(duration.Microseconds()/1000))
		// results.appendMetrics(orchestrator, timed, float64(duration.Microseconds()/1000))
	}
}

func deallocateBufferedEvents(orchestrator *src.Orchestrator, eventsBuffer []Event, results *EventResult, timed float64, unprocessedDeallocs map[src.ServiceID]Event) {
	log.Info("attempt to process deallocate events")
	log.Info("deallocation events:", eventsBuffer)
	for _, dEvent := range eventsBuffer {
		// if _, ok := processedIDs[dEvent.EventID]; !ok {
		log.Info("deallocate event:", dEvent)
		// }
		eventID := dEvent.EventID
		if _, ok := orchestrator.RunningServices[eventID]; ok {
			orchestrator.Deallocate(dEvent.TargetDomainID, dEvent.TargetServiceID, eventID)
			// orchestrator.BasicNodeReclaim(dEvent.TargetDomainID)
			// orchestrator.UpgradeServiceIfEnabledIntervalBased(dEvent.TargetDomainID) // change this to only one domain.
			orchestrator.CloudBasicNodeReclaim()
			log.Info("service exists in the first round")
			log.Info("/////////////////////")
		} else {
			unprocessedDeallocs[dEvent.EventID] = dEvent
			log.Info("service does not exist in first round. rejected?")
		}
		results.appendMetrics(orchestrator, dEvent.EventTime, 0)
		// results.appendMetrics(orchestrator, timed, 0)

	}
}

func BufferAllocateEvents(events []Event, interval float64, orchestrator *src.Orchestrator) (*cnfg.ResultContext, error) {
	initTime := 0.0
	endTime := interval + initTime
	eventsBuffer := EventsBuffer{
		DeallocEvents: make(DeallocEvents, 0),
		AllocEvents:   make(AllocEvents, 0),
	}
	// processedIDs := make(map[src.ServiceID]bool)
	// qosPerCost, qos, cost, durations, eventTime := []float64{}, []float64{}, []float64{}, []float64{}, []float64{}
	results := NewEventResult()
	unprocessedDeallocs := map[src.ServiceID]Event{}
	// rejectedEvents := []Event{}

	log.Info("Total events:", len(events), "Interval:", interval)
	log.Info("interval:", interval)
	log.Info("initTime:", initTime)
	log.Info("endTime:", endTime)
	eventIndex := 0
	bufferedEventsa := 0

	for eventIndex < len(events) {
		event := events[eventIndex]
		if event.EventType == "deallocate" {
			// upgrade := false
			if _, ok := orchestrator.RunningServices[event.EventID]; ok {
				// svc := *orchestrator.RunningServices[event.EventID]
				// if svc.AllocationMode == "Standard" {
				// 	upgrade = true
				// }
				orchestrator.Deallocate(event.TargetDomainID, event.TargetServiceID, event.EventID)
				// if upgrade {
				// 	orchestrator.UpgradeServiceIfEnabled(orchestrator.Config.UpgradeHeuristic, svc, event.TargetDomainID) // change this to only one domain.
				// }
				// orchestrator.BasicNodeReclaim(event.TargetDomainID)
				// orchestrator.NodeReclaimIfEnabled(event.TargetDomainID)
			} else {
				unprocessedDeallocs[event.EventID] = event
				log.Info("service does not exist in first round. rejected?")
			}
			results.appendMetrics(orchestrator, event.EventTime, 0)
		} else {
			if event.EventTime >= initTime && event.EventTime <= endTime {
				eventsBuffer.AllocEvents = append(eventsBuffer.AllocEvents, event)
				bufferedEventsa++
				continue
			}
		}
		log.Info("events are buffered")
		if event.EventTime > endTime {
			unprocessedDeallocs = map[src.ServiceID]Event{}
			allocateBufferedEvents(orchestrator, eventsBuffer, results, endTime, unprocessedDeallocs)
			for id, _ := range orchestrator.Domains {
				orchestrator.BasicNodeReclaim(id)
			}
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

	if len(eventsBuffer.AllocEvents) > 0 {
		allocateBufferedEvents(orchestrator, eventsBuffer, results, initTime, unprocessedDeallocs)
	}
	return results.Results, nil
}

func ProcessEventsBaseline(events []Event, orchestrator *src.Orchestrator) (*cnfg.ResultContext, error) {

	qosPerCost := make([]float64, 0)
	qos := make([]float64, 0)
	cost := make([]float64, 0)
	durations := make([]float64, 0)
	eventTime := make([]float64, 0)
	acceptance := make([]float64, 0)
	// rejected := make([]float64, 0)
	test := 0
	e := 0

	for _, event := range events {
		eventID := event.EventID
		log.Info("event:", event)
		log.Info("service:", orchestrator.AllServices[event.TargetServiceID])
		if event.EventType == "allocate" {
			startTime := time.Now()
			allocated, err := orchestrator.AllocateBaselineOld(event.TargetDomainID, event.TargetServiceID, eventID)
			duration := time.Since(startTime)
			log.Info("Allocate:", allocated, orchestrator.QoS, orchestrator.Cost)
			log.Info("Time to allocate:", duration)
			if err != nil {
				log.Info(err)
			}
			if allocated {
				log.Info("/////////////////////")
				log.Info("service is allocated ")
				e++
				log.Info("/////////////////////")
			} else {
				log.Info("/////////////////////")
				log.Info("service is rejected ")
				// delete(orchestrator.RunningServices, eventID)
				log.Info("/////////////////////")

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

			if _, ok := orchestrator.RunningServices[eventID]; ok {
				orchestrator.Deallocate(event.TargetDomainID, event.TargetServiceID, eventID)

				orchestrator.CloudBasicNodeReclaim()

			} else {
				log.Info("Service does not exist. rejected?")
			}
			qosPerCost = append(qosPerCost, math.Round(float64(orchestrator.QoS)*1000/float64(orchestrator.Cost))/1000)
			qos = append(qos, float64(orchestrator.QoS))
			cost = append(cost, float64(orchestrator.Cost))
			eventTime = append(eventTime, float64(event.EventTime))
		}
	}
	log.Info("QoS per Cost: ", qosPerCost)
	acceptance = append(acceptance, float64(e/len(events)))
	log.Info("Durations: ", durations)
	return &cnfg.ResultContext{
		QosPerCost: qosPerCost,
		Qos:        qos,
		Cost:       cost,
		Durations:  durations,
		EventTime:  eventTime,
		Acceptance: acceptance,
	}, nil
}
