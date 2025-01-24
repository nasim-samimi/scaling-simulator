import pandas as pd
import os
import random
from services import *

############################################
## Constants
############################################
# TOTAL_DURATION = 10000
NUM_USERS = 100
USER_IDS = range( NUM_USERS )
NUM_SERVICES = 10
SERVICE_IDS = range(NUM_SERVICES )
NUM_DOMAINS = 10
DOMAIN_IDS = range(NUM_DOMAINS )
EVENTS_LENGTH=300

# MAX_ARRIVAL_TIME = TOTAL_DURATION/NUM_USERS
# MIN_ARRIVAL_TIME = MAX_ARRIVAL_TIME/10

ALWAYS_UP_USER = range(1, 5)  # this is number of users that have 100% up time

NUM_STATIONARY_USERS_PER_DOMAIN = range(1, 4)  # stationary to a domain

MAX_UP_TIME=100#int(TOTAL_DURATION/2)
MIN_UP_TIME=NUM_DOMAINS
UP_TIME_RANGE = range(NUM_DOMAINS, MAX_UP_TIME + 1, NUM_DOMAINS)
NUM_SERVICE_PER_USER = range(1, 10)
MOBILITY_RANGE = range(1, NUM_DOMAINS+1 )  # number of domains a user can move to

# INTER_ARRIVAL_TIME_RANGE = range(MIN_ARRIVAL_TIME, MAX_ARRIVAL_TIME+1, 1)  # no exact arrival time information

############################################
## DataFrames
############################################

Users=pd.DataFrame(columns=['UserID', 'UpTime','Services', 'Domains','UpTimePerDomain','Mobility','MinArrivalTime','MaxArrivalTime','TotalUtil']) # total util is the sum of all utilizations of the services times the up time
Users['UserID']=USER_IDS
Users.index = Users['UserID']

Events=pd.DataFrame(columns=['EventTime', 'EventType', 'ServiceID','DomainID'])

############################################
## Functions
############################################

def alwaysUpUsers():
    upUsers = random.choice(ALWAYS_UP_USER)
    alwaysUPUserIDs = random.sample(USER_IDS, upUsers)
    Users.loc[Users['UserID'].isin(alwaysUPUserIDs), 'UpTime'] = 100
    return alwaysUPUserIDs

def UserTiming():
    alwaysUpUserIDs = alwaysUpUsers()
    print(alwaysUpUserIDs)
    otherUserIDs = set(USER_IDS) - set(alwaysUpUserIDs)
    print(otherUserIDs)
    for userID in otherUserIDs:
        numServices = random.choice(NUM_SERVICE_PER_USER)
        sIDs = random.sample(range(NUM_SERVICES ), numServices)

        Users.at[userID, 'Services'] = sIDs
        numDomains = random.choice(MOBILITY_RANGE)
        dID = random.sample(range( NUM_DOMAINS ), numDomains)
        Users.at[userID, 'Domains'] = dID

        upTime=random.choice(UP_TIME_RANGE)
        Users.loc[userID, 'UpTime'] = upTime

        Users.loc[userID, 'UpTimePerDomain'] = upTime / numDomains
        Users.loc[userID, 'Mobility'] = numDomains * 100 / NUM_DOMAINS
        min_arrival_time = random.choice(range(int(upTime*1.25*100),int(upTime*9*100),10))/100
        # print('min arrival time',min_arrival_time)
        Users.loc[userID, 'MinArrivalTime'] = min_arrival_time
        max_arrival_time = random.choice(range(int(min_arrival_time*3*10),int(min_arrival_time*5*10)))/10
        Users.loc[userID, 'MaxArrivalTime'] = max_arrival_time
        Users.loc[userID,'TotalUtil']=Services.loc[Users.loc[userID,'Services'],'sTotalUtil'].sum()*upTime
        # print('max arrival time',max_arrival_time)



    for userID in alwaysUpUserIDs:
        Users.at[userID, 'Services'] = random.sample(range( NUM_SERVICES ), random.choice(NUM_SERVICE_PER_USER))
        numDomains = random.choice(MOBILITY_RANGE)
        Users.at[userID, 'Domains'] = random.sample(range(NUM_DOMAINS ), numDomains)
        Users.loc[userID, 'UpTimePerDomain'] = MIN_UP_TIME
        Users.loc[userID, 'Mobility'] = numDomains * 100 / NUM_DOMAINS
        min_arrival_time=numDomains*MIN_UP_TIME
        max_arrival_time=min_arrival_time
        Users.loc[userID, 'MinArrivalTime'] = min_arrival_time
        Users.loc[userID, 'MaxArrivalTime'] = max_arrival_time
        Users.loc[userID,'TotalUtil']=Services.loc[Users.loc[userID,'Services'],'sTotalUtil'].sum()*MIN_UP_TIME

def EventGenerator(eventsLength=EVENTS_LENGTH):
    Events=pd.DataFrame(columns=['EventTime', 'EventType', 'ServiceID','DomainID','EventID','TotalUtil'])
    EventsAbstract=pd.DataFrame(columns=['EventTime', 'EventType', 'ServiceID','DomainID','EventID','UpTime','TotalUtil'])
    eventType = []
    eventTime = []
    eventDomain=[]
    eventServiceID=[]
    eventID=[]
    arrival=0
    upTime=[]
    util=[]
    i=0
    totalUpTime=0
    for _,u in Users.iterrows():
        # print('user:',u)
        # print(Users.head())
        print('user:',u)
        arrival=random.choice(range(0,int(u['MaxArrivalTime']*10)))/10
        while arrival<eventsLength:
            eT=arrival
            totalUpTime=totalUpTime+u['UpTime']
            for d in u['Domains']:
                eventCount=i
                for s in u['Services']:
                    eventTime.append(round(eT,1))
                    eventType.append('allocate')
                    eventDomain.append(d)
                    eventServiceID.append(s)
                    eventID.append(eventCount)
                    eventCount=eventCount+1
                    upTime.append(u['UpTimePerDomain'])
                    util.append(Services.loc[s,'sTotalUtil']*u['UpTimePerDomain'])
                eT = eT + u['UpTimePerDomain']
                eventCount=i
                for s in u['Services']:
                    eventTime.append(round(eT,1))
                    eventType.append('deallocate')
                    eventDomain.append(d)
                    eventServiceID.append(s)
                    eventID.append(eventCount)
                    eventCount=eventCount+1
                    upTime.append(0)
                    util.append(0)
                i=eventCount
            if u['MinArrivalTime']!=u['MaxArrivalTime']:
                arrival=arrival+random.choice(range(int(u['MinArrivalTime']*10),int(u['MaxArrivalTime']*10)))/10
            else:
                arrival=arrival+u['MinArrivalTime']
            if arrival<eT:
                print('error in arrival time, arrival time is invalid')
                
    Events['EventTime']=eventTime
    # print(eventTime)
    Events['EventType']=eventType
    # print(eventType)
    Events['ServiceID']=eventServiceID
    # print(eventServiceID)
    Events['DomainID']=eventDomain
    Events['EventID']=eventID
    Events['TotalUtil']=util
    EventsAbstract=Events
    EventsAbstract['UpTime']=upTime
    EventsAbstract=EventsAbstract[EventsAbstract['EventType']=='allocate']
    # print(eventDomain)
    Events=Events.sort_values(by=['EventTime'])
    Events.to_csv('data/events_0.csv', index=False)
    Users.to_csv('data/users.csv', index=False)
    totalUtil=sum(util)
    return totalUpTime, EventsAbstract,totalUtil

if __name__ == '__main__':
    ServiceGenerator(TOTAL_SERVICES,importanceRange,sBandwidthRange,sCoresRange,0)
    UserTiming()
    print("users are generated")
    totalUpTime,abstract,util=EventGenerator(eventsLength=EVENTS_LENGTH)
    print("events are generated")
    print('total up time:',totalUpTime)
    print('head of events',Events.head())
    print('tail of events',Events.tail())
    print('abstract events',abstract)
    print('abstract events',abstract.tail())
    print('total util:',util)
    print('done')


    



