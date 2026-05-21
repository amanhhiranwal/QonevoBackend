package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	client     *s3.Client
	bucketName string
	region     string
}

func NewS3Service(
	bucketName string,
	region string,
) (*S3Service, error) {

	// =====================================
	// Load AWS Config
	// =====================================

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
	)

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3Service{
		client:     client,
		bucketName: bucketName,
		region:     region,
	}, nil
}

// =====================================
// Upload Single File
// =====================================

func (s *S3Service) UploadFile(
	file multipart.File,
	fileHeader *multipart.FileHeader,
	folder string,
) (string, error) {

	// =====================================
	// Validate File
	// =====================================

	if file == nil || fileHeader == nil {
		return "", fmt.Errorf("invalid file")
	}

	// =====================================
	// Allowed Extensions
	// =====================================

	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	extension := strings.ToLower(
		filepath.Ext(fileHeader.Filename),
	)

	if !allowedExtensions[extension] {

		return "", fmt.Errorf(
			"unsupported file type: %s",
			extension,
		)
	}

	// =====================================
	// File Size Validation
	// =====================================

	const maxFileSize = 10 << 20 // 10MB

	if fileHeader.Size > maxFileSize {

		return "", fmt.Errorf(
			"file size exceeds 10MB",
		)
	}

	// =====================================
	// Generate Unique File Name
	// =====================================

	fileName := fmt.Sprintf(
		"%d%s",
		time.Now().UnixNano(),
		extension,
	)

	key := fmt.Sprintf(
		"%s/%s",
		folder,
		fileName,
	)

	// =====================================
	// Upload To S3
	// =====================================

	_, err := s.client.PutObject(
		context.TODO(),
		&s3.PutObjectInput{
			Bucket:      aws.String(s.bucketName),
			Key:         aws.String(key),
			Body:        file,
			ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
		},
	)

	if err != nil {
		return "", err
	}

	// =====================================
	// Generate Public URL
	// =====================================

	fileURL := fmt.Sprintf(
		"https://%s.s3.%s.amazonaws.com/%s",
		s.bucketName,
		s.region,
		key,
	)

	return fileURL, nil
}

// =====================================
// Upload Multiple Files
// =====================================

func (s *S3Service) UploadMultipleFiles(
	files []*multipart.FileHeader,
	folder string,
) ([]string, error) {

	var uploadedURLs []string

	for _, fileHeader := range files {

		file, err := fileHeader.Open()

		if err != nil {
			return nil, err
		}

		defer file.Close()

		url, err := s.UploadFile(
			file,
			fileHeader,
			folder,
		)

		if err != nil {
			return nil, err
		}

		uploadedURLs = append(
			uploadedURLs,
			url,
		)
	}

	return uploadedURLs, nil
}

// =====================================
// Delete File
// =====================================

func (s *S3Service) DeleteFile(
	key string,
) error {

	_, err := s.client.DeleteObject(
		context.TODO(),
		&s3.DeleteObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(key),
		},
	)

	return err
}

// =====================================
// Get File Stream
// =====================================

func (s *S3Service) GetFile(
	key string,
) (io.ReadCloser, error) {

	result, err := s.client.GetObject(
		context.TODO(),
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(key),
		},
	)

	if err != nil {
		return nil, err
	}

	return result.Body, nil
}