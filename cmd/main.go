package main

import (
	"crm-backend/internal/admin"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	admin.InitSuperAdmin()
	r := chi.NewRouter()
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
