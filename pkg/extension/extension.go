package extension

import (
	"github.com/go-logr/logr"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

func (es *Extensions) PreReconciliationExtPoint(obj *unstructured.Unstructured, vals chartutil.Values, log logr.Logger) error {
	return es.Iterate(func(ext Extension) error {
		e, ok := ext.(PreReconciliationExtension)
		if !ok {
			return nil
		}
		return e.ExecPreReconciliationExtension(obj, vals, log)
	})
}

func (es *Extensions) PostReconciliationExtPoint(obj *unstructured.Unstructured, rel release.Release, log logr.Logger) error {
	return es.Iterate(func(ext Extension) error {
		e, ok := ext.(PostReconciliationExtension)
		if !ok {
			return nil
		}
		return e.ExecPostReconciliationExtension(obj, rel, log)
	})
}

type PreReconciliationExtension interface {
	ExecPreReconciliationExtension(obj *unstructured.Unstructured, vals chartutil.Values, log logr.Logger) error
}

type PostReconciliationExtension interface {
	ExecPostReconciliationExtension(obj *unstructured.Unstructured, rel release.Release, log logr.Logger) error
}
