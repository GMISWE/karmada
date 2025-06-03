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

	"github.com/karmada-io/karmada/cmd/gmi/agent/app/options"
	"github.com/karmada-io/karmada/pkg/gmi/agent/pkg"
	"github.com/karmada-io/karmada/pkg/sharedcli/klogflag"
	"github.com/karmada-io/karmada/pkg/util/names"
	"github.com/karmada-io/karmada/pkg/version/sharedcommand"
	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
)

var engine *pkg.Engine

func NewGmiAgentCommand(ctx context.Context) *cobra.Command {
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Use:  names.KarmadaGMIAgentComponentName,
		Long: `The gmi agent for gmi cloud. It is responsible for gmi cloud agent.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			// if err := opts.Validate(); err != nil {
			// 	return err
			// }
			if err := run(ctx); err != nil {
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

	cmd.AddCommand(sharedcommand.NewCmdVersion(names.KarmadaGMIAgentComponentName))
	cmd.Flags().AddFlagSet(genericFlagSet)
	cmd.Flags().AddFlagSet(logsFlagSet)

	// cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	// sharedcli.SetUsageAndHelpFunc(cmd, fss, cols)
	return cmd
}

func run(ctx context.Context) error {
	// start the agent
	engine = pkg.NewEngine(ctx, names.KarmadaGMIAgentComponentName)
	engine.Start()
	return nil
}
