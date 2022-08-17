// Copyright 2021 Google LLC
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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/GoogleCloudPlatform/gcping/internal/config"
	"github.com/patrickmn/go-cache"
)

var once sync.Once

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Serving on :%s", port)

	// Set up cache with default 5 minutes retentio
	cache := cache.New(5*time.Minute, 5*time.Minute)
	region := os.Getenv("REGION")
	if region == "" {
		region = "pong"
	}

	// Prefetch and cache endpoint map
	_, _ = endpointsCache(cache)

	// Serve / from files in kodata.
	kdp := os.Getenv("KO_DATA_PATH")
	if kdp == "" {
		log.Println("KO_DATA_PATH unset")
		kdp = "/var/run/ko/"
	}
	http.Handle("/", http.FileServer(http.Dir(kdp)))

	http.HandleFunc("/api/endpoints", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-store")
		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Strict-Transport-Security", "max-age=3600; includeSubdomains; preload")
		endpoints, err := endpointsCache(cache)
		if err != nil {
			w.WriteHeader(500)
		}
		err = json.NewEncoder(w).Encode(endpoints)
		if err != nil {
			w.WriteHeader(500)
		}
	})

	// Serve /api/ping with region response.
	http.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-store")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Strict-Transport-Security", "max-age=3600; includeSubdomains; preload")
		once.Do(func() {
			w.Header().Add("X-First-Request", "true")
		})
		fmt.Fprintln(w, region)
	})

	// Serve /ping with region response to fix issue#96 on older cli versions.
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-store")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Strict-Transport-Security", "max-age=3600; includeSubdomains; preload")
		once.Do(func() {
			w.Header().Add("X-First-Request", "true")
		})
		fmt.Fprintln(w, region)
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// endpointsCache is used to fetch a value from the local cache if available
// and from the Cloud Rur Admin API in case it is not available, or has expired.
func endpointsCache(c *cache.Cache) (map[string]config.Endpoint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	em, found := c.Get("endpoints")
	if !found {
		e, err := config.GenerateConfigFromCloudRunAPI(ctx)
		if err != nil {
			//Retrieve stale_endpoints when the fresh endpoints are unavailable
			em, found := c.Get("stale_endpoints")
			if !found {
				log.Println("cannot fetch endpoint map")
				return nil, fmt.Errorf("cannot fetch map")
			}
			return em.(map[string]config.Endpoint), nil
		}
		c.Set("endpoints", e, cache.DefaultExpiration)
		c.Set("stale_endpoints", e, 0)
		return e, nil
	}

	return em.(map[string]config.Endpoint), nil
}
