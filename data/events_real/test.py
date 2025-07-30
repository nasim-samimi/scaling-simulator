import pandas as pd

df = pd.read_csv("domain_events_cleaned.csv")
unique_ids = df['ID'].nunique()

print(f"✅ Number of unique IDs: {unique_ids}")