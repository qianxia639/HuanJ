package handler

import (
	db "Ice/db/sqlc"
	"Ice/internal/logs"
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var clients = sync.Map{} // 存储在线用户连接

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

	// 存储连接
	clients.Store(handler.CurrentUserInfo.ID, conn)
	defer clients.Delete(handler.CurrentUserInfo.ID)

	for {

		var msg db.Message
		if err := conn.ReadJSON(&msg); err != nil {
			logs.Error("消息读取失败:", err)
			break
		}

		logs.Infof("message: %+v", msg)

		switch msg.SendType {
		case 1: // 私聊
		case 2: // 群聊
		default:
			logs.Error("未知的发送类型")
		}
	}
}

func (handler *Handler) privateChatMessage(ctx context.Context, msg db.Message) error {
	// 验证好友关系
	if exists, _ := handler.Store.ExistsFriendship(ctx, &db.ExistsFriendshipParams{UserID: msg.SenderID, FriendID: msg.ReceiverID}); !exists {
		return fmt.Errorf("非好友无法发送消息")
	}

	// 存储消息
	args := &db.CreateMessageParams{
		SessionID:    fmt.Sprintf("user:%d:%d", msg.SenderID, msg.ReceiverID),
		SenderID:     msg.SenderID,
		ReceiverID:   msg.ReceiverID,
		SendType:     msg.SendType,
		ReceiverType: msg.ReceiverType,
		MessageType:  msg.MessageType,
		Content:      msg.Content,
	}
	if err := handler.Store.CreateMessage(ctx, args); err != nil {
		return err
	}

	// 消息推送
	if conn, ok := clients.Load(msg.ReceiverID); ok {
		if err := conn.(*websocket.Conn).WriteJSON(msg); err != nil {
			return err
		}
	}
	return nil
}

func (handler *Handler) groupChatMessage(ctx context.Context, msg db.Message) error {
	// 校验群员身份
	if exists, _ := handler.Store.ExistsGroupMember(ctx, &db.ExistsGroupMemberParams{
		UserID:  msg.SenderID,
		GroupID: msg.ReceiverID,
	}); !exists {
		return fmt.Errorf("不在群组中")
	}

	// 存储消息
	args := &db.CreateMessageParams{
		SessionID:    fmt.Sprintf("group:%d:%d", msg.SenderID, msg.ReceiverID),
		SenderID:     msg.SenderID,
		ReceiverID:   msg.ReceiverID,
		SendType:     msg.SendType,
		ReceiverType: msg.ReceiverType,
		MessageType:  msg.MessageType,
		Content:      msg.Content,
	}
	if err := handler.Store.CreateMessage(ctx, args); err != nil {
		return err
	}

	// 获取群成员
	members, _ := handler.Store.GetGroupMemberList(ctx, msg.ReceiverID)

	// 消息推送
	for _, memberId := range members {
		// 可选: 是否发送给自己
		// if memberId == msg.SenderID {
		// 	continue
		// }
		if conn, ok := clients.Load(memberId); ok {
			if err := conn.(*websocket.Conn).WriteJSON(msg); err != nil {
				return err
			}
		}
	}

	return nil
}
