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
	if err != nil {
		log.Fatal(err)
	}
	orchestrator := util.Initialise(config)
	results := new(cnfg.ResultContext)
	fmt.Println("../data/events/hightraffic/events_" + config.System.Addition + ".csv")
	events := util.LoadEventsFromCSV("../data/events/hightraffic/events_" + config.System.Addition + ".csv")
	if config.Orchestrator.Baseline {
		results, err = eve.ProcessEvents(events, orchestrator)
	} else {
		results, err = eve.BufferEvents(events, 3.0, orchestrator)
	}
	if err != nil {
		log.Fatal(err)
	}
	util.WriteResults(results, config)
}
