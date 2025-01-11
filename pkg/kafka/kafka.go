package kafka

import (
	"context"
	"fmt"
	model "kafka-demo/pkg/user"
	"time"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"github.com/rs/zerolog/log"
)

func KafkaGo(name string) string {
    result := "KafkaGo " + name
    return result
}

var producer sarama.SyncProducer
var consumer sarama.ConsumerGroup

type KafkaGoConsumer = sarama.ConsumerGroup

type KafkaNotificationMsg struct {
	User *model.User
	NotificationStatus string
}

// Messages statuses in kafka
const (
	NEW_USER_WELCOMING_NOTIFICATIONS = "NEW_USER_WELCOMING_NOTIFICATIONS"
	Notifications = "notifications" 
)



// setupKafkaProducer initializes a Kafka producer
func SetupKafkaProducer(brokerList []string) (sarama.SyncProducer, error) {
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll
    config.Producer.Return.Successes = true

    p, err := sarama.NewSyncProducer(brokerList, config)
    if err != nil {
        return nil, err
    }
    producer = otelsarama.WrapSyncProducer(config, p)
    return p, nil
}
func ProduceMessage(topic string, key string, message string, headers map[string]string) error {
    log.Info().Msg(".............ProduceMessage.........")
    var kafkaHeaders []sarama.RecordHeader
    for k, v := range headers {
        kafkaHeaders = append(kafkaHeaders, sarama.RecordHeader{
            Key:   []byte(k),
            Value: []byte(v),
        })
    }
    msg := &sarama.ProducerMessage{
        Topic:     topic,
        Key:       sarama.StringEncoder(key),
        Value:     sarama.StringEncoder(message),
        Headers:   kafkaHeaders,
        Timestamp: time.Now(),
        Metadata:  "some metadata",
    }
    log.Info().Msgf(".............msg.........:%v", msg)
    log.Info().Msgf(".............kafkaHeaders.........:%v", kafkaHeaders)

    _, _, err := producer.SendMessage(msg)
    if err != nil {
        log.Error().Msgf(".............msg.........:%v", err.Error())

    }
    return err
}

// SetupKafkaConsumer initializes a Kafka consumer group
func SetupKafkaConsumer(brokerList []string, groupID string) (sarama.ConsumerGroup, error) {
    config := sarama.NewConfig()
    config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
    config.Consumer.Offsets.Initial = sarama.OffsetOldest

    c, err := sarama.NewConsumerGroup(brokerList, groupID, config)
    if err != nil {
        return nil, err
    }
    consumer = c
    return c, nil
}

// ConsumeMessages starts consuming messages from the specified Kafka topics
func ConsumeMessages(topics []string, handler sarama.ConsumerGroupHandler) error {
    ctx := context.Background()
    wrappedHandler := otelsarama.WrapConsumerGroupHandler(handler)
    for {
        err := consumer.Consume(ctx, topics, wrappedHandler)
        if err != nil {
            log.Error().Err(err).Msg("Error consuming messages")
            return err
        }
    }
}

type MessageHandler interface {
    CustomHandle(key string, msg string, headers map[string]string) error
}

// ConsumerGroupHandler is a custom implementation of sarama.ConsumerGroupHandler
type ConsumerGroupHandler struct {
    handler    MessageHandler
    autoCommit bool
}

// NewConsumerGroupHandler creates a new ConsumerGroupHandler
func NewConsumerGroupHandler(handler MessageHandler) *ConsumerGroupHandler {
    return &ConsumerGroupHandler{handler: handler}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
    return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
    return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (h *ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        // Print the message value
        fmt.Printf("Message claimed: value = %s, timestamp = %v, topic = %s\n", string(msg.Value), msg.Timestamp, msg.Topic)
        // Custom handle the message
        key := string(msg.Key)
        value := string(msg.Value)
        headersMap := make(map[string]string)
        for _, header := range msg.Headers {
            headersMap[string(header.Key)] = string(header.Value)
        }
        err := h.handler.CustomHandle(key, value, headersMap)
        if err != nil {
            log.Error().Err(err).Msg("Error in custom handle")
        }
        log.Info().Msg("Marking message as processed")
        sess.MarkMessage(msg, "")
    }
    return nil
}
