/*
Copyright 2025.

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

package controller

import (
	"context"

	issuesv1beta1 "github.com/CryptoRodeo/issues-operator/pkg/api/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// IssueReconciler reconciles a Issue object
type IssueReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=issues.konflux.dev,resources=issues,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=issues.konflux.dev,resources=issues/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=issues.konflux.dev,resources=issues/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Issue object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *IssueReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	issue := &issuesv1beta1.Issue{}
	err := r.Get(ctx, req.NamespacedName, issue)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if issue.Labels == nil {
		issue.Labels = make(map[string]string)
	}

	// Sync labels with spec fields
	needsUpdate := false

	// Sync severity
	if issue.Spec.Severity != "" {
		if issue.Labels["severity"] != issue.Spec.Severity {
			issue.Labels["severity"] = issue.Spec.Severity
			needsUpdate = true
		}
	}

	// Sync issueType
	if issue.Spec.IssueType != "" {
		if issue.Labels["issueType"] != issue.Spec.IssueType {
			issue.Labels["issueType"] = issue.Spec.IssueType
			needsUpdate = true
		}
	}

	// Sync resourceType
	if issue.Spec.Scope.ResourceType != "" {
		if issue.Labels["resourceType"] != issue.Spec.Scope.ResourceType {
			issue.Labels["resourceType"] = issue.Spec.Scope.ResourceType
			needsUpdate = true
		}
	}

	// Sync resourceName
	if issue.Spec.Scope.ResourceName != "" {
		if issue.Labels["resourceName"] != issue.Spec.Scope.ResourceName {
			issue.Labels["resourceName"] = issue.Spec.Scope.ResourceName
			needsUpdate = true
		}
	}

	// Sync resourceNamespace
	if issue.Spec.Scope.ResourceNamespace != "" {
		if issue.Labels["resourceNamespace"] != issue.Spec.Scope.ResourceNamespace {
			issue.Labels["resourceNamespace"] = issue.Spec.Scope.ResourceNamespace
			needsUpdate = true
		}
	}

	if needsUpdate {
		logger.Info("Updating issue labels", "name", issue.Name)
		if err := r.Update(ctx, issue); err != nil {
			logger.Error(err, "Failed to update issue labels")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IssueReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		For(&issuesv1beta1.Issue{}).
		Complete(r)
}
