# User Service

## Overview

User Services is a part of the Interesting Chats system, providing user management features including registration, authentication, friend management, and user profile handling. This service interacts with Kafka for processing user-related events.

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
- **REST API**: Allows for user management and interactions.
- **Search Functionality**: Users can search for other users based on criteria.
- **Subscriber Management**: Users can view their subscribers and manage subscriber requests.
- **Kafka Integration**: Processes events related to user actions in real-time via Kafka.

## Getting Started
### Prerequisites
- Go (version 1.22.2 or later)
- A valid config.json file [(see Configuration section)](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/user_services/config.json)

## Installation
1. Clone the repository:

```bash
   git clone https://github.com/ssq0-0/InterestingChats.git
   cd InterestingChats/backend/microservice/user_services
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
### 1. POST /registration
- **Description**: Handles user registration.
- **Request Body**:
    ```json
    {
      "email": "user@example.com",
      "username": "newuser",
      "password": "securepassword"
    }
    ```
- **Response**:
  - **Success (201 Created)**:
    ```json
    {
      "Errors": null,
      "Data": {
        "tokens": {
          "accessToken": "access.token.here",
          "refreshToken": "refresh.token.here"
        },
        "user": {
          "id": 1,
          "email": "user@example.com",
          "username": "newuser",
          "avatar": null
        }
      }
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "Errors": ["Invalid user data"],
      "Data": null
    }
    ```

### 2. POST /login
- **Description**: Handles user login.
- **Request Body**:
    ```json
    {
      "email": "user@example.com",
      "password": "securepassword"
    }
    ```
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "Errors": null,
      "Data": {
        "tokens": {
          "accessToken": "access.token.here",
          "refreshToken": "refresh.token.here"
        },
        "user": {
          "id": 1,
          "email": "user@example.com",
          "username": "newuser",
          "avatar": null
        }
      }
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "Errors": ["Invalid credentials"],
      "Data": null
    }
    ```

### 3. GET /my_profile
- **Description**: Retrieves the logged-in user's profile.
- **Headers**:
  - `Authorization: Bearer <accessToken>`
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "Errors": null,
      "Data": {
        "id": 1,
        "email": "user@example.com",
        "username": "newuser",
        "avatar": null
      }
    }
    ```
  - **Error (401 Unauthorized)**:
    ```json
    {
      "Errors": ["Unauthorized"],
      "Data": null
    }
    ```

### 4. GET /user_profile
- **Description**: Retrieves another user's profile.
- **Query Parameters**:
  - `user_id`: The ID of the user whose profile is being requested.
- **Headers**:
  - `Authorization: Bearer <accessToken>`
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "Errors": null,
      "Data": {
        "id": 2,
        "email": "friend@example.com",
        "username": "frienduser",
        "avatar": null
      }
    }
    ```
  - **Error (404 Not Found)**:
    ```json
    {
      "Errors": ["User not found"],
      "Data": null
    }
    ```

### 5. PATCH /changeData
- **Description**: Allows the user to update their information.
- **Request Body**:
    ```json
    {
      "type": "username",
      "Data": "updatedusername",
      "user_id": 1
    }
    ```
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "Errors": null,
      "Data": "User data updated successfully."
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "Errors": ["Invalid data"],
      "Data": null
    }
    ```

### 6. GET /searchUsers
- **Description**: Handles searching for users.
- **Query Parameters**:
  - `query`: The search term for users.
- **Headers**:
  - `Authorization: Bearer <accessToken>`
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "Errors": null,
      "Data": [
        {
          "id": 2,
          "email": "friend@example.com",
          "username": "frienduser",
          "avatar": null
        }
      ]
    }
    ```
  - **Error (404 Not Found)**:
    ```json
    {
      "Errors": ["No users found"],
      "Data": null
    }
    ```

### 7. POST /requestToFriendShip
- **Description**: Sends a friend request.
- **Request Body**:
    ```json
    {
      "user_id": 1,
      "friend_id": 2
    }
    ```
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "Errors": null,
      "Data": "Friend request sent."
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "Errors": ["Invalid friend request data"],
      "Data": null
    }
    ```

### 8. POST /acceptFriendShip
- **Description**: Accepts a friend request.
- **Request Body**:
    ```json
    {
      "user_id": 2,
      "friend_id": 1
    }
    ```
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "Errors": null,
      "Data": "Friend request accepted."
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "Errors": ["Invalid request data"],
      "Data": null
    }
    ```

### 9. GET /getFriends
- **Description**: Retrieves the list of user's friends.
- **Headers**:
  - `Authorization: Bearer <accessToken>`
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "Errors": null,
      "Data": [
        {
          "id": 2,
          "email": "friend@example.com",
          "username": "frienduser",
          "avatar": null
        }
      ]
    }
    ```
  - **Error (404 Not Found)**:
    ```json
    {
      "Errors": ["No friends found"],
      "Data": null
    }
    ```

### 10. GET /getSubscribers
- **Description**: Retrieves the list of user's subscribers.
- **Headers**:
  - `Authorization: Bearer <accessToken>`
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "Errors": null,
      "Data": [
        {
          "id": 3,
          "email": "subscriber@example.com",
          "username": "subscriberuser",
          "avatar": null
        }
      ]
    }
    ```
  - **Error (404 Not Found)**:
    ```json
    {
      "Errors": ["No subscribers found"],
      "Data": null
    }
    ```

### 11. DELETE /deleteFriend
- **Description**: Removes a user from the friend list.
- **Request Body**:
    ```json
    {
      "user_id": 1,
      "friend_id": 2
    }
    ```
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "Errors": null,
      "Data": "Friend removed."
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "Errors": ["Invalid data"],
      "Data": null
    }
    ```

### 12. DELETE /deleteFriendRequest
- **Description**: Cancels a friend request.
- **Request Body**:
    ```json
    {
      "user_id": 1,
      "friend_id": 2
    }
    ```
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "Errors": null,
      "Data": "Friend request deleted."
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "Errors": ["Invalid data"],
      "Data": null
    }
    ```

### 13. POST /saveImage
- **Description**: Allows the user to upload their avatar.
- **Request Body**: (Multipart/form-data)
    - File: `<avatar_file>`
- **Response**:
  - **Success (200 OK)**:
    ```json
    {
      "Errors": null,
      "Data": {
        "file_url": {
          "Errors": null,
          "temporary_url": "temporary.url.here",
          "static_url": "static.url.here"
        }
      }
    }
    ```
  - **Error (400 Bad Request)**:
    ```json
    {
      "Errors": ["Invalid file"],
      "Data": null
    }
    ```

### Notes
- All responses contain an `Errors` field indicating the presence of any errors and a `Data` field containing useful data if no errors occurred.
- Ensure to pass the user ID in the request headers when calling the `GET /getNotification` endpoint.


## Contributing

Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/YourFeature`).
3. Make your changes and commit them (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/YourFeature`).
5. Open a Pull Request.

## License
This project is licensed under the MIT License - see the [LICENSE](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/LICENSE) file for details.
