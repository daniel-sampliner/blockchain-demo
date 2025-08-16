// SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
//
// SPDX-License-Identifier: GLWTPL

package controller

import (
	"context"
	"fmt"
	"maps"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/util/retry"
	"k8s.io/utils/ptr"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	webappv1alpha1 "github.com/daniel-sampliner/blockchain-demo/racecourse-operator/api/v1alpha1"
)

const (
	// typeAvailableRacecourse represents the status of the Deployment reconciliation
	typeAvailableRacecourse = "Available"
)

var commonLabels = map[string]string{"app.kubernetes.io/name": "racecourse"}

// RacecourseReconciler reconciles a Racecourse object
type RacecourseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.my.domain,resources=racecourses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.my.domain,resources=racecourses/finalizers,verbs=update
// +kubebuilder:rbac:groups=webapp.my.domain,resources=racecourses/scale,verbs=get;update;patch
// +kubebuilder:rbac:groups=webapp.my.domain,resources=racecourses/status,verbs=get;update;patch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Racecourse object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.20.2/pkg/reconcile
func (r *RacecourseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	racecourse := &webappv1alpha1.Racecourse{}
	err := r.Get(ctx, req.NamespacedName, racecourse)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Racecourse resource not found")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get Racecourse")
		return ctrl.Result{}, err
	}

	if err = r.reconcileDeployment(ctx, req, racecourse); err != nil {
		log.Error(err, "Failed to reconcile Racecourse Deployment")
		return ctrl.Result{}, err
	}

	if err = r.reconcileService(ctx, racecourse); err != nil {
		log.Error(err, "Failed to reconcile Racecourse Service")
		return ctrl.Result{}, err
	}

	if err = r.reconcileIngress(ctx, racecourse); err != nil {
		log.Error(err, "Failed to reconcile Racecourse Ingress")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Minute}, nil
}

// reconcileDeployment reconciles the Deployment resource managed by Racecourse
func (r *RacecourseReconciler) reconcileDeployment(ctx context.Context, req ctrl.Request, racecourse *webappv1alpha1.Racecourse) error {
	log := logf.FromContext(ctx)

	var err error

	// Set status as Unknown when no status is available
	if len(racecourse.Status.Conditions) == 0 {
		meta.SetStatusCondition(&racecourse.Status.Conditions, metav1.Condition{
			Type:    typeAvailableRacecourse,
			Status:  metav1.ConditionUnknown,
			Reason:  "Reconciling",
			Message: "Starting reconciliation",
		})
		if err = r.Status().Update(ctx, racecourse); err != nil {
			log.Error(err, "Failed to update Racecourse status")
			return err
		}

		if err := r.Get(ctx, req.NamespacedName, racecourse); err != nil {
			log.Error(err, "Failed to re-fetch Racecourse")
			return err
		}
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      racecourse.Name,
			Namespace: racecourse.Namespace,
		},
	}
	var op controllerutil.OperationResult
	if err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		op, err = ctrl.CreateOrUpdate(ctx, r.Client, deployment, func() error {
			if err := ctrl.SetControllerReference(racecourse, deployment, r.Scheme); err != nil {
				log.Error(err, "Failed to set Deployment controller reference")
				return err
			}

			image := "localhost/racecourse:latest"

			var desiredReplicas int32 = 0
			if racecourse.Spec.Replicas != nil {
				desiredReplicas = *racecourse.Spec.Replicas
			}

			if deployment.ObjectMeta.Labels == nil {
				deployment.ObjectMeta.Labels = make(map[string]string, len(commonLabels))
			}
			maps.Copy(deployment.Labels, commonLabels)

			deployment.Spec.Replicas = ptr.To(desiredReplicas)
			deployment.Spec.Selector = &metav1.LabelSelector{MatchLabels: commonLabels}

			deployment.Spec.Template = corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: commonLabels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           image,
						Name:            "racecourse",
						ImagePullPolicy: corev1.PullNever,
						Ports: []corev1.ContainerPort{{
							ContainerPort: 3000,
							Name:          "webapp",
						}},
					}},
				},
			}

			return nil
		})
		return err
	}); err != nil {
		log.Error(err, "Failed to CreateOrUpdate Deployment")
		meta.SetStatusCondition(&racecourse.Status.Conditions, metav1.Condition{
			Type:   typeAvailableRacecourse,
			Status: metav1.ConditionFalse,
			Reason: "Reconciling",
			Message: fmt.Sprintf(
				"Failed to %s Deployment for the custom resource (%s): (%s)",
				op,
				racecourse.Name,
				err,
			),
		})
		if err := r.Status().Update(ctx, racecourse); err != nil {
			log.Error(err, "Failed to update Racecourse status")
			return err
		}

		return err
	}

	if err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		racecourse := &webappv1alpha1.Racecourse{}
		if err := r.Get(ctx, req.NamespacedName, racecourse); err != nil {
			log.Error(err, "Failed to re-fetch Racecourse")
			return err
		}

		meta.SetStatusCondition(&racecourse.Status.Conditions, metav1.Condition{
			Type:   typeAvailableRacecourse,
			Status: metav1.ConditionTrue,
			Reason: "Reconciling",
			Message: fmt.Sprintf(
				"Deployment for custom resource (%s) with %d replicas %s successfully",
				racecourse.Name,
				*deployment.Spec.Replicas,
				op,
			),
		})
		racecourse.Status.Replicas = deployment.Status.Replicas
		racecourse.Status.Selector = metav1.FormatLabelSelector(deployment.Spec.Selector)
		return r.Status().Update(ctx, racecourse)
	}); err != nil {
		log.Error(err, "Failed to update Racecourse status")
		return err
	}

	return nil
}

// reconcileService reconciles the Service resource managed by Racecourse
func (r *RacecourseReconciler) reconcileService(ctx context.Context, racecourse *webappv1alpha1.Racecourse) error {
	log := logf.FromContext(ctx)

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      racecourse.Name,
			Namespace: racecourse.Namespace,
		},
	}

	var err error
	if err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err = ctrl.CreateOrUpdate(ctx, r.Client, service, func() error {
			if err := ctrl.SetControllerReference(racecourse, service, r.Scheme); err != nil {
				log.Error(err, "Failed to set Service controller reference")
				return err
			}

			if service.ObjectMeta.Labels == nil {
				service.ObjectMeta.Labels = make(map[string]string, len(commonLabels))
			}
			maps.Copy(service.Labels, commonLabels)

			service.Spec.Selector = commonLabels
			service.Spec.Ports = []corev1.ServicePort{{
				Name:       "webapp",
				Port:       3000,
				TargetPort: intstr.FromString("webapp"),
			}}

			return nil
		})
		return err
	}); err != nil {
		log.Error(err, "Failed to CreateOrUpdate Service")
		return err
	}

	return nil
}

// reconcileIngress reconciles the Service resource managed by Racecourse
func (r *RacecourseReconciler) reconcileIngress(ctx context.Context, racecourse *webappv1alpha1.Racecourse) error {
	log := logf.FromContext(ctx)

	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      racecourse.Name,
			Namespace: racecourse.Namespace,
		},
	}

	var err error
	if err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err = ctrl.CreateOrUpdate(ctx, r.Client, ingress, func() error {
			if err := ctrl.SetControllerReference(racecourse, ingress, r.Scheme); err != nil {
				log.Error(err, "Failed to set Ingress controller reference")
				return err
			}

			if ingress.ObjectMeta.Labels == nil {
				ingress.ObjectMeta.Labels = make(map[string]string, len(commonLabels))
			}
			maps.Copy(ingress.Labels, commonLabels)

			host := ""
			if racecourse.Spec.IngressHost != nil {
				host = *racecourse.Spec.IngressHost
			}

			ingress.Spec.Rules = []networkingv1.IngressRule{{
				Host: host,
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{{
							Path:     "/",
							PathType: ptr.To(networkingv1.PathTypePrefix),
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: "racecourse",
									Port: networkingv1.ServiceBackendPort{
										Name: "webapp",
									},
								},
							},
						}},
					},
				}},
			}

			return nil
		})
		return err
	}); err != nil {
		log.Error(err, "Failed to CreateOrUpdate Ingress")
		return err
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RacecourseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1alpha1.Racecourse{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&networkingv1.Ingress{}).
		Complete(r)
}
