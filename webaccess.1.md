% webaccess(1) webaccess user manual | version 0.1.0 467b3ce
% R. S. Doiel
% 2026-01-05

# NAME

webaccess - User access manager for wsfn webserver

# SYNOPSIS

webaccess [OPTIONS]

webaccess VERB CONFIG_FILE [PARAMETER]

# DESCRIPTION

A nimble user access manager for the wsfn webserver.

webaccess is a command line utility for setting up and managing
user authentication and authorization for web services using the wsfn
web server framework.

# OPTIONS

-help [TOPIC]
: display help (this message) or help for TOPIC

-license
: display license

-version
: display version

-o FILE
: write output to FILE

# VERBS

init CONFIG_FILE
: creates an access control configuration file

add CONFIG_FILE USERNAME
: adds a user and prompts for password

update CONFIG_FILE USERNAME
: updates a user's password

remove CONFIG_FILE USERNAME
: removes a user

list CONFIG_FILE
: lists all users

test CONFIG_FILE USERNAME
: test a login for USERNAME (will prompt for password)

routes VERB CONFIG_FILE [ROUTE ...]
: manage routes (update, list, remove)

# EXAMPLES

Create access control file:

~~~
webaccess init access.toml
~~~

Add a user:

~~~
webaccess add access.toml alice
~~~

Remove a user:

~~~
webaccess remove access.toml bob
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

# TOPICS

Available help topics:

access-control    Access control concepts and configuration
users             Managing users
encryption        Password encryption options
config            Configuration file format

Use 'webaccess help TOPIC' or 'webaccess -help TOPIC' for more information.

