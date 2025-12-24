from results_cores import *

if __name__ == '__main__':
    metrics = ['cost', 'qos', 'qosPerCost']
    
    print("Generating QoS Ratio 2D comparison plots...")
    robustness_compare_qos(flags='allOpts')
    
    print("\nGenerating QoS Ratio 3D sheets...")
    for metric in metrics:
        robustness_qos_3d_sheets(dir1='improved/allOpts', dir2='baseline', metric=metric, flags='allOpts', nodesize=8)
    
    print("\n✅ All QoS ratio plots generated!")

