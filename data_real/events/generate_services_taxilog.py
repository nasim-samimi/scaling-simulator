import pandas as pd
import numpy as np
import random

# === PARAMETERS ===
INPUT_FILE = "data_real/events/taxi_log_domain_events_cleaned.csv"
OUTPUT_FILE = "data_real/events/taxi_log_final_events_with_services.csv"
SERVICE_POOL = list(range(100))  # service IDs 0–99
Services=pd.read_csv('data_real/services/services0.csv')
# === 1. Load Cleaned Domain Events ===
df = pd.read_csv(INPUT_FILE)
df = df.sort_values(['ID', 'time']).reset_index(drop=True)

# === 2. Assign Random Services to Each User ===
unique_ids = df['ID'].unique()
user_services = {
    uid: random.sample(SERVICE_POOL, random.randint(10, 20))
    for uid in unique_ids
}

# === 3. Process Allocate/Deallocate Pairs ===
final_rows = []
event_counter = 1
i = 0
while i < len(df) - 1:
    row = df.iloc[i]
    if row['event_type'] == 'allocate':
        # Check if the next row is a matching deallocate
        next_row = df.iloc[i + 1]
        if (
            next_row['event_type'] == 'deallocate' and
            next_row['ID'] == row['ID'] and
            next_row['domain'] == row['domain']
        ):
            uid = row['ID']
            services = user_services[uid]
            for sid in services:
                final_rows.append({
                    'EventTime': float(row['time']),
                    'EventType': 'allocate',
                    'ServiceID': sid,
                    'DomainID': int(row['domain']),
                    'EventID': event_counter,
                    'TotalUtil': Services.loc[sid,'sTotalUtil'],
                    'ID': uid
                })
                final_rows.append({
                    'EventTime': float(next_row['time']),
                    'EventType': 'deallocate',
                    'ServiceID': sid,
                    'DomainID': int(next_row['domain']),
                    'EventID': event_counter,
                    'TotalUtil': 0,
                    'ID': uid
                })
                event_counter += 1
            i += 2  # Skip the deallocate row we just paired
            continue
    i += 1  # Otherwise move to next row

# === 4. Save Output ===
final_df = pd.DataFrame(final_rows)
final_df.to_csv(OUTPUT_FILE, index=False)
print(f"✅ Final service events written to {OUTPUT_FILE}")
