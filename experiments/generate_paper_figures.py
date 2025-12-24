#!/usr/bin/env python3
"""
Generate only the figures used in the paper with new workload labels.
Output will go to plots_comparison/actual_workload/
"""

import sys
import os

# Temporarily change output directory
import results
original_plots = results.plots
results.plots = 'experiments/plots_comparison/actual_workload/'

# Import all the plotting functions
from results import *

# Create output directory
os.makedirs('experiments/plots_comparison/actual_workload/robustness/all/32/nodesize=8/mmRB/bestfit/', exist_ok=True)
os.makedirs('experiments/plots_comparison/actual_workload/robustness/all/128/', exist_ok=True)
os.makedirs('experiments/plots_comparison/actual_workload/robustness/all/nodesize=8/MMRB/bestfit/3D/', exist_ok=True)
os.makedirs('experiments/plots_comparison/actual_workload/robustness/all_interval_based/nodesize=8/mmRB/bestfit/3D/', exist_ok=True)

print("="*60)
print("Generating paper figures with actual workload values...")
print("="*60)
print(f"\nWorkload values on x-axis:")
for k, v in WORKLOAD_CORES_PER_DOMAIN.items():
    print(f"  addition={k} -> {v:.2f} cores/domain")
print("\n" + "="*60)

# Generate the figures used in the paper
nodesize = 8

print("\n1. Generating robustness_*_all.pdf (reallocation evaluation)...")
try:
    dir1 = 'improved/allOpts'
    robustness_max_scaling_size(dir1=dir1, metric='qos', flags='all', nodesize=nodesize)
    robustness_max_scaling_size(dir1=dir1, metric='cost', flags='all', nodesize=nodesize)
    robustness_max_scaling_size(dir1=dir1, metric='qosPerCost', flags='all', nodesize=nodesize)
    print("   ✓ Generated robustness_*_all.pdf")
except Exception as e:
    print(f"   ✗ Error: {e}")

print("\n2. Generating robustness_*_all_none_nodesizes.pdf...")
try:
    robustness_none_nodesizes(metric='qos')
    robustness_none_nodesizes(metric='cost')
    robustness_none_nodesizes(metric='qosPerCost')
    print("   ✓ Generated robustness_*_all_none_nodesizes.pdf")
except Exception as e:
    print(f"   ✗ Error: {e}")

print("\n3. Generating robustness_*_3D_all_sheets.pdf...")
try:
    robustness_compare_flags_nodesizes_3d_sheets(metric='qos', nodesize=nodesize)
    robustness_compare_flags_nodesizes_3d_sheets(metric='cost', nodesize=nodesize)
    robustness_compare_flags_nodesizes_3d_sheets(metric='qosPerCost', nodesize=nodesize)
    print("   ✓ Generated robustness_*_3D_all_sheets.pdf")
except Exception as e:
    print(f"   ✗ Error: {e}")

print("\n4. Generating robustness_*_3D_all_interval_based_sheets.pdf...")
try:
    robustness_interval_based_ths_3d_sheets(metric='qos', nodesize=nodesize)
    robustness_interval_based_ths_3d_sheets(metric='cost', nodesize=nodesize)
    robustness_interval_based_ths_3d_sheets(metric='qosPerCost', nodesize=nodesize)
    print("   ✓ Generated robustness_*_3D_all_interval_based_sheets.pdf")
except Exception as e:
    print(f"   ✗ Error: {e}")

print("\n5. Generating robustness_*_interval_based.pdf...")
try:
    robustness_interval_based(metric='qos', m=16, nodesize=nodesize)
    robustness_interval_based(metric='cost', m=16, nodesize=nodesize)
    robustness_interval_based(metric='qosPerCost', m=16, nodesize=nodesize)
    print("   ✓ Generated robustness_*_interval_based.pdf")
except Exception as e:
    print(f"   ✗ Error: {e}")

print("\n6. Generating robustness_*_interval_based_taxi.pdf...")
try:
    robustness_interval_based_real(metric='qos', m=16, nodesize=nodesize)
    robustness_interval_based_real(metric='cost', m=16, nodesize=nodesize)
    robustness_interval_based_real(metric='qosPerCost', m=16, nodesize=nodesize)
    print("   ✓ Generated robustness_*_interval_based_taxi.pdf")
except Exception as e:
    print(f"   ✗ Error: {e}")

print("\n" + "="*60)
print("DONE! New figures saved to:")
print("  experiments/plots_comparison/actual_workload/")
print("="*60)
print("\nX-axis now shows: 'Average Workload (cores/domain)'")
print("Values: 26.4, 31.7, 37.0, 42.2, 47.5, 52.8, 58.1, 63.4, 68.6, 73.9, 79.2")

