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
	"errors"
	"fmt"
	"strings"

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

// HostedZoneReconciler reconciles a HostedZone object
type HostedZoneReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=route53.aws.czan.io,resources=hostedzones,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=route53.aws.czan.io,resources=hostedzones/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=route53.aws.czan.io,resources=hostedzones/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HostedZone object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *HostedZoneReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Fetch the HostedZone instance
	hostedZone := &route53v1.HostedZone{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, hostedZone)
	if err != nil {
		if kerrors.IsNotFound(err) {
			r.Log.Info("HostedZone object not found")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		r.Log.Info("failed to retrieve HostedZone", "err", err)
		return ctrl.Result{}, err
	}

	if hostedZone.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(hostedZone, common.FinalizerName) {
			controllerutil.AddFinalizer(hostedZone, common.FinalizerName)
			err := r.Client.Update(context.TODO(), hostedZone)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		return r.reconcileDelete(hostedZone)
	}

	// Check if the zone exists already
	zoneExists, err := r53util.HostedZoneExists(hostedZone.Name)

	if err != nil {
		r.Log.Error(err, "failed to check for existing zone", "name", hostedZone.Name)
		return ctrl.Result{}, err
	}

	if *zoneExists {
		r.Log.Info("Zone exists", "name", hostedZone.Name)
	} else {
		// If not, create the zone
		err = r53util.CreateHostedZone(hostedZone.Name)
		if err != nil {
			r.Log.Error(err, "failed to create hosted zone", "name", hostedZone.Name)
			hostedZone.Status.Error = err.Error()
			r.Client.Status().Update(context.Background(), hostedZone)

			return ctrl.Result{}, err
		}
	}

	// update zone details
	err = r.updateZoneDetail(ctx, hostedZone)
	if err != nil {
		r.Log.Error(err, "failed to update zone details", "name", hostedZone.Name)
		return ctrl.Result{}, err
	}

	// handle zone delegation if field is non-empty
	if hostedZone.Spec.DelegateOf != (route53v1.HostedZoneParent{}) {
		r.Log.Info("Handling delegation", "name", hostedZone.Name, "reason", ".spec.delegateOf field is non-nil", "value", fmt.Sprintf("%+v", hostedZone.Spec.DelegateOf))
		nameServers, err := r53util.GetNameServers(hostedZone.Name)
		if err != nil {
			return ctrl.Result{}, err
		}

		r.Log.Info("Create zone delegation", "zone", hostedZone.Name)
		err = r53util.CreateZoneDelegation(
			hostedZone.Name,
			nameServers,
			hostedZone.Spec.DelegateOf.ZoneID,
			hostedZone.Spec.DelegateOf.RoleARN,
		)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *HostedZoneReconciler) updateZoneDetail(ctx context.Context, hostedZone *route53v1.HostedZone) error {
	// update zone details
	details, err := r53util.GetZoneByName(hostedZone.Name)
	if err != nil {
		r.Log.Error(err, "failed to get details for zone", "name", hostedZone.Name)
		return err
	}

	zoneID := strings.Replace(*details.HostedZones[0].Id, "/hostedzone/", "", -1)

	hostedZone.Status.Details = route53v1.HostedZoneDetails{
		ID:             zoneID,
		PrivateZone:    details.HostedZones[0].Config.PrivateZone,
		RecordSetCount: *details.HostedZones[0].ResourceRecordSetCount,
	}

	hostedZone.Status.Ready = true
	hostedZone.Status.Error = ""

	err = r.Client.Status().Update(ctx, hostedZone)
	if err != nil {
		return err
	}

	return nil
}

func (r *HostedZoneReconciler) reconcileDelete(hostedZone *route53v1.HostedZone) (reconcile.Result, error) {
	var err error
	ctx := context.Background()
	if hostedZone.DeletionTimestamp.IsZero() {
		r.Log.Info("Deletion timestamp is not set", "name", hostedZone.Name)
		return ctrl.Result{}, nil
	}

	zoneID := hostedZone.Status.Details.ID

	if zoneID == "" {
		return ctrl.Result{}, errors.New("hostedZone resource has no ID")
	}

	// Delete the zone
	err = r53util.DeleteHostedZone(hostedZone.Status.Details.ID)
	if err != nil {
		return ctrl.Result{}, err
	}

	controllerutil.RemoveFinalizer(hostedZone, common.FinalizerName)
	err = r.Client.Update(ctx, hostedZone)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HostedZoneReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&route53v1.HostedZone{}).
		Complete(r)
}
