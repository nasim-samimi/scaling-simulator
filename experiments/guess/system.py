import os
import pandas as pd
from users import *
from services import *
import math
import random
# we shouldn't be limited by the number of nodes in the cloud or the domain
# so we choose high numbers

NUM_CORES_PER_NODE=8

main_dir='data/'
PARTITIONING_H=['bestfit','worstfit']

REALLOCATION_H=["HBI","HCI","HBCI","HBIcC"]
NODE_SELECTION_H=["MinMin","MaxMax"]

MAX_BANDWIDTH_PER_CORE=95


def Heuristics(reallocationH,nodeSelectionH):
    if not os.path.exists(main_dir+'heuristics/'):
        os.mkdir(main_dir+'heuristics/')
    for nH in nodeSelectionH:
        for rH in reallocationH:
            for pH in PARTITIONING_H:
                df=pd.DataFrame(columns=['ReallocationHeuristic','NodeSelectionHeuristic','PartitioningHeuristic'])
                df['ReallocationHeuristic']=[rH]
                df['NodeSelectionHeuristic']=[nH]
                df['PartitioningHeuristic']=[pH]
                df.to_csv(main_dir+f'heuristics/{rH}_{nH}_{pH}.csv', index=False)

def computeNodeCores(d,opt):
    userIDs=[]
    sIDs=[]
    for _, u in Users.iterrows():
        if d in u['Domains']:
            # print(d)
            # print(userIDs)
            userIDs.append(u['UserID'])
    
    for IDs in userIDs:
        sIDs.extend(Users.loc[IDs,'Services'])
    schedule=pd.DataFrame(columns=['ServiceID','sCores','sBandwidth'])
    i=0
    for s in sIDs:
        schedule.loc[i,'sCores']=Services.loc[s,'sCores']
        schedule.loc[i,'sBandwidth']=Services.loc[s,'sBandwidth']
        schedule.loc[i,'ServiceID']=s
        i=i+1
    # print(schedule.to_string())
    # schedule=schedule.groupby('ServiceID', as_index=False).agg({'sCores': 'sum', 'sBandwidth': 'first'})
    # print(schedule)
    if opt[0]=='worstfit' and opt[1]=='MaxMax':
        nodes,_=WorstFitMaxMax(schedule)
        # print(nodes)
    # elif opt[0]=='firstfit' and opt[1]=='maxmax':
    #     nodes,_=FirstFitMaxMax(schedule)
    #     print(nodes)
    elif opt[0]=='bestfit' and opt[1]=='MaxMax':
        nodes,_=BestFitMaxMax(schedule)
        # print(nodes)
    elif opt[0]=='worstfit' and opt[1]=='MinMin':
        nodes,_=WorstFitMinMin(schedule)
        # print(nodes)
    elif opt[0]=='bestfit' and opt[1]=='MinMin':
        nodes,_=BestFitMinMin(schedule)
        # print(nodes)
    else:
        print('invalid heuristic')
        print(opt)
        # default on worstfit maxmax
        # nodes,_=WorstFitMaxMax(schedule)
    return nodes


def WorstFitMaxMax(df: pd.DataFrame, max_bandwidth_per_core=MAX_BANDWIDTH_PER_CORE, max_cores_per_node=NUM_CORES_PER_NODE):

    df = df.sort_values(by='sBandwidth', ascending=False).reset_index(drop=True)
    total_cores = math.ceil((df['sCores'] * df['sBandwidth']).sum() / max_bandwidth_per_core)
    initial_nodes = math.ceil(total_cores / max_cores_per_node)
    
    while True:
        # Initialize nodes and cores
        nodes = pd.DataFrame(columns=['cores','totalBandwidth'])
        for n in range(initial_nodes):
            nodes.loc[n,'cores']=[max_bandwidth_per_core] * max_cores_per_node
            nodes.loc[n,'totalBandwidth']=max_bandwidth_per_core*max_cores_per_node

        success = True  # Flag to indicate if allocation succeeds

        for _, row in df.iterrows():
            server_count = row['sCores']
            bandwidth = row['sBandwidth']

            nodes=nodes.sort_values(by='totalBandwidth',ascending=False).reset_index(drop=True)
            
            for i,node in nodes.iterrows():
                allocated = False
                allocated_cores = 0
                nodeCores=node['cores']
                available_cores = [core_index for core_index in range(len(nodeCores)) if nodeCores[core_index] >= bandwidth]
                if len(available_cores) < server_count:
                    continue
                # print("node",node)
                available_cores=sorted(available_cores, key=lambda x: nodeCores[x])
                # print(available_cores)
                for core_index in available_cores:
                    if nodeCores[core_index] >= bandwidth:
                        nodeCores[core_index] -= bandwidth
                        allocated_cores += 1
                    if allocated_cores == server_count:
                        allocated = True
                        nodes.loc[i,'cores']=nodeCores
                        nodes.loc[i,'totalBandwidth']=sum(nodeCores)
                        break
                
                if not allocated:
                    success = False
                    break
            
            if not success:
                break

        if success:
            # print('number of nodes nodes:', initial_nodes)
            return len(nodes), nodes
        else:
            initial_nodes += 1


def BestFitMaxMax(df: pd.DataFrame, max_bandwidth_per_core=MAX_BANDWIDTH_PER_CORE, max_cores_per_node=NUM_CORES_PER_NODE):
    df = df.sort_values(by='sBandwidth', ascending=False).reset_index(drop=True)
    total_cores = math.ceil((df['sCores'] * df['sBandwidth']).sum() / max_bandwidth_per_core)
    initial_nodes = math.ceil(total_cores / max_cores_per_node)
    
    while True:
        # Initialize nodes and cores
        nodes = pd.DataFrame(columns=['cores','totalBandwidth'])
        for n in range(initial_nodes):
            nodes.loc[n,'cores']=[max_bandwidth_per_core] * max_cores_per_node
            nodes.loc[n,'totalBandwidth']=max_bandwidth_per_core*max_cores_per_node

        success = True  # Flag to indicate if allocation succeeds

        for _, row in df.iterrows():
            server_count = row['sCores']
            bandwidth = row['sBandwidth']

            nodes=nodes.sort_values(by='totalBandwidth',ascending=False).reset_index(drop=True)
            
            for i,node in nodes.iterrows():
                allocated = False
                allocated_cores = 0
                nodeCores=node['cores']
                available_cores = [core_index for core_index in range(len(nodeCores)) if nodeCores[core_index] >= bandwidth]
                if len(available_cores) < server_count:
                    continue
                # print("node",node)
                available_cores=sorted(available_cores, key=lambda x: nodeCores[x], reverse=True) 
                # print(available_cores)
                for core_index in available_cores:
                    if nodeCores[core_index] >= bandwidth:
                        nodeCores[core_index] -= bandwidth
                        allocated_cores += 1
                    if allocated_cores == server_count:
                        allocated = True
                        nodes.loc[i,'cores']=nodeCores
                        nodes.loc[i,'totalBandwidth']=sum(nodeCores)
                        break
                
                if not allocated:
                    success = False
                    break
            
            if not success:
                break

        if success:
            # print('number of nodes nodes:', initial_nodes)
            return len(nodes), nodes
        else:
            initial_nodes += 1
    return

def WorstFitMinMin(df: pd.DataFrame, max_bandwidth_per_core=MAX_BANDWIDTH_PER_CORE, max_cores_per_node=NUM_CORES_PER_NODE):

    df = df.sort_values(by='sBandwidth', ascending=False).reset_index(drop=True)
    total_cores = math.ceil((df['sCores'] * df['sBandwidth']).sum() / max_bandwidth_per_core)
    initial_nodes = math.ceil(total_cores / max_cores_per_node)
    
    while True:
        nodes = pd.DataFrame(columns=['cores','totalBandwidth'])
        for n in range(initial_nodes):
            nodes.loc[n,'cores']=[max_bandwidth_per_core] * max_cores_per_node
            nodes.loc[n,'totalBandwidth']=max_bandwidth_per_core*max_cores_per_node

        success = True  # Flag to indicate if allocation succeeds

        for _, row in df.iterrows():
            server_count = row['sCores']
            bandwidth = row['sBandwidth']
            # print('service',row)
            nodes=nodes.sort_values(by='totalBandwidth').reset_index(drop=True)
            for i,node in nodes.iterrows():
                allocated = False
                allocated_cores = 0
                # print("node",node)
                nodeCores=node['cores']
                available_cores = [core_index for core_index in range(len(nodeCores)) if nodeCores[core_index] >= bandwidth]
                if len(available_cores) < server_count:
                    continue
                available_cores=sorted(available_cores, key=lambda x: nodeCores[x])
                # print(available_cores)
                for core_index in available_cores:
                    if nodeCores[core_index] >= bandwidth:
                        nodeCores[core_index] -= bandwidth
                        allocated_cores += 1
                    if allocated_cores == server_count:
                        allocated = True
                        nodes.loc[i,'cores']=nodeCores
                        nodes.loc[i,'totalBandwidth']=sum(nodeCores)
                        break
                
                if not allocated:
                    success = False
                    break
            
            if not success:
                break

        if success:
            # print('number of nodes nodes:', initial_nodes)
            return len(nodes), nodes
        else:
            initial_nodes += 1


def BestFitMinMin(df: pd.DataFrame, max_bandwidth_per_core=MAX_BANDWIDTH_PER_CORE, max_cores_per_node=NUM_CORES_PER_NODE):
    df = df.sort_values(by='sBandwidth', ascending=False).reset_index(drop=True)
    total_cores = math.ceil((df['sCores'] * df['sBandwidth']).sum() / max_bandwidth_per_core)
    initial_nodes = math.ceil(total_cores / max_cores_per_node)
    
    while True:
        # Initialize nodes and cores
        nodes = pd.DataFrame(columns=['cores','totalBandwidth'])
        for n in range(initial_nodes):
            nodes.loc[n,'cores']=[max_bandwidth_per_core] * max_cores_per_node
            nodes.loc[n,'totalBandwidth']=max_bandwidth_per_core*max_cores_per_node
        success = True  # Flag to indicate if allocation succeeds

        for _, row in df.iterrows():
            server_count = row['sCores']
            bandwidth = row['sBandwidth']

            nodes=nodes.sort_values(by='totalBandwidth').reset_index(drop=True)
            
            for i,node in nodes.iterrows():
                allocated = False
                allocated_cores = 0
                nodeCores=node['cores']
                available_cores = [core_index for core_index in range(len(nodeCores)) if nodeCores[core_index] >= bandwidth]
                if len(available_cores) < server_count:
                    continue
                # print("node",node)
                available_cores=sorted(available_cores, key=lambda x: nodeCores[x], reverse=True) 
                # print(available_cores)
                for core_index in available_cores:
                    if nodeCores[core_index] >= bandwidth:
                        nodeCores[core_index] -= bandwidth
                        allocated_cores += 1
                    if allocated_cores == server_count:
                        allocated = True
                        nodes.loc[i,'cores']=nodeCores
                        nodes.loc[i,'totalBandwidth']=sum(nodeCores)
                        break
                
                if not allocated:
                    success = False
                    break
            
            if not success:
                break

        if success:
            # print('number of nodes nodes:', initial_nodes)
            return len(nodes), nodes
        else:
            initial_nodes += 1
    return

def FirstFitMaxMax():
    return
# based on the simultaneous users in the system we can estimate the required resources
# sort the users based on their up time from longer to shorter and get the intersection of the period of each user with the next user. if intersection is empty, the procedure is stopped. The resources of the users are summed up. this gives the number of cores per domain.

# the number of cores can be devided into separate nodes where each node has 3 to 16 cores.
            
def domainNodes(opt):
    print(opt)
    for d in DOMAIN_IDS: 
        nodes=computeNodeCores(d,opt)
        print(nodes)
        nodeNames=[]
        for i in range(nodes):
            nodeNames.append(f'domain{d}_worker{i}')
        domainID=f'domain{d}'
        df=pd.DataFrame(columns=['NodeName', 'NumCores', 'PartitioningHeuristic','NodeSelectionHeuristic'])
        df['NodeName']= nodeNames
        df['PartitioningHeuristic']=[opt[0]]*nodes
        df['NodeSelectionHeuristic']=[opt[1]]*nodes
        df['NumCores']=[NUM_CORES_PER_NODE]*nodes
        if not os.path.exists(main_dir+'domainNodes/'):
            os.mkdir(main_dir+'domainNodes/')
        df.to_csv(main_dir+f'domainNodes/{opt[0]}/{opt[1]}/domainNodes{domainID}.csv', index=False)

if __name__ == '__main__':
    Heuristics(REALLOCATION_H,NODE_SELECTION_H)
    ServiceGenerator(TOTAL_SERVICES,importanceRange,sBandwidthRange,sCoresRange,0)
    UserTiming()
    print("users are generated")
    EventGenerator()
    print("events are generated")
    for opt0 in PARTITIONING_H:
        if not os.path.exists(main_dir+f'domainNodes/{opt0}'):
            os.mkdir(main_dir+f'domainNodes/{opt0}')
        for opt1 in NODE_SELECTION_H:
            if not os.path.exists(main_dir+f'domainNodes/{opt0}/{opt1}'):
                os.mkdir(main_dir+f'domainNodes/{opt0}/{opt1}')
            domainNodes([opt0,opt1])
    print('done')


