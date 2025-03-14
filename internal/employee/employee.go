package employee

type Employee struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	ShopID   int    `json:"shop_id"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
