import pandas as pd
import os
import random
from services import *
from users import *
from system import *
from interference import *

############################################
## Constants for users
############################################
# TOTAL_DURATION = 10000
NUM_USERS = 100
NUM_DOMAINS = 10
ALWAYS_UP_USER = range(1, 5)  # this is number of users that have 100% up time
NUM_STATIONARY_USERS_PER_DOMAIN = range(1, 4)  # stationary to a domain
MAX_UP_TIME=2000

MIN_UP_TIME=500
UP_TIME_RANGE = range(NUM_DOMAINS, MAX_UP_TIME + 1, NUM_DOMAINS)
NUM_SERVICE_PER_USER = range(1, 10)
MOBILITY_RANGE = range(1, NUM_DOMAINS+1 )  # number of domains a user can move to
DOMAIN_IDS = range(NUM_DOMAINS )


############################################
## Constants for services
############################################
NUM_SERVICES = 100
SERVICE_IDS=range(NUM_SERVICES)
importanceRange=range(1,NUM_SERVICES+1)
sBandwidthRange=range(5,81,1)
sCoresRange=range(1,8)#range(2,11,2) # max number of cores should be matched with the number of cores in the system


############################################
## Constants for system
############################################

MAX_BANDWIDTH_PER_CORE=100
NUM_CORES_PER_INIT_NODE=32
NUM_CORES_PER_SCALED_NODE=8

############################################
## Constants for Heuristics
############################################
PARTITIONING_H=['bestfit','worstfit']
NODE_SELECTION_H=["MinMin","MaxMax"]

############################################
## Constants for events
############################################
EVENTS_LENGTH=1000
ADDITION=[0.2,0.4,0.6,0.8,1.0,1.2,1.4,1.6,1.8,2.0]
# ADDITION=[0.8,0.9,1]
node_sizes=[8,16]

main_dir='data/'
events_dir='data/events/hightraffic'
############################################

if __name__=='__main__':
    services=ServiceGenerator(NUM_SERVICES,importanceRange,sBandwidthRange,sCoresRange,0,dir=main_dir)
    users=UserTiming(services,num_users=150,num_domains=NUM_DOMAINS)
    print("users are generated")
    EventGenerator(EVENTS_LENGTH,weight=0.4,Users=users,Services=services,dir=events_dir)
    print("events are generated")
    
    
    
    Users=pd.read_csv('data/users.csv')
    Users["Services"] = Users["Services"].apply(ast.literal_eval)
    Users["Domains"] = Users["Domains"].apply(ast.literal_eval)
    Services=pd.read_csv('data/services/services0.csv')
    events=pd.read_csv(f'{events_dir}/events_0.csv')
    totalUtil=(events['TotalUtil'].mul(events['UpTime'])).sum()

    for opt0 in PARTITIONING_H:
        for opt1 in NODE_SELECTION_H:
            for s in node_sizes:
                num_init_nodes=domainNodesLowerBound([opt0,opt1],main_dir+f'domainNodes{s}/{opt0}/{opt1}/Active/',Users=Users,Services=Services,num_cores=s,num_domains=NUM_DOMAINS)
                domainNodesUpperBound([opt0,opt1],main_dir+f'domainNodes{s}/{opt0}/{opt1}/Reserved/',Users=Users,Services=Services,num_cores=s,num_domains=NUM_DOMAINS,num_init_nodes=num_init_nodes,num_cores_init=s)
            num_init_nodes=domainNodesLowerBound([opt0,opt1],main_dir+f'domainNodesVariable/{opt0}/{opt1}/Active/',Users=Users,Services=Services,num_cores=NUM_CORES_PER_INIT_NODE,num_domains=NUM_DOMAINS)
            domainNodesUpperBound([opt0,opt1],main_dir+f'domainNodesVariable/{opt0}/{opt1}/Reserved/',Users=Users,Services=Services,num_cores=NUM_CORES_PER_SCALED_NODE,num_domains=NUM_DOMAINS,num_init_nodes=num_init_nodes,num_cores_init=NUM_CORES_PER_INIT_NODE)
    print('done with generating domain nodes')


    interference(ADDITION,totalUtil,events,Services,0.2,events_dir)