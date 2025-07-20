package repository

import (
	"github.com/jmoiron/sqlx"
)

type StorageInterface interface {
	InsertSub()
	GetSub()
	UpdateSub()
	DeleteSub()
	ListSub()
}

type Storage struct {
	DB *sqlx.DB
}
