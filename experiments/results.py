import pandas as pd
import os
import matplotlib.pyplot as plt
import numpy as np
import sys

figsize=(8,8)
# nodeHeus=['MinMin','MaxMax',]
nodeHeus=['MMRB','mmRB','MmRB','mMRB']
partitionHeus=['bestfit','worstfit']
reallocationHeus=["HBCI","HBI","HCI","HB","HC","HBC","LB","LC","LBC"]
main_dir = 'experiments/results/'
plots=main_dir+'plots/'
ADDITIONS=[0,0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1]
node_sizes = [8,12,16,20,24,28,32]
events_dir='data/events/hightraffic'


def plotfiles(main_dir,dir,addition,metric,dirs):
    leg=[]
    avg=[]
    print("current dir:",dir)
    
    for d in dirs:
        fulldir=f'{dir}{d}'
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
            leg.append(d)
    plt.legend(leg)

    plt.xlabel('Index')
    plt.ylabel(metric)
    plt.title('Line Plot of metric Data')
    plt.grid(True)
    if not os.path.exists(f'{main_dir}plots/baselines/addition={addition}'):
        os.makedirs(f'{main_dir}plots/baselines/addition={addition}')
    plt.savefig(f'{main_dir}plots/baselines/addition={addition}/{metric}_baseline.png')
    plt.close()
    
    return

def processfile(dir,addition,metric,averages,leg,times):
    avg=[]
    print("current dir:",dir)
    fulldir=f'{dir}'
    times.columns=['eventTime'] 
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
        plt.plot(times['eventTime'],qosPerCost[metric], label=f'{files}', linestyle='-')
        leg.append(files[:-4])
        averages.loc[len(averages)]=[avg_value,files[:-4]]
    plt.legend(leg)
    plt.xlabel('Index')
    plt.ylabel(metric)
    plt.title('Line Plot of metric Data')
    plt.grid(True)
    
    return

def processfiles(dir1='improved/',dir2='baseline/',metric='cost'):

    dirs=[]
    leg=[]
    dirs.append( main_dir+dir1)
    dirs.append( main_dir+dir2)
    plt.figure(figsize=figsize)
    averages=pd.DataFrame(columns=['averages','heuristics'])
    for dir in dirs:
        fulldir=f'{dir}{metric}/nodesize={node_size}/addition={addition}/{nodeHeu}/{partitionHeu}/'
        times=pd.read_csv(f'{dir}eventTime/nodesize={node_size}/addition={addition}/{nodeHeu}/{partitionHeu}/')
        processfile(fulldir,addition,metric,averages,leg,times)
        savingDir=f'{plots}nodesize={node_size}/addition={addition}/{metric}/'
    if not os.path.exists(savingDir):
        os.makedirs(savingDir)
    plt.savefig(f'{savingDir}{nodeHeu}_{partitionHeu}.png')
    averages.to_csv(f'{savingDir}{nodeHeu}_{partitionHeu}_averages.csv',index=False)
    plt.close()
    return averages

def runtimes(dir1='improved/',dir2='baseline/'):
    dirs=[]
    dirs.append( main_dir+dir1)
    dirs.append( main_dir+dir2)
    leg=[]
    avg=[]
    plt.figure(figsize=figsize)
    averages=pd.DataFrame(columns=['averages','heuristics'])
    for dir in dirs:
        fulldir=f'{dir}runtimes/nodesize={node_size}/addition={addition}/{nodeHeu}/{partitionHeu}/'
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
    plt.legend(leg)
    savingDir=f'{plots}nodesize={node_size}/addition={addition}/runs/'
    if not os.path.exists(savingDir):
        os.makedirs(savingDir)
    plt.savefig(f'{savingDir}{nodeHeu}_{partitionHeu}.png')
    averages.to_csv(f'{savingDir}{nodeHeu}_{partitionHeu}_averages.csv',index=False)
    plt.close()
    return


def time_based_avg(values,times):
    # print("values",values)
    # events=pd.read_csv(f'{events_dir}/events_{addition}.csv')
    events=times
    # add series of values as a new column to events dataframe because events does not have this column
    df = pd.DataFrame({'values': values, 'EventTime': events['EventTime']})
    # print("df",df)
    df.dropna(subset=["values"], inplace=True)
    df = df.loc[df.groupby("EventTime")["values"].idxmax()]

    df = df.sort_values(by="EventTime").reset_index(drop=True)
    # print(df["EventTime"].tolist())
    # print("df after grouping",df)
    df["EventTime"] = df["EventTime"].diff(periods=-1).abs().fillna(0).astype(float)
    
    # print("df after diff",df)
    df["total"] = df["EventTime"] * df["values"]
    # print("df after total",df)
    avg=df["total"].sum()/events['EventTime'].max()

    return avg


def compareBaselines(dir='baseline/'):

    # leg=[]
    
    # dirs=[]
    # for n in ['MinMin','MaxMax']:
    #     for p in ['bestfit','worstfit']:
    #         dirs.append(f'{n}/{p}/')
    # plt.figure(figsize=figsize)
    # for addition in ADDITIONS:
    #     for d in dirs:
    #         fulldir=f'{main_dir}{dir}runtimes/nodesize={node_size}/addition={addition}/{d}'
    #         for files in os.listdir(fulldir):
    #             print("file",files)
    #             if files=="":
    #                 print("empty")
    #                 print(fulldir)
    #                 continue
    #             runtimes=pd.read_csv(f'{fulldir}{files}',header=None)
    #             runtimes.columns=['runtimes']
    #             sorted_data = np.sort(runtimes['runtimes'])

    #             # Calculate the cumulative probabilities
    #             cdf = np.arange(1, len(sorted_data) + 1) / len(sorted_data)

    #             # Plot the CDF
    #             plt.plot(sorted_data, cdf, marker='.', linestyle='none')
    #             leg.append(d)
    #     plt.legend(leg)
    #     plt.xlabel('Runtime (ms)')
    #     plt.ylabel('CDF')
    #     plt.title('Cumulative Distribution Function (CDF) of Runtimes')
    #     plt.grid(True)
    #     if not os.path.exists(f'{main_dir}/plots/baselines/nodesize={node_size}/addition={addition}'):
    #         os.makedirs(f'{main_dir}/plots/baselines/nodesize={node_size}/addition={addition}')
    #     plt.savefig(f'{main_dir}/plots/baselines/nodesize={node_size}/addition={addition}/runs_baselines.png')
    #     plt.close()

    #     fulldir=f'{main_dir}{dir}qosPerCost/nodesize={node_size}/addition={addition}/'
    #     plotfiles(main_dir,fulldir,addition,'qosPerCost',dirs)
    #     fulldir=f'{main_dir}{dir}qos/nodesize={node_size}/addition={addition}/'
    #     plotfiles(main_dir,fulldir,addition,'qos' ,dirs)
    #     fulldir=f'{main_dir}{dir}cost/nodesize={node_size}/addition={addition}/'
    #     plotfiles(main_dir,fulldir,addition,'cost',dirs)

    avg=[]
    columns=[]
    averages=pd.DataFrame(columns=['averages','heuristics','addition'])
    
    #cost
    for node_size in node_sizes:
        for metric in ['cost','qos','qosPerCost']:
            leg=[]
            i=0
            averages=pd.DataFrame(columns=['averages','heuristics','addition'])
            for n in nodeHeus:
                for p in partitionHeus:
                    avg=[]
                    
                    
                    
                    for a in ADDITIONS:
                        fulldir=f'{main_dir}{dir}{metric}/nodesize={node_size}/addition={a}/{n}/{p}/'
                        Times_addr= f'{main_dir}{dir}eventTime/nodesize={node_size}/addition={a}/MaxMax/{p}/'
                # print(fulldir)
                        print(fulldir)
                        for files in os.listdir(fulldir):
                            Values = pd.read_csv(f'{fulldir}{files}', header=None)
                            Times=pd.read_csv(f'{Times_addr}{files}',header=None)
                            Times.columns=['EventTime']
                            Values.columns = [metric]

                            # avg_value = time_based_avg(Values[metric],a)
                            avg_value = time_based_avg(Values[metric],Times)
                            avg.append(avg_value)
                            averages.loc[i]=[avg_value,f'{n}-{p}-nodesize={node_size}',a]
                            i=i+1
        
                    print(leg)
                    leg.append(f'{n}-{p}-nodesize={node_size}')
            plt.figure(figsize=figsize)
            for l in leg:
                avgs = averages[averages['heuristics'] == l]
                avgs = avgs.sort_values(by='addition')
                plt.plot(avgs['addition'], avgs['averages'], marker='o')
            plt.grid(True)
            plt.xlabel('randomness')
            # plt.xticks(range(len(ADDITIONS)), ADDITIONS)
            plt.ylabel(metric)
            plt.title(f'Robustness comparison for {metric} - {n}-{p}')    
            plt.legend(leg)
            savingDir=f'{plots}/robustness/baselines/nodesize={node_size}/'
            if not os.path.exists(savingDir):
                os.makedirs(savingDir)
            plt.savefig(f'{savingDir}robustness_{metric}.png')
            plt.close()

    return

def robustness(dir1='improved/allOpts',dir2='baseline/',metric='cost',flags='allOpts',nodesize=8):

    avg=[]
    columns=[]
    averages=pd.DataFrame(columns=['averages','heuristics','addition'])
    dirs=[main_dir+dir1,main_dir+dir2]
    leg=[]
    #cost
    
    for n in nodeHeus:
        for p in partitionHeus:
            avg=[]
            
            i=0
            averages=pd.DataFrame(columns=['averages','heuristics','addition'])
            for a in ADDITIONS:
                leg=[]
                # for heuristics
                dir=dirs[0]
                fulldir=f'{dir}{metric}/nodesize={nodesize}/addition={a}/{n}/{p}/'
                Times_addr= f'{dir}eventTime/nodesize={nodesize}/addition={a}/{n}/{p}/'
                # print(fulldir)
                for files in os.listdir(fulldir):
                    Values = pd.read_csv(f'{fulldir}{files}', header=None)
                    Times=pd.read_csv(f'{Times_addr}{files}',header=None)
                    Times.columns=['EventTime']
                    Values.columns = [metric]

                    avg_value = time_based_avg(Values[metric],Times) 
                    averages.loc[i]=[avg_value,dir[20:]+files[:-4],a]
                    leg.append(dir[20:]+files[:-4])
                    i+=1
                # for baseline
                dir=dirs[1]
                fulldir=f'{dir}{metric}/nodesize={nodesize}/addition={a}/MaxMax/{p}/'
                Times_addr= f'{dir}eventTime/nodesize={nodesize}/addition={a}/MaxMax/{p}/'
                # print(fulldir)
                for files in os.listdir(fulldir):
                    Values = pd.read_csv(f'{fulldir}{files}', header=None)
                    Times=pd.read_csv(f'{Times_addr}{files}',header=None)
                    Times.columns=['EventTime']
                    Values.columns = [metric]

                    avg_value = time_based_avg(Values[metric],Times) 
                    averages.loc[i]=[avg_value,dir[20:]+files[:-4],a]
                    leg.append(dir[20:]+files[:-4])
                    i+=1
 
            print(leg)
            plt.figure(figsize=figsize)
            # colors = plt.cm.rainbow(np.linspace(0, 1, 15))
            heuristic_averages = averages.groupby('heuristics')['averages'].mean()
            if metric=='cost':
                top_heuristics = heuristic_averages.nsmallest(4).index
            else:
                top_heuristics = heuristic_averages.nlargest(4).index
            for l in leg:
                if 'baseline' in l:
                    b=l
                    break
            if b not in top_heuristics:
                top_heuristics=top_heuristics.append(pd.Index([b]))
            for l in top_heuristics:
                if 'baseline' in l:
                    marker='o'
                    linestyle='-'
                else:
                    marker='x'
                    linestyle='--'
                avgs = averages[averages['heuristics'] == l]
                avgs = avgs.sort_values(by='addition')
                plt.plot(avgs['addition'], avgs['averages'], marker=marker,linestyle=linestyle)
            plt.grid(True)
            plt.xlabel('randomness')
            # plt.xticks(range(len(ADDITIONS)), ADDITIONS)
            plt.ylabel(metric)
            plt.title(f'Robustness comparison for {metric} - {n}-{p}')    
            plt.legend(top_heuristics)
            savingDir=f'{plots}robustness/{flags}/nodesize={nodesize}/{n}/{p}/'
            if not os.path.exists(savingDir):
                os.makedirs(savingDir)
            plt.savefig(f'{savingDir}robustness_{metric}.png')
            plt.close()

    return




if __name__ == '__main__':
    dir1='improved/allOpts/'
    dir2='baseline/'
    # if len(sys.argv) != 5:    
    #     print('Usage: python3 results.py  <addition> <nodeHeu> <partitionHeu> ')
    #     sys.exit(1)


    # nodeHeu = sys.argv[2]
    # partitionHeu = sys.argv[3]
    # addition = sys.argv[1]
    # node_size=sys.argv[4]


    # print(nodeHeu, partitionHeu,addition)
    # print(f"{addition}.csv")
    # avgr=runtimes(dir1=dir1,dir2=dir2)
    # print(avgr)
    # avgqpc=processfiles(dir1=dir1,dir2=dir2,metric='qosPerCost')
    # print(avgqpc)
    # avgqos=processfiles(dir1=dir1,dir2=dir2,metric='qos')
    # print(avgqos)
    # # avgcost=cost(dir1=dir1,dir2=dir2)
    # avgcost=processfiles(dir1=dir1,dir2=dir2,metric='cost')
    # print(avgcost)  
    nodeHeus=['MaxMax','MinMin']
    compareBaselines(dir=dir2)


# def robustness_withflags(dir1='improved/',dir2='baseline/',metric='cost'):

#     avg=[]
#     columns=[]
#     averages=pd.DataFrame(columns=['averages','heuristics','addition'])
#     dirs=[main_dir+dir1,main_dir+dir2]
#     leg=[]
#     #cost
    
#     for n in nodeHeus:
#         for p in partitionHeus:
#             avg=[]
            
#             i=0
#             averages=pd.DataFrame(columns=['averages','heuristics','addition'])
#             for a in ADDITIONS:
#                 leg=[]
#                 for dir in dirs:
#                     fulldir=f'{dir}{metric}/nodesize={node_size}/addition={a}/{n}/{p}/'
#                     print(fulldir)
#                     for files in os.listdir(fulldir):
#                         qosPerCost = pd.read_csv(f'{fulldir}{files}', header=None)
#                         qosPerCost.columns = [metric]

#                         sorted_data = np.sort(qosPerCost[metric])
#                         max_index = np.argmax(sorted_data)
#                         avg_value = np.average(sorted_data)
#                         max_value = sorted_data[max_index]
#                         avg.append(avg_value)
#                         averages.loc[i]=[avg_value,files[:-4],a]
#                         leg.append(files[:-4])
#                         i+=1
 
#             print(leg)
#             plt.figure(figsize=figsize)
#             for l in leg:
                
#                 avgs = averages[averages['heuristics'] == l]
#                 avgs = avgs.sort_values(by='addition')
#                 plt.plot(avgs['addition'], avgs['averages'], marker='o')
#             plt.grid(True)
#             plt.xlabel('Additions')
#             # plt.xticks(range(len(ADDITIONS)), ADDITIONS)
#             plt.ylabel(metric)
#             plt.title(f'Robustness comparison for {metric} - {n}-{p}')    
#             plt.legend(leg)
#             savingDir=f'{plots}robustness_with_flags/{n}/{p}/'
#             if not os.path.exists(savingDir):
#                 os.makedirs(savingDir)
#             plt.savefig(f'{savingDir}robustness_{metric}.png')
#             plt.close()

# def qosPerCost(dir1='improved/',dir2='baseline/'):
#     dirs=[]
#     leg=[]
#     dirs.append( main_dir+dir1)
#     dirs.append( main_dir+dir2)

#     plt.figure(figsize=figsize)
#     averages=pd.DataFrame(columns=['averages','heuristics'])
#     for dir in dirs:
#         fulldir=f'{dir}qosPerCost/addition={addition}/{nodeHeu}/{partitionHeu}/'
#         times=pd.read_csv(f'{dir}eventTime/addition={addition}/{nodeHeu}/{partitionHeu}/')
#         processfile(fulldir,addition,'qosPerCost',averages,leg,times)
#     savingDir=f'{plots}addition={addition}/QpC/'
#     if not os.path.exists(savingDir):
#         os.makedirs(savingDir)
#     plt.savefig(f'{savingDir}{nodeHeu}_{partitionHeu}.png')
#     averages.to_csv(f'{savingDir}{nodeHeu}_{partitionHeu}_averages.csv',index=False)
#     plt.close()
#     return averages

# def qos(dir1='improved/',dir2='baseline/'):
    
#     dirs=[]
#     leg=[]
#     dirs.append( main_dir+dir1)
#     dirs.append( main_dir+dir2)
#     plt.figure(figsize=figsize)
#     averages=pd.DataFrame(columns=['averages','heuristics'])
#     for dir in dirs:
#         fulldir=f'{dir}qos/addition={addition}/{nodeHeu}/{partitionHeu}/'
#         times=pd.read_csv(f'{dir}eventTime/addition={addition}/{nodeHeu}/{partitionHeu}/')
#         processfile(fulldir,addition,'qos',averages,leg,times)
#     savingDir=f'{plots}addition={addition}/QoS/'
#     if not os.path.exists(savingDir):
#         os.makedirs(savingDir)
#     plt.savefig(f'{savingDir}{nodeHeu}_{partitionHeu}.png')
#     averages.to_csv(f'{savingDir}{nodeHeu}_{partitionHeu}_averages.csv',index=False)
#     plt.close()
#     return averages
# def cost(dir1='improved/',dir2='baseline/'):

#     dirs=[]
#     leg=[]
#     dirs.append( main_dir+dir1)
#     dirs.append( main_dir+dir2)
#     plt.figure(figsize=figsize)
#     averages=pd.DataFrame(columns=['averages','heuristics'])
#     for dir in dirs:
#         fulldir=f'{dir}cost/addition={addition}/{nodeHeu}/{partitionHeu}/'
#         times=pd.read_csv(f'{dir}eventTime/addition={addition}/{nodeHeu}/{partitionHeu}/')
#         processfile(fulldir,addition,'cost',averages,leg,times)
#         savingDir=f'{plots}addition={addition}/Cost/'
#     if not os.path.exists(savingDir):
#         os.makedirs(savingDir)
#     plt.savefig(f'{savingDir}{nodeHeu}_{partitionHeu}.png')
#     averages.to_csv(f'{savingDir}{nodeHeu}_{partitionHeu}_averages.csv',index=False)
#     plt.close()
#     return averages