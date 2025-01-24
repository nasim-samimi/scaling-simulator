import pandas as pd
import os
import sys

PARTITIONING_H=['bestfit','worstfit']

REALLOCATION_H=[]
NODE_SELECTION_H=["MinMin","MaxMax"]
addition=0
if len(sys.argv)>1:
    addition=sys.argv[1]

def run():
    for p in PARTITIONING_H:
        for n in NODE_SELECTION_H:
            
            df=pd.DataFrame(columns=['ReallocationHeuristic','NodeSelectionHeuristic','PartitioningHeuristic'])
            df['ReallocationHeuristic']=[""]
            df['NodeSelectionHeuristic']=[n]
            df['PartitioningHeuristic']=[p]
            df.to_csv('../data/heuristics.csv', index=False)
            print('running:',"",n,p)

            try: 
                os.system(f'go run main.go {addition} > log.txt')
                # os.system('cd ..')
                os.chdir('..')
                os.system(f'python3 experiments/results.py {n} {p} {addition}')
                os.chdir('scaling_baseline')
            except:
                print('Error:',"",n,p)
                return

if __name__=='__main__':
    run()                
            
            