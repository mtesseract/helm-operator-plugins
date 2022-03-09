package reconciler

import (
	"context"

	"github.com/operator-framework/helm-operator-plugins/pkg/extension"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type extensions struct {
	extensions []extension.Extension
}

func (es *extensions) register(e extension.Extension) {
	es.extensions = append(es.extensions, e)
}

func (es *extensions) len() int {
	return len(es.extensions)
}

func (es *extensions) get(idx int) extension.Extension {
	return es.extensions[idx]
}

func (es *extensions) forEach(f func(e extension.Extension) error) error {
	var err error
	for _, e := range es.extensions {
		err = f(e)
		if err != nil {
			return err
		}
	}
	return err
}

func (r *Reconciler) extPreReconcile(ctx context.Context, obj *unstructured.Unstructured) error {
	return r.extensions.forEach(func(ext extension.Extension) error {
		e, ok := ext.(extension.PreReconciliationExtension)
		if !ok {
			return nil
		}
		return e.PreReconcile(ctx, obj)
	})
}

func (r *Reconciler) extPreUninstall(ctx context.Context, obj *unstructured.Unstructured) error {
	return r.extensions.forEach(func(ext extension.Extension) error {
		e, ok := ext.(extension.PreUninstallExtension)
		if !ok {
			return nil
		}
		return e.PreUninstall(ctx, obj)
	})
}

func (r *Reconciler) extPostUninstall(ctx context.Context, obj *unstructured.Unstructured) error {
	return r.extensions.forEach(func(ext extension.Extension) error {
		e, ok := ext.(extension.PostUninstallExtension)
		if !ok {
			return nil
		}
		return e.PostUninstall(ctx, obj)
	})
}

func (r *Reconciler) extPostReconcile(ctx context.Context, obj *unstructured.Unstructured, rel release.Release, vals chartutil.Values) error {
	return r.extensions.forEach(func(ext extension.Extension) error {
		e, ok := ext.(extension.PostReconciliationExtension)
		if !ok {
			return nil
		}

		return e.PostReconcile(ctx, obj, rel, vals)
	})
}
