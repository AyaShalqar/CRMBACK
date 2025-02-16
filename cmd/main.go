package main

import (
	"crm-backend/internal/admin"
	"crm-backend/internal/db"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	database, err := db.NewDB("postgres://crm_user:crm_pass@localhost:5433/crm_db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	adminRepo := admin.NewRepository(database)

	err = adminRepo.Migrate()
	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	err = adminRepo.InitSuperAdmin()
	if err != nil {
		log.Fatal("Ошибка создания супер-админа:", err)
	}

	r := chi.NewRouter()

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
