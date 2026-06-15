% webserver Static Website with Dynamic API | version 0.1.0 5662b20
% R. S. Doiel
% 2026-01-05

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
webserver start webserver.toml
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
2. webserver matches prefix: /api/
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

Frontend on port 3000, backend on port 9000, webserver on port 8000:

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

