

# wsfn

_wsfn_ is a package for common web functions Caltech Library uses 
in various Golang based Caltech Library tools and services. The goal 
is to standardize our handling of web interactions.

+ wsfn.CORSPolicy is a structure for adding CORS headers to a http Handler
+ StaticRouter is a http Handler Function for working with static routes
+ RedirectRouter handles simple target prefix, destination prefix redirect handling
    + AddRedirectRoute adds a target prefix and destination prefix
    + HasRedirectRoutes return true if any redirect routes are configured
    + RedirectRouter uses the internal redirect data to handle redirects
+ ReverseProxy router lets front other web services.

An example **webserver** is also provided to demonstrate some of the
functionality available with this package. The **webserver** is
intended for instructional purposes only and shouldn't be used in a
production setting.

## Release Notes

- version: 0.1.0
- status: inactive
- released: 2026-01-05

- Added reverse proxy support.
- Upgraded Go to v1.26.4


### Authors

- Doiel, R. S.



### Maintainers

- https://orcid.org/0000-0003-0900-6903

## Software Requirements

- Go >= 1.26.4

### Software Suggestions

- GNU Make > 3.8
- Pandoc >= 3.9
- CMTools >= 0.0.45b



## Related resources


- [Download](https://github.com/caltechlibrary/wsfn/releases)
- [Getting Help, Reporting bugs](https://github.com/caltechlibrary/wsfn/issues)
- [LICENSE](https://caltechlibrary.github.io/wsfn/LICENSE)
- [Installation](INSTALL.md)
- [About](about.md)

