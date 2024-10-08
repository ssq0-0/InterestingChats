package utils

import (
	"context"
	"file_service/internal/consts"
	"file_service/internal/logger"
	"fmt"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

// Upload uploads a file to MinIO storage and returns presigned and constant URLs.
func Upload(minioClient minio.Client, file multipart.File, fileName, contentType, bucketName string, log logger.Logger) (string, string, string, error) {
	_, err := minioClient.PutObject(context.Background(), bucketName, fileName, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Errorf("Ошибка загрузки файла '%s': %v", fileName, err)
		return "", "", consts.ErrYandexInternalServerError, err
	}
	reqParams := make(url.Values)
	reqParams.Set("response-content-type", contentType)
	reqParams.Set("response-content-disposition", fmt.Sprintf("inline; filename=\"%s\"", fileName))

	expires := time.Hour * 120
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), bucketName, fileName, expires, reqParams)
	if err != nil {
		log.Errorf("Ошибка создания временной ссылки для файла '%s': %v", fileName, err)
		return "", "", consts.ErrYandexInternalServerError, err
	}
	log.Infof("Временная ссылка успешно создана: %s", presignedURL.String())

	constantURL := fmt.Sprintf("https://storage.yandexcloud.net/%s/%s", bucketName, fileName)
	log.Infof("Постоянная ссылка на объект: %s", constantURL)

	return presignedURL.String(), constantURL, "", nil
}

// Delete removes a file from MinIO storage.
func Delete(minioClient minio.Client, fileName, bucketName string, log logger.Logger) (string, error) {
	if err := minioClient.RemoveObject(context.Background(), bucketName, fileName, minio.RemoveObjectOptions{}); err != nil {
		log.Errorf("Ошибка удаления файла '%s': %v", fileName, err)
		return consts.ErrYandexInternalServerError, err
	}

	return "", nil
}
