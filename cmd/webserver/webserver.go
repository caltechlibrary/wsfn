//
// webserver.go - A simple web server for site development.
// Focus is on demonstrating the functionality provided by wsfn.go
// package.
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2023, Caltech
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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/wsfn"

	// 3rd Party packages
	"github.com/BurntSushi/toml"
)

// Flag options
var (
	helpText = `% {app_name}(1) {app_name} user manual | version {version} {release_hash}
% R. S. Doiel
% {release_date}

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS]

{app_name} [VERB PARAMETERS || CONFIG_NAME] [DOCROOT] [URL_TO_LISTEN_ON]

# DESCRIPTION

A nimble web server.

{app_name} is a command line utility for developing and testing 
static websites.  It uses Go's standard http libraries 
and can supports both http 1 and 2 out of the box.  It 
provides a minimal set of extra features useful for 
developing and testing web services that leverage static 
content. 

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

{app_name} is configured through a configuration file. You can
create an initialization file using the "init" action.
By default the created initialation file is "{app_name}".

# ACTION

The following actions are available

init
: creates a "webservice.toml" file. This is used by {app_name} for configuration.

start
: starts up the web service

htdocs
: sets the document root

cert_pem
: set the path to find cert.pem file for TLS

key_pem
: set the path to find the key.pem file for TLS

auth
: set auth type if used, e.g. Basic

access
: sets an external access file. The external access file is managed with the "webaccess" tool.

# EXAMPLES

Run web server using the content in the current directory
(assumes there is no "{app_name}" file in the working directory).

~~~
{app_name} start
~~~

Run web server using a specified directory

~~~
   {app_name} start /www/htdocs
~~~

Running web server using a "/etc/{app_name}" file for configuration.

~~~
   {app_name} start /etc/{app_name}
~~~

Running the web server using the basic setup of "/etc/{app_name}"
and overriding the default htdocs root and URL listened on

~~~
   {app_name} start /etc/{app_name} ./htdocs http://localhost:9011
~~~

Configure your web server with these steps

~~~
   {app_name} init webserver.toml
   {app_name} htdocs webserver.toml /var/www/htdocs
   {app_name} url webserver.toml https://www.example.edu:443
   {app_name} cert_pem webserver.toml /etc/certs/cert.pem
   {app_name} key_pem webserver.toml /etc/certs/key.pem
   {app_name} access webserver.toml /etc/wsfn/access.toml
~~~

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
)

// initWebService creates an initialization file.
func initWebService(args []string) error {
	var (
		err error
	)
	fName := "webservice.toml"
	switch {
	case len(args) > 1:
		return fmt.Errorf("Init expects a single filename ending in .toml or .json")
	case len(args) == 1:
		fName = args[0]
	}
	if _, err = os.Stat(fName); os.IsNotExist(err) == false {
		return fmt.Errorf("%q already exists", fName)
	}
	src := wsfn.DefaultInit()
	if strings.HasSuffix(fName, ".json") {
		o := new(wsfn.WebService)
		if _, err = toml.Decode(string(src), &o); err != nil {
			return err
		}
		src, err = json.MarshalIndent(o, "", "    ")
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(fName, src, 0660)
}

// setDocRootWebService sets the document root in an initialization file.
func setDocRootWebService(args []string) error {
	fName, docRoot := "", ""
	switch {
	case len(args) == 2:
		fName, docRoot = args[0], args[1]
	default:
		return fmt.Errorf("expecting web service filename and a single document root")
	}
	ws, err := wsfn.LoadWebService(fName)
	if err != nil {
		return err
	}
	ws.DocRoot = docRoot
	return ws.DumpWebService(fName)
}

// setAccessFile sets the access file
func setAccessFile(args []string) error {
	fName, accessFile := "", ""
	switch {
	case len(args) == 2:
		fName, accessFile = args[0], args[1]
	default:
		return fmt.Errorf("expecting web service filename and a single access filename")
	}
	if _, err := os.Stat(accessFile); os.IsNotExist(err) {
		return fmt.Errorf("%q, %s", accessFile, err)
	}
	ws, err := wsfn.LoadWebService(fName)
	if err != nil {
		return err
	}
	ws.AccessFile = accessFile
	return ws.DumpWebService(fName)
}

// setURL sets the scheme, hostname and port to listen on.
// If the scheme is https it sets the https configuration, if http
// sets the http configuration
func setURL(args []string) error {
	var (
		service *wsfn.Service
		fName   string
		uri     string
	)
	switch {
	case len(args) == 2:
		fName, uri = args[0], args[1]
	default:
		return fmt.Errorf("expecting web service filename and a single document root")
	}
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	ws, err := wsfn.LoadWebService(fName)
	if err != nil {
		return err
	}
	switch u.Scheme {
	case "https":
		if ws.Https == nil {
			ws.Https = new(wsfn.Service)
		}
		service = ws.Https
	case "http":
		if ws.Http == nil {
			ws.Http = new(wsfn.Service)
		}
		service = ws.Http
	default:
		return fmt.Errorf("%s is an unsupported scheme", u.Scheme)
	}
	service.Scheme = u.Scheme
	if strings.Contains(u.Host, ":") {
		p := strings.SplitN(u.Host, ":", 2)
		service.Host, service.Port = p[0], p[1]
	} else {
		service.Host = u.Host
	}
	switch u.Scheme {
	case "https":
		if service.Port == "" {
			service.Port = "443"
		}
		ws.Https = service
	case "http":
		if service.Port == "" {
			service.Port = "80"
		}
		ws.Http = service
	}
	return ws.DumpWebService(fName)
}

// Sets the cert.pem file to used for TLS
func setCertPEM(args []string) error {
	fName, certPEM := "", ""
	switch {
	case len(args) == 2:
		fName, certPEM = args[0], args[1]
	default:
		return fmt.Errorf("expecting web service filename and a path to cert.pem")
	}
	if _, err := os.Stat(certPEM); os.IsNotExist(err) {
		return fmt.Errorf("%q does not exist", certPEM)
	}
	ws, err := wsfn.LoadWebService(fName)
	if err != nil {
		return err
	}
	if ws.Https == nil {
		ws.Https = new(wsfn.Service)
	}
	ws.Https.CertPEM = certPEM
	return ws.DumpWebService(fName)
}

// Sets the key.pem file to used for TLS
func setKeyPEM(args []string) error {
	fName, keyPEM := "", ""
	switch {
	case len(args) == 2:
		fName, keyPEM = args[0], args[1]
	default:
		return fmt.Errorf("expecting web service filename and a path to key.pem")
	}
	if _, err := os.Stat(keyPEM); os.IsNotExist(err) {
		return fmt.Errorf("%q does not exist", keyPEM)
	}
	ws, err := wsfn.LoadWebService(fName)
	if err != nil {
		return err
	}
	if ws.Https == nil {
		ws.Https = new(wsfn.Service)
	}
	ws.Https.KeyPEM = keyPEM
	return ws.DumpWebService(fName)
}

func startService(args []string) error {
	var (
		cfg string
		ws  *wsfn.WebService
		err error
	)
	// check for local config
	if _, err := os.Stat("webserver.toml"); os.IsNotExist(err) == false {
		cfg = "webserver.toml"
	} else if _, err := os.Stat("webserver.json"); os.IsNotExist(err) == false {
		cfg = "webserver.json"
	}
	// Load a default configuration
	if cfg != "" {
		ws, err = wsfn.LoadWebService(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%q, %s\n", cfg, err)
			os.Exit(1)
		}
	} else {
		ws = wsfn.DefaultWebService()
	}
	// Adhoc overrides
	for _, arg := range args {
		switch {
		case strings.HasSuffix(arg, ".toml") || strings.HasSuffix(arg, ".json"):
			ws, err = wsfn.LoadWebService(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%q failed, %s\n", arg, err)
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
				ws.Http.Scheme = u.Scheme
				ws.Http.Host = host
				ws.Http.Port = port
			case "https":
				ws.Https.Scheme = u.Scheme
				ws.Https.Host = host
				ws.Https.Port = port
			default:
				fmt.Fprintf(os.Stderr, "Unsupported scheme %q", u.String())
				os.Exit(1)
			}
		default:
			ws.DocRoot = arg
		}
	}
	// Now we should be ready to run the web server
	if err = ws.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	appName := path.Base(os.Args[0])
	// NOTE: The following are set when version.go is generated
	version := wsfn.Version
	releaseDate := wsfn.ReleaseDate
	releaseHash := wsfn.ReleaseHash
	fmtHelp := wsfn.FmtHelp


	// Standard Options
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "version", false, "display version")
	flag.BoolVar(&quiet, "quiet", false, "suppress error messages")
	flag.StringVar(&outputFName, "o", "", "write output to filename")

	flag.Parse()
	args := flag.Args()

	// Setup IO
	var err error

	out := os.Stdout
	eout := os.Stderr

	
	// Process flags and update the environment as needed.
	if showHelp {
		fmt.Fprintf(out, "%s\n", fmtHelp(helpText, appName, version, releaseDate, releaseHash))
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintln(out, wsfn.LicenseText)
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintf(out, "%s %s %s\n", appName, version, releaseHash)
		os.Exit(0)
	}

	if outputFName != "" && outputFName != "-" {
		out, err = os.Create(outputFName)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
		defer out.Close()
	}

	if len(args) == 0 {
		fmt.Fprintf(eout, "Missing an action (e.g. start, init)\n")
		os.Exit(1)
	}
	verb, args := args[0], args[1:]
	switch verb {
	case "init":
		if err := initWebService(args); err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	case "htdocs":
		if err := setDocRootWebService(args); err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	case "url":
		if err := setURL(args); err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	case "cert_pem":
		if err := setCertPEM(args); err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	case "key_pem":
		if err := setKeyPEM(args); err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	case "access":
		if err := setAccessFile(args); err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
	case "start":
		if err := startService(args); err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(eout, "Unknown action %q\n", verb)
		os.Exit(1)
	}
}
