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
	"log"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
)

// RedirectService holds our redirect targets in an ordered list
// and a map to our applied routes.
type RedirectService struct {
	// Our map of redirect prefix to target replacement routes
	routes map[string]string
}

// HasRedirectRoutes returns true if redirects have been defined,
// false if not.
func (r *RedirectService) HasRedirectRoutes() bool {
	if len(r.routes) > 0 {
		return true
	}
	return false
}

// HasRoute returns true if the target route is defined
func (r *RedirectService) HasRoute(key string) bool {
	_, ok := r.routes[key]
	return ok
}

// Route takes a target and returns a destination and bool.
func (r *RedirectService) Route(key string) (string, bool) {
	destination, ok := r.routes[key]
	return destination, ok
}

// MakeRedirectService takes a m[string]string of redirects
// and loads it into our service's private routes attribute.
// It returns a new *RedirectService and error
func MakeRedirectService(m map[string]string) (*RedirectService, error) {
	r := new(RedirectService)
	if r.routes == nil {
		r.routes = make(map[string]string)
	}
	for k, v := range m {
		if err := r.AddRedirectRoute(k, v); err != nil {
			return r, err
		}
	}
	return r, nil
}

// AddRedirectRoute takes a target and a destination prefix
// and populates the internal datastructures to handle
// the redirecting target prefix to the destination prefix.
func (r *RedirectService) AddRedirectRoute(target, destination string) error {
	if r.routes == nil {
		r.routes = make(map[string]string)
	}
	prefixes := []string{}
	for key, _ := range r.routes {
		prefixes = append(prefixes, key)
	}
	sort.Strings(prefixes)
	// Make sure prefix has not been defined and don't collide
	for _, p := range prefixes {
		if strings.HasPrefix(p, target) || strings.HasPrefix(target, p) {
			return fmt.Errorf("targets %q and %q collide", target, p)
		}
	}
	r.routes[target] = destination
	return nil
}

// RedirectRouter handles redirect requests before passing on to the
// handler.
func (r *RedirectService) RedirectRouter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Do we have a redirect prefix in r.URL.Path
		for target, destination := range r.routes {
			if strings.HasPrefix(req.URL.Path, target) {
				// Clone our existing Request URL ...
				u, _ := url.Parse(req.URL.String())
				// Calculate a new path
				p := strings.TrimPrefix(u.Path, target)
				// Update our new path.
				u.Path = path.Join(destination, p)
				log.Printf("Redirecting %q to %q", req.URL.String(), u.String())
				// Send our redirect on its way!
				http.Redirect(w, req, u.String(), http.StatusMovedPermanently)
				return
			}
		}
		// If we make it this far, fall back to the default handler
		next.ServeHTTP(w, req)
	})
}
