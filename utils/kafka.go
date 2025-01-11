package utils

import (
	"encoding/json"
	"kafka-demo/pkg/kafka"
	"log"
)

// CustomMessageHandler implements kafka.MessageHandler
type CustomMessageHandler struct{}

func (h *CustomMessageHandler) CustomHandle(key string, msg string, headers map[string]string) error {
	newMessage := kafka.KafkaNotificationMsg{}

	err := json.Unmarshal([]byte(msg), &newMessage)
	
	if err != nil {
		log.Printf("Error parsing message: %v\n", err)
		return err
	}

	if newMessage.NotificationStatus == kafka.NEW_USER_WELCOMING_NOTIFICATIONS {


		// generate welcoming template
		subject , body := GenerateSystemTestEmail(
		newMessage.User.FirstName + " " + newMessage.User.LastName ,"Kafka-demo-by-john" )

		// send email to the user
		err := SendEmail(newMessage.User.Email , subject , body)

		if err != nil {
			return err
		}
	}


	return nil
}