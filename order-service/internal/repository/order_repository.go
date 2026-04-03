package repository

import (
	"context"
	"database/sql"

	"order-service/internal/domain"
)

type sqlOrderRepo struct {
	db *sql.DB
}

func NewSQLOrderRepository(db *sql.DB) domain.OrderRepository {
	return &sqlOrderRepo{db: db}
}

func (r *sqlOrderRepo) Create(ctx context.Context, o *domain.Order) error {
	query := `INSERT INTO orders (id, customer_id, item_name, amount, status, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, o.ID, o.CustomerID, o.ItemName, o.Amount, o.Status, o.CreatedAt)
	return err
}

func (r *sqlOrderRepo) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	query := `SELECT id, customer_id, item_name, amount, status, created_at FROM orders WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var o domain.Order
	err := row.Scan(&o.ID, &o.CustomerID, &o.ItemName, &o.Amount, &o.Status, &o.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *sqlOrderRepo) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

func (r *sqlOrderRepo) GetRecent(ctx context.Context, limit int) ([]*domain.Order, error) {
	query := "SELECT id, customer_id, item_name, amount, status, created_at FROM orders ORDER BY created_at DESC LIMIT $1"
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		o := &domain.Order{}
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.ItemName, &o.Amount, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
