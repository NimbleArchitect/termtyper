#include "sendkeys_linux.h"

#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>


#  define NeedFunctionPrototypes 1
#  include <X11/Xlib.h>
#  include <X11/keysym.h>
#  include <X11/extensions/XTest.h>
#  if XlibSpecificationRelease != 6
#      error Requires X11R6
#  endif

static char *displayname = NULL;
static int keysym = 0;
static Display *display = NULL;
static Window window = 0;


int OpenDisplay() {
    int focusreturn;
    //fprintf(stderr, "C:OpenDisplay:start\n");
    if (displayname == NULL) {
	    displayname = getenv("DISPLAY");
    }

    if (displayname == NULL) {
	    displayname = ":0.0";
    }

    display = XOpenDisplay(displayname);

    if (display == NULL)
    {
        //fprintf(stderr, "can't open display `%s'.\n", displayname);
    	return 1;
    }

    if (window == 0) {
        XGetInputFocus(display, &window, &focusreturn);
    }

    if(window == (Window)-1) {
        window = RootWindow(display, 0); // XXX nonzero screens?
    }

    return 0;
}

void SendKeyEvent(KeySym keysym, unsigned int shift) {
    // shift key down if needed
    if (shift == 1) {
        XTestFakeKeyEvent(display, XKeysymToKeycode(display, XK_Shift_L), True, CurrentTime);
        XSync(display, False);
    }

    // press the key
    if (keysym != NoSymbol) {
        XTestFakeKeyEvent(display, XKeysymToKeycode(display, keysym), True, CurrentTime);
        XSync(display, False);
        usleep(1000);
        XTestFakeKeyEvent(display, XKeysymToKeycode(display, keysym), False, CurrentTime);
        XSync(display, False);
    } 

    //shift key up if it was down
    if (shift == 1) {
        XTestFakeKeyEvent(display, XKeysymToKeycode(display, XK_Shift_L), False, CurrentTime);
        XSync(display, False);
    }
}

int SendAltTabKeys() {
    XInitThreads();
    //fprintf(stderr, "C:SendAltTabKeys:start\n");
    if (OpenDisplay() != 0) {
        return 0;
    }
    //fprintf(stderr, "C:SendAltTabKeys:sending keys\n");
    XTestFakeKeyEvent(display, XKeysymToKeycode(display, XK_Alt_L), True, CurrentTime);
    XSync(display, False);
    XTestFakeKeyEvent(display, XKeysymToKeycode(display, XK_Tab), True, CurrentTime);
    XSync(display, False);
    
    usleep( 12000 );

    XTestFakeKeyEvent(display, XKeysymToKeycode(display, XK_Tab), False, CurrentTime);
    XSync(display, False);
    XTestFakeKeyEvent(display, XKeysymToKeycode(display, XK_Alt_L), False, CurrentTime);
    XSync(display, False);
    //(stderr, "C:SendAltTabKeys:keys sent\n");
    XCloseDisplay(display);
    //fprintf(stderr, "C:SendAltTabKeys:display closed\n");
    return 0;
}

int Sendkey(const char *letter, int shift) {
    XInitThreads();
    if (OpenDisplay() != 0) {
        return 1;
    }

    if(shift == 1) {
        shift |= ShiftMask;
    }

    keysym = XStringToKeysym(letter);
    SendKeyEvent(keysym, shift);

    XCloseDisplay(display);

    return 0;
}