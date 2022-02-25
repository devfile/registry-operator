//
// Copyright (c) 2020-2022 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation

package registry

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
)

func GenerateDeployment(cr *registryv1alpha1.DevfileRegistry, scheme *runtime.Scheme, labels map[string]string) *appsv1.Deployment {
	replicas := int32(1)

	dep := &appsv1.Deployment{
		ObjectMeta: generateObjectMeta(cr.Name, cr.Namespace, labels),
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image:           cr.Spec.DevfileIndexImage,
							ImagePullPolicy: corev1.PullAlways,
							Name:            "devfile-registry",
							Ports: []corev1.ContainerPort{{
								ContainerPort: DevfileIndexPort,
							}},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("250m"),
									corev1.ResourceMemory: resource.MustParse("64Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("500m"),
									corev1.ResourceMemory: resource.MustParse("256Mi"),
								},
							},
							LivenessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/health",
										Port: intstr.FromInt(DevfileIndexPort),
									},
								},
							},
							ReadinessProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/health",
										Port: intstr.FromInt(DevfileIndexPort),
									},
								},
							},
							StartupProbe: &corev1.Probe{
								Handler: corev1.Handler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/viewer",
										Port: intstr.FromInt(DevfileIndexPort),
									},
								},
								InitialDelaySeconds: 30,
								PeriodSeconds:       1,
								TimeoutSeconds:      1,
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "viewer-config",
									MountPath: "/app/config",
									ReadOnly:  false,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "DEVFILE_VIEWER_ROOT",
									Value: "/viewer",
								},
								{
									Name:  "REGISTRY_NAME",
									Value: cr.Spec.Telemetry.RegistryName,
								},
								{
									Name:  "TELEMETRY_KEY",
									Value: cr.Spec.Telemetry.Key,
								},
							},
						},
						{
							Image: GetOCIRegistryImage(cr),
							Name:  "oci-registry",
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("100m"),
									corev1.ResourceMemory: resource.MustParse("64Mi"),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("500m"),
									corev1.ResourceMemory: resource.MustParse("256Mi"),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      DevfileRegistryVolumeName,
									MountPath: "/var/lib/registry",
								},
								{
									Name:      "config",
									MountPath: "/etc/docker/registry",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name:         DevfileRegistryVolumeName,
							VolumeSource: GetDevfileRegistryVolumeSource(cr),
						},
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: ConfigMapName(cr.Name),
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "registry-config.yml",
											Path: "config.yml",
										},
									},
								},
							},
						},
						{
							Name: "viewer-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: ConfigMapName(cr.Name),
									},
									Items: []corev1.KeyToPath{
										{
											Key:  "devfile-registry-hosts.json",
											Path: "devfile-registry-hosts.json",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	// Set Memcached instance as the owner and controller
	ctrl.SetControllerReference(cr, dep, scheme)
	return dep
}
