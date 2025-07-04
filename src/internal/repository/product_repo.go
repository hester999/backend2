package repository

import (
	"backend2/internal/apperr"
	"backend2/internal/entity"
	"database/sql"
	"errors"
	"fmt"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}

func (p *ProductRepo) CreateProduct(product entity.Product) (entity.Product, error) {
	query := `INSERT INTO product (id, name, category, supplier_id, image_id, price, available_stock, last_update_date)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := p.db.Exec(query,
		product.Id,
		product.Name,
		product.Category,
		product.SupplierId,
		product.ImageId,
		product.Price,
		product.AvailableStock,
		product.LastUpdate,
	)
	if err != nil {
		return entity.Product{}, fmt.Errorf("%w: %v", apperr.ErrProductInsert, err)
	}

	return product, nil
}

func (p *ProductRepo) GetProductById(id string) (entity.Product, error) {
	query := `SELECT id, name, category, supplier_id, image_id, price, available_stock, last_update_date FROM product WHERE id = $1`

	var product entity.Product
	err := p.db.QueryRow(query, id).Scan(
		&product.Id,
		&product.Name,
		&product.Category,
		&product.SupplierId,
		&product.ImageId,
		&product.Price,
		&product.AvailableStock,
		&product.LastUpdate,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Product{}, apperr.ErrProductNotFound
	}
	if err != nil {
		return entity.Product{}, fmt.Errorf("failed to scan product: %w", err)
	}
	return product, nil
}

func (p *ProductRepo) ReduceProduct(id string, count int) (entity.Product, error) {
	query := `UPDATE product SET available_stock = available_stock - $1 WHERE id = $2`

	res, err := p.db.Exec(query, count, id)
	if err != nil {
		return entity.Product{}, fmt.Errorf("%w: %v", apperr.ErrProductUpdate, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return entity.Product{}, fmt.Errorf("checking update rows: %w", err)
	}
	if rowsAffected == 0 {
		return entity.Product{}, apperr.ErrProductNotFound
	}

	product, err := p.GetProductById(id)
	if err != nil {
		return entity.Product{}, fmt.Errorf("%w: %v", apperr.ErrProductUpdate, err)
	}

	return product, nil
}

func (p *ProductRepo) GetProducts() ([]entity.Product, error) {
	query := `SELECT id, name, category, supplier_id, image_id, price, available_stock, last_update_date FROM product`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []entity.Product
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Category,
			&product.SupplierId,
			&product.ImageId,
			&product.Price,
			&product.AvailableStock,
			&product.LastUpdate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product row: %w", err)
		}
		products = append(products, product)
	}

	if len(products) == 0 {
		return nil, apperr.ErrProductNotFound
	}

	return products, nil
}

func (p *ProductRepo) DeleteProduct(id string) error {
	query := `DELETE FROM product WHERE id = $1`

	res, err := p.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", apperr.ErrProductDelete, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking delete rows: %w", err)
	}
	if rowsAffected == 0 {
		return apperr.ErrProductNotFound
	}

	return nil
}
