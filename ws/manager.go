package ws

import (
	"HuanJ/logs"
	"sync"
	"time"

	"go.uber.org/zap"
)

// 维护活动的客户端，并向客户端广播消息
type ConnectionManager struct {
	// 所有连接
	Clients    sync.Map       // key: *WsClient  value: struct{}
	Broadcast  chan []byte    // 广播消息管道
	Register   chan *WsClient // 注册管道
	Unregister chan *WsClient // 注销管道
	closed     chan struct{}
	closeOnce  sync.Once
}

var Manager = &ConnectionManager{
	Clients:    sync.Map{},
	Broadcast:  make(chan []byte),
	Register:   make(chan *WsClient),
	Unregister: make(chan *WsClient),
	closed:     make(chan struct{}),
}

func (cm *ConnectionManager) Run() {
	logs.Info("---监听管道通信---")
	for {
		select {
		case client := <-cm.Register:
			cm.handlerRegister(client)
		case client := <-cm.Unregister:
			cm.handlerUnregister(client)
		case message := <-cm.Broadcast:
			cm.handlerBroadcast(message)
		}
	}
}

// 注册
func (cm *ConnectionManager) handlerRegister(client *WsClient) {
	if _, ok := cm.Clients.LoadOrStore(client, struct{}{}); !ok {
		zap.L().Info("Client registered", zap.Int32("user_id", client.UserId))
	}
}

// 注销
func (cm *ConnectionManager) handlerUnregister(client *WsClient) {
	if _, ok := cm.Clients.LoadAndDelete(client); ok {
		close(client.Send)
		zap.L().Info("Client unregistered", zap.Int32("user_id", client.UserId))
	}
}

// 广播
func (cm *ConnectionManager) handlerBroadcast(message []byte) {
	cm.Clients.Range(func(key, value interface{}) bool {
		client := key.(*WsClient)
		select {
		case client.Send <- message:
			// 发送成功
		default:
			// 发送失败时异步注销客户端
			go func(c *WsClient) {
				select {
				case cm.Unregister <- c:
				case <-time.After(100 * time.Millisecond):
					zap.L().Warn("unregister timeout", zap.Int32("user_id", c.UserId))
				}
			}(client)
		}
		return true
	})
}

// 关闭连接
func (cm *ConnectionManager) Close() {
	cm.closeOnce.Do(func() {
		close(cm.closed)
		close(cm.Broadcast)
		close(cm.Register)
		close(cm.Unregister)
	})
}

// 获取客户端用户列表
// func (cm *ConnectionManager) GetUserClients(userId int32) []*WsClient {
// 	// clients, ok := cm.Clients.Load(userId)
// 	// if !ok {
// 	// 	return nil
// 	// }
// 	// m := clients.(map[*WsClient]struct{})
// 	// result := make([]*WsClient, 0, len(m))
// 	// for client := range m {
// 	// 	result = append(result, client)
// 	// }
// 	result := make([]*WsClient, 0)
// 	cm.Clients.Range(func(k, v interface{}) bool {

// 		val, ok := cm.Clients.Load(k)
// 		if !ok {
// 			return false
// 		}

// 		result = append(result, val.(*WsClient))

// 		return true
// 	})

// 	return result
// }
