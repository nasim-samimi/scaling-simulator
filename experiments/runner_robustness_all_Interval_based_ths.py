import pandas as pd
import os
import sys
from results import *
# import sleep


# flags=['intra_node_realloc','intra_node_reduced','interval_based'] #did we have more flags?
flags=['upgrade_service','intra_node_realloc','intra_node_reduced','interval_based'] #did we have more flags?
node_sizes = [8]
# intervals=[5,10,20,30,40,50,100]
intervals=[5,10,20,30,40,50]
intervals=[30]


if __name__=='__main__':
    # run()
    for n in node_sizes:
        for i in intervals:
            dir=f'interval={i}'
            dir2='baseline_interval_based'
            robustness_max_scaling_size_interval_based(dir1=dir,dir2=dir2,metric='cost',flags=f'all_interval_based/interval={i}',nodesize=n)            
            robustness_max_scaling_size_interval_based(dir1=dir,dir2=dir2,metric='qos',flags=f'all_interval_based/interval={i}',nodesize=n)
            robustness_max_scaling_size_interval_based(dir1=dir,dir2=dir2,metric='qosPerCost',flags=f'all_interval_based/interval={i}',nodesize=n)
        dir=f'improved/interval_based'
        robustness_interval_length_3d_sheets(dir1=dir,metric='cost',flags='all_interval_based',nodesize=n)            
        robustness_interval_length_3d_sheets(dir1=dir,metric='qos',flags='all_interval_based',nodesize=n)
        robustness_interval_length_3d_sheets(dir1=dir,metric='qosPerCost',flags='all_interval_based',nodesize=n)


            