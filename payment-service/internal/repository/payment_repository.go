package repository

import (
	"context"
	"database/sql"

	"payment-service/internal/domain"
)

type sqlPaymentRepo struct {
	db *sql.DB
}

func NewSQLPaymentRepository(db *sql.DB) domain.PaymentRepository {
	return &sqlPaymentRepo{db: db}
}

func (r *sqlPaymentRepo) Create(ctx context.Context, p *domain.Payment) error {
	query := `INSERT INTO payments (id, order_id, transaction_id, amount, status) 
              VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, p.ID, p.OrderID, p.TransactionID, p.Amount, p.Status)
	return err
}

func (r *sqlPaymentRepo) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
	query := `SELECT id, order_id, transaction_id, amount, status FROM payments WHERE order_id = $1`
	row := r.db.QueryRowContext(ctx, query, orderID)

	var p domain.Payment
	err := row.Scan(&p.ID, &p.OrderID, &p.TransactionID, &p.Amount, &p.Status)
	return &p, err
}
