package ws

import (
	"Rejuv/logs"
	"sync"
)

// 维护活动的客户端，并向客户端广播消息
type ConnectionManager struct {
	// 注册的客户端 key: userId(int32), value:*WsClient
	Clients sync.Map
	// 客户端发送来的消息
	Broadcast chan []byte
	// 注册来自客户端的请求
	Register chan *WsClient
	// 注销来自客户端的请求
	Unregister chan *WsClient
	mu         sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		Broadcast:  make(chan []byte),
		Register:   make(chan *WsClient),
		Unregister: make(chan *WsClient),
	}
}

func (cm *ConnectionManager) Run() {
	logs.Info("---监听管道通信---")
	for {
		select {
		case client := <-cm.Register:
			// clients, _ := cm.Clients.LoadOrStore(client.UserId, make(map[*WsClient]struct{}))
			// cm.mu.Lock()
			// clients.(map[*WsClient]struct{})[client] = struct{}{}
			// cm.mu.Unlock()
			cm.Clients.Store(client.UserId, client)
		case client := <-cm.Unregister:
			if _, ok := cm.Clients.Load(client.UserId); ok {
				// cm.mu.Lock()
				// delete(clients.(map[*WsClient]struct{}), client)
				// cm.Clients.Delete(client.UserId)
				// close(client.Send)
				// cm.mu.Unlock()
				cm.Clients.Delete(client.UserId)
				close(client.Send)
			}
		case _ = <-cm.Broadcast:
			cm.Clients.Range(func(k, v interface{}) bool {
				// val, ok := cm.Clients.Load(k)
				// if !ok {
				// 	return false
				// }

				// client := val.(map[*WsClient]struct{})

				// select {
				// 	case client
				// }
				return true
			})
		}
	}
}

// 获取客户端用户列表
func (cm *ConnectionManager) GetUserClients(userId int32) []*WsClient {
	// clients, ok := cm.Clients.Load(userId)
	// if !ok {
	// 	return nil
	// }
	// m := clients.(map[*WsClient]struct{})
	// result := make([]*WsClient, 0, len(m))
	// for client := range m {
	// 	result = append(result, client)
	// }
	result := make([]*WsClient, 0)
	cm.Clients.Range(func(k, v interface{}) bool {

		val, ok := cm.Clients.Load(k)
		if !ok {
			return false
		}

		result = append(result, val.(*WsClient))

		return true
	})

	return result
}
