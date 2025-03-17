import pandas as pd
import os
import sys
from results import *
# import sleep


flags=['interval_based','intra_node_realloc','intra_node_reduced','intra_domain_realloc'] #did we have more flags?
# flags=['upgrade_service','node_reclaim','intra_node_realloc','intra_node_reduced','intra_node_removed','interval_based'] #did we have more flags?
node_sizes = [8,12]
n=1

if __name__=='__main__':


    dir=f'improved/allOpts/'
    robustness_compare_nodesize(dir1=dir,metric='cost',flags='all')            
    robustness_compare_nodesize(dir1=dir,metric='qos',flags='all')
    robustness_compare_nodesize(dir1=dir,metric='qosPerCost',flags='all')

            