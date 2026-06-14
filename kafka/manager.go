package kafka

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Manager struct {
	Brokers []string
	Writers map[string]*kafka.Writer
	Readers map[string]*kafka.Reader
}

// Publish 向指定主题发送消息。
func (km *Manager) Publish(ctx context.Context, topic string, key, value []byte) error {
	writer, ok := km.Writers[topic]
	if !ok {
		return kafka.UnknownTopicOrPartition
	}
	return writer.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
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
			Addr:         kafka.TCP(brokers...),
			Topic:        t,
			Balancer:     &kafka.Hash{},
			MaxAttempts:  5,
			RequiredAcks: kafka.RequireAll,
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
			m, err := reader.FetchMessage(ctx)
			if err != nil {
				zap.L().Error("reader.ReadMessage failed", zap.String("topic", topic), zap.Error(err))
				break
			}
			if err := km.handleMessage(ctx, topic, m); err != nil {
				zap.L().Error("handle kafka message failed",
					zap.String("topic", topic),
					zap.ByteString("key", m.Key),
					zap.ByteString("value", m.Value),
					zap.Error(err),
				)
				continue
			}
			if err := reader.CommitMessages(ctx, m); err != nil {
				zap.L().Error("reader.CommitMessages failed", zap.String("topic", topic), zap.Error(err))
			}
		}
	} else {
		zap.L().Error("topic not found", zap.String("topic", topic))
	}

}

func (km *Manager) handleMessage(ctx context.Context, topic string, m kafka.Message) error {
	switch topic {
	case TopicLike:
		return handleLikeMessage(ctx, m)
	default:
		zap.L().Info("Received message", zap.String("topic", topic), zap.ByteString("value", m.Value))
		return nil
	}
}

func handleLikeMessage(ctx context.Context, m kafka.Message) error {
	event := new(models.PostLikeEvent)
	if err := json.Unmarshal(m.Value, event); err != nil {
		zap.L().Error("unmarshal like event failed", zap.ByteString("value", m.Value), zap.Error(err))
		return nil
	}
	return mysql.ApplyPostLikeEvent(ctx, event)
}

func (km *Manager) startFailedLikeEventRetry(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			km.retryFailedLikeEvents(ctx)
		}
	}
}

func (km *Manager) retryFailedLikeEvents(ctx context.Context) {
	events, err := mysql.ListFailedLikeEvents(ctx, 100)
	if err != nil {
		zap.L().Error("mysql.ListFailedLikeEvents failed", zap.Error(err))
		return
	}
	for _, event := range events {
		payload, err := json.Marshal(event)
		if err != nil {
			zap.L().Error("marshal failed like event failed", zap.Int64("eventID", event.EventID), zap.Error(err))
			continue
		}
		key := []byte(strconv.FormatInt(event.PostID, 10) + ":" + strconv.FormatInt(event.UserID, 10))
		if err := km.Publish(ctx, TopicLike, key, payload); err != nil {
			zap.L().Warn("retry publish like event failed", zap.Int64("eventID", event.EventID), zap.Error(err))
			continue
		}
		if err := mysql.DeleteFailedLikeEvent(ctx, event.EventID); err != nil {
			zap.L().Error("mysql.DeleteFailedLikeEvent failed", zap.Int64("eventID", event.EventID), zap.Error(err))
		}
	}
}
