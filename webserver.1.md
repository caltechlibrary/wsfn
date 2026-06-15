% webserver(1) webserver user manual | version 0.1.0 ce0a61f
% R. S. Doiel
% 2026-01-05

# NAME

webserver - A nimble web server

# SYNOPSIS

webserver [OPTIONS]

webserver [VERB [PARAMETERS] || CONFIG_FILE] [DOCROOT] [URL]

# DESCRIPTION

A nimble web server for developing and testing static websites.

webserver uses Go's standard HTTP libraries and supports both HTTP/1.1 and
HTTP/2 out of the box. It provides a minimal set of features useful for
developing and testing web services that leverage static content.

# OPTIONS

-help [TOPIC]
: display help (this message) or help for TOPIC

-license
: display license

-version
: display version

-o FILE
: write output to FILE

# VERBS

init CONFIG_FILE
: creates a configuration file

start [CONFIG_FILE] [DOCROOT] [URL]
: starts the web server

htdocs CONFIG_FILE DOCROOT
: sets the document root in configuration

url CONFIG_FILE URL
: sets the URL to listen on

cert_pem CONFIG_FILE PATH
: sets TLS certificate file path

key_pem CONFIG_FILE PATH
: sets TLS key file path

auth CONFIG_FILE TYPE
: sets authentication type (e.g., Basic)

access CONFIG_FILE FILE
: sets external access control file

# EXAMPLES

Run web server using the content in the current directory:

~~~
webserver start
~~~

Run web server using a specified directory:

~~~
   webserver start /www/htdocs
~~~

Run with specific configuration file:

~~~
   webserver start /etc/webserver
~~~

# TOPICS

Available help topics:

config-file        Configuration file format and options
reverse-proxy      Forward requests to backend services
static-website     Serving static websites
static-with-api    Example: Static site with dynamic API backend
tls               HTTPS/TLS configuration
auth              Authentication setup
cors              Cross-Origin Resource Sharing
redirects         URL redirection rules

Use 'webserver help TOPIC' or 'webserver -help TOPIC' for more information.

