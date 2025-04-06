package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const (
	TopicLike = "like" // 点赞主题
)

// createTopics 创建Kafka主题
// 该函数会检查每个broker节点是否可连接，并创建指定的主题
// 如果主题已存在，则不会重复创建
func createTopics(brokers []string, topics []string) error {
	client := kafka.Client{
		Addr: kafka.TCP(brokers...),
	}

	ctx := context.Background()

	var topicConfigs []kafka.TopicConfig
	for _, t := range topics {
		topicConfigs = append(topicConfigs, kafka.TopicConfig{
			Topic:             t,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
	}

	req := &kafka.CreateTopicsRequest{
		Topics: topicConfigs,
	}

	resp, err := client.CreateTopics(ctx, req)
	if err != nil {
		zap.L().Error("client.CreateTopics failed", zap.Error(err))
		return err
	}

	// 处理每个 topic 的错误（注意：Errors 是 map[string]error）
	for topicName, topicErr := range resp.Errors {
		if topicErr != nil {
			zap.L().Error("failed to create topic",
				zap.String("topic", topicName),
				zap.Error(topicErr),
			)
		} else {
			zap.L().Info("topic created", zap.String("topic", topicName))
		}
	}

	return nil
}
