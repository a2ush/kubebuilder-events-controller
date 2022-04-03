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
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cloudwatchlogs "github.com/a2ush/kubebuilder-events-controller/output"
)

// EventReconciler reconciles a Event object
type EventReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	CWClient *cloudwatchlogs.CloudWatchLogs
}

//+kubebuilder:rbac:groups=core,resources=events,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=events/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core,resources=events/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Event object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *EventReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	event := &corev1.Event{}
	err := r.Get(ctx, req.NamespacedName, event)
	if errors.IsNotFound(err) {
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Println(err, "unable to get Event", "name", req.NamespacedName)
		return ctrl.Result{}, err
	}

	log.Println("Send events to CloudWatch logs")
	r.CWClient.PutLogEvents(event)

	// /*
	// 	ObjectMeta.UID:
	// 	InvolvedObject.Namespace
	// 	InvolvedObject.Name
	// 	InvolvedObject.UID
	// 	Reason:
	// 	Message:
	// 	FirstTimestamp:
	// 	Source.Component:
	// 	Source.Host:
	// */
	// instance := event.DeepCopy()
	// fmt.Printf(
	// 	"new event :\n ObjectMeta.UID: %v\n InvolvedObject.Namespace: %v\n InvolvedObject.Name: %v\n InvolvedObject.UID: %v\n Reason: %v\n Message: %v\n FirstTimestamp: %v\n Source.Component: %v\n Source.Host: %v\n",
	// 	instance.ObjectMeta.UID,
	// 	instance.InvolvedObject.Namespace,
	// 	instance.InvolvedObject.Name,
	// 	instance.InvolvedObject.UID,
	// 	instance.Reason,
	// 	instance.Message,
	// 	instance.FirstTimestamp,
	// 	instance.Source.Component,
	// 	instance.Source.Host)

	return ctrl.Result{}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *EventReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Event{}).
		Complete(r)
}
