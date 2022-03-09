package extension

import (
	"context"

	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ReconcilerExtension interface {
	PreReconciliationExtension
	PostReconciliationExtension
}

type PreReconciliationExtension interface {
	PreReconcile(ctx context.Context, obj *unstructured.Unstructured) error
}

type PostReconciliationExtension interface {
	PostReconcile(ctx context.Context, obj *unstructured.Unstructured, release release.Release, vals chartutil.Values) error
}

type NoOpReconcilerExtension struct{}

func (e NoOpReconcilerExtension) PreReconcile(ctx context.Context, obj *unstructured.Unstructured) error {
	return nil
}

func (e NoOpReconcilerExtension) PostReconcile(ctx context.Context, obj *unstructured.Unstructured, release release.Release, vals chartutil.Values) error {
	return nil
}
