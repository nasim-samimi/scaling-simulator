package util

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	src "github.com/nasim-samimi/scaling-simulator/pkg/orchestrator"
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

func WriteResults(cost []float64, qos []float64, qosPerCost []float64, durations []float64, orchestrator *src.Orchestrator, addition string, baseFolder string) error {
	// Construct the base directory paths
	name := "addition=" + addition + "/" + string(orchestrator.Config.NodeHeuristic) + "/" + string(orchestrator.Config.PartitionHeuristic) + "/"
	name2 := string(orchestrator.Config.ReallocationHeuristic)

	// Define output subdirectories
	subDirs := []string{"runtimes", "qosPerCost", "qos", "cost"}

	// Create directories if they do not exist
	for _, subDir := range subDirs {
		dirPath := filepath.Join("../experiments/results", baseFolder, subDir, name)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return fmt.Errorf("error creating directory %s: %v", dirPath, err)
		}
	}

	// Define file paths for each result type
	filePaths := map[string][]float64{
		filepath.Join("../experiments/results", baseFolder, "qosPerCost", name, name2+".csv"): qosPerCost,
		filepath.Join("../experiments/results", baseFolder, "runtimes", name, name2+".csv"):   durations,
		filepath.Join("../experiments/results", baseFolder, "qos", name, name2+".csv"):        qos,
		filepath.Join("../experiments/results", baseFolder, "cost", name, name2+".csv"):       cost,
	}

	// Write results to CSV files
	for path, data := range filePaths {
		if err := WriteToCsv(path, data); err != nil {
			return fmt.Errorf("error writing to file %s: %v", path, err)
		}
	}

	return nil
}
