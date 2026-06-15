//
// reverse_proxy_test.go - Tests for reverse proxy functionality
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2024, Caltech
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
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
)

// ============================================================================
// ReverseProxyService Tests
// ============================================================================

func TestNewReverseProxyService(t *testing.T) {
	rps := NewReverseProxyService()
	if rps == nil {
		t.Fatal("Expected ReverseProxyService to be created")
	}
	if rps.routes == nil {
		t.Error("Expected routes map to be initialized")
	}
}

func TestReverseProxyService_HasReverseProxyRoutes(t *testing.T) {
	tests := []struct {
		name     string
		routes   map[string]*ReverseProxyTarget
		expected bool
	}{
		{
			name:     "empty routes",
			routes:   map[string]*ReverseProxyTarget{},
			expected: false,
		},
		{
			name: "single route",
			routes: map[string]*ReverseProxyTarget{
				"/api/": {TargetURL: mustParseURL("http://localhost:9000/")},
			},
			expected: true,
		},
		{
			name: "multiple routes",
			routes: map[string]*ReverseProxyTarget{
				"/api/":  {TargetURL: mustParseURL("http://localhost:9000/")},
				"/auth/": {TargetURL: mustParseURL("http://localhost:9001/")},
			},
			expected: true,
		},
		{
			name:     "nil routes",
			routes:   nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rps := &ReverseProxyService{routes: tt.routes}
			result := rps.HasReverseProxyRoutes()
			if result != tt.expected {
				t.Errorf("HasReverseProxyRoutes() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestReverseProxyService_AddReverseProxyRoute(t *testing.T) {
	tests := []struct {
		name        string
		prefix      string
		targetURL   string
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid route",
			prefix:     "/api/",
			targetURL:  "http://localhost:9000/",
			wantErr:    false,
		},
		{
			name:       "valid route without trailing slash",
			prefix:     "/api",
			targetURL:  "http://localhost:9000",
			wantErr:    false,
		},
		{
			name:       "invalid URL",
			prefix:     "/api/",
			targetURL:  "://invalid",
			wantErr:    true,
			errContains: "invalid URL",
		},
		{
			name:       "colliding route - prefix subset",
			prefix:     "/api",
			targetURL:  "http://localhost:9000/",
			wantErr:    true,
			errContains: "collide",
		},
		{
			name:       "colliding route - subset prefix",
			prefix:     "/api/users/",
			targetURL:  "http://localhost:9000/",
			wantErr:    true,
			errContains: "collide",
		},
		{
			name:       "empty prefix",
			prefix:     "",
			targetURL:  "http://localhost:9000/",
			wantErr:    true,
			errContains: "empty prefix",
		},
		{
			name:       "empty target URL",
			prefix:     "/api/",
			targetURL:  "",
			wantErr:    true,
			errContains: "empty target URL",
		},
		{
			name:       "non-http(s) scheme",
			prefix:     "/api/",
			targetURL:  "ftp://localhost:9000/",
			wantErr:    true,
			errContains: "unsupported scheme",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rps := NewReverseProxyService()
			
			// Add a base route for collision tests
			if strings.Contains(tt.name, "colliding") {
				if strings.Contains(tt.prefix, "/api") {
					rps.AddReverseProxyRoute("/api/", "http://localhost:9000/")
				}
			}
			
			err := rps.AddReverseProxyRoute(tt.prefix, tt.targetURL)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("AddReverseProxyRoute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.wantErr && tt.errContains != "" && err != nil {
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("Error message %q does not contain %q", err.Error(), tt.errContains)
				}
			}
			
			if !tt.wantErr {
				// Verify the route was added
				if len(rps.routes) == 0 {
					t.Error("Expected route to be added")
				}
				
				// Check that the target URL was parsed correctly
				if target, ok := rps.routes[tt.prefix]; ok {
					if target.TargetURL.String() != tt.targetURL {
						t.Errorf("Target URL = %v, expected %v", target.TargetURL.String(), tt.targetURL)
					}
					if target.Proxy == nil {
						t.Error("Expected proxy to be initialized")
					}
				}
			}
		})
	}
}

func TestReverseProxyService_Route(t *testing.T) {
	tests := []struct {
		name      string
		addPrefix string
		addURL    string
		lookUpPrefix string
		wantMatch bool
		wantURL   string
	}{
		{
			name:        "matching route",
			addPrefix:    "/api/",
			addURL:       "http://localhost:9000/",
			lookUpPrefix: "/api/",
			wantMatch:    true,
			wantURL:      "http://localhost:9000/",
		},
		{
			name:        "non-matching route",
			addPrefix:    "/api/",
			addURL:       "http://localhost:9000/",
			lookUpPrefix: "/other/",
			wantMatch:    false,
			wantURL:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rps := NewReverseProxyService()
			rps.AddReverseProxyRoute(tt.addPrefix, tt.addURL)
			
			target, ok := rps.Route(tt.lookUpPrefix)
			
			if ok != tt.wantMatch {
				t.Errorf("Route() match = %v, wantMatch %v", ok, tt.wantMatch)
				return
			}
			
			if tt.wantMatch && target.TargetURL.String() != tt.wantURL {
				t.Errorf("Route() URL = %v, wantURL %v", target.TargetURL.String(), tt.wantURL)
			}
		})
	}
}

// ============================================================================
// ReverseProxyRouter Tests
// ============================================================================

func TestReverseProxyRouter_ProxyRequest(t *testing.T) {
	// Create a test backend server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Backend-Test", "true")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Path: %s, Query: %s", r.URL.Path, r.URL.RawQuery)
	}))
	defer backend.Close()

	tests := []struct {
		name           string
		prefix         string
		targetURL      string
		requestPath    string
		requestQuery   string
		expectedPath   string
		expectProxy    bool
		expectStatus   int
		expectHeader   string
	}{
		{
			name:         "proxy with path stripping",
			prefix:       "/api/",
			targetURL:    backend.URL,
			requestPath:  "/api/users",
			requestQuery: "",
			expectedPath: "/users",
			expectProxy:  true,
			expectStatus: http.StatusOK,
			expectHeader: "true",
		},
		{
			name:         "proxy with query parameters",
			prefix:       "/api/",
			targetURL:    backend.URL,
			requestPath:  "/api/search",
			requestQuery: "q=test&page=1",
			expectedPath: "/search",
			expectProxy:  true,
			expectStatus: http.StatusOK,
			expectHeader: "true",
		},
		{
			name:         "non-matching path passes through",
			prefix:       "/api/",
			targetURL:    backend.URL,
			requestPath:  "/static/file.html",
			requestQuery: "",
			expectedPath: "/static/file.html",
			expectProxy:  false,
			expectStatus: http.StatusNotFound,
			expectHeader: "",
		},
		{
			name:         "root path with slash",
			prefix:       "/api/",
			targetURL:    backend.URL + "/",
			requestPath:  "/api/",
			requestQuery: "",
			expectedPath: "/",
			expectProxy:  true,
			expectStatus: http.StatusOK,
			expectHeader: "true",
		},
		{
			name:         "multiple route prefixes",
			prefix:       "/api/",
			targetURL:    backend.URL,
			requestPath:  "/api/v1/data",
			requestQuery: "",
			expectedPath: "/v1/data",
			expectProxy:  true,
			expectStatus: http.StatusOK,
			expectHeader: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rps := NewReverseProxyService()
			err := rps.AddReverseProxyRoute(tt.prefix, tt.targetURL)
			if err != nil {
				t.Fatalf("Failed to add route: %v", err)
			}

			router := rps.ReverseProxyRouter(http.NotFoundHandler())

			// Create request
			reqURL := fmt.Sprintf("http://localhost%s%s", tt.requestPath, tt.requestQuery)
			req, err := http.NewRequest("GET", reqURL, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(rr, req)

			// Check response
			if rr.Code != tt.expectStatus {
				t.Errorf("Status code = %v, expectStatus %v", rr.Code, tt.expectStatus)
			}

			if tt.expectProxy {
				headerValue := rr.Header().Get("X-Backend-Test")
				if headerValue != tt.expectHeader {
					t.Errorf("Backend header = %v, expectHeader %v", headerValue, tt.expectHeader)
				}
			}

			// For non-proxied requests, we should get 404 from our NotFoundHandler
			if !tt.expectProxy && rr.Code != http.StatusNotFound {
				t.Errorf("Non-proxied request should return 404, got %d", rr.Code)
			}
		})
	}
}

func TestReverseProxyRouter_MultipleRoutes(t *testing.T) {
	// Create test backends
	backend1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Backend1: %s", r.URL.Path)
	}))
	defer backend1.Close()

	backend2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Backend2: %s", r.URL.Path)
	}))
	defer backend2.Close()

	rps := NewReverseProxyService()
	rps.AddReverseProxyRoute("/api/", backend1.URL)
	rps.AddReverseProxyRoute("/auth/", backend2.URL)

	router := rps.ReverseProxyRouter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Static: %s", r.URL.Path)
	}))

	tests := []struct {
		path     string
		expected string
	}{
		{"/api/users", "Backend1: /users"},
		{"/auth/login", "Backend2: /login"},
		{"/static/file.html", "Static: /static/file.html"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://localhost"+tt.path, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Status code = %v, expected %v", rr.Code, http.StatusOK)
			}

			body := rr.Body.String()
			if body != tt.expected {
				t.Errorf("Body = %v, expected %v", body, tt.expected)
			}
		})
	}
}

func TestReverseProxyRouter_PreserveMethod(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Method", r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	rps := NewReverseProxyService()
	rps.AddReverseProxyRoute("/api/", backend.URL)

	router := rps.ReverseProxyRouter(http.NotFoundHandler())

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req, _ := http.NewRequest(method, "http://localhost/api/test", nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Status code = %v, expected %v for method %s", rr.Code, http.StatusOK, method)
			}

			headerValue := rr.Header().Get("X-Method")
			if headerValue != method {
				t.Errorf("Method = %v, expected %v", headerValue, method)
			}
		})
	}
}

func TestReverseProxyRouter_PreserveHeaders(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test-Header", r.Header.Get("X-Test-Header"))
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	rps := NewReverseProxyService()
	rps.AddReverseProxyRoute("/api/", backend.URL)

	router := rps.ReverseProxyRouter(http.NotFoundHandler())

	testHeaders := map[string]string{
		"X-Test-Header":   "test-value",
		"Authorization":   "Bearer token123",
		"Content-Type":   "application/json",
		"X-Custom-Header": "custom-value",
	}

	for header, value := range testHeaders {
		t.Run(header, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "http://localhost/api/test", nil)
			req.Header.Set(header, value)
			
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Status code = %v, expected %v", rr.Code, http.StatusOK)
			}

			responseHeader := rr.Header().Get("X-Test-Header")
			if header == "X-Test-Header" && responseHeader != value {
				t.Errorf("Header %s = %v, expected %v", header, responseHeader, value)
			}
		})
	}
}

func TestReverseProxyRouter_RequestBody(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.Header().Set("X-Body", string(body))
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	rps := NewReverseProxyService()
	rps.AddReverseProxyRoute("/api/", backend.URL)

	router := rps.ReverseProxyRouter(http.NotFoundHandler())

	testBodies := []string{
		"",
		"test data",
		`{"key": "value"}`,
		"large body " + strings.Repeat("x", 1000),
	}

	for _, body := range testBodies {
		t.Run(string(body[:min(20, len(body))]), func(t *testing.T) {
			req, _ := http.NewRequest("POST", "http://localhost/api/test", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Status code = %v, expected %v", rr.Code, http.StatusOK)
			}

			responseBody := rr.Header().Get("X-Body")
			if responseBody != body {
				t.Errorf("Body = %v, expected %v", responseBody, body)
			}
		})
	}
}

func TestReverseProxyRouter_ResponseHeaders(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Backend-Status", "ok")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Custom-Header", "custom-value")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status": "ok"}`)
	}))
	defer backend.Close()

	rps := NewReverseProxyService()
	rps.AddReverseProxyRoute("/api/", backend.URL)

	router := rps.ReverseProxyRouter(http.NotFoundHandler())

	req, _ := http.NewRequest("GET", "http://localhost/api/test", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Status code = %v, expected %v", rr.Code, http.StatusOK)
	}

	// Check that backend headers are preserved
	expectedHeaders := map[string]string{
		"X-Backend-Status": "ok",
		"Content-Type":    "application/json",
		"X-Custom-Header":  "custom-value",
	}

	for header, expected := range expectedHeaders {
		actual := rr.Header().Get(header)
		if actual != expected {
			t.Errorf("Header %s = %v, expected %v", header, actual, expected)
		}
	}
}

// ============================================================================
// WebService Integration Tests
// ============================================================================

func TestWebService_LoadReverseProxyFromTOML(t *testing.T) {
	// Create a temporary TOML file with reverse proxy config
	omlContent := `
htdocs = "."

[http]
host = "localhost"
port = "8000"

[reverse_proxy]
"/api/" = "http://localhost:9000/"
"/auth/" = "http://localhost:9001/"
`

	// Write to temp file
	tmpFile := createTempFile(t, omlContent, "*.toml")
	defer removeTempFile(t, tmpFile)

	ws, err := LoadWebService(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load web service: %v", err)
	}

	if ws.ReverseProxy == nil {
		t.Fatal("Expected ReverseProxy to be initialized")
	}

	if !ws.ReverseProxy.HasReverseProxyRoutes() {
		t.Error("Expected reverse proxy routes to be loaded")
	}

	// Check specific routes
	target, ok := ws.ReverseProxy.Route("/api/")
	if !ok {
		t.Error("Expected /api/ route to exist")
	} else {
		if target.TargetURL.String() != "http://localhost:9000/" {
			t.Errorf("API target URL = %v, expected http://localhost:9000/", target.TargetURL.String())
		}
	}

	target, ok = ws.ReverseProxy.Route("/auth/")
	if !ok {
		t.Error("Expected /auth/ route to exist")
	} else {
		if target.TargetURL.String() != "http://localhost:9001/" {
			t.Errorf("Auth target URL = %v, expected http://localhost:9001/", target.TargetURL.String())
		}
	}
}

func TestWebService_LoadReverseProxyFromJSON(t *testing.T) {
	// Create a temporary JSON file with reverse proxy config
	jsonContent := `{
		"htdocs": ".",
		"http": {
			"host": "localhost",
			"port": "8000"
		},
		"reverse_proxy": {
			"/api/": "http://localhost:9000/",
			"/auth/": "http://localhost:9001/"
		}
	}`

	// Write to temp file
	tmpFile := createTempFile(t, jsonContent, "*.json")
	defer removeTempFile(t, tmpFile)

	ws, err := LoadWebService(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load web service: %v", err)
	}

	if ws.ReverseProxy == nil {
		t.Fatal("Expected ReverseProxy to be initialized")
	}

	if !ws.ReverseProxy.HasReverseProxyRoutes() {
		t.Error("Expected reverse proxy routes to be loaded")
	}

	// Check specific routes
	target, ok := ws.ReverseProxy.Route("/api/")
	if !ok {
		t.Error("Expected /api/ route to exist")
	} else {
		if target.TargetURL.String() != "http://localhost:9000/" {
			t.Errorf("API target URL = %v, expected http://localhost:9000/", target.TargetURL.String())
		}
	}
}

func TestWebService_DumpReverseProxyToTOML(t *testing.T) {
	ws := DefaultWebService()
	ws.ReverseProxy = NewReverseProxyService()
	ws.ReverseProxy.AddReverseProxyRoute("/api/", "http://localhost:9000/")
	ws.ReverseProxy.AddReverseProxyRoute("/auth/", "http://localhost:9001/")

	tmpFile := createTempFile(t, "", "*.toml")
	defer removeTempFile(t, tmpFile)

	err := ws.DumpWebService(tmpFile)
	if err != nil {
		t.Fatalf("Failed to dump web service: %v", err)
	}

	// Read the file and verify it contains reverse proxy config
	content := readTempFile(t, tmpFile)
	
	if !strings.Contains(content, "[reverse_proxy]") {
		t.Error("Expected TOML to contain [reverse_proxy] section")
	}
	
	if !strings.Contains(content, "/api/") {
		t.Error("Expected TOML to contain /api/ route")
	}
	
	if !strings.Contains(content, "http://localhost:9000/") {
		t.Error("Expected TOML to contain target URL")
	}
}

func TestWebService_DumpReverseProxyToJSON(t *testing.T) {
	ws := DefaultWebService()
	ws.ReverseProxy = NewReverseProxyService()
	ws.ReverseProxy.AddReverseProxyRoute("/api/", "http://localhost:9000/")

	tmpFile := createTempFile(t, "", "*.json")
	defer removeTempFile(t, tmpFile)

	err := ws.DumpWebService(tmpFile)
	if err != nil {
		t.Fatalf("Failed to dump web service: %v", err)
	}

	// Read the file and verify it contains reverse proxy config
	content := readTempFile(t, tmpFile)
	
	if !strings.Contains(content, "reverse_proxy") {
		t.Error("Expected JSON to contain reverse_proxy field")
	}
	
	if !strings.Contains(content, "/api/") {
		t.Error("Expected JSON to contain /api/ route")
	}
	
	if !strings.Contains(content, "http://localhost:9000/") {
		t.Error("Expected JSON to contain target URL")
	}
}

func TestWebService_DefaultInitIncludesReverseProxyExample(t *testing.T) {
	initContent := DefaultInit()
	
	if !strings.Contains(string(initContent), "[reverse_proxy]") {
		t.Error("Expected default init to contain [reverse_proxy] example")
	}
	
	if !strings.Contains(string(initContent), "/api/") {
		t.Error("Expected default init to contain /api/ example")
	}
	
	if !strings.Contains(string(initContent), "http://localhost:9000/") {
		t.Error("Expected default init to contain example target URL")
	}
}

// ============================================================================
// Helper functions
// ============================================================================

func mustParseURL(u string) *url.URL {
	parsed, err := url.Parse(u)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse URL %s: %v", u, err))
	}
	return parsed
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func createTempFile(t *testing.T, content, pattern string) string {
	t.Helper()
	
	// Use the existing test file creation approach
	tmpFile, err := os.CreateTemp("", pattern)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	
	if content != "" {
		_, err = tmpFile.WriteString(content)
		if err != nil {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
			t.Fatalf("Failed to write to temp file: %v", err)
		}
	}
	
	tmpFile.Close()
	return tmpFile.Name()
}

func removeTempFile(t *testing.T, filename string) {
	t.Helper()
	os.Remove(filename)
}

func readTempFile(t *testing.T, filename string) string {
	t.Helper()
	
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}
	
	return string(content)
}
