//
// Package wsfn provides a common library of functions and structures for
// working with web services in Caltech Library projects and software.
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2021, Caltech
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
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"

	// 3rd Party packages
	"github.com/BurntSushi/toml"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/pbkdf2"
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
		// Check to see if we have a *.mjs JavaScript module.
		if ext := path.Ext(r.URL.Path); ext == ".mjs" {
			w.Header().Set("Content-Type", "text/javascript")
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

//
// NOTE: moved from redirects.go into wsfn.go
//

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

//
// NOTE: merged from cors.go into wsfn.go
//

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

// Handler accepts an http.Handler and returns a http.Handler. It
// Wraps the response with the CORS headers based on configuration
// in CORSPolicy struct. If cors is nil then it passes thru
// to next http.Handler unaltered.
func (cors *CORSPolicy) Handler(next http.Handler) http.Handler {
	if cors == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
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

//
// NOTE: Merged from access.go into wsfn.go
//

//
// access.go holds authentication related stucts and funcs.
// It includes those functions needed by the web service but
// also some funcs for things like generating/managing content
// of an access.toml file.
//

// Access holds the necessary configuration for doing
// basic auth authentication.
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication
// using Go's http.Request object.
type Access struct {
	// AuthType (e.g. Basic)
	AuthType string `json:"auth_type" toml:"auth_type"`
	// AuthName (e.g. string describing authorization, e.g. realm in basic auth)
	AuthName string `json:"auth_name" toml:"auth_name"`
	// Encryption is a string describing the encryption used
	// e.g. argon2id, pbkds2, md5 or sha512
	Encryption string `json:"encryption" toml:"encryption"`
	// Map holds a user to secret map. It is usually populated
	// after reading in the users file with LoadAccessTOML() or
	// LoadAccessJSON().
	Map map[string]*Secrets `json:"access" toml:"access"`
	// Routes is a list of URL path prefixes covered by
	// this Access control object.
	Routes []string `json:"routes" toml:"routes"`
}

type Secrets struct {
	// NOTE: salt is needed by Argon2 and pbkdb2.
	// If the toml/json file functions as the database then
	// this file MUST be kept safe with restricted permissions.
	// If not you just gave away your system a cracker.
	Salt []byte `json:"salt,omitempty" toml:"salt,omitempty"`
	// Key holds the salted hash ...
	Key []byte `json:"key, omitempty" toml:"key,omitempty"`
}

// LoadAccess loads a TOML or JSON access file.
func LoadAccess(fName string) (*Access, error) {
	switch {
	case strings.HasSuffix(fName, ".toml"):
		return loadAccessTOML(fName)
	case strings.HasSuffix(fName, ".json"):
		return loadAccessJSON(fName)
	default:
		return nil, fmt.Errorf("%q, unsupported format", fName)
	}
}

// loadAccessTOML loads a TOML acces file.
// and returns an Access struct and error.
func loadAccessTOML(accessTOML string) (*Access, error) {
	auth := new(Access)
	src, err := ioutil.ReadFile(accessTOML)
	if err != nil {
		return nil, err
	}
	if _, err := toml.Decode(string(src), &auth); err != nil {
		return nil, err
	}
	return auth, nil
}

// loadAccessJSON loads a JSON access file.
// and returns an Access struct and error.
func loadAccessJSON(accessJSON string) (*Access, error) {
	auth := new(Access)
	src, err := ioutil.ReadFile(accessJSON)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(src, &auth); err != nil {
		return nil, err
	}
	return auth, nil
}

// DumpAccess writes a access file.
func (a *Access) DumpAccess(fName string) error {
	switch {
	case strings.HasSuffix(fName, ".toml"):
		return a.dumpAccessTOML(fName)
	case strings.HasSuffix(fName, ".json"):
		return a.dumpAccessJSON(fName)
	default:
		return fmt.Errorf("%q, unsupported format", fName)
	}
}

// dumpAccessTOML writes a TOML access file.
func (a *Access) dumpAccessTOML(accessTOML string) error {
	buf := new(bytes.Buffer)
	tomlEncoder := toml.NewEncoder(buf)
	if err := tomlEncoder.Encode(a); err != nil {
		return err
	}
	return ioutil.WriteFile(accessTOML, buf.Bytes(), 0600)
}

// dumpAccessJSON writes an access.toml file.
func (a *Access) dumpAccessJSON(accessJSON string) error {
	src, err := json.MarshalIndent(a, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(accessJSON, src, 0600)
}

// UpdateAccess uses an *Access and username, password
// generates a salt and then adds username, salt
// and secret to .Map (creating one if needed)
func (a *Access) UpdateAccess(username string, password string) bool {
	if a.Map == nil {
		a.Map = make(map[string]*Secrets)
	}
	// Pick the preferred encryption if not set.
	if a.Encryption == "" {
		a.Encryption = "argon2id"
	}
	secret := new(Secrets)
	secret.Salt = make([]byte, 32)
	_, err := rand.Read(secret.Salt)
	if err != nil {
		return false
	}
	switch a.Encryption {
	case "argon2id":
		secret.Key = argon2.IDKey([]byte(password), secret.Salt, 1, 64*1024, 4, 32)
		a.Map[username] = secret
		return true
	case "pbkdf2":
		secret.Key = pbkdf2.Key([]byte(password), secret.Salt, 4097, 32, sha1.New)
		a.Map[username] = secret
		return true
	case "md5":
		h := md5.New()
		io.WriteString(h, password)
		secret.Key = h.Sum(nil)
		a.Map[username] = secret
		return true
	case "sha512":
		h := sha512.New()
		secret.Key = h.Sum([]byte(password))
		a.Map[username] = secret
		return true
	}
	// NOTE: We don't know the encryption scheme
	// so we fail to authenticate.
	return false
}

// RemoveAccess takes an *Access and username and
// deletes the username from .Map
// returns true if delete applied, false if user not found in map
func (a *Access) RemoveAccess(username string) bool {
	if _, ok := a.Map[username]; ok == true {
		delete(a.Map, username)
		return true
	}
	return false
}

// Login accepts username, password and ok boolean.
// Returns true if they match auth's settings false otherwise.
//
// How to choosing an appropriate hash method see
//
// https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
//
// md5 and sha512 are included for historic reasons
// They are NOT considered secure anymore as they are breakable
// with brute force using today's CPU/GPUs.
func (a *Access) Login(username string, password string) bool {
	var (
		u      *Secrets
		secret *Secrets
	)

	// Make sure we know about the user, others we can't validate
	if val, ok := a.Map[username]; ok {
		u = val
	} else {
		return false
	}
	secret = new(Secrets)
	switch a.Encryption {
	case "argon2id":
		secret.Key = argon2.IDKey([]byte(password), u.Salt, 1, 64*1024, 4, 32)
	case "pbkdf2":
		secret.Key = pbkdf2.Key([]byte(password), u.Salt, 4097, 32, sha1.New)
	case "md5":
		h := md5.New()
		io.WriteString(h, password)
		secret.Key = h.Sum(nil)
	case "sha512":
		h := sha512.New()
		secret.Key = h.Sum([]byte(password))
	default:
		// NOTE: We don't know the encryption scheme
		// so we fail to authenticate.
		return false
	}
	if bytes.Compare(secret.Key, u.Key) == 0 {
		return true
	}
	return false
}

// Checks to see if we have a defined route.
func (a *Access) isAccessRoute(p string) bool {
	for _, route := range a.Routes {
		if strings.HasPrefix(p, route) {
			return true
		}
	}
	return false
}

// GetUsername takes an Request object, inspects the headers
// and returns the username if possible, otherwise an error.
func (a *Access) GetUsername(r *http.Request) (string, error) {
	switch a.AuthType {
	case "basic":
		username, _, ok := r.BasicAuth()
		if ok == true {
			return username, nil
		}
		return "", fmt.Errorf("No user info found")
	default:
		return "", fmt.Errorf("Unsupported Auth Type")
	}
}

// Handler takes a handler and returns handler. If
// *Access is null it pass thru unchanged. Otherwise
// it applies the access policy.
func (a *Access) Handler(next http.Handler) http.Handler {
	if a == nil {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(res, req)
		})
	}
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if a.isAccessRoute(req.URL.Path) {
			res.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, a.AuthName))
			// Check to see if we've previously authenticated.
			username, password, ok := req.BasicAuth()
			if ok == false {
				http.Error(res, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if a.Login(username, password) == false {
				http.Error(res, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(res, req)
	})
}

// AccessHandler is a wrapping handler that checks if
// Access.Routes matches the req.URL.Path and if so
// applies access contraints. If *Access is nil then
// it just passes through to the next handler.
func AccessHandler(next http.Handler, a *Access) http.Handler {
	if a == nil {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(res, req)
		})
	}
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if a.isAccessRoute(req.URL.Path) {
			res.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, a.AuthName))
			// Check to see if we've previously authenticated.
			username, password, ok := req.BasicAuth()
			if ok == false {
				http.Error(res, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if a.Login(username, password) == false {
				http.Error(res, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(res, req)
	})
}

//
// NOTE: merged from defaults.go into wsfn.go
//

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

//
// NOTE: merged from json.go into wsfn.go
//

// jsonResponse enforces a common JSON response write handling.
// It takes a response writer and request plus a struct that can
// be converted to JSON.
func jsonResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	src, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Printf("json marshal error, %s %s", r.URL.Path, err)
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(src); err != nil {
		return
	}
	log.Printf("FIXME: Log successful requests here ... %s", r.URL.Path)
}

//
// NOTE: merged from license.go into wsfn.go
//

const (
	LicenseText = `Copyright (c) 2019, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`
)

//
// NOTE: merged from logger.go into wsfn.go
//

// RequestLogger logs the request based on the request object passed into
// it.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if len(q) > 0 {
			log.Printf("request Method: %s Path: %s RemoteAddr: %s UserAgent: %s Query: %+v\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), q)
		} else {
			log.Printf("request Method: %s Path: %s RemoteAddr: %s UserAgent: %s\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
		}
		next.ServeHTTP(w, r)
	})
}

// ResponseLogger logs the response based on a request, status and error
// message
func ResponseLogger(r *http.Request, status int, err error) {
	q := r.URL.Query()
	if len(q) > 0 {
		log.Printf("response Method: %s Path: %s RemoteAddr: %s UserAgent: %s Query: %+v Status: %d, %s %q\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), q, status, http.StatusText(status), err)
	} else {
		log.Printf("response Method: %s Path: %s RemoteAddr: %s UserAgent: %s Status: %d, %s %q\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), status, http.StatusText(status), err)
	}
}

//
// NOTE: merged from safefilesystem.go into wsfn.go
//

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

//
// NOTE: merged from server.go into wsfn.go
//

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

	// AccessFile holds a name of an access file to load and
	// populate .Access from.
	AccessFile string `json:"access_file,omitempty" toml:"access_file,omitempty"`

	// Access adds access related features to the service.
	// E.g. BasicAUTH support.
	Access *Access `json:"access,omitempty" toml:"access,omitempty"`

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
	CertPEM string `json:"cert_pem,omitempty" toml:"cert_pem,omitempty"`
	// KeyPEM describes the location of the key.pem used for TLS support
	KeyPEM string `json:"key_pem,omitempty" toml:"key_pem,omitempty"`
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

// LoadWebService loads a configuration file of *WebService
func LoadWebService(setup string) (*WebService, error) {
	var (
		ws  *WebService
		err error
	)

	switch {
	case strings.HasSuffix(setup, ".toml"):
		ws, err = loadWebServiceTOML(setup)
	case strings.HasSuffix(setup, ".json"):
		ws, err = loadWebServiceJSON(setup)
	default:
		err = fmt.Errorf("%q, unknown format.", setup)
	}
	if err != nil {
		return nil, err
	}
	// If AccessFile set is set overwrite .Access ...
	if ws.AccessFile != "" {
		ws.Access, err = LoadAccess(ws.AccessFile)
	}
	return ws, err
}

// loadWebServiceTOML loads a *WebService from a TOML file.
func loadWebServiceTOML(setup string) (*WebService, error) {
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

// loadWebServiceJSON loads a *WebService from a JSON file.
func loadWebServiceJSON(setup string) (*WebService, error) {
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

// DumpWebService writes a access file.
func (ws *WebService) DumpWebService(fName string) error {
	var (
		access *Access
		err    error
	)
	if ws.AccessFile != "" && ws.Access != nil {
		access = ws.Access
		ws.Access = nil
	}
	switch {
	case strings.HasSuffix(fName, ".toml"):
		err = ws.dumpWebServiceTOML(fName)
	case strings.HasSuffix(fName, ".json"):
		err = ws.dumpWebServiceJSON(fName)
	default:
		err = fmt.Errorf("%q, unsupported format", fName)
	}
	if access != nil {
		ws.Access = access
	}
	return err
}

// dumpWebServiceTOML writes a TOML file.
func (ws *WebService) dumpWebServiceTOML(fName string) error {
	buf := new(bytes.Buffer)
	tomlEncoder := toml.NewEncoder(buf)
	if err := tomlEncoder.Encode(ws); err != nil {
		return err
	}
	return ioutil.WriteFile(fName, buf.Bytes(), 0600)
}

// dumpWebServiceJSON writes a JSON file.
func (ws *WebService) dumpWebServiceJSON(fName string) error {
	src, err := json.MarshalIndent(ws, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fName, src, 0600)
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

	// Setup our Safe file system handler.
	fs, err := w.SafeFileSystem()
	if err != nil {
		return err
	}

	//FIXME: Figure out a better way to stack up handlers...
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(fs))

	// Run the configured services.
	switch {
	case w.Http != nil && w.Https != nil:
		// Run our http service in a go routine
		go func() {
			http.ListenAndServe(w.Http.Hostname(), RequestLogger(AccessHandler(mux, w.Access)))
		}()
		// Return our primary https service routine
		return http.ListenAndServeTLS(w.Https.Hostname(), w.Https.CertPEM, w.Https.KeyPEM, RequestLogger(AccessHandler(mux, w.Access)))
	case w.Https != nil:
		return http.ListenAndServeTLS(w.Https.Hostname(), w.Https.CertPEM, w.Https.KeyPEM, RequestLogger(AccessHandler(mux, w.Access)))
	case w.Http != nil:
		return http.ListenAndServe(w.Http.Hostname(), RequestLogger(AccessHandler(mux, w.Access)))
	default:
		return http.ListenAndServe(":8000", RequestLogger(AccessHandler(mux, w.Access)))
	}
}
