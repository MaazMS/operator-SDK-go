/*
Copyright 2021.

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

package controllers

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	test1v1alpha1 "github.com/MaazMS/operator-SDK-go/api/v1alpha1"
)

// VisitorsAppReconciler reconciles a VisitorsApp object
type VisitorsAppReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=test1.example.com,resources=visitorsapps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=test1.example.com,resources=visitorsapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=test1.example.com,resources=visitorsapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the VisitorsApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *VisitorsAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("visitorsapp", req.NamespacedName)

	// your logic here
	visitor := &test1v1alpha1.VisitorsApp{}
	err := r.Get(ctx, req.NamespacedName, visitor)
	if err != nil {
		if errors.IsNotFound(err) {
			// object is not found , return do not requeue
			log.Info("visitorApp resource not found. Ignoring since object must be deleted")
			// Reconcile successful - don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		// Reconcile failed due to error - requeue
		log.Error(err, "Failed to get visitorApp")
		return ctrl.Result{}, err
	}

	var result *ctrl.Result

	// == MySQL ==========
	result, err = r.ensureSecret(req, visitor, r.mysqlAuthSecret(visitor))
	if result != nil {
		return *result, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *VisitorsAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&test1v1alpha1.VisitorsApp{}).
		Complete(r)
}
