#ifndef _GREETER_H
#define _GREETER_H

#  define NeedFunctionPrototypes 1
#  include <X11/Xlib.h>
#  include <X11/keysym.h>
#  include <X11/extensions/XTest.h>
#  include <X11/Xmu/WinUtil.h>
#  include <gtk/gtk.h>
#  if XlibSpecificationRelease != 6
#      error Requires X11R6
#  endif

int handle_error(Display* display, XErrorEvent* error);
void print_window_name(Display* d, Window w);
Window GetFocusWindow(Display* display);
void SendKeyEvent(Display* display, KeySym keysym, unsigned int shift);
int Sendkey(const char *letter, int shift);
int SendAltTabKeys();
int LowerWindow(void* ptr);

#endif