//
//
// Copyright Red Hat
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

package registry

import (
	"reflect"
	"testing"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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
			name: "Case 3: Storage not set, default set to false",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{},
			},
			want: false,
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
			name: "Case 3: Storage not set, default set to false",
			cr: registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: crName,
				},
				Spec: registryv1alpha1.DevfileRegistrySpec{},
			},
			want: corev1.VolumeSource{},
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

func TestGetDevfileIndexMemoryLimit(t *testing.T) {
	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want resource.Quantity
	}{
		{
			name: "Case 1: Memory Limit size set in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					DevfileIndex: registryv1alpha1.DevfileRegistrySpecContainer{
						MemoryLimit: "5Gi",
					},
				},
			},
			want: resource.MustParse("5Gi"),
		},
		{
			name: "Case 2:  Memory Limit size not set in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					DevfileIndex: registryv1alpha1.DevfileRegistrySpecContainer{},
				},
			},
			want: resource.MustParse(DefaultDevfileIndexMemoryLimit),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			volSize := GetDevfileIndexMemoryLimit(&tt.cr)
			if volSize != tt.want {
				t.Errorf("TestGetDevfileIndexMemoryLimit error: storage size mismatch, expected: %v got: %v", tt.want, volSize)
			}
		})
	}

}

func TestGetOCIRegistryMemoryLimit(t *testing.T) {
	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want resource.Quantity
	}{
		{
			name: "Case 1: Memory Limit size set in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					OciRegistry: registryv1alpha1.DevfileRegistrySpecContainer{
						MemoryLimit: "5Gi",
					},
				},
			},
			want: resource.MustParse("5Gi"),
		},
		{
			name: "Case 2:  Memory Limit size not set in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					OciRegistry: registryv1alpha1.DevfileRegistrySpecContainer{},
				},
			},
			want: resource.MustParse(DefaultOCIRegistryMemoryLimit),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			volSize := GetOCIRegistryMemoryLimit(&tt.cr)
			if volSize != tt.want {
				t.Errorf("TestGetOCIRegistryMemoryLimit error: storage size mismatch, expected: %v got: %v", tt.want, volSize)
			}
		})
	}

}

func TestGetRegistryViewerMemoryLimit(t *testing.T) {
	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want resource.Quantity
	}{
		{
			name: "Case 1: Memory Limit size set in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					RegistryViewer: registryv1alpha1.DevfileRegistrySpecContainer{
						MemoryLimit: "5Gi",
					},
				},
			},
			want: resource.MustParse("5Gi"),
		},
		{
			name: "Case 2:  Memory Limit size not set in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					RegistryViewer: registryv1alpha1.DevfileRegistrySpecContainer{},
				},
			},
			want: resource.MustParse(DefaultRegistryViewerMemoryLimit),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			volSize := GetRegistryViewerMemoryLimit(&tt.cr)
			if volSize != tt.want {
				t.Errorf("TestGetRegistryViewerMemoryLimit error: storage size mismatch, expected: %v got: %v", tt.want, volSize)
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

func Test_getDevfileRegistrySpecContainer(t *testing.T) {
	tests := []struct {
		name         string
		quantity     string
		defaultValue string
		want         resource.Quantity
	}{
		{
			name:         "Case 1: DevfileRegistrySpecContainer given correct quantity",
			quantity:     "256Mi",
			defaultValue: "512Mi",
			want:         resource.MustParse("256Mi"),
		},
		{
			name:         "Case 2: DevfileRegistrySpecContainer given correct quantity",
			quantity:     "test",
			defaultValue: "512Mi",
			want:         resource.MustParse("512Mi"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getDevfileRegistrySpecContainer(tt.quantity, tt.defaultValue)
			if result != tt.want {
				t.Errorf("func TestgetDevfileRegistrySpecContainer(t *testing.T) {\n error: enablement value mismatch, expected: %v got: %v", tt.want, result)
			}
		})
	}
}

func TestGetK8sIngressClass(t *testing.T) {
	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want string
	}{
		{
			name: "Case 1: K8s ingress class set",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					K8s: registryv1alpha1.DevfileRegistrySpecK8sOnly{
						IngressClass: "test",
					},
				},
			},
			want: "test",
		},
		{
			name: "Case 2: K8s ingress class not set",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Telemetry: registryv1alpha1.DevfileRegistrySpecTelemetry{},
				},
			},
			want: DefaultK8sIngressClass,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetK8sIngressClass(&tt.cr)
			if result != tt.want {
				t.Errorf("func TestGetK8sIngressClass(t *testing.T) {\n error: enablement value mismatch, expected: %v got: %v", tt.want, result)
			}
		})
	}
}
