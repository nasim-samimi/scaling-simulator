import pandas as pd
import os
import matplotlib.pyplot as plt

df = pd.read_csv('data/events/hightraffic/events_0.5.csv')
print(df)


df_grouped = df.groupby("EventTime")["TotalUtil"].sum().reset_index()

# Step 3: Plot the histogram
plt.figure(figsize=(10, 5))
# plt.hist(df_grouped["TotalUtil"], weights=df_grouped["TotalUtil"], bins=20, edgecolor='black', alpha=0.7)
plt.plot(df_grouped["EventTime"], df_grouped["TotalUtil"])
plt.ylabel("Summed TotalUtil")
plt.xlabel("Event Time")
plt.title("Histogram of Summed TotalUtil per EventTime")
plt.grid(True)
plt.savefig("histogram.png")