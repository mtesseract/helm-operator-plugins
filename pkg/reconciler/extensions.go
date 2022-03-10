package reconciler

import (
	"context"

	"github.com/operator-framework/helm-operator-plugins/pkg/extension"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
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

func (r *Reconciler) extPreReconcile(ctx context.Context, obj *unstructured.Unstructured) error {
	return r.extensions.forEach(func(ext extension.ReconcilerExtension) error {
		e, ok := ext.(extension.PreReconciliationExtension)
		if !ok {
			return nil
		}
		return e.PreReconcile(ctx, obj)
	})
}

func (r *Reconciler) extPostReconcile(ctx context.Context, obj *unstructured.Unstructured, rel release.Release, vals chartutil.Values) error {
	return r.extensions.forEach(func(ext extension.ReconcilerExtension) error {
		e, ok := ext.(extension.PostReconciliationExtension)
		if !ok {
			return nil
		}

		return e.PostReconcile(ctx, obj, rel, vals)
	})
}
