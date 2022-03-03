package reconciler

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/operator-framework/helm-operator-plugins/pkg/extension"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
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

func (es *extensions) iterate(f func(e extension.Extension) error) error {
	var err error
	for _, e := range es.extensions {
		err = f(e)
		if err != nil {
			return err
		}
	}
	return err
}

func (es *extensions) loggerInto(l logr.Logger) {
	for _, ext := range es.extensions {
		inject.LoggerInto(l, ext)
	}
}

func (r *Reconciler) extPreReconcile(ctx context.Context, obj *unstructured.Unstructured, release release.Release, vals chartutil.Values) error {
	return r.extensions.iterate(func(ext extension.Extension) error {
		e, ok := ext.(extension.PreReconciliationExtension)
		if !ok {
			return nil
		}
		return e.PreReconcile(ctx, obj, release, vals)
	})
}

func (r *Reconciler) extPreDelete(ctx context.Context, obj *unstructured.Unstructured, release release.Release, vals chartutil.Values) error {
	return r.extensions.iterate(func(ext extension.Extension) error {
		e, ok := ext.(extension.PreDeletionExtension)
		if !ok {
			return nil
		}
		return e.PreDelete(ctx, obj, release, vals)
	})
}

func (r *Reconciler) extPostReconcile(ctx context.Context, obj *unstructured.Unstructured, rel release.Release, vals chartutil.Values) error {
	return r.extensions.iterate(func(ext extension.Extension) error {
		e, ok := ext.(extension.PostReconciliationExtension)
		if !ok {
			return nil
		}
		return e.PostReconcile(ctx, obj, rel, vals)
	})
}
