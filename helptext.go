//
// helptext.go - Help documentation for the wsfn package
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2024, Caltech
// All rights not granted herein are expressly reserved by Caltech
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package wsfn

// Library-level help topics for the wsfn package
// These are end-user documentation topics in Pandoc-compatible Markdown format.

const (
	// WsfnMainHelp is the main help topic for the wsfn library
	WsfnMainHelp = `% {app_name} Web Service Functions Library | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name} - Web Service Functions Library

# SYNOPSIS

Library for building web services in Go

# DESCRIPTION

The wsfn package provides common web service functionality for Caltech Library
projects and applications. It includes handlers, middleware, and utilities for
building robust web services with minimal boilerplate.

# FEATURES

- **Static File Serving**: Safe static file serving with dot-file protection
- **CORS Policy**: Configurable Cross-Origin Resource Sharing support
- **HTTP Redirects**: Flexible URL redirect management
- **Basic Authentication**: User authentication with multiple encryption options
- **Reverse Proxy**: Forward requests to backend services
- **Request Logging**: Built-in request logging middleware
- **Safe File System**: Protection against serving hidden/dot files
- **Content Type Management**: Custom MIME type mappings

# PACKAGE STRUCTURE

The wsfn package provides the following main components:

## WebService

The central configuration structure that defines a web service with:
- Document root for static files
- HTTP and/or HTTPS service configuration
- Authentication settings
- CORS policy
- Content type mappings
- Redirect rules
- Reverse proxy routes

## Middleware Handlers

Chainable HTTP handlers for common web service needs:
- StaticRouter: Static file handling with dot-file protection
- AccessHandler: Basic authentication
- RedirectRouter: URL redirects
- ReverseProxyRouter: Request forwarding to backends
- RequestLogger: Request logging
- CORSPolicy.Handler: CORS header management

# CONFIGURATION

Web services are typically configured through TOML or JSON files.

## TOML Configuration Example

~~~toml
htdocs = "./public"

[http]
host = "localhost"
port = "8000"

[https]
cert_pem = "/etc/certs/cert.pem"
key_pem = "/etc/certs/key.pem"
port = "8443"

[reverse_proxy]
"/api/" = "http://localhost:9000/"

[cors]
origin = "http://example.com"
allow_credentials = true
methods = [ "GET", "POST", "OPTIONS" ]
~~~

## JSON Configuration Example

~~~json
{
  "htdocs": "./public",
  "http": {
    "host": "localhost",
    "port": "8000"
  },
  "reverse_proxy": {
    "/api/": "http://localhost:9000/"
  }
}
~~~

# TOPICS

The following help topics are available for specific features:

reverse-proxy     Forward requests to backend services
static-serving    Serving static files with dot-file protection
authentication    Basic authentication and user management
cors             Cross-Origin Resource Sharing configuration
redirects        URL redirection rules
configuration    Configuration file format and options
tls              HTTPS/TLS configuration

Use these topics with your application's help system for detailed information.
`

	// ReverseProxyTopicHelp provides detailed documentation for the reverse proxy feature
	ReverseProxyTopicHelp = `% {app_name} Reverse Proxy | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

reverse-proxy - Forward requests to backend services

# SYNOPSIS

[reverse_proxy]
PREFIX = TARGET_URL

# DESCRIPTION

The reverse proxy feature allows {app_name} to forward HTTP requests matching
specific URL prefixes to other backend web services. This is useful for:

* Serving a static website while proxying API requests to a separate backend
* Consolidating multiple services under a single domain
* Adding authentication/authorization in front of backend services
* Testing frontend code against development backend services
* Creating API gateways

When a request matches a configured prefix, {app_name} strips the prefix
from the URL path and forwards the request to the target URL. The backend's
response is then returned to the client unchanged.

# PATH HANDLING

The prefix is stripped before forwarding. For example, with the configuration:

~~~toml
[reverse_proxy]
"/api/" = "http://backend:9000/"
~~~

The following request routing occurs:

| Incoming Request | Backend Request |
|-----------------|-----------------|
| /api/users      | /users          |
| /api/users/123 | /users/123      |
| /api/search?q=x | /search?q=x      |

# CONFIGURATION

Add a [reverse_proxy] section to your configuration file with prefix-to-URL
mappings:

~~~toml
[reverse_proxy]
"/api/" = "http://localhost:9000/"
"/auth/" = "http://localhost:9001/"
"/search/" = "http://search.example.com/"
~~~

Multiple prefixes can be configured. The **longest matching prefix** is used
when multiple prefixes could match a request.

# PREFIX RULES

- Prefixes must not be empty strings
- Prefixes should typically end with "/" (recommended but not required)
- Prefixes must not overlap (e.g., "/api" and "/api/users" cannot both exist)
- Target URLs must use http:// or https:// scheme

# VALIDATION

The following validation is performed when adding reverse proxy routes:

* **Empty prefix**: Rejected with error
* **Empty target URL**: Rejected with error
* **Invalid URL scheme**: Only http:// and https:// are accepted
* **Route collision**: Prefixes that overlap are rejected

# EXAMPLE: Static Website with Dynamic API

A common deployment pattern combines static file serving with reverse proxy:

## Configuration

~~~toml
htdocs = "./public"

[http]
host = "localhost"
port = "8000"

[reverse_proxy]
"/api/" = "http://localhost:9000/"
~~~

## Result

With this configuration:

| URL | Handler | Backend |
|-----|---------|---------|
| / | Static files | ./public/index.html |
| /about.html | Static files | ./public/about.html |
| /css/style.css | Static files | ./public/css/style.css |
| /api/users | Reverse proxy | http://localhost:9000/users |
| /api/posts/123 | Reverse proxy | http://localhost:9000/posts/123 |

## Request Flow

1. Browser requests: GET /api/users
2. {app_name} matches prefix: /api/
3. Strips prefix: /users
4. Forwards to: http://localhost:9000/users
5. Returns backend response to browser

# QUERY PARAMETERS

Query parameters are preserved exactly as received:

* Incoming: /api/search?q=test&page=1
* Backend: /search?q=test&page=1

# HTTP METHODS

All HTTP methods are supported and forwarded unchanged:
- GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
- Custom methods are also forwarded

# HEADERS

## Request Headers

All request headers are forwarded to the backend. You can use this for:
- Authentication tokens
- Content negotiation
- Custom application headers

## Response Headers

All response headers from the backend are returned to the client, with
the exception of headers that would conflict with the proxy operation.

# HTTPS BACKENDS

To proxy to HTTPS backends:

~~~toml
[reverse_proxy]
"/api/" = "https://api.example.com/"
~~~

TLS verification is performed according to Go's standard library behavior.

# LOAD BALANCING

For simple load balancing, proxy to a load balancer or use multiple backends:

~~~toml
[reverse_proxy]
"/api/" = "http://load-balancer:9000/"
~~~

# MULTIPLE BACKENDS

You can configure multiple backend services:

~~~toml
[reverse_proxy]
"/api/" = "http://api-service:9000/"
"/auth/" = "http://auth-service:9001/"
"/images/" = "http://image-service:9002/"
~~~

Each prefix can point to a different backend service.

# PERFORMANCE

Reverse proxy requests are handled efficiently:
- Connection pooling to backends
- Streaming request and response bodies
- Minimal memory overhead

# TROUBLESHOOTING

## Connection Refused

If the backend service is not running, clients will receive a 502 Bad Gateway
response. Ensure your backend services are running and accessible.

## Timeout Issues

Long-running requests may timeout. Consider:
- Increasing backend service timeouts
- Using connection pooling
- Optimizing backend response times

# SEE ALSO

static-serving, configuration, tls
`

	// StaticServingTopicHelp provides documentation for static file serving
	StaticServingTopicHelp = `% {app_name} Static File Serving | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

static-serving - Serve static files with dot-file protection

# SYNOPSIS

htdocs = DOCUMENT_ROOT

# DESCRIPTION

The static file serving feature provides efficient HTTP service for static
content such as HTML files, CSS stylesheets, JavaScript, images, and other
assets. It includes built-in protection against serving hidden files (dot-files).

# DOCUMENT ROOT

The document root is specified by the "htdocs" configuration setting:

~~~toml
htdocs = "./public"
~~~

If not specified, the current working directory is used.

The document root can also be specified via command line:

~~~bash
webserver start /var/www/html
~~~

# DOT-FILE PROTECTION

By default, {app_name} **does not serve** files or directories that start with
a dot (.). This includes:

* .htaccess files
* .git directories
* .env files
* Any file/directory with dot prefix

This protects sensitive files from accidental exposure.

## Configuration

Dot-file protection is always enabled and cannot be disabled.

## Custom Protection

For additional protection, use the access control features to restrict
access to specific paths.

# CONTENT TYPES

Common file extensions are mapped to appropriate MIME types:

| Extension | MIME Type |
|-----------|-----------|
| .html | text/html |
| .css | text/css |
| .js, .mjs | text/javascript |
| .json | application/json |
| .png | image/png |
| .jpg, .jpeg | image/jpeg |
| .gif | image/gif |
| .svg | image/svg+xml |
| .txt | text/plain |

## Custom Content Types

Add custom MIME type mappings in your configuration:

~~~toml
[content_types]
".xml" = "application/xml"
".rss" = "application/rss+xml"
".webmanifest" = "application/manifest+json"
~~~

# GZIPPED FILES

Files with .gz extension are served with Content-Encoding: gzip header:

* .json.gz → application/json with gzip encoding
* .js.gz → text/javascript with gzip encoding

# WEB ASSEMBLY

WebAssembly (.wasm) files are served with the correct MIME type:

~~~toml
[content_types]
".wasm" = "application/wasm"
~~~

# INDEX FILES

When a directory is requested, {app_name} does **not** automatically serve
index.html. You must request the file explicitly.

To add index file support, configure redirects:

~~~csv
# In redirects.csv
"/","/index.html"
~~~

# CACHING

Static files are served without Cache-Control headers by default. Add
them using the CORS configuration or a reverse proxy.

# PERFORMANCE

Static file serving is highly optimized:
- Uses Go's standard library http.FileServer
- Efficient file descriptor usage
- Minimal memory allocation
- Supports HTTP/2

# EXAMPLES

## Basic Static Site

~~~toml
htdocs = "./public"

[http]
port = "8000"
~~~

## Static Site with Custom Types

~~~toml
htdocs = "./public"

[http]
port = "8000"

[content_types]
".xml" = "application/xml"
".webmanifest" = "application/manifest+json"
~~~

# SEE ALSO

reverse-proxy, configuration, cors
`

	// AuthenticationTopicHelp provides documentation for basic authentication
	AuthenticationTopicHelp = `% {app_name} Authentication | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

authentication - Basic HTTP authentication for web services

# SYNOPSIS

[access]
auth_type = "Basic"
auth_name = "Protected Area"

# DESCRIPTION

The authentication feature provides Basic HTTP authentication for protecting
access to your web services. Users are defined in an access control file
and can be managed using the webaccess command line tool.

# AUTHENTICATION TYPES

Currently, only **Basic** authentication is supported.

## Basic Authentication

Basic authentication uses username/password credentials encoded in the
Authorization header. It is recommended for development and testing, or
for services behind additional protection (TLS, VPN, firewall).

# CONFIGURATION

## Web Service Configuration

In your webservice.toml:

~~~toml
access_file = "/etc/wsfn/access.toml"

[access]
auth_type = "Basic"
auth_name = "Restricted Area"
~~~

## Access Control File

The access control file defines users and their credentials:

~~~toml
auth_type = "Basic"
auth_name = "Restricted Area"
encryption = "argon2id"

[access.alice]
salt = [ ... ]
key = [ ... ]

[access.bob]
salt = [ ... ]
key = [ ... ]
~~~

# USER MANAGEMENT

Use the webaccess command to manage users:

## Add a User

~~~bash
webaccess add access.toml alice
# Prompts for password
~~~

## Update User Password

~~~bash
webaccess update access.toml alice
# Prompts for new password
~~~

## Remove a User

~~~bash
webaccess remove access.toml alice
~~~

## List Users

~~~bash
webaccess list access.toml
~~~

# PASSWORD ENCRYPTION

The following encryption methods are supported:

## argon2id (Recommended)

The default and most secure option. Uses the Argon2id algorithm which is
the winner of the Password Hashing Competition.

~~~toml
encryption = "argon2id"
~~~

## pbkdf2

Password-Based Key Derivation Function 2. NIST-approved.

~~~toml
encryption = "pbkdf2"
~~~

## sha512 (Legacy)

SHA-512 hash. Not recommended for new deployments as it lacks salting
in some configurations.

~~~toml
encryption = "sha512"
~~~

## md5 (Legacy - Not Recommended)

MD5 hash. **Not secure** - included for backward compatibility only.

~~~toml
encryption = "md5"
~~~

# PROTECTING SPECIFIC PATHS

By default, authentication applies to all requests if configured.
Use the routes setting to restrict authentication to specific paths:

~~~toml
[access]
auth_type = "Basic"
auth_name = "Admin Area"
routes = [ "/admin/", "/settings/" ]
~~~

With this configuration, only requests to /admin/ and /settings/ require
authentication.

# SECURITY CONSIDERATIONS

## Always Use HTTPS

Basic authentication sends credentials in base64-encoded form. Without
HTTPS, credentials can be intercepted. Always use TLS:

~~~toml
[https]
cert_pem = "/etc/certs/cert.pem"
key_pem = "/etc/certs/key.pem"
port = "443"

[access]
auth_type = "Basic"
auth_name = "Secure Area"
~~~

## Password Strength

Use strong, unique passwords for each user. The encryption methods
(especially argon2id) protect against brute force attacks, but strong
passwords are still essential.

## Access Control File Permissions

The access control file contains sensitive information (password hashes).
Ensure it has restricted permissions:

~~~bash
chmod 600 access.toml
chown root:root access.toml
~~~

## Rate Limiting

Consider adding rate limiting for authentication endpoints to prevent
brute force attacks.

# EXAMPLE: Protected API with Reverse Proxy

~~~toml
htdocs = "./public"
access_file = "/etc/wsfn/access.toml"

[https]
port = "443"

[access]
auth_type = "Basic"
auth_name = "API"

[reverse_proxy]
"/api/" = "http://localhost:9000/"
~~~

With this configuration:
- All requests require authentication
- Static files are served from ./public
- /api/ requests are proxied to the backend
- Both use HTTPS for security

# SEE ALSO

configuration, tls, webaccess
`

	// CorsTopicHelp provides documentation for CORS configuration
	CorsTopicHelp = `% {app_name} CORS Configuration | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

cors - Cross-Origin Resource Sharing configuration

# SYNOPSIS

[cors]
origin = ORIGIN
methods = [ METHOD, ... ]
headers = [ HEADER, ... ]

# DESCRIPTION

Cross-Origin Resource Sharing (CORS) is a mechanism that allows restricted
resources on a web page to be requested from another domain outside the
domain from which the first resource was served.

The CORS configuration in {app_name} allows you to control which origins
can access your web services and what HTTP methods and headers are allowed.

# CONFIGURATION

Add a [cors] section to your webservice.toml:

~~~toml
[cors]
origin = "http://example.com"
allow_credentials = true
methods = [ "GET", "POST", "OPTIONS" ]
headers = [ "Authorization", "Content-Type" ]
exposed_headers = [ "X-Custom-Header" ]
max_age = 86400
~~~

# SETTINGS

## origin

Specifies the allowed origin(s) for CORS requests. Can be:

* A single origin: "http://example.com"
* All origins: "*" (not recommended for credentials)
* Multiple origins: Not directly supported; use a proxy

~~~toml
origin = "http://example.com"
~~~

## allow_credentials

Whether to allow credentials (cookies, HTTP authentication) in CORS requests.

~~~toml
allow_credentials = true
~~~

## methods

List of allowed HTTP methods.

~~~toml
methods = [ "GET", "POST", "PUT", "DELETE", "OPTIONS" ]
~~~

## headers

List of allowed request headers.

~~~toml
headers = [ "Authorization", "Content-Type", "X-CSRF-Token" ]
~~~

## exposed_headers

List of response headers that can be exposed to the client.

~~~toml
headers = [ "X-Custom-Header", "X-Request-ID" ]
~~~

## max_age

How long (in seconds) the results of a preflight request can be cached.

~~~toml
max_age = 86400  # 24 hours
~~~

# DEFAULT CORS BEHAVIOR

If no [cors] section is specified, {app_name} applies a minimal CORS policy:

- No origin restriction (accepts all origins)
- Only GET method allowed
- No custom headers allowed
- No credentials allowed

# EXAMPLES

## Simple CORS for API

Allow requests from a specific frontend:

~~~toml
[cors]
origin = "https://myapp.example.com"
allow_credentials = true
methods = [ "GET", "POST", "PUT", "DELETE" ]
headers = [ "Authorization", "Content-Type" ]
~~~

## Open CORS for Public API

Allow requests from any origin (not recommended for sensitive data):

~~~toml
[cors]
origin = "*"
methods = [ "GET", "OPTIONS" ]
headers = [ "Content-Type" ]
~~~

## Complex CORS with Custom Headers

Support custom application headers:

~~~toml
[cors]
origin = "https://app.example.com"
allow_credentials = true
methods = [ "GET", "POST", "OPTIONS" ]
headers = [ "Authorization", "X-API-Key", "X-Request-ID" ]
exposed_headers = [ "X-RateLimit-Limit", "X-RateLimit-Remaining" ]
max_age = 3600
~~~

# PREFLIGHT REQUESTS

{app_name} automatically handles OPTIONS requests for preflight CORS
negotiation. No additional configuration is needed.

# CORS AND AUTHENTICATION

When using CORS with authentication:

1. Set allow_credentials = true
2. Ensure the origin is not "*"
3. The client must include credentials in the request

~~~toml
[cors]
origin = "https://app.example.com"
allow_credentials = true
~~~

# TESTING CORS

Use curl to test CORS headers:

~~~bash
# Simple request
curl -I -X OPTIONS http://localhost:8000/api/data \
  -H "Origin: http://example.com" \
  -H "Access-Control-Request-Method: GET"

# With credentials
curl -I -X OPTIONS http://localhost:8000/api/data \
  -H "Origin: http://example.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: authorization" \
  --cookie "session=abc123"
~~~

# TROUBLESHOOTING

## CORS Headers Not Present

Ensure:
1. The [cors] section is in your configuration
2. The request includes an Origin header
3. The origin matches the configured origin

## Preflight Failing

Check that:
1. OPTIONS method is in the allowed methods list
2. The requested headers are in the allowed headers list
3. The origin is allowed

# SEE ALSO

configuration, reverse-proxy
`

	// RedirectsTopicHelp provides documentation for URL redirects
	RedirectsTopicHelp = `% {app_name} URL Redirects | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

redirects - Configure URL redirects

# SYNOPSIS

[redirects]
SOURCE = DESTINATION

# DESCRIPTION

The redirects feature allows you to configure automatic HTTP redirects for
specific URL paths or patterns. When a request matches a configured source,
the client is redirected to the destination with a 301 Moved Permanently
status code.

# CONFIGURATION

## Inline Configuration

Add redirects directly in your webservice.toml:

~~~toml
[redirects]
"/old-path/" = "/new-path/"
"/legacy/" = "https://newsite.example.com/"
~~~

## External CSV File

For many redirects, use an external CSV file:

~~~toml
redirects_csv = "redirects.csv"
~~~

In redirects.csv:

~~~csv
# Comment lines start with #
/old-page.html,/new-page.html
/legacy/,/archived/
~~~

# CSV FILE FORMAT

The CSV file should have two columns:
1. Source path/prefix
2. Destination URL or path

## Rules

- Lines starting with # are treated as comments
- Empty lines are ignored
- Each line defines one redirect
- Paths are relative to the document root unless fully qualified

## Example CSV

~~~csv
# Old URLs to new locations
/about.html,/about/
/products/,/catalog/

# External redirects
/old-service/,https://newdomain.com/service/

# HTTP to HTTPS
http://example.com/,https://example.com/
~~~

# PATH HANDLING

## Prefix Matching

Redirects are **prefix-based**. A redirect for /old/ will match:
- /old/
- /old/page.html
- /old/subdir/file.txt

The matching portion is preserved in the destination:

~~~toml
[redirects]
"/old/" = "/new/"
~~~

Request: /old/page.html → Redirects to: /new/page.html

## Exact Matching

To match exactly, include the full path:

~~~toml
[redirects]
"/exact/path" = "/different/path"
~~~

Request: /exact/path → Redirects to: /different/path
Request: /exact/path/extra → No redirect (doesn't match exactly)

# REDIRECT TYPES

{app_name} currently only supports **301 Moved Permanently** redirects.
This tells browsers and search engines that the resource has permanently
moved to the new location.

# USE CASES

## Website Migration

Redirect old URLs to new structure:

~~~csv
/products/,/catalog/products/
/blog/,/news/
~~~

## Domain Changes

Redirect old domain to new domain:

~~~csv
/,https://newdomain.com/
~~~

## HTTP to HTTPS

Redirect HTTP requests to HTTPS:

~~~csv
http://example.com/,https://example.com/
~~~

## Trailing Slash Normalization

Ensure URLs always have trailing slashes:

~~~csv
/about,/about/
/contact,/contact/
~~~

# COMBINING WITH REVERSE PROXY

Redirects are processed **before** reverse proxy routing. This allows:

1. Redirect old API paths to new ones
2. Redirect to external URLs
3. Normalize URLs before proxying

Example:

~~~toml
[redirects]
"/api/v1/" = "/api/v2/"

[reverse_proxy]
"/api/" = "http://backend:9000/"
~~~

Request: /api/v1/users → Redirects to: /api/v2/users
Then proxied to: http://backend:9000/v2/users

# PERFORMANCE

Redirects are:
- Checked in order (first match wins)
- Applied before other middleware
- Very fast (simple string prefix matching)

# TROUBLESHOOTING

## Redirect Loop

Avoid circular redirects:

~~~csv
# BAD - creates loop
/a,/b
/b,/a
~~~

## Not Redirecting

Check that:
1. The source path is correct (case-sensitive)
2. The redirect is loaded (check configuration file)
3. The request matches the prefix exactly

# SEE ALSO

configuration, reverse-proxy
`

	// TlsTopicHelp provides documentation for TLS/HTTPS configuration
	TlsTopicHelp = `% {app_name} TLS/HTTPS Configuration | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

tls - HTTPS/TLS configuration for secure connections

# SYNOPSIS

[https]
cert_pem = CERTIFICATE_PATH
key_pem = PRIVATE_KEY_PATH
host = HOSTNAME
port = PORT

# DESCRIPTION

Transport Layer Security (TLS) provides encrypted communication between
clients and your web server. {app_name} supports TLS for HTTPS connections
using standard X.509 certificates.

# CONFIGURATION

Add an [https] section to your webservice.toml:

~~~toml
[https]
cert_pem = "/etc/certs/cert.pem"
key_pem = "/etc/certs/key.pem"
host = "0.0.0.0"
port = "443"
~~~

# SETTINGS

## cert_pem

Path to the X.509 certificate file in PEM format.

This file should contain:
- The server's certificate
- Any intermediate certificates (for certificate chain)

~~~toml
cert_pem = "/etc/letsencrypt/live/example.com/fullchain.pem"
~~~

## key_pem

Path to the private key file in PEM format.

This file should contain the server's private key.

~~~toml
key_pem = "/etc/letsencrypt/live/example.com/privkey.pem"
~~~

## host

The hostname or IP address to bind to.

~~~toml
host = "0.0.0.0"  # All interfaces
host = "localhost" # Localhost only
host = "example.com" # Specific hostname
~~~

## port

The TCP port to listen on for HTTPS connections.

~~~toml
port = "443"  # Standard HTTPS port
port = "8443" # Alternative port
~~~

# CERTIFICATE OPTIONS

## Self-Signed Certificates (Development)

Generate a self-signed certificate for development:

~~~bash
# Generate key and certificate
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem \
  -days 365 -nodes -subj "/CN=localhost"

# Configuration
[https]
cert_pem = "cert.pem"
key_pem = "key.pem"
port = "8443"
~~~

Access via: https://localhost:8443/

**Note**: Browsers will show security warnings for self-signed certificates.

## Let's Encrypt Certificates (Production)

Use certbot to obtain free certificates:

~~~bash
# Install certbot
sudo apt-get install certbot

# Obtain certificate
sudo certbot certonly --standalone -d example.com

# Configuration
[https]
cert_pem = "/etc/letsencrypt/live/example.com/fullchain.pem"
key_pem = "/etc/letsencrypt/live/example.com/privkey.pem"
port = "443"
~~~

## Renewal

Let's Encrypt certificates expire every 90 days. Set up automatic renewal:

~~~bash
sudo certbot renew --dry-run
~~~

# MIXED HTTP/HTTPS

You can run both HTTP and HTTPS simultaneously:

~~~toml
[http]
port = "80"

[https]
cert_pem = "/etc/certs/cert.pem"
key_pem = "/etc/certs/key.pem"
port = "443"
~~~

# REDIRECT HTTP TO HTTPS

To automatically redirect HTTP requests to HTTPS:

1. Configure both HTTP and HTTPS
2. Add a redirect for the root path

~~~toml
[http]
port = "80"

[https]
cert_pem = "/etc/certs/cert.pem"
key_pem = "/etc/certs/key.pem"
port = "443"

[redirects]
"http://example.com/" = "https://example.com/"
~~~

Or use your web server or load balancer to handle the redirect.

# PERFORMANCE

TLS connections have minimal overhead:
- Session resumption reduces handshake time
- Hardware acceleration (AES-NI) speeds up encryption
- Modern ciphers are efficient

# SECURITY BEST PRACTICES

## Use Strong Ciphers

{app_name} uses Go's default cipher suite which includes strong, modern ciphers.
No additional configuration is needed.

## Keep Certificates Updated

- Monitor certificate expiration dates
- Set up automatic renewal for Let's Encrypt
- Rotate certificates regularly

## Use Proper Permissions

Certificate and key files should have restricted permissions:

~~~bash
chmod 600 /etc/certs/key.pem
chown root:root /etc/certs/key.pem
~~~

## Use TLS 1.2 or Higher

Go's standard library automatically negotiates the highest available TLS version.
It supports TLS 1.0 through 1.3, preferring the highest version.

## Certificate Chain

Include intermediate certificates in cert.pem to ensure all clients can
verify your certificate.

# TROUBLESHOOTING

## SSL Handshake Failed

Check that:
1. Certificate file exists and is readable
2. Private key file exists and is readable
3. Certificate and key match
4. Certificate is not expired
5. Certificate is valid for the hostname

## Port Already In Use

Check for existing services on port 443:

~~~bash
sudo lsof -i :443
sudo netstat -tlnp | grep 443
~~~

## Browser Warnings

If browsers show security warnings:
1. Verify certificate is valid for the hostname
2. Check certificate chain is complete
3. Ensure certificate is not expired
4. Verify certificate is trusted by the browser

# TESTING

Test your HTTPS configuration:

~~~bash
# Check certificate
openssl s_client -connect localhost:443 -servername localhost

# Test with curl
curl -I https://localhost:443/

# Check SSL Labs rating
# (Run from a machine with internet access)
openssl s_client -connect example.com:443 | openssl x509 -noout -text
~~~

# SEE ALSO

configuration, reverse-proxy
`

	// ConfigurationTopicHelp provides documentation for configuration files
	ConfigurationTopicHelp = `% {app_name} Configuration | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

configuration - Web service configuration file format

# SYNOPSIS

Configuration through TOML or JSON files

# DESCRIPTION

{app_name} web services are configured through configuration files that
define the document root, networking settings, and feature configurations.

Both TOML and JSON formats are supported. TOML is recommended for human
editing and JSON for programmatic generation.

# CONFIGURATION FILE LOCATION

## Default Configuration

If no configuration file is specified, {app_name} looks for:
1. webserver.toml in the current directory
2. webserver.json in the current directory

If neither exists, a default configuration is used.

## Specifying Configuration File

On the command line:

~~~bash
webserver start /etc/webserver.toml
webserver start /path/to/config.json
~~~

## Configuration Actions

Use command line actions to modify configuration:

~~~bash
# Create default configuration file
webserver init myconfig.toml

# Set document root
webserver htdocs myconfig.toml /var/www/html

# Set URL to listen on
webserver url myconfig.toml https://example.com:443

# Set TLS certificate
webserver cert_pem myconfig.toml /etc/certs/cert.pem

# Set TLS key
webserver key_pem myconfig.toml /etc/certs/key.pem

# Set access control file
webserver access myconfig.toml /etc/wsfn/access.toml
~~~

# TOML FORMAT

TOML (Tom's Obvious, Minimal Language) is the recommended format for
human-readable configuration.

## Structure

~~~toml
# Comments start with #
key = "value"

[section]
key = "value"

[section.subsection]
key = "value"
~~~

## Data Types

- **Strings**: "value" or 'value'
- **Integers**: 123
- **Floats**: 123.45
- **Booleans**: true, false
- **Arrays**: [ "a", "b", "c" ]
- **Tables**: [section] (nested key-value pairs)

## Example TOML Configuration

~~~toml
htdocs = "./public"
access_file = "/etc/wsfn/access.toml"

[http]
host = "localhost"
port = "8000"

[https]
cert_pem = "/etc/certs/cert.pem"
key_pem = "/etc/certs/key.pem"
port = "443"

[reverse_proxy]
"/api/" = "http://localhost:9000/"

[cors]
origin = "http://example.com"
methods = [ "GET", "POST" ]

[content_types]
".xml" = "application/xml"
~~~

# JSON FORMAT

JSON (JavaScript Object Notation) is supported for programmatic configuration.

## Example JSON Configuration

~~~json
{
  "htdocs": "./public",
  "access_file": "/etc/wsfn/access.toml",
  "http": {
    "host": "localhost",
    "port": "8000"
  },
  "https": {
    "cert_pem": "/etc/certs/cert.pem",
    "key_pem": "/etc/certs/key.pem",
    "port": "443"
  },
  "reverse_proxy": {
    "/api/": "http://localhost:9000/"
  },
  "cors": {
    "origin": "http://example.com",
    "methods": ["GET", "POST"]
  },
  "content_types": {
    ".xml": "application/xml"
  }
}
~~~

# CONFIGURATION SECTIONS

## htdocs

String. The document root for static file serving.

~~~toml
htdocs = "./public"
~~~

## access_file

String. Path to the access control configuration file.

~~~toml
access_file = "/etc/wsfn/access.toml"
~~~

## redirects_csv

String. Path to a CSV file containing redirect rules.

~~~toml
redirects_csv = "redirects.csv"
~~~

## [http]

Table. HTTP service configuration.

| Key | Type | Description | Default |
|-----|------|-------------|---------|
| host | string | Hostname to bind to | localhost |
| port | string | Port to listen on | 8000 |
| scheme | string | Protocol | http |

~~~toml
[http]
host = "0.0.0.0"
port = "8000"
~~~

## [https]

Table. HTTPS service configuration.

| Key | Type | Description | Default |
|-----|------|-------------|---------|
| host | string | Hostname to bind to | localhost |
| port | string | Port to listen on | 8443 |
| cert_pem | string | Certificate file path | - |
| key_pem | string | Private key file path | - |
| scheme | string | Protocol | https |

~~~toml
[https]
cert_pem = "/etc/certs/cert.pem"
key_pem = "/etc/certs/key.pem"
port = "443"
~~~

## [access]

Table. Authentication configuration.

| Key | Type | Description | Default |
|-----|------|-------------|---------|
| auth_type | string | Authentication type | Basic |
| auth_name | string | Realm/description | - |
| routes | array | Paths to protect | - |

~~~toml
[access]
auth_type = "Basic"
auth_name = "Restricted Area"
routes = [ "/admin/", "/settings/" ]
~~~

## [cors]

Table. Cross-Origin Resource Sharing configuration.

See the cors help topic for details.

## [redirects]

Table. URL redirect rules.

| Key | Type | Description |
|-----|------|-------------|
| PREFIX | string | Destination URL |

~~~toml
[redirects]
"/old/" = "/new/"
"/legacy/" = "https://new.example.com/"
~~~

## [reverse_proxy]

Table. Reverse proxy routes.

| Key | Type | Description |
|-----|------|-------------|
| PREFIX | string | Backend URL |

~~~toml
[reverse_proxy]
"/api/" = "http://localhost:9000/"
~~~

## [content_types]

Table. Custom MIME type mappings.

| Key | Type | Description |
|-----|------|-------------|
| EXTENSION | string | MIME type |

~~~toml
[content_types]
".xml" = "application/xml"
".webmanifest" = "application/manifest+json"
~~~

# ENVIRONMENT VARIABLES

Configuration values can reference environment variables:

~~~toml
htdocs = "${WEBSITE_ROOT:-./public}"
~~~

Note: Environment variable support depends on your shell and how the
configuration is loaded.

# INCLUDING OTHER FILES

TOML supports including other files (Go 1.18+):

~~~toml
# main.toml
htdocs = "./public"

[http]
port = "8000"

# Include additional configuration
[include]
path = "reverse_proxy.toml"
~~~

# BEST PRACTICES

## Organization

- Keep configuration files in a known location (/etc, ./config)
- Use descriptive names (webserver-production.toml)
- Document configuration files with comments
- Version control configuration files

## Security

- Restrict permissions on configuration files
- Don't commit secrets to version control
- Use environment variables for sensitive data
- Validate configuration before starting

## Performance

- Minimize the number of configuration files
- Keep frequently changed settings separate
- Use includes for shared configuration

# VALIDATION

Configuration is validated when loaded. Common errors:

- **Invalid port**: Must be a valid port number
- **Missing certificate**: cert_pem must exist for HTTPS
- **Missing key**: key_pem must exist for HTTPS
- **Invalid URL**: reverse_proxy targets must be valid URLs
- **Route collision**: reverse_proxy prefixes must not overlap

# EXAMPLES

## Minimal Configuration

~~~toml
htdocs = "./public"

[http]
port = "8000"
~~~

## Production Configuration

~~~toml
htdocs = "/var/www/html"
access_file = "/etc/wsfn/access.toml"

[http]
host = "0.0.0.0"
port = "80"

[https]
cert_pem = "/etc/letsencrypt/live/example.com/fullchain.pem"
key_pem = "/etc/letsencrypt/live/example.com/privkey.pem"
port = "443"

[reverse_proxy]
"/api/" = "http://localhost:9000/"

[cors]
origin = "https://example.com"
allow_credentials = true

[content_types]
".xml" = "application/xml"
~~~

# SEE ALSO

reverse-proxy, tls, authentication, cors, redirects
`

	// Webserver help text constants for the webserver command
	// These are command-specific help topics

	// WebserverMainHelp is the main help topic for the webserver command
	WebserverMainHelp = `% {app_name}(1) {app_name} user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name} - A nimble web server

# SYNOPSIS

{app_name} [OPTIONS]

{app_name} [VERB [PARAMETERS] || CONFIG_FILE] [DOCROOT] [URL]

# DESCRIPTION

A nimble web server for developing and testing static websites.

{app_name} uses Go's standard HTTP libraries and supports both HTTP/1.1 and
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
{app_name} start
~~~

Run web server using a specified directory:

~~~
   {app_name} start /www/htdocs
~~~

Run with specific configuration file:

~~~
   {app_name} start /etc/{app_name}
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

Use '{app_name} help TOPIC' or '{app_name} -help TOPIC' for more information.
`

	// Webserver help text constants for the webserver command
	// These are command-specific help topics

	// WebserverConfigFileTopicHelp provides configuration file documentation
	WebserverConfigFileTopicHelp = `% {app_name} Configuration File | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

config-file - Configuration file format and options

# SYNOPSIS

{app_name} [VERB] [CONFIG_FILE]

# DESCRIPTION

{app_name} is configured through TOML or JSON files that define the
document root, networking settings, and feature configurations.

# CONFIGURATION FILE LOCATION

Default files (in current directory): webserver.toml, webserver.json

Specify file on command line:

~~~bash
{app_name} start /etc/webserver.toml
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
`

	// WebserverReverseProxyTopicHelp provides reverse proxy documentation
	WebserverReverseProxyTopicHelp = `% {app_name} Reverse Proxy | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

reverse-proxy - Forward requests to backend services

# SYNOPSIS

[reverse_proxy]
PREFIX = TARGET_URL

# DESCRIPTION

Forward HTTP requests matching URL prefixes to other backend services.
Useful for serving static websites while proxying API requests.

# CONFIGURATION

~~~toml
[reverse_proxy]
"/api/" = "http://localhost:9000/"
"/auth/" = "http://localhost:9001/"
~~~

# PATH HANDLING

Prefix is stripped before forwarding:

- /api/users -> proxied to http://backend:9000/users
- Query parameters preserved: /api/search?q=test -> /search?q=test
- All HTTP methods supported

# PREFIX RULES

- Must not be empty
- Should end with / (recommended)
- Must not overlap (e.g., /api and /api/users cannot both exist)
- Target URLs must use http:// or https:// scheme

# VALIDATION

- Empty prefix: Rejected
- Empty target URL: Rejected
- Invalid URL scheme: Only http:// and https:// accepted
- Route collision: Overlapping prefixes rejected

# EXAMPLE: Static Website with Dynamic API

~~~toml
htdocs = "./public"

[http]
port = "8000"

[reverse_proxy]
"/api/" = "http://localhost:9000/"
~~~

Result:
- / -> serves ./public/index.html
- /api/users -> proxied to http://localhost:9000/users

# MULTIPLE BACKENDS

~~~toml
[reverse_proxy]
"/api/" = "http://api-service:9000/"
"/auth/" = "http://auth-service:9001/"
"/images/" = "http://image-service:9002/"
~~~

# HTTPS BACKENDS

~~~toml
[reverse_proxy]
"/api/" = "https://api.example.com/"
~~~

# TROUBLESHOOTING

Backend Connection Refused:
- Verify backend service is running
- Check backend is listening on correct host/port
- Ensure backend is accessible from {app_name} server

# SEE ALSO

config-file, static-website, tls
`

	// WebserverStaticWebsiteTopicHelp provides static website serving documentation
	WebserverStaticWebsiteTopicHelp = `% {app_name} Static Website Serving | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

static-website - Serving static HTML/CSS/JS websites

# SYNOPSIS

{app_name} start [DOCROOT] [URL]

# DESCRIPTION

Serve static websites - HTML, CSS, JavaScript, images, and other assets
with built-in protection against serving hidden files (dot-files).

# BASIC USAGE

Serve current directory:

~~~
{app_name} start
~~~

Serve specific directory:

~~~
{app_name} start /var/www/html
~~~

# DOCUMENT ROOT

Set via command line or configuration:

~~~toml
htdocs = "./public"
~~~

# DOT-FILE PROTECTION

NOT SERVED: .htaccess, .git/, .env, .DS_Store, any file/directory starting with .

Protection is ALWAYS enabled, cannot be disabled.

# INDEX FILES

Directory requests do NOT automatically serve index.html.

Option 1: Request index.html directly

Option 2: Configure redirects:

~~~toml
redirects_csv = "redirects.csv"
~~~

In redirects.csv:

~~~csv
"/","/index.html"
~~~

# CONTENT TYPES

Built-in MIME types: .html, .css, .js, .json, .png, .jpg, .gif, .svg, .txt

Custom MIME types:

~~~toml
[content_types]
".xml" = "application/xml"
".wasm" = "application/wasm"
~~~

# GZIPPED FILES

Files with .gz extension served with Content-Encoding: gzip:

- data.json.gz -> Content-Type: application/json; Content-Encoding: gzip
- script.js.gz -> Content-Type: text/javascript; Content-Encoding: gzip

# PERFORMANCE

Highly optimized using Go's http.FileServer:
- Efficient file descriptor usage
- Minimal memory allocation
- Supports HTTP/2

# SEE ALSO

config-file, reverse-proxy, static-with-api
`

	// WebserverStaticWithApiTopicHelp provides static site with API documentation
	WebserverStaticWithApiTopicHelp = `% {app_name} Static Website with Dynamic API | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

static-with-api - Static website with reverse proxy to API backend

# SYNOPSIS

Static content + Dynamic API on same domain

# DESCRIPTION

Common pattern: serve static website while proxying API requests to
separate backend service. Allows single domain access with efficient
static file serving and specialized backend processing.

# PREREQUISITES

1. Static website files in a directory (e.g., ./public)
2. Backend API service running (e.g., localhost:9000)

# CONFIGURATION

~~~toml
htdocs = "./public"

[http]
port = "8000"

[reverse_proxy]
"/api/" = "http://localhost:9000/"
~~~

Start server:

~~~bash
{app_name} start webserver.toml
~~~

# REQUEST ROUTING

| Request URL | Handler | Backend URL |
|-------------|---------|-------------|
| / | Static files | ./public/index.html |
| /about.html | Static files | ./public/about.html |
| /api/users | Reverse proxy | http://localhost:9000/users |
| /api/posts/123 | Reverse proxy | http://localhost:9000/posts/123 |

# PATH STRIPPING

Prefix is stripped before forwarding:

1. Browser requests: GET /api/users
2. {app_name} matches prefix: /api/
3. Strips prefix: /users
4. Forwards to: http://localhost:9000/users
5. Returns backend response to browser

# QUERY PARAMETERS

Preserved exactly: /api/search?q=test&page=1 -> /search?q=test&page=1

# HTTP METHODS

All methods forwarded: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS

# REQUEST HEADERS

All headers forwarded: Authorization, Content-Type, X-*, Cookie

# RESPONSE HEADERS

All backend response headers returned to client.

# DEVELOPMENT ENVIRONMENT

Frontend on port 3000, backend on port 9000, {app_name} on port 8000:

~~~toml
htdocs = "./public"

[http]
port = "8000"

[reverse_proxy]
"/api/" = "http://localhost:9000/"
"/graphql/" = "http://localhost:4000/"
~~~

Access:
- Static site: http://localhost:8000/
- API: http://localhost:8000/api/
- GraphQL: http://localhost:8000/graphql/

# PRODUCTION ENVIRONMENT

Always use HTTPS in production:

~~~toml
htdocs = "/var/www/html"

[http]
port = "80"

[https]
cert_pem = "/etc/letsencrypt/live/example.com/fullchain.pem"
key_pem = "/etc/letsencrypt/live/example.com/privkey.pem"
port = "443"

[reverse_proxy]
"/api/" = "http://localhost:9000/"
~~~

# MULTIPLE BACKEND SERVICES

~~~toml
[reverse_proxy]
"/api/" = "http://api-service:9000/"
"/auth/" = "http://auth-service:9001/"
"/search/" = "http://search-service:9002/"
~~~

# HTTPS BACKENDS

~~~toml
[reverse_proxy]
"/api/" = "https://api.example.com/"
~~~

# SEE ALSO

reverse-proxy, config-file, tls
`

	// WebserverTlsTopicHelp provides TLS/HTTPS configuration documentation
	WebserverTlsTopicHelp = `% {app_name} TLS/HTTPS Configuration | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

tls - HTTPS/TLS configuration for secure connections

# SYNOPSIS

[https]
cert_pem = CERTIFICATE_PATH
key_pem = PRIVATE_KEY_PATH
host = HOSTNAME
port = PORT

# DESCRIPTION

Provides encrypted communication between clients and web server using
X.509 certificates.

# QUICK START

## Self-Signed Certificate (Development)

~~~bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem \
  -days 365 -nodes -subj "/CN=localhost"

{app_name} start webserver.toml
~~~

Configuration:

~~~toml
[https]
cert_pem = "cert.pem"
key_pem = "key.pem"
port = "8443"
~~~

Access: https://localhost:8443/

**Note**: Browsers will show security warnings for self-signed certificates.

## Let's Encrypt (Production)

~~~bash
sudo apt-get install certbot
sudo certbot certonly --standalone -d example.com
~~~

Configuration:

~~~toml
[https]
cert_pem = "/etc/letsencrypt/live/example.com/fullchain.pem"
key_pem = "/etc/letsencrypt/live/example.com/privkey.pem"
port = "443"
~~~

## Renewal

Let's Encrypt certificates expire every 90 days. Set up automatic renewal:

~~~bash
sudo certbot renew --dry-run
~~~

# MIXED HTTP/HTTPS

Run both simultaneously:

~~~toml
[http]
port = "80"

[https]
cert_pem = "/etc/certs/cert.pem"
key_pem = "/etc/certs/key.pem"
port = "443"
~~~

# REDIRECT HTTP TO HTTPS

To automatically redirect HTTP requests to HTTPS:

Option 1: Use redirects

~~~toml
[http]
port = "80"

[https]
port = "443"

[redirects]
"http://example.com/", "https://example.com/"
~~~

Option 2: Use reverse proxy (nginx, Apache)

# SECURITY BEST PRACTICES

- Use HTTPS for all production sites
- Keep certificates updated (Let's Encrypt: 90 days)
- Restrict permissions: chmod 600 key.pem
- Use strong cipher suites (Go defaults are good)

# TESTING

~~~bash
openssl s_client -connect localhost:443 -servername localhost
curl -I https://localhost:443/
~~~

# TROUBLESHOOTING

SSL Handshake Failed:
- Certificate/key file doesn't exist or isn't readable
- Certificate and key don't match
- Certificate is expired
- Certificate not valid for hostname

Port Already In Use:
~~~bash
sudo lsof -i :443
sudo netstat -tlnp | grep 443
~~~

Browser Warnings:
- Verify certificate is valid for the hostname
- Check certificate chain is complete
- Ensure certificate is not expired
- Verify certificate is trusted by the browser

# SEE ALSO

config-file, reverse-proxy, static-website
`

	// WebserverAuthTopicHelp provides authentication documentation
	WebserverAuthTopicHelp = `% {app_name} Authentication | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

auth - Basic HTTP authentication for web services

# SYNOPSIS

access_file = ACCESS_FILE

[access]
auth_type = "Basic"
auth_name = REALM
routes = [ PATH, ... ]

# DESCRIPTION

Basic HTTP authentication for protecting access to web services.
Users managed via webaccess command line tool.

# USER MANAGEMENT

~~~bash
webaccess init access.toml
webaccess add access.toml alice
webaccess update access.toml alice
webaccess remove access.toml alice
webaccess list access.toml
~~~

# PASSWORD ENCRYPTION

Supported methods (default: argon2id):

- argon2id: Most secure (recommended)
- pbkdf2: NIST-approved
- sha512: Legacy
- md5: Legacy - NOT RECOMMENDED

~~~toml
encryption = "argon2id"
~~~

# PROTECTING PATHS

Protect all requests (default if no routes):

~~~toml
[access]
auth_type = "Basic"
auth_name = "Restricted Area"
~~~

Protect specific paths only:

~~~toml
[access]
auth_type = "Basic"
auth_name = "Admin Area"
routes = [ "/admin/", "/settings/" ]
~~~

# SECURITY

ALWAYS use HTTPS with authentication:

~~~toml
[https]
cert_pem = "/etc/certs/cert.pem"
key_pem = "/etc/certs/key.pem"

[access]
auth_type = "Basic"
auth_name = "Secure Area"
~~~

- Use strong, unique passwords
- Restrict access file permissions: chmod 600 access.toml
- Consider rate limiting

# TESTING

~~~bash
curl -I http://localhost:8000/admin/
curl -I -u username:password http://localhost:8000/admin/
~~~

# SEE ALSO

config-file, webaccess, tls
`

	// WebserverCorsTopicHelp provides CORS configuration documentation
	WebserverCorsTopicHelp = `% {app_name} CORS Configuration | version {version} {release_hash}
% R. S. Doiel
% {release_date}

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
`

	// WebserverRedirectsTopicHelp provides URL redirects documentation
	WebserverRedirectsTopicHelp = `% {app_name} URL Redirects | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

redirects - Configure URL redirects

# SYNOPSIS

[redirects]
SOURCE = DESTINATION

redirects_csv = CSV_FILE

# DESCRIPTION

Automatic HTTP redirects (301 Moved Permanently) for URL paths or patterns.

# CONFIGURATION

## Inline

~~~toml
[redirects]
"/old-path/" = "/new-path/"
"/legacy/" = "https://newsite.example.com/"
~~~

## External CSV File

~~~toml
redirects_csv = "redirects.csv"
~~~

In redirects.csv:

~~~csv
# Comments start with #
/old-page.html,/new-page.html
/legacy/,/archived/
~~~

# PATH MATCHING

Prefix-based matching. Longest match wins.

- /old/ matches /old/, /old/page.html, /old/subdir/file.txt
- Matching portion preserved in destination

## Exact Matching

To match exactly, include the full path:

~~~toml
[redirects]
"/exact/path" = "/different/path"
~~~

Request: /exact/path -> Redirects to: /different/path
Request: /exact/path/extra -> No redirect (doesn't match exactly)

# REDIRECT TYPES

{app_name} currently only supports **301 Moved Permanently** redirects.
This tells browsers and search engines that the resource has permanently
moved to the new location.

# USE CASES

## Website Migration

Redirect old URLs to new structure:

~~~csv
/products/,/catalog/products/
/blog/,/news/
~~~

## Domain Changes

Redirect old domain to new domain:

~~~csv
/,https://newdomain.com/
~~~

## HTTP to HTTPS

~~~toml
[http]
port = "80"

[https]
port = "443"

[redirects]
"http://example.com/": "https://example.com/"
~~~

## Trailing Slash Normalization

Ensure URLs always have trailing slashes:

~~~csv
/about,/about/
/contact,/contact/
~~~

# COMBINING WITH REVERSE PROXY

Redirects processed BEFORE reverse proxy:

~~~toml
[redirects]
"/api/v1/" = "/api/v2/"

[reverse_proxy]
"/api/" = "http://backend:9000/"
~~~

Request: /api/v1/users -> Redirects to /api/v2/users -> Proxied to backend

# TROUBLESHOOTING

Redirect Not Working:
- Check source path spelling (case-sensitive)
- Verify redirect in configuration
- Ensure request path matches source

Redirect Loop:
- Avoid circular: /a -> /b, /b -> /a

# SEE ALSO

config-file, reverse-proxy, static-website
`

	// Webaccess help text constants for the webaccess command
	// These are command-specific help topics

	// WebaccessMainHelp is the main help topic for the webaccess command
	WebaccessMainHelp = `% {app_name}(1) {app_name} user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name} - User access manager for wsfn webserver

# SYNOPSIS

{app_name} [OPTIONS]

{app_name} VERB CONFIG_FILE [PARAMETER]

# DESCRIPTION

A nimble user access manager for the wsfn webserver.

{app_name} is a command line utility for setting up and managing
user authentication and authorization for web services using the wsfn
web server framework.

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
: creates an access control configuration file

add CONFIG_FILE USERNAME
: adds a user and prompts for password

update CONFIG_FILE USERNAME
: updates a user's password

remove CONFIG_FILE USERNAME
: removes a user

list CONFIG_FILE
: lists all users

test CONFIG_FILE USERNAME
: test a login for USERNAME (will prompt for password)

routes VERB CONFIG_FILE [ROUTE ...]
: manage routes (update, list, remove)

# EXAMPLES

Create access control file:

~~~
{app_name} init access.toml
~~~

Add a user:

~~~
{app_name} add access.toml alice
~~~

Remove a user:

~~~
{app_name} remove access.toml bob
~~~

List users defined in access.toml.

~~~
{app_name} list access.toml 
~~~

Test a login for Jane.Doe (will prompt for password)

~~~
{app_name} test access.toml Jane.Doe
~~~

Routes follow a similar pattern of update, list, remove.
(note you can update or remove more than one route at a time)

~~~
{app_name} routes update access.toml "/api/" "/private"

{app_name} routes list access.toml

{app_name} routes remove access.toml "/private/"
~~~

# TOPICS

Available help topics:

access-control    Access control concepts and configuration
users             Managing users
encryption        Password encryption options
config            Configuration file format

Use '{app_name} help TOPIC' or '{app_name} -help TOPIC' for more information.
`

	// WebaccessAccessControlTopicHelp provides access control documentation
	WebaccessAccessControlTopicHelp = `% {app_name} Access Control | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

access-control - Access control concepts and configuration

# SYNOPSIS

access_file = ACCESS_FILE_PATH

[access]
auth_type = "Basic"
auth_name = REALM
routes = [ PATH, ... ]

# DESCRIPTION

The access control feature provides Basic HTTP authentication for protecting
access to your web services. Users and passwords are stored in an access
control file and managed using this tool.

# ACCESS CONTROL FILE

The access control file is a TOML file that contains:
- Authentication settings
- User credentials (encrypted)
- Routes that require authentication

## Default Location

Typically specified in your webserver configuration:

~~~toml
access_file = "/etc/wsfn/access.toml"
~~~

## File Format

~~~toml
auth_type = "Basic"
auth_name = "Restricted Area"
encryption = "argon2id"

[access.alice]
salt = [ 1, 2, 3, ... ]
key = [ 1, 2, 3, ... ]

[access.bob]
salt = [ 1, 2, 3, ... ]
key = [ 1, 2, 3, ... ]
~~~

# AUTHENTICATION TYPES

Currently, only **Basic** authentication is supported.

## Basic Authentication

Standard HTTP Basic authentication using username/password.
Credentials are base64-encoded in the Authorization header.

**Always use HTTPS with Basic authentication!**

# AUTHENTICATION SCOPE

## Protect All Paths

If no routes are specified, authentication is required for ALL requests:

~~~toml
[access]
auth_type = "Basic"
auth_name = "Secure Area"
~~~

## Protect Specific Paths

Only requests matching the specified routes require authentication:

~~~toml
[access]
auth_type = "Basic"
auth_name = "Admin Area"
routes = [ "/admin/", "/settings/" ]
~~~

With this configuration:
- /admin/ and /settings/ require authentication
- All other paths are publicly accessible

# USER MANAGEMENT

Use this tool to manage users in the access control file.

## Creating the File

~~~bash
{app_name} init /etc/wsfn/access.toml
~~~

## Adding Users

~~~bash
{app_name} add /etc/wsfn/access.toml alice
~~~

You will be prompted for a password. The password is encrypted and
stored in the file.

## Updating Passwords

~~~bash
{app_name} update /etc/wsfn/access.toml alice
~~~

You will be prompted for a new password.

## Removing Users

~~~bash
{app_name} remove /etc/wsfn/access.toml alice
~~~

## Listing Users

~~~bash
{app_name} list /etc/wsfn/access.toml
~~~

# SECURITY BEST PRACTICES

## Always Use HTTPS

Basic authentication sends credentials in base64-encoded form.
Without HTTPS, credentials can be intercepted.

**Always configure HTTPS:**

~~~toml
[https]
cert_pem = "/etc/certs/cert.pem"
key_pem = "/etc/certs/key.pem"
port = "443"

[access]
auth_type = "Basic"
auth_name = "Secure Area"
~~~

## Strong Passwords

Use strong, unique passwords:
- Minimum 12 characters
- Mix of upper/lower case letters
- Include numbers and special characters
- Not dictionary words or common phrases

## File Permissions

The access control file contains sensitive information (password hashes).
Restrict permissions:

~~~bash
chmod 600 /etc/wsfn/access.toml
chown root:root /etc/wsfn/access.toml
~~~

## Backup

Regularly backup the access control file. If lost, all user access
will be lost.

## Rate Limiting

Consider rate limiting to prevent brute force attacks.

# EXAMPLES

## Complete Setup

1. Create access control file:

~~~bash
{app_name} init /etc/wsfn/access.toml
~~~

2. Add users:

~~~bash
{app_name} add /etc/wsfn/access.toml alice
{app_name} add /etc/wsfn/access.toml bob
~~~

3. Configure webserver:

~~~toml
access_file = "/etc/wsfn/access.toml"

[access]
auth_type = "Basic"
auth_name = "Restricted Area"
routes = [ "/admin/" ]
~~~

4. Start webserver with HTTPS:

~~~toml
[https]
cert_pem = "/etc/certs/cert.pem"
key_pem = "/etc/certs/key.pem"
port = "443"
~~~

# SEE ALSO

users, encryption, config, {app_name}(1)
`

	// WebaccessUsersTopicHelp provides user management documentation
	WebaccessUsersTopicHelp = `% {app_name} User Management | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

users - Managing users in the access control file

# SYNOPSIS

{app_name} VERB CONFIG_FILE [USERNAME]

# DESCRIPTION

User management commands allow you to add, update, remove, and list
users in an access control file. All operations require the access
control file path as the first argument.

# COMMANDS

## init - Create Access Control File

Creates a new access control file with default settings.

~~~bash
{app_name} init access.toml
~~~

This creates a file with:
- Basic authentication type
- Default encryption (argon2id)
- Empty user list

## add - Add a New User

Adds a new user to the access control file.

~~~bash
{app_name} add access.toml alice
~~~

The command will prompt you for a password for the user.
The password is encrypted and stored securely.

### Multiple Users

~~~bash
{app_name} add access.toml alice
{app_name} add access.toml bob
{app_name} add access.toml charlie
~~~

## update - Update User Password

Updates an existing user's password.

~~~bash
{app_name} update access.toml alice
~~~

You will be prompted for a new password.

## remove - Remove a User

Removes a user from the access control file.

~~~bash
{app_name} remove access.toml alice
~~~

This permanently removes the user. The action cannot be undone.

## list - List All Users

Lists all users in the access control file.

~~~bash
{app_name} list access.toml
~~~

Output format:

~~~
Users in access.toml:
  alice
  bob
  charlie
~~~

# BATCH OPERATIONS

## Add Multiple Users

Use a script to add multiple users:

~~~bash
#!/bin/bash
for user in alice bob charlie; do
  echo "Adding user: $user"
  {app_name} add access.toml "$user"
done
~~~

## Import Users from File

Read usernames from a file and add them:

~~~bash
#!/bin/bash
while read -r user; do
  {app_name} add access.toml "$user"
done < users.txt
~~~

# EXAMPLES

## Add First User

~~~bash
{app_name} init access.toml
{app_name} add access.toml admin
# Enter password when prompted
~~~

## Update Password

~~~bash
{app_name} update access.toml admin
# Enter new password when prompted
~~~

## Remove User

~~~bash
{app_name} remove access.toml olduser
~~~

## List All Users

~~~bash
{app_name} list access.toml
~~~

# SEE ALSO

access-control, encryption, config
`

	// WebaccessEncryptionTopicHelp provides encryption documentation
	WebaccessEncryptionTopicHelp = `% {app_name} Password Encryption | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

encryption - Password encryption options

# SYNOPSIS

encryption = METHOD

# DESCRIPTION

{app_name} supports multiple encryption methods for storing user passwords.
The encryption method determines how passwords are hashed and stored
in the access control file.

# ENCRYPTION METHODS

## argon2id (Recommended)

The default and most secure encryption method. Argon2id is the
winner of the Password Hashing Competition (PHC). It provides
Excellent protection against:
- Brute force attacks
- Timing attacks
- GPU/ASIC attacks

~~~toml
encryption = "argon2id"
~~~

**Parameters used:**
- 1 iteration
- 64MB memory
- 4 threads
- 32 byte output

## pbkdf2

Password-Based Key Derivation Function 2 with HMAC-SHA1.
NIST-approved and widely used.

~~~toml
encryption = "pbkdf2"
~~~

**Parameters used:**
- 4097 iterations
- SHA-1 hash function
- 32 byte output

## sha512 (Legacy)

SHA-512 hash. Not recommended for new deployments.

~~~toml
encryption = "sha512"
~~~

SHA-512 is cryptographically secure but lacks built-in salting
in some configurations.

## md5 (Legacy - NOT RECOMMENDED)

MD5 hash. **Not secure** - included only for backward compatibility.

~~~toml
encryption = "md5"
~~~

MD5 is cryptographically broken and should not be used for
new deployments.

# SETTING ENCRYPTION METHOD

The encryption method is set in the access control file:

~~~bash
# Create file with default (argon2id)
{app_name} init access.toml

# Or manually specify encryption
{app_name} init access.toml
# Then edit access.toml to change encryption method
~~~

# COMPARING METHODS

| Method | Security | Speed | Recommendation |
|--------|----------|-------|----------------|
| argon2id | Excellent | Slow | Recommended |
| pbkdf2 | Good | Medium | Good alternative |
| sha512 | Medium | Fast | Legacy only |
| md5 | Poor | Very Fast | Avoid |

## argon2id

- **Security**: Excellent
- **Memory**: 64MB (resistant to GPU/ASIC attacks)
- **Iterations**: 1
- **Threads**: 4 (parallel computation)
- **Output**: 32 bytes
- **Use**: Recommended for new deployments

## pbkdf2

- **Security**: Good
- **Iterations**: 4097
- **Hash**: SHA-1
- **Output**: 32 bytes
- **Use**: Legacy systems, NIST compliance

## sha512

- **Security**: Medium (no salt in some configs)
- **Output**: 64 bytes
- **Use**: Legacy compatibility only

## md5

- **Security**: Poor (cryptographically broken)
- **Output**: 16 bytes
- **Use**: Avoid - for backward compatibility only

# MIGRATING ENCRYPTION

To change the encryption method for existing users:

1. Create a new access control file with the new encryption method
2. Add all users to the new file
3. Replace the old file with the new one

~~~bash
# Create new file
{app_name} init access-new.toml

# Add users with new encryption
{app_name} add access-new.toml alice
{app_name} add access-new.toml bob

# Replace old file
mv access-new.toml access.toml
~~~

# VERIFYING ENCRYPTION

Check the encryption method in the access control file:

~~~bash
grep encryption access.toml
~~~

# SEE ALSO

access-control, users, config
`

	// WebaccessConfigTopicHelp provides configuration documentation
	WebaccessConfigTopicHelp = `% {app_name} Configuration | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

config - Configuration file format for webaccess

# SYNOPSIS

{app_name} VERB CONFIG_FILE

# DESCRIPTION

The webaccess tool uses TOML configuration files to store user
credentials and authentication settings. This topic describes
the configuration file format and options.

# CONFIGURATION FILE

The configuration file is a TOML file that contains authentication
settings and user credentials.

## Default Location

No default location. You must specify the configuration file for
all commands.

## Creating a Configuration File

~~~bash
{app_name} init access.toml
~~~

This creates a new file with default settings.

# FILE FORMAT

## Complete Example

~~~toml
auth_type = "Basic"
auth_name = "Restricted Area"
encryption = "argon2id"

[access.alice]
salt = [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32 ]
key = [ 1, 2, 3, ... ]

[access.bob]
salt = [ 1, 2, 3, ... ]
key = [ 1, 2, 3, ... ]
~~~

## Sections

### Main Section

| Key | Type | Description | Default |
|-----|------|-------------|---------|
| auth_type | string | Authentication type | "Basic" |
| auth_name | string | Realm/description | "" |
| encryption | string | Password encryption method | "argon2id" |

### User Sections

Each user has a section under [access]:

~~~toml
[access.USERNAME]
salt = [ ... ]
key = [ ... ]
~~~

# MAIN SETTINGS

## auth_type

The authentication type to use.

**Supported values:**
- "Basic" - HTTP Basic authentication (currently the only supported type)

~~~toml
auth_type = "Basic"
~~~

## auth_name

The realm or description for authentication. Displayed in authentication
prompts in browsers.

~~~toml
auth_name = "Restricted Area"
auth_name = "Admin Console"
auth_name = "Please authenticate"
~~~

## encryption

The password encryption method to use.

**Supported values:**
- "argon2id" - Recommended (default)
- "pbkdf2" - NIST-approved
- "sha512" - Legacy
- "md5" - Legacy, not recommended

~~~toml
encryption = "argon2id"
~~~

# USER SETTINGS

Each user section contains:

## salt

Random bytes used for password hashing. Must be 32 bytes for argon2id.

~~~toml
[access.alice]
salt = [ 1, 2, 3, ... 32 bytes ... ]
~~~

## key

The hashed password. Length varies by encryption method.

~~~toml
[access.alice]
key = [ 1, 2, 3, ... hashed password ... ]
~~~

# EXAMPLES

## Minimal Configuration

~~~toml
auth_type = "Basic"
auth_name = "Secure Area"
encryption = "argon2id"
~~~

## Configuration with Users

~~~toml
auth_type = "Basic"
auth_name = "Admin Area"
encryption = "argon2id"

[access.alice]
salt = [ ... ]
key = [ ... ]

[access.bob]
salt = [ ... ]
key = [ ... ]
~~~

# FILE LOCATIONS

## Recommended Locations

- /etc/wsfn/access.toml (system-wide)
- ./access.toml (project-specific)
- ~/.config/wsfn/access.toml (user-specific)

## Security Considerations

- Restrict file permissions: chmod 600
- Restrict ownership: chown root:root (for system files)
- Backup regularly
- Don't commit to version control

# TOML SYNTAX

TOML (Tom's Obvious, Minimal Language) is easy to read and write.

## Basic Syntax

~~~toml
# Comments start with #
key = "value"

[section]
key = "value"

[section.subsection]
key = "value"
~~~

## Data Types

- Strings: "value" or 'value'
- Integers: 123
- Floats: 123.45
- Booleans: true, false
- Arrays: [ 1, 2, 3 ]
- Tables: [section]

# SEE ALSO

access-control, users, encryption
`
)
