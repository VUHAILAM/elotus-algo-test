package files

import (
	"context"
	"hackathon/pkg/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ImageService struct {
	db *gorm.DB
}

func NewImageService(db *gorm.DB) *ImageService {
	return &ImageService{
		db: db,
	}
}

func (is *ImageService) Save(ctx context.Context, image models.Image) error {
	err := is.db.WithContext(ctx).Create(&image).Error
	if err != nil {
		logrus.Error("Save image metadata error", err)
		return err
	}

	return nil
}
