
# wsfn

_wsfn_ is a package on common web functions Caltech Library uses 
through several of the Caltech Library tools and services. The goal 
is to standardize our handling of web interactions.

+ wsfn.CORSPolicy is a structure for adding CORS headers to a http Handler
+ StaticRouter is a http Handler Function for working with static routes
+ RedirectRouter handles simple target prefix, destination prefix redirect handling
    + AddRedirectRoute adds a target prefix and destination prefix
    + HasRedirectRoutes return true if any redirect routes are configured
    + RedirectRouter uses the internal redirect data to handle redirects



