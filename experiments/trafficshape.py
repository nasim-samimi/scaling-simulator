import pandas as pd
import matplotlib.pyplot as plt

addition = [0, 0.2, 0.4, 0.6, 0.8, 1.0, 1.2, 1.4, 1.6, 1.8, 2.0]

for a in addition:
    # Load the dataset
    df = pd.read_csv(f'data/events/hightraffic/events_{a}.csv')

    # Sort by EventTime to process in order
    df = df.sort_values("EventTime")

    # Get unique DomainIDs
    domain_ids = df["DomainID"].unique()

    # Process each domain separately
    for domain in domain_ids:
        # Filter data for this DomainID
        df_domain = df[df["DomainID"] == domain]

        # Create a time series dictionary for this domain
        utilization_timeline = {}

        # Dictionary to track active events and their utilization within this domain
        active_events = {}

        for _, row in df_domain.iterrows():
            event_time = row["EventTime"]
            event_id = row["EventID"]
            util = row["TotalUtil"]
            event_type = row["EventType"]

            if event_type == "allocate":
                # Add the event utilization to the active events for this domain
                active_events[event_id] = util
            elif event_type == "deallocate" and event_id in active_events:
                # Remove the event when it is deallocated
                del active_events[event_id]

            # Compute total utilization at the current time step
            total_utilization = sum(active_events.values())

            # Store in the timeline
            utilization_timeline[event_time] = total_utilization

        # Convert the timeline dictionary to a DataFrame for plotting
        df_timeline = pd.DataFrame(list(utilization_timeline.items()), columns=["EventTime", "TotalUtil"])

        # Plot for this domain
        plt.figure(figsize=(10, 5))
        plt.plot(df_timeline["EventTime"], df_timeline["TotalUtil"], marker='o', linestyle='-')
        plt.ylabel("Summed TotalUtil (Active)")
        plt.xlabel("Event Time")
        plt.title(f"Total Utilization Over Time (Accumulated) - Domain {domain} (Addition {a})")
        plt.grid(True)
        plt.savefig(f"hist/histogram_accumulated_domain_{domain}_addition_{a}.png")
        plt.close()
