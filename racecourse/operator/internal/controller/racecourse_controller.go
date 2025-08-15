// SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
//
// SPDX-License-Identifier: GLWTPL

package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	webappv1alpha1 "github.com/daniel-sampliner/blockchain-demo/racecourse/operator/api/v1alpha1"
)

const (
	// typeAvailableRacecourse represents the status of the Deployment reconciliation
	typeAvailableRacecourse = "Available"
)

// RacecourseReconciler reconciles a Racecourse object
type RacecourseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=webapp.my.domain,resources=racecourses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.my.domain,resources=racecourses/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=webapp.my.domain,resources=racecourses/finalizers,verbs=update
// +kubebuilder:rbac:groups=race.wv,resources=racecourses/scale,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

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
			return ctrl.Result{}, err
		}

		if err := r.Get(ctx, req.NamespacedName, racecourse); err != nil {
			log.Error(err, "Failed to re-fetch Racecourse")
			return ctrl.Result{}, err
		}
	}

	found := &appsv1.Deployment{}
	err = r.Get(
		ctx,
		types.NamespacedName{
			Name:      racecourse.Name,
			Namespace: racecourse.Namespace,
		},
		found,
	)
	// Create new deployment if it does not exist
	if err != nil && apierrors.IsNotFound(err) {
		dep, err := r.deploymentForRacecourse(racecourse)
		if err != nil {
			log.Error(err, "Failed to define new Deployment resource for Racecourse")
			meta.SetStatusCondition(&racecourse.Status.Conditions, metav1.Condition{
				Type:   typeAvailableRacecourse,
				Status: metav1.ConditionFalse,
				Reason: "Reconciling",
				Message: fmt.Sprintf(
					"Failed to create Deployment for the custom resource (%s): (%s)",
					racecourse.Name,
					err,
				),
			})
			if err := r.Status().Update(ctx, racecourse); err != nil {
				log.Error(err, "Failed to update Racecourse status")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, err
		}

		log.Info(
			"Creating a new Deployment",
			"Deployment.Namespace",
			dep.Namespace,
			"Deployment.Name",
			dep.Name,
		)
		if err = r.Create(ctx, dep); err != nil {
			log.Error(
				err,
				"Failed to create a new Deployment",
				"Deployment.Namespace",
				dep.Namespace,
				"Deployment.Name",
				dep.Name,
			)
			return ctrl.Result{}, err
		}

		// Requeue after creating deployment
		return ctrl.Result{RequeueAfter: time.Second}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	var desiredReplicas int32 = 0
	if racecourse.Spec.Replicas != nil {
		desiredReplicas = *racecourse.Spec.Replicas
	}

	if found.Spec.Replicas == nil || *found.Spec.Replicas != desiredReplicas {
		found.Spec.Replicas = ptr.To(desiredReplicas)
		if err = r.Update(ctx, found); err != nil {
			log.Error(
				err,
				"Failed to update Deployment",
				"Deployment.Namespace",
				found.Namespace,
				"Deployment.Name",
				found.Name,
			)

			if err := r.Get(ctx, req.NamespacedName, racecourse); err != nil {
				log.Error(err, "Failed to re-fetch Racecourse")
				return ctrl.Result{}, err
			}

			meta.SetStatusCondition(&racecourse.Status.Conditions, metav1.Condition{
				Type:   typeAvailableRacecourse,
				Status: metav1.ConditionFalse,
				Reason: "Scaling",
				Message: fmt.Sprintf(
					"Failed to update the replicas for the custom resource (%s): (%s)",
					racecourse.Name,
					err,
				)},
			)
			if err := r.Status().Update(ctx, racecourse); err != nil {
				log.Error(err, "Failed to update Racecourse status")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, err
		}

		// Requeue after updating Deployment
		return ctrl.Result{RequeueAfter: time.Second}, nil
	}

	meta.SetStatusCondition(&racecourse.Status.Conditions, metav1.Condition{
		Type:   typeAvailableRacecourse,
		Status: metav1.ConditionTrue,
		Reason: "Reconciling",
		Message: fmt.Sprintf(
			"Deployment for custom resource (%s) with %d replicas created successfully",
			racecourse.Name,
			desiredReplicas,
		),
	})
	racecourse.Status.Replicas = found.Status.Replicas
	racecourse.Status.Selector = metav1.FormatLabelSelector(found.Spec.Selector)

	if err := r.Status().Update(ctx, racecourse); err != nil {
		log.Error(err, "Failed to update Racecourse status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// deploymentForRacecourse returns a Racecourse Deployment object
func (r *RacecourseReconciler) deploymentForRacecourse(racecourse *webappv1alpha1.Racecourse) (*appsv1.Deployment, error) {
	image := "localhost/racecourse:latest"
	labels := map[string]string{"app.kubernetes.io/name": "racecourse"}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      racecourse.Name,
			Namespace: racecourse.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: racecourse.Spec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
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
			},
		},
	}

	if err := ctrl.SetControllerReference(racecourse, dep, r.Scheme); err != nil {
		return nil, err
	}
	return dep, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RacecourseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1alpha1.Racecourse{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
