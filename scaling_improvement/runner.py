import pandas as pd
import os
import sys
# import sleep

PARTITIONING_H=['bestfit','worstfit']

REALLOCATION_H=["HBCI","HBI","HCI","HB","HC","HBC","LB","LC","LBC"]
# REALLOCATION_H=["HBCI"]
NODE_SELECTION_H=["MinMin","MaxMax"]
addition=0
if len(sys.argv) > 1:
    addition = sys.argv[1]

def run():
    for p in PARTITIONING_H:
        for n in NODE_SELECTION_H:
            for r in REALLOCATION_H:
                df=pd.DataFrame(columns=['ReallocationHeuristic','NodeSelectionHeuristic','PartitioningHeuristic'])
                df['ReallocationHeuristic']=[r]
                df['NodeSelectionHeuristic']=[n]
                df['PartitioningHeuristic']=[p]
                df.to_csv('../data/heuristics.csv', index=False)
                print('running:',r,n,p)

                try: 
                    os.system(f'go run main.go {addition} > log.txt')
                    # os.system('cd ..')
                except:
                    print('Error:',"",n,p)
                    return
            os.chdir('..')
            os.system(f'python3 experiments/results.py {n} {p} {addition}')
            os.chdir('scaling_improvement')
            # sleep(5)

if __name__=='__main__':
    run()                
            
            