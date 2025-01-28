package util

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	src "github.com/nasim-samimi/scaling-simulator/pkg/orchestrator"
)

func WriteToCsv(filePath string, records []float64) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Unable to create output file %s, %v", filePath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		err := writer.Write([]string{strconv.FormatFloat(record, 'f', -1, 64)})
		if err != nil {
			log.Fatalf("Unable to write to CSV file %s, %v", filePath, err)
		}
	}
}

func WriteResults(cost []float64, qos []float64, qosPerCost []float64, durations []float64, orchestrator *src.Orchestrator, addition string) error {
	// write the results to csv files

	name := "addition=" + addition + "/" + string(orchestrator.Config.NodeHeuristic) + "/" + string(orchestrator.Config.PartitionHeuristic) + "/"
	name2 := string(orchestrator.Config.ReallocationHeuristic)
	//first check if directory exists
	if _, err := os.Stat("../experiments/results/improved/runtimes/" + name); os.IsNotExist(err) {
		// create directory
		os.MkdirAll("../experiments/results/improved/runtimes/"+name, os.ModePerm)
	}
	if _, err := os.Stat("../experiments/results/improved/qosPerCost/" + name); os.IsNotExist(err) {
		// create directory
		os.MkdirAll("../experiments/results/improved/qosPerCost/"+name, os.ModePerm)
	}
	if _, err := os.Stat("../experiments/results/improved/qos/" + name); os.IsNotExist(err) {
		// create directory
		os.MkdirAll("../experiments/results/improved/qos/"+name, os.ModePerm)
	}
	if _, err := os.Stat("../experiments/results/improved/cost/" + name); os.IsNotExist(err) {
		// create directory
		os.MkdirAll("../experiments/results/improved/cost/"+name, os.ModePerm)
	}
	WriteToCsv("../experiments/results/improved/qosPerCost/"+name+name2+".csv", qosPerCost)
	WriteToCsv("../experiments/results/improved/runtimes/"+name+name2+".csv", durations)
	WriteToCsv("../experiments/results/improved/qos/"+name+name2+".csv", qos)
	WriteToCsv("../experiments/results/improved/cost/"+name+name2+".csv", cost)
	return nil
}
