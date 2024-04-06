package imageupload

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"

	"google.golang.org/api/option"
)

var ctx = context.Background()
var opt = option.WithCredentialsFile("pkg/service/image-upload/serviceAccountKey.json") // Replace with your service account key path
var config = &firebase.Config{StorageBucket: "golangwithfirebase.appspot.com"}          // Replace with your Firebase Storage bucket name
// TODO : Change the serviceAccountKey.json with your file

func UploadService(file io.Reader, foldername string, filename string) (string, error) {

	path := "hotelmanagement/" + foldername + "/" + filename

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return "", err
	}

	client, err := app.Storage(ctx)
	if err != nil {
		return "", err
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return "", err
	}

	object := bucket.Object(path)

	wc := object.NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	// Set file ACL to public-read
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://storage.googleapis.com/golangwithfirebase.appspot.com/hotelmanagement/%s/%s", foldername, filename)

	return url, nil
}

func DeleteService(filename string) error {

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return err
	}

	client, err := app.Storage(ctx)
	if err != nil {
		return err
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return err
	}

	object := bucket.Object(filename)

	if err := object.Delete(ctx); err != nil {
		return err
	}

	return nil
}

func ImageUpload() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, handler, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()
		name := strings.ReplaceAll(handler.Filename, " ", "")
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), name)
		mediaLink, err := UploadService(file, "temp", filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"media_link": mediaLink})
	}
}

func ImageDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawurl := c.PostForm("url")

		url := strings.TrimPrefix(rawurl, "https://storage.googleapis.com/golangwithfirebase.appspot.com/")
		// it delete based on the path from root folder to filename
		if err := DeleteService(url); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully."})
	}
}
