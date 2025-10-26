import pandas as pd

df = pd.read_csv("data_real/events/taxi_log_domain_events_cleaned.csv")
unique_ids = df['ID'].nunique()

print(f"✅ Number of unique IDs: {unique_ids}")