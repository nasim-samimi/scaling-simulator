import pandas as pd
import os
import sys
from results import *
# import sleep

PARTITIONING_H=['bestfit','worstfit']

# REALLOCATION_H=["HBCI","HBI","HCI","HB","HC","HBC","LB","LC","LBC"]
REALLOCATION_H=["HBCI"]
NODE_SELECTION_H=["MinMin","MaxMax"]
addition=0
if len(sys.argv) > 1:
    addition = sys.argv[1]

def run():
    for p in PARTITIONING_H:
        for n in NODE_SELECTION_H:
                print('running:',n,p)

                try: 
                    os.system(f'python3 experiments/results.py {addition} {n} {p} ')
                except:
                    print('Error:',"",n,p)
                    return

if __name__=='__main__':
    # run()    
    robustness(metric='cost')            
    robustness(metric='qos')
    robustness(metric='qosPerCost')
            
            