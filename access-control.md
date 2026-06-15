% webaccess Access Control | version 0.1.0 5662b20
% R. S. Doiel
% 2026-01-05

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
webaccess init /etc/wsfn/access.toml
~~~

## Adding Users

~~~bash
webaccess add /etc/wsfn/access.toml alice
~~~

You will be prompted for a password. The password is encrypted and
stored in the file.

## Updating Passwords

~~~bash
webaccess update /etc/wsfn/access.toml alice
~~~

You will be prompted for a new password.

## Removing Users

~~~bash
webaccess remove /etc/wsfn/access.toml alice
~~~

## Listing Users

~~~bash
webaccess list /etc/wsfn/access.toml
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
webaccess init /etc/wsfn/access.toml
~~~

2. Add users:

~~~bash
webaccess add /etc/wsfn/access.toml alice
webaccess add /etc/wsfn/access.toml bob
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

users, encryption, config, webaccess(1)

