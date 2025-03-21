import pandas as pd
import os
import sys
from results import *
# import sleep


flags=['intra_node_reduced','intra_node_realloc','interval_based'] #did we have more flags?
# flags=['upgrade_service','node_reclaim','intra_node_realloc','intra_node_reduced','intra_node_removed','interval_based'] #did we have more flags?
node_sizes = [8,16]


if __name__=='__main__':
    # run()
    for n in node_sizes:
        dir=f'improved/allOpts'
        robustness_compare_node_core_selection(dir1=dir,metric='cost',flags='all',nodesize=n)            
        robustness_compare_node_core_selection(dir1=dir,metric='qos',flags='all',nodesize=n)
        robustness_compare_node_core_selection(dir1=dir,metric='qosPerCost',flags='all',nodesize=n)


            