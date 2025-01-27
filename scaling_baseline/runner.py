import pandas as pd
import os
import sys
# import sleep

PARTITIONING_H=['bestfit','worstfit']

# REALLOCATION_H=["HBCI","HBI","HCI","HB","HC","HBC","LB","LC","LBC"]
REALLOCATION_H=["HBCI"]
NODE_SELECTION_H=["MinMin","MaxMax"]
ADDITION=[0,0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1]
def run():
    for a in ADDITION:
        for p in PARTITIONING_H:
            for n in NODE_SELECTION_H:
                    df=pd.DataFrame(columns=['ReallocationHeuristic','NodeSelectionHeuristic','PartitioningHeuristic'])
                    df['ReallocationHeuristic']=[""]
                    df['NodeSelectionHeuristic']=[n]
                    df['PartitioningHeuristic']=[p]
                    df.to_csv('../data/heuristics.csv', index=False)
                    print('running:',n,p)

                    try: 
                        os.system(f'go run main.go {a} > log.txt')
                        # os.system('cd ..')
                    except:
                        print('Error:',"",n,p)
                        return
                # sleep(5)

if __name__=='__main__':
    run()                
            
            