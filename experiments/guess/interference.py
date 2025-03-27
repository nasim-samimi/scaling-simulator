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
MAX_UP_TIME=10000#int(TOTAL_DURATION/2)
MIN_UP_TIME=8000


def generateRandomUser(Users,Services,addedUtil,addition,events:pd.DataFrame):
    u=0
    Users['TotalUtil']=Users['UpTime']*Users['TotalUtil']
    potentialUsers=Users[Users['TotalUtil']<=addedUtil]
    potentialUsers=potentialUsers[potentialUsers['UpTime']>0]
    potentialUsers=potentialUsers.sort_values(by='TotalUtil',ascending=True)
    extraEvents=pd.DataFrame(columns=events.columns)
    ind=0
    eventCount=len(events)+1
    print("in generate random user")
    while u<addedUtil:
        randUsersID=random.choice(potentialUsers.index)
        randUser=Users.loc[randUsersID]
        numAppearance=random.choice(range(1,5))
        randFirstAppearance=random.choice(range(0,int(events['EventTime'].max()*0.4)))
        arrival=randFirstAppearance
        for n in range(numAppearance):
            eT=arrival
            u=u+randUser['TotalUtil']*randUser['UpTime']
            if u>addedUtil:
                break
            for d in randUser['Domains']:
                i=eventCount
                for s in randUser['Services']:
                    extraEvents.loc[ind,'EventTime']=round(eT,2)
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
                    extraEvents.loc[ind,'EventTime']=round(eT,2)
                    extraEvents.loc[ind,'EventType']='deallocate'
                    extraEvents.loc[ind,'DomainID']=d
                    extraEvents.loc[ind,'ServiceID']=s
                    extraEvents.loc[ind,'EventID']=i
                    extraEvents.loc[ind,'TotalUtil']=Services.loc[s,'sTotalUtil']*randUser['UpTimePerDomain']
                    extraEvents.loc[ind,'UpTime']=0
                    i=i+1
                    ind=ind+1
                    
                eventCount=i
            if randUser['MinArrivalTime']!=randUser['MaxArrivalTime']:
                arrival=arrival+random.choice(range(int(randUser['MinArrivalTime']*10),int(randUser['MaxArrivalTime']*5)))/10
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

def generateRandService(addedUtil,events:pd.DataFrame,Services:pd.DataFrame):
    u=0
    extraEvents=pd.DataFrame(columns=events.columns)
    ind=0
    i=len(events)+1
    max_time=events['EventTime'].max()
    range1 = range(int(max_time * 0.1), int(max_time * 0.2))
    range2 = range(int(max_time * 0.35), int(max_time * 0.55))
    range3=range(int(max_time * 0.7), int(max_time * 0.8))
    # range3=range(int(max_time * 0.9), int(max_time * 1))
    c_range=list(range1)+list(range2)+list(range3)

    
    while u<addedUtil:
        randServiceID=random.choice(SERVICE_IDS)
        randService=Services.loc[randServiceID]
        if randService['sTotalUtil']==0:
            continue

        #     ##############################################
        # valid_events = events[(events['EventTime'].isin(c_range)) & (events['EventType'] == 'deallocate')]

        # if valid_events.empty:
        #     continue  # Skip iteration if no valid times found
        # valid_events = valid_events.sample(frac=1).reset_index(drop=True)

        # # Select a random deallocate event
        # chosen_event = valid_events.iloc[0]  # Pick the first after shuffling

        # # Extract EventTime and DomainID from the selected event
        # randFirstAppearance = chosen_event["EventTime"] + 0.1
        # randDomain = chosen_event["DomainID"]
        # #####################################################
        randFirstAppearance=random.choice(c_range) # the 0.8 is to make sure that the service is not added at the end of the events
        randDomain=random.choice(range(NUM_DOMAINS))
        #######################################################

        randUpTime=random.choice(range(MIN_UP_TIME,math.ceil(MAX_UP_TIME)))
        util=randService['sTotalUtil']*randUpTime
        if util>addedUtil-u:
            # randUpTime=randUpTime-1
            # if randUpTime<=0:
            #     randUpTime=1
            #     util=randService['sTotalUtil']*randUpTime
            #     break
            break
        u=u+util
        arrival=randFirstAppearance
        
        eT=arrival
        
        extraEvents.loc[ind,'EventTime']=round(eT,1)
        extraEvents.loc[ind,'EventType']='allocate'
        extraEvents.loc[ind,'DomainID']=randDomain
        extraEvents.loc[ind,'ServiceID']=randServiceID
        extraEvents.loc[ind,'EventID']=i
        extraEvents.loc[ind,'TotalUtil']=randService['sTotalUtil']
        # print('rand up time:',randUpTime)
        extraEvents.loc[ind,'UpTime']=randUpTime
        ind=ind+1
        eT = eT + randUpTime
    
        extraEvents.loc[ind,'EventTime']=round(eT,1)
        extraEvents.loc[ind,'EventType']='deallocate'
        extraEvents.loc[ind,'DomainID']=randDomain
        extraEvents.loc[ind,'ServiceID']=randServiceID
        extraEvents.loc[ind,'EventID']=i
        extraEvents.loc[ind,'TotalUtil']=randService['sTotalUtil']
        extraEvents.loc[ind,'UpTime']=0
        i=i+1
        ind=ind+1
                
        # randServiceID=random.sample(SERVICE_IDS,4)[0]
    print(extraEvents['TotalUtil'].sum())
    print('added util:',addedUtil)
    return extraEvents.sort_values(by='EventTime')

def interference(additions,totalUtil,events,Services,additionStep,events_dir):
    addedUtil=totalUtil*additionStep
    newEvents=events
    for addition in additions: 
        newUserEvents=generateRandService(addedUtil,newEvents,Services)
        newEvents=pd.concat([newEvents,newUserEvents])
        newEvents.sort_values(by='EventTime',inplace=True)
        event_type_order = {"allocate": 0, "deallocate": 1}

# Sort by EventType first, then by EventTime
        # newEvents.sort_values(
        #     by=['EventType', 'EventTime'], 
        #     ascending=[True, True],  # Ascending for EventType (allocate first), then for EventTime
        #     key=lambda x: x.map(event_type_order) if x.name == 'EventType' else x  # Apply custom sorting for EventType
        # )
        newEvents.to_csv(f'{events_dir}/events_{addition}.csv',index=False)

        print('done')
