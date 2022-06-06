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
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	route53v1 "github.com/sgryczan/r53-hz-controller/api/v1"
	r53util "github.com/sgryczan/r53-hz-controller/pkg/aws"
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

	// Check if the zone exists already
	zoneExists, err := r53util.HostedZoneExists(hostedZone.Name)

	if err != nil {
		r.Log.Error(err, "failed to check for existing zone", "name", hostedZone.Name)
		return ctrl.Result{}, err
	}

	if *zoneExists {
		r.Log.Info("Zone exists", "name", hostedZone.Name)
		return ctrl.Result{}, nil
	}

	// If not, create the zone
	err = r53util.CreateHostedZone(hostedZone.Name)
	if err != nil {
		r.Log.Error(err, "failed to create hosted zone", "name", hostedZone.Name)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HostedZoneReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&route53v1.HostedZone{}).
		Complete(r)
}