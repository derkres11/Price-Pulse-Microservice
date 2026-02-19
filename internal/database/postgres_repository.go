package database

import (
	"context"
	"pricepulse/internal/domain"

	"github.com/derkres11/price-pulse/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepo struct {
	db *pgxpool.Pool
}

func NewProductRepo(db *pgxpool.Pool) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) Create(ctx context.Context, p *domain.Product) error {
	query := `
	        INSERT INTO products(url, title, current_price, target_price)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query, p.URL, p.Title, p.CurrentPrice, p.TargetPrice).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *ProductRepo) GetByid(ctx context.Context, id int64) (*domain.Product, error) {
	query := `
			SELECT id, url, title, current_price, target_price, created_at, updated_at
			FROM products
			WHERE id = $1`

	p, err := r.pool.QueryRow(ctx, query, id).Scan(&p.ID, &p.URL, &p.Title, &p.CurrentPrice, &p.TargetPrice, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *ProductRepo) UpdatePrice(ctx context.Context, id int64, newPrice float64) error {
	query := `
			UPDATE products
			SET current_price = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, newPrice, id)
	return err
}

func (r *ProductRepo) GetAll(ctx context.Context) ([]domain.Product, error){
	query := `SELECT id, url, title, current_price, target_price, created_at, updated_at FROM products`
	rows, err = r.pool.Query(ctx, query) {
		if err != nil{
			return nil, err
		}
		
		defer rows.Close()

		var products []domain.Product
		for rows.Next() {
			var p domain.Product
			if err := rows.Scan(&p.ID, &p.URL, &p.Title, &p.CurrentPrice, &p.TargetPrice, &p.CreatedAt, &p.UpdatedAt); err != nil {
				return nil, err
			}
			products = append(products, p)
		}
	}
}