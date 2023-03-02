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

package registry

import (
	"fmt"

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
	allowPrivilegeEscalation := false
	runAsNonRoot := true
	runAsUser := int64(1001)
	runAsGroup := int64(2001)
	fsGroup := int64(3001)

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
							SecurityContext: &corev1.SecurityContext{
								AllowPrivilegeEscalation: &allowPrivilegeEscalation,
								RunAsNonRoot:             &runAsNonRoot,
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{"ALL"},
								},
								SeccompProfile: &corev1.SeccompProfile{
									Type: "RuntimeDefault",
								},
							},
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
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/health",
										Port: intstr.FromInt(DevfileIndexPort),
									},
								},
								InitialDelaySeconds: 15,
								PeriodSeconds:       10,
								TimeoutSeconds:      3,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/health",
										Port: intstr.FromInt(DevfileIndexPort),
									},
								},
								InitialDelaySeconds: 15,
								PeriodSeconds:       10,
								TimeoutSeconds:      3,
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
							SecurityContext: &corev1.SecurityContext{
								AllowPrivilegeEscalation: &allowPrivilegeEscalation,
								RunAsNonRoot:             &runAsNonRoot,
								Capabilities: &corev1.Capabilities{
									Drop: []corev1.Capability{"ALL"},
								},
								SeccompProfile: &corev1.SeccompProfile{
									Type: "RuntimeDefault",
								},
							},
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
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/v2",
										Port: intstr.FromInt(OCIServerPort),
									},
								},
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
								TimeoutSeconds:      3,
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/v2",
										Port: intstr.FromInt(OCIServerPort),
									},
								},
								InitialDelaySeconds: 3,
								PeriodSeconds:       10,
								TimeoutSeconds:      3,
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
		dep.Spec.Template.Spec.Containers[0].StartupProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: "/viewer",
					Port: intstr.FromInt(RegistryViewerPort),
				},
			},
			InitialDelaySeconds: 30,
			PeriodSeconds:       10,
			TimeoutSeconds:      3,
		}
		dep.Spec.Template.Spec.Containers = append(dep.Spec.Template.Spec.Containers, corev1.Container{
			Image:           GetRegistryViewerImage(cr),
			ImagePullPolicy: corev1.PullAlways,
			Name:            "registry-viewer",
			SecurityContext: &corev1.SecurityContext{
				AllowPrivilegeEscalation: &allowPrivilegeEscalation,
				RunAsNonRoot:             &runAsNonRoot,
				Capabilities: &corev1.Capabilities{
					Drop: []corev1.Capability{"ALL"},
				},
				SeccompProfile: &corev1.SeccompProfile{
					Type: "RuntimeDefault",
				},
			},
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
			Env: []corev1.EnvVar{
				{
					Name:  "NEXT_PUBLIC_ANALYTICS_WRITE_KEY",
					Value: cr.Spec.Telemetry.RegistryViewerWriteKey,
				},
				{
					Name: "DEVFILE_REGISTRIES",
					Value: fmt.Sprintf(`
					[
						{
							"name": "%s",
							"url": "http://localhost:8080",
							"fqdn": "%s"
						}
					]`, cr.ObjectMeta.Name, cr.Status.URL),
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "viewer-env-file",
					MountPath: "/app/.env.production",
					SubPath:   ".env.production",
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
							Path: ".env.production",
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

	// Enables podspec security context if storage is enabled
	if cr.Spec.Storage.Enabled == nil || *cr.Spec.Storage.Enabled {
		dep.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{
			RunAsNonRoot: &runAsNonRoot,
			RunAsUser:    &runAsUser,
			RunAsGroup:   &runAsGroup,
			FSGroup:      &fsGroup,
		}
	}

	// Set DevfileRegistry instance as the owner and controller
	_ = ctrl.SetControllerReference(cr, dep, scheme)
	return dep
}
