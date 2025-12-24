import pandas as pd
import os
import numpy as np

# Calculate time-weighted average
def time_based_avg(values, times):
    if len(values) == 0 or len(times) == 0:
        return 0
    time_diffs = times['EventTime'].diff().fillna(0)
    weighted_sum = (values * time_diffs).sum()
    total_time = times['EventTime'].iloc[-1] - times['EventTime'].iloc[0]
    return weighted_sum / total_time if total_time > 0 else 0

main_dir = 'experiments/results/'
ADDITIONS = [0, 0.2, 0.4, 0.6, 0.8, 1.0, 1.2, 1.4, 1.6, 1.8, 2.0]
nodesize = 8
max_size = 32

print("="*80)
print("CLOUD COST SENSITIVITY ANALYSIS (max_scaling_threshold=32)")
print("="*80)

cloud_costs = [1, 2, 3, 4, 5]
metrics = ['cost', 'qos', 'qosPerCost']

for metric in metrics:
    print(f"\n{'='*80}")
    print(f"METRIC: {metric.upper()}")
    print(f"{'='*80}")
    
    for c in cloud_costs:
        dirs = main_dir + "cloud_cost/"
        
        # Collect all data points
        improved_values = []
        baseline_values = []
        
        for a in ADDITIONS:
            # Improved
            dir_imp = dirs + f'c{c}/allOpts/max_scaling_threshold={max_size}/'
            fulldir = f'{dir_imp}{metric}/nodesize={nodesize}/addition={a}/mmRB/bestfit/'
            times_addr = f'{dir_imp}eventTime/nodesize={nodesize}/addition={a}/mmRB/bestfit/'
            
            if os.path.exists(fulldir):
                for files in os.listdir(fulldir):
                    if files.endswith('.csv'):
                        values = pd.read_csv(f'{fulldir}{files}', header=None)
                        times = pd.read_csv(f'{times_addr}{files}', header=None)
                        times.columns = ['EventTime']
                        values.columns = [metric]
                        
                        if 'qos' in metric:
                            avg_value = time_based_avg(values[metric]/100, times)
                        else:
                            avg_value = time_based_avg(values[metric], times)
                        improved_values.append(avg_value)
            
            # Baseline
            dir_base = dirs + f'c{c}/baseline/{max_size}/'
            fulldir = f'{dir_base}{metric}/nodesize={nodesize}/addition={a}/None/bestfit/'
            times_addr = f'{dir_base}eventTime/nodesize={nodesize}/addition={a}/None/bestfit/'
            
            if os.path.exists(fulldir):
                for files in os.listdir(fulldir):
                    if files.endswith('.csv'):
                        values = pd.read_csv(f'{fulldir}{files}', header=None)
                        times = pd.read_csv(f'{times_addr}{files}', header=None)
                        times.columns = ['EventTime']
                        values.columns = [metric]
                        
                        if 'qos' in metric:
                            avg_value = time_based_avg(values[metric]/100, times)
                        else:
                            avg_value = time_based_avg(values[metric], times)
                        baseline_values.append(avg_value)
        
        if improved_values and baseline_values:
            imp_avg = np.mean(improved_values)
            base_avg = np.mean(baseline_values)
            
            if metric == 'cost':
                improvement = ((base_avg - imp_avg) / base_avg) * 100
            else:
                improvement = ((imp_avg - base_avg) / base_avg) * 100
            
            print(f"\nCloud Cost Multiplier c{c}:")
            print(f"  Baseline avg: {base_avg:.4f}")
            print(f"  Improved avg: {imp_avg:.4f}")
            print(f"  Improvement: {improvement:.2f}%")

print("\n" + "="*80)
print("QoS RATIO SENSITIVITY ANALYSIS (max_scaling_threshold=32)")
print("="*80)

qos_ratios = [0.2, 0.5, 0.8]

for metric in metrics:
    print(f"\n{'='*80}")
    print(f"METRIC: {metric.upper()}")
    print(f"{'='*80}")
    
    for qr in qos_ratios:
        dirs = main_dir + "QoS/"
        
        # Collect all data points
        improved_values = []
        baseline_values = []
        
        for a in ADDITIONS:
            # Improved
            dir_imp = dirs + f'{qr}/allOpts/max_scaling_threshold={max_size}/'
            fulldir = f'{dir_imp}{metric}/nodesize={nodesize}/addition={a}/mmRB/bestfit/'
            times_addr = f'{dir_imp}eventTime/nodesize={nodesize}/addition={a}/mmRB/bestfit/'
            
            if os.path.exists(fulldir):
                for files in os.listdir(fulldir):
                    if files.endswith('.csv'):
                        values = pd.read_csv(f'{fulldir}{files}', header=None)
                        times = pd.read_csv(f'{times_addr}{files}', header=None)
                        times.columns = ['EventTime']
                        values.columns = [metric]
                        
                        if 'qos' in metric:
                            avg_value = time_based_avg(values[metric]/100, times)
                        else:
                            avg_value = time_based_avg(values[metric], times)
                        improved_values.append(avg_value)
            
            # Baseline
            dir_base = dirs + f'{qr}/baseline/{max_size}/'
            fulldir = f'{dir_base}{metric}/nodesize={nodesize}/addition={a}/None/bestfit/'
            times_addr = f'{dir_base}eventTime/nodesize={nodesize}/addition={a}/None/bestfit/'
            
            if os.path.exists(fulldir):
                for files in os.listdir(fulldir):
                    if files.endswith('.csv'):
                        values = pd.read_csv(f'{fulldir}{files}', header=None)
                        times = pd.read_csv(f'{times_addr}{files}', header=None)
                        times.columns = ['EventTime']
                        values.columns = [metric]
                        
                        if 'qos' in metric:
                            avg_value = time_based_avg(values[metric]/100, times)
                        else:
                            avg_value = time_based_avg(values[metric], times)
                        baseline_values.append(avg_value)
        
        if improved_values and baseline_values:
            imp_avg = np.mean(improved_values)
            base_avg = np.mean(baseline_values)
            
            if metric == 'cost':
                improvement = ((base_avg - imp_avg) / base_avg) * 100
            else:
                improvement = ((imp_avg - base_avg) / base_avg) * 100
            
            print(f"\nQoS Ratio {qr}:")
            print(f"  Baseline avg: {base_avg:.4f}")
            print(f"  Improved avg: {imp_avg:.4f}")
            print(f"  Improvement: {improvement:.2f}%")

print("\n" + "="*80)

