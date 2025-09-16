package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/platonso/order-viewer/internal/domain"
)

type PostgresRepo struct {
	DB *pgxpool.Pool
}

func NewPostgresRepo(ctx context.Context, connStr string) (*PostgresRepo, error) {
	db, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return &PostgresRepo{DB: db}, nil
}

func (r *PostgresRepo) Save(ctx context.Context, order *domain.Order) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	orderQuery := `
		INSERT INTO orders (order_uid,track_number, entry, locale, internal_signature, 
        	customer_id, delivery_service, shardkey, sm_id,date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`
	_, err = tx.Exec(ctx, orderQuery,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	if err != nil {
		// Проверка, является ли ошибка PostgreSQL ошибкой уникальности (23505 - unique_violation)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrOrderAlreadyExists
		}
		return fmt.Errorf("failed to insert order: %w", err)
	}

	deliveryQuery := `
		INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
`
	_, err = tx.Exec(ctx, deliveryQuery,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)

	if err != nil {
		return fmt.Errorf("failed to insert delivery: %w", err)
	}

	paymentQuery := `
		INSERT INTO payments (order_uid, transaction, request_id, currency, provider, 
			 amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`
	_, err = tx.Exec(ctx, paymentQuery,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err != nil {
		return fmt.Errorf("failed to insert payment: %w", err)
	}

	itemQuery := `
		INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, 
        	sale, size, total_price, nm_id, brand, status)
    	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
`
	for _, item := range order.Items {
		_, err = tx.Exec(ctx, itemQuery,
			order.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			return fmt.Errorf("failed to insert item: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

func (r *PostgresRepo) FindByID(ctx context.Context, orderUID string) (*domain.Order, error) {
	orderQuery := `
		SELECT order_uid, track_number, entry, locale, internal_signature, customer_id,
		       delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders 
		WHERE order_uid = $1
	`

	var order domain.Order
	err := r.DB.QueryRow(ctx, orderQuery, orderUID).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, fmt.Errorf("failed to query order: %w", err)
	}

	deliveryQuery := `
		SELECT name, phone, zip, city, address, region, email
		FROM deliveries
		WHERE order_uid = $1
`

	err = r.DB.QueryRow(ctx, deliveryQuery, orderUID).Scan(
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.Zip,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query delivery: %w", err)
	}

	paymentQuery := `
		SELECT p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, 
		       p.bank, p.delivery_cost, p.goods_total, p.custom_fee
		FROM payments p
		WHERE order_uid = $1
`

	err = r.DB.QueryRow(ctx, paymentQuery, orderUID).Scan(
		&order.Payment.Transaction,
		&order.Payment.RequestID,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDt,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query payment: %w", err)
	}

	itemsQuery := `
		SELECT i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale, i.size, 
		       i.total_price, i.nm_id, i.brand, i.status
		FROM items i
		WHERE order_uid = $1
	`

	rows, err := r.DB.Query(ctx, itemsQuery, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	var items []domain.Item
	for rows.Next() {
		var item domain.Item
		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating items: %w", err)
	}

	order.Items = items

	return &order, nil
}

func (r *PostgresRepo) Close() {
	if r.DB != nil {
		r.DB.Close()
	}
}
