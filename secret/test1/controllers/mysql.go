package controllers

import (
	test1v1alpha1 "github.com/MaazMS/operator-SDK-go/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// mysqlAuthSecret
func (r *VisitorsAppReconciler) mysqlAuthSecret(v *test1v1alpha1.VisitorsApp) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mysql-auth",
			Namespace: v.Namespace,
		},
		Type: "Opaque",
		StringData: map[string]string{
			"username": "visitors-user",
			"password": "visitors-pass",
		},
	}
	controllerutil.SetControllerReference(v, secret, r.Scheme)
	return secret
}
