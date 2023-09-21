/*
Copyright 2020-2023 Red Hat, Inc.

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

package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/devfile/registry-operator/tests/integration/pkg/config"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForPodRunningByLabelWithNamespace waits for the pod matching the specified label in a specified namespace to become running
// An error is returned if the pod does not exist or the timeout is reached
func (w *K8sClient) WaitForPodRunningByLabelWithNamespace(label string, namespace string) (deployed bool, err error) {
	timeout := time.After(6 * time.Minute)
	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-timeout:
			return false, errors.New("timed out")
		case <-tick:
			err := w.WaitForRunningPodBySelector(namespace, label, 3*time.Minute)
			if err == nil {
				return true, nil
			}
		}
	}
}

// WaitForPodFailedByLabelWithNamespace waits for the pod matching the specified label in a specified namespace to become running
// An error is returned if the pod does not exist or the timeout is reached
func (w *K8sClient) WaitForPodFailedByLabelWithNamespace(label string, namespace string) (deployed bool, err error) {
	timeout := time.After(6 * time.Minute)
	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-timeout:
			return false, errors.New("timed out")
		case <-tick:
			err := w.WaitForFailedPodBySelector(namespace, label, 3*time.Minute)
			if err == nil {
				return true, nil
			}
		}
	}
}

// WaitForPodFailedByLabel waits for the pod matching the specified label to become running
// An error is returned if the pod does not exist or the timeout is reached
func (w *K8sClient) WaitForPodFailedByLabel(label string) (deployed bool, err error) {
	return w.WaitForPodFailedByLabelWithNamespace(label, config.Namespace)
}

// WaitForPodRunningByLabel waits for the pod matching the specified label to become running
// An error is returned if the pod does not exist or the timeout is reached
func (w *K8sClient) WaitForPodRunningByLabel(label string) (deployed bool, err error) {
	return w.WaitForPodRunningByLabelWithNamespace(label, config.Namespace)
}

// WaitForRunningPodBySelector waits up to timeout seconds for all pods in 'namespace' with given 'selector' to enter running state.
// Returns an error if no pods are found or not all discovered pods enter running state.
func (w *K8sClient) WaitForRunningPodBySelector(namespace, selector string, timeout time.Duration) error {
	podList, err := w.ListPods(namespace, selector)
	if err != nil {
		return err
	}
	if len(podList.Items) == 0 {
		fmt.Println("Pod not created yet with selector " + selector + " in namespace " + namespace)

		return fmt.Errorf("Pod not created yet in %s with label %s", namespace, selector)
	}

	for _, pod := range podList.Items {
		fmt.Println("Pod " + pod.Name + " created in namespace " + namespace + "...Checking startup data.")
		if err := w.waitForPodRunning(namespace, pod.Name, timeout); err != nil {
			return err
		}
	}

	return nil
}

// WaitForFailedPodBySelector waits up to timeout seconds for all pods in 'namespace' with given 'selector' to enter running state.
// Returns an error if no pods are found or not all discovered pods enter running state.
func (w *K8sClient) WaitForFailedPodBySelector(namespace, selector string, timeout time.Duration) error {
	podList, err := w.ListPods(namespace, selector)
	if err != nil {
		return err
	}
	if len(podList.Items) == 0 {
		fmt.Println("No pods for " + selector + " in namespace " + namespace)

		return nil
	}

	for _, pod := range podList.Items {
		fmt.Println("Pod " + pod.Name + " created in namespace " + namespace + "...Checking for failure.")
		if err := w.waitForPodFailing(namespace, pod.Name, timeout); err != nil {
			return err
		}
	}

	return nil
}

// ListPods returns the list of currently scheduled or running pods in `namespace` with the given selector
func (w *K8sClient) ListPods(namespace, selector string) (*v1.PodList, error) {
	listOptions := metav1.ListOptions{LabelSelector: selector}
	podList, err := w.Kube().CoreV1().Pods(namespace).List(context.TODO(), listOptions)

	if err != nil {
		return nil, err
	}
	return podList, nil
}

// Poll up to timeout seconds for pod to enter running state.
// Returns an error if the pod never enters the running state.
func (w *K8sClient) waitForPodRunning(namespace, podName string, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, w.isPodRunning(podName, namespace))
}

// Poll up to timeout seconds for pod to enter running state.
// Returns an error if the pod never enters the running state.
func (w *K8sClient) waitForPodFailing(namespace, podName string, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, w.isPodFailing(podName, namespace))
}

// return a condition function that indicates whether the given pod is
// currently running
func (w *K8sClient) isPodRunning(podName, namespace string) wait.ConditionFunc {
	return func() (bool, error) {
		pod, _ := w.Kube().CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
		age := time.Since(pod.GetCreationTimestamp().Time).Seconds()

		switch pod.Status.Phase {
		case v1.PodRunning:
			fmt.Println("Pod started after", age, "seconds")
			return true, nil
		case v1.PodFailed, v1.PodSucceeded:
			return false, nil
		}
		return false, nil
	}
}

// return a condition function that indicates whether the given pod is
// currently running
func (w *K8sClient) isPodFailing(podName, namespace string) wait.ConditionFunc {
	return func() (bool, error) {
		pod, _ := w.Kube().CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})

		if pod.Status.Phase == v1.PodRunning {
			return false, nil
		} else {
			fmt.Printf("Pod terminated %s\n", podName)
			return true, nil
		}
	}
}
