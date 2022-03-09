/*
Copyright 2020 The Operator-SDK Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hook

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/operator-framework/helm-operator-plugins/pkg/extension"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type PreHookFunc struct {
	f func(context.Context, *unstructured.Unstructured, logr.Logger) error
}

func NewPreHookFunc(f func(context.Context, *unstructured.Unstructured, logr.Logger) error) *PreHookFunc {
	wrappedF := func(ctx context.Context, obj *unstructured.Unstructured, log logr.Logger) error {
		err := f(ctx, obj, log)
		if err != nil {
			log.Error(err, "pre-release hook failed")
		}
		return nil
	}

	return &PreHookFunc{f: wrappedF}
}

func NewPostHookFunc(f func(context.Context, *unstructured.Unstructured, release.Release, chartutil.Values, logr.Logger) error) *PostHookFunc {
	wrappedF := func(ctx context.Context, obj *unstructured.Unstructured, rel release.Release, vals chartutil.Values, log logr.Logger) error {
		err := f(ctx, obj, rel, vals, log)
		if err != nil {
			log.Error(err, "post-release hook failed", "name", rel.Name, "version", rel.Version)
		}
		return nil
	}

	return &PostHookFunc{f: wrappedF}
}

type PreHook interface {
	extension.PreReconciliationExtension
}

type PostHook interface {
	extension.PostReconciliationExtension
}

func (h *PreHookFunc) PreReconcile(ctx context.Context, obj *unstructured.Unstructured) error {
	log := logr.FromContextOrDiscard(ctx)
	return h.f(ctx, obj, log)
}

var _ extension.PreReconciliationExtension = (*PreHookFunc)(nil)

type PostHookFunc struct {
	f func(context.Context, *unstructured.Unstructured, release.Release, chartutil.Values, logr.Logger) error
}

func (h *PostHookFunc) PostReconcile(ctx context.Context, obj *unstructured.Unstructured, rel release.Release, vals chartutil.Values) error {
	log := logr.FromContextOrDiscard(ctx)
	return h.f(ctx, obj, rel, vals, log)
}

var _ extension.PostReconciliationExtension = (*PostHookFunc)(nil)
