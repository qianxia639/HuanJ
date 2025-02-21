package handler

import (
	"Ice/internal/logs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func (handler *Handler) wsHandler(ctx *gin.Context) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	logs.Info("upgrader success")

	defer conn.Close()

	for {
		var msg struct {
			TargetID int64  `json:"receiver_id"` // 用户ID或群组ID
			Content  string `json:"content"`
		}

		if err := conn.ReadJSON(&msg); err != nil {
			logs.Error("消息读取失败:", err)
			break
		}

		logs.Info("message: ", msg)

		if err := conn.WriteJSON(gin.H{
			"from":    msg.TargetID,
			"content": msg.Content,
		}); err != nil {
			logs.Error("实时消息推送失败: ", err)
		}
	}
}
