import os
import pandas as pd

def events(totalDuration, numEvents, timeScale):
    
    b=0
    e=totalDuration
    objects=[]
    sortedEvents=[]
    for i in range(numEvents):
        objects.append('event'+str(i))
        sortedEvents.append('event'+str(i))

    return