from results_cores import *

if __name__ == '__main__':
    metrics = ['cost', 'qos', 'qosPerCost']
    
    print("Generating Cloud Cost 2D comparison plots...")
    robustness_compare_cloudcost(flags='allOpts')
    
    print("\nGenerating Cloud Cost 3D sheets...")
    for metric in metrics:
        robustness_cloud_cost_3d_sheets(dir1='improved/allOpts', dir2='baseline', metric=metric, flags='allOpts', nodesize=8)
    
    print("\n✅ All cloud cost plots generated!")

