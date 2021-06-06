package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/GoogleCloudPlatform/gcping/internal/config"
)

var once sync.Once

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
	http.Handle("/", http.FileServer(http.Dir(kdp)))

	http.HandleFunc("/endpoints", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-store")
		w.Header().Add("Content-Type", "application/json;charset=utf-8")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Strict-Transport-Security", "max-age=3600; includeSubdomains; preload")
		err := json.NewEncoder(w).Encode(config.AllEndpoints)
		if err != nil {
			w.WriteHeader(500)
		}
	})

	// Serve /ping with region response.
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
