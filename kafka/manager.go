package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Manager struct {
	Brokers []string
	Writers map[string]*kafka.Writer
	Readers map[string]*kafka.Reader
}

// createManager 创建Kafka管理器
func createManager(brokers []string, topic []string, groupID string) (*Manager, error) {
	writers := make(map[string]*kafka.Writer)
	readers := make(map[string]*kafka.Reader)

	// 判断每个broker节点是否可连接
	for _, broker := range brokers {
		conn, err := kafka.Dial("tcp", broker)
		if err != nil {
			zap.L().Error("kafka.Dial failed", zap.String("broker", broker), zap.Error(err))
			return nil, err
		}
		conn.Close()
	}

	// 创建kafka主题
	if err := createTopics(brokers, topic); err != nil {
		zap.L().Error("createTopics failed", zap.Error(err))
		return nil, err
	}

	// 初始化生产者和消费者
	for _, t := range topic {
		writers[t] = &kafka.Writer{
			Addr: kafka.TCP(brokers...),
			Topic: t,
			Balancer: &kafka.LeastBytes{},
		}

		readers[t] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    t,
			GroupID:  groupID,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		})
	}

	return &Manager{Brokers: brokers, Writers: writers, Readers: readers}, nil

}

// 启动kafka消费者
func (km *Manager) startConsumer(ctx context.Context, topic string) {
	if reader, ok := km.Readers[topic]; ok {
		for {
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				zap.L().Error("reader.ReadMessage failed", zap.String("topic", topic), zap.Error(err))
				break
			}
			zap.L().Info("Received message", zap.String("topic", topic), zap.ByteString("value", m.Value))
		}
	} else {
		zap.L().Error("topic not found", zap.String("topic", topic))
	}

}


