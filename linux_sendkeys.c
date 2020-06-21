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

//static char    *progname       = NULL;
static char    *displayname    = NULL;
static Window   window         = 0;
static Display *display        = NULL;
//static char     keyname[1024];
//static int      shift          = 0;
static int      keysym         = 0;


int	MyErrorHandler(Display *my_display, XErrorEvent *event)
{
    fprintf(stderr, "XSendEvent(0x%lx) failed.\n", window);
    return 1;
}

void SendEvent(XKeyEvent *event)
{
    XSync(display, False);
    XSetErrorHandler(MyErrorHandler);
    XSendEvent(display, window, True, KeyPressMask, (XEvent*)event);
    XSync(display, False);
    XSetErrorHandler(NULL);
}

void SendKeyPressedEvent(KeySym keysym, unsigned int shift)
{
    XKeyEvent		event;

    // Meta not yet implemented (Alt used instead ;->)
    int meta_mask=0;

    event.display	= display;
    event.window	= window;
    event.root		= RootWindow(display, 0); // XXX nonzero screens?
    event.subwindow	= None;
    event.time		= CurrentTime;
    event.x		= 1;
    event.y		= 1;
    event.x_root	= 1;
    event.y_root	= 1;
    event.same_screen	= True;
    event.type		= KeyPress;
    event.state		= 0;

    //
    // press down shift keys one at a time...
    //

    if (shift & ShiftMask) {
        //fprintf(stderr, "shift mask" );
        event.keycode = XKeysymToKeycode(display, XK_Shift_L);
        SendEvent(&event);
        event.state |= ShiftMask;
    }
    if (shift & ControlMask) {
        //fprintf(stderr, "ctrl mask" );
        event.keycode = XKeysymToKeycode(display, XK_Control_L);
        SendEvent(&event);
        event.state |= ControlMask;
    }

    if (shift & Mod1Mask) {
        //fprintf(stderr, "alt mask" );
        event.keycode = XKeysymToKeycode(display, XK_Alt_L);
        SendEvent(&event);
        event.state |= Mod1Mask;
    }
    if (shift & meta_mask) {
        event.keycode = XKeysymToKeycode(display, XK_Meta_L);
        SendEvent(&event);
        event.state |= meta_mask;
    }

    //
    //  Now with shift keys held down, send event for the key itself...
    //

    //usleep( 40000 );
    //fprintf(stderr, "sym: 0x%x\n", keysym);
    if (keysym != NoSymbol) {
        event.keycode = XKeysymToKeycode(display, keysym);
        // fprintf(stderr, "code: 0x%x, %d\n", event.keycode, event.keycode );
        SendEvent(&event);

        event.type = KeyRelease;
        SendEvent(&event);
    }

    //
    // Now release all the shift keys...
    //

    if (shift & ShiftMask) {
        event.keycode = XKeysymToKeycode(display, XK_Shift_L);
        SendEvent(&event);
        event.state &= ~ShiftMask;
    }
    if (shift & ControlMask) {
        event.keycode = XKeysymToKeycode(display, XK_Control_L);
        SendEvent(&event);
        event.state &= ~ControlMask;
    }
    if (shift & Mod1Mask) {
        event.keycode = XKeysymToKeycode(display, XK_Alt_L);
        SendEvent(&event);
        event.state &= ~Mod1Mask;
    }
    if (shift & meta_mask) {
        event.keycode = XKeysymToKeycode(display, XK_Meta_L);
        SendEvent(&event);
        event.state &= ~meta_mask;
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
    
    return 0;
}

int Sendkey(const char *letter, int shift, int alt) {
    int Junk;

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

    if(window == 0)
        XGetInputFocus(display, &window, &Junk);

    if(window == (Window)-1)
    {
        window = RootWindow(display, 0); // XXX nonzero screens?
    }

    if(shift == 1) {
        shift |= ShiftMask;
    }
    if(alt == 1) {
        shift |= Mod1Mask;
    }

    keysym = XStringToKeysym(letter);
    SendKeyPressedEvent(keysym, shift);

    XCloseDisplay(display);

    return 0;
}