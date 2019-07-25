//
// webaccess.go - Generates/Manages a "access.toml" file.
// for use with wsfn basic auth services. "access.toml" is
// analogous to Apache's htpasswd file but uses a different
// format including identifying the authentication setup of the
// web service instanciated with wsfn.
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
	"os"
	"strings"

	// X packages
	"golang.org/x/crypto/ssh/terminal"

	// Caltech Library packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/wsfn"
)

// Flag options
var (
	description = `

SYNOPSIS

	a nimble web server user access file manager

%s is a command line utility for setting up/managing
user access to web services built on wsfn.

CONFIGURATION

%s generates/manages a TOML file usually named "access.toml". 
The following is an example.

` + "```" + `toml
` + "```" + `
`

	examples = `
Create an empty "access.toml" file.

   %s access.toml init

Add user id "Jane.Doe" to "access.toml".
The access program prompts for a password. It will 
create "access.toml" if it doesn't exist.

   %s access.toml updte Jane.Doe

Remove "Jane.Doe" from access.toml.

   %s access.toml remove Jane.Doe

List users defined in access.toml.

   %s access.toml list

Test a login for Jane.Doe (will prompt for password)

   %s access.toml test Jane.Doe
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

func initAccess(fName string) error {
	a := new(wsfn.Access)
	a.AuthType = "basic"
	a.Encryption = "argon2id"
	switch {
	case strings.HasSuffix(fName, ".toml"):
		return a.DumpAccessTOML(fName)
	case strings.HasSuffix(fName, ".json"):
		return a.DumpAccessJSON(fName)
	default:
		return fmt.Errorf("Unknown format of file.")
	}
}

func updateAccess(fName, username, password string) error {
	// See if fName exists
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		if err := initAccess(fName); err != nil {
			return err
		}
	}
	switch {
	case strings.HasSuffix(fName, ".toml"):
		a, err := wsfn.LoadAccessTOML(fName)
		if err != nil {
			return err
		}
		if a.UpdateAccess(username, password) == false {
			return fmt.Errorf("Failed to update %s", username)
		}
		return a.DumpAccessTOML(fName)
	case strings.HasSuffix(fName, ".json"):
		a, err := wsfn.LoadAccessJSON(fName)
		if err != nil {
			return err
		}
		if a.UpdateAccess(username, password) == false {
			return fmt.Errorf("Failed to update %s", username)
		}
		return a.DumpAccessJSON(fName)
	default:
		return fmt.Errorf("Unknown format of file.")
	}
}

func removeAccess(fName, username string) error {
	switch {
	case strings.HasSuffix(fName, ".toml"):
		a, err := wsfn.LoadAccessTOML(fName)
		if err != nil {
			return err
		}
		if a.RemoveAccess(username) == false {
			return fmt.Errorf("Failed to find %s", username)
		}
		return a.DumpAccessTOML(fName)
	case strings.HasSuffix(fName, ".json"):
		a, err := wsfn.LoadAccessJSON(fName)
		if err != nil {
			return err
		}
		if a.RemoveAccess(username) == false {
			return fmt.Errorf("Failed to find %s", username)
		}
		return a.DumpAccessJSON(fName)
	default:
		return fmt.Errorf("Unknown format of file.")
	}
}

func listAccess(fName string) error {
	var (
		a   *wsfn.Access
		err error
	)
	switch {
	case strings.HasSuffix(fName, ".toml"):
		a, err = wsfn.LoadAccessTOML(fName)
		if err != nil {
			return err
		}
	case strings.HasSuffix(fName, ".json"):
		a, err = wsfn.LoadAccessJSON(fName)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown format of file.")
	}
	for key, _ := range a.Map {
		if key != "" {
			fmt.Fprintf(os.Stdout, "%s\n", key)
		}
	}
	return nil
}

func testAccess(fName, username, password string) error {
	var (
		a   *wsfn.Access
		err error
	)
	// See if fName exists
	if _, err = os.Stat(fName); os.IsNotExist(err) {
		return err
	}
	switch {
	case strings.HasSuffix(fName, ".toml"):
		a, err = wsfn.LoadAccessTOML(fName)
		if err != nil {
			return err
		}
	case strings.HasSuffix(fName, ".json"):
		a, err = wsfn.LoadAccessJSON(fName)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unknown format of file.")
	}
	if a.Login(username, password) == false {
		return fmt.Errorf("Failed to authenticate %s", username)
	}
	return nil
}

func main() {
	app := cli.NewCli(wsfn.Version)
	appName := app.AppName()

	// Document non-option parameters
	app.SetParams(`ACCESS_TOML VERB [PARAMETER]`)

	// Add Help Docs
	app.AddHelp("license", []byte(fmt.Sprintf(wsfn.LicenseText, appName, wsfn.Version)))
	app.AddHelp("description", []byte(fmt.Sprintf(description, appName, appName)))
	app.AddHelp("examples", []byte(fmt.Sprintf(examples, appName, appName, appName, appName, appName)))

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

	fName, verb, userid := "", "", ""
	switch len(args) {
	case 3:
		fName, verb, userid = args[0], args[1], args[2]
	case 2:
		fName, verb, userid = args[0], args[1], ""
	default:
		fmt.Fprintf(os.Stderr, "Incomplete command, try %s -help", appName)
		os.Exit(1)
	}

	switch verb {
	case "init":
		err = initAccess(fName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	case "update":
		fmt.Fprintf(os.Stdout, "Enter a password:\n")
		password, err := terminal.ReadPassword(0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		if err = updateAccess(fName, userid, string(password)); err != nil {
			fmt.Fprintf(os.Stderr, "update failed, %s\n", err)
			os.Exit(1)
		}
	case "remove":
		if err = removeAccess(fName, userid); err != nil {
			fmt.Fprintf(os.Stderr, "remove failed, %s\n", err)
			os.Exit(1)
		}
	case "list":
		if err = listAccess(fName); err != nil {
			fmt.Fprintf(os.Stderr, "list failed, %s\n", err)
			os.Exit(1)
		}
	case "test":
		fmt.Fprintf(os.Stdout, "Enter a password:\n")
		password, err := terminal.ReadPassword(0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		if err = testAccess(fName, userid, string(password)); err != nil {
			fmt.Fprintf(os.Stderr, "test failed, %s\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "OK\n")
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command, try %s -help", appName)
		os.Exit(1)
	}
}
