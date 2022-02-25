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
	"reflect"
	"testing"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsTLSEnabled(t *testing.T) {
	tlsEnabled := true
	tlsDisabled := false

	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want bool
	}{
		{
			name: "Case 1: TLS enabled in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					TLS: registryv1alpha1.DevfileRegistrySpecTLS{
						Enabled: &tlsEnabled,
					},
				},
			},
			want: true,
		},
		{
			name: "Case 2: TLS disabled in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					TLS: registryv1alpha1.DevfileRegistrySpecTLS{
						Enabled: &tlsDisabled,
					},
				},
			},
			want: false,
		},
		{
			name: "Case 3: TLS not set, default set to true",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlsSetting := IsTLSEnabled(&tt.cr)
			if tlsSetting != tt.want {
				t.Errorf("TestIsTLSEnabled error: tls value mismatch, expected: %v got: %v", tt.want, tlsSetting)
			}
		})
	}

}

func TestIsStorageEnabled(t *testing.T) {
	storageEnabled := true
	storageDisabled := false

	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want bool
	}{
		{
			name: "Case 1: Storage enabled in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Storage: registryv1alpha1.DevfileRegistrySpecStorage{
						Enabled: &storageEnabled,
					},
				},
			},
			want: true,
		},
		{
			name: "Case 2: Storage disabled in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Storage: registryv1alpha1.DevfileRegistrySpecStorage{
						Enabled: &storageDisabled,
					},
				},
			},
			want: false,
		},
		{
			name: "Case 3: Storage not set, default set to true",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlsSetting := IsStorageEnabled(&tt.cr)
			if tlsSetting != tt.want {
				t.Errorf("TestIsStorageEnabled error: storage value mismatch, expected: %v got: %v", tt.want, tlsSetting)
			}
		})
	}

}

func TestGetDevfileRegistryVolumeSource(t *testing.T) {
	storageEnabled := true
	storageDisabled := false
	crName := "devfileregistry-test"

	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want corev1.VolumeSource
	}{
		{
			name: "Case 1: Storage enabled in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: crName,
				},
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Storage: registryv1alpha1.DevfileRegistrySpecStorage{
						Enabled: &storageEnabled,
					},
				},
			},
			want: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: PVCName(crName),
				},
			},
		},
		{
			name: "Case 2: Storage disabled in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: crName,
				},
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Storage: registryv1alpha1.DevfileRegistrySpecStorage{
						Enabled: &storageDisabled,
					},
				},
			},
			want: corev1.VolumeSource{},
		},
		{
			name: "Case 3: Storage not set, default set to true",
			cr: registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: crName,
				},
				Spec: registryv1alpha1.DevfileRegistrySpec{},
			},
			want: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: PVCName(crName),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlsSetting := GetDevfileRegistryVolumeSource(&tt.cr)
			if !reflect.DeepEqual(tlsSetting, tt.want) {
				t.Errorf("TestGetDevfileRegistryVolumeSource error: storage source mismatch, expected: %v got: %v", tt.want, tlsSetting)
			}
		})
	}

}

func TestGetDevfileRegistryVolumeSize(t *testing.T) {
	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want string
	}{
		{
			name: "Case 1: Volume size set in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Storage: registryv1alpha1.DevfileRegistrySpecStorage{
						RegistryVolumeSize: "5Gi",
					},
				},
			},
			want: "5Gi",
		},
		{
			name: "Case 2: Volume size not set in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{},
			},
			want: DefaultDevfileRegistryVolumeSize,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			volSize := getDevfileRegistryVolumeSize(&tt.cr)
			if volSize != tt.want {
				t.Errorf("TestGetDevfileRegistryVolumeSize error: storage size mismatch, expected: %v got: %v", tt.want, volSize)
			}
		})
	}

}

func TestIsTelemetryEnabled(t *testing.T) {
	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want bool
	}{
		{
			name: "Case 1: Telemetry key not specified",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Telemetry: registryv1alpha1.DevfileRegistrySpecTelemetry{
						RegistryName: "test",
					},
				},
			},
			want: false,
		},
		{
			name: "Case 2: Telemetry key specified",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Telemetry: registryv1alpha1.DevfileRegistrySpecTelemetry{
						RegistryName: "test",
						Key:          "abcdef",
					},
				},
			},
			want: true,
		},
		{
			name: "Case 3: Telemetry key specified but is empty",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Telemetry: registryv1alpha1.DevfileRegistrySpecTelemetry{
						RegistryName: "test",
						Key:          "",
					},
				},
			},
			want: false,
		},
		{
			name: "Case 4: Telemetry object is empty",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enabled := IsTelemetryEnabled(&tt.cr)
			if enabled != tt.want {
				t.Errorf("func TestIsTelemetryEnabled(t *testing.T) {\n error: enablement value mismatch, expected: %v got: %v", tt.want, enabled)
			}
		})
	}

}
