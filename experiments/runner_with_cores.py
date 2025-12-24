#!/usr/bin/env python3
"""
Runner that uses results_cores.py to generate plots with Average Core Demand on x-axis.
Output will go to experiments/plots_with_cores/
"""

import pandas as pd
import os
import sys
from results_cores import *

node_sizes = [8]

if __name__=='__main__':
    for n in node_sizes:
        dir1='improved/allOpts'
        
        print(f"Generating plots with Average Core Demand on x-axis...")
        print(f"Output directory: {plots}")
        print("")
        
        robustness_max_scaling_size(dir1=dir1,metric='cost',flags='all',nodesize=n)            
        robustness_max_scaling_size(dir1=dir1,metric='qos',flags='all',nodesize=n)
        robustness_max_scaling_size(dir1=dir1,metric='qosPerCost',flags='all',nodesize=n)
        
        print("\nDone! Check experiments/plots_with_cores/")

