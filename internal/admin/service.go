package admin

import (
	"fmt"
)

var users []User // времанное база данных

func InitSuperAdmin() {

	for _, user := range users {
		if user.Role == "superadmin" {
			fmt.Println("Супер-админ уже существует")
			return
		}
	}
}

// type User struct{
// 	 ID int
// 	 FirstName string
// 	 LastName string
// 	 Email string
// 	 PhoneNumber string
// 	 Role string
// }

// func (User) getUserById (userId int) User {
// 	return User{ID: 1, FirstName: "Dimash", LastName: "Arystambek", Email: "Hello@gmail.com"}
// }
