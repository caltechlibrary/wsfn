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
	"io/ioutil"
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

%s is a command line utility for developing and testing 
static websites.  It uses Go's standard http libraries 
and can supports both http 1 and 2 out of the box.  It 
provides a minimal set of extra features useful for 
developing and testing web services that leverage static 
content. 

CONFIGURATION

%s is configured through a configuration file. You can
create an initialization file using the "init" action.
By default the created initialation file is "%s".

ACTIONS

The following actions are available

+ init     creates a %q file.
+ htdocs   sets the document root
+ cert_pem set the path to find cert.pem file for TLS
+ key_pem  set the path to find the key.pem file for TLS
+ auth     set auth type if used, e.g. Basic
+ access   sets an external access file

To configure an external access file use the "webaccess" tool.

`

	examples = `
Run web server using the content in the current directory
(assumes there is no "%s" file in the working directory).

   %s

Run web server using a specified directory

   %s /www/htdocs

Running web server using a "/etc/%s" file for configuration.

   %s /etc/%s

Running the web server using the basic setup of "/etc/%s"
and overriding the default htdocs root and URL listened on

   %s start /etc/%s ./htdocs http://localhost:9011

Configure your web server with

   %s init webserver.toml
   %s htdocs webserver.toml /var/www/htdocs
   %s url webserver.toml https://www.example.edu:443
   %s cert_pem webserver.toml /etc/certs/cert.pem
   %s key_pem webserver.toml /etc/certs/key.pem
   %s access webserver.toml /etc/wsfn/access.toml

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
	fName := "webservice.toml"
	switch {
	case len(args) > 1:
		return fmt.Errorf("Init expects a single filename ending in .toml or .json")
	case len(args) == 1:
		fName = args[0]
	}
	src := wsfn.DefaultInit()
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

func main() {
	app := cli.NewCli(wsfn.Version)
	appName := app.AppName()
	configName := strings.TrimSuffix(strings.ToLower(appName), ".exe") + ".toml"

	// Document non-option parameters
	app.SetParams(`[VERB PARAMETERS || CONFIG_NAME]`, `[DOCROOT]`, `[URL_TO_LISTEN_ON]`)

	// Add Help Docs
	app.AddHelp("license", []byte(fmt.Sprintf(wsfn.LicenseText, appName, wsfn.Version)))
	app.AddHelp("description", []byte(fmt.Sprintf(description, appName, appName, configName, configName)))
	app.AddHelp("examples", []byte(fmt.Sprintf(examples, appName, appName, appName, appName, appName, configName, configName, appName, configName, appName, appName, appName, appName, appName, appName)))

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "example", false, "display example(s)")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "generate man page")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")

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

	// setup from command line
	webservice := wsfn.DefaultWebService()
	if _, err := os.Stat("webserver.toml"); os.IsNotExist(err) == false {
		webservice, err = wsfn.LoadWebService("webserver.toml")
		if err != nil {
			fmt.Fprintf(os.Stderr, "toml parse %s, %s\n", "webserver.toml", err)
			os.Exit(1)
		}
	} else if _, err := os.Stat("webserver.json"); os.IsNotExist(err) == false {
		webservice, err = wsfn.LoadWebService("webserver.json")
		if err != nil {
			fmt.Fprintf(os.Stderr, "json parse %s, %s\n", "webserver.json", err)
			os.Exit(1)
		}
	}
	if len(args) == 0 {
		// Now we should be ready to run the web server
		if err = webservice.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}
	switch args[0] {
	case "init":
		if err := initWebService(args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	case "htdocs":
		if err := setDocRootWebService(args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	case "url":
		if err := setURL(args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	case "cert_pem":
		if err := setCertPEM(args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	case "key_pem":
		if err := setKeyPEM(args[1:]); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}
	// Adhoc command overrides
	for _, arg := range args {
		switch {
		case strings.HasSuffix(arg, ".toml") || strings.HasSuffix(arg, ".json"):
			webservice, err = wsfn.LoadWebService(arg)
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
