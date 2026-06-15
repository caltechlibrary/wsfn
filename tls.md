% webserver TLS/HTTPS Configuration | version 0.1.0 5662b20
% R. S. Doiel
% 2026-01-05

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

webserver start webserver.toml
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

