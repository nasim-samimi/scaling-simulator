package util

import (
	"fmt"

	cnfg "github.com/nasim-samimi/scaling-simulator/pkg/config"
	src "github.com/nasim-samimi/scaling-simulator/pkg/orchestrator"
)

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

	baseFolder := "../data/domainNodes" + config.System.NodeSize + "/" + string(partitionHeuristic) + "/" + string(nodeHeuristic)
	domains := LoadDomains(baseFolder)

	// initialise the orchestrator
	orchestrator := src.NewOrchestrator(&config.Orchestrator, cloud, domains, svcs)
	fmt.Println("the addition from config file:", config.System.Addition)
	return orchestrator
}
