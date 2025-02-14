import pandas as pd
import os
import sys
from results import *
# import sleep


flags=['interval_based'] #did we have more flags?
# flags=['upgrade_service','node_reclaim','intra_node_realloc','intra_node_reduced','intra_node_removed'] #did we have more flags?
node_sizes = [8,12]


if __name__=='__main__':
    # run()
    for f in flags:
        for n in node_sizes: 
            if 'intra' in f:

                dir=f'improved/with_{f}/'
                robustness(dir1=dir,metric='cost',flags=f'with_{f}',nodesize=n)            
                robustness(dir1=dir,metric='qos',flags=f'with_{f}',nodesize=n)
                robustness(dir1=dir,metric='qosPerCost',flags=f'with_{f}',nodesize=n)
            else:
                dir=f'improved/with_{f}/'   
                robustness(dir1=dir,metric='cost',flags=f'with_{f}',nodesize=n)            
                robustness(dir1=dir,metric='qos',flags=f'with_{f}',nodesize=n)
                robustness(dir1=dir,metric='qosPerCost',flags=f'with_{f}',nodesize=n)

            