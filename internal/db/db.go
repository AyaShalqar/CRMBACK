package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	Conn *pgx.Conn
}

func NewDB(dsn string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе: %w", err)
	}

	fmt.Println("Подключение к базе установлено")
	return &DB{Conn: conn}, nil
}

func (db *DB) Close() {
	db.Conn.Close(context.Background())
	fmt.Println("Соединение с базой закрыто")
}
