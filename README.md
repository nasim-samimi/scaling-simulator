<!-- generate users and services -->
<!-- generate system according to the users -->
<!-- test the users -->
<!-- baseline algorithm -->

                               

for non_interval_based:
python3 runner_none.py                                 
python3 runner_with_flags_ths.py
python3 runner_allopts.py                              
python3 runner_baseline_ths.py

for interval_based:
python3 runner_none.py  
python3 runner_baseline_ths.py
python3 runner_allopts_interval_based_ths.py 