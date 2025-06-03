/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/05/29
 * @Desc    : nfd
 */

package nfd

import (
	"context"
	"time"

	"github.com/karmada-io/karmada/pkg/gmi/base/types"
	"github.com/karmada-io/karmada/pkg/util"
	"github.com/piaobeizu/titan/service"
	"github.com/piaobeizu/titan/utils/client"
	"github.com/sirupsen/logrus"
)

type NFDService struct {
	ctx              context.Context
	intervalDiscover time.Duration // nfd discover interval
	intervalCheck    time.Duration // nfd check interval
	wsClient         *client.WSClient
}

func NewNFDService(ctx context.Context, ws *client.WSClient) *NFDService {
	intervalDiscover := util.GetEnv("NFD_INTERVAL_DISCOVER", "200ms") // 200ms
	intervalCheck := util.GetEnv("NFD_INTERVAL_CHECK", "60s")         // 60s
	intervalDiscoverDuration, err := time.ParseDuration(intervalDiscover)
	if err != nil {
		logrus.Fatalf("Failed to parse NFD_INTERVAL_DISCOVER: %v", err)
	}
	intervalCheckDuration, err := time.ParseDuration(intervalCheck)
	if err != nil {
		logrus.Fatalf("Failed to parse NFD_INTERVAL_CHECK: %v", err)
	}
	return &NFDService{
		ctx:              ctx,
		intervalDiscover: intervalDiscoverDuration,
		intervalCheck:    intervalCheckDuration,
		wsClient:         ws,
	}
}

func (nfd *NFDService) DiscoverFeatures() {
	tickerDiscover := time.NewTicker(nfd.intervalDiscover)
	defer tickerDiscover.Stop()
	tickerCheck := time.NewTicker(nfd.intervalCheck)
	defer tickerCheck.Stop()

	printNo := 0
	for {
		select {
		case <-nfd.ctx.Done():
			return
		case <-tickerCheck.C:
			if err := nfd.checkNode(); err != nil {
				logrus.Errorf("Failed to check node: %v", err)
			}
		case <-tickerDiscover.C:
			node := nfd.discoverNode()
			if printNo%10 == 0 {
				logrus.Infof("send nfd ws msg: %p", &node)
			}
			nfd.wsClient.SendMessage(client.WSMSG_TYPE_JSON, service.WSMessage[types.Node]{
				Type:      types.WSMSG_TYPE_REPORT_NF,
				Timestamp: time.Now().UTC().Unix(),
				Data:      node,
			})
			printNo++
		}
	}
}

// checkNode check node status
func (nfd *NFDService) checkNode() error {
	// TODO: check gpu status

	return nil
}

func (nfd *NFDService) discoverNode() types.Node {
	hostDetector := NewHostDetector()
	gpuDetector := NewGPUDetector()

	host := hostDetector.DiscoverHost()
	gpus := gpuDetector.DiscoverGPUs()

	return types.Node{
		UUID: hostDetector.FetchNodeUUID(),
		Host: host,
		GPUs: gpus,
	}
}
