import os
import pandas as pd
import random

#accept a range of values for services
# variables: total number of services, importance, bandwidth, cores
totalServices=500
importanceRange=range(1,totalServices+1)
sBandwidthRange=range(5,81,5)
sCoresRange=range(2,11,2) # max number of cores should be matched with the number of cores in the system


def custom_probability(x):
    return 0.7 if x < 40 else 0.3 


def reduced(bandwidth:list,cores:list) -> tuple :
    # require some strategy to generate the reduced parameters
    # reduced strategy has impact on qos
    rEBandwidth=[x / 2 for x in bandwidth]
    rECores=[x / 2 for x in cores]
    rCBandwidth=[x / 2 for x in bandwidth]
    rCCores=[x / 2 for x in cores]
    return rEBandwidth,rECores,rCBandwidth,rCCores

def service(numServices,importance,sBandwidth,sCores,number): # it returns random values for the services for now
    seed=0
    serviceImportance=random.sample(importance,numServices)
    weights = [1.5 if x < 40 else 0.5 for x in sBandwidth]  # Higher weight for numbers below 40
    weights = [w / sum(weights) for w in weights]
    serviceBandwidth=random.choices(sBandwidth,k=numServices,weights=weights)
    serviceCores=random.choices(sCores,k=numServices)
    rEBandwidth,rECores,rCBandwidth,rCCores=reduced(serviceBandwidth,serviceCores)
    df=pd.DataFrame(columns=['ServiceID', 'Importance', 'sBandwidth', 'sCores', 'rEBandwidth', 'rECores', 'rCBandwidth', 'rCCores'])
    df['ServiceID']=range(numServices)
    df['Importance']=serviceImportance
    df['sBandwidth']=serviceBandwidth
    df['sCores']=serviceCores
    df['rEBandwidth']=rEBandwidth
    df['rECores']=rECores
    df['rCBandwidth']=rCBandwidth
    df['rCCores']=rCCores
    if not os.path.exists('data/services/'):
        os.mkdir('data/services/')
    df.to_csv(f'data/services/services{number}.csv', index=False)
    return

if __name__ == '__main__':
    numServiceSets=10
    for i in range(numServiceSets):
        service(totalServices,importanceRange,sBandwidthRange,sCoresRange,i)
    print('done')

