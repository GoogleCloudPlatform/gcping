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
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/gcping/internal/config"
	"github.com/GoogleCloudPlatform/gcping/internal/httphandler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Serving on :%s", port)

	region := os.Getenv("REGION")
	if region == "" {
		region = "pong"
	}

	// Serve / from files in kodata.
	kdp := os.Getenv("KO_DATA_PATH")
	if kdp == "" {
		log.Println("KO_DATA_PATH unset")
		kdp = "/var/run/ko/"
	}

	handler := httphandler.New(&httphandler.Options{
		Region:     region,
		StaticRoot: http.Dir(kdp),
		Endpoints:  config.AllEndpoints,
	})

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("ListenAndServe(): %v", err)
	}
	log.Print("Exiting.")
}
