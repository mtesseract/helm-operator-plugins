package extension

import (
	"context"

	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Extension interface{}

type PreReconciliationExtension interface {
	PreReconcile(ctx context.Context, obj *unstructured.Unstructured, release release.Release, vals chartutil.Values) error
}
type PreDeletionExtension interface {
	PreDelete(ctx context.Context, obj *unstructured.Unstructured, release release.Release, vals chartutil.Values) error
}

type PostReconciliationExtension interface {
	PostReconcile(ctx context.Context, obj *unstructured.Unstructured, release release.Release, vals chartutil.Values) error
}
