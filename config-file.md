% webserver Configuration File | version 0.1.0 5662b20
% R. S. Doiel
% 2026-01-05

# NAME

config-file - Configuration file format and options

# SYNOPSIS

webserver [VERB] [CONFIG_FILE]

# DESCRIPTION

webserver is configured through TOML or JSON files that define the
document root, networking settings, and feature configurations.

# CONFIGURATION FILE LOCATION

Default files (in current directory): webserver.toml, webserver.json

Specify file on command line:

~~~bash
webserver start /etc/webserver.toml
~~~

# TOML EXAMPLE

~~~toml
htdocs = "./public"
access_file = "/etc/wsfn/access.toml"

[http]
host = "0.0.0.0"
port = "80"

[https]
cert_pem = "/etc/certs/cert.pem"
key_pem = "/etc/certs/key.pem"
port = "443"

[reverse_proxy]
"/api/" = "http://localhost:9000/"

[cors]
origin = "https://example.com"
methods = [ "GET", "POST", "OPTIONS" ]
headers = [ "Authorization", "Content-Type" ]

[content_types]
".xml" = "application/xml"
~~~

# MAIN SETTINGS

htdocs = DOCUMENT_ROOT
access_file = ACCESS_FILE_PATH
redirects_csv = REDIRECTS_FILE_PATH

# HTTP/HTTPS SERVICE

[http]
host = HOSTNAME  # default: localhost
port = PORT      # default: 8000

[https]
cert_pem = CERTIFICATE_PATH  # required
key_pem = PRIVATE_KEY_PATH   # required
port = PORT                  # default: 8443

# FEATURES

[reverse_proxy]
PREFIX = TARGET_URL

[cors]
origin = ORIGIN
methods = [ METHOD, ... ]
headers = [ HEADER, ... ]

[redirects]
SOURCE = DESTINATION

[content_types]
EXTENSION = MIME_TYPE

[access]
auth_type = "Basic"
auth_name = REALM
routes = [ PATH, ... ]

