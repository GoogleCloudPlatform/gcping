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

package config

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEndpointsFromServer(t *testing.T) {
	const endpointsPath = "/api/endpoints"
	testCases := []struct {
		name    string
		body    string
		code    int
		wantErr bool
		want    map[string]Endpoint
	}{
		{
			name:    "valid",
			code:    http.StatusOK,
			body:    `{"test-region":{"URL":"https://test-region","Region":"test-region","RegionName":"Test Region"}}`,
			wantErr: false,
			want: map[string]Endpoint{
				"test-region": {
					URL:        "https://test-region",
					Region:     "test-region",
					RegionName: "Test Region",
				},
			},
		},
		{
			name:    "no results",
			code:    http.StatusOK,
			body:    "{}",
			wantErr: false,
			want:    map[string]Endpoint{},
		},
		{
			name:    "empty string",
			code:    http.StatusOK,
			wantErr: true,
		},
		{
			name:    "syntax error",
			code:    http.StatusOK,
			body:    "invalid",
			wantErr: true,
		},
		{
			name:    "server error",
			code:    http.StatusInternalServerError,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var fakeHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
				if got, want := r.Method, http.MethodGet; got != want {
					t.Errorf("Handler: method got %s, want %s", got, want)
				}
				if got, want := r.URL.EscapedPath(), endpointsPath; got != want {
					t.Errorf("Handler: url got %s, want %s", got, want)
				}

				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(tc.code)
				_, err := io.WriteString(w, tc.body)
				if err != nil {
					t.Errorf("Handler: WriteString() failed: %v", err)
				}
			}

			ts := httptest.NewServer(fakeHandler)
			t.Cleanup(ts.Close)
			got, err := EndpointsFromServer(context.Background(), ts.URL+endpointsPath)
			if got := (err != nil); got != tc.wantErr {
				t.Errorf("EndpointsFromServer(): got error %v, want %v", got, tc.wantErr)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("EndpointsFromServer() = (-want, +got):\n%s", diff)
			}
		})
	}
}
