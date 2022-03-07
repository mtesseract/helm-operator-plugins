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
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
)

type PreHookFunc struct {
	f   func(context.Context, *unstructured.Unstructured, release.Release, chartutil.Values, logr.Logger) error
	log *logr.Logger
}

func NewPreHookFunc(f func(context.Context, *unstructured.Unstructured, release.Release, chartutil.Values, logr.Logger) error) *PreHookFunc {
	log := logr.Discard()
	wrappedF := func(ctx context.Context, obj *unstructured.Unstructured, rel release.Release, vals chartutil.Values, log logr.Logger) error {
		err := f(ctx, obj, rel, vals, log)
		if err != nil {
			log.Error(err, "pre-release hook failed")
		}
		return nil
	}

	return &PreHookFunc{f: wrappedF, log: &log}
}

func (h *PreHookFunc) InjectLogger(l logr.Logger) error {
	h.log = &l
	return nil
}

var _ inject.Logger = (*PreHookFunc)(nil)

func NewPostHookFunc(f func(context.Context, *unstructured.Unstructured, release.Release, chartutil.Values, logr.Logger) error) *PostHookFunc {
	log := logr.Discard()
	wrappedF := func(ctx context.Context, obj *unstructured.Unstructured, rel release.Release, vals chartutil.Values, log logr.Logger) error {
		err := f(ctx, obj, rel, vals, log)
		if err != nil {
			log.Error(err, "post-release hook failed", "name", rel.Name, "version", rel.Version)
		}
		return nil
	}

	return &PostHookFunc{f: wrappedF, log: &log}
}

func (h *PostHookFunc) InjectLogger(l logr.Logger) error {
	h.log = &l
	return nil
}

var _ inject.Logger = (*PostHookFunc)(nil)

type PreHook interface {
	extension.PreReconciliationExtension
}

type PostHook interface {
	extension.PostReconciliationExtension
}

func (h *PreHookFunc) PreReconcile(ctx context.Context, obj *unstructured.Unstructured, release release.Release, vals chartutil.Values) error {
	log := h.log
	if log == nil {
		sink := logr.Discard()
		log = &sink
	}
	return h.f(ctx, obj, release, vals, *log)
}

var _ extension.PreReconciliationExtension = (*PreHookFunc)(nil)

type PostHookFunc struct {
	f   func(context.Context, *unstructured.Unstructured, release.Release, chartutil.Values, logr.Logger) error
	log *logr.Logger
}

func (h *PostHookFunc) PostReconcile(ctx context.Context, obj *unstructured.Unstructured, rel release.Release, vals chartutil.Values) error {
	log := h.log
	if log == nil {
		sink := logr.Discard()
		log = &sink
	}
	return h.f(ctx, obj, rel, vals, *log)
}

var _ extension.PostReconciliationExtension = (*PostHookFunc)(nil)
