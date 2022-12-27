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
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/GoogleCloudPlatform/gcping/internal/config"
	"github.com/google/go-cmp/cmp"
)

func TestPing(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		wantHeader http.Header
	}{
		{
			name: "first time",
			wantHeader: http.Header{
				"Content-Type":                {"text/plain; charset=utf-8"},
				"Cache-Control":               {"no-store"},
				"Access-Control-Allow-Origin": {"*"},
				"X-First-Request":             {"true"},
				"Strict-Transport-Security":   {"max-age=3600; includeSubdomains; preload"},
			},
		},
		{
			name: "second time",
			wantHeader: http.Header{
				"Content-Type":                {"text/plain; charset=utf-8"},
				"Cache-Control":               {"no-store"},
				"Access-Control-Allow-Origin": {"*"},
				"Strict-Transport-Security":   {"max-age=3600; includeSubdomains; preload"},
			},
		},
	}

	handler := New(&Options{Region: "test-region"})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "https://gcping.com/api/ping", nil)
			w := httptest.NewRecorder()
			handler.HandlePing(w, req)
			resp := w.Result()
			t.Cleanup(func() { resp.Body.Close() })

			if got, want := resp.StatusCode, http.StatusOK; got != want {
				t.Errorf("HandlePing() Status Code: got %v, want %v", got, want)
			}
			if diff := cmp.Diff(tc.wantHeader, resp.Header); diff != "" {
				t.Errorf("HandlePing() Header (-want, +got):\n%s", diff)
			}
			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Failed to read response body: %v", err)
			}
			if want := "test-region\n"; string(got) != want {
				t.Errorf("HandlePing() = got %s, want %s", string(got), want)
			}
		})
	}
}

func TestRouting(t *testing.T) {
	t.Parallel()

	fakeStatic := fstest.MapFS{
		"index.html": &fstest.MapFile{
			Data: []byte("<html></html>"),
		},
	}
	handler := New(&Options{
		StaticRoot: http.FS(fakeStatic),
		Region:     "test-region",
		Endpoints: map[string]config.Endpoint{
			"test-region": {
				URL:        "https://test-region",
				Region:     "test-region",
				RegionName: "Test Region",
			},
		},
	})
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)
	client := ts.Client()

	testCases := []struct {
		path     string
		wantCode int
		wantBody string
	}{
		{"/", http.StatusOK, "<html></html>"},
		{"/index.html", http.StatusOK, "<html></html>"},
		{"/api", http.StatusNotFound, "404 page not found\n"},
		{"/api/ping", http.StatusOK, "test-region\n"},
		{"/api/endpoints", http.StatusOK, `{"test-region":{"URL":"https://test-region","Region":"test-region","RegionName":"Test Region"}}` + "\n"},
		{"/ping", http.StatusOK, "test-region\n"},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			t.Parallel()

			resp, err := client.Get(ts.URL + tc.path)
			if err != nil {
				t.Fatalf("Get() failed: %v", err)
			}
			t.Cleanup(func() { resp.Body.Close() })
			if got := resp.StatusCode; got != tc.wantCode {
				t.Errorf("Get() Status Code = got %d, want %d", got, tc.wantCode)
			}
			got, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("ReadAll() failed for response body: %v", err)
			}
			if diff := cmp.Diff(tc.wantBody, string(got)); diff != "" {
				t.Errorf("Get() response body = (-want, +got):\n%s", diff)
			}
		})
	}
}
