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
def runtimes(dir1='heuristics/',dir2='baseline/'):
    main_dir = 'experiments/results/'
    dirs=[]
    dirs.append( main_dir+dir1)
    dirs.append( main_dir+dir2)
    leg=[]
    plt.figure(figsize=(10, 6))
    for dir in dirs:
        for files in os.listdir(dir):
            if nodeHeu in files and partitionHeu in files and f"{addition}.csv" in files:
                if 'runtimes' in files:
                    print(files)
                    runtimes=pd.read_csv(f'{dir}{files}',header=None)
                    runtimes.columns=['runtimes']
                    sorted_data = np.sort(runtimes['runtimes'])

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
    return

def qosPerCost():
    main_dir = 'experiments/results/'
    dirs=[]
    dirs.append( main_dir+'heuristics/')
    dirs.append( main_dir+'baseline/')
    leg=[]
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
                    max_value = sorted_data[max_index]
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
    return

if __name__ == '__main__':
    runtimes(dir1='heuristics/',dir2='baseline/')
    qosPerCost()
