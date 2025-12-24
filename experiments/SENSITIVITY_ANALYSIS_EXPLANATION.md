# Sensitivity Analysis Experiments

## 1. Cloud Cost Sensitivity (c1, c2, c3, c4, c5)

### What it tests:
- How the system behaves when **cloud resources cost 1×, 2×, 3×, 4×, 5× more** than edge resources
- Default assumption: Cloud cost = 1× Edge cost
- Sensitivity: What happens when cloud is 2×, 3×, 4×, 5× more expensive?

### Expected Trends:

#### **Cost Metric:**
- As cloud cost multiplier increases (c1 → c5):
  - **Total cost should INCREASE** (cloud resources become more expensive)
  - **Gap between improved and baseline should WIDEN** (improved algorithm avoids expensive cloud better)

#### **QoS Metric:**
- As cloud cost multiplier increases:
  - **QoS might DECREASE slightly** (system avoids cloud to save cost, might reject more services)
  - **Improved should maintain better QoS** than baseline (better resource management)

#### **QoS per Cost Metric:**
- As cloud cost multiplier increases:
  - **QoS/Cost ratio should DECREASE** (paying more for similar QoS)
  - **Improved should show better QoS/Cost** than baseline (more efficient)

### Interpretation:
- If cloud is very expensive (c4, c5), the algorithm should:
  - Use Edge nodes as much as possible
  - Only use Cloud when absolutely necessary
  - Scale up Edge domains more aggressively
  - Accept some QoS degradation to avoid cloud costs

---

## 2. QoS Ratio Sensitivity (0.2, 0.5, 0.8)

### What it tests:
- **Reduced QoS value** as a fraction of Standard QoS
- When a service runs in "reduced mode", it gets reduced QoS value
- Values: 0.2 = 20% of standard, 0.5 = 50% of standard, 0.8 = 80% of standard

### Expected Trends:

#### **Cost Metric:**
- As QoS ratio increases (0.2 → 0.8):
  - **Cost might DECREASE** (reduced mode is "less bad", more acceptable to use)
  - **More services can run in reduced mode without penalty**

#### **QoS Metric:**
- As QoS ratio increases (0.2 → 0.8):
  - **Total QoS should INCREASE** (reduced mode gives more QoS)
  - **System more willing to use reduced mode** (smaller QoS penalty)

#### **QoS per Cost Metric:**
- As QoS ratio increases:
  - **QoS/Cost ratio should INCREASE** (getting better QoS for same cost)
  - **More flexibility in resource management**

### Interpretation:
- **Low ratio (0.2)**: Reduced mode is very penalized
  - System tries hard to avoid reduced mode
  - More scaling, higher cost
  - Better standard QoS but fewer services

- **High ratio (0.8)**: Reduced mode is acceptable
  - System more willing to use reduced mode
  - Less aggressive scaling
  - Lower cost but more services in reduced mode

---

## Figure Organization:

### For each sensitivity parameter:
- **16, 32, 64, 128** = max scaling threshold values
- **Inside each folder:**
  - `robustness_{cost,qos,qosPerCost}_all_{cloud_cost|QoS}.pdf` = 2D comparison plots
  - `robustness_{cloud_cost|QoS}_{cost,qos,qosPerCost}_sheets.pdf` = 3D surface plots

### How to read 2D plots:
- X-axis: Average Core Demand (workload intensity)
- Y-axis: Cost, QoS, or QoS/Cost
- Different lines: Different cloud costs (c1-c5) or QoS ratios (0.2, 0.5, 0.8)
- Improved vs Baseline for each parameter

### How to read 3D plots:
- X-axis: Average Core Demand (workload)
- Y-axis: Cloud Cost (1-5) or QoS Ratio (0.2-0.8)
- Z-axis: Cost, QoS, or QoS/Cost
- Surfaces show how metric varies with both workload AND sensitivity parameter

---

## Usage in Thesis:

These sensitivity analyses demonstrate:
1. **Robustness** - algorithm performance under different cost assumptions
2. **Flexibility** - how system adapts to different QoS requirements
3. **Practical applicability** - real systems have varying cloud costs and QoS needs

Recommended for **"Additional Experimental Results"** or **"Sensitivity Analysis"** section in thesis.

