# termtyper

An autotyping assistant that eases the burdern of manually typing long commands

## Features
* full keyboard control, so you never have to reach for the mouse
* custom argument support with default values
* copy from the clipboard to specified argument using shortcut keys alt+1, alt+2, etc
* import and export commands to/from json

## How it works
Press your selected keyboard combination to open termtyper, search for your selected command using the arrow keys to navigate press enter to select the command and enter again to type the command into your active window

## Why?
I have way too many terminal commands to remember and I got fed up storing them all in a text editor.

## Getting Started
termtype doesn't have any commands stored by default so you will have to add them yourself to do so open the app and click on the new button (alt + n), once saved you can search the commands by name.

## Installation

### Download

clone from GitHub
```
git clone https://github.com/NimbleArchitect/termtyper.git
cd termtyper
```


### Building
you will need to insall the following dependicies first

#### Fedora
```
sudo dnf install gtk3-devel webkit2gtk3-devel libxkbcommon-devel libXtst-devel libxkbcommon-x11-devel xorg-x11-xkb-utils-devel libpng-devel xsel xclip
go get -d
```

#### Ubuntu
```
sudo apt-get install libwebkit2gtk-4.0
```

then you can cd into the source folder and build with
```
go get -d
go build
```

### Arch
```
sudo pacman -S go pkg-config sqlite gcc
```
then cd into the source folder and run the following
```
export CGO_CFLAGS="-g -O2 -Wno-return-local-addr"
go get -d
go build
```


### Windows

you will need git mingw and golang installed then run the following to build
```
    set GCO_ENABLED=1
    go get -d -v
    go build -ldflags="-H windowsgui" 
```

### OSX


## Configuration


### Command line

-n allows creating items from stdin, use with the following alias to save the previous run command ```alias ns=history |tail -n2 | head -n1 | cut -d " " -f3- | /path/to/termtyper -n```


## Built With
golang

## Contributing


## Versioning


## Authors

* **NimbleArchitect** - **Initial work**

## Acknowledgments
