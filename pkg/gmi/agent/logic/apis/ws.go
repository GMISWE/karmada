/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/05/29
 * @Desc    : gmi agent apis
 */

package apis

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/karmada-io/karmada/pkg/gmi/base/types"
	"github.com/piaobeizu/titan/service"
)

func (h *WSHandler) ReportNF(c *gin.Context, msg []byte) (any, error) {
	var message service.WSMessage[types.Node]
	if err := json.Unmarshal(msg, &message); err != nil {
		return nil, err
	}
	h.GMIAgent.SendNfdWSMsg(message.Data)
	return nil, nil
}
