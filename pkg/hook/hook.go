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
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type HookInput struct {
	Ctx                context.Context
	Log                logr.Logger
	Obj                *unstructured.Unstructured // current state of the resource object
	HelmRelease        *release.Release
	HelmValues         map[string]interface{} // for convenience: merged helmValues + helmConfig + overrideValues
	DeletionInProgress bool
}

type PreHook interface {
	Exec(in *HookInput) error
}

type PreHookFunc func(in *HookInput) error

func (f PreHookFunc) Exec(in *HookInput) error {
	return f(in)
}

type PostHook interface {
	Exec(in *HookInput) error
}

type PostHookFunc func(in *HookInput) error

func (f PostHookFunc) Exec(in *HookInput) error {
	return f(in)
}
