/*
Copyright 2020-2022 Red Hat, Inc.

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

package cluster

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// Borrowed from https://github.com/devfile/devworkspace-operator/blob/master/internal/cluster/info.go

// IsOpenShift returns true if the operator is running on an OpenShift cluster
func IsOpenShift() (bool, error) {
	kubeCfg, err := config.GetConfig()
	if err != nil {
		return false, err
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(kubeCfg)
	if err != nil {
		return false, err
	}
	apiList, err := discoveryClient.ServerGroups()
	if err != nil {
		return false, err
	}
	if findAPIGroup(apiList.Groups, "route.openshift.io") == nil {
		return false, nil
	} else {
		return true, nil
	}
}

func findAPIGroup(source []metav1.APIGroup, apiName string) *metav1.APIGroup {
	for i := 0; i < len(source); i++ {
		if source[i].Name == apiName {
			return &source[i]
		}
	}
	return nil
}
