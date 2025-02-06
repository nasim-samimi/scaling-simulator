import pandas as pd
import os
import matplotlib.pyplot as plt
addition=[0,0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1]
for a in addition:
    df = pd.read_csv(f'data/events/hightraffic/events_{a}.csv')
    df_grouped = df.groupby("EventTime")["TotalUtil"].sum().reset_index()
    plt.figure(figsize=(10, 5))
    plt.plot(df_grouped["EventTime"], df_grouped["TotalUtil"])
    plt.ylabel("Summed TotalUtil")
    plt.xlabel("Event Time")
    plt.title("Histogram of Summed TotalUtil per EventTime")
    plt.grid(True)
    plt.savefig(f"histogram_{a}.png")
