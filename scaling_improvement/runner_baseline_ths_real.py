import pandas as pd
import os
import sys
import yaml
import subprocess
import itertools
# import sleep

PARTITIONING_H=['bestfit']

# REALLOCATION_H=["HBCI"]
NODE_SELECTION_H=["None"]
# NODE_SELECTION_H=["MinMin"]
# ADDITION=[0,0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1]
ADDITION=[0.1,0.25,0.5,0.75,1.0]
# ADDITION=[0.7,0.8,0.9,1]

edge_node_cost=1
cloud_node_cost=3


results_dir = "baseline_real"

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
        "baseline":True,
        "node_size": 8,
    },
    "system": {
        "init_node_size": 16,
        "scaled_node_size": 8,
        "node_size": 8
    }
}

# Define mutually exclusive options (only one can be enabled at a time)
# node_sizes = [8,12,16,20,24,28,32]
node_sizes = [8,16]
max_scaling_cores=[16,32,64,128]#[128,96,64,32,16,512,256,200]
data_dir="data_real"
# Generate all parameter combinations
def generate_param_combinations():
    for size in node_sizes:
        for th in max_scaling_cores:
            max_scaling_threshold=th/size
            for addition in ADDITION:
                for p in PARTITIONING_H:
                    for n in NODE_SELECTION_H:
                            # If the option is a reallocation strategy, iterate through heuristics    
                        # Create a new config with only one flag enabled
                        config=fixed_config.copy()

                        config["system"]["addition"] = addition
                        config["system"]["results_dir"] = f'{results_dir}/{th}'
                        config["system"]["data_dir"] = data_dir
                        config["orchestrator"]["partition_heuristic"]=p
                        config["orchestrator"]["node_heuristic"]=n
                        config["orchestrator"]["edge_node_cost"]=edge_node_cost
                        config["orchestrator"]["cloud_node_cost"]=cloud_node_cost
                        config["orchestrator"]["node_size"]=size
                        config["orchestrator"]["max_scaling_threshold"]=max_scaling_threshold


                        # Write to config.yaml
                        print(config)
                        with open("config.yaml", "w") as file:
                            yaml.dump(config, file, default_flow_style=False)


                        # Run the Go script
                        os.system(f'go run main.go > log.txt')
                        # return


    print("All parameter combinations processed.")
    return



if __name__=='__main__':
    # run()                
    generate_param_combinations()
            
            