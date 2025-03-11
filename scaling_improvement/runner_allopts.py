import pandas as pd
import os
import sys
import yaml
import subprocess
import itertools
# import sleep

PARTITIONING_H=['bestfit','worstfit']

REALLOCATION_H=["LB","LC","LBC","LBCI","LBI","LCI","LRED","LI"]
REALLOCATION_INTRANODE_H=["HBCI","HBI","HCI","HB","HC","HBC","LB","LC","LBC","LBCI","LBI","LCI"]
REALLOCATION_INTRADOMAIN_H=["HBCI","HBI","HCI","HB","HC","HBC","LB","LC","LBC","LBCI","LBI","LCI"]
REALLOCATION_REDUCED_H=["LB","LC","LBC","LBCI","LBI","LCI","LRED","LI"]
# REALLOCATION_REDUCED_H=["LI"]
# REALLOCATION_REMOVED_H=["LB","LC","LBC","LBCI","LBI","LCI","LREM","LI"]
REALLOCATION_REMOVED_H=["LI"]
# REALLOCATION_H=["HBCI"]
# NODE_SELECTION_H=["MinMin","MaxMax"]
NODE_SELECTION_H=["MMRB","mmRB","MmRB","mMRB"]
# ADDITION=[0]
ADDITION=[0,0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1]
edge_node_cost=1
cloud_node_cost=3
# thresholds=[80,100]
threshold=100
node_sizes = [8,16]
max_scaling_cores=[16,32,64]
# divide cores by node size.
max_scaling_threshold=[1,2,4]

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
        "upgrade_service":False,
        "node_reclaim":False,
        "intra_node_realloc":True,
        "intra_domain_realloc":False,
        "intra_node_reduced":True,
        "intra_node_removed":False,
        "interval_based":True,
        "max_scaling_threshold":20,
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
    "interval_based"
]

# Generate all parameter combinations
def generate_param_combinations():
    for size in node_sizes:
        for th in max_scaling_cores:
            max_scaling_threshold=th/size
            for addition in ADDITION:
                for p in PARTITIONING_H:
                    for n in NODE_SELECTION_H:
                            # If the option is a reallocation strategy, iterate through heuristics

                        for heuristic in REALLOCATION_H:
                            # Create a new config with only one flag enabled
                            config=fixed_config.copy()
                            # config["orchestrator"] = {key: False for key in exclusive_options}  # Disable all first
                            # config["orchestrator"][exclusive_option] = True  # Enable only the current one
                            config["orchestrator"][f"reallocation_heuristic"] = heuristic  # Assign heuristic
                            config["orchestrator"]["domain_node_threshold"] = threshold

                            config["system"]["addition"] = addition
                            config["system"]["results_dir"] = f'{results_dir}/allOpts/{th}=max_scaling_threshold'
                            config["orchestrator"]["partition_heuristic"]=p
                            config["orchestrator"]["node_heuristic"]=n
                            config["orchestrator"]["edge_node_cost"]=edge_node_cost
                            config["orchestrator"]["cloud_node_cost"]=cloud_node_cost
                            config["system"]["node_size"]=size
                            config["orchestrator"]["max_scaling_threshold"]=max_scaling_threshold

                            # Write to config.yaml
                            print(config)
                            with open("config.yaml", "w") as file:
                                yaml.dump(config, file, default_flow_style=False)

                            print(f"Generated config.yaml with: all options, Heuristic = {heuristic}, Addition = {addition}, Results Dir = {results_dir}")

                            # Run the Go script
                            os.system(f'go run main.go > log.txt')

    print("All parameter combinations processed.")
    return

# fixed_config_allOpts = {
#     "orchestrator": {
#         "domain_node_threshold": 100,
#         "cloud_node_cost": 3,
#         "edge_node_cost": 1,
#         "partition_heuristic": "bestfit",
#         "node_heuristic": "MaxMax",
#         "reallocation_heuristic": "HB",
#         "upgrade_service":True,
#         "node_reclaim":True,
#         "intra_node_realloc":True,
#         "intra_domain_realloc":True,
#         "intra_node_reduced":True
#     },
#     "system": {
#         "init_node_size": 16,
#         "scaled_node_size": 8,
#         "addition": 0,
#     }
# }

# # these are to be changed 
# INTRANODE_REALLOCATION_H=["HBCI","HBI","HCI","HB","HC","HBC","LB","LC","LBC"]
# INTRADOMAIN_REALLOCATION_H=["HBCI","HBI","HCI","HB","HC","HBC","LB","LC","LBC"]
# INTRANODE_REDUCED_REALLOCATION_H=["HBCI","HBI","HCI","HB","HC","HBC","LB","LC","LBC"]
# def generate_param_full_option():
#     for addition in ADDITION:
#         for p in PARTITIONING_H:
#             for n in NODE_SELECTION_H:
#                     # If the option is a reallocation strategy, iterate through heuristics
#                 for intranodeheuristic in INTRANODE_REALLOCATION_H:
#                     for intradomainheuristic in INTRADOMAIN_REALLOCATION_H:
#                         for intranodereducedheuristic in INTRANODE_REDUCED_REALLOCATION_H:
#                     # Create a new config with only one flag enabled
#                             config=fixed_config.copy()
#                             # config["orchestrator"][f"reallocation_heuristic"] = heuristic  # Assign heuristic
#                             config["orchestrator"]["intra_node_realloc_heu"]=intranodeheuristic
#                             config["orchestrator"]["intra_domain_realloc_heu"]=intradomainheuristic
#                             config["orchestrator"]["intra_node_reduced_heu"]=intranodereducedheuristic

#                             config["system"]["addition"] = addition
#                             config["system"]["results_dir"] = f'{results_dir}/allOpts'
#                             config["orchestrator"]["partition_heuristic"]=p
#                             config["orchestrator"]["node_heuristic"]=n
#                             config["orchestrator"]["edge_node_cost"]=edge_node_cost
#                             config["orchestrator"]["cloud_node_cost"]=cloud_node_cost

#                             # Write to config.yaml
#                             print(config)
#                             with open("config.yaml", "w") as file:
#                                 yaml.dump(config, file, default_flow_style=False)

#                             # Run the Go script
#                             os.system(f'go run main.go > log.txt')


#     return
if __name__=='__main__':
    # run()                
    generate_param_combinations()
            
            