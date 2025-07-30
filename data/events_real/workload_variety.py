import pandas as pd
import numpy as np
import os

# === PARAMETERS ===
INPUT_FILE = "data/events_real/final_events_with_services.csv"
OUTPUT_DIR = "data/events_real/sparse_datasets"
SPARSITY_LEVELS = [1.0, 0.75, 0.5, 0.25, 0.1]  # % of users to keep

# === 1. Load data and sort by time ===
df = pd.read_csv(INPUT_FILE)
df = df.sort_values(by='EventTime').reset_index(drop=True)

# === 2. Get unique user IDs ===
all_ids = df['ID'].unique()
print(f"Total users: {len(all_ids)}")

# === 3. Create output directory ===
os.makedirs(OUTPUT_DIR, exist_ok=True)

# === 4. Generate sparser datasets ===
for level in SPARSITY_LEVELS:
    n_users = int(len(all_ids) * level)
    selected_ids = np.random.choice(all_ids, n_users, replace=False)
    sparse_df = df[df['ID'].isin(selected_ids)].copy()
    output_file = os.path.join(OUTPUT_DIR, f"events_{int(level * 100)}pct.csv")
    sparse_df.to_csv(output_file, index=False)
    print(f"✅ Saved {output_file} with {n_users} users and {len(sparse_df)} rows.")
