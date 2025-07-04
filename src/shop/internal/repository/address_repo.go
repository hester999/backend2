package repository

import (
	"backend2/internal/entity"
	"database/sql"
	"errors"
	"fmt"
)

type AddressRepo struct {
	db *sql.DB
}

func NewAddressRepo(db *sql.DB) *AddressRepo {
	return &AddressRepo{db: db}
}

func (a *AddressRepo) Save(address entity.Address) (entity.Address, error) {

	query := `insert into address (id, country,city,street) values ($1, $2, $3, $4)`

	_, err := a.db.Exec(query, address.ID, address.Country, address.City, address.Street)

	if err != nil {
		return entity.Address{}, fmt.Errorf("error save addres: %w", err)

	}
	return address, nil
}

func (a *AddressRepo) Update(address entity.Address) (entity.Address, error) {

	query := `UPDATE address set country = $1, city = $2, street = $3 WHERE id = $4`

	_, err := a.db.Exec(query, address.Country, address.City, address.Street, address.ID)
	if err != nil {
		return entity.Address{}, fmt.Errorf("error update addres: %w", err)
	}
	return address, nil
}

func (a *AddressRepo) Delete(address entity.Address) error {
	query := `DELETE FROM address WHERE id = $1`
	_, err := a.db.Exec(query, address.ID)
	if err != nil {
		return fmt.Errorf("error delete addres: %w", err)
	}
	return nil
}
func (a *AddressRepo) GetById(address entity.Address) (entity.Address, error) {
	query := `SELECT * FROM address WHERE id = $1`
	var addr entity.Address
	row := a.db.QueryRow(query, address.ID).Scan(&addr.ID, &addr.Country, &addr.City, &addr.Street)

	if errors.Is(row, sql.ErrNoRows) {
		return entity.Address{}, fmt.Errorf("address not found %w", sql.ErrNoRows)
	}
	return addr, nil

}
