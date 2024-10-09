# Notification Service

## Overview

The Notification Service is a Golang-based server responsible for handling notifications within the application. It processes incoming notification messages from a Kafka topic and provides a RESTful API for clients to retrieve and manage notifications.

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [API](#api)
- [Configuration](#configuration)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Features
- **Kafka Integration**: Consumes notification messages from Kafka topics.
- **RESTful API**: Provides endpoints to retrieve and update notifications.
- **Notification Handling**: Processes notifications for different user actions and updates their status.

## Getting Started
### Prerequisites
- Go (version 1.22.2 or later)
- A valid `config.json` file containing the necessary service configurations.

### Installation
1. Clone the repository:

    ```bash
    git clone https://github.com/ssq0-0/InterestingChats.git
    cd InterestingChats/backend/microservice/notifications
    ```

2. Set up your environment: Ensure that you have Docker and Make installed. Verify your Docker installation by running:

    ```bash
    docker --version
    ```

    or for Make:

    ```bash
    make --version
    ```

3. Start building the Docker container:

    ```bash
    make up
    ```

## API

The Notification Service provides several endpoints for managing notifications. Below is the description of the available API methods.

### 1. GET /getNotification
- **Description**: Retrieves notifications for a user based on their user ID.
- **Headers**:
  - `X-User-ID: <userID>`
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "Errors": null,
      "Data": [
        {
          "ID": 1,
          "UserID": 123,
          "SenderID": 456,
          "Type": "friend_request",
          "Message": "You have a new friend request.",
          "Time": "2023-10-10T12:34:56Z",
          "IsRead": false
        }
      ]
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "Errors": ["Invalid user ID"],
      "Data": null
    }
    ```

### 2. PATCH /readNotification
- **Description**: Marks a notification as read.
- **Request Body**:
    ```json
    [
      {
        "ID": 1
      }
    ]
    ```
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "Errors": null,
      "Data": "Notification marked as read."
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "Errors": ["Invalid notification data"],
      "Data": null
    }
    ```

### Notes
- All responses contain an `Errors` field indicating the presence of any errors and a `Data` field containing useful data if no errors occurred.
- Ensure to pass the user ID in the request headers when calling the `GET /getNotification` endpoint.

## Configuration
Your `config.json` file should include the necessary configurations for connecting to the Kafka service and the server settings. Here is a basic example:

```json
{
  "Server": {
    "Host": "localhost",
    "Port": "8008"
  },
  "Kafka": {
    "BootstrapServers": "127.0.0.1:9092",
    "GroupID": "notification_service",
    "Topics": ["friends_operation"],
    "AutoOffsetReset": "earliest",
    "EnableAutoCommit": true,
    "AutoCommitIntervalMs": 1000
  },
  "Services": {
    "db_service": {
      "Protocol": "http",
      "Host": "localhost",
      "Port": "8002"
    }
  }
}
```

## Contributing

Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/YourFeature`).
3. Make your changes and commit them (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/YourFeature`).
5. Open a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
