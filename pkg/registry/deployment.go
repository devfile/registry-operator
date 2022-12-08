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
							Env: []corev1.EnvVar{
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
					},
				},
			},
		},
	}

	// Set Registry Viewer if headless is false, else run headless mode
	if !IsHeadlessEnabled(cr) {
		dep.Spec.Template.Spec.Containers = append(dep.Spec.Template.Spec.Containers, corev1.Container{
			Image:           GetRegistryViewerImage(cr),
			ImagePullPolicy: corev1.PullAlways,
			Name:            "registry-viewer",
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
						Path: "/viewer",
						Port: intstr.FromInt(RegistryViewerPort),
					},
				},
			},
			ReadinessProbe: &corev1.Probe{
				Handler: corev1.Handler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/viewer",
						Port: intstr.FromInt(RegistryViewerPort),
					},
				},
			},
			StartupProbe: &corev1.Probe{
				Handler: corev1.Handler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/viewer",
						Port: intstr.FromInt(RegistryViewerPort),
					},
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "viewer-env-file",
					MountPath: "/app/apps/registry-viewer/.env.local",
					SubPath:   ".env.local",
				},
			},
		})
		dep.Spec.Template.Spec.Volumes = append(dep.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: "viewer-env-file",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: ConfigMapName(cr.Name),
					},
					Items: []corev1.KeyToPath{
						{
							Key:  ".env.registry-viewer",
							Path: ".env.local",
						},
					},
				},
			},
		})
	} else {
		// Set environment variable to run index server in headless mode
		dep.Spec.Template.Spec.Containers[0].Env = append(dep.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
			Name:  "REGISTRY_HEADLESS",
			Value: "true",
		})
	}

	// Set DevfileRegistry instance as the owner and controller
	_ = ctrl.SetControllerReference(cr, dep, scheme)
	return dep
}
