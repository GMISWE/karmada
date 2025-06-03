/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/05/27
 * @Desc    : gmi agent for gmi cloud
 */

package app

import (
	"context"
	"fmt"

	"github.com/karmada-io/karmada/cmd/gmi/console/app/options"
	"github.com/karmada-io/karmada/pkg/gmi/console/pkg"
	"github.com/karmada-io/karmada/pkg/sharedcli/klogflag"
	"github.com/karmada-io/karmada/pkg/util/names"
	"github.com/karmada-io/karmada/pkg/version"
	"github.com/karmada-io/karmada/pkg/version/sharedcommand"
	"github.com/spf13/cobra"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog/v2"
)

var engine *pkg.Engine

func NewGmiConsoleCommand(ctx context.Context) *cobra.Command {
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Use:  names.KarmadaGmiConsoleComponentName,
		Long: `The gmi console for gmi cloud. It is responsible for gmi cloud console.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			// if err := opts.Validate(); err != nil {
			// 	return err
			// }
			if err := run(ctx, opts); err != nil {
				return err
			}
			return nil
		},
	}

	fss := cliflag.NamedFlagSets{}

	genericFlagSet := fss.FlagSet("generic")
	opts.AddFlags(genericFlagSet)

	// Set klog flags
	logsFlagSet := fss.FlagSet("logs")
	klogflag.Add(logsFlagSet)

	cmd.AddCommand(sharedcommand.NewCmdVersion(names.KarmadaGmiConsoleComponentName))
	cmd.Flags().AddFlagSet(genericFlagSet)
	cmd.Flags().AddFlagSet(logsFlagSet)

	// cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	// sharedcli.SetUsageAndHelpFunc(cmd, fss, cols)
	return cmd
}

func run(ctx context.Context, opts *options.Options) error {
	klog.Infof("karmada-gmi-storage version: %s", version.Get())

	restConfig, err := clientcmd.BuildConfigFromFlags(opts.Master, opts.KubeConfig)
	if err != nil {
		return fmt.Errorf("error building kubeconfig: %s", err.Error())
	}
	restConfig.QPS, restConfig.Burst = opts.KubeAPIQPS, opts.KubeAPIBurst

	dynamicClientSet := dynamic.NewForConfigOrDie(restConfig)
	kubeClientSet := kubernetes.NewForConfigOrDie(restConfig)

	// start the agent
	engine = pkg.NewEngine(ctx, names.KarmadaGmiConsoleComponentName, kubeClientSet, dynamicClientSet)
	engine.Start()

	return nil
}
