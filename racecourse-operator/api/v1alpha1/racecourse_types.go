// SPDX-FileCopyrightText: 2025 Daniel Sampliner <samplinerD@gmail.com>
//
// SPDX-License-Identifier: GLWTPL

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RacecourseSpec defines the desired state of Racecourse.
type RacecourseSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Replicas is the number of replicas in the deployment
	Replicas *int32 `json:"replicas,omitempty"`

	// IngressHost is the host to use in the Ingress resource
	IngressHost *string `json:"ingressHost,omitempty"`
}

// RacecourseStatus defines the observed state of Racecourse.
type RacecourseStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Conditions []metav1.Condition `json:"conditions,omitempty"`
	Replicas   int32              `json:"replicas"`
	Selector   string             `json:"selector"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.selector
// +kubebuilder:printcolumn:name="Replicas",type=integer,JSONPath=`.spec.replicas`
// +kubebuilder:printcolumn:name="IngressHost",type=integer,JSONPath=`.spec.ingressPath`

// Racecourse is the Schema for the racecourses API.
type Racecourse struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RacecourseSpec   `json:"spec,omitempty"`
	Status RacecourseStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RacecourseList contains a list of Racecourse.
type RacecourseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Racecourse `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Racecourse{}, &RacecourseList{})
}
