# Redis Service

## Overview

The Redis Service is a Golang-based microservice designed to manage user sessions and relationships using Redis as the primary data store. This service integrates with Kafka for event-driven processing, allowing it to handle user data efficiently and respond to various events in real-time.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [API](#api)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Session Management**: Retrieve and manage user sessions stored in Redis.
- **Friend List Management**: Handle friends and subscribers of users with CRUD operations.
- **Kafka Integration**: Listen to Kafka topics for user-related events, allowing real-time updates in the Redis database.
- **RESTful API**: Provides a set of endpoints for easy interaction with client applications.

## Getting Started

### Prerequisites

- Go (version 1.22.2 or later)
- A running instance of Redis
- Kafka instance with necessary topics configured

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/ssq0-0/InterestingChats.git
    cd InterestingChats/backend/microservice/redis
    ```

2. Set up your environment: Make sure you have Docker and Make installed. Verify your Docker installation by running:

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

The Redis Service provides several endpoints for managing user sessions and relationships. Below is the description of the available API methods.

### 1. GET /getSession

- **Description**: Retrieves the session for a user based on the provided userID.
- **Query Parameters**:
  - `userID`: The ID of the user whose session is being requested.
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "Errors": null,
      "Data": {
        "session_id": "string",
        "id": 1,
        "username": "string",
        "email": "string",
        "avatar": "string"
      }
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "Errors": ["Internal server error"],
      "Data": null
    }
    ```

### 2. GET /getFriendList

- **Description**: Retrieves the list of friends for a user based on the provided userID.
- **Query Parameters**:
  - `userID`: The ID of the user whose friends are being requested.
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "Errors": null,
      "Data": [
        {
          "id": 1,
          "username": "string",
          "email": "string",
          "avatar": "string"
        }
      ]
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "Errors": ["failed get friends list"],
      "Data": null
    }
    ```

### 3. GET /getSubscribers

- **Description**: Retrieves the list of subscribers for a user based on the provided userID.
- **Query Parameters**:
  - `userID`: The ID of the user whose subscribers are being requested.
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "Errors": null,
      "Data": [
        {
          "id": 1,
          "username": "string",
          "email": "string",
          "avatar": "string"
        }
      ]
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "Errors": ["failed get subscribers list"],
      "Data": null
    }
    ```

### Notes
- All responses contain an `Errors` field indicating the presence of any errors and a `Data` field containing useful data if no errors occurred.
- Ensure that all requests include the necessary parameters for successful responses.

## Contributing

Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/YourFeature`).
3. Make your changes and commit them (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/YourFeature`).
5. Open a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
