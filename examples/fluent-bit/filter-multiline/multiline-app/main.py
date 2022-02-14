import os
import time
file1 = open('/test.log', 'r')
Lines = file1.readlines()
 
count = 0

for i in range(10):
    print("app running normally...")
    time.sleep(1)

# Strips the newline character
for line in Lines:
    count += 1
    print(line.rstrip())
print(count)
print("app terminated.")