% webserver Reverse Proxy | version 0.1.0 5662b20
% R. S. Doiel
% 2026-01-05

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
- Ensure backend is accessible from webserver server

# SEE ALSO

config-file, static-website, tls

