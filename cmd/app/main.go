package main

import (
	"fmt"
	"kafka-demo/api/grpc"
	"kafka-demo/internal"
	"kafka-demo/pkg/kafka"
	"kafka-demo/utils"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)



func main() {

	// Connect to MongoDB
	_, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// close db connection
	defer db.DisconnectDB()

	// Access a specific collection
	users := db.GetCollection(db.USERS_COLL)
	notifications_log := db.GetCollection(db.NOTIFICATIONS_LOG_COLL)

	if users == nil || notifications_log == nil {
		log.Fatalf("Failed to connect to MongoDB")
		return 
	}

	err = godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return 
	}
	Port := os.Getenv("PORT")

	if Port == "" {
		log.Fatalf("PORT environment variable is not set in .env file")
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", Port))

	if err != nil {
		log.Fatalf("Failed to listen on port 8080: %v", err)
	}
	
	fmt.Println("server is running on port" , listener.Addr() , listener.Addr().Network())

	broker := os.Getenv("KAFKA_BROKERS")
	consumerGroup := os.Getenv("KAFKA_GROUP_ID")
	topic := os.Getenv("KAFKA_TOPIC")

	if broker == "" || consumerGroup == "" || topic == "" {
		log.Fatalf("Failed to get kafka variables")
	}

	brokers := []string {broker}

	// setup kafka consumer and producer 
	_ , err = kafka.SetupKafkaConsumer(brokers, "notifications-one")

	_ , err = kafka.SetupKafkaProducer(brokers)

	go func () {

		handler := kafka.NewConsumerGroupHandler(&utils.CustomMessageHandler{})
		err := kafka.ConsumeMessages([]string{topic}, handler)
		if err != nil {
			log.Fatalf("Error consuming messages: %v", err)
		}
	}()

	if err != nil {
		log.Fatalf("kafka service failed %v", err )
	}
	
	err = grpc.Serve(listener)

	if err != nil {
		log.Fatalf("Failed to listen on port 8080: %v", err)
	}

}

