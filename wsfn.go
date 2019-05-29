//
// Package wsfn provides a common library of functions and structures for
// working with web services in Caltech Library projects and software.
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2017, Caltech
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
	"log"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
)

const Version = `v0.0.3`

// IsDotPath checks to see if a path is requested with a dot file (e.g. docs/.git/* or docs/.htaccess)
func IsDotPath(p string) bool {
	for _, part := range strings.Split(path.Clean(p), "/") {
		if strings.HasPrefix(part, "..") == false && strings.HasPrefix(part, ".") == true && len(part) > 1 {
			return true
		}
	}
	return false
}

type CORSPolicy struct {
	Origin           string
	Options          string
	Headers          string
	ExposedHeaders   string
	AllowCredentials string
}

var (
	// An internal ordered list of keys in redirectRoutes map
	redirectPrefixes = []string{}
	// Our map of redirect prefix to target replacement routes
	redirectRoutes = map[string]string{}
)

// Handle accepts an http.Handler and returns a http.Handler. It
// Wraps the response with the CORS headers based on configuration
// in CORSPolicy struct.
func (cors *CORSPolicy) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cors.Origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", cors.Origin)
		}
		if cors.Options != "" {
			w.Header().Set("Access-Control-Allow-Methods", cors.Options)
		}
		if cors.Headers != "" {
			w.Header().Set("Access-Control-Allow-Headers", cors.Headers)
		}
		if cors.ExposedHeaders != "" {
			w.Header().Set("Access-Control-Expose-Headers", cors.ExposedHeaders)
		}
		if cors.AllowCredentials != "" {
			w.Header().Set("Access-Control-Allow-Credentials", cors.AllowCredentials)
		}
		// Bailout if we ahve an OPTIONS preflight request
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequestLogger logs the request based on the request object passed into
// it.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if len(q) > 0 {
			log.Printf("Request: %s Path: %s RemoteAddr: %s UserAgent: %s Query: %+v\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), q)
		} else {
			log.Printf("Request: %s Path: %s RemoteAddr: %s UserAgent: %s\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
		}
		next.ServeHTTP(w, r)
	})
}

// ResponseLogger logs the response based on a request, status and error
// message
func ResponseLogger(r *http.Request, status int, err error) {
	q := r.URL.Query()
	if len(q) > 0 {
		log.Printf("Response: %s Path: %s RemoteAddr: %s UserAgent: %s Query: %+v Status: %d, %s %q\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), q, status, http.StatusText(status), err)
	} else {
		log.Printf("Response: %s Path: %s RemoteAddr: %s UserAgent: %s Status: %d, %s %q\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), status, http.StatusText(status), err)
	}
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

		// If we make it this far, fall back to the default handler
		next.ServeHTTP(w, r)
	})
}

// HasRedirectRoutes returns true if redirects have been defined,
// false if not.
func HasRedirectRoutes() bool {
	if len(redirectPrefixes) > 0 || len(redirectRoutes) > 0 {
		return true
	}
	return false
}

// AddRedirectRoute takes a target and a destination prefix
// and populates the internal datastructures to handle
// the redirecting target prefix to the destination prefix.
func AddRedirectRoute(target, destination string) error {
	// Make sure prefix has not been defined
	for _, p := range redirectPrefixes {
		if strings.HasPrefix(p, target) || strings.HasPrefix(target, p) {
			return fmt.Errorf("targets %q and %q collide", target, p)
		}
	}
	redirectRoutes[target] = destination
	redirectPrefixes = append(redirectPrefixes, target)
	sort.Strings(redirectPrefixes)
	return nil
}

// isRedirectTarget
func isRedirectTarget(srcPath, targetPath string) bool {
	// Do a course level check first ...
	ok, _ := path.Match(targetPath, srcPath)
	return ok
}

// RedirectRouter handles redirect requests before passing on to the
// handler.
func RedirectRouter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do we have a redirect prefix in r.URL.Path
		for _, target := range redirectPrefixes {
			if strings.HasPrefix(r.URL.Path, target) {
				// Update our path to use new prefix
				if destination, ok := redirectRoutes[target]; ok == true {
					// Clone our existing Request URL ...
					u, _ := url.Parse(r.URL.String())
					// Calculate a new path
					p := strings.TrimPrefix(u.Path, target)
					// Update our new path.
					u.Path = path.Join(destination, p)
					log.Printf("Redirecting %q to %q", r.URL.String(), u.String())
					// Send our redirect on its way!
					http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
					return
				}
			}
		}
		// If we make it this far, fall back to the default handler
		next.ServeHTTP(w, r)
	})
}
