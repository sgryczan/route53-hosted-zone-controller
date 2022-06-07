/*
Copyright 2022.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HostedZoneSpec defines the desired state of HostedZone
type HostedZoneSpec struct {
	// AWSAccountID indicates the AWS Account in which the zone resides
	// +kubebuilder:validation:Optional
	// +nullable
	AWSAccountID string `json:"awsAccountID,omitempty"`

	// DelegateOf indicates if this hosted zone is a delegate of another hosted zone.
	// +kubebuilder:validation:Optional
	// +nullable
	DelegateOf HostedZoneParent `json:"delegateOf,omitempty"`

	// RecurseDelete. Indicates if all records in zone should be deleted when zone is deleted.
	// +kubebuilder:validation:Optional
	// +nullable
	RecurseDelete bool `json:"recurseDelete,omitempty"`
}

// HostedZoneParent represents a parent hosted zone in an AWS Account
type HostedZoneParent struct {
	AWSAccountID string `json:"awsAccountID"`
	ZoneID       string `json:"zoneID"`
	RoleARN      string `json:"roleARN"`
}

// HostedZoneStatus defines the observed state of HostedZone
type HostedZoneStatus struct {
	// +kubebuilder:validation:Optional
	// +nullable
	Details HostedZoneDetails `json:"details"`
	// +kubebuilder:validation:Optional
	// +nullable
	Ready bool `json:"ready"`
	// +kubebuilder:validation:Optional
	// +nullable
	Error string `json:"error"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName="hz"
//+kubebuilder:printcolumn:name="Id",type=string,JSONPath=`.status.details.hostedZoneID`

// HostedZone is the Schema for the hostedzones API
type HostedZone struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HostedZoneSpec   `json:"spec,omitempty"`
	Status HostedZoneStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HostedZoneList contains a list of HostedZone
type HostedZoneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HostedZone `json:"items"`
}

// HostedZoneDetails contains observed details about the hosted zone
type HostedZoneDetails struct {
	ID             string `json:"hostedZoneID,omitempty"`
	PrivateZone    bool   `json:"privateZone,omitempty"`
	RecordSetCount int64  `json:"recordSetCount,omitempty"`
}

func init() {
	SchemeBuilder.Register(&HostedZone{}, &HostedZoneList{})
}
