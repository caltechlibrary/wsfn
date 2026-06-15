% webserver CORS Configuration | version 0.1.0 5662b20
% R. S. Doiel
% 2026-01-05

# NAME

cors - Cross-Origin Resource Sharing configuration

# SYNOPSIS

[cors]
origin = ORIGIN
methods = [ METHOD, ... ]
headers = [ HEADER, ... ]
allow_credentials = BOOLEAN

# DESCRIPTION

Control which origins can access your web services and what HTTP methods
and headers are allowed.

# CONFIGURATION

~~~toml
[cors]
origin = "https://myapp.example.com"
allow_credentials = true
methods = [ "GET", "POST", "PUT", "DELETE", "OPTIONS" ]
headers = [ "Authorization", "Content-Type" ]
exposed_headers = [ "X-Request-ID" ]
max_age = 86400
~~~

# SETTINGS

origin: Allowed origin (string, not * with credentials)
methods: Allowed HTTP methods (array)
headers: Allowed request headers (array)
exposed_headers: Exposed response headers (array)
allow_credentials: Allow credentials (boolean, default: false)
max_age: Preflight cache time in seconds (integer)

# DEFAULT BEHAVIOR

If no [cors] section: origin=*, methods=[GET], allow_credentials=false

# EXAMPLES

Simple CORS for API:

~~~toml
[cors]
origin = "https://myapp.example.com"
methods = [ "GET", "POST", "OPTIONS" ]
~~~

CORS with credentials:

~~~toml
[cors]
origin = "https://myapp.example.com"
allow_credentials = true
methods = [ "GET", "POST", "OPTIONS" ]
headers = [ "Authorization" ]
~~~

Open CORS (not recommended for credentials):

~~~toml
[cors]
origin = "*"
methods = [ "GET", "OPTIONS" ]
~~~

# PREFLIGHT REQUESTS

Automatically handled. OPTIONS must be in methods list.

# CORS AND AUTHENTICATION

With authentication:

1. allow_credentials = true
2. origin cannot be "*"
3. Client must include credentials

# CORS AND REVERSE PROXY

CORS headers added before proxying:

~~~toml
[reverse_proxy]
"/api/" = "http://localhost:9000/"

[cors]
origin = "https://myapp.example.com"
allow_credentials = true
~~~

# TESTING

~~~bash
curl -I http://localhost:8000/api/data \
  -H "Origin: http://example.com"

curl -I -X OPTIONS http://localhost:8000/api/data \
  -H "Origin: http://example.com" \
  -H "Access-Control-Request-Method: POST"
~~~

# TROUBLESHOOTING

CORS Headers Not Present:
- Verify [cors] section in configuration
- Check request includes Origin header
- Ensure origin matches configured origin

Preflight Failing:
- OPTIONS must be in methods
- Requested headers must be in headers

# SEE ALSO

config-file, reverse-proxy, authentication

