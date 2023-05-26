Installation
============

**webaccess** and **webserver** are command line programs demonstrating the features of wsfn Go package. They provide a simple static file web server and access control tool (using BasicAuth). They are targetted at localhost development and not intended to be a production web server.

Quick install with curl
-----------------------

There is an experimental installer.sh script that can be run with the
following command to install lastest table release. This may work for
macOS, Linux and if you're using Windows with the Unix subsystem.

~~~
curl https://caltechlibrary.github.io/wsfn/installer.sh | sh
~~~

Below are generalized instructions for installation of a release.

Compiled version
----------------

Compiled versions are available for macOS (Intel and M1 processors as macos-amd64 or macos-arm64), Linux (amd64 process, linux-amd64), Windows (amd64 and arm64 processor, windows-amd64 and windows-arm64) and Rapsberry Pi (arm7 processor, raspbian-arm7)

VERSION\_NUMBER is a [symantic version number](http://semver.org/) (e.g.
`v0.0.10`)

For all the released version go to the project page on Github and click latest release

> <https://github.com/caltechlibrary/wsfn/releases/latest>

| Platform    | Zip Filename                                 |
|-------------|----------------------------------------------|
| Windows     | wsfn-VERSION_NUMBER-Windows-x86_64.zip |
| Windows     | wsfn-VERSION_NUMBER-Windows-arm64.zip |
| macOS       | wsfn-VERSION_NUMBER-macOS-x86_64.zip  |
| macOS       | wsfn-VERSION_NUMBER-macOS-arm64.zip  |
| Linux/Intel | wsfn-VERSION_NUMBER-Linux-x86_64.zip   |
| Linux/ARM64 | wsfn-VERSION_NUMBER-Linux-aarch64.zip   |
| Raspberry Pi ARM 7 | wsfn-VERSION_NUMBER-RaspberryPiOS-arm7.zip |

The basic recipe
----------------

- Find the Zip file listed matching the architecture you're running
  and download it
      - (e.g. if you're on a Windows 10 laptop/Surface with a Intel
        style CPU you'd choose the Zip file with "Windows-x86_64" in the
        name).
- Download the zip file and unzip the file.
- Copy the contents of the folder named "bin" to a folder that is in
  your path
      - (e.g. "\$HOME/bin" is common).
- Adjust your PATH if needed
      - (e.g. `export PATH="\$HOME/bin:\$PATH"`)
- Test

### macOS

1.  Download the zip file
2.  Unzip the zip file
3.  Copy the executables to $HOME/bin (or a folder in your path)
4.  Make sure the new location in in our path
5.  Test

Here's an example of the commands run in the Terminal App after downloading the zip file.

#### Intel Hardware

``` shell
    cd Downloads/
    unzip wsfn-*-macOS-x86_64.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    webserver -version
```

#### M1 (ARM64) Hardware

``` shell
    cd Downloads/
    unzip wsfn-*-macOS-arm64.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    webserver -version
```


### Windows

1.  Download the zip file
2.  Unzip the zip file
3.  Copy the executables to $HOME/bin (or a folder in your path)
4.  Test

Here's an example of the commands run in from the Bash shell on Windows 10 after downloading the zip file.

#### Intel Hardware

Most machines running Windows in 2023 are running using Intel style processors. The exceptions are some developer boxes targetting ARM CPU and some surface tablets that run on ARM processors.

``` shell
    cd Downloads/
    unzip wsfn-*-Windows-x86_64.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    webserver -version
```


#### ARM64 Hardware

``` shell
    cd Downloads/
    unzip wsfn-*-Windows-arm64.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    webserver -version
```

Windows on ARM is relatively rare (in 2023). There are a few Windows for ARM developer boxes out and some Microsoft Surface tablets use an ARM processor.

To find out what type of processor you are running on Windows you can type "systeminfo" into the command prompt or search in the taskbar for "systeminfo" then click on the System Information menu item. You're looking for an entry for "processor". It may look something like

~~~
Processor(s):              1 Processor(s) Installed.
                           [01]: Intel64 Family 6 Model 142 Stepping 12 GenuineIntel ~1803 Mhz
~~

If you see "Intel" somewhere in the description choose the Zip file with "windos" and "amd64" in the name.


### Linux

1.  Download the zip file
2.  Unzip the zip file
3.  Copy the executables to $HOME/bin (or a folder in your path)
4.  Test

Here's an example of the commands run in from the Bash shell after downloading the zip file.

``` shell
    cd Downloads/
    unzip wsfn-*-Linux-x86_64.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    webserver -version
```

### Raspberry Pi

Released version is for a Raspberry Pi 2 or later use (i.e. requires ARM
7 support).

1.  Download the zip file
2.  Unzip the zip file
3.  Copy the executables to $HOME/bin (or a folder in your path)
4.  Test

Here's an example of the commands run in from the Bash shell after downloading the zip file.

``` shell
    cd Downloads/
    unzip wsfn-*-raspbian-arm7.zip
    mkdir -p $HOME/bin
    cp -v bin/* $HOME/bin/
    export PATH=$HOME/bin:$PATH
    webserver -version
```


Compiling from source
---------------------

*wsfn* is "go gettable". Use the "go get" command to download the dependant packages as well as *wsfn*'s source code.

``` shell
    go get -u github.com/caltechlibrary/wsfn/...
```

Or clone the repstory and then compile

``` shell
    cd
    git clone https://github.com/caltechlibrary/wsfn \
        src/github.com/caltechlibrary/wsfn
    cd src/github.com/caltechlibrary/wsfn
    make
    make test
    make install
```
