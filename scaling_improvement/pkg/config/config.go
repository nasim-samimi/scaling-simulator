package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type OrchestratorConfig struct {
	UpgradeService        bool      `yaml:"upgrade_service"`
	NodeReclaim           bool      `yaml:"node_reclaim"`
	IntraNodeRealloc      bool      `yaml:"intra_node_realloc"`
	IntraNodeReallocHeu   Heuristic `yaml:"intra_node_realloc_heu"`
	IntraDomainRealloc    bool      `yaml:"intra_domain_realloc"`
	IntraDomainReallocHeu Heuristic `yaml:"intra_domain_realloc_heu"`
	IntraNodeReduced      bool      `yaml:"intra_node_reduced"`
	IntraNodeReducedHeu   Heuristic `yaml:"intra_node_reduced_heu"`
	IntraNodeRemoved      bool      `yaml:"intra_node_removed"`
	DomainNodeThreshold   float64   `yaml:"domain_node_threshold"`
	EdgeReduced           bool      `yaml:"edge_reduced"`
	CloudNodeCost         Cost      `yaml:"cloud_node_cost"`
	EdgeNodeCost          Cost      `yaml:"edge_node_cost"`
	IntervalBased         bool      `yaml:"interval_based"`

	PartitionHeuristic    Heuristic `yaml:"partition_heuristic"`
	NodeHeuristic         Heuristic `yaml:"node_heuristic"`
	ReallocationHeuristic Heuristic `yaml:"reallocation_heuristic"`
	Baseline              bool      `yaml:"baseline"`
}
type SystemConfig struct {
	InitNodeSize   uint64 `yaml:"init_node_size"`
	ScaledNodeSize uint64 `yaml:"scaled_node_size"`
	Addition       string `yaml:"addition"`
	ResultsDir     string `yaml:"results_dir"`
	NodeSize       string `yaml:"node_size"`
}

type Config struct {
	Orchestrator OrchestratorConfig `yaml:"orchestrator"`
	System       SystemConfig       `yaml:"system"`
}

func LoadConfig(filePath string) (*Config, error) {
	// Read the YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %v", err)
	}
	// Parse the YAML into the Config struct
	var config Config
	if err := yaml.Unmarshal([]byte(data), &config); err != nil {
		fmt.Printf("Error unmarshaling YAML: %v\n", err)
		return nil, err
	}
	fmt.Println("raw data:", config)

	return &config, nil
}
