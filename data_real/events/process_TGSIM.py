import pandas as pd
import numpy as np
from sklearn.cluster import KMeans
import math

# === PARAMETERS ===
INPUT_FILE = "data/events_real/TGSIM.csv"  # <-- replace with your actual file name
OUTPUT_FILE = "data/events_real/events_real.csv"
CLEANED_EVENTS_FILE = "data/events_real/domain_events_cleaned.csv"
HYSTERESIS_TIME = 1.0  # seconds — you can adjust this
NUM_DOMAINS = 10
DISTANCE_THRESHOLD = 100  # units same as your x/y (e.g., meters)

# === 1. LOAD DATA ===
df = pd.read_csv(INPUT_FILE)
df = df[['ID', 'time', 'xloc_kf', 'yloc_kf']].dropna()
df['time'] = pd.to_numeric(df['time'], errors='coerce')
df = df.dropna(subset=['time'])

# === 2. CLUSTER SPACE TO CREATE DOMAINS ===
coords = df[['xloc_kf', 'yloc_kf']]
kmeans = KMeans(n_clusters=NUM_DOMAINS, random_state=42).fit(coords)
domain_centers = kmeans.cluster_centers_

# === 3. ASSIGN DOMAIN IF WITHIN THRESHOLD ===
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

# === 4. DETECT ALLOCATE & DEALLOCATE EVENTS ===
df = df.sort_values(by=['ID', 'time'])
df['prev_domain'] = df.groupby('ID')['domain'].shift(1)

events = []

for i, row in df.iterrows():
    current = row['domain']
    previous = row['prev_domain']
    obj_ID = row['ID']
    timestamp = row['time']

    if pd.isna(previous) and not pd.isna(current):
        # First time connecting
        events.append({'ID': obj_ID, 'time': timestamp, 'domain': current, 'event_type': 'allocate'})
    elif not pd.isna(previous):
        if pd.isna(current):
            # Disconnected from a domain
            events.append({'ID': obj_ID, 'time': timestamp, 'domain': previous, 'event_type': 'deallocate'})
        elif current != previous:
            # Switched domain: deallocate old, allocate new
            events.append({'ID': obj_ID, 'time': timestamp, 'domain': previous, 'event_type': 'deallocate'})
            events.append({'ID': obj_ID, 'time': timestamp, 'domain': current, 'event_type': 'allocate'})


# === 5. SAVE RESULT ===
events_df = pd.DataFrame(events)

cleaned_events = []
connection_state = {}

# Sort first
events_df = events_df.sort_values(["ID", "time"]).reset_index(drop=True)

for idx in range(len(events_df)):
    row = events_df.loc[idx]
    obj_id = row["ID"]
    domain = row["domain"]
    event_type = row["event_type"]
    time = row["time"]

    # Get next row if exists
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
        # Lookahead: is next event an immediate reallocate to the same domain?
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
# Convert and save
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
print(f"Cleaned event list saved to {CLEANED_EVENTS_FILE}")

