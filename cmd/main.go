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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Получаем строку подключения к БД из переменной окружения
func getPostgresDSN() string {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://crm_user:crm_pass@localhost:5433/crm_db"
	}
	return dsn
}

func main() {
	// 1. Инициализируем БД
	database, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// 2. Запускаем миграции
	if err := runMigrations(database); err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	// 3. Инициализируем репозитории, сервисы, хендлеры
	adminRepo := admin.NewRepository(database)
	adminService := admin.NewService(adminRepo)
	adminHandler := admin.NewHandler(adminService)

	// Создаём супер-админа после миграций
	if err := adminRepo.InitSuperAdmin(); err != nil {
		log.Fatal("Ошибка создания супер-админа:", err)
	}

	shopRepo := shop.NewRepository(database)
	shopService := shop.NewService(shopRepo)
	shopHandler := shop.NewHandler(shopService)

	employeeRepo := employee.NewRepository(database)
	employeeService := employee.NewService(employeeRepo)
	employeeHandler := employee.NewHandler(employeeService)

	// 4. Настраиваем маршруты с CORS
	r := setupRoutes(adminHandler, shopHandler, employeeHandler)

	// 5. Запускаем сервер с graceful shutdown
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

	// Ожидаем SIGINT (Ctrl+C) или SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	// Завершаем сервер за 5 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("Server exiting")
}

// Подключение к PostgreSQL
func initDB() (*db.DB, error) {
	dsn := getPostgresDSN()
	database, err := db.NewDB(dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе: %w", err)
	}
	return database, nil
}

// Запуск всех миграций
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

// Настройка роутов с поддержкой CORS
func setupRoutes(
	adminHandler *admin.Handler,
	shopHandler *shop.Handler,
	employeeHandler *employee.Handler,
) *chi.Mux {
	r := chi.NewRouter()

	// 🌟 Добавляем CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Фронтенд
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // 5 минут кэширования CORS
	}))

	// Полезные middleware
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	// Авторизация
	r.Post("/admin/login", adminHandler.Login)

	// Управление пользователями
	r.Route("/admin/users", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/", adminHandler.GetUsers)
		r.Post("/", adminHandler.CreateUser)
		r.Put("/{id}", adminHandler.UpdateUser)
		r.Delete("/{id}", adminHandler.DeleteUser)
	})

	// Управление магазинами
	r.Route("/admin/shops", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Post("/", shopHandler.CreateShopHandler)
		r.Get("/", shopHandler.GetShopsHandler)
	})

	// Управление магазинами владельца
	r.Route("/owner/shops", func(r chi.Router) {
		r.Use(auth.AuthMiddleware)
		r.Get("/", shopHandler.GetShopsByOwner)

		// Управление сотрудниками магазина
		r.Route("/{id}/employees", func(r chi.Router) {
			r.Post("/", employeeHandler.AddEmployee)
			r.Get("/", employeeHandler.GetEmployeesByShop)
			r.Delete("/{employee_id}", employeeHandler.RemoveEmployee)
		})
	})

	return r
}
