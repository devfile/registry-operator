/*
Copyright 2022 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"fmt"
	"github.com/devfile/registry-operator/pkg/test"
	"github.com/hashicorp/go-multierror"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevfileRegistriesValidateURL(t *testing.T) {

	testServer := test.GetNewUnstartedTestServer()
	assert.NotNil(t, testServer)
	testServer.Start()
	defer testServer.Close()

	tests := []struct {
		name              string
		devfileRegistries []DevfileRegistryService
		wantErr           []string
	}{
		{
			name: "Registries list with valid registries",
			devfileRegistries: []DevfileRegistryService{
				{
					Name: devfileStagingRegistryName,
					URL:  devfileStagingRegistryURL,
				},
			},
			wantErr: []string{},
		},
		{
			name: "Registries list with invalid registries",
			devfileRegistries: []DevfileRegistryService{
				{
					Name: "Bad URL",
					URL:  "https://registry.stage.devfilex.io",
				},
			},
			wantErr: []string{
				fmt.Sprintf(InvalidRegistry, "https://registry.stage.devfilex.io"),
			},
		},
		{
			name: "Registries list with duplicate names and URLs",
			devfileRegistries: []DevfileRegistryService{
				{
					Name: devfileStagingRegistryName,
					URL:  devfileStagingRegistryURL,
				},
				{
					Name: devfileStagingRegistryName,
					URL:  devfileStagingRegistryURL,
				},
			},
			wantErr: []string{
				fmt.Sprintf(dupRegName, devfileStagingRegistryName),
				fmt.Sprintf(dupURLName, devfileStagingRegistryURL),
			},
		},
		{
			name: "Registries list with valid v1 and v2 indices",
			devfileRegistries: []DevfileRegistryService{
				{
					Name: devfileStagingRegistryName,
					URL:  devfileStagingRegistryURL,
				},
				{
					Name: "local v2index server",
					URL:  testServer.URL,
				},
			},
			wantErr: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURLs(tt.devfileRegistries)
			if merr, ok := err.(*multierror.Error); ok && tt.wantErr != nil {
				assert.Equal(t, len(tt.wantErr), len(merr.Errors), fmt.Sprintf("Errors do not match = %v, want %v", err, tt.wantErr))
				for _, testErr := range tt.wantErr {
					assert.ErrorContains(t, err, testErr)
				}
			} else {
				assert.Equal(t, nil, err, "Error should be nil")
			}
		})
	}
}
