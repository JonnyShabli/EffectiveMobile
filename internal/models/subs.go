package models

import (
	"database/sql"
	"time"
)

type Subscription struct {
	Sub_id       int          `json:"sub_id" db:"sub_id"`
	Service_name string       `json:"service_name" db:"service_name"`
	Price        int          `json:"price" db:"price"`
	User_id      string       `json:"user_id" db:"user_id"`
	Start_date   string       `json:"start_date" db:"start_date"`
	Created_at   time.Time    `json:"created_at" db:"created_at"`
	Updated_at   time.Time    `json:"updated_at" db:"updated_at"`
	Deleted_at   sql.NullTime `json:"deleted_at" db:"deleted_at"`
}

type SubscriptionDTO struct {
	Sub_id       int    `json:"sub_id" db:"sub_id"`
	Service_name string `json:"service_name" db:"service_name"`
	Price        int    `json:"price" db:"price"`
	User_id      string `json:"user_id" db:"user_id"`
	Start_date   string `json:"start_date" db:"start_date"`
}
