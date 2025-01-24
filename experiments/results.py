import pandas as pd
import os
import matplotlib.pyplot as plt
import numpy as np
import sys

if len(sys.argv) < 4:
    print('Usage: python3 results.py <nodeHeu> <partitionHeu> <addition>')
    sys.exit(1)

nodeHeu = sys.argv[1]
partitionHeu = sys.argv[2]
addition = sys.argv[3]
print(nodeHeu, partitionHeu,addition)
print(f"{addition}.csv")
def runtimes(dir1='improved/',dir2='baseline/'):
    main_dir = 'experiments/results/'
    dirs=[]
    dirs.append( main_dir+dir1)
    dirs.append( main_dir+dir2)
    leg=[]
    avg=[]
    plt.figure(figsize=(10, 6))
    for dir in dirs:
        for files in os.listdir(dir):
            if nodeHeu in files and partitionHeu in files and f"{addition}.csv" in files:
                if 'runtimes' in files:
                    print(files)
                    runtimes=pd.read_csv(f'{dir}{files}',header=None)
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
                    plt.xlabel('Runtime (ms)')
                    plt.ylabel('CDF')
                    plt.title('Cumulative Distribution Function (CDF) of Runtimes')
                    plt.grid(True)
                    leg.append(files[9:-4])
    print(leg.append( 'Baseline'))
    plt.legend(leg)
    plt.savefig(f'{main_dir}{addition}_runs_{nodeHeu}_{partitionHeu}.png')
    plt.close()
    return avg

def qosPerCost(dir1='improved/',dir2='baseline/'):
    main_dir = 'experiments/results/'
    dirs=[]
    dirs.append( main_dir+dir1)
    dirs.append( main_dir+dir2)
    leg=[]
    avg=[]
    plt.figure(figsize=(10, 6))
    for dir in dirs:
        for files in os.listdir(dir):
            if nodeHeu in files and partitionHeu in files and f"{addition}.csv" in files:
                if 'qosPerCost' in files:
                    # Read the CSV file
                    qosPerCost = pd.read_csv(f'{dir}{files}', header=None)
                    qosPerCost.columns = ['qosPerCost']

                    # Sort the data
                    sorted_data = np.sort(qosPerCost['qosPerCost'])
                    max_index = np.argmax(sorted_data)
                    avg_value = np.average(sorted_data)
                    max_value = sorted_data[max_index]
                    # avg_value = sorted_data[avg_index]
                    avg.append(avg_value)
                    # Plot the data as a line plot
                    plt.plot(qosPerCost['qosPerCost'], label=f'{files}', linestyle='-')
                    # plt.scatter(max_index, max_value, color='red')  # Mark the point
                    # plt.annotate(f'Max: ({max_value:.2f})',
                    #              xy=(max_index, max_value),
                    #              xytext=(max_index - 800, max_value + 0.5),  # Adjusted position
                    #              arrowprops=dict(facecolor='black', arrowstyle='->'),
                    #              fontsize=8)
                    # plt.plot(sorted_data, label=f'{files}', linestyle='-')
                    plt.xlabel('Index')
                    plt.ylabel('QoS/Cost')
                    plt.title('Line Plot of QoS/Cost Data')
                    plt.grid(True)
                    plt.legend()
                    leg.append(files[11:-4])
    print(leg.append( 'Baseline'))
    plt.legend(leg)
    # plt.show()
    plt.savefig(f'{main_dir}{addition}_QpC_{nodeHeu}_{partitionHeu}.png')
    plt.close()
    return avg

def qos(dir1='improved/',dir2='baseline/'):
    main_dir = 'experiments/results/'
    dirs=[]
    dirs.append( main_dir+dir1)
    dirs.append( main_dir+dir2)
    leg=[]
    avg=[]
    plt.figure(figsize=(10, 6))
    for dir in dirs:
        for files in os.listdir(dir):
            if nodeHeu in files and partitionHeu in files and f"{addition}.csv" in files:
                if 'qos_' in files:
                    # Read the CSV file
                    qos = pd.read_csv(f'{dir}{files}', header=None)
                    qos.columns = ['qos']

                    # Sort the data
                    sorted_data = np.sort(qos['qos'])
                    max_index = np.argmax(sorted_data)
                    avg_value = np.average(sorted_data)
                    max_value = sorted_data[max_index]
                    # avg_value = sorted_data[avg_index]
                    avg.append(avg_value)
                    # Plot the data as a line plot
                    plt.plot(qos['qos'], label=f'{files}', linestyle='-')
                    # plt.scatter(max_index, max_value, color='red')  # Mark the point
                    # plt.annotate(f'Max: ({max_value:.2f})',
                    #              xy=(max_index, max_value),
                    #              xytext=(max_index - 800, max_value + 0.5),  # Adjusted position
                    #              arrowprops=dict(facecolor='black', arrowstyle='->'),
                    #              fontsize=8)
                    # plt.plot(sorted_data, label=f'{files}', linestyle='-')
                    plt.xlabel('Index')
                    plt.ylabel('QoS')
                    plt.title('Line Plot of QoS Data')
                    plt.grid(True)
                    plt.legend()
                    leg.append(files[:-4])
    print(leg.append( 'Baseline'))
    plt.legend(leg)
    # plt.show()
    plt.savefig(f'{main_dir}{addition}_QoS_{nodeHeu}_{partitionHeu}.png')
    plt.close()
    return avg
def cost(dir1='improved/',dir2='baseline/'):
    main_dir = 'experiments/results/'
    dirs=[]
    dirs.append( main_dir+dir1)
    dirs.append( main_dir+dir2)
    leg=[]
    avg=[]
    plt.figure(figsize=(10, 6))
    for dir in dirs:
        for files in os.listdir(dir):
            if nodeHeu in files and partitionHeu in files and f"{addition}.csv" in files:
                if 'cost' in files:
                    # Read the CSV file
                    cost = pd.read_csv(f'{dir}{files}', header=None)
                    cost.columns = ['cost']

                    # Sort the data
                    sorted_data = np.sort(cost['cost'])
                    max_index = np.argmax(sorted_data)
                    avg_value = np.average(sorted_data)
                    max_value = sorted_data[max_index]
                    # avg_value = sorted_data[avg_index]
                    avg.append(avg_value)
                    # Plot the data as a line plot
                    plt.plot(cost['cost'], label=f'{files}', linestyle='-')
                    # plt.scatter(max_index, max_value, color='red')  # Mark the point
                    # plt.annotate(f'Max: ({max_value:.2f})',
                    #              xy=(max_index, max_value),
                    #              xytext=(max_index - 800, max_value + 0.5),  # Adjusted position
                    #              arrowprops=dict(facecolor='black', arrowstyle='->'),
                    #              fontsize=8)
                    # plt.plot(sorted_data, label=f'{files}', linestyle='-')
                    plt.xlabel('Index')
                    plt.ylabel('Cost')
                    plt.title('Line Plot of Cost Data')
                    plt.grid(True)
                    plt.legend()
                    leg.append(files[:-4])
    print(leg.append( 'Baseline'))
    plt.legend(leg)
    # plt.show()
    plt.savefig(f'{main_dir}{addition}_Cost_{nodeHeu}_{partitionHeu}.png')
    plt.close()
    return avg

if __name__ == '__main__':
    dir1='improved/'
    dir2='baseline/'
    avgr=runtimes(dir1=dir1,dir2=dir2)
    print(avgr)
    avgqpc=qosPerCost(dir1=dir1,dir2=dir2)
    print(avgqpc)
    avgqos=qos(dir1=dir1,dir2=dir2)
    print(avgqos)
    avgcost=cost(dir1=dir1,dir2=dir2)
    print(avgcost)
