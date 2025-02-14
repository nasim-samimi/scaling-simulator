package events

import (
	"fmt"
	"math"
	"time"

	cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"
	src "github.com/nasim-samimi/scaling-simulator/pkg/orchestrator"
)

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
		fmt.Println("event:", event)
		fmt.Println("service:", orchestrator.AllServices[event.TargetServiceID])
		if event.EventType == "allocate" {
			startTime := time.Now()
			allocated, err := orchestrator.Allocate(event.TargetDomainID, event.TargetServiceID, eventID)
			duration := time.Since(startTime)
			fmt.Println("Allocate:", allocated, orchestrator.QoS, orchestrator.Cost)
			fmt.Println("Time to allocate:", duration)
			if err != nil {
				fmt.Println(err)
			}
			if allocated {
				fmt.Println("/////////////////////")
				fmt.Println("service is allocated ")
				fmt.Println("/////////////////////")
			} else {
				fmt.Println("/////////////////////")
				fmt.Println("service is rejected ")
				// delete(orchestrator.RunningServices, eventID)
				fmt.Println("/////////////////////")

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
			fmt.Println("/////////////////////")
			fmt.Println("Deallocate")
			if _, ok := orchestrator.RunningServices[eventID]; ok {
				orchestrator.Deallocate(event.TargetDomainID, event.TargetServiceID, eventID)
			}
			orchestrator.UpgradeServiceIfEnabled()
			orchestrator.NodeReclaimIfEnabled(event.TargetDomainID)
			qosPerCost = append(qosPerCost, math.Round(float64(orchestrator.QoS)*1000/float64(orchestrator.Cost))/1000)
			qos = append(qos, float64(orchestrator.QoS))
			cost = append(cost, float64(orchestrator.Cost))
			eventTime = append(eventTime, float64(event.EventTime))
			fmt.Println("/////////////////////")
		}
	}
	fmt.Println("QoS per Cost: ", qosPerCost)

	fmt.Println("Durations: ", durations)
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
	processedIDs := make(map[src.ServiceID]bool)
	qosPerCost := make([]float64, 0)
	qos := make([]float64, 0)
	cost := make([]float64, 0)
	durations := make([]float64, 0)
	eventTime := make([]float64, 0)
	unprocessedDeallocs := make([]Event, 0)
	fmt.Println("length of events:", len(events))
	fmt.Println("interval:", interval)
	fmt.Println("initTime:", initTime)
	fmt.Println("endTime:", endTime)
	for _, event := range events {
		if event.EventTime >= initTime && event.EventTime <= endTime {
			if event.EventType == "allocate" {
				eventsBuffer.AllocEvents = append(eventsBuffer.AllocEvents, event)
			}
			if event.EventType == "deallocate" {
				eventsBuffer.DeallocEvents = append(eventsBuffer.DeallocEvents, event)
			}
			continue
		}
		fmt.Println("events are buffered")
		if event.EventTime > endTime {
			fmt.Println("attempt to process deallocate events")
			unprocessedDeallocs = nil
			for e, dEvent := range eventsBuffer.DeallocEvents {
				if _, ok := processedIDs[dEvent.EventID]; !ok {
					unprocessedDeallocs = append(unprocessedDeallocs, eventsBuffer.DeallocEvents[e])
					continue
				}
				eventID := dEvent.EventID
				if _, ok := orchestrator.RunningServices[eventID]; ok {
					orchestrator.Deallocate(dEvent.TargetDomainID, dEvent.TargetServiceID, eventID)
					qosPerCost = append(qosPerCost, math.Round(float64(orchestrator.QoS)*1000/float64(orchestrator.Cost))/1000)
					qos = append(qos, float64(orchestrator.QoS))
					cost = append(cost, float64(orchestrator.Cost))
					// eventTime = append(eventTime, float64(dEvent.EventTime))
					eventTime = append(eventTime, float64(endTime))
					fmt.Println("/////////////////////")
				}

			}
			// sort events based on qos
			for _, aEvent := range eventsBuffer.AllocEvents {
				eventID := aEvent.EventID
				startTime := time.Now()
				allocated, err := orchestrator.Allocate(aEvent.TargetDomainID, aEvent.TargetServiceID, eventID)
				duration := time.Since(startTime)
				fmt.Println("Allocate:", allocated, orchestrator.QoS, orchestrator.Cost)
				fmt.Println("Time to allocate:", duration)
				if err != nil {
					fmt.Println(err)
				}
				if allocated {
					fmt.Println("/////////////////////")
					fmt.Println("service is allocated ")
					fmt.Println("/////////////////////")
				} else {
					fmt.Println("/////////////////////")
					fmt.Println("service is rejected ")
					// delete(orchestrator.RunningServices, eventID)
					fmt.Println("/////////////////////")

				}
				qosPerCost = append(qosPerCost, math.Round(float64(orchestrator.QoS)*1000/float64(orchestrator.Cost))/1000)
				qos = append(qos, float64(orchestrator.QoS))
				cost = append(cost, float64(orchestrator.Cost))
				durations = append(durations, float64(duration.Microseconds())/1000)
				// eventTime = append(eventTime, float64(aEvent.EventTime))
				eventTime = append(eventTime, float64(endTime))
				processedIDs[eventID] = true
			}
			if len(unprocessedDeallocs) > 0 {
				for _, dEvent := range unprocessedDeallocs {
					eventID := dEvent.EventID
					if _, ok := orchestrator.RunningServices[eventID]; ok {
						orchestrator.Deallocate(dEvent.TargetDomainID, dEvent.TargetServiceID, eventID)
						qosPerCost = append(qosPerCost, math.Round(float64(orchestrator.QoS)*1000/float64(orchestrator.Cost))/1000)
						qos = append(qos, float64(orchestrator.QoS))
						cost = append(cost, float64(orchestrator.Cost))
						// eventTime = append(eventTime, float64(dEvent.EventTime))
						eventTime = append(eventTime, float64(endTime))
						fmt.Println("/////////////////////")
					}

				}
				unprocessedDeallocs = make([]Event, 0)
			}
			initTime = endTime
			endTime = interval + initTime
			eventsBuffer = EventsBuffer{
				DeallocEvents: make(DeallocEvents, 0),
				AllocEvents:   make(AllocEvents, 0),
			}
			fmt.Println("initTime:", initTime)
			fmt.Println("endTime:", endTime)

		}
		if event.EventTime <= endTime {
			if event.EventType == "allocate" {
				eventsBuffer.AllocEvents = append(eventsBuffer.AllocEvents, event)
			}
			if event.EventType == "deallocate" {
				eventsBuffer.DeallocEvents = append(eventsBuffer.DeallocEvents, event)
			}
		}
	}

	return &cnfg.ResultContext{
		QosPerCost: qosPerCost,
		Qos:        qos,
		Cost:       cost,
		Durations:  durations,
		EventTime:  eventTime,
	}, nil
}
