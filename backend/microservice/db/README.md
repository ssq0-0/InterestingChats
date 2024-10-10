# Database Microservice

## Overview

The InterestingChats database microservice is part of the InterestingChats application. It handles user authentication, chat management, and data interactions with a database. The service utilizes Kafka for asynchronous message processing, improving scalability and performance.

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [API](#api)
- [Contributing](#contributing)
- [License](#license)
- 
## Technologies Used

- Go (Golang)
- PostgreSQL (or any SQL database)
- Kafka for messaging
- Gorilla Mux for routing

## Getting Started
### Prerequisites
- Go (version 1.22.2 or later)
- A valid config.json file [(see Configuration section)](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/db/config.json)

## Installation
1. Clone the repository:

```bash
   git clone https://github.com/ssq0-0/InterestingChats.git
   cd InterestingChats/backend/microservice/auth_service
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
### User Management
#### 1. User Registration

- **Endpoint:** `/registration`
- **Method:** `POST`
- **Description:** Handles user registration requests.
- **Request Body:**
  ```json
  {
    "username": "string",
    "email": "string",
    "password": "string"
  }
  ```
- **Response:**
  - **201 Created**
    ```json
    {
      "errors": null,
      "data": {
        "user_id": "int",
        "username": "string",
        "email": "string"
      }
    }
    ```
  - **400 Bad Request** (if input is invalid)

#### 2. User Login

- **Endpoint:** `/login`
- **Method:** `POST`
- **Description:** Handles user login requests.
- **Request Body:**
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": {
        "user_id": "int",
        "username": "string"
      }
    }
    ```
  - **400 Bad Request** (if email or password is incorrect)


#### 3. Check User Existence

- **Endpoint:** `/checkUser`
- **Method:** `GET`
- **Description:** Verifies if a user exists based on the provided user ID.
- **Query Parameters:**
  - `userID` (int, required): The ID of the user to check.
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": "int"
    }
    ```
  - **400 Bad Request** (if user not found)

#### 4. Get User Profile Information

- **Endpoint:** `/profileInfo`
- **Method:** `GET`
- **Description:** Retrieves detailed user profile information.
- **Query Parameters:**
  - `userID` (int, required): The ID of the user whose information is to be retrieved.
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": {
        "user_id": "int",
        "username": "string",
        "email": "string",
        // Additional profile fields
      }
    }
    ```
  - **400 Bad Request** (if user not found)

#### 5. Search Users

- **Endpoint:** `/searchUsers`
- **Method:** `GET`
- **Description:** Searches for users based on input parameters.
- **Query Parameters:**
  - `symbols` (string, required): The search criteria.
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": [
        {
          "user_id": "int",
          "username": "string",
          "email": "string"
        }
      ]
    }
    ```

#### 6. Change User Data

- **Endpoint:** `/changeUserData`
- **Method:** `POST`
- **Description:** Modifies user data in the system.
- **Request Body:**
  ```json
  {
    "type": "string", // "username" or "email"
    "data": "string",
    "user_id": "int"
  }
  ```
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": "successful changed"
    }
    ```
  - **400 Bad Request** (if input is invalid)

#### 7. Upload User Photo

- **Endpoint:** `/uploadPhoto`
- **Method:** `POST`
- **Description:** Handles photo uploads by users.
- **Request Body:**
  ```json
  {
    "user_id": "int",
    "photo": "base64_string"
  }
  ```
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": null
    }
    ```

### Friendship Operations
#### 8. Request to Friend

- **Endpoint:** `/requestToFriendShip`
- **Method:** `POST`
- **Description:** Sends a friendship request to another user.
- **Request Body:**
  ```json
  {
    "user_id": "int", // Requesting user ID
    "friend_id": "int" // User to send request to
  }
  ```
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": {
        "user_id": "int",
        "friend_id": "int"
      }
    }
    ```

#### 9. Accept Friendship Request

- **Endpoint:** `/acceptFriendShip`
- **Method:** `POST`
- **Description:** Accepts a friendship request from another user.
- **Request Body:**
  ```json
  {
    "user_id": "int", // User accepting the request
    "friend_id": "int" // User who sent the request
  }
  ```
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": {
        "user_id": "int",
        "friend_id": "int"
      }
    }
    ```


#### 10. Delete Friend

- **Endpoint:** `/deleteFriend`
- **Method:** `DELETE`
- **Description:** Removes a user from the friend list.
- **Request Body:**
  ```json
  {
    "user_id": "int",
    "friend_id": "int"
  }
  ```
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": {
        "user_id": "int",
        "friend_id": "int"
      }
    }
    ```


#### 11. Delete Friendship Request

- **Endpoint:** `/deleteFriendRequest`
- **Method:** `DELETE`
- **Description:** Cancels a friendship request.
- **Request Body:**
  ```json
  {
    "user_id": "int",
    "friend_id": "int"
  }
  ```
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": {
        "user_id": "int",
        "friend_id": "int"
      }
    }
    ```


#### 12. Get Friends List

- **Endpoint:** `/getFriends`
- **Method:** `GET`
- **Description:** Retrieves the list of a user's friends.
- **Query Parameters:**
  - `user_id` (int, required): The ID of the user whose friends are to be retrieved.
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": [
        {
          "user_id": "int",
          "username": "string"
        }
      ]
    }
    ```

#### 13. Get Subscribers List

- **Endpoint:** `/getSubs`
- **Method:** `GET`
- **Description:** Retrieves the list of a user's subscribers.
- **Query Parameters:**
  - `user_id` (int, required): The ID of the user whose subscribers are to be retrieved.
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": [
        {
          "user_id": "int",
          "username": "string"
        }
      ]
    }
    ```
### Notification Management
#### 1. Add Notification

- **Endpoint:** `/addNotification`
- **Method:** `POST`
- **Description:** Adds a new notification for a specified user.
- **Request Body:**
  ```json
  {
    "user_id": "int", // ID of the user receiving the notification
    "sender_id": "int", // ID of the user sending the notification
    "type": "string", // Type of notification (e.g., "message", "friend_request")
    "message": "string" // Content of the notification
  }
  ```
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": "notification add"
    }
    ```
  - **400 Bad Request** (if input is invalid)

---

#### 2. Get Notifications

- **Endpoint:** `/getNotifications`
- **Method:** `GET`
- **Description:** Retrieves all notifications for a specified user.
- **Query Parameters:**
  - `userID` (int, required): The ID of the user whose notifications are to be retrieved.
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": [
        {
          "id": "int",
          "user_id": "int",
          "sender_id": "int",
          "type": "string",
          "message": "string",
          "time": "string", // ISO 8601 format
          "is_read": "bool"
        }
      ]
    }
    ```
  - **400 Bad Request** (if userID is not found)

---

#### 3. Read Notifications

- **Endpoint:** `/readNotifications`
- **Method:** `POST`
- **Description:** Marks notifications as read for the specified user.
- **Request Body:**
  ```json
  {
    "user_id": "int", // ID of the user
    "notification_ids": ["int"] // Array of notification IDs to be marked as read
  }
  ```
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": "successful read notifications!"
    }
    ```
  - **400 Bad Request** (if userID or notification IDs are invalid)


### Chat Management

#### 1. Create Chat

- **Endpoint:** `/createChat`
- **Method:** `POST`
- **Description:** Creates a new chat room.
- **Request Body:**
  ```json
  {
    "creator": "int", // ID of the user creating the chat
    "chat_name": "string", // Name of the chat
    "members": ["int"], // List of user IDs to be added to the chat
    "hashtags": [{"hashtag": "string"}] // List of hashtags associated with the chat
  }
  ```
- **Response:**
  - **201 Created**
    ```json
    {
      "errors": null,
      "data": {
        "id": "int", // ID of the created chat
        "creator": "int",
        "chat_name": "string",
        "members": ["int"],
        "hashtags": [{"id": "int", "hashtag": "string"}]
      }
    }
    ```
  - **400 Bad Request** (if input is invalid)

#### 2. Get Chat

- **Endpoint:** `/getChat`
- **Method:** `GET`
- **Description:** Retrieves detailed information about a specific chat room.
- **Query Parameters:**
  - `chatID` (int, required): ID of the chat room to retrieve.
  - `userID` (int, required): ID of the user requesting the chat information.
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": {
        "id": "int",
        "creator": "int",
        "chat_name": "string",
        "members": [{"id": "int", "username": "string"}],
        "messages": [{"id": "int", "body": "string", "time": "string", "user": {"id": "int", "username": "string"}}]
      }
    }
    ```
  - **400 Bad Request** (if input is invalid)

#### 3. Get User Chats

- **Endpoint:** `/getUserChats`
- **Method:** `GET`
- **Description:** Retrieves all chats a user is a member of.
- **Query Parameters:**
  - `userID` (int, required): The ID of the user.
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": [
        {
          "id": "int",
          "chat_name": "string",
          "creator": "int",
          "members": ["int"] // List of member IDs
        }
      ]
    }
    ```
  - **400 Bad Request** (if userID is not found)

#### 4. Get All Chats

- **Endpoint:** `/getAllChats`
- **Method:** `GET`
- **Description:** Retrieves a list of all available chat rooms.
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": [
        {
          "id": "int",
          "chat_name": "string",
          "creator": "int",
          "members": ["int"]
        }
      ]
    }
    ```
  - **400 Bad Request** (if an error occurs)

#### 5. Search Chat

- **Endpoint:** `/searchChat`
- **Method:** `GET`
- **Description:** Searches for chat rooms by specific keywords.
- **Query Parameters:**
  - `chatName` (string, required): Keyword(s) to search for in chat names.
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": [
        {
          "id": "int",
          "chat_name": "string",
          "creator": "int"
        }
      ]
    }
    ```
  - **400 Bad Request** (if input is invalid)

#### 6. Delete Chat

- **Endpoint:** `/deleteChat`
- **Method:** `DELETE`
- **Description:** Deletes a chat room.
- **Query Parameters:**
  - `chatID` (int, required): ID of the chat room to delete.
- **Response:**
  - **204 No Content**
    ```json
    {
      "errors": null,
      "data": null
    }
    ```
  - **400 Bad Request** (if chatID is not found)

#### 7. Delete Member

- **Endpoint:** `/deleteMember`
- **Method:** `DELETE`
- **Description:** Removes a user from a chat room.
- **Request Body:**
  ```json
  {
    "user_id": "int", // ID of the user to remove
    "chat_id": "int" // ID of the chat from which the user is to be removed
  }
  ```
- **Response:**
  - **204 No Content**
    ```json
    {
      "errors": null,
      "data": null
    }
    ```
  - **400 Bad Request** (if input is invalid)

#### 8. Add Members

- **Endpoint:** `/addMembers`
- **Method:** `POST`
- **Description:** Adds one or more members to a chat room.
- **Request Body:**
  ```json
  {
    "chat_id": "int", // ID of the chat
    "user_id": "int" // ID of the user to add
  }
  ```
- **Response:**
  - **202 Accepted**
    ```json
    {
      "errors": null,
      "data": "successful added member"
    }
    ```
  - **400 Bad Request** (if input is invalid)

#### 9. Join to Chat

- **Endpoint:** `/joinToChat`
- **Method:** `POST`
- **Description:** Allows a user to join a chat room.
- **Request Body:**
  ```json
  {
    "chat_id": "int", // ID of the chat
    "user_id": "int" // ID of the user joining
  }
  ```
- **Response:**
  - **202 Accepted**
    ```json
    {
      "errors": null,
      "data": "successful joined chat"
    }
    ```
  - **400 Bad Request** (if input is invalid)

#### 10. Get Author

- **Endpoint:** `/getAuthor`
- **Method:** `GET`
- **Description:** Retrieves information about the creator of a chat room.
- **Query Parameters:**
  - `chatID` (int, required): ID of the chat room.
  - `userID` (int, required): ID of the user requesting author information.
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": {
        "user_id": "int" // ID of the chat creator
      }
    }
    ```
  - **400 Bad Request** (if input is invalid)

#### 11. Change Chat Name

- **Endpoint:** `/changeChatName`
- **Method:** `PATCH`
- **Description:** Updates the name of a chat room.
- **Request Body:**
  ```json
  {
    "chat_id": "int", // ID of the chat to change
    "chat_name": "string" // New name of the chat
  }
  ```
- **Response:**
  - **200 OK**
    ```json
    {
      "errors": null,
      "data": "new_chat_name" // Updated chat name
    }
    ```
  - **400 Bad Request** (if input is invalid)

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
