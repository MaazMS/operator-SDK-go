package controllers

import (
	"context"
	test1v1alpha1 "github.com/MaazMS/operator-sdk-go/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

const (
	backendPort        = 8000
	backendServicePort = 30685
	backendImage       = "jdob/visitors-service:1.0.0"
)

func (r *VisitorsAppReconciler) backendDeployment(visitor *test1v1alpha1.VisitorsApp) *appsv1.Deployment {
	labels := labels(visitor, "backend")
	size := visitor.Spec.Size

	userSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "mysql-auth"},
			Key:                  "username",
		},
	}

	passwordSecret := &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{Name: "mysql-auth"},
			Key:                  "password",
		},
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      visitor.Name + "-backend",
			Namespace: visitor.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           backendImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "visitors-service",
						Ports: []corev1.ContainerPort{{
							ContainerPort: backendPort,
							Name:          "visitors",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "MYSQL_DATABASE",
								Value: "visitors",
							},
							{
								Name:  "MYSQL_SERVICE_HOST",
								Value: "mysql-service",
							},
							{
								Name:      "MYSQL_USERNAME",
								ValueFrom: userSecret,
							},
							{
								Name:      "MYSQL_PASSWORD",
								ValueFrom: passwordSecret,
							},
						},
					}},
				},
			},
		},
	}
	// Set VisitorsApp instance as the owner and controller
	controllerutil.SetControllerReference(visitor, dep, r.Scheme)
	return dep
}
func (r *VisitorsAppReconciler) backendService(visitor *test1v1alpha1.VisitorsApp) *corev1.Service {
	labels := labels(visitor, "backend")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      visitor.Name + "-backend-service",
			Namespace: visitor.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       backendPort,
				TargetPort: intstr.FromInt(backendPort),
				NodePort:   30685,
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	r.Log.Info("Service Spec", "Service.Name", s.ObjectMeta.Name)

	controllerutil.SetControllerReference(visitor, s, r.Scheme)
	return s

}

func (r *VisitorsAppReconciler) updateBackendStatus(visitor *test1v1alpha1.VisitorsApp) error {
	ctx := context.Background()
	visitor.Status.BackendImage = backendImage
	err := r.Status().Update(ctx, visitor)
	return err
}

func (r *VisitorsAppReconciler) handleBackendChanges(visitor *test1v1alpha1.VisitorsApp) (*ctrl.Result, error) {

	ctx := context.Background()
	found := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{Name: visitor.Name + "-backend", Namespace: visitor.Namespace}, found)

	if err != nil {
		return &ctrl.Result{RequeueAfter: 5 * time.Second}, err
	}
	size := visitor.Spec.Size

	if size != *found.Spec.Replicas {

		found.Spec.Replicas = &size
		err = r.Update(ctx, found)

		if err != nil {
			r.Log.Error(err, "Failed to update Deployment.", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return &ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return &ctrl.Result{Requeue: true}, nil
	}
	return nil, err
}
