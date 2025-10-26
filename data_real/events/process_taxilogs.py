import pandas as pd
import numpy as np
import os
import glob
from datetime import datetime
from sklearn.cluster import KMeans
import math

# === CONFIG ===
INPUT_FOLDER = "data_real/events/taxi_log/"  # UPDATE this
CLEANED_EVENTS_FILE = "data_real/events/taxi_log_domain_events_cleaned.csv"
NUM_DOMAINS = 10
DISTANCE_THRESHOLD = 100  # meters or arbitrary units

# === 1. Read and Combine All TXT Files ===
all_files = sorted(glob.glob(os.path.join(INPUT_FOLDER, "*.txt")))
# print(all_files)

dfs = []
i=0
for file in all_files:
    df = pd.read_csv(
        file,
        header=None,  # ⬅️ No column names in the file
        names=["ID", "datetime", "xloc_kf", "yloc_kf"]  # ⬅️ Define them explicitly
    )
    # print(df)
    df['time'] = pd.to_datetime(df['datetime'], errors='coerce')
    df = df.dropna(subset=['time'])  # remove any unparseable timestamps
    df['time'] = df['time'].astype(np.int64) // 10**6  # convert to UNIX timestamp
    dfs.append(df[['ID', 'time', 'xloc_kf', 'yloc_kf']])
    i+=1
    if i==5000:
        break

df = pd.concat(dfs).dropna()
print(f"✅ Loaded {len(df)} records from {len(all_files)} files.")

# === 2. Cluster to Define Domain Locations ===
coords = df[['xloc_kf', 'yloc_kf']]
kmeans = KMeans(n_clusters=NUM_DOMAINS, random_state=42).fit(coords)
domain_centers = kmeans.cluster_centers_

# === 3. Assign Domain Based on Proximity ===
def get_domain(x, y):
    min_dist = float('inf')
    domain = None
    for i, (dx, dy) in enumerate(domain_centers):
        dist = math.sqrt((x - dx)**2 + (y - dy)**2)
        if dist < min_dist:
            min_dist = dist
            domain = i
    return domain if min_dist <= DISTANCE_THRESHOLD else np.nan

df['domain'] = df.apply(lambda row: get_domain(row['xloc_kf'], row['yloc_kf']), axis=1)

# === 4. Detect Allocate & Deallocate Events ===
df = df.sort_values(by=['ID', 'time'])
df['prev_domain'] = df.groupby('ID')['domain'].shift(1)

events = []
for i, row in df.iterrows():
    current = row['domain']
    previous = row['prev_domain']
    obj_ID = row['ID']
    timestamp = row['time']

    if pd.isna(previous) and not pd.isna(current):
        events.append({'ID': obj_ID, 'time': timestamp, 'domain': current, 'event_type': 'allocate'})
    elif not pd.isna(previous):
        if pd.isna(current):
            events.append({'ID': obj_ID, 'time': timestamp, 'domain': previous, 'event_type': 'deallocate'})
        elif current != previous:
            events.append({'ID': obj_ID, 'time': timestamp, 'domain': previous, 'event_type': 'deallocate'})
            events.append({'ID': obj_ID, 'time': timestamp, 'domain': current, 'event_type': 'allocate'})

events_df = pd.DataFrame(events)

# === 5. Clean Redundant Alloc/Dealloc ===
cleaned_events = []
connection_state = {}

events_df = events_df.sort_values(["ID", "time"]).reset_index(drop=True)

for idx in range(len(events_df)):
    row = events_df.loc[idx]
    obj_id = row["ID"]
    domain = row["domain"]
    event_type = row["event_type"]
    time = row["time"]

    next_row = events_df.loc[idx + 1] if idx + 1 < len(events_df) else None
    current_domain = connection_state.get(obj_id, None)

    if event_type == "allocate":
        if current_domain != domain:
            cleaned_events.append({
                "ID": obj_id,
                "time": time,
                "domain": domain,
                "event_type": "allocate"
            })
            connection_state[obj_id] = domain

    elif event_type == "deallocate":
        is_redundant = (
            next_row is not None and
            next_row["ID"] == obj_id and
            next_row["event_type"] == "allocate" and
            next_row["domain"] == domain
        )

        if current_domain == domain and not is_redundant:
            cleaned_events.append({
                "ID": obj_id,
                "time": time,
                "domain": domain,
                "event_type": "deallocate"
            })
            connection_state[obj_id] = None

# === 6. Save Output ===
cleaned_df = pd.DataFrame(cleaned_events)
def drop_if_last_is_allocate(group):
    if group.iloc[-1]['event_type'] == 'allocate':
        return group.iloc[:-1]  # drop last row
    else:
        return group  # keep all

# Sort by time and apply the function per user
cleaned_df = (
    cleaned_df.sort_values(['ID', 'time'])
    .groupby('ID', group_keys=False)  # preserve flat structure
    .apply(drop_if_last_is_allocate)
    .reset_index(drop=True)
)
cleaned_df.to_csv(CLEANED_EVENTS_FILE, index=False)
print(f"✅ Cleaned event list saved to {CLEANED_EVENTS_FILE}")
