//
// webserver.go - A simple web server for site development.
// Focus is on demonstrating the functionality provided by wsfn.go
// package.
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2019, Caltech
// All rights not granted herein are expressly reserved by Caltech
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/wsfn"
)

// Flag options
var (
	description = `

SYNOPSIS

	a nimble web server

%s is a command line utility for developing and testing static websites.
It uses Go's standard http libraries and can supports both http 1 and 2
out of the box.  It provides a minimal set of extra features useful
for developing and testing web services that leverage static content.

CONFIGURATION

%s can be configured through a TOML file named "ws.toml". 
The following is an example.

` + "```" + `toml
	# 
	# Minimal webserver.toml file example. 
	#

	# Setting up standard http support
	[http]
	host = "localhost"
	port = "8000"

	# 
	# Setup your document root for the website.
	# It is relative to the current working directory
	# unless a path is fully specified.
	htdocs = "/var/www/htdocs"

` + "```" + `
`

	examples = `
Run web server using the content in the current directory
(assumes there is no "ws.toml" file in the working directory).

   %s

Run web server using a specified directory

   %s /www/htdocs

Running web server using a "/etc/webserver.toml" file for configuration.

   %s /etc/webserver.toml

Running the web server using the basic setup of "/etc/webserver.toml"
and overriding the default htdocs root and URL listened on

   %s /etc/websrever.toml ./htdocs http://localhost:9011
`

	// Standard options
	showHelp         bool
	showVersion      bool
	showLicense      bool
	showExamples     bool
	outputFName      string
	generateMarkdown bool
	generateManPage  bool
	quiet            bool

	// Application option(s)
	initFile bool
)

func main() {
	app := cli.NewCli(wsfn.Version)
	appName := app.AppName()

	// Document non-option parameters
	app.SetParams(`[WEBSERVER_TOML] [DOCROOT] [URL_TO_LISTEN_ON]`)

	// Add Help Docs
	app.AddHelp("license", []byte(fmt.Sprintf(wsfn.LicenseText, appName, wsfn.Version)))
	app.AddHelp("description", []byte(fmt.Sprintf(description, appName, appName, appName)))
	app.AddHelp("examples", []byte(fmt.Sprintf(examples, appName, appName, appName)))

	// Standard Options
	app.BoolVar(&showHelp, "h", false, "display help")
	app.BoolVar(&showHelp, "help", false, "display help")
	app.BoolVar(&showLicense, "l", false, "display license")
	app.BoolVar(&showLicense, "license", false, "display license")
	app.BoolVar(&showVersion, "v", false, "display version")
	app.BoolVar(&showVersion, "version", false, "display version")
	app.BoolVar(&showExamples, "example", false, "display example(s)")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "generate man page")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")

	// Application Options are set via TOML file.
	app.BoolVar(&initFile, "init", false, "generate a default TOML file")

	app.Parse()
	args := app.Args()

	// Setup IO
	var err error

	app.Eout = os.Stderr

	app.Out, err = cli.Create(outputFName, os.Stdout)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(outputFName, app.Out)

	// Process flags and update the environment as needed.
	if generateMarkdown {
		app.GenerateMarkdown(app.Out)
		os.Exit(0)
	}
	if generateManPage {
		app.GenerateManPage(app.Out)
		os.Exit(0)
	}
	if showHelp || showExamples {
		if len(args) > 0 {
			fmt.Fprintln(app.Out, app.Help(args...))
		} else {
			app.Usage(app.Out)
		}
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintln(app.Out, app.License())
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintln(app.Out, app.Version())
		os.Exit(0)
	}

	// Application option(s)
	if initFile {
		fmt.Fprintf(os.Stdout, "%s\n", wsfn.DefaultInit())
		os.Exit(0)
	}

	// setup from command line
	webservice := wsfn.DefaultWebService()
	if _, err := os.Stat("webserver.toml"); os.IsNotExist(err) == false {
		webservice, err = wsfn.LoadTOML("webserver.toml")
		if err != nil {
			fmt.Fprintf(os.Stderr, "toml parse %s, %s\n", "webserver.toml", err)
			os.Exit(1)
		}
	} else if _, err := os.Stat("webserver.json"); os.IsNotExist(err) == false {
		webservice, err = wsfn.LoadTOML("webserver.json")
		if err != nil {
			fmt.Fprintf(os.Stderr, "json parse %s, %s\n", "webserver.json", err)
			os.Exit(1)
		}
	}
	if len(args) > 0 {
		for _, arg := range args {
			switch {
			case strings.HasSuffix(arg, ".toml"):
				webservice, err = wsfn.LoadTOML(arg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "toml parse %s, %s\n", arg, err)
					os.Exit(1)
				}
			case strings.HasSuffix(arg, ".json"):
				webservice, err = wsfn.LoadJSON(arg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "toml parse %s, %s\n", arg, err)
					os.Exit(1)
				}
			case strings.Contains(arg, "://"):
				u, err := url.Parse(arg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "url parse %s, %s\n", arg, err)
					os.Exit(1)
				}
				host := u.Host
				port := ""
				if strings.Contains(u.Host, ":") {
					p := strings.SplitN(u.Host, ":", 2)
					host, port = p[0], p[1]
				}
				switch u.Scheme {
				case "http":
					webservice.Http.Scheme = u.Scheme
					webservice.Http.Host = host
					webservice.Http.Port = port
				case "https":
					webservice.Https.Scheme = u.Scheme
					webservice.Https.Host = host
					webservice.Https.Port = port
				default:
					fmt.Fprintf(os.Stderr, "Unsupported scheme %q", u.String())
					os.Exit(1)
				}
			default:
				webservice.DocRoot = arg
			}
		}
	}
	// Now we should be ready to run the web server
	if err = webservice.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
