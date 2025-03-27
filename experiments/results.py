import pandas as pd
import os
import matplotlib.pyplot as plt
import numpy as np
import sys
import math

figsize=(5,4)
linewidth=2.6
fontsize=11
fontsizetxt=13

# nodeHeus=['MMRB','mmRB','MmRB','mMRB']
# nodeHeus=['MMRB','mmRB','MmRB','mMRB','MaxMax','MinMin']
# partitionHeus=['bestfit','worstfit']
partitionHeus=['bestfit']
reallocationHeus=["HBCI","HBI","HCI","HB","HC","HBC","LB","LC","LBC"]
main_dir = 'experiments/results/'
plots='experiments/plots/'
# ADDITIONS=[0,0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1]
ADDITIONS=[0,0.2,0.4,0.6,0.8,1.0,1.2,1.4,1.6,1.8,2.0]
ADDITIONS_LABEL=[0,0.3,0.6,0.9,1.2,1.5,1.8,2.1,2.4,2.7,3]
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
    plt.savefig(f'{main_dir}plots/baselines/addition={addition}/{metric}_baseline.pdf', bbox_inches='tight')
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
    plt.savefig(f'{savingDir}{nodeHeu}_{partitionHeu}.pdf')
    averages.to_csv(f'{savingDir}{nodeHeu}_{partitionHeu}_averages.csv',index=False)
    plt.close()
    return averages

def runtimes_enhanced(dir1='improved/',dir2='baseline/',nodesize=8,flags='allOpts'):
    avg=[]
    columns=[]
    averages=pd.DataFrame(columns=['averages','heuristics','addition'])
    dirs=[main_dir+dir1,main_dir+dir2]
    leg=[]
    avg=[]
    max_size=[16,32,64,128]
    plt.figure(figsize=(5,3))
            
    i=0
    averages=pd.DataFrame(columns=['averages','heuristics'])
    for m in max_size:
        Values = pd.DataFrame(columns=['runtimes'])
        for a in ADDITIONS:

            # for heuristics
            dir=dirs[0]+f'/max_scaling_threshold={m}/'
            fulldir=f'{dir}runtimes/nodesize={nodesize}/addition={a}/mmRB/bestfit/'
            # print(fulldir)
            for files in os.listdir(fulldir):
                new_values = pd.read_csv(f'{fulldir}{files}', header=None, names=['runtimes'])
                Values = pd.concat([Values, new_values], axis=0, ignore_index=True)
                

        leg.append(f'max_scaling_threshold={m}')
        i+=1
            # for baseline
        threshold = np.percentile(Values['runtimes'], 99)
        Values = Values[Values['runtimes'] <= threshold]
        sorted_data = np.sort(Values['runtimes'])*1000
        # max_index = np.argmax(sorted_data)
        # avg_value = np.average(sorted_data)
        # max_value = sorted_data[max_index]
        # # avg_value = sorted_data[avg_index]
        # avg.append(avg_value)
        # Calculate the cumulative probabilities
        cdf = np.arange(1, len(sorted_data) + 1) / len(sorted_data)

        # Plot the CDF
        plt.plot(sorted_data, cdf,markersize=3)
        print(leg)
        
        # colors = plt.cm.rainbow(np.linspace(0, 1, 15))


        plt.grid(True)
        plt.xlabel('Runtime (\u00B5s)')
        plt.ylabel('CDF')
        # plt.title(f'Robustness comparison for {metric} - {n}-{p}')    
        plt.legend(leg)
    savingDir=f'{plots}robustness/{flags}/nodesize={nodesize}/mmRB/bestfit/'
    if not os.path.exists(savingDir):
        os.makedirs(savingDir)
    plt.savefig(f"robustness_runtimes_enhanced_{flags}.png", dpi=300, bbox_inches='tight')

    # plt.savefig(f'{savingDir}robustness_runtimes_{flags}.pdf')
    plt.close()

def runtimes(dir1='improved/',dir2='baseline/',nodesize=8,flags='allOpts'):
    avg=[]
    columns=[]
    averages=pd.DataFrame(columns=['averages','heuristics','addition'])
    dirs=[main_dir+dir1,main_dir+dir2]
    leg=[]
    avg=[]
    max_size=[16,32,64,128]
    plt.figure(figsize=(5,3))
            
    i=0
    averages=pd.DataFrame(columns=['averages','heuristics'])
    for m in max_size:
        Values = pd.DataFrame(columns=['runtimes'])
        for a in ADDITIONS:

            # for heuristics
            dir=dirs[0]+f'/max_scaling_threshold={m}/'
            fulldir=f'{dir}runtimes/nodesize={nodesize}/addition={a}/mmRB/bestfit/'
            # print(fulldir)
            for files in os.listdir(fulldir):
                new_values = pd.read_csv(f'{fulldir}{files}', header=None, names=['runtimes'])
                Values = pd.concat([Values, new_values], axis=0, ignore_index=True)
                

        leg.append(f'max_scaling_threshold={m}')
        i+=1
            # for baseline
        threshold = np.percentile(Values['runtimes'], 99)
        Values = Values[Values['runtimes'] <= threshold]
        sorted_data = np.sort(Values['runtimes'])*1000
        # max_index = np.argmax(sorted_data)
        # avg_value = np.average(sorted_data)
        # max_value = sorted_data[max_index]
        # # avg_value = sorted_data[avg_index]
        # avg.append(avg_value)
        # Calculate the cumulative probabilities
        cdf = np.arange(1, len(sorted_data) + 1) / len(sorted_data)

        # Plot the CDF
        plt.plot(sorted_data, cdf,markersize=3)
        print(leg)
        
        # colors = plt.cm.rainbow(np.linspace(0, 1, 15))


        plt.grid(True)
        plt.xlabel('Runtime (\u00B5s)')
        plt.ylabel('CDF')
        # plt.title(f'Robustness comparison for {metric} - {n}-{p}')    
        plt.legend(leg)
    savingDir=f'{plots}robustness/{flags}/nodesize={nodesize}/mmRB/bestfit/'
    if not os.path.exists(savingDir):
        os.makedirs(savingDir)
    plt.savefig(f"robustness_runtimes_all_{flags}.png", dpi=300, bbox_inches='tight')

    # plt.savefig(f'{savingDir}robustness_runtimes_{flags}.pdf')
    plt.close()

def time_based_avg(values,times):
    # print("values",values)
    # events=pd.read_csv(f'{events_dir}/events_{addition}.csv')
    events=times
    # add series of values as a new column to events dataframe because events does not have this column
    df = pd.DataFrame({'values': values, 'EventTime': events['EventTime']})
    # print("df",df)
    # df["values"] = pd.to_numeric(df["values"], errors='coerce')
    # Values.fillna(0,inplace=True)
    # df.ffill()#(subset=["values"], inplace=True)
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
    #     plt.savefig(f'{main_dir}/plots/baselines/nodesize={node_size}/addition={addition}/runs_baselines.pdf')
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
                        Times_addr= f'{main_dir}{dir}eventTime/nodesize={node_size}/addition={a}/Max/{p}/'
                # print(fulldir)
                        print(fulldir)
                        for files in os.listdir(fulldir):
                            Values = pd.read_csv(f'{fulldir}{files}', header=None)
                            Times=pd.read_csv(f'{Times_addr}{files}',header=None)
                            Times.columns=['EventTime']
                            Values.columns = [metric]

                            # avg_value = time_based_avg(Values[metric],a)
                            if 'qos' in metric:
                                avg_value = time_based_avg(Values[metric]/100, Times)
                            else:
                                avg_value = time_based_avg(Values[metric], Times)
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
            plt.xlabel('random workload')
            # plt.xticks(range(len(ADDITIONS)), ADDITIONS)
            plt.ylabel(metric)
            # plt.title(f'Robustness comparison for {metric} - {n}-{p}')    
            plt.legend(leg)
            savingDir=f'{plots}/robustness/baselines/nodesize={node_size}/'
            if not os.path.exists(savingDir):
                os.makedirs(savingDir)
            plt.savefig(f'{savingDir}robustness_{metric}.pdf', bbox_inches='tight')
            plt.close()

    return

def robustness(dir1='improved/allOpts',dir2='baseline/',metric='cost',flags='allOpts',nodesize=8):

    avg=[]
    columns=[]
    averages=pd.DataFrame(columns=['averages','heuristics','addition'])
    dirs=[main_dir+dir1,main_dir+dir2]
    leg=[]
    nodeHeus=['mmRB']#,'MmRB','mMRB','MMRB']
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
                    files="HQ.csv"
                    Values = pd.read_csv(f'{fulldir}{files}', header=None)
                    Times=pd.read_csv(f'{Times_addr}{files}',header=None)
                    Times.columns=['EventTime']
                    Values.columns = [metric]

                    if 'qos' in metric:
                        avg_value = time_based_avg(Values[metric]/100, Times)
                    else:
                        avg_value = time_based_avg(Values[metric], Times)
                    averages.loc[i]=[avg_value,files[:-4],a]
                    leg.append(files[:-4])
                    i+=1
                # for baseline
                dir=dirs[1]
                fulldir=f'{dir}{metric}/nodesize={nodesize}/addition={a}/Max/{p}/'
                Times_addr= f'{dir}eventTime/nodesize={nodesize}/addition={a}/Max/{p}/'
                # print(fulldir)
                for files in os.listdir(fulldir):
                    Values = pd.read_csv(f'{fulldir}{files}', header=None)
                    Times=pd.read_csv(f'{Times_addr}{files}',header=None)
                    Times.columns=['EventTime']
                    Values.columns = [metric]

                    if 'qos' in metric:
                        avg_value = time_based_avg(Values[metric]/100, Times)
                    else:
                        avg_value = time_based_avg(Values[metric], Times) 
                    averages.loc[i]=[avg_value,'baseline',a]
                    leg.append('baseline')
                    i+=1
 
            print(leg)
            plt.figure(figsize=figsize)
            # colors = plt.cm.rainbow(np.linspace(0, 1, 15))
            heuristic_averages = averages.groupby('heuristics')['averages'].mean()
            if metric=='cost':
                top_heuristics = heuristic_averages.nsmallest(7).index
            else:
                top_heuristics = heuristic_averages.nlargest(7).index
            for l in leg:
                if 'baseline' in l:
                    b=l
                    break
            if b not in top_heuristics:
                top_heuristics=top_heuristics.append(pd.Index([b]))
            for l in top_heuristics:
                if 'baseline' in l:
                    marker='o'
                    linestyle='--'
                    avgs = averages[averages['heuristics'] == l]
                    avgs = avgs.sort_values(by='addition')
                    plt.plot(avgs['addition'], avgs['averages'], marker=marker,linestyle=linestyle, label=l,linewidth=linewidth,markersize=8,color='black')
                else:
                    marker='x'
                    linestyle='-'
                    avgs = averages[averages['heuristics'] == l]
                    avgs = avgs.sort_values(by='addition')
                    plt.plot(avgs['addition'], avgs['averages'], marker=marker,linestyle=linestyle, label=l,linewidth=linewidth,markersize=8)
            plt.grid(True)
            plt.xlabel('random workload')
            # plt.xticks(range(len(ADDITIONS)), ADDITIONS)
            plt.ylabel(metric)
            # plt.title(f'Robustness comparison for {metric} - {n}-{p}')    
            plt.legend(top_heuristics)
            savingDir=f'{plots}robustness/{flags}/nodesize={nodesize}/{n}/{p}/'
            if not os.path.exists(savingDir):
                os.makedirs(savingDir)
            plt.savefig(f'{savingDir}robustness_{metric}_{flags}.pdf', bbox_inches='tight')
            plt.close()

    return



def robustness_compare_node_core_selection(dir1='improved/allOpts',dir2='baseline',metric='cost',flags='allOpts',nodesize=8):

    avg=[]
    columns=[]
    averages=pd.DataFrame(columns=['averages','heuristics','addition'])
    dirs=[main_dir+dir1,main_dir+dir2]
    legend=[]
    max_size=[16,32,64]
    #cost
    nodeHeus=['MMRB','mmRB','MmRB','mMRB']
    
    for m in max_size:
        plt.figure(figsize=figsize)
        baselineadded=False
        for n in nodeHeus:
            for p in partitionHeus:
                avg=[]
                
                i=0
                averages=pd.DataFrame(columns=['averages','heuristics','addition'])
                for a in ADDITIONS:
                    # for heuristics
                    leg=[]
                    dir=dirs[0]+f'/max_scaling_threshold={m}/'
                    fulldir=f'{dir}{metric}/nodesize={nodesize}/addition={a}/{n}/{p}/'
                    Times_addr= f'{dir}eventTime/nodesize={nodesize}/addition={a}/{n}/{p}/'
                    # print(fulldir)
                    for files in os.listdir(fulldir):
                        Values = pd.read_csv(f'{fulldir}{files}', header=None)
                        Times=pd.read_csv(f'{Times_addr}{files}',header=None)
                        Times.columns=['EventTime']
                        Values.columns = [metric]

                        if 'qos' in metric:
                            avg_value = time_based_avg(Values[metric]/100, Times)
                        else:
                            avg_value = time_based_avg(Values[metric], Times)
                        averages.loc[i]=[avg_value,f'{n}_{p}',a]
                        leg.append(f'{n}_{p}')
                        i+=1
                    # for baseline
                    if not baselineadded:
                        dir=dirs[1]+f'/{m}/'
                        fulldir=f'{dir}{metric}/nodesize={nodesize}/addition={a}/None/bestfit/'
                        Times_addr= f'{dir}eventTime/nodesize={nodesize}/addition={a}/None/bestfit/'
                        # print(fulldir)
                        for files in os.listdir(fulldir):
                            Values = pd.read_csv(f'{fulldir}{files}', header=None)
                            Times=pd.read_csv(f'{Times_addr}{files}',header=None)
                            Times.columns=['EventTime']
                            Values.columns = [metric]

                            if 'qos' in metric:
                                avg_value = time_based_avg(Values[metric]/100, Times)
                            else:
                                avg_value = time_based_avg(Values[metric], Times)
                            averages.loc[i]=[avg_value,'baseline',a]
                            leg.append('baseline')
                            i+=1
    
                # print(leg)
                baselineadded=True
                
                # colors = plt.cm.rainbow(np.linspace(0, 1, 15))
                if averages.empty:
                    continue
                heuristic_averages = averages.groupby('heuristics')['averages'].mean()
                if metric=='cost':
                    top_heuristics = heuristic_averages.nsmallest(4).index
                else:
                    top_heuristics = heuristic_averages.nlargest(4).index
                for l in leg:
                    if 'baseline' in l:
                        b=l
                        if b not in top_heuristics:
                            top_heuristics=top_heuristics.append(pd.Index([b]))
                        break
                
                for l in top_heuristics:
                    avgs = averages[averages['heuristics'] == l]
                    avgs = avgs.sort_values(by='addition')
                    legend.append(l)
                    if 'baseline' in l:
                        marker='x'
                        linestyle='-'
                        plt.plot(avgs['addition'], avgs['averages'], marker=marker,linestyle=linestyle, label=l,linewidth=linewidth,markersize=8,color='black')
                    else:
                        marker='o'
                        linestyle='--'
                        plt.plot(avgs['addition'], avgs['averages'], marker=marker,linestyle=linestyle, label=l,linewidth=linewidth,markersize=8)
                    
                    # plt.plot(avgs['addition'], avgs['averages'], marker=marker,linestyle=linestyle, label=l,linewidth=linewidth,markersize=8)
        plt.grid(True)
        plt.xlabel('random workload',fontsize=fontsizetxt)
        plt.xticks(fontsize=fontsize)
        plt.yticks(fontsize=fontsize)
        plt.ylabel(metric,fontsize=fontsizetxt)
        plt.legend(legend,fontsize=fontsizetxt)
        savingDir=f'{plots}robustness/{flags}/{m}/nodesize={nodesize}/'
        if not os.path.exists(savingDir):
            os.makedirs(savingDir)
        plt.savefig(f'{savingDir}robustness_{metric}_all_nodeheus.pdf', bbox_inches='tight')
        plt.close()


def robustness_compare_nodesize(dir1='improved/allOpts',dir2='baseline',metric='cost',flags='allOpts'):

    avg=[]
    columns=[]
    averages=pd.DataFrame(columns=['averages','heuristics','addition'])
    dirs=[main_dir+dir1,main_dir+dir2]
    legend=[]
    max_size=[16,32,64,128]
    nodesizes=[8,16]
    #cost
    
    for m in max_size:    
        plt.figure(figsize=figsize)
        for nodesize in nodesizes:
        
            avg=[]
            
            i=0
            averages=pd.DataFrame(columns=['averages','heuristics','addition'])
            for a in ADDITIONS:
                # for heuristics
                leg=[]
                dir=dirs[0]+f'/max_scaling_threshold={m}/'
                fulldir=f'{dir}{metric}/nodesize={nodesize}/addition={a}/mmRB/bestfit/'
                Times_addr= f'{dir}eventTime/nodesize={nodesize}/addition={a}/mmRB/bestfit/'
                # print(fulldir)
                for files in os.listdir(fulldir):
                    Values = pd.read_csv(f'{fulldir}{files}', header=None)
                    Times=pd.read_csv(f'{Times_addr}{files}',header=None)
                    Times.columns=['EventTime']
                    Values.columns = [metric]

                    if 'qos' in metric:
                        avg_value = time_based_avg(Values[metric]/100, Times)
                    else:
                        avg_value = time_based_avg(Values[metric], Times)
                    averages.loc[i]=[avg_value,f'node_size={nodesize}',a]
                    leg.append(f'node_size={nodesize}')
                    i+=1
                # for baseline
                dir=dirs[1]+f'/{m}/'
                fulldir=f'{dir}{metric}/nodesize={nodesize}/addition={a}/None/bestfit/'
                Times_addr= f'{dir}eventTime/nodesize={nodesize}/addition={a}/None/bestfit/'
                # print(fulldir)
                for files in os.listdir(fulldir):
                    Values = pd.read_csv(f'{fulldir}{files}', header=None)
                    Times=pd.read_csv(f'{Times_addr}{files}',header=None)
                    Times.columns=['EventTime']
                    Values.columns = [metric]

                    if 'qos' in metric:
                        avg_value = time_based_avg(Values[metric]/100, Times)
                    else:
                        avg_value = time_based_avg(Values[metric], Times)
                    averages.loc[i]=[avg_value,f'baseline_nodesize={nodesize}',a]
                    leg.append(f'baseline_nodesize={nodesize}')
                    i+=1

            # print(leg)
            
            # colors = plt.cm.rainbow(np.linspace(0, 1, 15))
            if averages.empty:
                continue
            heuristic_averages = averages.groupby('heuristics')['averages'].mean()
            if metric=='cost':
                top_heuristics = heuristic_averages.nsmallest(4).index
            else:
                top_heuristics = heuristic_averages.nlargest(4).index
            for l in leg:
                if 'baseline' in l:
                    b=l
                    if b not in top_heuristics:
                        top_heuristics=top_heuristics.append(pd.Index([b]))
                    break
            
            for l in top_heuristics:
                avgs = averages[averages['heuristics'] == l]
                avgs = avgs.sort_values(by='addition')
                legend.append(l)
                if 'baseline' in l:
                    marker='x'
                    linestyle='-'
                    plt.plot(ADDITIONS, avgs['averages'], marker=marker,linestyle=linestyle, label=l,linewidth=linewidth,markersize=8)
                else:
                    marker='o'
                    linestyle='--'
                    plt.plot(ADDITIONS, avgs['averages'], marker=marker,linestyle=linestyle, label=l,linewidth=linewidth,markersize=8)
                
                # plt.plot(avgs['addition'], avgs['averages'], marker=marker,linestyle=linestyle, label=l,linewidth=linewidth,markersize=8)
        plt.grid(True)
        plt.xlabel('random workload',fontsize=fontsizetxt)
        plt.xticks(fontsize=fontsize)
        plt.yticks(fontsize=fontsize)
        plt.ylabel(metric,fontsize=fontsizetxt)
        plt.legend(legend,fontsize=fontsizetxt)
        savingDir=f'{plots}robustness/{flags}/{m}/'
        if not os.path.exists(savingDir):
            os.makedirs(savingDir)
        plt.savefig(f'{savingDir}robustness_{metric}_all_nodesizes.pdf', bbox_inches='tight')
        print(f'{savingDir}robustness_{metric}_all_nodesizes.pdf')
        plt.close()

def compute_cost(qospercost:pd.DataFrame,qos:pd.DataFrame):
    cost=qos/qospercost
    cost=cost.fillna(0)
    return cost
def robustness_max_scaling_size(dir1='improved/allOpts',dir2='baseline',metric='cost',flags='allOpts',nodesize=8):
    nodeHeus=['mmRB']#,'MmRB','mMRB','MMRB']
    avg=[]
    columns=[]
    averages=pd.DataFrame(columns=['averages','heuristics','addition'])
    dirs=[main_dir+dir1,main_dir+dir2]
    leg=[]
    max_size=[16,32,64,128]#[512,128,200,96,256,64,32,16]
    #cost
    valid_heu=["LBCI","LRED"]

    for m in max_size:
        for n in nodeHeus:
            for p in partitionHeus:
                avg=[]
                
                i=0
                averages=pd.DataFrame(columns=['averages','heuristics','addition'])
                for a in ADDITIONS:
                    leg=[]
                    # for heuristics
                    dir=dirs[0]+f'/max_scaling_threshold={m}'
                    fulldir=f'{dir}/{metric}/nodesize={nodesize}/addition={a}/{n}/{p}/'
                    Times_addr= f'{dir}/eventTime/nodesize={nodesize}/addition={a}/{n}/{p}/'
                    # print(fulldir)
                    for files in os.listdir(fulldir):
                        if "intra_node_reall" in flags:
                            if "LB"  not in files and "LC" not in files and "LBC" not in files and "HC" not in files:
                                continue
                        if "reduced" in flags:
                            if "LIHBC"  not in files and "LRED" not in files and "LI" not in files:
                                continue
                        
                        Values = pd.read_csv(f'{fulldir}{files}', header=None,dtype={f'{metric}': float})
                        # qos = pd.read_csv(f'{dir}/qos/nodesize={nodesize}/addition={a}/{n}/{p}/{files}', header=None)
                        # qospercost = pd.read_csv(f'{dir}/qosPerCost/nodesize={nodesize}/addition={a}/{n}/{p}/{files}', header=None)
                        # if 'cost' in metric:
                        #     Values=compute_cost(qospercost=qospercost,qos=qos)
                        Times=pd.read_csv(f'{Times_addr}{files}',header=None)
                        Times.columns=['EventTime']
                        Values.columns = [metric]
                        # Values=Values.ffill()
                        
                        # Values[metric] = pd.to_numeric(Values[metric], errors='coerce')
                        # Values.fillna(0,inplace=True)

                        # print(fulldir)

                        
                        avg_value = time_based_avg(Values[metric], Times)
                        if 'qos' in metric:
                            avg_value=avg_value/100
                        # print(flags)
                        if 'with' in flags:
                            l=files[:-4]
                        else:
                            l='enhanced baseline'
                        averages.loc[i]=[avg_value,l,a]
                        leg.append(l)
                        i+=1
                    # for baseline
                    dir=dirs[1]+f'/{m}'
                    fulldir=f'{dir}/{metric}/nodesize={nodesize}/addition={a}/None/{p}/'
                    Times_addr= f'{dir}/eventTime/nodesize={nodesize}/addition={a}/None/{p}/'
                    # print(fulldir)
                    for files in os.listdir(fulldir):
                        Values = pd.read_csv(f'{fulldir}{files}', header=None)
                        Times=pd.read_csv(f'{Times_addr}{files}',header=None)
                        Times.columns=['EventTime']
                        Values.columns = [metric]

                        avg_value = time_based_avg(Values[metric], Times)
                        if 'qos' in metric:
                            avg_value = avg_value/100
                            
                        averages.loc[i]=[avg_value,'baseline',a]
                        leg.append('baseline')
                        i+=1
    
                print(leg)
                print(fulldir)
                plt.figure(figsize=figsize)
                # colors = plt.cm.rainbow(np.linspace(0, 1, 15))
                heuristic_averages = averages.groupby('heuristics')['averages'].mean()
                if metric=='cost':
                    top_heuristics = heuristic_averages.nsmallest(4).index
                else:
                    top_heuristics = heuristic_averages.nlargest(4).index
                for l in leg:
                    if 'baseline' in l and not 'enhanced' in l:
                        b=l
                        break
                if b not in top_heuristics:
                    top_heuristics=top_heuristics.append(pd.Index([b]))
                for l in top_heuristics:
                    if 'baseline' in l and not 'enhanced' in l:
                        marker='x'
                        linestyle='-'
                        avgs = averages[averages['heuristics'] == l]
                        avgs = avgs.sort_values(by='addition')
                        plt.plot(avgs['addition'], avgs['averages'], marker=marker,linestyle=linestyle, label=l,linewidth=linewidth,markersize=8,color='black')
                    else:
                        marker='o'
                        linestyle='--'
                        avgs = averages[averages['heuristics'] == l]
                        avgs = avgs.sort_values(by='addition')
                        plt.plot(avgs['addition'], avgs['averages'], marker=marker,linestyle=linestyle, label=l,linewidth=linewidth,markersize=8)
                plt.grid(True)
                plt.xlabel('random workload',fontsize=fontsizetxt)
                plt.ylabel(metric,fontsize=fontsizetxt)   
                plt.xticks(fontsize=fontsize)
                plt.yticks(fontsize=fontsize)
                # plt.title(f'Robustness comparison for {metric} - {n}-{p}')    
                plt.legend(top_heuristics,fontsize=fontsizetxt)
                savingDir=f'{plots}robustness/{flags}/{m}/nodesize={nodesize}/{n}/{p}/'
                if not os.path.exists(savingDir):
                    os.makedirs(savingDir)
                plt.savefig(f'{savingDir}robustness_{metric}_{flags}.pdf', bbox_inches='tight')
                plt.close()

    return


import os
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D
from scipy.interpolate import griddata  # For surface interpolation

def robustness_max_scaling_size_3d_sheets(dir1='improved/allOpts', dir2='baseline', metric='cost', flags='allOpts', nodesize=8):
    dirs = [main_dir + dir1, main_dir + dir2]
    max_size = np.array([ 16,32,64,128])  # Different max scaling sizes
    nodeHeus=['MMRB','mmRB','MmRB','mMRB']
    for n in nodeHeus:
        for p in partitionHeus:
            all_averages = pd.DataFrame(columns=['averages', 'heuristics', 'addition', 'max_size'])

            for m in max_size:
                i = 0
                averages = pd.DataFrame(columns=['averages', 'heuristics', 'addition'])
                
                for a in ADDITIONS:
                    leg = []
                    
                    # Process the improved heuristics
                    dir = dirs[0] + f'/max_scaling_threshold={m}'
                    fulldir = f'{dir}/{metric}/nodesize={nodesize}/addition={a}/{n}/{p}/'
                    Times_addr = f'{dir}/eventTime/nodesize={nodesize}/addition={a}/{n}/{p}/'
                    
                    for files in os.listdir(fulldir):
                        Values = pd.read_csv(f'{fulldir}{files}', header=None)
                        Times = pd.read_csv(f'{Times_addr}{files}', header=None)
                        Times.columns = ['EventTime']
                        Values.columns = [metric]

                        if 'qos' in metric:
                            avg_value = time_based_avg(Values[metric]/100, Times)
                        else:
                            avg_value = time_based_avg(Values[metric], Times)
                        heuristic_label = 'improved' if 'all' in flags else files[:-4]
                        averages.loc[i] = [avg_value, heuristic_label, a]
                        leg.append(heuristic_label)
                        i += 1

                    # Process the baseline heuristics
                    dir = dirs[1] + f'/{m}'
                    fulldir = f'{dir}/{metric}/nodesize={nodesize}/addition={a}/None/{p}/'
                    Times_addr = f'{dir}/eventTime/nodesize={nodesize}/addition={a}/None/{p}/'

                    for files in os.listdir(fulldir):
                        Values = pd.read_csv(f'{fulldir}{files}', header=None)
                        Times = pd.read_csv(f'{Times_addr}{files}', header=None)
                        Times.columns = ['EventTime']
                        Values.columns = [metric]

                        if 'qos' in metric:
                            avg_value = time_based_avg(Values[metric]/100, Times)
                        else:
                            avg_value = time_based_avg(Values[metric], Times)
                        averages.loc[i] = [avg_value, 'baseline', a]
                        leg.append('baseline')
                        i += 1
                
                # Store all values including max_size for 3D plotting
                averages['max_size'] = m
                all_averages = pd.concat([all_averages, averages], ignore_index=True)

            # Identify top heuristics
            heuristic_averages = all_averages.groupby('heuristics')['averages'].mean()
            if metric == 'cost':
                top_heuristics = heuristic_averages.nsmallest(4).index
            else:
                top_heuristics = heuristic_averages.nlargest(4).index
            
            # Ensure baseline is included in the plot
            baseline_label = 'baseline'
            if baseline_label not in top_heuristics:
                top_heuristics = top_heuristics.append(pd.Index([baseline_label]))

            # Create a new 3D figure for each `(n, p)` combination
            fig = plt.figure(figsize=(6,5))
            ax = fig.add_subplot(111, projection='3d')

            # Distinguishable colors for each heuristic
            colormaps = ['viridis', 'coolwarm', 'plasma', 'cividis', 'spring']
            color_idx = 0

            # Plot separate surfaces for each heuristic
            for l in top_heuristics:
                X_vals = []
                Y_vals = []
                Z_vals = []

                # Collect data for the current heuristic
                for m in max_size:
                    avgs = all_averages[(all_averages['heuristics'] == l) & (all_averages['max_size'] == m)]
                    if avgs.empty:
                        continue
                    
                    avgs = avgs.sort_values(by='addition')
                    X_vals.extend(avgs['addition'].values)
                    Y_vals.extend([m] * len(avgs))  # Y-axis is max_size
                    Z_vals.extend(avgs['averages'].values)

                if len(X_vals) < 3:
                    continue  # Skip heuristics with too few data points for interpolation

                # Convert to NumPy arrays
                X_vals = np.array(X_vals)
                Y_vals = np.array(Y_vals)
                Z_vals = np.array(Z_vals)

                # Create a grid for smooth interpolation
                grid_x, grid_y = np.meshgrid(np.linspace(X_vals.min(), X_vals.max(), 30),
                                             np.linspace(Y_vals.min(), Y_vals.max(), 30))

                # Interpolate the data for the surface
                grid_z = griddata((X_vals, Y_vals), Z_vals, (grid_x, grid_y), method='cubic')

                # Select a color map
                cmap = colormaps[color_idx % len(colormaps)]
                color_idx += 1  # Cycle through color maps

                # Plot the heuristic surface
                ax.plot_surface(grid_x, grid_y, grid_z, cmap=cmap, alpha=0.75, edgecolor='k', label=l)

                # Add legend manually
                if 'cost' in metric:
                    if l=='baseline':
                        # ax.plot_wireframe(grid_x, grid_y, grid_z, color='black', linewidth=1.5)  # Wireframe for baseline
                        t=1200
                    else:
                        t=1000
                    ax.text(X_vals.max()-0.3, Y_vals.mean(), t , l, color='black', fontsize=10, fontweight='bold', bbox=dict(facecolor='white', alpha=0.7))
                if 'qos' in metric:
                    if l=='baseline':
                        # ax.plot_wireframe(grid_x, grid_y, grid_z, color='black', linewidth=1.5)  # Wireframe for baseline
                        t=200
                    else:
                        t=250
                    ax.text(X_vals.max(), Y_vals.min(), t , l, color='black', fontsize=10, fontweight='bold', bbox=dict(facecolor='white', alpha=0.7))
                if 'qosPerCost' in metric:
                    if l=='baseline':
                        # ax.plot_wireframe(grid_x, grid_y, grid_z, color='black', linewidth=1.5)  # Wireframe for baseline
                        t=0.09
                    else:
                        t=0.12
                    ax.text(X_vals.max(), Y_vals.min(), t , l, color='black', fontsize=10, fontweight='bold', bbox=dict(facecolor='white', alpha=0.7))
                # ax.text(X_vals.mean(), Y_vals.mean(), Z_vals.max(), l, color='black', fontsize=10, fontweight='bold')
                

            # Set 3D plot labels and grid
            ax.set_xlabel('random workload',labelpad=10,fontsize=fontsizetxt)
            ax.set_ylabel('Max Scaling Size',fontsize=fontsizetxt)
            ax.set_zlabel(metric,fontsize=fontsizetxt)
            
            if 'qos' in metric:
                ax.view_init( azim=150)
            if 'qosPer' in metric:
                ax.view_init( azim=60) 
                for tick in ax.get_xticklabels():
                    tick.set_rotation(60) 
            # ax.set_title(f'3D Robustness Comparison for {metric}\nNode: {n}, Partition: {p}')
            ax.grid(True)
            
            
            if 'cost' in metric:
                ax.view_init(elev=22, azim=120) 
                for tick in ax.get_xticklabels():
                    tick.set_rotation(60) 

            # xticks = ax.get_xticks()
            # ax.set_xticks(xticks)
            # ax.set_xticklabels([f'{tick:.0f}' for tick in xticks])

            # Save the plot in a separate folder for each `(n, p)`
            savingDir = f'{plots}robustness/{flags}/nodesize={nodesize}/{n}/{p}/3D/'
            if not os.path.exists(savingDir):
                os.makedirs(savingDir)
            
            plt.savefig(f'{savingDir}robustness_{metric}_3D_{flags}_sheets.pdf', bbox_inches='tight',pad_inches=0.35)
            plt.close()  # Clear the plot for the next iteration


def robustness_interval_length_3d_sheets(dir1='improved/allOpts', dir2='baseline', metric='cost', flags='allOpts', nodesize=8):
    dirs = [main_dir + dir1, main_dir + dir2]
    interval_length = np.array([ 5,10,20,30,40])  # Different max scaling sizes
    nodeHeus=['mmRB']
    m=16
    for n in nodeHeus:
        for p in partitionHeus:
            all_averages = pd.DataFrame(columns=['averages', 'heuristics', 'addition', 'interval_length'])

            for intl in interval_length:
                i = 0
                averages = pd.DataFrame(columns=['averages', 'heuristics', 'addition'])
                
                for a in ADDITIONS:
                    leg = []
                    
                    # Process the improved heuristics
                    dir = dirs[0] + f'/interval={intl}/max_scaling_threshold={m}'
                    fulldir = f'{dir}/{metric}/nodesize={nodesize}/addition={a}/{n}/{p}/'
                    Times_addr = f'{dir}/eventTime/nodesize={nodesize}/addition={a}/{n}/{p}/'
                    
                    for files in os.listdir(fulldir):
                        Values = pd.read_csv(f'{fulldir}{files}', header=None)
                        Times = pd.read_csv(f'{Times_addr}{files}', header=None)
                        Times.columns = ['EventTime']
                        Values.columns = [metric]

                        if 'qos' in metric:
                            avg_value = time_based_avg(Values[metric]/100, Times)
                        else:
                            avg_value = time_based_avg(Values[metric], Times)
                        heuristic_label = 'improved' if 'all' in flags else files[:-4]
                        averages.loc[i] = [avg_value, heuristic_label, a]
                        leg.append(heuristic_label)
                        i += 1

                    # Process the baseline heuristics
                    dir = dirs[1] + f'/{m}'
                    fulldir = f'{dir}/{metric}/nodesize={nodesize}/addition={a}/None/{p}/'
                    Times_addr = f'{dir}/eventTime/nodesize={nodesize}/addition={a}/None/{p}/'

                    for files in os.listdir(fulldir):
                        Values = pd.read_csv(f'{fulldir}{files}', header=None)
                        Times = pd.read_csv(f'{Times_addr}{files}', header=None)
                        Times.columns = ['EventTime']
                        Values.columns = [metric]

                        if 'qos' in metric:
                            avg_value = time_based_avg(Values[metric]/100, Times)
                        else:
                            avg_value = time_based_avg(Values[metric], Times)
                        averages.loc[i] = [avg_value, 'baseline', a]
                        leg.append('baseline')
                        i += 1
                
                # Store all values including max_size for 3D plotting
                averages['interval_length'] = intl
                all_averages = pd.concat([all_averages, averages], ignore_index=True)
            
            if 'cost' in metric:
                print(all_averages)
            # Identify top heuristics
            heuristic_averages = all_averages.groupby('heuristics')['averages'].mean()
            if metric == 'cost':
                top_heuristics = heuristic_averages.nsmallest(4).index
            else:
                top_heuristics = heuristic_averages.nlargest(4).index
            
            # Ensure baseline is included in the plot
            baseline_label = 'baseline'
            if baseline_label not in top_heuristics:
                top_heuristics = top_heuristics.append(pd.Index([baseline_label]))

            # Create a new 3D figure for each `(n, p)` combination
            fig = plt.figure(figsize=(6,5))
            ax = fig.add_subplot(111, projection='3d')

            # Distinguishable colors for each heuristic
            colormaps = ['viridis', 'coolwarm', 'plasma', 'cividis', 'spring']
            color_idx = 0

            # Plot separate surfaces for each heuristic
            for l in top_heuristics:
                X_vals = []
                Y_vals = []
                Z_vals = []

                # Collect data for the current heuristic
                for intl in interval_length:
                    avgs = all_averages[(all_averages['heuristics'] == l) & (all_averages['interval_length'] == intl)]
                    if avgs.empty:
                        continue
                    
                    avgs = avgs.sort_values(by='addition')
                    X_vals.extend(avgs['addition'].values)
                    Y_vals.extend([intl] * len(avgs))  # Y-axis is max_size
                    Z_vals.extend(avgs['averages'].values)

                if len(X_vals) < 3:
                    continue  # Skip heuristics with too few data points for interpolation

                # Convert to NumPy arrays
                X_vals = np.array(X_vals)
                Y_vals = np.array(Y_vals)
                Z_vals = np.array(Z_vals)

                # Create a grid for smooth interpolation
                grid_x, grid_y = np.meshgrid(np.linspace(X_vals.min(), X_vals.max(), 30),
                                             np.linspace(Y_vals.min(), Y_vals.max(), 30))

                # Interpolate the data for the surface
                grid_z = griddata((X_vals, Y_vals), Z_vals, (grid_x, grid_y), method='cubic')

                # Select a color map
                cmap = colormaps[color_idx % len(colormaps)]
                color_idx += 1  # Cycle through color maps

                # Plot the heuristic surface
                ax.plot_surface(grid_x, grid_y, grid_z, cmap=cmap, alpha=0.75, edgecolor='k', label=l)

                # Add legend manually
                if 'cost' in metric:
                    if l=='baseline':
                        # ax.plot_wireframe(grid_x, grid_y, grid_z, color='black', linewidth=1.5)  # Wireframe for baseline
                        t=175
                    else:
                        t=170
                    ax.text(X_vals.min(), Y_vals.min(), t , l, color='black', fontsize=10, fontweight='bold', bbox=dict(facecolor='white', alpha=0.7))
                if 'qos' in metric:
                    if l=='baseline':
                        # ax.plot_wireframe(grid_x, grid_y, grid_z, color='black', linewidth=1.5)  # Wireframe for baseline
                        t=54
                    else:
                        t=57
                    ax.text(X_vals.min(), Y_vals.max(), t , l, color='black', fontsize=10, fontweight='bold', bbox=dict(facecolor='white', alpha=0.7))
                if 'qosPerCost' in metric:
                    if l=='baseline':
                        # ax.plot_wireframe(grid_x, grid_y, grid_z, color='black', linewidth=1.5)  # Wireframe for baseline
                        t=0.3
                    else:
                        t=0.32
                    ax.text(X_vals.min(), Y_vals.max(), t , l, color='black', fontsize=10, fontweight='bold', bbox=dict(facecolor='white', alpha=0.7))
                # ax.text(X_vals.mean(), Y_vals.mean(), Z_vals.max(), l, color='black', fontsize=10, fontweight='bold')
                

            # Set 3D plot labels and grid
            ax.set_xlabel('random workload',labelpad=10,fontsize=fontsizetxt)
            ax.set_ylabel('Interval Length',fontsize=fontsizetxt)
            ax.set_zlabel(metric,fontsize=fontsizetxt)
            # if 'cost' in metric:
            #     ax.view_init(elev=22, azim=120) 
            # if 'qos' in metric:
            #     ax.view_init( azim=150)
            # if 'qosPer' in metric:
            #     ax.view_init( azim=60) 
            # ax.set_title(f'3D Robustness Comparison for {metric}\nNode: {n}, Partition: {p}')
            ax.grid(True)
            for tick in ax.get_xticklabels():
                tick.set_rotation(30)  # Or try 45 for steeper angle
                # tick.set_horizontalalignment('right') 


            # xticks = ax.get_xticks()
            # ax.set_xticks(xticks)
            # ax.set_xticklabels([f'{tick:.0f}' for tick in xticks])

            # Save the plot in a separate folder for each `(n, p)`
            savingDir = f'{plots}robustness/all_interval_based/nodesize={nodesize}/{n}/{p}/3D/'
            if not os.path.exists(savingDir):
                os.makedirs(savingDir)
            
            plt.savefig(f'{savingDir}robustness_{metric}_3D_{flags}_sheets.pdf', bbox_inches='tight',pad_inches=0.35)
            plt.close()  # Clear the plot for the next iteration

def robustness_max_scaling_size_interval_based(dir1='improved/allOpts',dir2='baseline',metric='cost',flags='allOpts',nodesize=8):
    nodeHeus=['mmRB']#,'MmRB','mMRB','MMRB']
    avg=[]
    columns=[]
    averages=pd.DataFrame(columns=['averages','heuristics','addition'])
    dirs=[main_dir+dir1,main_dir+dir2]
    leg=[]
    max_size=[16,32,64,128]#[512,128,200,96,256,64,32,16]
    #cost
    valid_heu=["LBCI","LRED"]

    for m in max_size:
        for n in nodeHeus:
            for p in partitionHeus:
                avg=[]
                
                i=0
                averages=pd.DataFrame(columns=['averages','heuristics','addition'])
                for a in ADDITIONS:
                    leg=[]
                    # for heuristics
                    dir=dirs[0]+f'/max_scaling_threshold={m}'
                    fulldir=f'{dir}/{metric}/nodesize={nodesize}/addition={a}/{n}/{p}/'
                    Times_addr= f'{dir}/eventTime/nodesize={nodesize}/addition={a}/{n}/{p}/'
                    # print(fulldir)
                    for files in os.listdir(fulldir):
                        # if "intra_node_reall" in flags:
                        #     if "LB"  not in files and "LC" not in files and "LBC" not in files and "HC" not in files:
                        #         continue
                        # if "reduced" in flags:
                        #     if "LBCI"  not in files and "LRED" not in files:
                        #         continue
                        
                        Values = pd.read_csv(f'{fulldir}{files}', header=None,dtype={f'{metric}': float})
                        # qos = pd.read_csv(f'{dir}/qos/nodesize={nodesize}/addition={a}/{n}/{p}/{files}', header=None)
                        # qospercost = pd.read_csv(f'{dir}/qosPerCost/nodesize={nodesize}/addition={a}/{n}/{p}/{files}', header=None)
                        # if 'cost' in metric:
                        #     Values=compute_cost(qospercost=qospercost,qos=qos)
                        Times=pd.read_csv(f'{Times_addr}{files}',header=None)
                        Times.columns=['EventTime']
                        Values.columns = [metric]
                        # Values=Values.ffill()
                        
                        # Values[metric] = pd.to_numeric(Values[metric], errors='coerce')
                        # Values.fillna(0,inplace=True)

                        # print(fulldir)

                        
                        avg_value = time_based_avg(Values[metric], Times)
                        if 'qos' in metric:
                            avg_value=avg_value/100
                        # print(flags)
                        l='interval-based'
                        averages.loc[i]=[avg_value,l,a]
                        leg.append(l)
                        i+=1
                    # for baseline
                    # if not baselineprocessed:
                    
                    dir=dirs[1]+f'/{m}'
                    fulldir=f'{dir}/{metric}/nodesize={nodesize}/addition={a}/None/{p}/'
                    Times_addr= f'{dir}/eventTime/nodesize={nodesize}/addition={a}/None/{p}/'
                    # print(fulldir)
                    for files in os.listdir(fulldir):
                        Values = pd.read_csv(f'{fulldir}{files}', header=None)
                        Times=pd.read_csv(f'{Times_addr}{files}',header=None)
                        Times.columns=['EventTime']
                        Values.columns = [metric]

                        avg_value = time_based_avg(Values[metric], Times)
                        if 'qos' in metric:
                            avg_value = avg_value/100
                            
                        averages.loc[i]=[avg_value,'baseline',a]
                        leg.append('baseline')
                        i+=1
                    dir=f'{main_dir}improved/none_interval_based/max_scaling_threshold={m}'
                    fulldir=f'{dir}/{metric}/nodesize={nodesize}/addition={a}/{n}/{p}/'
                    Times_addr= f'{dir}/eventTime/nodesize={nodesize}/addition={a}/{n}/{p}/'
                    # print(fulldir)
                    for files in os.listdir(fulldir):
                        Values = pd.read_csv(f'{fulldir}{files}', header=None)
                        Times=pd.read_csv(f'{Times_addr}{files}',header=None)
                        Times.columns=['EventTime']
                        Values.columns = [metric]

                        avg_value = time_based_avg(Values[metric], Times)
                        if 'qos' in metric:
                            avg_value = avg_value/100
                            
                        averages.loc[i]=[avg_value,'enhanced baseline',a]
                        leg.append('enhanced baseline')
                        i+=1
    
                print(leg)
                print(fulldir)
                plt.figure(figsize=figsize)
                # colors = plt.cm.rainbow(np.linspace(0, 1, 15))
                heuristic_averages = averages.groupby('heuristics')['averages'].mean()
                if metric=='cost':
                    top_heuristics = heuristic_averages.nsmallest(4).index
                else:
                    top_heuristics = heuristic_averages.nlargest(4).index
                # for l in leg:
                #     if 'baseline' in l:
                #         b=l
                #         break
                # if b not in top_heuristics:
                #     top_heuristics=top_heuristics.append(pd.Index([b]))
                for l in top_heuristics:
                    if 'baseline' in l and 'enhanced baseline' not in l:
                        marker='x'
                        linestyle='-'
                        avgs = averages[averages['heuristics'] == l]
                        avgs = avgs.sort_values(by='addition')
                        plt.plot(avgs['addition'], avgs['averages'], marker=marker,linestyle=linestyle, label=l,linewidth=linewidth,markersize=8,color='black')
                    elif 'enhanced baseline' in l:
                        marker='x'
                        linestyle='-'
                        avgs = averages[averages['heuristics'] == l]
                        avgs = avgs.sort_values(by='addition')
                        plt.plot(avgs['addition'], avgs['averages'], marker=marker,linestyle=linestyle, label=l,linewidth=linewidth,markersize=8)
                    else:
                        marker='o'
                        linestyle='--'
                        avgs = averages[averages['heuristics'] == l]
                        avgs = avgs.sort_values(by='addition')
                        plt.plot(avgs['addition'], avgs['averages'], marker=marker,linestyle=linestyle, label=l,linewidth=linewidth,markersize=8)
                plt.grid(True)
                plt.xlabel('random workload',fontsize=fontsizetxt)
                plt.ylabel(metric,fontsize=fontsizetxt)   
                plt.xticks(fontsize=fontsize)
                plt.yticks(fontsize=fontsize)
                plt.legend(top_heuristics,fontsize=fontsizetxt)
                savingDir=f'{plots}robustness/{flags}/{m}/nodesize={nodesize}/{n}/{p}/'
                if not os.path.exists(savingDir):
                    os.makedirs(savingDir)
                plt.savefig(f'{savingDir}robustness_{metric}_interval_based.pdf', bbox_inches='tight')
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
    
    # print(avgr)
    # avgqpc=processfiles(dir1=dir1,dir2=dir2,metric='qosPerCost')
    # print(avgqpc)
    # avgqos=processfiles(dir1=dir1,dir2=dir2,metric='qos')
    # print(avgqos)
    # # avgcost=cost(dir1=dir1,dir2=dir2)
    # avgcost=processfiles(dir1=dir1,dir2=dir2,metric='cost')
    # print(avgcost)  
    nodeHeus=['Max','Min']
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
#             plt.savefig(f'{savingDir}robustness_{metric}.pdf')
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
#     plt.savefig(f'{savingDir}{nodeHeu}_{partitionHeu}.pdf')
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
#     plt.savefig(f'{savingDir}{nodeHeu}_{partitionHeu}.pdf')
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
#     plt.savefig(f'{savingDir}{nodeHeu}_{partitionHeu}.pdf')
#     averages.to_csv(f'{savingDir}{nodeHeu}_{partitionHeu}_averages.csv',index=False)
#     plt.close()
#     return averages


