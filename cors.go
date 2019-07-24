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
	"net/http"
	"strings"
)

// CORSPolicy defines the policy elements for our CORS settings.
type CORSPolicy struct {
	// Origin usually would be set the hostname of the service.
	Origin string `json:"origin,omitempty" toml:"origin,omitempty"`
	// Options to include in the policy (e.g. GET, POST)
	Options []string `json:"options,omitempty" toml:"options,omitempty"`
	// Headers to include in the policy
	Headers []string `json:"headers,omitempty" toml:"headers,omitempty"`
	// ExposedHeaders to include in the policy
	ExposedHeaders []string `json:"exposed_headers,omitempty" toml:"exposed_headers,omitempty"`
	// AllowCredentials header handling in the policy either true or not set
	AllowCredentials bool `json:"allow_credentials,omitempty" toml:"allow_credentials,omitempty"`
}

// Handle accepts an http.Handler and returns a http.Handler. It
// Wraps the response with the CORS headers based on configuration
// in CORSPolicy struct.
func (cors *CORSPolicy) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cors.Origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", cors.Origin)
		}
		if len(cors.Options) > 0 {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(cors.Options, ","))
		}
		if len(cors.Headers) > 0 {
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(cors.Headers, ","))
		}
		if len(cors.ExposedHeaders) > 0 {
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(cors.ExposedHeaders, ","))
		}
		if cors.AllowCredentials == true {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		// Bailout if we ahve an OPTIONS preflight request
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
