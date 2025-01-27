import pandas as pd
import os
from users import *
import math
import random
import ast

# this script is used to add random behaviour to the generated events.

# add randomness of 50% to the events
# compute while generating the events.
# increase some random user uptime by 10 to 20%
# compute how much has been increased
# add the rest of randomness by random unexpected users.

def unexpectedServices():
    return

# def addRandomness(addition, totalUpTime, events:pd.DataFrame): # addition is in percentage
#     NewUpTime=totalUpTime*(1+addition)
#     addedUpTime=addition*totalUpTime
#     newUsersUpTime=generateRandomUpTime(addedUpTime,events['EventTime'].max())
#     newEventsCount=len(newUsersUpTime)
#     newEventsTime=random.sample(range(0,events['EventTime'].max()),newEventsCount)
#     newEvents=pd.DataFrame(columns=['EventTime', 'EventType', 'ServiceID','DomainID','EventID'])
#     newEventsIDs=range(events['EventID'].max()+1,events['EventID'].max()+newEventsCount+1)
#     newEvents['EventTime']=newEventsTime
#     newEvents['EventType']='allocate'

def generateRandomUser(addedUtil,addition,events:pd.DataFrame):
    u=0
    potentialUsers=Users[Users['TotalUtil']<=addedUtil]
    extraEvents=pd.DataFrame(columns=events.columns)
    ind=0
    eventCount=len(events)+1
    while u<addedUtil:
        randUsersID=random.choice(potentialUsers.index)
        randUser=Users.loc[randUsersID]
        numAppearance=random.choice(range(1,5))
        randFirstAppearance=random.choice(range(0,int(events['EventTime'].max())))
        arrival=randFirstAppearance
        for n in range(numAppearance):
            eT=arrival
            u=u+randUser['TotalUtil']
            if u>addedUtil:
                break
            for d in randUser['Domains']:
                i=eventCount
                for s in randUser['Services']:
                    extraEvents.loc[ind,'EventTime']=round(eT,1)
                    extraEvents.loc[ind,'EventType']='allocate'
                    extraEvents.loc[ind,'DomainID']=d
                    extraEvents.loc[ind,'ServiceID']=s
                    extraEvents.loc[ind,'EventID']=i
                    extraEvents.loc[ind,'UpTime']=randUser['UpTimePerDomain']
                    extraEvents.loc[ind,'TotalUtil']=Services.loc[s,'sTotalUtil']*randUser['UpTimePerDomain']
                    i=i+1
                    ind=ind+1
                eT = eT + randUser['UpTimePerDomain']
                i=eventCount
                
                for s in randUser['Services']:
                    extraEvents.loc[ind,'EventTime']=round(eT,1)
                    extraEvents.loc[ind,'EventType']='deallocate'
                    extraEvents.loc[ind,'DomainID']=d
                    extraEvents.loc[ind,'ServiceID']=s
                    extraEvents.loc[ind,'EventID']=i
                    extraEvents.loc[ind,'TotalUtil']=0
                    extraEvents.loc[ind,'UpTime']=0
                    i=i+1
                    ind=ind+1
                    
                eventCount=i
            if randUser['MinArrivalTime']!=randUser['MaxArrivalTime']:
                arrival=arrival+random.choice(range(int(randUser['MinArrivalTime']*10),int(randUser['MaxArrivalTime']*10)))/10
            else:
                arrival=arrival+randUser['MinArrivalTime']
            if arrival>events['EventTime'].max():
                break
            if arrival<eT:
                print('error in arrival time, arrival time is invalid')
            eventCount=eventCount+1
    return extraEvents.sort_values(by='EventTime')

def generateRandomUpTime(addedUpTime,maxTime):
    t=0
    upTimes=[]
    remainingTime=addedUpTime
    while (t<addedUpTime):
        if addedUpTime-t<2:
            upTimes[-1]=upTimes[-1]+addedUpTime-t
            break
        upTimes.append(random.choice(range(2,min(remainingTime,maxTime))))
        t=t+upTimes[-1]
        remainingTime=remainingTime-upTimes[-1]
        
    return upTimes

def generateRandService(addedUtil,addition,events:pd.DataFrame,Services:pd.DataFrame):
    u=0
    extraEvents=pd.DataFrame(columns=events.columns)
    ind=0
    i=len(events)+1
    while u<addedUtil:
        randServiceID=random.choice(SERVICE_IDS)
        randService=Services.loc[randServiceID]
        randFirstAppearance=random.choice(range(1,int(events['EventTime'].max())))
        randDomain=random.choice(DOMAIN_IDS)
        randUpTime=random.choice(range(math.ceil(MIN_UP_TIME/NUM_DOMAINS),math.ceil(MAX_UP_TIME/NUM_DOMAINS)))
        util=randService['sTotalUtil']*randUpTime
        while util>addedUtil-u:
            randUpTime=randUpTime-1
            if randUpTime<=0:
                randUpTime=1
                util=randService['sTotalUtil']*randUpTime
                break
        u=u+util
        arrival=randFirstAppearance
        
        eT=arrival
        
        extraEvents.loc[ind,'EventTime']=round(eT,1)
        extraEvents.loc[ind,'EventType']='allocate'
        extraEvents.loc[ind,'DomainID']=randDomain
        extraEvents.loc[ind,'ServiceID']=randServiceID
        extraEvents.loc[ind,'EventID']=i
        extraEvents.loc[ind,'TotalUtil']=util
        print('rand up time:',randUpTime)
        extraEvents.loc[ind,'UpTime']=randUpTime
        ind=ind+1
        eT = eT + randUpTime
    
        extraEvents.loc[ind,'EventTime']=round(eT,1)
        extraEvents.loc[ind,'EventType']='deallocate'
        extraEvents.loc[ind,'DomainID']=randDomain
        extraEvents.loc[ind,'ServiceID']=randServiceID
        extraEvents.loc[ind,'EventID']=i
        extraEvents.loc[ind,'TotalUtil']=0
        extraEvents.loc[ind,'UpTime']=0
        i=i+1
        ind=ind+1
                
        # randServiceID=random.sample(SERVICE_IDS,4)[0]
    print(extraEvents['TotalUtil'].sum())
    print('added util:',addedUtil)
    return extraEvents.sort_values(by='EventTime')


if __name__=='__main__':
    Users=pd.read_csv('data/users.csv')
    addition=1
    Users["Services"] = Users["Services"].apply(ast.literal_eval)
    Users["Domains"] = Users["Domains"].apply(ast.literal_eval)
    Services=pd.read_csv('data/services/services0.csv')
    events=pd.read_csv('data/events_0.csv')
    totalUtil=events['TotalUtil'].sum()
    addedUtil=totalUtil*addition

    newUsersEvents=generateRandomUser(addedUtil,addition,events)
    print(newUsersEvents)
    print(events)
    # newEvents=generateRandService(addedUtil,0.5,events,Services)
    # print(newEvents)
    newEvents=pd.concat([events,newUsersEvents])
    newEvents.sort_values(by='EventTime',inplace=True)
    newEvents.to_csv(f'data/events_{addition}.csv',index=False)

    # newEvents.to_csv('data/events/events.csv',index=False)
    print('done')

