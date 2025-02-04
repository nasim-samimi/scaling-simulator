package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	"time"

	cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"
	src "github.com/nasim-samimi/scaling-simulator/pkg/orchestrator"
	util "github.com/nasim-samimi/scaling-simulator/pkg/util"
)

func processEvents(orchestrator *src.Orchestrator, addition string) (*cnfg.ResultContext, error) {

	events := util.LoadEventsFromCSV("../data/events/events_" + addition + ".csv")
	qosPerCost := make([]float64, 0)
	qos := make([]float64, 0)
	cost := make([]float64, 0)
	durations := make([]float64, 0)
	eventTime := make([]float64, 0)
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
				qosPerCost = append(qosPerCost, math.Round(float64(orchestrator.QoS)*1000/float64(orchestrator.Cost))/1000)
				qos = append(qos, float64(orchestrator.QoS))
				cost = append(cost, float64(orchestrator.Cost))
				durations = append(durations, float64(duration.Microseconds())/1000)
				eventTime = append(eventTime, float64(event.EventTime))
				fmt.Println("/////////////////////")
				fmt.Println("service is allocated ")
				fmt.Println("/////////////////////")
			} else {
				fmt.Println("/////////////////////")
				fmt.Println("service is rejected ")
				fmt.Println("/////////////////////")
			}
			test++
			// if test == 50 {
			// 	break
			// }
		}
		if event.EventType == "deallocate" {
			fmt.Println("/////////////////////")
			fmt.Println("Deallocate")
			orchestrator.Deallocate(event.TargetDomainID, event.TargetServiceID, eventID)
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

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the configuration YAML file")
	flag.Parse()
	config, err := cnfg.LoadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	orchestrator := util.Initialise(config)
	results, err := processEvents(orchestrator, config.System.Addition)
	if err != nil {
		log.Fatal(err)
	}
	util.WriteResults(results, config)
}
