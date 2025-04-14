package employee

import "time"

type Employee struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	ShopID   int       `json:"shop_id"`
	Position string    `json:"position,omitempty"`
	HiredAt  time.Time `json:"hired_at,omitempty"`
}

type AddEmployeeRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Position  string `json:"position,omitempty"`
}
