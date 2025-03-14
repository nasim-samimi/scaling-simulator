package util

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"
)

func WriteToCsv(filePath string, records []float64) error {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Unable to create output file %s, %v", filePath, err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		err := writer.Write([]string{strconv.FormatFloat(record, 'f', -1, 64)})
		if err != nil {
			log.Fatalf("Unable to write to CSV file %s, %v", filePath, err)
			return err
		}
	}
	return nil
}

func WriteResults(results *cnfg.ResultContext, config *cnfg.Config) error {
	// Construct the base directory paths
	// nodeSize := strconv.Itoa(int(config.System.NodeSize))
	nodeSize := config.System.NodeSize
	name := "/nodesize=" + nodeSize + "/addition=" + config.System.Addition + "/" + string(config.Orchestrator.NodeHeuristic) + "/" + string(config.Orchestrator.PartitionHeuristic) + "/"
	reallocName := ""
	// if orchestrator.Config.IntraDomainRealloc {
	// 	reallocName = string(orchestrator.Config.IntraDomainReallocHeu)
	// }
	// if orchestrator.Config.IntraNodeRealloc {
	// 	reallocName = string(orchestrator.Config.IntraNodeReallocHeu)
	// }
	// if orchestrator.Config.IntraNodeReduced {
	// 	reallocName = string(orchestrator.Config.IntraNodeReducedHeu)
	// }
	if config.Orchestrator.IntraDomainRealloc || config.Orchestrator.IntraNodeRealloc || config.Orchestrator.IntraNodeReduced || config.Orchestrator.IntraNodeRemoved {
		if config.Orchestrator.UpgradeService{
			reallocName=string(config.Orchestrator.ReallocationHeuristic)+"_"+string(config.Orchestrator.UpgradeHeuristic)
		}
		reallocName = string(config.Orchestrator.ReallocationHeuristic)
	}
	if reallocName ==""{
		if config.Orchestrator.UpgradeService{
			reallocName=string(config.Orchestrator.UpgradeHeuristic)
		}
	}
	

	if reallocName == "" {
		reallocName = "improved"
	}
	// Define output subdirectories
	subDirs := []string{"runtimes", "qosPerCost", "qos", "cost", "eventTime"}

	// Create directories if they do not exist
	baseFolder := config.System.ResultsDir
	for _, subDir := range subDirs {
		dirPath := filepath.Join("../experiments/results", baseFolder, subDir, name)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return fmt.Errorf("error creating directory %s: %v", dirPath, err)
		}
	}

	// Define file paths for each result type
	filePaths := map[string][]float64{
		filepath.Join("../experiments/results", baseFolder, "qosPerCost", name, reallocName+".csv"): results.QosPerCost,
		filepath.Join("../experiments/results", baseFolder, "runtimes", name, reallocName+".csv"):   results.Durations,
		filepath.Join("../experiments/results", baseFolder, "qos", name, reallocName+".csv"):        results.Qos,
		filepath.Join("../experiments/results", baseFolder, "cost", name, reallocName+".csv"):       results.Cost,
		filepath.Join("../experiments/results", baseFolder, "eventTime", name, reallocName+".csv"):  results.EventTime,
	}

	// Write results to CSV files
	for path, data := range filePaths {
		log.Println("Writing to file: ", path)
		if err := WriteToCsv(path, data); err != nil {
			return fmt.Errorf("error writing to file %s: %v", path, err)
		}
	}

	return nil
}
