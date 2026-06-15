//
// webaccess.go - Generates/Manages a "access.toml" file.
// for use with wsfn basic auth services. "access.toml" is
// analogous to Apache's htpasswd file but uses a different
// format including identifying the authentication setup of the
// web service instanciated with wsfn.
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
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"

	// X packages
	"golang.org/x/crypto/ssh/terminal"

	// Caltech Library packages
	"github.com/caltechlibrary/wsfn"
)

// Flag options
var (
	// Standard options
	showHelp    bool
	showVersion bool
	showLicense bool
	outputFName string
	quiet       bool
)

func initAccess(fName string) error {
	if fName == "" {
		fName = "access.toml"
	}
	if _, err := os.Stat(fName); os.IsNotExist(err) == false {
		return fmt.Errorf("%q already exists", fName)
	}
	a := new(wsfn.Access)
	a.AuthType = "basic"
	a.Encryption = "argon2id"
	return a.DumpAccess(fName)
}

func updateAccess(fName, username, password string) error {
	a, err := wsfn.LoadAccess(fName)
	if err != nil {
		return err
	}
	if a.UpdateAccess(username, password) == false {
		return fmt.Errorf("Failed to update %s", username)
	}
	return a.DumpAccess(fName)
}

func removeAccess(fName, username string) error {
	a, err := wsfn.LoadAccess(fName)
	if err != nil {
		return err
	}
	if a.RemoveAccess(username) == false {
		return fmt.Errorf("Failed to find %s", username)
	}
	return a.DumpAccess(fName)
}

func listAccess(fName string) error {
	var (
		a   *wsfn.Access
		err error
	)
	a, err = wsfn.LoadAccess(fName)
	if err != nil {
		return err
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
	a, err = wsfn.LoadAccess(fName)
	if err != nil {
		return err
	}
	if a.Login(username, password) == false {
		return fmt.Errorf("Failed to authenticate %s", username)
	}
	return nil
}

func listRoutes(a *wsfn.Access) error {
	for _, route := range a.Routes {
		fmt.Fprintf(os.Stdout, "%s\n", route)
	}
	return nil
}

func updateRoutes(fName string, a *wsfn.Access, args []string) error {
	for _, arg := range args {
		if strings.HasPrefix(arg, "/") == false {
			arg = "/" + arg
		}
		if strings.HasSuffix(arg, "/") == false {
			arg += "/"
		}
		for _, route := range a.Routes {
			if strings.HasPrefix(arg, route) || strings.HasPrefix(route, arg) {
				return fmt.Errorf("%q collides with %q", arg, route)
			}
		}
		a.Routes = append(a.Routes, arg)
		sort.Strings(a.Routes)
	}
	return a.DumpAccess(fName)
}

func removeRoutes(fName string, a *wsfn.Access, args []string) error {
	for _, arg := range args {
		routeFound := false
		if strings.HasPrefix(arg, "/") == false {
			arg = "/" + arg
		}
		for i, route := range a.Routes {
			if strings.Compare(arg, route) == 0 {
				a.Routes = append(a.Routes[:i], a.Routes[i+1:]...)
				routeFound = true
			}
		}
		if routeFound == false {
			return fmt.Errorf("Could not find route %q", arg)
		}
	}
	sort.Strings(a.Routes)
	return a.DumpAccess(fName)
}

func manageRoutes(args []string) error {
	var (
		verb  string
		fName string
	)
	switch len(args) {
	case 0:
		return fmt.Errorf("update, list, or remove?")
	case 1:
		return fmt.Errorf("missing access filename")
	case 2:
		verb, fName = args[0], args[1]
		args = []string{}
	default:
		verb, fName, args = args[0], args[1], args[2:]
	}
	a, err := wsfn.LoadAccess(fName)
	if err != nil {
		return err
	}
	switch verb {
	case "list":
		return listRoutes(a)
	case "update":
		return updateRoutes(fName, a, args)
	case "remove":
		return removeRoutes(fName, a, args)
	default:
		return fmt.Errorf("Unknown route action, %q", verb)
	}
}

// showHelpTopic displays help for a specific topic or the main help
func showHelpTopic(w io.Writer, topic, appName, version, releaseDate, releaseHash string) {
	var helpText string
	
	switch topic {
	case "", "help":
		helpText = wsfn.WebaccessMainHelp
	case "access-control", "access":
		helpText = wsfn.WebaccessAccessControlTopicHelp
	case "users":
		helpText = wsfn.WebaccessUsersTopicHelp
	case "encryption":
		helpText = wsfn.WebaccessEncryptionTopicHelp
	case "config", "configuration":
		helpText = wsfn.WebaccessConfigTopicHelp
	default:
		// Unknown topic - show main help with error message
		fmt.Fprintf(w, "Unknown help topic: %s\n\n", topic)
		fmt.Fprintf(w, "%s\n", wsfn.FmtHelp(wsfn.WebaccessMainHelp, appName, version, releaseDate, releaseHash))
		return
	}
	
	fmt.Fprintf(w, "%s\n", wsfn.FmtHelp(helpText, appName, version, releaseDate, releaseHash))
}

func main() {
	appName := path.Base(os.Args[0])
	// NOTE: The following are set when version.go is generated
	version := wsfn.Version
	releaseDate := wsfn.ReleaseDate
	releaseHash := wsfn.ReleaseHash
	fmtHelp := wsfn.FmtHelp

	// Handle -help TOPIC before flag parsing
	// We need to check os.Args directly to support -help topic syntax
	helpTopic := ""
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			if os.Args[i] == "-help" && i+1 < len(os.Args) {
				// Check if next arg is not a flag (doesn't start with -)
				if len(os.Args[i+1]) > 0 && os.Args[i+1][0] != '-' {
					helpTopic = os.Args[i+1]
					// Display help for topic and exit immediately
					showHelpTopic(os.Stdout, helpTopic, appName, version, releaseDate, releaseHash)
					os.Exit(0)
				} else {
					// Just -help, remove it so flag doesn't see it
					newArgs := make([]string, 0, len(os.Args)-1)
					newArgs = append(newArgs, os.Args[:i]...)
					newArgs = append(newArgs, os.Args[i+1:]...)
					os.Args = newArgs
					break
				}
			}
		}
	}

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

	//in := os.Stdin
	out := os.Stdout
	eout := os.Stderr

	// Process flags and update the environment as needed.
	if showHelp {
		showHelpTopic(out, helpTopic, appName, version, releaseDate, releaseHash)
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprint(out, wsfn.LicenseText)
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



	verb, fName, userid := "", "", ""
	switch len(args) {
	case 3:
		verb, fName, userid = args[0], args[1], args[2]
	case 2:
		verb, fName, userid = args[0], args[1], ""
	case 1:
		verb, fName, userid = args[0], "", ""
		if strings.Compare(verb, "routes") == 0 {
			fmt.Fprintf(eout, "Missing action and parameters\ntry %s -h\n", appName)
			os.Exit(1)
		}
	case 0:
		fmt.Fprintf(eout, "%s\n", fmtHelp(wsfn.WebaccessMainHelp, appName, version, releaseDate, releaseHash))
		os.Exit(1)
	default:
		verb, fName, userid = args[0], "", ""
		if strings.Compare(verb, "routes") != 0 {
			fmt.Fprintf(eout, "To many parameters, try %s -help\n", appName)
			os.Exit(1)
		}
	}

	switch verb {
		case "help":
			topic := ""
			if len(args) > 1 {
				topic = args[1]
			}
			showHelpTopic(out, topic, appName, version, releaseDate, releaseHash)
			os.Exit(0)
	case "init":
		err = initAccess(fName)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
	case "update":
		fmt.Fprintf(os.Stdout, "Enter a password:\n")
		password, err := terminal.ReadPassword(0)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
		if err = updateAccess(fName, userid, string(password)); err != nil {
			fmt.Fprintf(eout, "update failed, %s\n", err)
			os.Exit(1)
		}
	case "remove":
		if err = removeAccess(fName, userid); err != nil {
			fmt.Fprintf(eout, "remove failed, %s\n", err)
			os.Exit(1)
		}
	case "list":
		if err = listAccess(fName); err != nil {
			fmt.Fprintf(eout, "list failed, %s\n", err)
			os.Exit(1)
		}
	case "test":
		fmt.Fprintf(os.Stdout, "Enter a password:\n")
		password, err := terminal.ReadPassword(0)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			os.Exit(1)
		}
		if err = testAccess(fName, userid, string(password)); err != nil {
			fmt.Fprintf(eout, "test failed, %s\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "OK\n")
	case "routes":
		if err = manageRoutes(args[1:]); err != nil {
			fmt.Fprintf(eout, "%s %s, failed\n%s\n", appName,
				strings.Join(args, " "), err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(eout, "Unknown action %q, try %s -help\n", verb, appName)
		os.Exit(1)
	}
}
