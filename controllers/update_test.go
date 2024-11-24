package controllers

import (
	"testing"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func TestHeadlessStatusOutdated(t *testing.T) {
	var r *DevfileRegistryReconciler
	var cr *registryv1alpha1.DevfileRegistry

	r = &DevfileRegistryReconciler{}
	cr = &registryv1alpha1.DevfileRegistry{
		Spec: registryv1alpha1.DevfileRegistrySpec{
			Headless: func(b bool) *bool { return &b }(true),
		},
	}

	type args struct {
		cr  *registryv1alpha1.DevfileRegistry
		dep *appsv1.Deployment
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Cr headless true and deployment env REGISTRY_HEADLESS true",
			args: args{
				cr: cr,
				dep: func() *appsv1.Deployment {
					dep := &appsv1.Deployment{}
					dep.Spec.Template.Spec.Containers = []corev1.Container{
						{
							Env: []corev1.EnvVar{{Name: "REGISTRY_HEADLESS", Value: "true"}},
						},
					}
					return dep
				}(),
			},
			want: false,
		},
		{
			name: "Cr headless true and deployment env REGISTRY_HEADLESS false",
			args: args{
				cr: &registryv1alpha1.DevfileRegistry{
					Spec: registryv1alpha1.DevfileRegistrySpec{
						Headless: func(b bool) *bool { return &b }(true),
					},
				},
				dep: func() *appsv1.Deployment {
					dep := &appsv1.Deployment{}
					dep.Spec.Template.Spec.Containers = []corev1.Container{
						{
							Env: []corev1.EnvVar{{Name: "REGISTRY_HEADLESS", Value: "false"}},
						},
					}
					return dep
				}(),
			},
			want: true,
		},
		{
			name: "Cr headless false and deployment env REGISTRY_HEADLESS false",
			args: args{
				cr: &registryv1alpha1.DevfileRegistry{
					Spec: registryv1alpha1.DevfileRegistrySpec{
						Headless: func(b bool) *bool { return &b }(false),
					},
				},
				dep: func() *appsv1.Deployment {
					dep := &appsv1.Deployment{}
					dep.Spec.Template.Spec.Containers = []corev1.Container{
						{
							Env: []corev1.EnvVar{{Name: "REGISTRY_HEADLESS", Value: "false"}},
						},
					}
					return dep
				}(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := r.isHeadlessStatusOutdated(tt.args.cr, tt.args.dep); got != tt.want {
				t.Errorf("isHeadlessStatusOutdated() = %v, want %v", got, tt.want)
			}
		})
	}
}
