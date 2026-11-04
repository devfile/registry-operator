//
//
// Copyright Red Hat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-logr/logr"
	authenticationv1 "k8s.io/api/authentication/v1"
	authorizationv1 "k8s.io/api/authorization/v1"
	"k8s.io/client-go/rest"
)

func TestSecureMetricsAuthorization(t *testing.T) {
	// Model the API server's TokenReview and SubjectAccessReview endpoints.
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/apis/authentication.k8s.io/v1/tokenreviews":
			var review authenticationv1.TokenReview
			if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
				t.Error(err)
				w.WriteHeader(400)
				return
			}
			review.Status.Authenticated = review.Spec.Token == "reader" || review.Spec.Token == "denied"
			if review.Status.Authenticated {
				review.Status.User.Username = review.Spec.Token
			}
			if err := json.NewEncoder(w).Encode(review); err != nil {
				t.Error(err)
			}
		case "/apis/authorization.k8s.io/v1/subjectaccessreviews":
			var review authorizationv1.SubjectAccessReview
			if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
				t.Error(err)
				w.WriteHeader(400)
				return
			}
			attrs := review.Spec.NonResourceAttributes
			review.Status.Allowed = review.Spec.User == "reader" && attrs != nil && attrs.Path == "/metrics" && attrs.Verb == "get"
			if err := json.NewEncoder(w).Encode(review); err != nil {
				t.Error(err)
			}
		default:
			t.Errorf("unexpected API request: %s", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer api.Close()
	options := newMetricsOptions(":8443", true, nil)
	if !options.SecureServing || options.FilterProvider == nil {
		t.Fatal("HTTPS metrics must require authentication and authorization")
	}
	filter, err := options.FilterProvider(&rest.Config{Host: api.URL}, api.Client())
	if err != nil {
		t.Fatal(err)
	}
	handler, err := filter(logr.Discard(), http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, token string
		status      int
	}{
		{"anonymous", "", http.StatusUnauthorized},
		// controller-runtime v0.19 reports invalid bearer token errors as HTTP 500.
		{"invalid token", "invalid", http.StatusInternalServerError},
		{"authenticated without permission", "denied", http.StatusForbidden},
		{"authorized reader", "reader", http.StatusOK},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "https://metrics.example/metrics", nil)
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, req)
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d: %s", response.Code, tc.status, response.Body.String())
			}
		})
	}
}
