/*
Copyright 2024.

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
    "fmt"
	"path"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/errors"
    "k8s.io/apimachinery/pkg/types"
    "sigs.k8s.io/controller-runtime/pkg/client"
    // "sigs.k8s.io/controller-runtime/pkg/log"

    danav1alpha1 "github.com/TalDebi/namespacelabel-assignment.git/api/v1alpha1"
)

// NamespaceLabelReconciler reconciles a NamespaceLabel object
type NamespaceLabelReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=dana.dana.io,resources=namespacelabels,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dana.dana.io,resources=namespacelabels/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=dana.dana.io,resources=namespacelabels/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NamespaceLabel object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile
func (r *NamespaceLabelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// log := log.Log.WithValues("namespacelabel", req.NamespacedName)

	// Fetch the NamespaceLabel instance
	var namespaceLabel danav1alpha1.NamespaceLabel
	if err := r.Get(ctx, req.NamespacedName, &namespaceLabel); err != nil {
		if errors.IsNotFound(err) {
			// Handle deletion
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Ensure only one NamespaceLabel per namespace
	existingLabels := &danav1alpha1.NamespaceLabelList{}
	if err := r.List(ctx, existingLabels, client.InNamespace(req.Namespace)); err != nil {
		return ctrl.Result{}, err
	}
	if len(existingLabels.Items) > 1 {
		return ctrl.Result{}, fmt.Errorf("only one NamespaceLabel allowed per namespace")
	}

	// Fetch the namespace
	var namespace corev1.Namespace
	if err := r.Get(ctx, types.NamespacedName{Name: req.Namespace}, &namespace); err != nil {
		return ctrl.Result{}, err
	}

	// Protect existing labels
	protectedLabels := []string{"kubernetes.io/*", "mycompany.com/protected"}
	labelsToUpdate := filterLabels(namespaceLabel.Spec.Labels, protectedLabels)

	// Update namespace labels
	namespace.Labels = mergeLabels(namespace.Labels, labelsToUpdate)
	if err := r.Update(ctx, &namespace); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func filterLabels(labels map[string]string, protected []string) map[string]string {
    filtered := make(map[string]string)
    for key, value := range labels {
        isProtected := false
        for _, p := range protected {
            if matched, _ := path.Match(p, key); matched {
                isProtected = true
                break
            }
        }
        if !isProtected {
            filtered[key] = value
        }
    }
    return filtered
}

func mergeLabels(existing, updates map[string]string) map[string]string {
    for key, value := range updates {
        existing[key] = value
    }
    return existing
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceLabelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&danav1alpha1.NamespaceLabel{}).
		Complete(r)
}
