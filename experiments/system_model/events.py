import os
import pandas as pd
import random


totalDuration=100000
timeSteps=1
serviceIDs=range(1,501)
domainIDs=range(1,11)
eventDuration=range(10,51,5)
eventTypes=['allocate']*8
eventTypes.append(['deallocate']*2)

t=0
i=0
event=[]
eventTime=[]
eventType=[]
allocatedServices=[]
df=pd.DataFrame(columns=['EventTime', 'EventType', 'ServiceID','DomainID'])
for t in range(0,totalDuration,timeSteps):
    sID=random.choice(serviceIDs)
    dID=random.choice(domainIDs)
    duration=random.choice(eventDuration)
    df.loc[i]=[t,'allocate',sID,dID]
    df.loc[i+1]=[t+duration,'deallocate',sID,dID]
    i+=2
    df=df.sort_values(by='EventTime')
df.to_csv('data/events.csv', index=False)


    