package router

import (
	"github.com/gin-gonic/gin"
	"media-service/services/api_response"
	"media-service/services/common"
	"media-service/services/errors"
	"net/http"
	"os"
	"path/filepath"
)

type MediaHandler struct {
}

func NewMediaHandler() MediaHandler {
	return MediaHandler{}
}

func MediaRoutes(router *gin.RouterGroup, handler MediaHandler) {
	routerAuth := router.Group("media")
	{
		// màn hình đăng ký tài xế
		routerAuth.POST("/upload-media", handler.UploadImage)
	}
}

func (h MediaHandler) UploadImage(ctx *gin.Context) {

	// Get the image file from the request
	file, err := ctx.FormFile("upload-image")
	if err != nil {
		err = errors.ErrorResponse("No image uploaded", errors.ValidationError)
		_ = ctx.Error(err)
		return
	}

	// Get the absolute path of the current directory
	dir, err := os.Getwd()
	if err != nil {
		err = errors.ErrorResponse("Error getting current directory", errors.UnknownError)
		_ = ctx.Error(err)
		return
	}

	imageID := common.GenerateIdentityID()
	// Define the path to save the image
	path := filepath.Join(dir, "images", imageID+".jpg")

	// Save the image to the path
	if err = ctx.SaveUploadedFile(file, path); err != nil {
		err = errors.ErrorResponse("Error saving image", errors.UnknownError)
		_ = ctx.Error(err)
		return
	}

	// Get the host from the .env file
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost" // default value if HOST is not set
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default value if PORT is not set
	}

	// Create the image URL
	imageUrl := "http://" + host + ":" + port + "/images/" + imageID + ".jpg"
	imageSaveUrl := api_response.ImageSaveURL{ImageURL: imageUrl}

	rs := api_response.ImageSaveResponse{
		Meta: api_response.NewMetaData(),
		Data: imageSaveUrl,
	}

	// Return the URL of the image
	ctx.JSON(http.StatusOK, rs)
}
