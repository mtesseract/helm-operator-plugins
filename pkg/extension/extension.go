package extension

import (
	"context"

	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
)

type ReconcilerExtension interface {
	PreReconciliationExtension
	PostReconciliationExtension
}

type PreReconciliationExtension interface {
	PreReconcile(ctx context.Context, reconciliationContext *Context, obj *unstructured.Unstructured) error
}

type PostReconciliationExtension interface {
	PostReconcile(ctx context.Context, reconciliationContext *Context, obj *unstructured.Unstructured) error
}

type NoOpReconcilerExtension struct{}

type Context struct {
	KubernetesConfig *rest.Config
	HelmRelease      *release.Release
	HelmValues       chartutil.Values
}

func (c *Context) GetHelmRelease() release.Release {
	if c == nil || c.HelmRelease == nil {
		return release.Release{}
	}
	return *c.HelmRelease
}

func (c *Context) GetHelmValues() chartutil.Values {
	if c == nil {
		return nil
	}
	return c.HelmValues
}

func (c *Context) GetKubernetesConfig() *rest.Config {
	if c == nil {
		return nil
	}
	return c.KubernetesConfig
}

func (e NoOpReconcilerExtension) PreReconcile(ctx context.Context, reconciliationContext *Context, obj *unstructured.Unstructured) error {
	return nil
}

func (e NoOpReconcilerExtension) PostReconcile(ctx context.Context, reconciliationContext *Context, obj *unstructured.Unstructured) error {
	return nil
}
