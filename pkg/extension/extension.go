package extension

import (
	"github.com/go-logr/logr"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Extension interface{}

type PreReconciliationExtension interface {
	ExecPreReconciliationExtension(obj *unstructured.Unstructured, vals chartutil.Values, log logr.Logger) error
}

type PostReconciliationExtension interface {
	ExecPostReconciliationExtension(obj *unstructured.Unstructured, rel release.Release, log logr.Logger) error
}
