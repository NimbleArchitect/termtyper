#include "linux_sendkeys.h"

#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#ifdef WIN32
#  include "windows_x11.h"
#else
#  define NeedFunctionPrototypes 1
#  include <X11/Xlib.h>
#  include <X11/keysym.h>
#  include <X11/extensions/XTest.h>
#  if XlibSpecificationRelease != 6
#      error Requires X11R6
#  endif
#endif

static char *displayname = NULL;
static int keysym = 0;
static Display *display = NULL;
static Window window = 0;


int	KeyErrorHandler(Display *my_display, XErrorEvent *keyevent)
{
    fprintf(stderr, "XSendEvent(0x%lx) failed.\n", window);
    return 1;
}

void SendEvent(XKeyEvent *keyevent)
{
    XSync(display, False);
    XSetErrorHandler(KeyErrorHandler);
    XSendEvent(display, window, True, KeyPressMask, (XEvent*)keyevent);
    XSync(display, False);
    XSetErrorHandler(NULL);
}

void SendKeyEvent(KeySym keysym, unsigned int shift) {
    XKeyEvent keyevent;

    keyevent.display = display;
    keyevent.root = RootWindow(display, 0);
    keyevent.same_screen = True;
    keyevent.state	= 0;
    keyevent.subwindow	= None;
    keyevent.time = CurrentTime;
    keyevent.type = KeyPress;
    keyevent.window = window;
    keyevent.x	= 1;
    keyevent.x_root = 1;
    keyevent.y	= 1;
    keyevent.y_root = 1;

    // shift key down if needed
    if (shift == 1) {
        keyevent.keycode = XKeysymToKeycode(display, XK_Shift_L);
        SendEvent(&keyevent);
        keyevent.state = 1;
    }

    // press the key
    if (keysym != NoSymbol) {
        keyevent.keycode = XKeysymToKeycode(display, keysym);
        SendEvent(&keyevent);
        keyevent.type = KeyRelease;
        SendEvent(&keyevent);
    }

    //shift key up if it was down
    if (shift == 1) {
        keyevent.keycode = XKeysymToKeycode(display, XK_Shift_L);
        SendEvent(&keyevent);
        keyevent.state = 0;
    }
}

int SendAltTab() {
    if(displayname == NULL)
    	displayname = getenv("DISPLAY");

    if(displayname == NULL)
	    displayname = ":0.0";

    display = XOpenDisplay(displayname);

    if(display == NULL)
    {
        fprintf(stderr, "can't open display `%s'.\n", displayname);
    	exit(1);
    }

    XTestFakeKeyEvent(display, XKeysymToKeycode(display, XK_Alt_L), True, CurrentTime);
    XSync(display, False);
    XTestFakeKeyEvent(display, XKeysymToKeycode(display, XK_Tab), True, CurrentTime);
    XSync(display, False);
    usleep( 12000 );
    XTestFakeKeyEvent(display, XKeysymToKeycode(display, XK_Tab), False, CurrentTime);
    XSync(display, False);
    XTestFakeKeyEvent(display, XKeysymToKeycode(display, XK_Alt_L), False, CurrentTime);
    XSync(display, False);

    XCloseDisplay(display);
    return 0;
}

int Sendkey(const char *letter, int shift) {
    int Junk;

    if (displayname == NULL) {
	    displayname = getenv("DISPLAY");
    }

    if (displayname == NULL) {
	    displayname = ":0.0";
    }

    display = XOpenDisplay(displayname);

    if (display == NULL)
    {
        fprintf(stderr, "can't open display `%s'.\n", displayname);
    	exit(1);
    }

    if (window == 0) {
        XGetInputFocus(display, &window, &Junk);
    }

    if(window == (Window)-1) {
        window = RootWindow(display, 0); // XXX nonzero screens?
    }

    if(shift == 1) {
        shift |= ShiftMask;
    }

    keysym = XStringToKeysym(letter);
    SendKeyEvent(keysym, shift);

    XCloseDisplay(display);

    return 0;
}