import pandas as pd
import os
import sys
from results import *
# import sleep

PARTITIONING_H=['bestfit','worstfit']

# REALLOCATION_H=["HBCI","HBI","HCI","HB","HC","HBC","LB","LC","LBC"]
REALLOCATION_H=["HBCI"]
NODE_SELECTION_H=["MinMin","MaxMax"]
flags=['upgrade_service','node_reclaim','intra_node_realloc','intra_domain_realloc','intra_node_reduced'] #did we have more flags?
node_thresholds=[80,100]

if __name__=='__main__':
    # run()
    for f in flags: 
        if 'intra' in f:
            for t in node_thresholds:
                dir=f'improved/with_{f}_threshold_{t}/'
                robustness(dir1=dir,metric='cost',flags=f'with_{f}_threshold_{t}')            
                robustness(dir1=dir,metric='qos',flags=f'with_{f}_threshold_{t}')
                robustness(dir1=dir,metric='qosPerCost',flags=f'with_{f}_threshold_{t}')
        else:
            dir=f'improved/with_{f}/'   
            robustness(dir1=dir,metric='cost',flags=f'with_{f}')            
            robustness(dir1=dir,metric='qos',flags=f'with_{f}')
            robustness(dir1=dir,metric='qosPerCost',flags=f'with_{f}')

            