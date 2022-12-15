package service

import (
	"gorm.io/gorm"
)

type ServiceRepository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{DB: db}
}
