//
//
// Copyright 2022 Red Hat, Inc.
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

package test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	indexSchema "github.com/devfile/registry-support/index/generator/schema"
)

// ListType is used in common test functions to distinguish the flow of logic between the Cluster and Namespace CRs
type ListType string

const (
	ClusterListType   ListType = "ClusterDevfileRegistriesList"
	NamespaceListType ListType = "DevfileRegistriesList"
	ApiVersion                 = "registry.devfile.io/v1alpha1"
	Timeout                    = time.Second * 10
	Interval                   = time.Millisecond * 250
)

// GetNewUnstartedTestServer is a mock test index server that supports just the v2 index schema
func GetNewUnstartedTestServer() *httptest.Server {
	notFilteredV2Index := []indexSchema.Schema{
		{
			Name: "v2index",
		},
	}

	testServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data []indexSchema.Schema
		var err error
		// Current support is to allow user to specify the entire index which will include stacks and samples
		if strings.Contains(r.URL.String(), "/v2index/all") {
			data = notFilteredV2Index
		} else {
			http.NotFound(w, r)
		}

		bytes, err := json.MarshalIndent(&data, "", "  ")
		if err != nil {
			log.Fatalf("Unexpected error while doing json marshal: %v", err)
			return
		}

		_, err = w.Write(bytes)
		if err != nil {
			log.Fatalf("Unexpected error while writing data: %v", err)
		}
	}))

	return testServer
}
