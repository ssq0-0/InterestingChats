# File Service

A file management service for saving and deleting images using MinIO as the object storage service. This application provides an HTTP server with endpoints to handle image uploads and deletions.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [License](#license)
- 
## Features

- Upload images with presigned URLs.
- Delete images from storage.
- Built with Go and utilizes the MinIO SDK for object storage.


## Getting Started
### Prerequisites
- Go (version 1.22.2 or later)
- A valid config.json file [(see Configuration section)](https://github.com/ssq0-0/InterestingChats/blob/main/backend/microservice/file_service/config.json)

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

### 1. Save Image

- **Endpoint:** `/saveImage`
- **Method:** `POST`
- **Description:** Uploads an image to MinIO.
- **Request Body:** Form data containing the image file.
- **Response:**
  - **200 OK**
    ```json
    {
      "Errors": null,
      "TemporaryLink": "string", // Presigned URL for temporary access
      "StaticLink": "string" // Permanent URL for direct access
    }
    ```

### 2. Delete Image

- **Endpoint:** `/deleteImage`
- **Method:** `DELETE`
- **Description:** Deletes an image from MinIO.
- **Query Parameters:**
  - `fileName` (string, required): Name of the file to be deleted.
- **Response:**
  - **200 OK**
    ```json
    {
      "Errors": null,
      "Data": "successful deleted"
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
