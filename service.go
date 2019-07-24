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
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

// WebService describes all the configuration and
// capabilities of running a wsfn based web service.
type WebService struct {
	// This is the document root for static file services
	// If an empty string then assume current working directory.
	DocRoot string `json:"htdocs" toml:"htdocs"`
	// Https describes an Https service
	Https *Service `json:"https,omitempty" toml:"https,omitempty"`
	// Http describes an Http service
	Http *Service `json:"http,omitempty" toml:"http,omitempty"`

	// BasicAUTH describes the configuration of BasicAUTH protected paths
	//BasicAUTH *BasicAuth `json:"basic_auth,omitempty" toml:"basic_auth,omiteempty"`

	// CORS describes the CORS policy for the web services
	CORS *CORSPolicy `json:"cors,omitempty" toml:"cors,omitempty"`

	// ContentTypes describes a file extension mapped to a single
	// MimeType.
	ContentTypes map[string]string `json:"content_types,omitempty" toml:"content_types,omitempty"`

	// RedirectsCSV is the filename/path to a CSV file describing
	// redirects.
	RedirectsCSV string `json:"redirects_csv,omitempty" toml:"redirects_csv,omitempty"`

	// Redirects describes a target path to destination path.
	// Normally this is populated by a redirects.csv file.
	Redirects map[string]string `json:"redirects,omitempty" toml:"redirects,omitempty"`

	// ReverseProxy descibes the path web paths that are sent
	// to another proxied URL.
	ReverseProxy map[string]string `json:"reverse_proxy,omitempty" toml:"reverse_proxy,omitempty"`
}

// Service holds the description needed to startup a service
// e.g. https, http.
type Service struct {
	// Scheme holds the protocol to use, defaults to http if not set.
	Scheme string `json:"scheme,omitempty" toml:"scheme,omitempty"`
	// Host is the hostname to use, if empty "localhost" is assumed"
	Host string `json:"host,omitempty" toml:"host,omitempty"`
	// Port is a string holding the port number to listen on
	// An empty strings defaults port to 8000
	Port string `json:"port,omitempty" toml:"port,omitempty"`
	// CertPEM describes the location of cert.pem used for TLS support
	CertPEM string `json:"cert_pem" toml:"cert_pem"`
	// KeyPEM describes the location of the key.pem used for TLS support
	KeyPEM string `json:"key_pem" toml:"key_pem"`
}

// String renders an URL version of *Service.
func (s *Service) String() string {
	r := []string{}
	if s.Scheme != "" {
		r = append(r, s.Scheme, "://")
	}
	r = append(r, s.Hostname())
	return strings.Join(r, "")
}

// Hostname returns a host+port from a *Service
func (s *Service) Hostname() string {
	r := []string{}
	if s.Host != "" {
		r = append(r, s.Host)
	}
	if s.Port != "" {
		r = append(r, ":", s.Port)
	}
	return strings.Join(r, "")
}

// LoadTOML loads a configuration of *WebService from a TOML file.
func LoadTOML(setup string) (*WebService, error) {
	src, err := ioutil.ReadFile(setup)
	if err != nil {
		return nil, err
	}
	w := new(WebService)
	if _, err := toml.Decode(string(src), &w); err != nil {
		return nil, err
	}
	if w.DocRoot == "" {
		w.DocRoot = "."
	}
	if w.Http != nil {
		w.Http.Scheme = "http"
	}
	if w.Https != nil {
		w.Https.Scheme = "https"
	}
	return w, nil
}

// LoadJSON loads a configuration of *WebService from a JSON file.
func LoadJSON(setup string) (*WebService, error) {
	src, err := ioutil.ReadFile(setup)
	if err != nil {
		return nil, err
	}
	w := new(WebService)
	if err := json.Unmarshal(src, &w); err != nil {
		return nil, err
	}
	if w.DocRoot == "" {
		w.DocRoot = "."
	}
	if w.Http != nil {
		w.Http.Scheme = "http"
	}
	if w.Https != nil {
		w.Https.Scheme = "https"
	}
	return w, nil
}

// Run() starts a web service(s) described in the *WebService struct.
func (w *WebService) Run() error {
	var err error
	if w.DocRoot == "" {
		w.DocRoot, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	log.Printf("Document root %s", w.DocRoot)
	if w.Http != nil {
		log.Printf("Listening for %s", w.Http.String())
	}
	if w.Https != nil {
		log.Printf("Listening for %s", w.Https.String())
	}
	fs := w.SafeFileSystem()
	mux := http.NewServeMux()
	//FIXME: Figure out how I want to stack up my handlers...

	// Setup our default file service handler.
	mux.Handle("/", RequestLogger(http.FileServer(fs)))

	// Run the configured services.
	switch {
	case w.Http != nil && w.Https != nil:
		// Run our http service in a go routine
		go func(addr string, handler http.Handler) {
			http.ListenAndServe(addr, handler)
		}(w.Http.Hostname(), mux)
		// Return our primar https service routine
		return http.ListenAndServeTLS(w.Https.Hostname(), w.Https.CertPEM, w.Https.KeyPEM, mux)
	case w.Https != nil:
		return http.ListenAndServeTLS(w.Https.Hostname(), w.Https.CertPEM, w.Https.KeyPEM, mux)
	case w.Http != nil:
		return http.ListenAndServe(w.Http.Hostname(), mux)
	default:
		return http.ListenAndServe(":8000", mux)
	}
}
