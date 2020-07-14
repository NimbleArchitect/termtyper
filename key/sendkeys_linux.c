#include "sendkeys_linux.h"

#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#  define NeedFunctionPrototypes 1
#  include <X11/Xlib.h>
#  include <X11/keysym.h>
#  include <X11/extensions/XTest.h>
#  include <X11/Xmu/WinUtil.h>
#  if XlibSpecificationRelease != 6
#      error Requires X11R6
#  endif

static int keysym = 0;


Bool xerror = False;

void print_window_name(Display* d, Window w);
Window GetFocusWindow(Display* display);
void SendKeyEvent(Display* display, KeySym keysym, unsigned int shift);

Display* OpenDisplay() {
    Display* display = NULL;
    char *displayname = NULL;

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
        fprintf(stderr, "can't open display `%s'.\n", displayname);
    	return NULL;
    }

    return display;

}

void SendKeyEvent(Display* display, KeySym keysym, unsigned int shift) {
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
    Display* display = NULL;
    //Window window = 0;

    XInitThreads();
    //fprintf(stderr, "C:SendAltTabKeys:start\n");
    display = OpenDisplay();
    if (display == NULL) {
        printf("Failed to open display");
        return 0;
    }

    //window = GetFocusWindow(display);


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
    Display* display = NULL;

    XInitThreads();
    display = OpenDisplay();
    if (display == NULL) {
        printf("Failed to open display");
        return 0;
    }

    //window = GetFocusWindow(display);


    if(shift == 1) {
        shift |= ShiftMask;
    }

    keysym = XStringToKeysym(letter);
    SendKeyEvent(display, keysym, shift);

    XCloseDisplay(display);

    return 0;
}

int handle_error(Display* display, XErrorEvent* error){
  printf("ERROR: X11 error\n");
  xerror = True;
  return 1;
}

void print_window_name(Display* d, Window w){
  XTextProperty prop;
  Status s;

  printf("window name:\n");

  s = XGetWMName(d, w, &prop); // see man
  if(!xerror && s){
    int count = 0, result;
    char **list = NULL;
    result = XmbTextPropertyToTextList(d, &prop, &list, &count); // see man
    if(result == Success){
      printf("\t%s\n", list[0]);
    }else{
      printf("ERROR: XmbTextPropertyToTextList\n");
    }
  }else{
    printf("ERROR: XGetWMName\n");
  }
}


Window GetFocusWindow(Display* display) {
    Window window = 0;
    int focusreturn = 0;

    XGetInputFocus(display, &window, &focusreturn);

    if(window == (Window)-1) {
        window = RootWindow(display, 0); // XXX nonzero screens?
    }

    return window;
}

int LowerWindow() {
    //sdosent work...why?
    Display* display = NULL;
    Window window = 0;

    XInitThreads();
    display = OpenDisplay();
    XSetErrorHandler(handle_error);
    if (display == NULL) {
        printf("Failed to open display");
        return 0;
    }

    window = GetFocusWindow(display);

    //print_window_name(display, window);
    XLowerWindow(display, window);
    usleep( 12000 );
    XSync(display, False);
    XCloseDisplay(display);

    return 0;
}