% webserver URL Redirects | version 0.1.0 5662b20
% R. S. Doiel
% 2026-01-05

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

webserver currently only supports **301 Moved Permanently** redirects.
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

