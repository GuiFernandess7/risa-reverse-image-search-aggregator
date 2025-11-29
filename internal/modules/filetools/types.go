package filetools

import (
	"mime/multipart"

	"gorm.io/gorm"
)

type ImageHandler struct {
	DB *gorm.DB
}

type UploadInput struct {
	Engine string `form:"engine" validate:"required"`
	File   multipart.File
	Header *multipart.FileHeader
}
