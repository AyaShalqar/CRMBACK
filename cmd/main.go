package main

import (
	"context"
	"crm-backend/internal/admin"
	"crm-backend/internal/auth"
	"crm-backend/internal/db"
	"crm-backend/internal/employee"
	"crm-backend/internal/shop"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "crm-backend/docs" // Импорт сгенерированной документации

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/swag" // для генерации Swagger
)

// @title CRM Backend API
// @version 1.0
// @description API для CRM системы управления магазинами
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http

func getPostgresDSN() string {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://crm_user:crm_pass@localhost:5433/crm_db"
	}
	return dsn
}

func main() {

	database, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := runMigrations(database); err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	adminRepo := admin.NewRepository(database)
	adminService := admin.NewService(adminRepo)
	adminHandler := admin.NewHandler(adminService)

	if err := adminRepo.InitSuperAdmin(); err != nil {
		log.Fatal("Ошибка создания супер-админа:", err)
	}

	shopRepo := shop.NewRepository(database)
	shopService := shop.NewService(shopRepo)
	shopHandler := shop.NewHandler(shopService)

	employeeRepo := employee.NewRepository(database)
	employeeService := employee.NewService(employeeRepo)
	employeeHandler := employee.NewHandler(employeeService)

	authRepo := auth.NewRepository(database)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	r := setupRoutes(adminHandler, shopHandler, employeeHandler, authHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		fmt.Println("Server running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("Server exiting")
}

func initDB() (*db.DB, error) {
	dsn := getPostgresDSN()
	database, err := db.NewDB(dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе: %w", err)
	}
	return database, nil
}

func runMigrations(database *db.DB) error {
	adminRepo := admin.NewRepository(database)
	if err := adminRepo.Migrate(); err != nil {
		return fmt.Errorf("Ошибка миграции admin/users: %w", err)
	}

	shopRepo := shop.NewRepository(database)
	if err := shopRepo.Migrate(); err != nil {
		return fmt.Errorf("Ошибка миграции shops: %w", err)
	}
	if err := shopRepo.MigrateItems(); err != nil {
		return fmt.Errorf("Ошибка миграции items: %w", err)
	}

	employeeRepo := employee.NewRepository(database)
	if err := employeeRepo.Migrate(); err != nil {
		return fmt.Errorf("Ошибка миграции employees: %w", err)
	}

	return nil
}

func setupRoutes(
	adminHandler *admin.Handler,
	shopHandler *shop.Handler,
	employeeHandler *employee.Handler,
	authHandler *auth.Handler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"}, // Фронтенд
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Post("/auth/login", adminHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/auth/me", authHandler.Me)
	})

	r.Route("/admin/users", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/", adminHandler.GetUsers)
		r.Post("/", adminHandler.CreateUser)
		r.Put("/{id}", adminHandler.UpdateUser)
		r.Delete("/{id}", adminHandler.DeleteUser)
	})

	r.Route("/admin/shops", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Post("/", shopHandler.CreateShopHandler)
		r.Get("/", shopHandler.GetShopsHandler)
	})

	r.Route("/owner/shops", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/", shopHandler.GetShopsByOwner)

		r.Route("/{id}/employees", func(r chi.Router) {
			r.Post("/", employeeHandler.AddEmployee)
			r.Get("/", employeeHandler.GetEmployeesByShop)
			r.Delete("/{employee_id}", employeeHandler.RemoveEmployee)
		})
	})

	return r
}
