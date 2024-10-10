# API Gateway

This project is an API Gateway for the Interesting Chats application, built using Go and the Fiber web framework. It acts as a proxy, routing requests to various services and handling authentication.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Authentication Middleware](#authentication-middleware)
- [Contributing](#contributing)
- [License](#license)

## Features

- Proxy requests to user, auth, notification, and chat services.
- Middleware for user authentication.
- CORS support.
- Easy-to-extend structure for adding new routes.

## Getting Started
### Prerequisites
- Go (version 1.22.2 or later)
- A valid `config.json` file containing the necessary service configurations.

### Installation
1. Clone the repository:

    ```bash
    git clone https://github.com/ssq0-0/InterestingChats.git
    cd InterestingChats/backend/microservice/api_gateway
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
## API Endpoints

The following endpoints are available through the API Gateway:

### Public Endpoints

- **POST /registration**: Register a new user.
- **POST /login**: Log in an existing user.
- **POST /refreshToken**: Refresh the authentication token.

### Protected Endpoints

All protected endpoints require a valid `Authorization` header.

- **GET /my_profile**: Get the profile of the authenticated user.
- **GET /user_profile**: Get a specific user's profile.
- **PATCH /changeData**: Update user data.
- **GET /searchUsers**: Search for users.
- **POST /auth**: Authenticate the user.

#### Friends Management

- **GET /getFriends**: Retrieve a list of friends.
- **GET /getSubscribers**: Retrieve a list of subscribers.
- **POST /requestToFriendShip**: Send a friend request.
- **POST /acceptFriendShip**: Accept a friend request.
- **DELETE /deleteFriend**: Remove a friend.
- **DELETE /deleteFriendRequest**: Cancel a friend request.

#### Notifications

- **GET /getNotification**: Get notifications.
- **PATCH /readNotification**: Mark a notification as read.

#### File Management

- **POST /saveImage**: Save an image for the user.

#### Chat Management

- **GET /getChat**: Get chat details.
- **GET /getChatBySymbol**: Search for chat by symbol.
- **GET /getAllChats**: Retrieve all chats.
- **GET /getUserChats**: Get chats for the authenticated user.
- **POST /joinToChat**: Join an existing chat.
- **POST /createChat**: Create a new chat.
- **DELETE /deleteChat**: Delete a chat.
- **DELETE /leaveChat**: Leave a chat.
- **POST /addMember**: Add a member to a chat.
- **DELETE /deleteMember**: Remove a member from a chat.
- **PATCH /changeChatName**: Change the name of a chat.
- **PATCH /setTag**: Set tags for a chat.
- **GET /getTags**: Get tags for a chat.
- **DELETE /deleteTags**: Delete tags from a chat.

## Authentication Middleware

The API Gateway includes middleware for user authentication. It checks for an `Authorization` header in incoming requests and validates the token against the auth service. If the token is valid, it extracts the user ID and makes it available in the request context.

## Contributing

Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/YourFeature`).
3. Make your changes and commit them (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature/YourFeature`).
5. Open a Pull Request.

## License
This project is licensed under the MIT License - see the [LICENSE](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/LICENSE) file for details.
