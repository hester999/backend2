package repository

import (
	"backend2/internal/apperr"
	"backend2/internal/entity"
	"database/sql"
	"errors"
	"fmt"
)

type SupplierRepo struct {
	db *sql.DB
}

func NewSupplier(db *sql.DB) *SupplierRepo {
	return &SupplierRepo{db: db}
}

func (s *SupplierRepo) CreateSupplier(supplier entity.Supplier) (entity.Supplier, error) {
	query := `INSERT INTO supplier (id, name, address_id, phone_number) VALUES ($1, $2, $3, $4)`

	_, err := s.db.Exec(query, supplier.Id, supplier.Name, supplier.AddressId, supplier.PhoneNumber)
	if err != nil {
		return entity.Supplier{}, fmt.Errorf("%w: %v", apperr.ErrSupplierInsert, err)
	}
	return supplier, nil
}

func (s *SupplierRepo) GetSupplierById(id string) (entity.Supplier, error) {
	query := `SELECT id, name, address_id, phone_number FROM supplier WHERE id = $1`
	var supplier entity.Supplier
	err := s.db.QueryRow(query, id).Scan(
		&supplier.Id,
		&supplier.Name,
		&supplier.AddressId,
		&supplier.PhoneNumber,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Supplier{}, apperr.ErrSupplierNotFound
	}
	if err != nil {
		return entity.Supplier{}, fmt.Errorf("error getting supplier: %w", err)
	}
	return supplier, nil
}

func (s *SupplierRepo) UpdateSupplier(supplier entity.Supplier) (entity.Supplier, error) {
	query := `UPDATE supplier SET name = $1, address_id = $2, phone_number = $3 WHERE id = $4`

	res, err := s.db.Exec(query, supplier.Name, supplier.AddressId, supplier.PhoneNumber, supplier.Id)
	if err != nil {
		return entity.Supplier{}, fmt.Errorf("%w: %v", apperr.ErrSupplierUpdate, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return entity.Supplier{}, fmt.Errorf("checking update rows: %w", err)
	}
	if rowsAffected == 0 {
		return entity.Supplier{}, apperr.ErrSupplierNotFound
	}

	newSupplier, err := s.GetSupplierById(supplier.Id)
	if err != nil {
		return entity.Supplier{}, fmt.Errorf("error getting updated supplier: %w", err)
	}
	return newSupplier, nil
}

func (s *SupplierRepo) DeleteSupplierById(id string) error {
	query := `DELETE FROM supplier WHERE id = $1`

	res, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", apperr.ErrSupplierDelete, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking delete rows: %w", err)
	}
	if rowsAffected == 0 {
		return apperr.ErrSupplierNotFound
	}

	return nil
}

func (s *SupplierRepo) GetAllSuppliers() ([]entity.Supplier, error) {
	query := `SELECT id, name, address_id, phone_number FROM supplier`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error getting suppliers: %w", err)
	}
	defer rows.Close()

	var suppliers []entity.Supplier
	for rows.Next() {
		var supplier entity.Supplier
		err := rows.Scan(&supplier.Id, &supplier.Name, &supplier.AddressId, &supplier.PhoneNumber)
		if err != nil {
			return nil, fmt.Errorf("error scanning supplier: %w", err)
		}
		suppliers = append(suppliers, supplier)
	}

	if len(suppliers) == 0 {
		return nil, apperr.ErrSupplierNotFound
	}

	return suppliers, nil
}
