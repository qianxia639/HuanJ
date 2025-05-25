package ws

import (
	"HuanJ/logs"
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

type RedisConnectionManager struct {
	// 注册的客户端 key: userId(int32), value:*WsClient
	Client *redis.Client
	// 客户端发送来的消息
	Broadcast chan []byte
	// 注册管道
	Register chan *WsClient
	// 注销管道
	Unregister chan *WsClient
	mu         sync.RWMutex
	ctx        context.Context
}

func NewRedisConnectionManager(redisClient *redis.Client) *RedisConnectionManager {
	return &RedisConnectionManager{
		Client:     redisClient,
		Broadcast:  make(chan []byte),
		Register:   make(chan *WsClient),
		Unregister: make(chan *WsClient),
		ctx:        context.Background(),
	}
}

func (m *RedisConnectionManager) Run() {
	logs.Info("---监听管道通信---")
	for {
		select {
		case client := <-m.Register:
			m.RegisterClient(client)
		case client := <-m.Unregister:
			m.UnregisterClient(client)
		case message := <-m.Broadcast:
			m.BroadcastMessage(message)
		}
	}
}

// 注册客户端
func (m *RedisConnectionManager) RegisterClient(client *WsClient) error {
	// 序列化客户端消息
	clientData, err := json.Marshal(client)
	if err != nil {
		return err
	}
	// 使用事务确保原子性
	pipe := m.Client.TxPipeline()
	pipe.HSet(m.ctx, "ws_clients", client.UserId, clientData)
	pipe.SAdd(m.ctx, "active_ws_clients", client.UserId)
	_, err = pipe.Exec(m.ctx)
	return err
}

// 注销客户端
func (m *RedisConnectionManager) UnregisterClient(client *WsClient) error {
	pipe := m.Client.TxPipeline()
	pipe.HDel(m.ctx, "ws_clients", fmt.Sprintf("%d", client.UserId))
	pipe.SRem(m.ctx, "active_ws_clients", client.UserId)
	_, err := pipe.Exec(m.ctx)
	return err
}

// 广播消息
func (m *RedisConnectionManager) BroadcastMessage(message []byte) error {
	// 获取所有活跃客户端的userId
	userIds, err := m.Client.SMembers(m.ctx, "active_ws_clients").Result()
	if err != nil {
		return err
	}

	// 遍历userId,发送消息
	for _, userId := range userIds {
		clientData, err := m.Client.HGet(m.ctx, "ws_clients", userId).Result()
		if err != nil {
			continue
		}
		var client WsClient
		if err := json.Unmarshal([]byte(clientData), &client); err != nil {
			continue
		}

		// 发送给消息

	}

	return nil
}
