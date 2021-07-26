
# Installation

*wsfn* is a Go package for building web services.  It includes 
**webserver** as an example command line program run from a 
shell like Bash. 

For all the released version go to the project page on Github and 
click latest release

>    https://github.com/caltechlibrary/wsfn/releases/latest

You will see a list of filenames is in the form of `wsfn-VERSION_NO-PLATFORM_NAME.zip`.

> VERSION_NUMBER is a [symantic version number](http://semver.org/) (e.g. v0.1.2)

> PLATFROM_NAME is a description of a platform (e.g. windows-amd64, macosx-amd64).

Compiled versions are available for Mac OS X (amd64 processor, macosx-amd64), 
Linux (amd64 processor, linux-amd64), Windows (amd64 processor, windows-amd64) 
and Rapsberry Pi (ARM7 processor, raspbian-arm7).

| Platform    | Zip Filename                            |
|-------------|-----------------------------------------|
| Windows     | wsfn-VERSION_NUMBER-windows-amd64.zip |
| Mac OS X    | wsfn-VERSION_NUMBER-macos-amd64.zip  |
| Linux/Intel | wsfn-VERSION_NUMBER-linux-amd64.zip   |
| Raspbery Pi | wsfn-VERSION_NUMBER-raspberry_pi_os-arm7.zip |


## The basic recipe 

+ Download the zip file matching your platform 
+ Unzip it 
+ Copy the contents of the "bin" folder to a folder in your shell's path (e.g. $HOME/bin). 
+ Adjust you PATH if needed
+ test to see if it worked


### Mac OS X

1. Download the zip file
2. Unzip the zip file
3. Copy the executables to $HOME/bin (or a folder in your path)
4. Test

Here's an example of the commands run in the Terminal App after 
downloading the zip file.

```shell
    cd Downloads/
    unzip wsfn-*-macos-amd64.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    webserver -version
```

### Windows

1. Download the zip file
2. Unzip the zip file
3. Copy the executables to $HOME/bin (or a folder in your path)
4. Test

Here's an example of the commands run in from the Bash shell on 
Windows 10 after downloading the zip file.

```shell
    cd Downloads/
    unzip wsfn-*-windows-amd64.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    webserver -version
```


### Linux 

1. Download the zip file
2. Unzip the zip file
3. Copy the executables to $HOME/bin (or a folder in your path)
4. Test

Here's an example of the commands run in from the Bash shell after
downloading the zip file.

```shell
    cd Downloads/
    unzip wsfn-*-linux-amd64.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    webserver -version
```


### Raspberry Pi

Released version is for a Raspberry Pi 2 or later use (i.e. requires ARM 7 support).

1. Download the zip file
2. Unzip the zip file
3. Copy the executables to $HOME/bin (or a folder in your path)
4. Test

Here's an example of the commands run in from the Bash shell after
downloading the zip file.

```shell
    cd Downloads/
    unzip wsfn-*-raspberry_pi_os-arm7.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    webserver -version
```


## Compiling from source

_wsfn_ is "go gettable".  Use the "go get" command to download the dependant packages
as well as _wsfn_'s source code.

```shell
    go get -u github.com/caltechlibrary/pkgassets/...
    go get -u github.com/caltechlibrary/wsfn/...
```

Or clone the repstory and then compile

```shell
    cd
    git clone https://github.com/caltechlibrary/pkgassets src/github.com/caltechlibrary/pkgassets
    cd src/github.com/caltechlibrary/pkgassets
    make
    make test
    make install
    cd
    git clone https://github.com/caltechlibrary/wsfn src/github.com/caltechlibrary/wsfn
    cd src/github.com/caltechlibrary/wsfn
    make
    make test
    make install
```


