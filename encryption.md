% webaccess Password Encryption | version 0.1.0 5662b20
% R. S. Doiel
% 2026-01-05

# NAME

encryption - Password encryption options

# SYNOPSIS

encryption = METHOD

# DESCRIPTION

webaccess supports multiple encryption methods for storing user passwords.
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
webaccess init access.toml

# Or manually specify encryption
webaccess init access.toml
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
webaccess init access-new.toml

# Add users with new encryption
webaccess add access-new.toml alice
webaccess add access-new.toml bob

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

