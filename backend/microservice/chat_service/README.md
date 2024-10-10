## Overview

This service is a WebSocket-based chat server with Kafka integration to handle messages. The server manages multiple chat rooms, allowing users to join, send messages, and interact in real-time. Kafka is used to ensure reliable message handling and persistence.

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

- **WebSocket Support**: Real-time communication for chat messages.
- **Kafka Integration**: Message publishing for reliability and scaling.
- **REST API**: Allows for chat management and interactions.
- **Multi-User Chat**: Support for adding/removing members from chats.
- **Search system**: Chat search system by hash tags and names

## Getting Started
### Prerequisites
- Go (version 1.22.2 or later)
- A valid config.json file [(see Configuration section)](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/chat_service/config.json)

## Installation
1. Clone the repository:

```bash
   git clone https://github.com/ssq0-0/InterestingChats.git
   cd InterestingChats/backend/microservice/chat_service
```

2. Set up your environment: Make sure you have Docker and Make installed. Verify your Docker installation by running:

```bash
   docker --version
```

or for make:
```bash
make --version
```

3. Start building the Docker container
```bash
make up
```


## API

The Chat Service provides several endpoints for managing chat functionalities. Below is the description of the available API methods.

### 1. GET /getAllChats
- **Description:** Retrieves all chats available in the system and sends the response to the client.
- **Headers:**
  - `X-User-ID: <userID>`
- **Response:**
  - **Success (200 OK):**
    ```json
    {
      "Errors": null,
      "Data": [
        {
          "id": 1,
          "creator": 123,
          "chat_name": "Chat 1",
          "members": {
            "123": {"id": 123, "username": "user1"},
            "456": {"id": 456, "username": "user2"}
          },
          "messages": []
        },
        {
          "id": 2,
          "creator": 456,
          "chat_name": "Chat 2",
          "members": {
            "123": {"id": 123, "username": "user1"},
            "789": {"id": 789, "username": "user3"}
          },
          "messages": []
        }
      ]
    }
    ```
  - **Error (400 Bad Request):**
    ```json
    {
      "Errors": ["Unable to retrieve chats"],
      "Data": null
    }
    ```

### 2. GET /getUserChats
- **Description:** Retrieves chats specifically for the authenticated user and sends the response to the client.
- **Headers:**
  - `X-User-ID: <userID>`
- **Response:**
  - **Success (200 OK):**
    ```json
    {
      "Errors": null,
      "Data": [
        {
          "id": 1,
          "chat_name": "User Chat 1",
          "members": ["user1", "user2"]
        },
        {
          "id": 2,
          "chat_name": "User Chat 2",
          "members": ["user3", "user4"]
        }
      ]
    }
    ```
  - **Error (400 Bad Request):**
    ```json
    {
      "Errors": ["Unable to retrieve user chats"],
      "Data": null
    }
    ```

### 3. GET /getChat
- **Description:** Retrieves the chat history for a specific chat and sends the response to the client.
- **Headers:**
  - `X-User-ID: <userID>`
- **Query Parameters:**
  - `chatID=<chatID>`
- **Response:**
  - **Success (200 OK):**
    ```json
    {
      "Errors": null,
      "Data": {
        "id": 1,
        "chat_name": "Chat 1",
        "messages": [
          {
            "id": 1,
            "body": "Hello!",
            "time": "2024-10-10T12:00:00Z",
            "user": {
              "id": 123,
              "username": "user1"
            }
          },
          {
            "id": 2,
            "body": "Hi!",
            "time": "2024-10-10T12:05:00Z",
            "user": {
              "id": 456,
              "username": "user2"
            }
          }
        ]
      }
    }
    ```
  - **Error (400 Bad Request):**
    ```json
    {
      "Errors": ["Unable to retrieve chat history"],
      "Data": null
    }
    ```

### 4. GET /getChatBySymbol
- **Description:** Performs a search for chats based on a query and sends the response to the client.
- **Headers:**
  - `X-User-ID: <userID>`
- **Query Parameters:**
  - `query=<searchQuery>`
- **Response:**
  - **Success (200 OK):**
    ```json
    {
      "Errors": null,
      "Data": [
        {
          "id": 1,
          "chat_name": "Chat 1"
        },
        {
          "id": 2,
          "chat_name": "Chat 2"
        }
      ]
    }
    ```
  - **Error (400 Bad Request):**
    ```json
    {
      "Errors": ["Search query failed"],
      "Data": null
    }
    ```

### 5. POST /createChat
- **Description:** Creates a new chat with the provided information and sends the response to the client.
- **Headers:**
  - `X-User-ID: <userID>`
- **Request Body:**
    ```json
    {
      "chat_name": "New Chat",
      "members": [123, 456]
    }
    ```
- **Response:**
  - **Success (201 Created):**
    ```json
    {
      "Errors": null,
      "Data": {
        "id": 3,
        "chat_name": "New Chat",
        "members": [123, 456]
      }
    }
    ```
  - **Error (400 Bad Request):**
    ```json
    {
      "Errors": ["Unable to create chat"],
      "Data": null
    }
    ```

### 6. DELETE /deleteChat
- **Description:** Deletes a chat identified by its ID and sends the response to the client.
- **Headers:**
  - `X-User-ID: <userID>`
- **Query Parameters:**
  - `chatID=<chatID>`
- **Response:**
  - **Success (200 OK):**
    ```json
    {
      "Errors": null,
      "Data": "successful delete"
    }
    ```
  - **Error (400 Bad Request):**
    ```json
    {
      "Errors": ["Unable to delete chat"],
      "Data": null
    }
    ```

### 7. POST /addMember
- **Description:** Adds members to a chat and sends the response to the client.
- **Headers:**
  - `X-User-ID: <userID>`
- **Request Body:**
    ```json
    {
      "chatID": 1,
      "members": [456, 789]
    }
    ```
- **Response:**
  - **Success (200 OK):**
    ```json
    {
      "Errors": null,
      "Data": "Members added successfully"
    }
    ```
  - **Error (400 Bad Request):**
    ```json
    {
      "Errors": ["Unable to add members"],
      "Data": null
    }
    ```

### 8. DELETE /deleteMember
- **Description:** Removes members from a chat and sends the response to the client.
- **Headers:**
  - `X-User-ID: <userID>`
- **Request Body:**
    ```json
    {
      "chatID": 1,
      "members": [456]
    }
    ```
- **Response:**
  - **Success (200 OK):**
    ```json
    {
      "Errors": null,
      "Data": "Members removed successfully"
    }
    ```
  - **Error (400 Bad Request):**
    ```json
    {
      "Errors": ["Unable to remove members"],
      "Data": null
    }
    ```

### 9. POST /joinToChat
- **Description:** Allows a user to join a chat and sends the response to the client.
- **Headers:**
  - `X-User-ID: <userID>`
- **Request Body:**
    ```json
    {
      "chatID": 1
    }
    ```
- **Response:**
  - **Success (200 OK):**
    ```json
    {
      "Errors": null,
      "Data": "Successfully joined the chat"
    }
    ```
  - **Error (400 Bad Request):**
    ```json
    {
      "Errors": ["Unable to join chat"],
      "Data": null
    }
    ```

### 10. PATCH /setTag
- **Description:** Assigns a tag to a chat and sends the response to the client.
- **Headers:**
  - `X-User-ID: <userID>`
- **Request Body:**
    ```json
    {
      "chatID": 1,
      "tag": "Important"
    }
    ```
- **Response:**
  - **Success (200 OK):**
    ```json
    {
      "Errors": null,
      "Data": "Tag added successfully"
    }
    ```
  - **Error (400 Bad Request):**
    ```json
    {
      "Errors": ["Unable to add tag"],
      "Data": null
    }
    ```

### 11. GET /getTags
- **Description:** Retrieves tags associated with chats and sends the response to the client.
- **Headers:**
  - `X-User-ID: <userID>`
- **Response:**
  - **Success (200 OK):**
    ```json
    {
      "Errors": null,
      "Data": [
        {
          "id": 1,
          "hashtag": "Important"
        },
        {
          "id": 2,
          "hashtag": "Urgent"
        }
      ]
    }
    ```
  - **Error (400 Bad Request):**
    ```json


    {
      "Errors": ["Unable to retrieve tags"],
      "Data": null
    }
    ```

### 12. DELETE /deleteTag
- **Description:** Removes a tag from a chat and sends the response to the client.
- **Headers:**
  - `X-User-ID: <userID>`
- **Request Body:**
    ```json
    {
      "chatID": 1,
      "tag": "Important"
    }
    ```
- **Response:**
  - **Success (200 OK):**
    ```json
    {
      "Errors": null,
      "Data": "Tag removed successfully"
    }
    ```
  - **Error (400 Bad Request):**
    ```json
    {
      "Errors": ["Unable to remove tag"],
      "Data": null
    }
    ```

### Notes
- All responses contain an `Errors` field indicating the presence of any errors and a `Data` field containing useful data if no errors occurred.
- Ensure that tokens are passed in the headers or requests in the correct format for successful authentication.

## Contributing
Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/YourFeature`).
3. Make your changes and commit them (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/YourFeature`).
5. Open a Pull Request.

## License
This project is licensed under the MIT License - see the [LICENSE](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/LICENSE) file for details.
