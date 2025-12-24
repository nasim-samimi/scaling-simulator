# Two Versions of Results Plotting

I've created two versions so you can choose which x-axis to use:

## Option 1: Original "Extra Workload %" (results.py)

**File:** `results.py`  
**Runners:** `runner_robustness_*.py` (your existing runners)  
**Output:** `experiments/plots/`

**X-axis:**
- Label: "extra workload (%)"
- Values: 0, 20, 40, 60, 80, 100, 120, 140, 160, 180, 200

**Use this if:** You want to keep the original format

---

## Option 2: "Average Core Demand" (results_cores.py) ⭐ NEW

**File:** `results_cores.py`  
**Runners:**
- `runner_cores_all_ths.py` - Main robustness plots + 3D sheets
- `runner_cores_compare_nodesize.py` - Node size comparison
- `runner_cores_compare_nodecore.py` - Node core selection comparison
- `runner_cores_all_interval_based_ths.py` - Interval-based with 3D sheets
- `runner_cores_all_interval_based_ths_real.py` - Interval-based with real data

**Output:** `experiments/plots_with_cores/`

**X-axis:**
- Label: "Average Core Demand"  
- Values: 264, 317, 370, 422, 475, 528, 581, 634, 686, 739, 792

**Use this if:** You want to show actual average core demand (clearer for readers)

---

## How to Generate Plots with Core Demand

```bash
cd /Users/nsamimidehkordi/scaling-simulator
source /opt/homebrew/anaconda3/bin/activate base

# Run the specific runners you need:
python experiments/runner_cores_all_ths.py                      # Main plots + 3D
python experiments/runner_cores_compare_nodesize.py              # Node size comparison
python experiments/runner_cores_compare_nodecore.py              # Node core comparison
python experiments/runner_cores_all_interval_based_ths.py        # Interval-based
python experiments/runner_cores_all_interval_based_ths_real.py   # Interval-based (real data)
```

Plots will be saved to: `experiments/plots_with_cores/`

---

## What Changed in results_cores.py

1. **Added workload mapping:**
```python
WORKLOAD_CORES_SYSTEMWIDE = {
    0: 264.08,    # baseline average core demand
    0.2: 316.82,
    ...
    2.0: 791.87   # 3× baseline
}
```

2. **Changed x-axis label:**
- From: `"extra workload (%)"`
- To: `"Average Core Demand"`

3. **Changed x-axis values:**
- From: `avgs['addition'] * 100` (gives 0, 20, 40, ...)
- To: `avgs['addition'].map(WORKLOAD_CORES_SYSTEMWIDE)` (gives 264, 317, 370, ...)

4. **Changed output directory:**
- From: `plots='experiments/plots/'`
- To: `plots='experiments/plots_with_cores/'`

---

## For the Paper

With the new version, you can write:

> "We measure workload as average core demand across the system. The baseline configuration generates an average demand of **264 cores**. To evaluate algorithm robustness, we scale the workload from 264 to 792 cores (3× baseline) by injecting additional service requests..."

This is much clearer than "0% to 200% extra workload"!

---

## Which Version to Use?

**I recommend Option 2 (Average Core Demand)** because:
- ✅ Concrete and understandable
- ✅ No confusing "extra" vs "total"
- ✅ Readers immediately know the actual system load
- ✅ Can compare directly to system capacity

But you have both versions available!

