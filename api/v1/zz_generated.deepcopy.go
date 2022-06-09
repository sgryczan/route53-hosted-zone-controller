//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostedZone) DeepCopyInto(out *HostedZone) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostedZone.
func (in *HostedZone) DeepCopy() *HostedZone {
	if in == nil {
		return nil
	}
	out := new(HostedZone)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HostedZone) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostedZoneDetails) DeepCopyInto(out *HostedZoneDetails) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostedZoneDetails.
func (in *HostedZoneDetails) DeepCopy() *HostedZoneDetails {
	if in == nil {
		return nil
	}
	out := new(HostedZoneDetails)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostedZoneList) DeepCopyInto(out *HostedZoneList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]HostedZone, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostedZoneList.
func (in *HostedZoneList) DeepCopy() *HostedZoneList {
	if in == nil {
		return nil
	}
	out := new(HostedZoneList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *HostedZoneList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostedZoneParent) DeepCopyInto(out *HostedZoneParent) {
	*out = *in
	out.HostedZoneRef = in.HostedZoneRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostedZoneParent.
func (in *HostedZoneParent) DeepCopy() *HostedZoneParent {
	if in == nil {
		return nil
	}
	out := new(HostedZoneParent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostedZoneSpec) DeepCopyInto(out *HostedZoneSpec) {
	*out = *in
	out.DelegateOf = in.DelegateOf
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostedZoneSpec.
func (in *HostedZoneSpec) DeepCopy() *HostedZoneSpec {
	if in == nil {
		return nil
	}
	out := new(HostedZoneSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HostedZoneStatus) DeepCopyInto(out *HostedZoneStatus) {
	*out = *in
	out.Details = in.Details
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HostedZoneStatus.
func (in *HostedZoneStatus) DeepCopy() *HostedZoneStatus {
	if in == nil {
		return nil
	}
	out := new(HostedZoneStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRecord) DeepCopyInto(out *ResourceRecord) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRecord.
func (in *ResourceRecord) DeepCopy() *ResourceRecord {
	if in == nil {
		return nil
	}
	out := new(ResourceRecord)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceRecord) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRecordList) DeepCopyInto(out *ResourceRecordList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ResourceRecord, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRecordList.
func (in *ResourceRecordList) DeepCopy() *ResourceRecordList {
	if in == nil {
		return nil
	}
	out := new(ResourceRecordList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceRecordList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRecordSet) DeepCopyInto(out *ResourceRecordSet) {
	*out = *in
	if in.ResourceRecords != nil {
		in, out := &in.ResourceRecords, &out.ResourceRecords
		*out = make([]ResourceRecordValue, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRecordSet.
func (in *ResourceRecordSet) DeepCopy() *ResourceRecordSet {
	if in == nil {
		return nil
	}
	out := new(ResourceRecordSet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRecordSpec) DeepCopyInto(out *ResourceRecordSpec) {
	*out = *in
	in.RecordSet.DeepCopyInto(&out.RecordSet)
	out.HostedZoneRef = in.HostedZoneRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRecordSpec.
func (in *ResourceRecordSpec) DeepCopy() *ResourceRecordSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceRecordSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRecordStatus) DeepCopyInto(out *ResourceRecordStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRecordStatus.
func (in *ResourceRecordStatus) DeepCopy() *ResourceRecordStatus {
	if in == nil {
		return nil
	}
	out := new(ResourceRecordStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRecordValue) DeepCopyInto(out *ResourceRecordValue) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRecordValue.
func (in *ResourceRecordValue) DeepCopy() *ResourceRecordValue {
	if in == nil {
		return nil
	}
	out := new(ResourceRecordValue)
	in.DeepCopyInto(out)
	return out
}
