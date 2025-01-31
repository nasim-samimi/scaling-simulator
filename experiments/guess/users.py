import pandas as pd
import os
import random
from services import *

############################################
## Constants
############################################
# TOTAL_DURATION = 10000
NUM_DOMAINS = 10





ALWAYS_UP_USER = range(1, 5)  # this is number of users that have 100% up time

NUM_STATIONARY_USERS_PER_DOMAIN = range(1, 4)  # stationary to a domain

 # number of domains a user can move to

# INTER_ARRIVAL_TIME_RANGE = range(MIN_ARRIVAL_TIME, MAX_ARRIVAL_TIME+1, 1)  # no exact arrival time information

############################################
## DataFrames
############################################

Users=pd.DataFrame(columns=['UserID', 'UpTime','Services', 'Domains','UpTimePerDomain','Mobility','MinArrivalTime','MaxArrivalTime','TotalUtil']) # total util is the sum of all utilizations of the services times the up time

Users.index = Users['UserID']

Events=pd.DataFrame(columns=['EventTime', 'EventType', 'ServiceID','DomainID'])

############################################
## Functions
############################################

def alwaysUpUsers():
    upUsers = random.choice(ALWAYS_UP_USER)
    alwaysUPUserIDs = range(0, upUsers)
    Users.loc[Users['UserID'].isin(alwaysUPUserIDs), 'UpTime'] = 100
    return alwaysUPUserIDs

def UserTiming(Services,num_users,num_domains=NUM_DOMAINS):
    Users['UserID']=range(num_users)
    MAX_UP_TIME=100#int(TOTAL_DURATION/2)
    MIN_UP_TIME=num_domains
    UP_TIME_RANGE = range(num_domains, MAX_UP_TIME + 1, num_domains)
    NUM_SERVICE_PER_USER = range(1, 10)
    MOBILITY_RANGE = range(1, num_domains+1 ) 


    num_services=len(Services)
    alwaysUpUserIDs = alwaysUpUsers()
    print(alwaysUpUserIDs)
    otherUserIDs = range(len(alwaysUpUserIDs), num_users)
    # print('other user ids',otherUserIDs)
    # print('always up user ids',alwaysUpUserIDs)
    print(otherUserIDs)
    for userID in otherUserIDs:
        numServicesPerUser = random.choice(NUM_SERVICE_PER_USER)
        sIDs = random.sample(range(num_services ), numServicesPerUser)
        # print('sids',sIDs)
        # print('userID',userID)
        Users.at[userID, 'Services'] = sIDs
        numDomains = random.choice(MOBILITY_RANGE)
        dID = random.sample(range( num_domains ), numDomains)
        Users.at[userID, 'Domains'] = dID

        upTime=random.choice(UP_TIME_RANGE)
        Users.loc[userID, 'UpTime'] = upTime

        Users.loc[userID, 'UpTimePerDomain'] = upTime / numDomains
        Users.loc[userID, 'Mobility'] = numDomains * 100 / num_domains
        min_arrival_time = random.choice(range(int(upTime*1.25*100),int(upTime*9*100),10))/100
        # print('min arrival time',min_arrival_time)
        Users.loc[userID, 'MinArrivalTime'] = min_arrival_time
        max_arrival_time = random.choice(range(int(min_arrival_time*3*10),int(min_arrival_time*5*10)))/10
        Users.loc[userID, 'MaxArrivalTime'] = max_arrival_time
        Users.loc[userID,'TotalUtil']=Services.loc[Users.loc[userID,'Services'],'sTotalUtil'].sum()
        # print('max arrival time',max_arrival_time)



    for userID in alwaysUpUserIDs:
        Users.at[userID, 'Services'] = random.sample(range( num_services ), random.choice(NUM_SERVICE_PER_USER))
        numDomains = random.choice(MOBILITY_RANGE)
        Users.at[userID, 'Domains'] = random.sample(range(NUM_DOMAINS ), numDomains)
        Users.loc[userID, 'UpTimePerDomain'] = MIN_UP_TIME
        Users.loc[userID, 'Mobility'] = numDomains * 100 / NUM_DOMAINS
        min_arrival_time=numDomains*MIN_UP_TIME
        max_arrival_time=min_arrival_time
        Users.loc[userID, 'MinArrivalTime'] = min_arrival_time
        Users.loc[userID, 'MaxArrivalTime'] = max_arrival_time
        Users.loc[userID,'TotalUtil']=Services.loc[Users.loc[userID,'Services'],'sTotalUtil'].sum()*MIN_UP_TIME

def EventGenerator(eventsLength):
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
    Events.to_csv('data/events/events_0.csv', index=False)
    Users.to_csv('data/users.csv', index=False)
    totalUtil=sum(util)
    return totalUpTime, EventsAbstract,totalUtil


    



