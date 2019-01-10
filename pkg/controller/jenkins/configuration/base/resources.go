package base

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *ReconcileJenkinsBaseConfiguration) createResource(obj metav1.Object) error {
	runtimeObj, ok := obj.(runtime.Object)
	if !ok {
		return fmt.Errorf("is not a %T a runtime.Object", obj)
	}

	// Set Jenkins instance as the owner and controller
	if err := controllerutil.SetControllerReference(r.jenkins, obj, r.scheme); err != nil {
		return err
	}

	return r.client.Create(context.TODO(), runtimeObj)
}

func (r *ReconcileJenkinsBaseConfiguration) updateResource(obj metav1.Object) error {
	runtimeObj, ok := obj.(runtime.Object)
	if !ok {
		return fmt.Errorf("is not a %T a runtime.Object", obj)
	}

	// set Jenkins instance as the owner and controller, don't check error(can be already set)
	_ = controllerutil.SetControllerReference(r.jenkins, obj, r.scheme)

	return r.client.Update(context.TODO(), runtimeObj)
}

func (r *ReconcileJenkinsBaseConfiguration) createOrUpdateResource(obj metav1.Object) error {
	runtimeObj, ok := obj.(runtime.Object)
	if !ok {
		return fmt.Errorf("is not a %T a runtime.Object", obj)
	}

	// set Jenkins instance as the owner and controller, don't check error(can be already set)
	_ = controllerutil.SetControllerReference(r.jenkins, obj, r.scheme)

	err := r.client.Create(context.TODO(), runtimeObj)
	if err != nil && errors.IsAlreadyExists(err) {
		return r.updateResource(obj)
	} else if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}

	return nil
}
