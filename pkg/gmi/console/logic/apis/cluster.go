/*
 @Version : 1.0
 @Author  : wangxiaokang
 @Email   : 'xiaokang.w@gmicloud.ai'
 @Time    : 2025/06/03 12:03:10
 Desc     : report cluster info to console
*/

package apis

import (
	"github.com/gin-gonic/gin"
)

type ReportRequest struct {
	ClusterName   string `json:"cluster_name"`
	ClusterID     string `json:"cluster_id"`
	ClusterType   string `json:"cluster_type"`
	ClusterStatus string `json:"cluster_status"`
}

func (h *Handler) Report(c *gin.Context) {
	var req ReportRequest
	h.parse(c, &req)

	h.success(req)
}
