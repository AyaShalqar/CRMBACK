package shop

import (
	"context"
	"crm-backend/internal/db"
	"fmt"
)

type Repository struct {
	db *db.DB
}

func NewRepository(db *db.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Migrate() error {
	_, err := r.db.Conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS shops (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			owner_id INT REFERENCES users(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return fmt.Errorf("ошибка миграции shops: %w", err)
	}
	fmt.Println("Миграция shops выполнена успешно")
	return nil

}

func (r *Repository) CreateShop(ctx context.Context, shop Shop) error {
	_, err := r.db.Conn.Exec(ctx, `
		INSERT INTO shops (name, description, owner_id)
		VALUES ($1, $2, $3)
	`, shop.Name, shop.Description, shop.OwnerID)

	if err != nil {
		return fmt.Errorf("ошибка создания магазина: %w", err)
	}
	return nil
}

func (r *Repository) GetShops(ctx context.Context) ([]Shop, error) {
	rows, err := r.db.Conn.Query(ctx, `SELECT id, name, description, owner_id FROM shops`)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения магазинов: %w", err)
	}
	defer rows.Close()

	var shops []Shop
	for rows.Next() {
		var shop Shop
		if err := rows.Scan(&shop.ID, &shop.Name, &shop.Description, &shop.OwnerID); err != nil {
			return nil, fmt.Errorf("ошибка чтения магазина: %w", err)
		}
		shops = append(shops, shop)
	}

	return shops, nil
}

func (r *Repository) UpdateShop(ctx context.Context, shop Shop) error {
	_, err := r.db.Conn.Exec(ctx, `
		UPDATE shops SET name = $1, description = $2 WHERE id = $3
	`, shop.Name, shop.Description, shop.ID)

	if err != nil {
		return fmt.Errorf("ошибка обновления магазина: %w", err)
	}
	return nil
}

func (r *Repository) DeleteShop(ctx context.Context, id int) error {
	_, err := r.db.Conn.Exec(ctx, `DELETE FROM shops WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления магазина: %w", err)
	}
	return nil
}

func (r *Repository) GetShopsByOwner(ctx context.Context, ownerID int) ([]Shop, error) {
	rows, err := r.db.Conn.Query(ctx, `
	SELECT id, name, description, owner_id FROM shops WHERE owner_id = $1
	`, ownerID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения магазинов владельца: %w", err)
	}
	defer rows.Close()

	var shops []Shop

	for rows.Next() {
		var shop Shop
		if err := rows.Scan(&shop.ID, &shop.Name, &shop.Description, &shop.OwnerID); err != nil {
			return nil, fmt.Errorf("ошибка чтения магазина: %w", err)
		}
		shops = append(shops, shop)
	}
	return shops, nil
}

func (r *Repository) MigrateItems() error {
	_, err := r.db.Conn.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS items (
		id SERIAL PRIMARY KEY,
		shop_id INT REFERENCES shops(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL,
		brand VARCHAR(255) NOT NULL,
		category VARCHAR(255),
		size VARCHAR(50),
		purchase_price NUMERIC(10, 2) DEFAULT 0,
		sale_price NUMERIC(10, 2) DEFAULT 0,
		photo_url TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	`)
	if err != nil {
		return fmt.Errorf("ошибка миграции items: %w", err)
	}
	fmt.Println("Миграция items выполнена успешно")
	return nil
}

func (r *Repository) CreateItem(ctx context.Context, item *Item) error {
	query := `
	INSERT INTO items 
		(shop_id, name, brand, category, size, purchase_price, sale_price, photo_url, created_at, updated_at) 
	VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	RETURNING id
`
}
