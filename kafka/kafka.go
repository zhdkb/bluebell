package kafka

import (
	"context"

	"go.uber.org/zap"
)

var manager *Manager

// Init 初始化Kafka管理器
func Init(brokers []string) error {
	topics := []string{TopicLike}
	var err error
	manager, err = createManager(brokers, topics, "my-group")
	if err != nil {
		zap.L().Error("createManager failed")
		return err
	}
	zap.L().Info("createManager success", zap.Strings("brokers", brokers))

	// 启动消费者
	for _, topic := range topics {
		go manager.startConsumer(context.Background(), topic)
	}

	return nil
}


// GetManager 获取Kafka管理器
func GetManager() *Manager {
	if manager == nil {
		zap.L().Error("manager is nil")
		return nil
	}
	return manager
}


// Close 关闭Kafka管理器
func (km *Manager) Close() {
	for _, writer := range km.Writers {
		if err := writer.Close(); err != nil {
			zap.L().Error("writer.Close failed", zap.Error(err))
		}
	}

	for _, reader := range km.Readers {
		if err := reader.Close(); err != nil {
			zap.L().Error("reader.Close failed", zap.Error(err))
		}
	}
	zap.L().Info("Kafka manager closed")
}
