package handlers

import (
	"context"
	"file_service/internal/logger"
	"fmt"
	"net/http"
	"time"

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
	// Ограничиваем размер загружаемого файла до 10MB
	r.ParseMultipartForm(10 << 20)

	// Получаем файл из формы
	file, handler, err := r.FormFile("image")
	if err != nil {
		is.log.Infof("Ошибка получения файла: %v", err)
		http.Error(w, "Ошибка получения файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Загрузка файла в Yandex Object Storage
	info, err := is.minioClient.PutObject(context.Background(), is.bucketName, handler.Filename, file, -1, minio.PutObjectOptions{ContentType: "image/png"})
	if err != nil {
		is.log.Infof("Ошибка загрузки файла: %v", err)
		http.Error(w, "Ошибка загрузки файла", http.StatusInternalServerError)
		return
	}
	is.log.Infof("req info: %+v", info)
	expires := time.Hour * 1 // Временная ссылка действительна 1 час
	presignedURL, err := is.minioClient.PresignedGetObject(context.Background(), is.bucketName, handler.Filename, expires, nil)
	if err != nil {
		is.log.Infof("Ошибка создания временной ссылки: %v", err)
		http.Error(w, "Ошибка создания ссылки", http.StatusInternalServerError)
		return
	}

	is.log.Infof("Файл успешно загружен: %s, Presigned URL: %s", handler.Filename, presignedURL.String())

	fmt.Fprintf(w, "Файл успешно загружен: %s\n", handler.Filename)
	fmt.Fprintf(w, "Ссылка на изображение: %s\n", presignedURL.String())
}
