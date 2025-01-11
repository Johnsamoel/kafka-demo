# Kafka Demo App

Kafka Demo App is a simple demonstration application that showcases the integration of multiple technologies like **Kafka**, **Golang**, **gRPC**, **MongoDB**, **Sarama**, and **Gmail** for sending notifications to registered users. It follows a clean **Repository Design Pattern** and is containerized using **Docker** and **Docker Compose** for easy deployment.

---

## Key Features

- **User Notification System**: Sends email notifications to newly registered users using Gmail.
- **Kafka Integration**: Utilizes Kafka for message queueing and handling asynchronous tasks.
- **gRPC API**: Implements gRPC for efficient communication between services.
- **MongoDB**: Stores user and notification logs in a MongoDB database.
- **Sarama**: Leverages Sarama library for Kafka producers and consumers.
- **Repository Design Pattern**: Maintains clean and scalable code structure.
- **Dockerized**: Fully containerized with Docker and Docker Compose for streamlined deployment.

---

## Technologies Used

1. **Golang** - Main programming language.
2. **Apache Kafka** - For message streaming.
3. **gRPC** - Remote procedure call framework.
4. **MongoDB** - Database for storing user and notification information.
5. **Sarama** - Kafka library for Go.
6. **Gmail** - Used for sending email notifications.
7. **Docker & Docker Compose** - For containerization and orchestration.

---

## Prerequisites

- **Docker** and **Docker Compose** installed on your machine.
- A Gmail account to send email notifications (requires an App Password for Gmail).
- Ensure the following ports are available:
  - **MongoDB**: `27017`
  - **Kafka**: `9092`
  - **gRPC**: `8080`

---

## How It Works

1. A user registers via a gRPC API.
2. The application sends a Kafka message to a specific topic.
3. A Kafka consumer processes the message and triggers the email notification system.
4. The notification is sent to the user's registered email using Gmail.
5. Notification logs are stored in MongoDB.

---

## Setup and Installation
Clone the repository:

```bash
Copy code
git clone https://github.com/yourusername/kafka-demo.git
cd kafka-demo
```
Create a .env file based on the provided .env.example file and fill in the required variables:
```
env
Copy code
PORT=8080
KAFKA_BROKERS=localhost:9092
KAFKA_GROUP_ID=notifications-one
KAFKA_TOPIC=notifications
MONGO_URI=mongodb://localhost:27017
GMAIL_USERNAME=your-gmail-address
GMAIL_APP_PASSWORD=your-gmail-app-password
Build and run the application using Docker Compose:
```
```bash
Copy code
docker-compose up --build
```

## Project Structure

```plaintext
kafka-demo/
├── api/             
│   ├── grpc/              # gRPC API implementation
│       ├── pb             # result of proto bufer files
│       ├── proto          # all proto buffer files
│       ├── server         # start grpc server
│       ├── user           # implement user grpc api's
│ 
├── cmd/    
│   ├── app           
│       ├── main.go        # app entry point
│
├── internal/              # Internal libraries
│   ├── mongo/             # mongo methods and connections
│   ├── user/              # user mongo repository
│
├── pkg/
│   ├── kafka/         
│       ├── kafka.go        # Kafka producer and consumer logic
│
│── user/          
│   ├── Model/              # User model
│   ├── repo/               # User repository
│   ├── service/            # User service
│
├── utils/                  # helper methods
├── docker-compose.yml      # Docker Compose configuration
├── .env                    # Environment variables (not committed to Git)
├── README.md               # This file

