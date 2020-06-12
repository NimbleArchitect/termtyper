#!/usr/bin/python3

import sys
from time import sleep
from pynput.keyboard import Key, Controller

keyboard = Controller()

keyboard.press(Key.alt)
keyboard.press(Key.tab)
keyboard.release(Key.tab)
keyboard.release(Key.alt)

sleep(1.2)
line = sys.stdin.readline()
if str.strip(line) == "{ENTER}":
    keyboard.type("\n")
else:
    keyboard.type(str.rstrip(line))
