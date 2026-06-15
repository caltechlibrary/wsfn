% webaccess Configuration | version 0.1.0 5662b20
% R. S. Doiel
% 2026-01-05

# NAME

config - Configuration file format for webaccess

# SYNOPSIS

webaccess VERB CONFIG_FILE

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
webaccess init access.toml
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

