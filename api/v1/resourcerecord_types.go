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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ResourceRecordSpec defines the desired state of ResourceRecord
type ResourceRecordSpec struct {
	RecordSet  ResourceRecordSet      `json:"recordSet"`
	HostedZone corev1.ObjectReference `json:"hostedZone,omitempty"`
}

type ResourceRecordSet struct {
	Name            string                `json:"name"`
	Type            string                `json:"type"`
	TTL             int64                 `json:"ttl"`
	ResourceRecords []ResourceRecordValue `json:"resourceRecords"`
}

type ResourceRecordValue struct {
	Value string `json:"value"`
}

// ResourceRecordStatus defines the observed state of ResourceRecord
type ResourceRecordStatus struct {
	Ready bool `json:"ready"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ResourceRecord is the Schema for the resourcerecords API
type ResourceRecord struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceRecordSpec   `json:"spec,omitempty"`
	Status ResourceRecordStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ResourceRecordList contains a list of ResourceRecord
type ResourceRecordList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceRecord `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceRecord{}, &ResourceRecordList{})
}
