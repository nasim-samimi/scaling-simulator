import pandas as pd
import os
import sys
from results_cores import *

flags=['interval_based','intra_node_realloc','intra_node_reduced','intra_domain_realloc']
node_sizes = [8,12]
n=1

if __name__=='__main__':
    dir=f'improved/allOpts'
    robustness_compare_nodesize(dir1=dir,metric='cost',flags='all')            
    robustness_compare_nodesize(dir1=dir,metric='qos',flags='all')
    robustness_compare_nodesize(dir1=dir,metric='qosPerCost',flags='all')

    dir=f'improved/allOpts'
    robustness_compare_nodesize_none(dir1=dir,metric='cost',flags='all')            
    robustness_compare_nodesize_none(dir1=dir,metric='qos',flags='all')
    robustness_compare_nodesize_none(dir1=dir,metric='qosPerCost',flags='all')

