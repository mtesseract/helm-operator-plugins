package extension

import (
	"context"

	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Extension interface{}

type PreReconciliationExtension interface {
	PreReconcile(ctx context.Context, obj *unstructured.Unstructured) error
}

type PostReconciliationExtension interface {
	PostReconcile(ctx context.Context, obj *unstructured.Unstructured, release release.Release, vals chartutil.Values) error
}

type PreUninstallExtension interface {
	PreUninstall(ctx context.Context, obj *unstructured.Unstructured) error
}

type PostUninstallExtension interface {
	PostUninstall(ctx context.Context, obj *unstructured.Unstructured) error
}
