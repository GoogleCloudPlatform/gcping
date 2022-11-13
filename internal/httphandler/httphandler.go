// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/GoogleCloudPlatform/gcping/internal/config"
)

// Options contains parameters for Handler.
type Options struct {
	// StaticRoot is the root for static serving content.
	StaticRoot http.FileSystem
	// Region is the region where the instance runs (e.g. us-west1).
	Region string
	// Endpoints is a list of available endpoints.
	Endpoints map[string]config.Endpoint
}

// Handler is a http.Handler implementation
type Handler struct {
	Options
	once    sync.Once
	handler http.Handler
}

// New returns a new intance of Handler based on opt.
func New(opts *Options) *Handler {
	s := &Handler{
		Options: *opts,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.StaticHandler())

	mux.HandleFunc("/api/endpoints", s.HandleEndpoints)
	mux.HandleFunc("/api/ping", s.HandlePing)

	// Serve /ping with region response to fix issue#96 on older cli versions.
	mux.HandleFunc("/ping", s.HandlePing)
	s.handler = mux
	return s
}

// ServeHTTP implements http.Handler.
func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

// StaticHandler returns a handler for static files at StaticRoot.
func (s *Handler) StaticHandler() http.HandlerFunc {
	h := http.FileServer(s.StaticRoot)
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: add HTST header to static files.
		// addHTSTHeader(w)
		h.ServeHTTP(w, r)
	}
}

// HandleEndpoints returns a list of available endpoints as JSON.
func (s *Handler) HandleEndpoints(w http.ResponseWriter, r *http.Request) {
	addHeaders(w)
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s.Endpoints); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// HandlePing returns the current region as a response for ping.
func (s *Handler) HandlePing(w http.ResponseWriter, r *http.Request) {
	addHeaders(w)
	s.once.Do(func() {
		w.Header().Add("X-First-Request", "true")
	})
	fmt.Fprintln(w, s.Region)
}

func addHTSTHeader(w http.ResponseWriter) {
	w.Header().Add("Strict-Transport-Security", "max-age=3600; includeSubdomains; preload")
}

func addHeaders(w http.ResponseWriter) {
	addHTSTHeader(w)
	w.Header().Add("Cache-Control", "no-store")
	w.Header().Add("Access-Control-Allow-Origin", "*")
}
