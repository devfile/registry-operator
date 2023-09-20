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

package controllers

import (
	"context"
	"fmt"

	"github.com/devfile/registry-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
)

// SetValidateDevfileRegistriesConditionAndUpdateCR sets the condition of the cluster devfile registries list validation
func (r *ClusterDevfileRegistriesListReconciler) SetValidateDevfileRegistriesConditionAndUpdateCR(ctx context.Context, req ctrl.Request, clusterDevfileRegistriesList *v1alpha1.ClusterDevfileRegistriesList, validateError error) error {
	log := ctrl.LoggerFrom(ctx)
	var (
		condition metav1.Condition
		err       error
	)

	if validateError == nil {
		condition = metav1.Condition{
			Type:    typeValidateDevfileRegistries,
			Status:  metav1.ConditionTrue,
			Reason:  "Ready",
			Message: allRegistriesReachable,
		}
	} else {
		condition = metav1.Condition{
			Type:    typeValidateDevfileRegistries,
			Status:  metav1.ConditionFalse,
			Reason:  "NotReady",
			Message: fmt.Sprintf("Cluster devfile registries list failed to validate: %v", validateError),
		}
	}

	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		validateDevfileRegistriesAndUpdateCondition(clusterDevfileRegistriesList.Spec.DevfileRegistries, condition, func(c metav1.Condition) {
			meta.SetStatusCondition(&clusterDevfileRegistriesList.Status.Conditions, c)
		})
		return r.Status().Update(ctx, clusterDevfileRegistriesList)
	})

	if err != nil {
		log.Error(err, "Unable to update cluster devfile registries list")
		return err
	}

	return nil
}
