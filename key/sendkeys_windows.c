// +build windows

#include "sendkeys_windows.h"

#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define WINVER 0x0500
#include <windows.h>


int Sendkey(const char *letter, int shift) {
    // Create an array of generic keyboard INPUT structures
    INPUT ip[10];
    for (n=0 ; n<10 ; ++n)
    {
        ip.type = INPUT_KEYBOARD;
        ip.ki.wScan = 0;
        ip.ki.time = 0;
        ip.ki.dwFlags = 0; // 0 for key press
        ip.ki.dwExtraInfo = 0;
    }
    
    // Now we need to adjust each INPUT structure so that
    // the whole sequence is equivalent to typing "hello".
    // Note that every second structure must be changed
    // from a key press to a key release.
    
    // INPUT structures to press and release the 'H' key
    ip[0].ki.wVk = 'H';
    ip[1].ki.wVk = 'H';
    ip[1].ki.dwFlags = KEYEVENTF_KEYUP;
    // INPUT structures to press and release the 'E' key
    ip[2].ki.wVk = 'E';
    ip[3].ki.wVk = 'E';
    ip[3].ki.dwFlags = KEYEVENTF_KEYUP;
    // INPUT structures to press and release the 'L' key
    ip[4].ki.wVk = 'L';
    ip[5].ki.wVk = 'L';
    ip[5].ki.dwFlags = KEYEVENTF_KEYUP;
    // INPUT structures to press and release the 'L' key
    ip[6].ki.wVk = 'L';
    ip[7].ki.wVk = 'L';
    ip[7].ki.dwFlags = KEYEVENTF_KEYUP;
    // INPUT structures to press and release the 'O' key
    ip[8].ki.wVk = 'O';
    ip[9].ki.wVk = 'O'; 
    ip[9].ki.dwFlags = KEYEVENTF_KEYUP;
    
    // Send all 10 INPUT structures to type the word
    SendInput(10, ip, sizeof(INPUT));
     
    // The End
    return 0;
}