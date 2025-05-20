/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2025/05/07 16:02:21
 Desc     :
*/

package app

import (
	"context"
	"flag"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/karmada-io/karmada/cmd/gmi-operator/app/options"
	"github.com/karmada-io/karmada/pkg/sharedcli/klogflag"
	"github.com/karmada-io/karmada/pkg/sharedcli/profileflag"
	"github.com/karmada-io/karmada/pkg/util/names"
	"github.com/karmada-io/karmada/pkg/version"
	"github.com/karmada-io/karmada/pkg/version/sharedcommand"

	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
)

func NewGmiInitCommand(ctx context.Context) *cobra.Command {
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Use:  names.KarmadaAgentComponentName,
		Long: `The karmada-gmi-init is a initialization tool of karmada to init the gmi inference engine environment.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			// validate options
			if errs := opts.Validate(); len(errs) != 0 {
				return errs.ToAggregate()
			}
			if err := run(ctx, opts); err != nil {
				return err
			}
			return nil
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	fss := cliflag.NamedFlagSets{}

	genericFlagSet := fss.FlagSet("generic")
	genericFlagSet.AddGoFlagSet(flag.CommandLine)
	// opts.AddFlags(genericFlagSet, controllers.ControllerNames())

	// Set klog flags
	logsFlagSet := fss.FlagSet("logs")
	klogflag.Add(logsFlagSet)

	cmd.AddCommand(sharedcommand.NewCmdVersion(names.KarmadaAgentComponentName))
	cmd.Flags().AddFlagSet(genericFlagSet)
	cmd.Flags().AddFlagSet(logsFlagSet)

	// cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	// sharedcli.SetUsageAndHelpFunc(cmd, fss, cols)

	return cmd
}

func run(ctx context.Context, opts *options.Options) error {
	klog.Infof("karmada-scheduler version: %s", version.Get())

	profileflag.ListenAndServe(opts.ProfileOpts)

	restConfig, err := clientcmd.BuildConfigFromFlags(opts.Master, opts.KubeConfig)
	if err != nil {
		return fmt.Errorf("error building kubeconfig: %s", err.Error())
	}
	restConfig.QPS, restConfig.Burst = opts.KubeAPIQPS, opts.KubeAPIBurst

	// dynamicClientSet := dynamic.NewForConfigOrDie(restConfig)
	// kubeClientSet := kubernetes.NewForConfigOrDie(restConfig)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		cancel()
	}()

	// start event watcher
	// watcher.Start(restConfig, kubeClientSet, dynamicClientSet)
	return nil
}
