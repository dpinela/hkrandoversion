# hkrandoversion

hkrandoversion is an utility for determining the version number of a Hollow Knight
randomizer DLL without having to launch it first.

## How to install

In a terminal, run

    go get github.com/dpinela/hkrandoversion@latest

or get a prebuilt executable from the Releases tab.

## How to use

In a terminal, run

    hkrandoversion path/to/rando.dll

If it recognizes the DLL as a randomizer DLL, it will print its version
in the same format that is displayed in-game.