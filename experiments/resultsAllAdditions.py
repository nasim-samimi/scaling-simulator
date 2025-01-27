import pandas as pd
import os
import matplotlib.pyplot as plt
import numpy as np
import sys
from results import runtimes, qosPerCost, qos, cost, compareBaselines

figsize=(15, 10)
nodeHeus=['MinMin','MaxMax']
partitionHeus=['bestfit','worstfit']
additions=[0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1.0]



if __name__ == '__main__':
    dir1='improved/'
    dir2='baseline/'

    for addition in additions: 
        for n in nodeHeus:
            for p in partitionHeus:
                nodeHeu=n
                partitionHeu=p
                print(nodeHeu, partitionHeu,addition)
                print(f"{addition}.csv")
                avgr=runtimes(dir1=dir1,dir2=dir2)
                print(avgr)
                avgqpc=qosPerCost(dir1=dir1,dir2=dir2)
                print(avgqpc)
                avgqos=qos(dir1=dir1,dir2=dir2)
                print(avgqos)
                avgcost=cost(dir1=dir1,dir2=dir2)
                print(avgcost)
        compareBaselines(dir=dir2)

