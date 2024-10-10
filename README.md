# Interesting Chats

## Overview

**Interesting Chats** is an online chat application designed to facilitate discussions based on various interests. The project consists of multiple microservices, each handling different functionalities of the chat application, ensuring scalability and maintainability.


## Microservices Description

1. **API Gateway**: Serves as the entry point for client requests, routing them to the appropriate microservices. [Read features](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/api_gateway/README.md#features)
2. **Auth Service**: Manages user authentication and authorization. [Read features](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/auth_service/README.md#features)
3. **Chat Service**: Handles real-time chat functionalities, including sending and receiving messages. [Read features](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/chat_service/README.md#features)
4. **Database Service**: Interacts with the database for storing user data, chat histories, and notifications. [Read features](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/db/README.md)
5. **File Service**: Manages file uploads and storage for user-generated content. [Read features](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/file_service/README.md#feature)
6. **Kafka Service**: Manages message brokering between microservices using Kafka. [Read features](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/kafka/README.md)
7. **Notification Service**: Handles notifications for user actions and system events. [Read features](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/notification_service/README.md#features)
8. **Redis Service**: Caches frequently accessed data to improve performance. [Read features](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/redis/README.md#features)
9. **User Services**: Manages user-related functionalities such as profiles, friends, and settings. [Read features](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/user_services/README.md#features)

## Technologies Used
- Go
- Docker
- Kafka
- Redis
- PostgreSQL
- HTML/CSS/JavaScript for the frontend

## Architecture Diagram
[![2024-10-10-20-11-25.png](https://i.postimg.cc/mggtxsnZ/2024-10-10-20-11-25.png)](https://postimg.cc/QHvX5Rsv)
## Getting Started

To run the application locally, follow these steps:
Go to the detailed installation instructions for each service:

1. [**API Gateway**.](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/api_gateway/README.md#getting-started)
2. [**Auth Service**.](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/auth_service/README.md#getting-started)
3. [**Chat Service**.](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/chat_service/README.md#getting-started)
4. [**Database Service**.](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/db/README.md#getting-started)
5. [**File Service**.](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/file_service/README.md)
6. [**Kafka Service**.](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/kafka/README.md)
7. [**Notification Service**.](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/notification_service/README.md#getting-started)
8. [**Redis Service**.](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/redis/README.md#getting-started)
9. [**User Services**.](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/user_services/README.md#getting-started)

## Contributing
Contributions are welcome! Please create a pull request or open an issue to discuss potential changes.

## License
This project is licensed under the MIT License - see the LICENSE file for details.

## Contact
For inquiries or feedback, please reach out to @qss1161 in Telegram.

## Acknowledgments
- [Go](https://golang.org/) for the backend
- [Docker](https://www.docker.com/) for containerization
- [Kafka](https://kafka.apache.org/) for message brokering
