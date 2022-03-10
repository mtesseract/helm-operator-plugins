package extension

import (
	"context"

	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
)

type ReconcilerExtension interface {
	BeginReconciliationExtension
	EndReconciliationExtension
}

type BeginReconciliationExtension interface {
	BeginReconcile(ctx context.Context, reconciliationContext *Context, obj *unstructured.Unstructured) error
}

type EndReconciliationExtension interface {
	EndReconcile(ctx context.Context, reconciliationContext *Context, obj *unstructured.Unstructured) error
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

func (e NoOpReconcilerExtension) BeginReconcile(ctx context.Context, reconciliationContext *Context, obj *unstructured.Unstructured) error {
	return nil
}

func (e NoOpReconcilerExtension) EndReconcile(ctx context.Context, reconciliationContext *Context, obj *unstructured.Unstructured) error {
	return nil
}
