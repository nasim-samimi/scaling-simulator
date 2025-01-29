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

if __name__=='__main__':
    # run()
    for f in flags: 
        dir=f'improved/with_{f}'   
        robustness(dir1=dir,metric='cost')            
        robustness(dir1=dir,metric='qos')
        robustness(dir1=dir,metric='qosPerCost')

            