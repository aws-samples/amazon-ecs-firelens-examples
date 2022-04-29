import os
import sys
from time import sleep
file1 = open('/houndofbaskerville.txt', 'r')
Lines = file1.readlines()

INFINITE = False
if 'INFINITE' in os.environ:
    INFINITE = True

iterate = True

while iterate:
    iterate = INFINITE
    # print the whole text to stdout
    count = 0
    for line in Lines:
        count += 1
        print(line.rstrip(), end='')
    print("")
    print(count)
    sleep(5)

    # print the whole text to stderr
    count = 0
    for line in Lines:
        count += 1
        print(line.rstrip(), end='', file=sys.stderr)
    print("", file=sys.stderr)
    print(count)
    sleep(5)