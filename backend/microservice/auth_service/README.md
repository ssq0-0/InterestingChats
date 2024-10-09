# Auth service
## Overview

This Auth Service is a Golang-based authentication server designed to handle user authentication through JSON Web Tokens (JWT). It provides functionalities for user authorization, token generation, and token refresh, ensuring secure access to protected resources.
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
- JWT-Based Authentication: Secure user authentication using JWT tokens.
- Token Generation: Creates access and refresh tokens for users.
- Token Refresh: Validates and issues new access tokens using refresh tokens.
- RESTful API: Implements a set of RESTful endpoints for easy integration with client applications.

## Getting Started
### Prerequisites
- Go (version 1.22.2 or later)
- A valid config.json file [(see Configuration section)](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/auth_service/config.json)

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

The Auth Service provides several endpoints for managing user authentication. Below is the description of the available API methods.

### 1. POST /auth
- **Description:** Validates the JWT token from the Authorization header and returns the user ID if the token is valid.
- **Headers:**
  - `Authorization: Bearer <token>`
- **Response:**
  - **Success (200 OK):**
    ```json
    {
      "Errors": null,
      "Data": <userID>
    }
    ```
  - **Error (400 Bad Request):**
    ```json
    {
      "Errors": ["Internal server error"],
      "Data": null
    }
    ```

### 2. POST /generate_tokens
- **Description:** Generates and returns new access and refresh tokens for the user.
- **Request Body:**
    ```json
    {
      "email": "user@example.com",
      "username": "username",
      "password": "password"
    }
    ```
- **Response:**
  - **Success (200 OK):**
    ```json
    {
      "Errors": null,
      "Data": {
        "tokens": {
          "accessToken": "<accessToken>",
          "refreshToken": "<refreshToken>"
        }
      }
    }
    ```
  - **Error (400 Bad Request):**
    ```json
    {
      "Errors": ["Internal server error"],
      "Data": null
    }
    ```

### 3. POST /refreshToken
- **Description:** Validates the provided refresh token and generates a new access token.
- **Request Parameters:**
  - `refreshToken=<token>`
- **Response:**
  - **Success (200 OK):**
    ```json
    {
      "Errors": null,
      "Data": "<newAccessToken>"
    }
    ```
  - **Error (400 Bad Request):**
    ```json
    {
      "Errors": ["Invalid token"],
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
This project is licensed under the MIT License - see the [LICENSE](auth_service/LICENSE) file for details.
