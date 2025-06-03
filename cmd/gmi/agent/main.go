/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/05/27
 * @Desc    : gmi agent for gmi cloud
 */

package main

import (
	"os"

	"github.com/karmada-io/karmada/cmd/gmi/agent/app"
	"k8s.io/component-base/cli"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

func main() {
	ctx := controllerruntime.SetupSignalHandler()
	cmd := app.NewGmiAgentCommand(ctx)
	code := cli.Run(cmd)
	os.Exit(code)
}
