import os
import pandas as pd
from users import *
from services import *
import math
import random
# we shouldn't be limited by the number of nodes in the cloud or the domain
# so we choose high numbers

NUM_CORES_PER_SCALED_NODE=8
NUM_CORES_PER_INIT_NODE=32

main_dir='data/'
PARTITIONING_H=['bestfit','worstfit']

REALLOCATION_H=["HBI","HCI","HBCI","HBIcC"]
NODE_SELECTION_H=["MinMin","MaxMax"]

MAX_BANDWIDTH_PER_CORE=100
EVENTS_LENGTH=1000


# deriving the lower bound.
def computeNodeCoresLowerbound(d,opt):
    userIDs=[]
    sIDs=[]
    domainUsers=Users[Users['Domains'].apply(lambda x: d in x)] # what does this do?
    # compute the prtion of up time per max arrival time per user
    totalUtil = 0
    for _, user in domainUsers.iterrows():
        totalUtil += user['UpTime'] / user['MaxArrivalTime']
        # print('total util:',totalUtil)
        # print('user up time:',user['UpTime'])
        # print('max arrival time:',user['MaxArrivalTime'])
    OverlappingUsers=math.ceil(totalUtil)
    
    # print('overlapping users:',OverlappingUsers)
    # print('max domain users:',len(domainUsers))
    # print('domain users:',domainUsers)
    selectedUsers = domainUsers.sort_values(by='TotalUtil', ascending=False).head(OverlappingUsers)
    userIDs = selectedUsers['UserID'].tolist()
    
    for IDs in userIDs:
        sIDs.extend(Users.loc[IDs,'Services'])
    schedule=pd.DataFrame(columns=['ServiceID','sCores','sBandwidth'])
    i=0
    for s in sIDs:
        schedule.loc[i,'sCores']=Services.loc[s,'sCores']
        schedule.loc[i,'sBandwidth']=Services.loc[s,'sBandwidth']
        schedule.loc[i,'ServiceID']=s
        i=i+1

    if opt[0]=='worstfit' and opt[1]=='MaxMax':
        nodes,_=WorstFitMaxMax(schedule,max_cores_per_node=NUM_CORES_PER_INIT_NODE)

    elif opt[0]=='bestfit' and opt[1]=='MaxMax':
        nodes,_=BestFitMaxMax(schedule,max_cores_per_node=NUM_CORES_PER_INIT_NODE)
        # print(nodes)
    elif opt[0]=='worstfit' and opt[1]=='MinMin':
        nodes,_=WorstFitMinMin(schedule,max_cores_per_node=NUM_CORES_PER_INIT_NODE)
        # print(nodes)
    elif opt[0]=='bestfit' and opt[1]=='MinMin':
        nodes,_=BestFitMinMin(schedule,max_cores_per_node=NUM_CORES_PER_INIT_NODE)
        # print(nodes)
    else:
        print('invalid heuristic')
        print(opt)
        # default on worstfit maxmax
        # nodes,_=WorstFitMaxMax(schedule)
    return nodes


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

    if opt[0]=='worstfit' and opt[1]=='MaxMax':
        nodes,_=WorstFitMaxMax(schedule,max_cores_per_node=NUM_CORES_PER_SCALED_NODE)

    elif opt[0]=='bestfit' and opt[1]=='MaxMax':
        nodes,_=BestFitMaxMax(schedule,max_cores_per_node=NUM_CORES_PER_SCALED_NODE)
        # print(nodes)
    elif opt[0]=='worstfit' and opt[1]=='MinMin':
        nodes,_=WorstFitMinMin(schedule,max_cores_per_node=NUM_CORES_PER_SCALED_NODE)
        # print(nodes)
    elif opt[0]=='bestfit' and opt[1]=='MinMin':
        nodes,_=BestFitMinMin(schedule,max_cores_per_node=NUM_CORES_PER_SCALED_NODE)
        # print(nodes)
    else:
        print('invalid heuristic')
        print(opt)
        # default on worstfit maxmax
        # nodes,_=WorstFitMaxMax(schedule)
    return nodes


def WorstFitMaxMax(df: pd.DataFrame, max_bandwidth_per_core=MAX_BANDWIDTH_PER_CORE, max_cores_per_node=NUM_CORES_PER_SCALED_NODE):

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


def BestFitMaxMax(df: pd.DataFrame, max_bandwidth_per_core=MAX_BANDWIDTH_PER_CORE, max_cores_per_node=NUM_CORES_PER_SCALED_NODE):
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

def WorstFitMinMin(df: pd.DataFrame, max_bandwidth_per_core=MAX_BANDWIDTH_PER_CORE, max_cores_per_node=NUM_CORES_PER_SCALED_NODE):

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


def BestFitMinMin(df: pd.DataFrame, max_bandwidth_per_core=MAX_BANDWIDTH_PER_CORE, max_cores_per_node=NUM_CORES_PER_SCALED_NODE):
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
        nodeNames=[]
        for i in range(nodes):
            nodeNames.append(f'domain{d}_worker{i}_r')
        domainID=f'domain{d}'
        df=pd.DataFrame(columns=['NodeName', 'NumCores', 'PartitioningHeuristic','NodeSelectionHeuristic'])
        df['NodeName']= nodeNames
        df['PartitioningHeuristic']=[opt[0]]*nodes
        df['NodeSelectionHeuristic']=[opt[1]]*nodes
        df['NumCores']=[NUM_CORES_PER_INIT_NODE]*nodes
        if not os.path.exists(main_dir+'domainNodes/'):
            os.mkdir(main_dir+'domainNodes/')
        df.to_csv(main_dir+f'domainNodes/{opt[0]}/{opt[1]}/domainNodes{domainID}.csv', index=False)

def domainNodesLowerBound(opt):
    print(opt)
    for d in DOMAIN_IDS: 
        nodes=computeNodeCoresLowerbound(d,opt)
        nodeNames=[]
        for i in range(nodes):
            nodeNames.append(f'domain{d}_worker{i}_a')
        domainID=f'domain{d}'
        df=pd.DataFrame(columns=['NodeName', 'NumCores', 'PartitioningHeuristic','NodeSelectionHeuristic'])
        df['NodeName']= nodeNames
        df['PartitioningHeuristic']=[opt[0]]*nodes
        df['NodeSelectionHeuristic']=[opt[1]]*nodes
        df['NumCores']=[NUM_CORES_PER_INIT_NODE]*nodes
        if not os.path.exists(main_dir+'domainNodes/'):
            os.mkdir(main_dir+'domainNodes/')
        df.to_csv(main_dir+f'domainNodes/{opt[0]}/{opt[1]}/domainNodes{domainID}.csv', index=False)

if __name__ == '__main__':
    # Heuristics(REALLOCATION_H,NODE_SELECTION_H)
    Services=ServiceGenerator(NUM_SERVICES,importanceRange,sBandwidthRange,sCoresRange,0)
    UserTiming(Services)
    print("check if the users have a total utilisation")
    print(Users.head())
    print("users are generated")
    EventGenerator(EVENTS_LENGTH)
    print("events are generated")
    for opt0 in PARTITIONING_H:
        if not os.path.exists(main_dir+f'domainNodes/{opt0}'):
            os.mkdir(main_dir+f'domainNodes/{opt0}')
        for opt1 in NODE_SELECTION_H:
            if not os.path.exists(main_dir+f'domainNodes/{opt0}/{opt1}'):
                os.mkdir(main_dir+f'domainNodes/{opt0}/{opt1}')
            domainNodes([opt0,opt1])
    print('done')
    for opt0 in PARTITIONING_H:
        for opt1 in NODE_SELECTION_H:
            if not os.path.exists(main_dir+f'domainNodes/Active/{opt0}/{opt1}'): # or Reserved
                os.mkdir(main_dir+f'domainNodes/Active/{opt0}/{opt1}')
            domainNodesLowerBound([opt0,opt1])
    print('done')



