import pandas as pd
import os
import sys
from results import *
# import sleep


flags=['upgrade_service','intra_node_realloc','intra_node_reduced','interval_based'] #did we have more flags?
# flags=['upgrade_service','intra_node_realloc','intra_node_reduced','interval_based'] #did we have more flags?
node_sizes = [8]


if __name__=='__main__':
    # run()
    for n in node_sizes:
                dir=f'improved/allOpts/'
                runtimes(dir1=dir,flags='all',nodesize=n)            


            