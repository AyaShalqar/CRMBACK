# 🧾 CRM Backend

**CRM Backend** — это серверная часть системы для управления офлайн-магазинами одежды. Она позволяет создавать магазины, добавлять сотрудников и товары, управлять ими и следить за активностью.

Проект построен на **Go**, использует **PostgreSQL** и реализует авторизацию через **JWT**, а также документирован с помощью **Swagger**.

---

## 🔧 Стек технологий

- 🟦 Go (Golang)
- 🐘 PostgreSQL
- 📦 pgx (PostgreSQL драйвер)
- 🧪 Chi router + middleware
- 🔐 JWT (golang-jwt)
- 🔒 bcrypt (для хеширования паролей)
- 📑 Swagger (swaggo)
- 🐳 Docker + Docker Compose

---

## 📚 Функциональность

### 🔐 Аутентификация
- `POST /auth/login` — логин по email и паролю
- `GET /auth/me` — текущий пользователь по токену

### 👑 Админ (`superadmin`)
- `POST /admin/users` — создать пользователя
- `GET /admin/users` — список пользователей
- `PUT /admin/users/{id}` — обновить пользователя
- `DELETE /admin/users/{id}` — удалить пользователя
- `POST /admin/shops` — создать магазин
- `GET /admin/shops` — список всех магазинов

### 🏪 Владелец (`owner`)
- `GET /owner/shops` — получить свои магазины
- `POST /owner/shops/{id}/employees` — добавить сотрудника
- `GET /owner/shops/{id}/employees` — список сотрудников
- `DELETE /owner/shops/{id}/employees/{employee_id}` — удалить сотрудника

### 📦 Товары
- `POST /owner/shops/{shopID}/items` — добавить товар
- `GET /owner/shops/{shopID}/items` — список товаров
- `GET /owner/shops/{shopID}/items/{itemID}` — получить товар
- `PUT /owner/shops/{shopID}/items/{itemID}` — обновить товар
- `DELETE /owner/shops/{shopID}/items/{itemID}` — удалить товар

---

## 🧑‍💼 Роли пользователей

- `superadmin` — администратор всей системы
- `owner` — владелец одного или нескольких магазинов
- `employee` — сотрудник магазина

---

## 🛠️ Установка и запуск

### 📄 1. Создай файл `.env`

```env
POSTGRES_USER=crm_user
POSTGRES_PASSWORD=crm_pass
POSTGRES_DB=crm_db

🐳 2. Запусти Docker
docker-compose up --build
API будет доступен по адресу:
📍 http://localhost:8080
Swagger-документация:
📄 http://localhost:8080/swagger/index.html

🔑 JWT Авторизация

После логина:

POST /auth/login
{
  "email": "admin@crm.kz",
  "password": "superAdmin123"
}
Вы получите:

{
  "token": "..."
}
Передавайте его в заголовках:

Authorization: Bearer <token>
🔍 Swagger

Для генерации документации:

go install github.com/swaggo/swag/cmd/swag@latest
swag init
Открыть в браузере:
🧭 http://localhost:8080/swagger/index.html

🗂 Структура проекта

crm-backend/
├── internal/
│   ├── admin/        // управление пользователями
│   ├── auth/         // JWT, middleware
│   ├── db/           // подключение к PostgreSQL
│   ├── employee/     // сотрудники магазинов
│   ├── shop/         // магазины и товары
├── main.go           // запуск сервера
├── Dockerfile
├── docker-compose.yml
├── README.md
└── docs/             // Swagger
📥 Примеры запросов

Добавить магазин:
POST /admin/shops
Authorization: Bearer <superadmin_token>

{
  "name": "Zara Mega",
  "description": "Магазин одежды",
  "owner_id": 2
}
Добавить сотрудника:
POST /owner/shops/1/employees
Authorization: Bearer <owner_token>

{
  "first_name": "Али",
  "last_name": "Ибрагимов",
  "email": "ali@shop.kz",
  "password": "pass123",
  "position": "кассир"
}
📄 Документация

PDF-документация проекта доступна в файле:
📎 CRM_Backend_README.pdf

🧠 Автор

Zeinaddin Zurgambayev
Backend-разработчик | Golang | Казахстан 🇰🇿
LinkedIn https://www.linkedin.com/in/zeinaddin-zurgambaev-5b5a392bb/

📌 Планы на будущее

 Заказы и корзины
 Статистика и отчёты
 Роли и права доступа (RBAC)
 UI для владельца и админа
 Отправка email-уведомлений
