
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "sendkeys_linux.h"

static int keysym = 0;
Bool xerror = False;


Display* OpenDisplay() {
    Display* display = NULL;
    char *displayname = NULL;

    fprintf(stderr, "C:OpenDisplay:start\n");
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
    XSetErrorHandler(handle_error);

    XInitThreads();
    fprintf(stderr, "C:SendAltTabKeys:start\n");
    display = OpenDisplay();
    if (display == NULL) {
        printf("Failed to open display");
        return 0;
    }

    //window = GetFocusWindow(display);


    fprintf(stderr, "C:SendAltTabKeys:sending keys\n");
    XTestFakeKeyEvent(display, XKeysymToKeycode(display, XK_Alt_L), True, CurrentTime);
    XSync(display, False);
    XTestFakeKeyEvent(display, XKeysymToKeycode(display, XK_Tab), True, CurrentTime);
    XSync(display, False);
    
    usleep( 12000 );

    XTestFakeKeyEvent(display, XKeysymToKeycode(display, XK_Tab), False, CurrentTime);
    XSync(display, False);
    XTestFakeKeyEvent(display, XKeysymToKeycode(display, XK_Alt_L), False, CurrentTime);
    XSync(display, False);
    fprintf(stderr, "C:SendAltTabKeys:keys sent\n");
    XCloseDisplay(display);
    fprintf(stderr, "C:SendAltTabKeys:display closed\n");
    return 0;
}

int Sendkey(const char *letter, int shift) {
    Display* display = NULL;

    XSetErrorHandler(handle_error);
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

unsigned int getWindowCount( Display *display, Window parent_window, int depth )
{
    Window  root_return;
    Window  parent_return;
    Window *children_list = NULL;
    Window top_window;
    unsigned int list_length = 0;
    //unsigned int total_list_length = 0;

    // query the window list recursively, until each window reports no sub-windows
    printf( "getWindowCount() - Window %lu\n", parent_window );
    if ( 0 != XQueryTree( display, parent_window, &root_return, &parent_return, &children_list, &list_length ) )
    {
        //printf( "getWindowCount() - %s    %d window handle returned\n", indent, list_length );
        if ( list_length > 0 && children_list != NULL )
        {
            // for ( int i=0; i<list_length; i++)
            // {
            //     // But typically the "top-level" window is not what the user sees, so query any children
            //     // Only the odd window has child-windows.  XEyes does.
            //     if ( children_list[i] != 0 )
            //     {
            //         unsigned int child_length = getWindowCount( display, children_list[i], depth+1 );
            //         total_list_length += child_length;  
            //     }
            // }
            top_window = children_list[list_length -3];
            XRaiseWindow(display, top_window);
            printf( "getWindowCount() - Window %lu\n", top_window );
            XSync(display, False);
            XFree( children_list ); // cleanup
        }
    }

    return list_length; 
}


Window GetFocusWindow(Display* display) {
    Window window = 0;
    int focusreturn = 0;

    XGetInputFocus(display, &window, &focusreturn);

    if(window == (Window)-1) {
        printf("getinputfocus failed, using rootwindo instead");
        window = RootWindow(display, 0); // XXX nonzero screens?
    }

    return window;
}

int LowerWindow(void* ptr) {

    gtk_window_iconify (ptr);

    return 0;
}
