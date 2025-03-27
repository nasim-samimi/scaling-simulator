import pandas as pd
import os
import sys
from results import *
# import sleep


# flags=['intra_node_realloc','intra_node_reduced','interval_based'] #did we have more flags?
flags=['upgrade_service','intra_node_realloc','intra_node_reduced','interval_based'] #did we have more flags?
node_sizes = [8]


if __name__=='__main__':
    # run()
    for n in node_sizes:
        dir=f'improved/none_interval_based'
        dir2='baseline_interval_based'
        robustness_max_scaling_size(dir1=dir,dir2=dir2,metric='cost',flags='none_interval_based',nodesize=n)            
        robustness_max_scaling_size(dir1=dir,dir2=dir2,metric='qos',flags='none_interval_based',nodesize=n)
        robustness_max_scaling_size(dir1=dir,dir2=dir2,metric='qosPerCost',flags='none_interval_based',nodesize=n)


            