import pandas as pd
import os
import sys
from results_cores import *

flags=['upgrade_service','intra_node_realloc','intra_node_reduced','interval_based']
node_sizes = [8]

if __name__=='__main__':
    for n in node_sizes:
        dir=f'improved/allOpts'
        robustness_max_scaling_size(dir1=dir,metric='cost',flags='all',nodesize=n)            
        robustness_max_scaling_size(dir1=dir,metric='qos',flags='all',nodesize=n)
        robustness_max_scaling_size(dir1=dir,metric='qosPerCost',flags='all',nodesize=n)

        robustness_max_scaling_size_3d_sheets(dir1=dir,metric='cost',flags='all',nodesize=n)            
        robustness_max_scaling_size_3d_sheets(dir1=dir,metric='qos',flags='all',nodesize=n)
        robustness_max_scaling_size_3d_sheets(dir1=dir,metric='qosPerCost',flags='all',nodesize=n)

