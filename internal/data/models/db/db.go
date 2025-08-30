package db

import (
	"gorm.io/gorm"
)

type DB struct {
	Source *gorm.DB
}

func NewDb(db *gorm.DB) *DB {
	return &DB{Source: db}
}
