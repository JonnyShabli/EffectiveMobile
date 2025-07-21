package models

type Subscription struct {
	Sub_id       string `json:"sub_id" db:"sub_id"`
	Price        int    `json:"price" db:"price"`
	Service_name string `json:"service_name" db:"service_name"`
	User_id      string `json:"user_id" db:"user_id"`
	Start_date   string `json:"start_date" db:"start_date"`
	Created_at   int32  `json:"created_at" db:"created_at"`
	Updated_at   int32  `json:"updated_at" db:"updated_at"`
	Deleted_at   int32  `json:"deleted_at" db:"deleted_at"`
}

type SubscriptionDTO struct {
	price        int    `json:"price" db:"price"`
	service_name string `json:"service_name" db:"service_name"`
	user_id      string `json:"user_id" db:"user_id"`
	start_date   string `json:"start_date" db:"start_date"`
}
