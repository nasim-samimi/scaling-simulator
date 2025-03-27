import pandas as pd
import os
import sys
import yaml
import subprocess
import itertools
# import sleep

PARTITIONING_H=['bestfit']

REALLOCATION_H=["LB","LC","LBC","LBCI","LBI","LCI","LRED","LI"]
REALLOCATION_INTRANODE_H=["HBCI","HBI","HCI","HB","HC","HBC","LB","LC","LBC","LBCI","LBI","LCI"]
REALLOCATION_INTRADOMAIN_H=["HBCI","HBI","HCI","HB","HC","HBC","LB","LC","LBC","LBCI","LBI","LCI"]
REALLOCATION_REDUCED_H=["LB","LC","LBC","LBCI","LBI","LCI","LRED","LI"]
# REALLOCATION_REDUCED_H=["LI"]
# REALLOCATION_REMOVED_H=["LB","LC","LBC","LBCI","LBI","LCI","LREM","LI"]
REALLOCATION_REMOVED_H=["LI"]
# REALLOCATION_H=["HBCI"]
NODE_SELECTION_H=["mmRB"]
# NODE_SELECTION_H=["MMRB","mmRB","MmRB","mMRB"]
# ADDITION=[0]
# ADDITION=[0,0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1]
ADDITION=[0,0.2,0.4,0.6,0.8,1.0,1.2,1.4,1.6,1.8,2.0]
edge_node_cost=1
cloud_node_cost=3
# thresholds=[80,100]
threshold=100
node_sizes = [8]
max_scaling_cores=[16,32,64,128]#,64,128]#[512,128,200,96,256,64,32,16]
# divide cores by node size.
max_scaling_threshold=[1,2,4]
# interval_length=[5,10,20,30,40,50,100]
interval_length=[5,10,20,30,40,50]

results_dir = "improved"
data_dir="data_interval"
filename="improved"

# Define fixed configuration
fixed_config = {
    "orchestrator": {
        "domain_node_threshold": 100,
        "cloud_node_cost": 3,
        "edge_node_cost": 1,
        "partition_heuristic": "bestfit",
        "node_heuristic": "MaxMax",
        "reallocation_heuristic": "HB",
        "upgrade_service":False,
        "node_reclaim":False,
        "intra_node_realloc":False,
        "intra_domain_realloc":False,
        "intra_node_reduced":False,
        "intra_node_removed":False,
        "interval_based":True,
        "max_scaling_threshold":20,
        "node_size": 8,
    },
    "system": {
        "init_node_size": 16,
        "scaled_node_size": 8,
        "node_size": 8
    }
}

# Define mutually exclusive options (only one can be enabled at a time)
exclusive_options = [
    # "upgrade_service",
    # "node_reclaim",
    # "intra_domain_realloc",
    # "intra_node_realloc",
    # "intra_node_reduced",
    # "intra_node_removed"
    # "interval_based"
]

# Generate all parameter combinations
def generate_param_combinations():
    for size in node_sizes:
        for i in interval_length:
            for th in max_scaling_cores:
                max_scaling_threshold=th/size
                for addition in ADDITION:
                    for p in PARTITIONING_H:
                        for n in NODE_SELECTION_H:
                                # If the option is a reallocation strategy, iterate through heuristics

                            # for heuristic in REALLOCATION_H:
                                # Create a new config with only one flag enabled
                            config=fixed_config.copy()
                            # config["orchestrator"] = {key: False for key in exclusive_options}  # Disable all first
                            # config["orchestrator"][exclusive_option] = True  # Enable only the current one
                            # config["orchestrator"][f"reallocation_heuristic"] = heuristic  # Assign heuristic
                            config["orchestrator"]["intraintra_node_realloc_heu"]="LB"
                            config["orchestrator"]["intra_node_reduced_heu"]="LRED"
                            config["orchestrator"]["domain_node_threshold"] = threshold

                            config["system"]["addition"] = addition
                            config["system"]["results_dir"] = f'{results_dir}/allOpts_interval_based/interval={i}/max_scaling_threshold={th}'
                            config["system"]["data_dir"] = data_dir
                            config["orchestrator"]["partition_heuristic"]=p
                            config["orchestrator"]["node_heuristic"]=n
                            config["orchestrator"]["edge_node_cost"]=edge_node_cost
                            config["orchestrator"]["cloud_node_cost"]=cloud_node_cost
                            config["orchestrator"]["node_size"]=size
                            config["orchestrator"]["max_scaling_threshold"]=max_scaling_threshold
                            config["orchestrator"]["upgrade_heuristic"]="HQ"
                            config["system"]["file_name"]=filename
                            config["system"]["interval_length"]=i

                            # Write to config.yaml
                            print(config)
                            with open("config.yaml", "w") as file:
                                yaml.dump(config, file, default_flow_style=False)

                            print(f"Generated config.yaml with: all options, Addition = {addition}, Results Dir = {results_dir}")

                            # Run the Go script
                            os.system(f'go run main.go > log.txt')

    print("All parameter combinations processed.")
    return

if __name__=='__main__':
    # run()                
    generate_param_combinations()
            
            