package handler

import (
	db "HuanJ/db/sqlc"
	"HuanJ/logs"
	"HuanJ/ws"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func (handler *Handler) wssHandler(connManager *ws.ConnectionManager, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		zap.L().Error("Websocket upgrade failed", zap.Error(err))
		return
	}

	zap.L().Info("Websocket upgrader success...")

	client := &ws.WsClient{
		UserId:      handler.CurrentUserInfo.ID,
		ConnManager: connManager,
		Conn:        conn,
		Send:        make(chan []byte, 256),
	}
	client.ConnManager.Register <- client

	defer func() {
		client.ConnManager.Unregister <- client
		conn.Close()
	}()

	go client.WritePump()
	go client.ReadPump()
}

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

	// 存储客户端连接信息到RedisClient
	userKey := fmt.Sprintf("ws_client:%d", handler.CurrentUserInfo.ID)
	if err := handler.RedisClient.SAdd(ctx, userKey, conn.RemoteAddr().String()).Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// 客户端断开时清理RedisClient中的连接信息
	defer func() {
		if err := handler.RedisClient.Del(ctx, userKey).Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}()

	// 处理消息
	for {

		var msg db.Message
		if err := conn.ReadJSON(&msg); err != nil {
			logs.Error("消息读取失败:", err)
			break
		}

		logs.Infof("message: %+v", msg)

		// switch msg.SendType {
		// case 1: // 私聊
		// case 2: // 群聊
		// default:
		// 	logs.Error("未知的发送类型")
		// }
	}
}

func (handler *Handler) wsHandlerV2(ctx *gin.Context) {

	// 用户认证
	auth := ctx.Request.Header.Get(authorizationHeader)
	payload, err := handler.Token.VerifyToken(auth)
	if err != nil {
		handler.Error(ctx, http.StatusUnauthorized, err.Error())
		zap.L().Error("Toekn校验失败", zap.Error(err))
		return
	}
	var loginUserInfo db.LoginUserInfo
	key := "user:" + payload.Username
	err = handler.RedisClient.Get(ctx, key).Scan(&loginUserInfo)
	if err != nil {
		zap.L().Error("缓存获取失败", zap.Error(err))
		handler.ServerError(ctx)
		return
	}

	// 升级 HTTP 连接为 Websocket
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		handler.Error(ctx, http.StatusBadRequest, "Websocket ungrade failed")
		zap.L().Error("Websocket upgrade failed",
			zap.Error(err),
			zap.String("remote", ctx.Request.RemoteAddr),
			zap.String("path", ctx.Request.URL.Path),
		)
		return
	}

	zap.L().Info("Websocket upgrader success...",
		zap.String("remote", ctx.Request.RemoteAddr),
		zap.String("username", loginUserInfo.Username),
	)

	client := &ws.WsClient{
		UserId: loginUserInfo.ID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}
	client.ConnManager.Register <- client

	defer func() {
		// 防止重复注销
		client.ConnManager.Unregister <- client
		conn.Close()
	}()

	go client.WritePump()
	go client.ReadPump()
}

func (handler *Handler) privateChatMessage(ctx context.Context, msg db.Message) error {
	// 验证好友关系
	if exists, _ := handler.Store.ExistsFriendship(ctx, &db.ExistsFriendshipParams{
		UserID:   msg.SenderID,
		FriendID: msg.ReceiverID,
	}); !exists {
		return fmt.Errorf("非好友无法发送消息")
	}

	// 存储消息
	// args := &db.CreateMessageParams{
	// 	SessionID:    fmt.Sprintf("user:%d:%d", msg.SenderID, msg.ReceiverID),
	// 	SenderID:     msg.SenderID,
	// 	ReceiverID:   msg.ReceiverID,
	// 	SendType:     msg.SendType,
	// 	ReceiverType: msg.ReceiverType,
	// 	MessageType:  msg.MessageType,
	// 	Content:      msg.Content,
	// }
	// if err := handler.Store.CreateMessage(ctx, args); err != nil {
	// 	return err
	// }

	// 消息推送
	// if conn, ok := clients.Load(msg.ReceiverID); ok {
	// 	if err := conn.(*websocket.Conn).WriteJSON(msg); err != nil {
	// 		return err
	// 	}
	// }
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
	// args := &db.CreateMessageParams{
	// 	SessionID:    fmt.Sprintf("group:%d:%d", msg.SenderID, msg.ReceiverID),
	// 	SenderID:     msg.SenderID,
	// 	ReceiverID:   msg.ReceiverID,
	// 	SendType:     msg.SendType,
	// 	ReceiverType: msg.ReceiverType,
	// 	MessageType:  msg.MessageType,
	// 	Content:      msg.Content,
	// }
	// if err := handler.Store.CreateMessage(ctx, args); err != nil {
	// 	return err
	// }

	// 获取群成员
	members, _ := handler.Store.GetGroupMemberList(ctx, msg.ReceiverID)
	_ = members
	// 消息推送
	// for _, memberId := range members {
	// 	// 可选: 是否发送给自己
	// 	if memberId == msg.SenderID {
	// 		continue
	// 	}
	// 	if conn, ok := clients.Load(memberId); ok {
	// 		if err := conn.(*websocket.Conn).WriteJSON(msg); err != nil {
	// 			return err
	// 		}
	// 	}
	// }

	return nil
}
