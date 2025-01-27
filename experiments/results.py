import pandas as pd
import os
import matplotlib.pyplot as plt
import numpy as np
import sys

figsize=(15, 10)
nodeHeus=['MinMin','MaxMax']
partitionHeus=['bestfit','worstfit']
main_dir = 'experiments/results/'
plots=main_dir+'plots/'

def plotfiles(main_dir,dir,addition,metric):
    leg=[]
    avg=[]
    print("current dir:",dir)
    fulldir=f'{dir}'
    for files in os.listdir(fulldir):
                # Read the CSV file
        if files=="":
            print("empty")
            print(fulldir)
            continue
        qosPerCost = pd.read_csv(f'{fulldir}{files}', header=None)
        qosPerCost.columns = [metric]

        # Sort the data
        sorted_data = np.sort(qosPerCost[metric])
        max_index = np.argmax(sorted_data)
        avg_value = np.average(sorted_data)
        max_value = sorted_data[max_index]
            # avg_value = sorted_data[avg_index]
        avg.append(avg_value)

        # Plot the data as a line plot
        plt.plot(qosPerCost[metric], label=f'{files}', linestyle='-')
        leg.append(files[:-4])
    plt.legend(leg)

    plt.xlabel('Index')
    plt.ylabel(metric)
    plt.title('Line Plot of metric Data')
    plt.grid(True)
    
    return

def processfiles(dir,addition,metric,averages,leg):
    avg=[]
    print("current dir:",dir)
    fulldir=f'{dir}'
    for files in os.listdir(fulldir):
                # Read the CSV file
        if files=="":
            print("empty")
            print(fulldir)
            continue
        qosPerCost = pd.read_csv(f'{fulldir}{files}', header=None)
        qosPerCost.columns = [metric]

        # Sort the data
        sorted_data = np.sort(qosPerCost[metric])
        max_index = np.argmax(sorted_data)
        avg_value = np.average(sorted_data)
        max_value = sorted_data[max_index]
            # avg_value = sorted_data[avg_index]
        avg.append(avg_value)

        # Plot the data as a line plot
        plt.plot(qosPerCost[metric], label=f'{files}', linestyle='-')
        leg.append(files[:-4])
        averages.loc[len(averages)]=[avg_value,files[:-4]]
    plt.legend(leg)
    averages.loc[len(averages)]=[np.average(avg), 'Baseline']
    plt.xlabel('Index')
    plt.ylabel(metric)
    plt.title('Line Plot of metric Data')
    plt.grid(True)
    
    return



def runtimes(dir1='improved/',dir2='baseline/'):
    dirs=[]
    dirs.append( main_dir+dir1)
    dirs.append( main_dir+dir2)
    leg=[]
    avg=[]
    plt.figure(figsize=figsize)
    averages=pd.DataFrame(columns=['averages','heuristics'])
    for dir in dirs:
        fulldir=f'{dir}runtimes/addition={addition}/{nodeHeu}/{partitionHeu}/'
        for files in os.listdir(fulldir):
            
            print(files)
            runtimes=pd.read_csv(f'{fulldir}{files}',header=None)
            runtimes.columns=['runtimes']
            sorted_data = np.sort(runtimes['runtimes'])
            max_index = np.argmax(sorted_data)
            avg_value = np.average(sorted_data)
            max_value = sorted_data[max_index]
            # avg_value = sorted_data[avg_index]
            avg.append(avg_value)
            # Calculate the cumulative probabilities
            cdf = np.arange(1, len(sorted_data) + 1) / len(sorted_data)

            # Plot the CDF
            plt.plot(sorted_data, cdf, marker='.', linestyle='none')
            leg.append(files[:-4])
            averages.loc[len(averages)]=[avg_value,files[:-4]]
    leg.append( 'Baseline')
    
    plt.xlabel('Runtime (ms)')
    plt.ylabel('CDF')
    plt.title('Cumulative Distribution Function (CDF) of Runtimes')
    plt.grid(True)
    averages.loc[len(averages)]=[np.average(avg), 'Baseline']
    plt.legend(leg)
    savingDir=f'{plots}addition={addition}/runs/'
    if not os.path.exists(savingDir):
        os.makedirs(savingDir)
    plt.savefig(f'{savingDir}{nodeHeu}_{partitionHeu}.png')
    plt.close()
    return

def qosPerCost(dir1='improved/',dir2='baseline/'):
    dirs=[]
    leg=[]
    dirs.append( main_dir+dir1)
    dirs.append( main_dir+dir2)

    plt.figure(figsize=figsize)
    averages=pd.DataFrame(columns=['averages','heuristics'])
    for dir in dirs:
        fulldir=f'{dir}qosPerCost/addition={addition}/{nodeHeu}/{partitionHeu}/'
        processfiles(fulldir,addition,'qosPerCost',averages,leg)
    savingDir=f'{plots}addition={addition}/QpC/'
    if not os.path.exists(savingDir):
        os.makedirs(savingDir)
    plt.savefig(f'{savingDir}{nodeHeu}_{partitionHeu}.png')
    plt.close()
    return averages

def qos(dir1='improved/',dir2='baseline/'):
    
    dirs=[]
    leg=[]
    dirs.append( main_dir+dir1)
    dirs.append( main_dir+dir2)
    plt.figure(figsize=figsize)
    averages=pd.DataFrame(columns=['averages','heuristics'])
    for dir in dirs:
        fulldir=f'{dir}qos/addition={addition}/{nodeHeu}/{partitionHeu}/'
        processfiles(fulldir,addition,'qos',averages,leg)
    savingDir=f'{plots}addition={addition}/QoS/'
    if not os.path.exists(savingDir):
        os.makedirs(savingDir)
    plt.savefig(f'{savingDir}{nodeHeu}_{partitionHeu}.png')
    plt.close()
    return averages
def cost(dir1='improved/',dir2='baseline/'):

    dirs=[]
    leg=[]
    dirs.append( main_dir+dir1)
    dirs.append( main_dir+dir2)
    plt.figure(figsize=figsize)
    averages=pd.DataFrame(columns=['averages','heuristics'])
    for dir in dirs:
        fulldir=f'{dir}cost/addition={addition}/{nodeHeu}/{partitionHeu}/'
        processfiles(fulldir,addition,'cost',averages,leg)
        savingDir=f'{plots}addition={addition}/Cost/'
    if not os.path.exists(savingDir):
        os.makedirs(savingDir)
    plt.savefig(f'{savingDir}{nodeHeu}_{partitionHeu}.png')
    plt.close()
    return averages



def compareBaselines(dir='baseline/'):

    leg=[]
    
    dirs=[]
    for n in ['MinMin','MaxMax']:
        for p in ['bestfit','worstfit']:
            dirs.append(f'{n}/{p}/')
    plt.figure(figsize=figsize)
    for d in dirs:
        fulldir=f'{main_dir}{dir}runtimes/addition={addition}/{d}'
        for files in os.listdir(fulldir):
            print("file",files)
            if files=="":
                print("empty")
                print(fulldir)
                continue
            runtimes=pd.read_csv(f'{fulldir}{files}',header=None)
            runtimes.columns=['runtimes']
            sorted_data = np.sort(runtimes['runtimes'])

            # Calculate the cumulative probabilities
            cdf = np.arange(1, len(sorted_data) + 1) / len(sorted_data)

            # Plot the CDF
            plt.plot(sorted_data, cdf, marker='.', linestyle='none')
            leg.append(d)
    plt.legend(leg)
    plt.xlabel('Runtime (ms)')
    plt.ylabel('CDF')
    plt.title('Cumulative Distribution Function (CDF) of Runtimes')
    plt.grid(True)
    plt.savefig(f'{main_dir}{addition}_runs_baselines.png')
    plt.close()

    for d in dirs:
        fulldir=f'{main_dir}{dir}qosPerCost/addition={addition}/{d}'
        plotfiles(main_dir,fulldir,addition,'qosPerCost')
        fulldir=f'{main_dir}{dir}qos/addition={addition}/{d}'
        plotfiles(main_dir,fulldir,addition,'qos_' )
        fulldir=f'{main_dir}{dir}cost/addition={addition}/{d}'
        plotfiles(main_dir,fulldir,addition,'cost')

    
if __name__ == '__main__':
    dir1='improved/'
    dir2='baseline/'
    if len(sys.argv) != 4:    
        print('Usage: python3 results.py  <addition> <nodeHeu> <partitionHeu> ')
        sys.exit(1)


    nodeHeu = sys.argv[2]
    partitionHeu = sys.argv[3]
    addition = sys.argv[1]
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
