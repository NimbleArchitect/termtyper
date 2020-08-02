#include "sendkeys_windows.h"

int Sendkey(const char* letter) {
    //if shift set wscan to MapVirtualKey( VK_LSHIFT, 0 );
    // Create an array of generic keyboard INPUT structures
    int c = 0;
    INPUT ip[4];
    for (int n=0; n<4; ++n) {
        ip[n].type = INPUT_KEYBOARD;
        ip[n].ki.wScan = 0;
        ip[n].ki.wVk = 0;
        ip[n].ki.time = 0;
        ip[n].ki.dwFlags = KEYEVENTF_SCANCODE; // 0 for key press
        ip[n].ki.dwExtraInfo = 0;
    }

    SHORT mapkey = VkKeyScan(*letter);
    int shift = (mapkey & 0x100) > 0;
    
    if (shift == 1) {
        // fprintf(stderr, "** shift\n");
        ip[c].ki.dwFlags = KEYEVENTF_SCANCODE;
        ip[c].ki.wScan = MapVirtualKey( VK_LSHIFT, 0 );
        c = 1;
    }

    // fprintf(stderr, "** %c\n",*letter);

    ip[c].ki.wScan = MapVirtualKey( LOBYTE(mapkey), 0);
    c++;
    //release key
    ip[c].ki.wScan = MapVirtualKey( LOBYTE(mapkey), 0);
    ip[c].ki.dwFlags =  KEYEVENTF_SCANCODE | KEYEVENTF_KEYUP;

    if (shift == 1) {
        // fprintf(stderr, "** none\n");
        ip[3].ki.dwFlags = KEYEVENTF_SCANCODE | KEYEVENTF_KEYUP;
        ip[3].ki.wScan = MapVirtualKey( VK_LSHIFT, 0 );
    }

    if (shift == 1) {
        SendInput(4, ip, sizeof(INPUT));
    } else {
        SendInput(2, ip, sizeof(INPUT));
    }
    
    return 0;
}

int SendAltTab(void) {
    // Create an array of generic keyboard INPUT structures
    INPUT ip[4];
    for (int n=0; n<4; ++n) {
        ip[n].type = INPUT_KEYBOARD;
        ip[n].ki.wScan = 0;
        ip[n].ki.wVk = 0;
        ip[n].ki.time = 0;
        ip[n].ki.dwFlags = KEYEVENTF_SCANCODE; // 0 for key press
        ip[n].ki.dwExtraInfo = 0;
    }

    ip[0].ki.dwFlags = KEYEVENTF_SCANCODE; // 0 for key press
    ip[0].ki.wScan = 0x38; //press alt

    ip[1].ki.dwFlags = 0; // 0 for key press
    ip[1].ki.wVk = VK_TAB;; //press tab

    ip[2].ki.dwFlags = KEYEVENTF_KEYUP;
    ip[2].ki.wVk = VK_TAB; //release tab
    
    ip[3].ki.wScan = 0x38;
    ip[3].ki.dwFlags = KEYEVENTF_SCANCODE | KEYEVENTF_KEYUP;

    SendInput(4, ip, sizeof(INPUT));

    return 0;
}
