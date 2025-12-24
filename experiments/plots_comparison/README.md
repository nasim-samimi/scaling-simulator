# Figure Comparison: Old vs New Labels

This folder contains a comparison of the figures used in the paper with different x-axis labels.

## Summary of Changes

**OLD X-axis label:** `extra workload (%)`  
**OLD X-axis values:** 0, 20, 40, 60, 80, 100, 120, 140, 160, 180, 200

**NEW X-axis label:** `Workload Multiplier`  
**NEW X-axis values:** 1.0, 1.2, 1.4, 1.6, 1.8, 2.0, 2.2, 2.4, 2.6, 2.8, 3.0

## What This Means

### Baseline Workload
- **264 cores** average system-wide demand
- **26.4 cores/domain** (10 domains)
- **36.8 million core-seconds** total over 38.7-hour simulation

### At Different Multipliers
| Multiplier | Old Label | Avg Cores/Domain | Description |
|------------|-----------|------------------|-------------|
| 1.0× | 0% extra | 26.4 cores | Baseline (normal operation) |
| 1.2× | 20% extra | 31.7 cores | Light overload |
| 1.4× | 40% extra | 37.0 cores | Moderate overload |
| 1.6× | 60% extra | 42.3 cores | |
| 1.8× | 80% extra | 47.5 cores | |
| 2.0× | 100% extra | 52.8 cores | Double baseline |
| 2.2× | 120% extra | 58.1 cores | |
| 2.4× | 140% extra | 63.4 cores | |
| 2.6× | 160% extra | 68.7 cores | |
| 2.8× | 180% extra | 73.9 cores | |
| 3.0× | 200% extra | 79.2 cores | Triple baseline (severe overload) |

## Folders

### old_labels/
Contains the original figures from the paper with:
- X-axis: "extra workload (%)"
- Values: 0, 20, 40, ..., 200

### new_labels/  
Contains the updated figures with:
- X-axis: "Workload Multiplier"
- Values: 1.0, 1.2, 1.4, ..., 3.0

## Figures Included

All 18 figures used in the paper:

### Figure: flags (Reallocation evaluation)
- `robustness_qos_all.pdf`
- `robustness_cost_all.pdf`
- `robustness_qosPerCost_all.pdf`

### Figure: all-nodesizes (Node size comparison)
- `robustness_qos_all_none_nodesizes.pdf`
- `robustness_cost_all_none_nodesizes.pdf`
- `robustness_qosPerCost_all_none_nodesizes.pdf`

### Figure: all (3D scaling size comparison)
- `robustness_qos_3D_all_sheets.pdf`
- `robustness_cost_3D_all_sheets.pdf`
- `robustness_qosPerCost_3D_all_sheets.pdf`

### Figure: all-interval-based (3D interval length comparison)
- `robustness_qos_3D_all_interval_based_sheets.pdf`
- `robustness_cost_3D_all_interval_based_sheets.pdf`
- `robustness_qosPerCost_3D_all_interval_based_sheets.pdf`

### Figure: interval-based (Interval-based method)
- `robustness_qos_interval_based.pdf`
- `robustness_cost_interval_based.pdf`
- `robustness_qosPerCost_interval_based.pdf`

### Figure: interval-based-taxi (T-Drive dataset results)
- `robustness_qos_interval_based_taxi.pdf`
- `robustness_cost_interval_based_taxi.pdf`
- `robustness_qosPerCost_interval_based_taxi.pdf`

## How to Use

1. Compare side-by-side: Open corresponding PDFs from `old_labels/` and `new_labels/`
2. The only difference is the x-axis label and values
3. The data plotted is identical - just presented more clearly

## Recommendation

Use the **new_labels** version in the paper because:
✅ "Workload Multiplier" is self-explanatory  
✅ No confusion about "extra" vs "total"  
✅ Baseline (1×) is clear reference point  
✅ Can explain concrete workload (26.4 cores/domain at 1×)  
✅ Standard presentation in systems research  

## Updated Paper Text

Replace confusing "extra workload" with clear "workload multiplier":

> "We measure workload as average CPU demand over time. The baseline configuration 
> generates an average demand of 26.4 cores per domain. We scale the workload from 
> 1× to 3× baseline by injecting additional service requests, resulting in average 
> demands ranging from 26.4 to 79.2 cores per domain."

