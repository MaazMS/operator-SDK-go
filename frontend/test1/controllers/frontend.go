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
	frontendPort        = 3000
	frontendServicePort = 30686
	frontendImage       = "jdob/visitors-webui:1.0.0"
)

func (r *VisitorsAppReconciler) frontendDeployment(visitor *test1v1alpha1.VisitorsApp) *appsv1.Deployment {

	labels := labels(visitor, "frontend")
	size := int32(1)
	// If the header was specified, add it as an env variable
	env := []corev1.EnvVar{}
	if visitor.Spec.Title != "" {
		env = append(env, corev1.EnvVar{
			Name:  "REACT_APP_TITLE",
			Value: visitor.Spec.Title,
		})
	}
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      visitor.Name + "-frontend",
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
						Image: frontendImage,
						Name:  "visitors-webui",
						Ports: []corev1.ContainerPort{{
							ContainerPort: frontendPort,
							Name:          "visitors",
						}},
						Env: env,
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(visitor, dep, r.Scheme)
	return dep
}

func (r *VisitorsAppReconciler) frontendService(visitor *test1v1alpha1.VisitorsApp) *corev1.Service {
	labels := labels(visitor, "frontend")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      visitor.Name + "-frontend-service",
			Namespace: visitor.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       frontendPort,
				TargetPort: intstr.FromInt(frontendPort),
				NodePort:   frontendServicePort,
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	r.Log.Info("Service Spec", "Service.Name", s.ObjectMeta.Name)
	controllerutil.SetControllerReference(visitor, s, r.Scheme)
	return s
}
func (r *VisitorsAppReconciler) updateFrontendStatus(visitor *test1v1alpha1.VisitorsApp) error {
	ctx := context.Background()
	visitor.Status.FrontendImage = frontendImage
	err := r.Status().Update(ctx, visitor)
	return err
}

func (r *VisitorsAppReconciler) handleFrontendChanges(visitor *test1v1alpha1.VisitorsApp) (*ctrl.Result, error) {
	ctx := context.Background()
	found := &appsv1.Deployment{}

	err := r.Get(ctx, types.NamespacedName{Name: visitor.Name + "-frontend", Namespace: visitor.Namespace}, found)

	if err != nil {
		return &ctrl.Result{RequeueAfter: 5 * time.Second}, err
	}
	title := visitor.Spec.Title
	existing := (*found).Spec.Template.Spec.Containers[0].Env[0].Value

	if title != existing {
		(*found).Spec.Template.Spec.Containers[0].Env[0].Value = title
		err = r.Update(ctx, found)

		if err != nil {

			r.Log.Error(err, "Failed to update Deployment.", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return &ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return &ctrl.Result{Requeue: true}, err
	}
	return nil, err
}
