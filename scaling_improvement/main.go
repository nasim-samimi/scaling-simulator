package main

import (
	"flag"
	"fmt"
	"log"

	cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"
	eve "github.com/nasim-samimi/scaling-simulator/pkg/events"
	util "github.com/nasim-samimi/scaling-simulator/pkg/util"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the configuration YAML file")
	flag.Parse()
	config, err := cnfg.LoadConfig(*configPath)
	fmt.Println("../" + config.System.DataDir + "/events/hightraffic/events_" + config.System.Addition + ".csv")
	if err != nil {
		log.Fatal(err)
	}
	orchestrator := util.Initialise(config)
	results := new(cnfg.ResultContext)
	fmt.Println("../" + config.System.DataDir + config.System.Addition + ".csv")
	events := util.LoadEventsFromCSV("../" + config.System.DataDir + "/events/hightraffic/events_" + config.System.Addition + ".csv")
	if config.Orchestrator.Baseline {
		results, err = eve.ProcessEventsBaseline(events, orchestrator)
	} else {
		if config.Orchestrator.IntervalBased {
			events := util.LoadEventsFromCSV("../" + config.System.DataDir + "/events/hightraffic/events_" + config.System.Addition + ".csv")
			results, err = eve.BufferEvents(events, config.System.IntervalLength, orchestrator) // 20ms was working before
			// results, err = eve.BufferAllocateEvents(events, 5.0, orchestrator)
		} else {
			results, err = eve.ProcessEvents(events, orchestrator)
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	util.WriteResults(results, config)
}
