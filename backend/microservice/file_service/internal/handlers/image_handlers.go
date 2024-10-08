package handlers

import (
	"file_service/internal/logger"
	"file_service/internal/models"
	"file_service/internal/services"
	"net/http"

	"github.com/minio/minio-go/v7"
)

type ImageService struct {
	log         logger.Logger
	minioClient *minio.Client
	bucketName  string
}

func NewImageService(log logger.Logger, minioClient *minio.Client, bucketName string) *ImageService {
	return &ImageService{
		log:         log,
		minioClient: minioClient,
		bucketName:  bucketName,
	}
}

func (is *ImageService) SaveImage(w http.ResponseWriter, r *http.Request) {
	response, clientErr, err := services.AddAvatar(r, is.log, is.minioClient, is.bucketName)
	if err != nil {
		is.HandleError(w, http.StatusBadRequest, []string{clientErr}, err)
		return
	}

	is.SendRespond(w, http.StatusOK, response)
}

func (is *ImageService) DeleteImage(w http.ResponseWriter, r *http.Request) {
	if clientErr, err := services.DeleteImage(r, is.log, is.minioClient, is.bucketName); err != nil {
		is.HandleError(w, http.StatusBadRequest, []string{clientErr}, err)
		return
	}

	// TODO: change status code to DELETED and refactor respond func
	is.SendRespond(w, http.StatusOK, &models.Response{
		Errors: nil,
		Data:   "successful deleted",
	})
}
