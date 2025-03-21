import os
import pandas as pd
import random

#accept a range of values for services
# variables: total number of services, importance, bandwidth, cores
NUM_SERVICES = 100
importanceRange=range(1,NUM_SERVICES+1)
sBandwidthRange=range(5,81,1)
sCoresRange=range(1,8)#range(2,11,2) # max number of cores should be matched with the number of cores in the system

SERVICE_IDS=range(NUM_SERVICES)

Services=pd.DataFrame(columns=['ServiceID', 'Importance', 'sBandwidth', 'sCores', 'rEBandwidth', 'rECores', 'rCBandwidth', 'rCCores', 'sTotalUtil'])
Services['ServiceID']=SERVICE_IDS
Services.index = Services['ServiceID']


def reduced(bandwidth:list,cores:list) -> tuple :
    # require some strategy to generate the reduced parameters
    # reduced strategy has impact on qos
    rEBandwidth=[x / 2 for x in bandwidth]
    rECores=[1 for x in cores]
    rCBandwidth=[x/2 for x in bandwidth]
    rCCores=[min(x-1,1) for x in cores]
    return rEBandwidth,rECores,rCBandwidth,rCCores

def ServiceGenerator(numServices,importance,sBandwidth,sCores,number): # it returns random values for the services for now
    seed=0
    serviceImportance=(random.sample(importance,numServices))
    weights = [1.5 if x < 40 else 0.5 for x in sBandwidth]  # Higher weight for numbers below 40
    weights = [w / sum(weights) for w in weights]
    serviceBandwidth=random.choices(sBandwidth,k=numServices,weights=weights)
    serviceCores=random.choices(sCores,k=numServices)
    serviceUtil=pd.Series(serviceBandwidth)*pd.Series(serviceCores)
    rEBandwidth,rECores,rCBandwidth,rCCores=reduced(serviceBandwidth,serviceCores)
    Services['Importance']=serviceImportance
    Services['sBandwidth']=serviceBandwidth
    Services['sCores']=serviceCores
    Services['rEBandwidth']=rEBandwidth
    Services['rECores']=rECores
    Services['rCBandwidth']=rCBandwidth
    Services['rCCores']=rCCores
    print('service util:',serviceUtil)
    Services['sTotalUtil']=serviceUtil.to_list()

    if not os.path.exists('data/services/'):
        os.mkdir('data/services/')
    Services.to_csv(f'data/services/services{number}.csv', index=False)
    return Services

