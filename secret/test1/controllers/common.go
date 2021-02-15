package controllers

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	test1v1alpha1 "github.com/MaazMS/operator-SDK-go/api/v1alpha1"
)

/// ==================================
func (r *VisitorsAppReconciler) ensureSecret(req ctrl.Request,
	visitor *test1v1alpha1.VisitorsApp,
	s *corev1.Secret) (*ctrl.Result, error) {

	ctx := context.Background()

	found := &corev1.Secret{}
	err := r.Get(ctx, types.NamespacedName{Name: s.Name, Namespace: visitor.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Create the secret
		r.Log.Info("creating new secret", "secret.Namespace", s.Namespace, "secret.Name", s.Name)
		err = r.Create(ctx, s)
		if err != nil {

			r.Log.Error(err, "Failed to create new Secret", "Secret.Namespace", s.Namespace, "Secret.Name", s.Name)
			return &ctrl.Result{}, err
		} // Secret created successfully - return and requeue
		return &ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		r.Log.Error(err, "Failed to get Secret")
		return &ctrl.Result{}, err
	}
	return nil, err
}
