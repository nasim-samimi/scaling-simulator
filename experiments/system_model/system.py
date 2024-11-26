import os
import pandas as pd

# we shouldn't be limited by the number of nodes in the cloud or the domain
# so we choose high numbers
numCloudNodes=100
numDomain=10
numDomainNodes=20
numCoreRange=range(10,21,5)
print(numCoreRange)

main_dir='data/'
partitioningOpt=['bestfit','firstfit','worstfit']


def cloudNodes(numCloudNodes, numCoreRange, partitioningOpt):
    df=pd.DataFrame(columns=['NodeName', 'NumCores', 'PartitioningHeuristic'])
    cloudNodeNames=[]
    for i in range(numCloudNodes):
        cloudNodeNames.append(f'worker{i}')
    # create a directory to store the cloud node information
    if not os.path.exists(main_dir+'cloudNodes/'):
        os.mkdir(main_dir+'cloudNodes/')
    for opt in partitioningOpt:
        if not os.path.exists(main_dir+'cloudNodes/'+opt):
            os.mkdir(main_dir+'cloudNodes/'+opt)
        # os.mkdir(main_dir+'cloudNodes/'+opt)
        for numCores in numCoreRange:
            # os.mkdir(main_dir+'cloudNodes/'+opt+'/numCores='+str(numCores))
            if not os.path.exists(main_dir+f'cloudNodes/{opt}/numCores={numCores}'):
                os.mkdir(main_dir+f'cloudNodes/{opt}/numCores={numCores}')
    
            df['NodeName']= cloudNodeNames
            df['PartitioningHeuristic']=[opt]*numCloudNodes
            df['NumCores']=[numCores]*numCloudNodes
            print(df)
            df.to_csv(main_dir+f'cloudNodes/{opt}/numCores={numCores}/cloudNodes.csv', index=False)
    return

def domainNodes(numNodes, numCoresRange, partitioningOpt,numDomains):
    df=pd.DataFrame(columns=['NodeName', 'NumCores', 'PartitioningHeuristic'])
    nodeNames=[]
    for i in range(numNodes):
        nodeNames.append(f'worker{i}')
    # create a directory to store the domain node information
    if not os.path.exists(main_dir+'domainNodes/'):
        os.mkdir(main_dir+'domainNodes/')
    for opt in partitioningOpt:
        if not os.path.exists(main_dir+'domainNodes/'+opt):
            os.mkdir(main_dir+'domainNodes/'+opt)
        for numCores in numCoresRange:
            if not os.path.exists(main_dir+'domainNodes/'+opt+'/numCores='+str(numCores)):
                os.mkdir(main_dir+'domainNodes/'+opt+'/numCores='+str(numCores))
            for domain in range(numDomains):
                domainID=f'domain{domain}'
                df['NodeName']= nodeNames
                df['PartitioningHeuristic']=[opt]*numNodes
                df['NumCores']=[numCores]*numNodes
                print(df)
                df.to_csv(main_dir+f'domainNodes/{opt}/numCores={numCores}/domainNodes{domainID}.csv', index=False)
    return

if __name__ == '__main__':
    cloudNodes(numCloudNodes, numCoreRange, partitioningOpt)
    domainNodes(numDomainNodes, numCoreRange, partitioningOpt,numDomain)
    print('done')