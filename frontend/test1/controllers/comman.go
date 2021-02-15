package controllers

import (
	"context"
	test1v1alpha1 "github.com/MaazMS/operator-sdk-go/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

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

func (r *VisitorsAppReconciler) ensureDeployment(req ctrl.Request,
	visitor *test1v1alpha1.VisitorsApp,
	dep *appsv1.Deployment) (*ctrl.Result, error) { //

	ctx := context.Background() //

	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: dep.Name, Namespace: visitor.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		r.Log.Info("creating new deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {

			r.Log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return &ctrl.Result{}, err
		} // Secret created successfully - return and requeue
		return &ctrl.Result{Requeue: true}, nil

	} else if err != nil {
		r.Log.Error(err, "Failed to get Deployment")
		return &ctrl.Result{}, err
	}
	return nil, err

}

func (r *VisitorsAppReconciler) ensureService(req ctrl.Request, visitor *test1v1alpha1.VisitorsApp, s *corev1.Service) (*ctrl.Result, error) {

	ctx := context.Background()

	found := &corev1.Service{}
	err := r.Get(ctx, types.NamespacedName{Name: s.Name, Namespace: visitor.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {

		// Create the secret
		r.Log.Info("Creating a new Service", "Service.Namespace", s.Namespace, "services.Name", s.Name)
		err = r.Create(ctx, s)
		if err != nil {

			r.Log.Error(err, "Failed to create new service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
			return &ctrl.Result{}, err
		} // Secret created successfully - return and requeue
		return &ctrl.Result{Requeue: true}, nil

	} else if err != nil {
		r.Log.Error(err, "failed to get Service")
		return &ctrl.Result{}, err
	}
	return nil, err
}
func labels(v *test1v1alpha1.VisitorsApp, tier string) map[string]string {
	return map[string]string{"app": "visitors", "visitorssite_cr": v.Name, "tier": tier}
}
