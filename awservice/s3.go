package awservice

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func ConnectS3(cfg aws.Config) *s3.Client {
	s3Client := s3.NewFromConfig(cfg)
	_, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		log.Fatalf("unable to load S3 buckets, %v", err)
	}
	log.Printf("[CONNECTED] to AWS S3 service")
	return s3Client
}

func UploadFileToS3(s3Client *s3.Client, c *gin.Context) {

	bucketName := os.Getenv("AWS_BUCKET")

	// Read the uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read uploaded file"})
		return
	}
	defer file.Close()

	// Use the filename from the uploaded file
	filename := header.Filename

	// Upload to S3
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
		return
	}

	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, filename)
	c.JSON(http.StatusOK, gin.H{"url": fileURL})
}

func DownloadFileFromS3(s3Client *s3.Client, c *gin.Context) {
	err := godotenv.Load()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load .env file"})
		return
	}
	bucketName := os.Getenv("AWS_BUCKET")

	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	expiration := time.Duration(5) * time.Minute
	presignClient := s3.NewPresignClient(s3Client)

	presignedURL, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	}, s3.WithPresignExpires(expiration))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate presigned URL: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": presignedURL.URL})
}
