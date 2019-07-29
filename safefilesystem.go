//
// Package wsfn provides a common library of functions and structures for
// working with web services in Caltech Library projects and software.
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
package wsfn

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

//
// This is loosely based on Go's example of a web server that
// avoids serving dot files.
// See https://golang.org/pkg/net/http/#example_FileServer_dotFileHiding
//

// hasDotPrefix checks a path for containing either ., .. prefixes
// in a path.
func hasDotPrefix(s string) bool {
	parts := strings.Split(s, "/")
	for _, p := range parts {
		if strings.HasPrefix(p, ".") {
			return true
		}
	}
	return false
}

// SafeFile are ones that do NOT have a "." as a prefix
// on the path.
type SafeFile struct {
	http.File
}

// SafeFileSystem is used to hide dot file paths from
// our web services.
type SafeFileSystem struct {
	http.FileSystem
}

// Readdir wraps SafeFile method checks first if we
// have a dot path problem before use http.File.Readdir.
func (sf SafeFile) Readdir(n int) ([]os.FileInfo, error) {
	// Get a raw list of files.
	ls, err := sf.File.Readdir(n)
	if err != nil {
		return nil, err
	}
	infoList := []os.FileInfo{}
	for _, info := range ls {
		if strings.HasPrefix(info.Name(), ".") == false {
			infoList = append(infoList, info)
		}
	}
	return infoList, nil
}

// Open is a wrapper around the Open method of the embedded
// SafeFileSystem. It serves a 403 permision error when name has
// a file or directory who's path parts is a dot file prefix.
func (fs SafeFileSystem) Open(p string) (http.File, error) {
	if hasDotPrefix(p) {
		// If dot file setup to return a 403 response by
		// passing an OS level file permission error
		return nil, os.ErrPermission
	}
	// If we got this fare we can open the file safely.
	fp, err := fs.FileSystem.Open(p)
	if err != nil {
		return nil, err
	}
	return SafeFile{fp}, err
}

///
// SafeFileSystem returns a new safe file system using
// the *WebService.DocRoot as the directory source.
//
// Example usage:
//
// ws := wsfn.LoadTOML("web-service.toml")
// fs, err := ws.SafeFileSystem()
// if err != nil {
//     log.Fatalf("%s\n", err)
// }
// http.Handle("/", http.FileServer(ws.SafeFileSystem()))
// log.Fatal(http.ListenAndService(ws.Http.Hostname(), nil))
//
func (w *WebService) SafeFileSystem() (SafeFileSystem, error) {
	if w.DocRoot == "" {
		w.DocRoot = "."
	}
	if info, err := os.Stat(w.DocRoot); err != nil {
		return SafeFileSystem{}, err
	} else if info.IsDir() == false {
		return SafeFileSystem{}, fmt.Errorf("%q is not a directory", w.DocRoot)
	}
	return SafeFileSystem{http.Dir(w.DocRoot)}, nil
}

//
// MakeSafeFileSystem without a *WebService takes a doc root
// and returns a SafeFileSystem struct.
//
// Example usage:
//
// fs, err := MakeSafeFileSystem("/var/www/htdocs")
// if err != nil {
//     log.Fatalf("%s\n", err)
// }
// http.Handle("/", http.FileServer(fs))
// log.Fatal(http.ListenAndService(":8000", nil))
//
func MakeSafeFileSystem(docRoot string) (SafeFileSystem, error) {
	if docRoot == "" {
		return SafeFileSystem{}, fmt.Errorf("document root not set")
	}
	if info, err := os.Stat(docRoot); err != nil {
		return SafeFileSystem{}, err
	} else if info.IsDir() == false {
		return SafeFileSystem{}, fmt.Errorf("%q is not a directory", docRoot)
	}
	return SafeFileSystem{http.Dir(docRoot)}, nil
}
