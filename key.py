#!/usr/bin/python3

import sys
from pynput.keyboard import Key, Controller

keyboard = Controller()
line = sys.stdin.readline()
if str.strip(line) == "{ENTER}":
    keyboard.type("\n")
else:
    keyboard.type(str.rstrip(line))
