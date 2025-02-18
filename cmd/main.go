package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"crm-backend/internal/admin"
	"crm-backend/internal/db"
)

func main() {
	database, err := db.NewDB("postgres://crm_user:crm_pass@localhost:5433/crm_db")
	if err != nil {
		log.Fatal("Ошибка подключения к базе:", err)
	}
	defer database.Close()

	adminRepo := admin.NewRepository(database)
	adminService := admin.NewService(adminRepo)
	adminHandler := admin.NewHandler(adminService)

	// Миграция и супер-админ
	if err := adminRepo.Migrate(); err != nil {
		log.Fatal("Ошибка миграции:", err)
	}
	if err := adminRepo.InitSuperAdmin(); err != nil {
		log.Fatal("Ошибка создания супер-админа:", err)
	}

	// Роуты
	r := chi.NewRouter()
	r.Post("/admin/users", adminHandler.CreateUser)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
