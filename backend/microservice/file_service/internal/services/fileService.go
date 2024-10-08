package services

import (
	"file_service/internal/logger"
	"file_service/internal/models"
	"file_service/internal/utils"
	"net/http"

	"github.com/minio/minio-go/v7"
)

// AddAvatar processes the upload of an avatar image.
func AddAvatar(r *http.Request, log logger.Logger, minioClient *minio.Client, bucketName string) (*models.FileResponse, string, error) {
	file, _, fileName, contentType, clientErr, err := utils.ParseRequestFormFile(r)
	if err != nil {
		return nil, clientErr, err
	}

	presignedURL, constantURL, clientErr, err := utils.Upload(*minioClient, file, fileName, contentType, bucketName, log)
	if err != nil {
		return nil, clientErr, err
	}

	return &models.FileResponse{Errors: nil, TemporaryLink: presignedURL, StaticLink: constantURL}, "", nil
}

// DeleteImage processes the deletion of an image file.
func DeleteImage(r *http.Request, log logger.Logger, minioClient *minio.Client, bucketName string) (string, error) {
	if clientErr, err := utils.Delete(*minioClient, r.URL.Query().Get("fileName"), bucketName, log); err != nil {
		return clientErr, err
	}

	return "", nil
}
