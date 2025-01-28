package config

import (
	"fmt"
	"os"

	"sigs.k8s.io/yaml"
)

type Config struct {
	Addition              string  `yaml:"addition"`
	UpgradeService        bool    `yaml:"upgrade_service"`
	NodeReclaim           bool    `yaml:"node_reclaim"`
	IntraNodeRealloc      bool    `yaml:"intra_node_realloc"`
	IntraNodeReallocHeu   string  `yaml:"intra_node_realloc_heu"`
	IntraDomainRealloc    bool    `yaml:"intra_domain_realloc"`
	IntraDomainReallocHeu string  `yaml:"intra_domain_realloc_heu"`
	IntraNodeReduced      bool    `yaml:"intra_node_reduced"`
	IntraNodeReducedHeu   string  `yaml:"intra_node_reduced_heu"`
	DomainNodeThreshold   float64 `yaml:"domain_node_threshold"`
	EdgeReduced           bool    `yaml:"edge_reduced"`
	CloudNodeCost         float64 `yaml:"cloud_node_cost"`
	EdgeNodeCost          float64 `yaml:"edge_node_cost"`
	InitNodeSize          uint64  `yaml:"init_node_size"`
	ScaledNodeSize        uint64  `yaml:"scaled_node_size"`
	PartitionHeuristic    string  `yaml:"partition_heuristic"`
	NodeHeuristic         string  `yaml:"node_selection_heuristic"`
	ReallocationHeuristic string  `yaml:"reallocation_heuristic"`
}

func LoadConfig(filePath string) (*Config, error) {
	// Read the YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %v", err)
	}

	// Parse the YAML into the Config struct
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %v", err)
	}

	return &config, nil
}
