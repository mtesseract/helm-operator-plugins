package extension

import (
	"context"

	"github.com/go-logr/logr"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
)

type Extension interface{}

type Extensions struct {
	extensions []Extension
}

func (es *Extensions) Register(e Extension) {
	es.extensions = append(es.extensions, e)
}

func (es *Extensions) Len() int {
	return len(es.extensions)
}

func (es *Extensions) Get(idx int) Extension {
	return es.extensions[idx]
}

func (es *Extensions) Iterate(f func(e Extension) error) error {
	var err error
	for _, e := range es.extensions {
		err = f(e)
		if err != nil {
			return err
		}
	}
	return err
}

func (es *Extensions) LoggerInto(l logr.Logger) {
	for _, ext := range es.extensions {
		inject.LoggerInto(l, ext)
	}
}

type PreReconciliationExtension interface {
	ExecPreReconciliationExtension(ctx context.Context, obj *unstructured.Unstructured, release release.Release, vals chartutil.Values) error
}

func (es *Extensions) PreReconciliationExtPoint(ctx context.Context, obj *unstructured.Unstructured, release release.Release, vals chartutil.Values) error {
	return es.Iterate(func(ext Extension) error {
		e, ok := ext.(PreReconciliationExtension)
		if !ok {
			return nil
		}
		return e.ExecPreReconciliationExtension(ctx, obj, release, vals)
	})
}

type PreDeletionExtension interface {
	ExecPreDeletionExtension(ctx context.Context, obj *unstructured.Unstructured, vals chartutil.Values) error
}

func (es *Extensions) PreDeletionExtensionExtPoint(ctx context.Context, obj *unstructured.Unstructured, vals chartutil.Values) error {
	return es.Iterate(func(ext Extension) error {
		e, ok := ext.(PreDeletionExtension)
		if !ok {
			return nil
		}
		return e.ExecPreDeletionExtension(ctx, obj, vals)
	})
}

type PostReconciliationExtension interface {
	ExecPostReconciliationExtension(ctx context.Context, obj *unstructured.Unstructured, rel release.Release) error
}

func (es *Extensions) PostReconciliationExtPoint(ctx context.Context, obj *unstructured.Unstructured, rel release.Release) error {
	return es.Iterate(func(ext Extension) error {
		e, ok := ext.(PostReconciliationExtension)
		if !ok {
			return nil
		}
		return e.ExecPostReconciliationExtension(ctx, obj, rel)
	})
}
