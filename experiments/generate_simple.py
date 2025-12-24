#!/usr/bin/env python3
"""
Simple test to generate one figure with new labels
"""
import os
import sys

# Add current directory to path
sys.path.insert(0, os.path.dirname(__file__))

# Import everything from results
from results import *

print("="*60)
print("Workload mapping (WORKLOAD_CORES_PER_DOMAIN):")
for k, v in WORKLOAD_CORES_PER_DOMAIN.items():
    print(f"  {k} -> {v:.2f} cores/domain")
print("="*60)

# Try generating one simple figure
print("\nTesting: robustness_max_scaling_size...")
try:
    robustness_max_scaling_size(
        dir1='improved/allOpts',
        dir2='baseline',
        metric='qos',
        flags='all',
        nodesize=8
    )
    print("✓ SUCCESS!")
except Exception as e:
    print(f"✗ Error: {e}")
    import traceback
    traceback.print_exc()

print("\n" + "="*60)
print("Check: experiments/plots/robustness/all/...")

