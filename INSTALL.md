Installation for development of **wsfn**
===========================================

**wsfn** A golang package providing simple static http services.

Quick install with curl or irm
------------------------------

There is an experimental installer.sh script that can be run with the following command to install latest table release. This may work for macOS, Linux and if you’re using Windows with the Unix subsystem. This would be run from your shell (e.g. Terminal on macOS).

~~~shell
curl https://caltechlibrary.github.io/wsfn/installer.sh | sh
~~~

This will install the programs included in wsfn in your `$HOME/bin` directory.

If you are running Windows 10 or 11 use the Powershell command below.

~~~ps1
irm https://caltechlibrary.github.io//installer.ps1 | iex
~~~

Installing from source
----------------------

### Required software


### Steps

1. git clone https://github.com/caltechlibrary/wsfn
2. Change directory into the `wsfn` directory
3. Make to build, test and install

~~~shell
git clone https://github.com/caltechlibrary/wsfn
cd wsfn
make
make test
make install
~~~

