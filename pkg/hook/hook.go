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
	f   func(context.Context, *unstructured.Unstructured, chartutil.Values, logr.Logger) error
	log *logr.Logger
}

func NewPreHookFunc(f func(context.Context, *unstructured.Unstructured, chartutil.Values, logr.Logger) error) PreHookFunc {
	log := logr.Discard()
	return PreHookFunc{f: f, log: &log}
}

func NewPostHookFunc(f func(context.Context, *unstructured.Unstructured, release.Release, logr.Logger) error) PostHookFunc {
	log := logr.Discard()
	return PostHookFunc{f: f, log: &log}
}

type PreHook interface {
	extension.PreReconciliationExtension
}

type PostHook interface {
	extension.PostReconciliationExtension
}

func (h PreHookFunc) ExecPreReconciliationExtension(ctx context.Context, obj *unstructured.Unstructured, vals chartutil.Values) error {
	log := h.log
	if log == nil {
		sink := logr.Discard()
		log = &sink
	}
	return h.f(ctx, obj, vals, *log)
}

type PostHookFunc struct {
	f   func(context.Context, *unstructured.Unstructured, release.Release, logr.Logger) error
	log *logr.Logger
}

func (h PostHookFunc) ExecPostReconciliationExtension(ctx context.Context, obj *unstructured.Unstructured, rel release.Release) error {
	log := h.log
	if log == nil {
		sink := logr.Discard()
		log = &sink
	}
	return h.f(ctx, obj, rel, *log)
}
