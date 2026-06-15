Installation **wsfn**
============================

**wsfn** _wsfn_ is a package for common web functions Caltech Library uses 
in various Golang based Caltech Library tools and services. The goal 
is to standardize our handling of web interactions.

+ wsfn.CORSPolicy is a structure for adding CORS headers to a http Handler
+ StaticRouter is a http Handler Function for working with static routes
+ RedirectRouter handles simple target prefix, destination prefix redirect handling
    + AddRedirectRoute adds a target prefix and destination prefix
    + HasRedirectRoutes return true if any redirect routes are configured
    + RedirectRouter uses the internal redirect data to handle redirects
+ ReverseProxy router lets front other web services.

An example **webserver** is also provided to demonstrate some of the
functionality available with this package. The **webserver** is
intended for instructional purposes only and shouldn't be used in a
production setting.

Quick install with curl or irm
------------------------------

There is an experimental installer.sh script that can be run with the following command to install latest table release. This may work for macOS, Linux and if you’re using Windows with the Unix subsystem. This would be run from your shell (e.g. Terminal on macOS).

~~~shell
curl https://caltechlibrary.github.io/wsfn/installer.sh | sh
~~~

This will install the programs included in wsfn in your `$HOME/bin` directory.

If you are running Windows 10 or 11 use the Powershell command below.

~~~ps1
irm https://caltechlibrary.github.io/wsfn/installer.ps1 | iex
~~~

### If your are running macOS or Windows

You may get security warnings if you are using macOS or Windows. See the notes for the specific operating system you’re using to fix issues.

- [INSTALL_NOTES_macOS.md](INSTALL_NOTES_macOS.md)
- [INSTALL_NOTES_Windows.md](INSTALL_NOTES_Windows.md)

Installing from source
----------------------

### Required software

- Go >= 1.26.4

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

