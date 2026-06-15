% webserver Static Website Serving | version 0.1.0 5662b20
% R. S. Doiel
% 2026-01-05

# NAME

static-website - Serving static HTML/CSS/JS websites

# SYNOPSIS

webserver start [DOCROOT] [URL]

# DESCRIPTION

Serve static websites - HTML, CSS, JavaScript, images, and other assets
with built-in protection against serving hidden files (dot-files).

# BASIC USAGE

Serve current directory:

~~~
webserver start
~~~

Serve specific directory:

~~~
webserver start /var/www/html
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

