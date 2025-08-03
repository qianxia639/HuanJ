package ws

import (
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// 允许写消息的时间
	writeWait = 10 * time.Second

	// 允许读取下一条pong消息的时间
	pongWait = 60 * time.Second

	// 向另一方发送ping. 必须小于 pongWait.
	pingPeriod = (pongWait * 9) / 10

	// 允许的消息大小
	maxMessageSize = 512
)

type WsClient struct {
	UserId      int32
	ConnManager *ConnectionManager
	Conn        *websocket.Conn
	// 发送消息的缓冲管道
	Send chan []byte
}

// ReadPump 将消息从websocket连接发送到连接管理
func (c *WsClient) ReadPump() {
	defer func() {
		c.ConnManager.Unregister <- c
		c.Conn.Close()
	}()

	// 设置读取限制
	c.Conn.SetReadLimit(maxMessageSize)
	// 设置读取超时
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zap.L().Error("websocket连接已关闭", zap.Error(err))
			}
			break
		}
		// 消息处理
		HandlerMessage(c, message)
	}
}

// WritePump 将消息从连接管理发送到websocket连接
func (c *WsClient) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// 连接管理关闭了管道
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				zap.L().Error("获取下一个消息失败", zap.Error(err))
				return
			}
			w.Write(message)

			/// 将排队的聊天消息添加到当前websocket消息中
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				zap.L().Error("写入ping消息失败", zap.Error(err))
				return
			}
		}
	}
}

func HandlerMessage(client *WsClient, message []byte) {
	// 解析消息格式
	// 根据消息类型处理私聊/群聊
	// 获取目标用户的客户端连接并发送消息
	// 群聊则查询群成员并遍历发送
	// client.ConnManager.Broadcast <- message
}

// 私聊
// func SendPrivateMessage(sender *WsClient, toId int32, message []byte) {
// 	clients := sender.ConnManager.GetUserClients(toId)
// 	for _, client := range clients {
// 		client.Send <- message
// 	}
// 	// 持久化消息
// }

// 群聊
func SendGroupMessage(sender *WsClient, groupId int32, message []byte) {
	// 获取群组成员，然后遍历
	// clients := sender.ConnManager.GetUserClients(toId)
	// for _, client := range clients {
	// 	client.Send <- message
	// }
	// 持久化消息
}
