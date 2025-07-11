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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	context "context"
	time "time"

	apistopov1alpha1 "github.com/karmada-io/karmada/pkg/apis/topo/v1alpha1"
	versioned "github.com/karmada-io/karmada/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/karmada-io/karmada/pkg/generated/informers/externalversions/internalinterfaces"
	topov1alpha1 "github.com/karmada-io/karmada/pkg/generated/listers/topo/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// CloudInformer provides access to a shared informer and lister for
// Clouds.
type CloudInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() topov1alpha1.CloudLister
}

type cloudInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewCloudInformer constructs a new informer for Cloud type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewCloudInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredCloudInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredCloudInformer constructs a new informer for Cloud type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredCloudInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TopoV1alpha1().Clouds().List(context.Background(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TopoV1alpha1().Clouds().Watch(context.Background(), options)
			},
			ListWithContextFunc: func(ctx context.Context, options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TopoV1alpha1().Clouds().List(ctx, options)
			},
			WatchFuncWithContext: func(ctx context.Context, options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.TopoV1alpha1().Clouds().Watch(ctx, options)
			},
		},
		&apistopov1alpha1.Cloud{},
		resyncPeriod,
		indexers,
	)
}

func (f *cloudInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredCloudInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *cloudInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apistopov1alpha1.Cloud{}, f.defaultInformer)
}

func (f *cloudInformer) Lister() topov1alpha1.CloudLister {
	return topov1alpha1.NewCloudLister(f.Informer().GetIndexer())
}
