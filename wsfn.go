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
	"path"
	"strings"
)

// IsDotPath checks to see if a path is requested with a dot file (e.g. docs/.git/* or docs/.htaccess)
func IsDotPath(p string) bool {
	for _, part := range strings.Split(path.Clean(p), "/") {
		if strings.HasPrefix(part, "..") == false && strings.HasPrefix(part, ".") == true && len(part) > 1 {
			return true
		}
	}
	return false
}

// StaticRouter scans the request object to either add a .html extension
// or prevent serving a dot file path
func StaticRouter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Methods", "GET")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		// Stop here if its Preflighted OPTIONS request
		if r.Method == "OPTIONS" {
			return
		}

		// If given a dot file path, send forbidden
		if IsDotPath(r.URL.Path) == true {
			http.Error(w, "Forbidden", 403)
			ResponseLogger(r, 403, fmt.Errorf("Forbidden, requested a dot path"))
			return
		}
		// Check if we have a gzipped JSON file
		if strings.HasSuffix(r.URL.Path, ".json.gz") || strings.HasSuffix(r.URL.Path, ".js.gz") {
			w.Header().Set("Content-Encoding", "gzip")
		}
		// Check to see if we have a *.wasm file, then make sure
		// we have the correct headers.
		if ext := path.Ext(r.URL.Path); ext == ".wasm" {
			w.Header().Set("Content-Type", "application/wasm")
		}
		// Check to see if we have a JS module file, then make sure
		// we have the correct headers
		if ext := path.Ext(r.URL.Path); (ext == ".mjs") || (ext == ".js") {
			w.Header().Set("Content-Type", "text/javascript")
		}

		// If we make it this far, fall back to the default handler
		next.ServeHTTP(w, r)
	})
}
