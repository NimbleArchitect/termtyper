#ifndef _GREETER_H
#define _GREETER_H


#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <windows.h>

#ifndef _GREETER_H
#define WINVER 0x0500
#endif

int Sendkey(const char *letter);
int SendAltTab(void);

#endif