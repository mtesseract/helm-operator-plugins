package reconciler

import (
	"context"

	"github.com/operator-framework/helm-operator-plugins/pkg/extension"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type extensions []extension.ReconcilerExtension

func (es extensions) forEach(f func(e extension.ReconcilerExtension) error) error {
	var err error
	for _, e := range es {
		err = f(e)
		if err != nil {
			return err
		}
	}
	return err
}

func (r *Reconciler) extBeginReconcile(ctx context.Context, reconciliationContext *extension.Context, obj *unstructured.Unstructured) error {
	return r.extensions.forEach(func(ext extension.ReconcilerExtension) error {
		e, ok := ext.(extension.BeginReconciliationExtension)
		if !ok {
			return nil
		}
		return e.BeginReconcile(ctx, reconciliationContext, obj)
	})
}

func (r *Reconciler) extEndReconcile(ctx context.Context, reconciliationContext *extension.Context, obj *unstructured.Unstructured) error {
	return r.extensions.forEach(func(ext extension.ReconcilerExtension) error {
		e, ok := ext.(extension.EndReconciliationExtension)
		if !ok {
			return nil
		}

		return e.EndReconcile(ctx, reconciliationContext, obj)
	})
}
