package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"crm-backend/internal/admin"
	"crm-backend/internal/auth"
	"crm-backend/internal/db"
	"crm-backend/internal/employee"
	"crm-backend/internal/shop"
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

	if err := adminRepo.Migrate(); err != nil {
		log.Fatal("Ошибка миграции:", err)
	}
	if err := adminRepo.InitSuperAdmin(); err != nil {
		log.Fatal("Ошибка создания супер-админа:", err)
	}

	r := chi.NewRouter()

	r.Post("/admin/login", adminHandler.Login)

	r.Route("/admin/users", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/", adminHandler.GetUsers)
		r.Post("/", adminHandler.CreateUser)
		r.Put("/{id}", adminHandler.UpdateUser)
		r.Delete("/{id}", adminHandler.DeleteUser)
	})
	shopRepo := shop.NewRepository(database)

	if err := shopRepo.Migrate(); err != nil {
		log.Fatal("Ошибка миграции shops:", err)
	}
	shopService := shop.NewService(shopRepo)
	shopHandler := shop.NewHandler(shopService)

	r.Route("/admin/shops", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Post("/", shopHandler.CreateShopHandler)
		r.Get("/", shopHandler.GetShopsHandler)
	})
	employeeRepo := employee.NewRepository(database)
	if err := employeeRepo.Migrate(); err != nil {
		log.Fatal("Ошибка миграции employees:", err)
	}
	employeeService := employee.NewService(employeeRepo)
	employeeHandler := employee.NewHandler(employeeService)

	r.Route("/owner/shops", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/", shopHandler.GetShopsByOwner)
		r.Route("/{id}/employees", func(r chi.Router) {
			r.Post("/", employeeHandler.AddEmployee)
			r.Get("/", employeeHandler.GetEmployeesByShop)
			r.Delete("/{employee_id}", employeeHandler.RemoveEmployee)
		})
	})
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
