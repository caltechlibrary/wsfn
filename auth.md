% webserver Authentication | version 0.1.0 5662b20
% R. S. Doiel
% 2026-01-05

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

