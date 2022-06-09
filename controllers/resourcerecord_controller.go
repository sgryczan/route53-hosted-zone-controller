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

package controllers

import (
	"context"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	route53v1 "github.com/sgryczan/r53-hz-controller/api/v1"
	r53util "github.com/sgryczan/r53-hz-controller/pkg/aws"
	"github.com/sgryczan/r53-hz-controller/pkg/common"
)

// ResourceRecordReconciler reconciles a ResourceRecord object
type ResourceRecordReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=route53.aws.czan.io,resources=resourcerecords,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=route53.aws.czan.io,resources=resourcerecords/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=route53.aws.czan.io,resources=resourcerecords/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ResourceRecord object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *ResourceRecordReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Fetch the ResourceRecord instance
	resourceRecord := &route53v1.ResourceRecord{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, resourceRecord)
	if err != nil {
		if kerrors.IsNotFound(err) {
			r.Log.Info("ResourceRecord object not found")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		r.Log.Info("failed to retrieve ResourceRecord", "err", err)
		return ctrl.Result{}, err
	}

	if resourceRecord.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(resourceRecord, common.FinalizerName) {
			controllerutil.AddFinalizer(resourceRecord, common.FinalizerName)
			err := r.Client.Update(context.TODO(), resourceRecord)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		return r.reconcileDelete(resourceRecord)
	}

	// Get the hosted zone ID
	zoneID, err := r53util.GetZoneIDByName(resourceRecord.Spec.HostedZone.Name)
	if err != nil {
		r.Log.Error(err, "failed to retreive hosted zone id", "name", resourceRecord.Name)
		return ctrl.Result{}, err
	}

	err = r53util.UpdateRecordSet(*zoneID, &resourceRecord.Spec.RecordSet)
	if err != nil {
		r.Log.Error(err, "failed to update record set", "name", resourceRecord.Name)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ResourceRecordReconciler) reconcileDelete(resourceRecord *route53v1.ResourceRecord) (reconcile.Result, error) {
	var err error
	ctx := context.Background()
	if resourceRecord.DeletionTimestamp.IsZero() {
		r.Log.Info("Deletion timestamp is not set", "name", resourceRecord.Name)
		return ctrl.Result{}, nil
	}

	// Get the hosted zone ID
	zoneID, err := r53util.GetZoneIDByName(resourceRecord.Spec.HostedZone.Name)
	if err != nil {
		r.Log.Error(err, "failed to retreive hosted zone id", "name", resourceRecord.Name)
		return ctrl.Result{}, err
	}

	// Delete the record
	err = r53util.DeleteRecordSet(*zoneID, &resourceRecord.Spec.RecordSet)
	if err != nil {
		return ctrl.Result{}, err
	}

	controllerutil.RemoveFinalizer(resourceRecord, common.FinalizerName)
	err = r.Client.Update(ctx, resourceRecord)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ResourceRecordReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&route53v1.ResourceRecord{}).
		Complete(r)
}
