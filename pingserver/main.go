// Copyright 2019 Google Inc. All Rights Reserved.
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

package main

import (
	"log"
	"net/http"
	"os"
)

// greetHandler responds with an HTTP 200 status code and a message when the service is running.
func greetHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("hello"))
}

// pingHandler responds with an HTTP 200 status code and a message when the service is running.
func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("pong"))
}

// setupHandlers sets the routes for the server.
func setupHandlers() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", greetHandler)
	mux.HandleFunc("/ping", pingHandler)
	return mux
}

func main() {
	mux := setupHandlers()
	port := os.Getenv("PORT") // makes it portable to Cloud Run
	if port == "" {
		port = "80"
	}
	log.Printf("Server listening on *:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
