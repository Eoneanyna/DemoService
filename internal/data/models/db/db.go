package db

import "gorm.io/gorm"

type DB struct {
	Source *gorm.DB
}
