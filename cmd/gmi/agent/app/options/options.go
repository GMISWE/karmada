/*
Copyright 2021 The Karmada Authors.

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

package options

import (
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	utilfeature "k8s.io/apiserver/pkg/util/feature"

	clusterscheme "github.com/karmada-io/karmada/pkg/apis/cluster/scheme"
	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	pkgfeatures "github.com/karmada-io/karmada/pkg/features"
	"github.com/karmada-io/karmada/pkg/sharedcli/profileflag"
)

const defaultEtcdPathPrefix = "/gmi-agent"

// Options contains everything necessary to create and run aggregated-apiserver.
type Options struct {
	Etcd           *genericoptions.EtcdOptions
	SecureServing  *genericoptions.SecureServingOptionsWithLoopback
	Authentication *genericoptions.DelegatingAuthenticationOptions
	Authorization  *genericoptions.DelegatingAuthorizationOptions
	Audit          *genericoptions.AuditOptions
	Features       *genericoptions.FeatureOptions
	CoreAPI        *genericoptions.CoreAPIOptions

	// KubeAPIQPS is the QPS to use while talking with karmada-apiserver.
	KubeAPIQPS float32
	// KubeAPIBurst is the burst to allow while talking with karmada-apiserver.
	KubeAPIBurst int

	ProfileOpts profileflag.Options
}

// NewOptions returns a new Options.
func NewOptions() *Options {
	o := &Options{
		Etcd:           genericoptions.NewEtcdOptions(storagebackend.NewDefaultConfig(defaultEtcdPathPrefix, clusterscheme.Codecs.LegacyCodec(schema.GroupVersion{Group: clusterv1alpha1.GroupVersion.Group, Version: clusterv1alpha1.GroupVersion.Version}))),
		SecureServing:  genericoptions.NewSecureServingOptions().WithLoopback(),
		Authentication: genericoptions.NewDelegatingAuthenticationOptions(),
		Authorization:  genericoptions.NewDelegatingAuthorizationOptions(),
		Audit:          genericoptions.NewAuditOptions(),
		Features:       genericoptions.NewFeatureOptions(),
		CoreAPI:        genericoptions.NewCoreAPIOptions(),
	}
	o.Etcd.StorageConfig.EncodeVersioner = runtime.NewMultiGroupVersioner(schema.GroupVersion{Group: clusterv1alpha1.GroupVersion.Group, Version: clusterv1alpha1.GroupVersion.Version}, schema.GroupKind{Group: clusterv1alpha1.GroupName})
	return o
}

// AddFlags adds flags to the specified FlagSet.
func (o *Options) AddFlags(flags *pflag.FlagSet) {
	o.Etcd.AddFlags(flags)
	o.SecureServing.AddFlags(flags)
	o.Authentication.AddFlags(flags)
	o.Authorization.AddFlags(flags)
	o.Audit.AddFlags(flags)
	o.Features.AddFlags(flags)
	o.CoreAPI.AddFlags(flags)

	flags.Lookup("kubeconfig").Usage = "Path to karmada control plane kubeconfig file."

	flags.Float32Var(&o.KubeAPIQPS, "kube-api-qps", 40.0, "QPS to use while talking with karmada-apiserver.")
	flags.IntVar(&o.KubeAPIBurst, "kube-api-burst", 60, "Burst to use while talking with karmada-apiserver.")
	_ = utilfeature.DefaultMutableFeatureGate.Add(pkgfeatures.DefaultFeatureGates)
	utilfeature.DefaultMutableFeatureGate.AddFlag(flags)
	o.ProfileOpts.AddFlags(flags)
}

// Complete fills in fields required to have valid data.
func (o *Options) Complete() error {
	return nil
}

// // Run runs the gmi-agent with options. This should never exit.
// func (o *Options) Run(ctx context.Context) error {
// 	klog.Infof("karmada-gmi-agent version: %s", version.Get())

// 	profileflag.ListenAndServe(o.ProfileOpts)

// 	config, err := o.Config()
// 	if err != nil {
// 		return err
// 	}

// 	restConfig := config.GenericConfig.ClientConfig
// 	restConfig.QPS, restConfig.Burst = o.KubeAPIQPS, o.KubeAPIBurst
// 	secretLister := config.GenericConfig.SharedInformerFactory.Core().V1().Secrets().Lister()
// 	config.GenericConfig.EffectiveVersion = utilversion.NewEffectiveVersion("1.0")

// 	server, err := config.Complete().New("gmi-agent", restConfig, secretLister)
// 	if err != nil {
// 		return err
// 	}

// 	server.GenericAPIServer.AddPostStartHookOrDie("start-gmi-agent-informers", func(context genericapiserver.PostStartHookContext) error {
// 		config.GenericConfig.SharedInformerFactory.Start(context.Done())
// 		return nil
// 	})

// 	return server.GenericAPIServer.PrepareRun().RunWithContext(ctx)
// }

// // Config returns config for the api server given Options
// func (o *Options) Config() (*customapiserver.Config, error) {
// 	// TODO have a "real" external address
// 	if err := o.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{netutils.ParseIPSloppy("127.0.0.1")}); err != nil {
// 		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
// 	}

// 	o.Features = &genericoptions.FeatureOptions{EnableProfiling: false}

// 	serverConfig := genericapiserver.NewRecommendedConfig(clusterscheme.Codecs)
// 	serverConfig.LongRunningFunc = customLongRunningRequestCheck(sets.NewString("watch", "proxy", "stream"),
// 		sets.NewString("attach", "exec", "proxy", "log", "portforward"))
// 	serverConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(generatedopenapi.GetOpenAPIDefinitions, openapi.NewDefinitionNamer(clusterscheme.Scheme))
// 	serverConfig.OpenAPIV3Config = genericapiserver.DefaultOpenAPIV3Config(generatedopenapi.GetOpenAPIDefinitions, openapi.NewDefinitionNamer(clusterscheme.Scheme))
// 	serverConfig.OpenAPIConfig.Info.Title = "Gmi-Agent"
// 	if err := o.applyTo(serverConfig); err != nil {
// 		return nil, err
// 	}

// 	config := &customapiserver.Config{
// 		GenericConfig: serverConfig,
// 		ExtraConfig:   customapiserver.ExtraConfig{},
// 	}
// 	return config, nil
// }

// func (o *Options) applyTo(config *genericapiserver.RecommendedConfig) error {
// 	if err := o.Etcd.ApplyTo(&config.Config); err != nil {
// 		return err
// 	}
// 	if err := o.SecureServing.ApplyTo(&config.Config.SecureServing, &config.Config.LoopbackClientConfig); err != nil {
// 		return err
// 	}
// 	if err := o.Authentication.ApplyTo(&config.Config.Authentication, config.SecureServing, config.OpenAPIConfig); err != nil {
// 		return err
// 	}
// 	if err := o.Authorization.ApplyTo(&config.Config.Authorization); err != nil {
// 		return err
// 	}
// 	if err := o.Audit.ApplyTo(&config.Config); err != nil {
// 		return err
// 	}
// 	if err := o.CoreAPI.ApplyTo(config); err != nil {
// 		return err
// 	}
// 	kubeClient, err := kubernetes.NewForConfig(config.ClientConfig)
// 	if err != nil {
// 		return err
// 	}
// 	if err = o.Features.ApplyTo(&config.Config, kubeClient, config.SharedInformerFactory); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // disable `deprecation` check until the underlying genericfilters.BasicLongRunningRequestCheck starts using generic Set.
// //
// //nolint:staticcheck
// func customLongRunningRequestCheck(longRunningVerbs, longRunningSubresources sets.String) apirequest.LongRunningRequestCheck {
// 	return func(r *http.Request, requestInfo *apirequest.RequestInfo) bool {
// 		reqClone := r.Clone(context.Background())
// 		p := reqClone.URL.Path
// 		currentParts := lifted.SplitPath(p)
// 		if isClusterProxy(currentParts) {
// 			currentParts = currentParts[6:]
// 			reqClone.URL.Path = "/" + strings.Join(currentParts, "/")
// 			requestInfo = lifted.NewRequestInfo(reqClone)
// 		}

// 		return genericfilters.BasicLongRunningRequestCheck(longRunningVerbs, longRunningSubresources)(r, requestInfo)
// 	}
// }

// func isClusterProxy(pathParts []string) bool {
// 	// cluster/proxy url path format: /apis/cluster.karmada.io/v1alpha1/clusters/{cluster}/proxy/...
// 	if len(pathParts) < 6 {
// 		return false
// 	}

// 	if pathParts[1] == "cluster.karmada.io" && pathParts[5] == "proxy" {
// 		return true
// 	}

// 	if pathParts[1] == "cluster.karmada.io" && len(pathParts) >= 7 &&
// 		pathParts[5] == "proxy" && pathParts[6] == "gpu" && pathParts[7] == "stream" {
// 		return true
// 	}

// 	return false
// }
