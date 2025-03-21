package util

import (
	"fmt"
	"os"

	cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"
	src "github.com/nasim-samimi/scaling-simulator/pkg/orchestrator"
	"github.com/sirupsen/logrus"
)

var logr = logrus.New()

// Automatically runs when the package is imported
func init() {
	logr.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logr.SetOutput(os.Stdout)       // Log to console
	logr.SetLevel(logrus.InfoLevel) // Set log level
}

func Initialise(config *cnfg.Config) *src.Orchestrator {

	CloudNodes, reservedCloudNodes := LoadCloudFromCSV("../data/cloud_b.csv")
	cloud := src.NewCloud(CloudNodes, reservedCloudNodes)
	// read domain csv files in domain folder
	svcs := LoadSVCFromCSV("../data/services/services0.csv")
	// nodeHeuristic, reallocHeuristic, partitionHeuristic := util.LoadHeuristicFromCSV("../data/heuristics.csv")
	nodeHeuristic := config.Orchestrator.NodeHeuristic
	partitionHeuristic := config.Orchestrator.PartitionHeuristic
	fmt.Println("configs", config)
	fmt.Println("initialising the orchestrator")
	fmt.Println("nodeHeuristic:", nodeHeuristic)
	fmt.Println("partitionHeuristic:", partitionHeuristic)
	fmt.Println("node size:", config.System.NodeSize)

	baseFolder := "../data/domainNodes" + config.System.NodeSize + "/" + string(partitionHeuristic) + "/MaxMax" //+ string(nodeHeuristic) /TODO: select MAXMAX for now
	domains := LoadDomains(baseFolder)

	// initialise the orchestrator
	orchestrator := src.NewOrchestrator(&config.Orchestrator, cloud, domains, svcs)
	fmt.Println("the addition from config file:", config.System.Addition)
	return orchestrator
}
