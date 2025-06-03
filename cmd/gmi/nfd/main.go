/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/05/26
 * @Desc    : node feature discovery for gmi cloud
 */

package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/google/uuid"
	"github.com/karmada-io/karmada/pkg/gmi/nfd"
	"github.com/karmada-io/karmada/pkg/util"
	"github.com/piaobeizu/titan/log"
	"github.com/piaobeizu/titan/utils"
	"github.com/piaobeizu/titan/utils/client"
	"github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 1024)
			n := runtime.Stack(buf, true)
			// print trace stack
			logrus.Errorf("panic error: %v\n%s", err, string(buf[:n]))
			os.Exit(1)
		}
	}()

	// add time sync
	timeClient, err := utils.NewTimeClient(utils.DefaultConfig())
	if err != nil {
		logrus.Errorf("Failed to create time client: %v", err)
		return
	}
	timeClient.StartDaemon()

	ctx, cancel := context.WithCancel(context.Background())

	log.InitLog("gmi-nfd", "debug")

	logrus.Info("starting gmi nfd")
	// start websocket client
	node_name := util.GetEnv("NODE_NAME", uuid.New().String())
	url := fmt.Sprintf("%s?client_id=%s", util.GetEnv("GMI_AGENT_URL", "ws://localhost:8080/ws"), node_name)

	// set websocket config
	config := client.DefaultConfig(url)
	ws := client.NewClient(ctx, node_name, config)
	if err := ws.Connect(); err != nil {
		logrus.Errorf("Failed to connect to nfd: %v", err)
		return
	}

	// start nfd service
	logrus.Info("start nfd service")
	nfd := nfd.NewNFDService(ctx, ws)
	go nfd.DiscoverFeatures()

	// wait for websocket client to be closed
	ws.Wait()
	cancel()
}
