import pandas as pd
import os
import sys
import yaml
import subprocess
import itertools
# import sleep

PARTITIONING_H=['bestfit','worstfit']

REALLOCATION_H=["HBCI","HBI","HCI","HB","HC","HBC","LB","LC","LBC"]
# REALLOCATION_H=["HBCI"]
NODE_SELECTION_H=["MinMin","MaxMax"]
ADDITION=[0,0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1]
upgrade_service_options = [True, False]
node_reclaim_options = [True, False]
intra_node_realloc_options = [True, False]
edge_node_cost=1
cloud_node_cost=3
thresholds=[80,100]


results_dir = "improved"

# Define fixed configuration
fixed_config = {
    "orchestrator": {
        "domain_node_threshold": 100,
        "cloud_node_cost": 3,
        "edge_node_cost": 1,
        "partition_heuristic": "bestfit",
        "node_heuristic": "MaxMax",
        "reallocation_heuristic": "HB",
        "domain_node_threshold": 100,
        "upgrade_service":False,
        "node_reclaim":False,
    },
    "system": {
        "init_node_size": 16,
        "scaled_node_size": 8
    }
}

# Define mutually exclusive options (only one can be enabled at a time)
exclusive_options = [
    
    "intra_node_realloc",
    "intra_domain_realloc",
    "intra_node_reduced"
]

# Generate all parameter combinations
def generate_param_combinations():
    for exclusive_option in exclusive_options:
        for addition in ADDITION:
            for p in PARTITIONING_H:
                for n in NODE_SELECTION_H:
                        # If the option is a reallocation strategy, iterate through heuristics
                    if "realloc" in exclusive_option or "intra_node_reduced" in exclusive_option:
                        for threshold in thresholds:
                            for heuristic in REALLOCATION_H:
                                # Create a new config with only one flag enabled
                                config=fixed_config.copy()
                                config["orchestrator"] = {key: False for key in exclusive_options}  # Disable all first
                                config["orchestrator"][exclusive_option] = True  # Enable only the current one
                                config["orchestrator"][f"reallocation_heuristic"] = heuristic  # Assign heuristic
                                config["orchestrator"]["domain_node_threshold"] = threshold

                                config["system"]["addition"] = addition
                                config["system"]["results_dir"] = f'{results_dir}/with_{exclusive_option}_threshold_{threshold}'
                                config["orchestrator"]["partition_heuristic"]=p
                                config["orchestrator"]["node_heuristic"]=n
                                config["orchestrator"]["edge_node_cost"]=edge_node_cost
                                config["orchestrator"]["cloud_node_cost"]=cloud_node_cost

                                # Write to config.yaml
                                print(config)
                                with open("config.yaml", "w") as file:
                                    yaml.dump(config, file, default_flow_style=False)

                                print(f"Generated config.yaml with: {exclusive_option} = True, Heuristic = {heuristic}, Addition = {addition}, Results Dir = {results_dir}")

                                # Run the Go script
                                os.system(f'go run main.go > log.txt')
                    else:
                        # Non-reallocation cases (directly set the option to True)
                        config=fixed_config.copy()
                        
                        config["orchestrator"] = {key: False for key in exclusive_options}  # Disable all first
                        config["orchestrator"][exclusive_option] = True  # Enable only the current one

                        config["system"]["addition"] = addition
                        config["system"]["results_dir"] = f'{results_dir}/with_{exclusive_option}'
                        config["orchestrator"]["partition_heuristic"]=p
                        config["orchestrator"]["node_heuristic"]=n  
                        config["orchestrator"]["edge_node_cost"]=edge_node_cost
                        config["orchestrator"]["cloud_node_cost"]=cloud_node_cost

                        # Write to config.yaml
                        print(config)
                        with open("config.yaml", "w") as file:
                            yaml.dump(config, file, default_flow_style=False)

                        print(f"Generated config.yaml with: {exclusive_option} = True, Addition = {addition}, Results Dir = {results_dir}")

                        # Run the Go script
                        os.system(f'go run main.go > log.txt')

    print("All parameter combinations processed.")
    return

fixed_config_allOpts = {
    "orchestrator": {
        "domain_node_threshold": 100,
        "cloud_node_cost": 3,
        "edge_node_cost": 1,
        "partition_heuristic": "bestfit",
        "node_heuristic": "MaxMax",
        "reallocation_heuristic": "HB",
        "upgrade_service":True,
        "node_reclaim":True,
        "intra_node_realloc":True,
        "intra_domain_realloc":True,
        "intra_node_reduced":True
    },
    "system": {
        "init_node_size": 16,
        "scaled_node_size": 8,
        "addition": 0,
    }
}

# these are to be changed 
INTRANODE_REALLOCATION_H=["HBCI","HBI","HCI","HB","HC","HBC","LB","LC","LBC"]
INTRADOMAIN_REALLOCATION_H=["HBCI","HBI","HCI","HB","HC","HBC","LB","LC","LBC"]
INTRANODE_REDUCED_REALLOCATION_H=["HBCI","HBI","HCI","HB","HC","HBC","LB","LC","LBC"]
def generate_param_full_option():
    for addition in ADDITION:
        for p in PARTITIONING_H:
            for n in NODE_SELECTION_H:
                    # If the option is a reallocation strategy, iterate through heuristics
                for intranodeheuristic in INTRANODE_REALLOCATION_H:
                    for intradomainheuristic in INTRADOMAIN_REALLOCATION_H:
                        for intranodereducedheuristic in INTRANODE_REDUCED_REALLOCATION_H:
                    # Create a new config with only one flag enabled
                            config=fixed_config.copy()
                            # config["orchestrator"][f"reallocation_heuristic"] = heuristic  # Assign heuristic
                            config["orchestrator"]["intra_node_realloc_heu"]=intranodeheuristic
                            config["orchestrator"]["intra_domain_realloc_heu"]=intradomainheuristic
                            config["orchestrator"]["intra_node_reduced_heu"]=intranodereducedheuristic

                            config["system"]["addition"] = addition
                            config["system"]["results_dir"] = f'{results_dir}/allOpts'
                            config["orchestrator"]["partition_heuristic"]=p
                            config["orchestrator"]["node_heuristic"]=n
                            config["orchestrator"]["edge_node_cost"]=edge_node_cost
                            config["orchestrator"]["cloud_node_cost"]=cloud_node_cost

                            # Write to config.yaml
                            print(config)
                            with open("config.yaml", "w") as file:
                                yaml.dump(config, file, default_flow_style=False)

                            # Run the Go script
                            os.system(f'go run main.go > log.txt')


    return
if __name__=='__main__':
    # run()                
    generate_param_combinations()
            
            