//
// Copyright (c) 2019-2020 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//

package util

import (
	"crypto/tls"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

// Poll up to timeout seconds for pod to enter running state.
// Returns an error if the pod never enters the running state.
func WaitForServer(url string, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, func() (bool, error) {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get(url)
		if err != nil {
			return false, err
		}
		if resp.StatusCode == 200 {
			return true, nil
		}
		return false, nil
	})
}
