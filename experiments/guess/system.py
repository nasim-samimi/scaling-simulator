import os
import pandas as pd
from users import *
from services import *
import math
import random
import ast
# we shouldn't be limited by the number of nodes in the cloud or the domain
# so we choose high numbers

NUM_CORES_PER_SCALED_NODE=8
NUM_CORES_PER_INIT_NODE=32

main_dir='data/'
PARTITIONING_H=['bestfit','worstfit']

NODE_SELECTION_H=["MinMin","MaxMax"]

MAX_BANDWIDTH_PER_CORE=100
EVENTS_LENGTH=1000


# deriving the lower bound.
def computeNodeCoresLowerbound(d,opt,num_cores,Users,Services):
    userIDs=[]
    sIDs=[]
    print(Users.head())
    # domainUsers=Users[Users['Domains'].apply(lambda x: d in x)] # what does this do?
    domainUsers = Users[Users['Domains'].apply(lambda x: isinstance(x, (list, str)) and d in x)]

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
    selectedUsers = domainUsers.sort_values(by='TotalUtil', ascending=True).head(OverlappingUsers)
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
        nodes,_=WorstFitMaxMax(schedule,num_cores_per_scaled_node=num_cores,num_init_nodes=0, num_cores_per_init_node=0)

    elif opt[0]=='bestfit' and opt[1]=='MaxMax':
        nodes,_=BestFitMaxMax(schedule,num_cores_per_scaled_node=num_cores,num_init_nodes=0, num_cores_per_init_node=0)
        # print(nodes)
    elif opt[0]=='worstfit' and opt[1]=='MinMin':
        nodes,_=WorstFitMinMin(schedule,num_cores_per_scaled_node=num_cores,num_init_nodes=0, num_cores_per_init_node=0)
        # print(nodes)
    elif opt[0]=='bestfit' and opt[1]=='MinMin':
        nodes,_=BestFitMinMin(schedule,num_cores_per_scaled_node=num_cores,num_init_nodes=0, num_cores_per_init_node=0)
        # print(nodes)
    else:
        print('invalid heuristic')
        print(opt)
        # default on worstfit maxmax
        # nodes,_=WorstFitMaxMax(schedule)
    return nodes


def computeNodeCoresUpperbound(d,opt,num_cores_scaled,Users,Services,num_init_nodes,num_cores_init):
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
        nodes,_=WorstFitMaxMax(schedule,num_cores_per_scaled_node=num_cores_scaled,num_init_nodes=num_init_nodes, num_cores_per_init_node=num_cores_init)

    elif opt[0]=='bestfit' and opt[1]=='MaxMax':
        nodes,_=BestFitMaxMax(schedule,num_cores_per_scaled_node=num_cores_scaled,num_init_nodes=num_init_nodes, num_cores_per_init_node=num_cores_init)
        # print(nodes)
    elif opt[0]=='worstfit' and opt[1]=='MinMin':
        nodes,_=WorstFitMinMin(schedule,num_cores_per_scaled_node=num_cores_scaled,num_init_nodes=num_init_nodes, num_cores_per_init_node=num_cores_init)
        # print(nodes)
    elif opt[0]=='bestfit' and opt[1]=='MinMin':
        nodes,_=BestFitMinMin(schedule,num_cores_per_scaled_node=num_cores_scaled,num_init_nodes=num_init_nodes, num_cores_per_init_node=num_cores_init)
        # print(nodes)
    else:
        print('invalid heuristic')
        print(opt)
        # default on worstfit maxmax
        # nodes,_=WorstFitMaxMax(schedule)
    return nodes

def domainNodesUpperBound(opt,dir,num_init_nodes,Users=Users,Services=Services,num_cores=NUM_CORES_PER_SCALED_NODE,num_domains=NUM_DOMAINS,num_cores_init=NUM_CORES_PER_INIT_NODE):
    print(opt)
    domain_ids=range(num_domains)
    j=0
    for d in domain_ids: 
        nodes=computeNodeCoresUpperbound(d,opt,num_cores,Users=Users,Services=Services,num_init_nodes=num_init_nodes[j],num_cores_init=num_cores_init)
        nodeNames=[]
        for i in range(nodes):
            nodeNames.append(f'domain{d}_worker{i}_r')
        domainID=f'domain{d}'
        df=pd.DataFrame(columns=['NodeName', 'NumCores', 'PartitioningHeuristic','NodeSelectionHeuristic'])
        df['NodeName']= nodeNames
        df['PartitioningHeuristic']=[opt[0]]*nodes
        df['NodeSelectionHeuristic']=[opt[1]]*nodes
        df['NumCores']=[num_cores]*nodes
        if not os.path.exists(dir):
            os.makedirs(dir)
        df.to_csv(f'{dir}domainNodes{domainID}.csv', index=False)
        j=j+1

def domainNodesLowerBound(opt,dir,Users,Services,num_cores=NUM_CORES_PER_INIT_NODE,num_domains=NUM_DOMAINS):
    print(opt)
    domain_ids=range(num_domains)
    # print("num cores:",num_cores)
    nodes_cores=[]
    for d in domain_ids: 
        nodes=computeNodeCoresLowerbound(d,opt,num_cores,Users=Users,Services=Services)
        nodes_cores.append(nodes)
        nodeNames=[]
        for i in range(nodes):
            nodeNames.append(f'domain{d}_worker{i}_a')
        domainID=f'domain{d}'
        df=pd.DataFrame(columns=['NodeName', 'NumCores', 'PartitioningHeuristic','NodeSelectionHeuristic'])
        df['NodeName']= nodeNames
        df['PartitioningHeuristic']=[opt[0]]*nodes
        df['NodeSelectionHeuristic']=[opt[1]]*nodes
        df['NumCores']=[num_cores]*nodes
        if not os.path.exists(dir):
            os.makedirs(dir)
        df.to_csv(f'{dir}domainNodes{domainID}.csv', index=False)
    return nodes_cores

def WorstFitMaxMax(df: pd.DataFrame, max_bandwidth_per_core=MAX_BANDWIDTH_PER_CORE, num_cores_per_scaled_node=NUM_CORES_PER_SCALED_NODE,num_cores_per_init_node=NUM_CORES_PER_INIT_NODE, num_init_nodes=1):

    df = df.sort_values(by='sBandwidth', ascending=False).reset_index(drop=True)
    total_cores = math.ceil((df['sCores'] * df['sBandwidth']).sum() / max_bandwidth_per_core)
    initial_nodes = math.ceil((total_cores-(num_cores_per_init_node*num_init_nodes) )/ num_cores_per_scaled_node)
    if initial_nodes <= 0:
        return 0, pd.DataFrame(columns=['cores','totalBandwidth'])
    
    while True:
        # Initialize nodes and cores
        nodes = pd.DataFrame(columns=['cores','totalBandwidth'])
        for n in range(num_init_nodes):
            nodes.loc[n,'cores']=[max_bandwidth_per_core] * num_cores_per_init_node
            nodes.loc[n,'totalBandwidth']=max_bandwidth_per_core*num_cores_per_init_node

        for n in range(num_init_nodes,num_init_nodes+initial_nodes):
            nodes.loc[n,'cores']=[max_bandwidth_per_core] * num_cores_per_scaled_node
            nodes.loc[n,'totalBandwidth']=max_bandwidth_per_core*num_cores_per_scaled_node

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
            return initial_nodes, nodes
        else:
            initial_nodes += 1


def BestFitMaxMax(df: pd.DataFrame, max_bandwidth_per_core=MAX_BANDWIDTH_PER_CORE, num_cores_per_scaled_node=NUM_CORES_PER_SCALED_NODE,num_cores_per_init_node=NUM_CORES_PER_INIT_NODE, num_init_nodes=1):
    df = df.sort_values(by='sBandwidth', ascending=False).reset_index(drop=True)
    total_cores = math.ceil((df['sCores'] * df['sBandwidth']).sum() / max_bandwidth_per_core)
    initial_nodes = math.ceil((total_cores-(num_cores_per_init_node*num_init_nodes) )/ num_cores_per_scaled_node)
    if initial_nodes <= 0:
        return 0, pd.DataFrame(columns=['cores','totalBandwidth'])
    
    while True:
        # Initialize nodes and cores
        nodes = pd.DataFrame(columns=['cores','totalBandwidth'])
        for n in range(num_init_nodes):
            nodes.loc[n,'cores']=[max_bandwidth_per_core] * num_cores_per_init_node
            nodes.loc[n,'totalBandwidth']=max_bandwidth_per_core*num_cores_per_init_node

        for n in range(num_init_nodes,num_init_nodes+initial_nodes):
            nodes.loc[n,'cores']=[max_bandwidth_per_core] * num_cores_per_scaled_node
            nodes.loc[n,'totalBandwidth']=max_bandwidth_per_core*num_cores_per_scaled_node

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
            return initial_nodes, nodes
        else:
            initial_nodes += 1
    return

def WorstFitMinMin(df: pd.DataFrame, max_bandwidth_per_core=MAX_BANDWIDTH_PER_CORE, num_cores_per_scaled_node=NUM_CORES_PER_SCALED_NODE,num_cores_per_init_node=NUM_CORES_PER_INIT_NODE, num_init_nodes=1):

    df = df.sort_values(by='sBandwidth', ascending=False).reset_index(drop=True)
    total_cores = math.ceil((df['sCores'] * df['sBandwidth']).sum() / max_bandwidth_per_core)
    initial_nodes = math.ceil((total_cores-(num_cores_per_init_node*num_init_nodes)) / num_cores_per_scaled_node)
    if initial_nodes <= 0:
        return 0, pd.DataFrame(columns=['cores','totalBandwidth'])
    
    while True:
        # Initialize nodes and cores
        nodes = pd.DataFrame(columns=['cores','totalBandwidth'])
        for n in range(num_init_nodes):
            nodes.loc[n,'cores']=[max_bandwidth_per_core] * num_cores_per_init_node
            nodes.loc[n,'totalBandwidth']=max_bandwidth_per_core*num_cores_per_init_node

        for n in range(num_init_nodes,num_init_nodes+initial_nodes):
            nodes.loc[n,'cores']=[max_bandwidth_per_core] * num_cores_per_scaled_node
            nodes.loc[n,'totalBandwidth']=max_bandwidth_per_core*num_cores_per_scaled_node

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
            return initial_nodes, nodes
        else:
            initial_nodes += 1


def BestFitMinMin(df: pd.DataFrame, num_cores_per_scaled_node,num_cores_per_init_node, num_init_nodes, max_bandwidth_per_core=MAX_BANDWIDTH_PER_CORE):
    df = df.sort_values(by='sBandwidth', ascending=False).reset_index(drop=True)
    total_cores = math.ceil((df['sCores'] * df['sBandwidth']).sum() / max_bandwidth_per_core)
    # print("total cores:",total_cores)
    # print("num cores per scaled node:",num_cores_per_scaled_node)
    initial_nodes = math.ceil((total_cores-(num_cores_per_init_node*num_init_nodes)) / num_cores_per_scaled_node)
    # print("initial nodes:",initial_nodes)
    if initial_nodes <= 0:
        return 0, pd.DataFrame(columns=['cores','totalBandwidth'])
    
    while True:
        # Initialize nodes and cores
        nodes = pd.DataFrame(columns=['cores','totalBandwidth'])
        for n in range(num_init_nodes):
            nodes.loc[n,'cores']=[max_bandwidth_per_core] * num_cores_per_init_node
            nodes.loc[n,'totalBandwidth']=max_bandwidth_per_core*num_cores_per_init_node

        for n in range(num_init_nodes,num_init_nodes+initial_nodes):
            nodes.loc[n,'cores']=[max_bandwidth_per_core] * num_cores_per_scaled_node
            nodes.loc[n,'totalBandwidth']=max_bandwidth_per_core*num_cores_per_scaled_node
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
            return initial_nodes, nodes
        else:
            initial_nodes += 1
    return

def FirstFitMaxMax():
    return
# based on the simultaneous users in the system we can estimate the required resources
# sort the users based on their up time from longer to shorter and get the intersection of the period of each user with the next user. if intersection is empty, the procedure is stopped. The resources of the users are summed up. this gives the number of cores per domain.

# the number of cores can be devided into separate nodes where each node has 3 to 16 cores.
            


if __name__ == '__main__':
    node_sizes=[8,12,16,20,24,28,32]
    Services=pd.read_csv('data/services/services0.csv')
    Users=pd.read_csv('data/users.csv')
    Users["Services"] = Users["Services"].apply(ast.literal_eval)
    Users["Domains"] = Users["Domains"].apply(ast.literal_eval)
    # for opt0 in PARTITIONING_H:
    #     for opt1 in NODE_SELECTION_H:
    #         num_init_nodes=domainNodesLowerBound([opt0,opt1],main_dir+f'domainNodestest/{opt0}/{opt1}/Active/')
    #         print(num_init_nodes)
    #         domainNodesUpperBound([opt0,opt1],main_dir+f'domainNodestest/{opt0}/{opt1}/Reserved/',num_init_nodes=num_init_nodes,num_cores_init=NUM_CORES_PER_INIT_NODE)
    for opt0 in PARTITIONING_H:
        for opt1 in NODE_SELECTION_H:
            for s in node_sizes:
                num_init_nodes=domainNodesLowerBound([opt0,opt1],main_dir+f'domainNodes{s}/{opt0}/{opt1}/Active/',num_cores=s,num_domains=NUM_DOMAINS)
                domainNodesUpperBound([opt0,opt1],main_dir+f'domainNodes{s}/{opt0}/{opt1}/Reserved/',num_cores=s,num_domains=NUM_DOMAINS,num_init_nodes=num_init_nodes,num_cores_init=s)
            num_init_nodes=domainNodesLowerBound([opt0,opt1],main_dir+f'domainNodesVariable/{opt0}/{opt1}/Active/',num_cores=NUM_CORES_PER_INIT_NODE,num_domains=NUM_DOMAINS)
            domainNodesUpperBound([opt0,opt1],main_dir+f'domainNodesVariable/{opt0}/{opt1}/Reserved/',num_cores=NUM_CORES_PER_SCALED_NODE,num_domains=NUM_DOMAINS,num_init_nodes=num_init_nodes,num_cores_init=NUM_CORES_PER_INIT_NODE)
    print('done with generating domain nodes')
    print('done')



