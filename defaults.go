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

// DefaultService is http, port 8000 on localhost.
func DefaultService() *Service {
	h := new(Service)
	h.Scheme = "http"
	h.Host = "localhost"
	h.Port = "8000"
	return h
}

// DefaultWebService setups to listen for http://localhost:8000
// with the htdocs of the current working directory.
func DefaultWebService() *WebService {
	w := new(WebService)
	w.DocRoot = "."
	w.Http = DefaultService()
	return w
}

// DefaultInit generates a default TOML initialization file.
func DefaultInit() []byte {
	return []byte(`
#
# A TOML file example for configuring **webserver**.
# Comments start with "#"
#

# 
# Setup your document root for the website.
# This must be before the other entries.
#
# It is relative to the current working directory
# unless a path is fully specified. A period or 
# empty string will set it to the current working 
# directory.
htdocs = "htdocs"

#
# If using access restrictions (e.g. basic auth)
# set the file for managing access.
# Uncomment to use.
#
#access_file = "access.toml"

#
# Use redirects in a separate file (e.g. JSON, TOML, CSV).
# Uncomment to use.
#
#redirects_file = "redirects.csv"

#
# Managing content types in a separate file (e.g. JSON, TOML, CSV)
# Uncomment to use.
#
#content_types_file = "content-types.csv"

# Setting up standard http support
[http]
host = "localhost"
port = "8000"

# Setting up HTTPS scheme support, uncomment for https support
#[https]
#cert_pem = "etc/certs/cert_pem"
#key_pem = "etc/certs/key_pem"
#host = "localhost"
#port = "8443"

#
# CORS policy configuration example adpated from 
# Mozilla website.
# See https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
#
# Uncomment to use.
#[cors]
#Access_Control_Origin = "http://foo.example:8000"
#Access_Control_Allow_Credentials = true
#Access_Control_Methods = [ "POST", "GET" ]
#Access_Control_Allow_Headers = [ "X-PINGPONG", "Content-Type" ]
#Access_Control_Max_Age = 86400

#
# Managing file extensions to mime types in the
# file.
#
# Uncomment to use.
#[content_types]
#".json" = "application/json"
#".toml" = "text/plain+x-toml"

#
# Managing redirects in this file.
#
# Uncomment to use.
#[redirects]
#"http://localhost:8000/" = "https://localhost:8443/"
#"/bad-path/" = "/good-path/"

#
# Managin reverse-proxy in this file.
#
# Uncomment to use.
#[reverse_proxy]
#"/api/" = "http://localhost:9000/"

`)
}
