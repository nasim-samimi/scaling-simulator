import pandas as pd
import os
import sys
from results import *
# import sleep


# flags=['interval_based'] #did we have more flags?
flags=['intra_node_realloc','intra_node_reduced']#,'upgrade_service','interval_based'] #did we have more flags?
node_sizes = [8]


if __name__=='__main__':
    # run()
    for n in node_sizes:
        for f in flags:

            dir=f'improved/with_{f}'
            robustness_max_scaling_size(dir1=dir,metric='cost',flags=f'with_{f}',nodesize=n)            
            robustness_max_scaling_size(dir1=dir,metric='qos',flags=f'with_{f}',nodesize=n)
            robustness_max_scaling_size(dir1=dir,metric='qosPerCost',flags=f'with_{f}',nodesize=n)


            