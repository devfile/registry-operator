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

	indexSchema "github.com/devfile/registry-support/index/generator/schema"
	registryLibrary "github.com/devfile/registry-support/registry-library/library"
	"github.com/hashicorp/go-multierror"
)

const (
	dupRegName      = "Duplicate registry name %s in registries list.  Ensure name is unique \n"
	dupURLName      = "Duplicate registry URL %s in registries list.  Ensure URL is unique \n"
	InvalidRegistry = "Devfile %s Registry is either invalid or unavailable, unable to add to the DevfileRegistryService list. Ensure you provide a valid Devfile Registry URL \n"
)

func validateURLs(devfileRegistries []DevfileRegistryService) (errors error) {
	processedName := make(map[string]bool)
	processedURL := make(map[string]bool)
	//validate URLs
	for i := range devfileRegistries {
		registry := devfileRegistries[i]
		url := registry.URL
		name := registry.Name

		//validate any duplicate registry names or urls
		if _, ok := processedName[name]; ok {

			err := fmt.Errorf(dupRegName, name)
			errors = multierror.Append(errors, err)
		}
		processedName[name] = true

		if _, ok := processedURL[url]; ok {
			err := fmt.Errorf(dupURLName, url)
			errors = multierror.Append(errors, err)
		}
		processedURL[url] = true

		err := IsRegistryValid(registry.SkipTLSVerify, url)
		if err != nil {
			errors = multierror.Append(errors, err)
		}
	}

	return errors
}

// IsRegistryValid determines if the given DevfileRegistryService.URL returns a
// well-formed v1 or v2 index schema

func IsRegistryValid(skipTLSVerify bool, url string) error {
	//call the registry library to determine if URL is a valid devfile registry
	registryOptions := registryLibrary.RegistryOptions{}
	if skipTLSVerify {
		registryOptions.SkipTLSVerify = skipTLSVerify
	}

	//Validate that url is a supported registry
	//try with a v1 index
	_, err := registryLibrary.GetRegistryIndex(url, registryOptions, indexSchema.SampleDevfileType, indexSchema.StackDevfileType)
	if err != nil {
		//try with a v2index
		registryOptions.NewIndexSchema = true
		_, err = registryLibrary.GetRegistryIndex(url, registryOptions, indexSchema.SampleDevfileType, indexSchema.StackDevfileType)
		if err != nil {
			err := fmt.Errorf(InvalidRegistry, url)
			return err
		}
	}

	return nil
}

// IsNamespaceValid determines if given namespace for deployment
// is valid.
func IsNamespaceValid(namespace string) error {
	if namespace == "default" {
		return fmt.Errorf("devfile registry deployment namespace should never be 'default'.")
	}

	return nil
}
