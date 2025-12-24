import pandas as pd
import os
import sys
from results_cores import *

flags=['intra_node_reduced','intra_node_realloc','interval_based']
node_sizes = [8,16]

if __name__=='__main__':
    for n in node_sizes:
        dir=f'improved/allOpts'
        robustness_compare_node_core_selection(dir1=dir,metric='cost',flags='all',nodesize=n)            
        robustness_compare_node_core_selection(dir1=dir,metric='qos',flags='all',nodesize=n)
        robustness_compare_node_core_selection(dir1=dir,metric='qosPerCost',flags='all',nodesize=n)

