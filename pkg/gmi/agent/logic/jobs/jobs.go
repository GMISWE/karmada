/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/05/29
 * @Desc    : gmi jobs
 */

package jobs

import (
	"github.com/karmada-io/karmada/pkg/gmi/agent/pkg/core"
	"github.com/sirupsen/logrus"
)

func (j *Job) SyncResources() {
	resource := core.NewResource()

	topo := resource.CalTopo()

	logrus.Infof("sync resources to karmada: %p", &topo)
	// url := fmt.Sprintf("%s/api/v1/gpu", j.Args["karmada-control-plane-addr"].(string))
	// for _, node := range nodes {
	// json, _ := json.Marshal(node)
	// resp, err := http.Post(url, "application/json", bytes.NewBuffer(json))
	// if err != nil {
	// 	logrus.Errorf("push gpu info to karmada failed: %s", err)
	// }
	// defer resp.Body.Close()
	// }
}
