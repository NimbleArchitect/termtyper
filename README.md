# termtyper

An autotyping assistant that eases the burdern of manually typing long commands

## Features
* full keyboard control, so you never have to reach for the mouse
* custom argument support with default values

## How it works
Press your selected keyboard combination to open termtyper, search for your selected command using the arrow keys to navigate press enter to select the command and enter again to type the command into your active window

## Why?
I have way too many terminal commands to remember and I got fed up with using a text editor to remember all the commands.

## Getting Started
termtype doesn't have any commands stored by default so you will have to add them yourself to do so open the app and click on the new button (alt + n), once saved you can search the commands by name.

## Installation

### Download

clone from GitHub
```
git clone https://github.com/NimbleArchitect/termtyper.git
cd GpioSentry
```


### Building
you will need to insall the following dependicies first
```
sudo dnf install gtk3-devel webkit2gtk3-devel libxkbcommon-devel libXtst-devel libxkbcommon-x11-devel xorg-x11-xkb-utils-devel libpng-devel xsel xclip
go get github.com/zserge/webview
go get github.com/atotto/clipboard
go get github.com/mattn/go-sqlite3
```
then you can cd into the source folder and build with
```
go build
```

## Configuration


### Command line

none yet ;)


## Built With
golang

## Contributing


## Versioning


## Authors

* **NimbleArchitect** - **Initial work**

## Acknowledgments
