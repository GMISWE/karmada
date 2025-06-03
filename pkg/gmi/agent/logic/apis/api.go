/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/01/21 21:39:25
 Desc     :
*/

package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/karmada-io/karmada/pkg/gmi/agent/logic/config"
	"github.com/karmada-io/karmada/pkg/gmi/agent/pkg/core"
	"github.com/piaobeizu/titan/service"
)

type Handler struct {
	service.ApiHandler
}

type WSHandler struct {
	service.WSHandler[any]
	GMIAgent *core.GmiAgent
}

func (h *Handler) Demo(c *gin.Context) {
	h.success("hello world")
}

func (h *Handler) response(code int, message string, data any) service.Response {
	h.Response.Code, h.Response.Message, h.Response.Data = code, message, data
	return h.Response
}

func (h *Handler) success(data any) service.Response {
	h.Response.Code, h.Response.Message, h.Response.Data = 0, "success", data
	return h.Response
}

func (h *Handler) fail(data any) service.Response {
	h.Response.Code, h.Response.Message, h.Response.Data = 1, "failed", data
	return h.Response
}

func (h *WSHandler) response(conn *websocket.Conn, msg service.WSMessage[any]) {
	conn.WriteJSON(msg)
}

func (h *Handler) parse(c *gin.Context, data any) {
	if err := c.ShouldBindJSON(&data); err != nil {
		c.AbortWithStatusJSON(200, h.response(config.ERR_API_USER, "parse request body error", nil))
		panic(err)
	}
}
