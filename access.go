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
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/BurntSushi/toml"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/pbkdf2"
)

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
