import time
from datetime import datetime

time_format = "%H:%M:%S"
sleep_duration = 1000

print("Current Time =", datetime.now().strftime(time_format))

print(f"sleeping for {sleep_duration}s...")

time.sleep(sleep_duration)

print("Current Time =", datetime.now().strftime(time_format))

print("hello world!")