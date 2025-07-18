/*
Copyright The Karmada Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	context "context"

	storagev1alpha1 "github.com/karmada-io/karmada/pkg/apis/storage/v1alpha1"
	scheme "github.com/karmada-io/karmada/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// JuicefsesGetter has a method to return a JuicefsInterface.
// A group's client should implement this interface.
type JuicefsesGetter interface {
	Juicefses() JuicefsInterface
}

// JuicefsInterface has methods to work with Juicefs resources.
type JuicefsInterface interface {
	Create(ctx context.Context, juicefs *storagev1alpha1.Juicefs, opts v1.CreateOptions) (*storagev1alpha1.Juicefs, error)
	Update(ctx context.Context, juicefs *storagev1alpha1.Juicefs, opts v1.UpdateOptions) (*storagev1alpha1.Juicefs, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, juicefs *storagev1alpha1.Juicefs, opts v1.UpdateOptions) (*storagev1alpha1.Juicefs, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*storagev1alpha1.Juicefs, error)
	List(ctx context.Context, opts v1.ListOptions) (*storagev1alpha1.JuicefsList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *storagev1alpha1.Juicefs, err error)
	JuicefsExpansion
}

// juicefses implements JuicefsInterface
type juicefses struct {
	*gentype.ClientWithList[*storagev1alpha1.Juicefs, *storagev1alpha1.JuicefsList]
}

// newJuicefses returns a Juicefses
func newJuicefses(c *StorageV1alpha1Client) *juicefses {
	return &juicefses{
		gentype.NewClientWithList[*storagev1alpha1.Juicefs, *storagev1alpha1.JuicefsList](
			"juicefses",
			c.RESTClient(),
			scheme.ParameterCodec,
			"",
			func() *storagev1alpha1.Juicefs { return &storagev1alpha1.Juicefs{} },
			func() *storagev1alpha1.JuicefsList { return &storagev1alpha1.JuicefsList{} },
		),
	}
}
