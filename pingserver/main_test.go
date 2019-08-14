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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGreetHandler(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		{
			name: "successful greet",
			want: []byte("hello"),
		},
	}

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(greetHandler)
			handler.ServeHTTP(w, r)

			resp := w.Result()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Unexpected status code %d", resp.StatusCode)
			}

			if diff := cmp.Diff(w.Body.Bytes(), tt.want); diff != "" {
				t.Errorf("greetHandler returned unexpected body (-got +want):\n%s", diff)
			}
		})
	}
}

func TestPingHandler(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		{
			name: "successful ping",
			want: []byte("pong"),
		},
	}

	r, err := http.NewRequest(http.MethodGet, "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(pingHandler)
			handler.ServeHTTP(w, r)

			resp := w.Result()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Unexpected status code %d", resp.StatusCode)
			}

			if diff := cmp.Diff(w.Body.Bytes(), tt.want); diff != "" {
				t.Errorf("pingHandler returned unexpected body (-got +want):\n%s", diff)
			}
		})
	}
}

func TestSetupHandler(t *testing.T) {
	tests := []struct {
		name           string
		endpoint       string
		wantStatusCode int
	}{
		{
			name:           "successful greet - HTTP 200",
			endpoint:       "/",
			wantStatusCode: 200,
		},
		{
			name:           "successful ping - HTTP 200",
			endpoint:       "/ping",
			wantStatusCode: 200,
		},
		{
			name:           "wrong endpoint - HTTP 404",
			endpoint:       "/foobar",
			wantStatusCode: 404,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r, err := http.NewRequest(http.MethodGet, tc.endpoint, nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			handler := setupHandlers()
			handler.ServeHTTP(w, r)

			resp := w.Result()
			if resp.StatusCode != tc.wantStatusCode {
				t.Errorf("Unexpected status code %d", resp.StatusCode)
			}

		})
	}
}
