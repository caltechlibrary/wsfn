% webserver(1) webserver user manual | version 0.0.14 0b9ce83
% R. S. Doiel
% 2026-01-05

# NAME

webserver

# SYNOPSIS

webserver [OPTIONS]

webserver [VERB PARAMETERS || CONFIG_NAME] [DOCROOT] [URL_TO_LISTEN_ON]

# DESCRIPTION

A nimble web server.

webserver is a command line utility for developing and testing 
static websites.  It uses Go's standard http libraries 
and can supports both http 1 and 2 out of the box.  It 
provides a minimal set of extra features useful for 
developing and testing web services that leverage static 
content. 

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-o
: write output to filename


# CONFIG_FILE

webserver is configured through a configuration file. You can
create an initialization file using the "init" action.
By default the created initialation file is "webserver".

# ACTION

The following actions are available

init
: creates a "webservice.toml" file. This is used by webserver for configuration.

start
: starts up the web service

htdocs
: sets the document root

cert_pem
: set the path to find cert.pem file for TLS

key_pem
: set the path to find the key.pem file for TLS

auth
: set auth type if used, e.g. Basic

access
: sets an external access file. The external access file is managed with the "webaccess" tool.

# EXAMPLES

Run web server using the content in the current directory
(assumes there is no "webserver" file in the working directory).

~~~
webserver start
~~~

Run web server using a specified directory

~~~
   webserver start /www/htdocs
~~~

Running web server using a "/etc/webserver" file for configuration.

~~~
   webserver start /etc/webserver
~~~

Running the web server using the basic setup of "/etc/webserver"
and overriding the default htdocs root and URL listened on

~~~
   webserver start /etc/webserver ./htdocs http://localhost:9011
~~~

Configure your web server with these steps

~~~
   webserver init webserver.toml
   webserver htdocs webserver.toml /var/www/htdocs
   webserver url webserver.toml https://www.example.edu:443
   webserver cert_pem webserver.toml /etc/certs/cert.pem
   webserver key_pem webserver.toml /etc/certs/key.pem
   webserver access webserver.toml /etc/wsfn/access.toml
~~~


