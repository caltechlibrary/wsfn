% webaccess(1) webaccess user manual | version 0.0.12 aee0451
% R. S. Doiel
% 2025-02-27

# NAME

webaccess

# SYNOPSIS

webaccess [OPTIONS]

webaccess VERB CONFIG_FILE [PARAMETER]

# DESCRIPTION

A nimble user access manager for the wsfn webserver.

webaccess is a command line utility for setting up/managing
user access to web services built on wsfn.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-o
: write output to filename


# CONFIG_FILE

webaccess provides a command line interface for managing
an access file. It provides the ability to 
setup users as well as protected routes.

# EXAMPLES

Create an empty "access.toml" file.

~~~
webaccess init access.toml
~~~

Add user id "Jane.Doe" to "access.toml".
The access program prompts for a password. 

~~~
webaccess update access.toml Jane.Doe
~~~

Remove "Jane.Doe" from access.toml.

~~~
webaccess remove access.toml Jane.Doe
~~~

List users defined in access.toml.

~~~
webaccess list access.toml 
~~~

Test a login for Jane.Doe (will prompt for password)

~~~
webaccess test access.toml Jane.Doe
~~~

Routes follow a similar pattern of update, list, remove.
(note you can update or remove more than one route at a time)

~~~
webaccess routes update access.toml "/api/" "/private"

webaccess routes list access.toml

webaccess routes remove access.toml "/private/"
~~~


