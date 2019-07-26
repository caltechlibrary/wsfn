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
# Lines starting with "#" are comments.
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


# Setting up HTTPS scheme support, uncomment for https support
#[https]
#cert_pem = "etc/certs/cert_pem"
#key_pem = "etc/certs/key_pem"
#host = "localhost"
#port = "8443"

# Setting up standard http support
[http]
host = "localhost"
port = "8000"

#
# Configure HTTP Basic AUTH
# Example if uncommented would protect the /api/ path.
#
# paths = is a list of protected web server paths to protect
# passwords = is a path to a file in Apache htpasswd format to
# use for Baisc AUTH.
#
#[basic_auth]
#paths = [ "/api/" ]
#passwords = "etc/ws-api-passwords"

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
# Mapping file extensions to mime types
#
# Uncomment to use.
#[content_types]
#".json" = "application/json"
#".toml" = "text/plain+x-toml"

#
# Redirects are specified in CSV file format.
# first column is the target, second the destination
#
# Uncomment to use.
#
#redirects_csv = "redirects.csv"

#
# Redirects are specified in this file.
#
# Uncomment and edit to use.
#[redirects]
#"http://localhost:8000/" = "https://localhost:8443/"
#"/bad-path/" = "/good-path/"

#
# reverse-proxy examples
#
# Uncomment and edit to use.
#[reverse_proxy]
#"/api/" = "http://localhost:9000/"

# To added access configuration using webaccess tool.
#[access]
# ...
`)
}
