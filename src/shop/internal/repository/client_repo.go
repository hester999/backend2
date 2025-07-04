package repository

import (
	"backend2/internal/apperr"
	"backend2/internal/entity"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type ClientRepo struct {
	db *sql.DB
}

func NewClientRepo(db *sql.DB) *ClientRepo {
	return &ClientRepo{
		db: db,
	}
}

func (c *ClientRepo) CreateClient(newClient entity.Client) (entity.Client, error) {

	query := `
        INSERT INTO client (
            id, client_name, client_surname, birthday, gender, registration_date, address_id
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := c.db.Exec(
		query,
		newClient.Id,
		newClient.ClientName,
		newClient.ClientSurname,
		newClient.BirthDate,
		newClient.Gender,
		newClient.RegistrationDate,
		newClient.AddressId,
	)
	if err != nil {
		log.Print(newClient.Id)
		return entity.Client{}, fmt.Errorf("%w: %v", apperr.ErrInsertFailed, err)
	}
	newClient, err = c.GetClientById(newClient.Id)
	if err != nil {
		return entity.Client{}, fmt.Errorf("%w: %v", apperr.ErrInsertFailed, err)
	}
	return newClient, nil
}

func (c *ClientRepo) UpdateClient(id string, newClient entity.Client) (entity.Client, error) {
	query := `
        	UPDATE address
		SET country = $1, city = $2, street = $3
		WHERE id = (
			SELECT address_id FROM client WHERE id = $4
		)
	`

	result, err := c.db.Exec(
		query,
		newClient.Address.Country,
		newClient.Address.City,
		newClient.Address.Street,
		id,
	)
	if err != nil {
		return entity.Client{}, fmt.Errorf("%w: %v", apperr.ErrUpdateFailed, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return entity.Client{}, fmt.Errorf("checking update rows: %w", err)
	}
	if rowsAffected == 0 {
		return entity.Client{}, apperr.ErrClientNotFound
	}

	return newClient, nil
}

func (c *ClientRepo) DeleteClient(id string) error {
	query := `DELETE FROM client WHERE id = $1`
	result, err := c.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", apperr.ErrDeleteFailed, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking delete rows: %w", err)
	}
	if rowsAffected == 0 {
		return apperr.ErrClientNotFound
	}

	return nil
}

func (c *ClientRepo) GetClients(name, surName string) ([]entity.Client, error) {

	query := `
		SELECT client.id , client_name, client_surname, birthday, gender, registration_date, address_id, address.id as id, address.country as country, address.city as city, address.street as street
				FROM client
				inner join address  on address.id = client.address_id
		WHERE client_name = $1 AND client_surname = $2
	`

	rows, err := c.db.Query(query, name, surName)
	if err != nil {
		return nil, fmt.Errorf("failed to get clients: %w", err)
	}
	defer rows.Close()

	var clients []entity.Client
	for rows.Next() {
		var client entity.Client
		err := rows.Scan(
			&client.Id,
			&client.ClientName,
			&client.ClientSurname,
			&client.BirthDate,
			&client.Gender,
			&client.RegistrationDate,
			&client.AddressId,
			&client.Address.ID,
			&client.Address.Country,
			&client.Address.City,
			&client.Address.Street,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan client row: %w", err)
		}
		clients = append(clients, client)
	}

	if len(clients) == 0 {

		return make([]entity.Client, 0), apperr.ErrClientNotFound
	}

	return clients, nil
}

func (c *ClientRepo) GetAllClients(limit, offset int) ([]entity.Client, error) {
	var (
		query string
		rows  *sql.Rows
		err   error
	)

	if limit == 0 && offset == 0 {
		query = `SELECT client.id , client_name, client_surname, birthday, gender, registration_date, address_id, address.id as id, address.country as country, address.city as city, address.street as street
				FROM client
				inner join address  on address.id = client.address_id`
		rows, err = c.db.Query(query)
	} else if limit > 0 && offset > 0 {
		query = `SELECT client.id , client_name, client_surname, birthday, gender, registration_date, address_id, address.id as id, address.country as country, address.city as city, address.street as street
				FROM client
				inner join address  on address.id = client.address_id LIMIT $1 OFFSET $2`
		rows, err = c.db.Query(query, limit, offset)
	} else if limit > 0 {
		query = `SELECT client.id , client_name, client_surname, birthday, gender, registration_date, address_id, address.id as id, address.country as country, address.city as city, address.street as street
				FROM client
				inner join address  on address.id = client.address_id LIMIT $1`
		rows, err = c.db.Query(query, limit)
	} else {
		query = `SELECT client.id , client_name, client_surname, birthday, gender, registration_date, address_id, address.id as id, address.country as country, address.city as city, address.street as street
				FROM client
				inner join address  on address.id = client.address_id OFFSET $1`
		rows, err = c.db.Query(query, offset)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get clients: %w", err)
	}
	defer rows.Close()

	var clients []entity.Client
	for rows.Next() {
		var client entity.Client
		err := rows.Scan(
			&client.Id,
			&client.ClientName,
			&client.ClientSurname,
			&client.BirthDate,
			&client.Gender,
			&client.RegistrationDate,
			&client.AddressId,
			&client.Address.ID,
			&client.Address.Country,
			&client.Address.City,
			&client.Address.Street,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan client row: %w", err)
		}
		clients = append(clients, client)
	}

	if len(clients) == 0 {
		return nil, apperr.ErrClientNotFound
	}

	return clients, nil
}

func (c *ClientRepo) GetClientById(id string) (entity.Client, error) {
	query := `
		SELECT client.id , client_name, client_surname, birthday, gender, registration_date, address_id, address.id as id, address.country as country, address.city as city, address.street as street
				FROM client
				inner join address  on address.id = client.address_id
		WHERE client.id = $1
	`
	var client entity.Client
	err := c.db.QueryRow(query, id).Scan(
		&client.Id,
		&client.ClientName,
		&client.ClientSurname,
		&client.BirthDate,
		&client.Gender,
		&client.RegistrationDate,
		&client.AddressId,
		&client.Address.ID,
		&client.Address.Country,
		&client.Address.City,
		&client.Address.Street,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Client{}, apperr.ErrClientNotFound
	}
	if err != nil {
		return entity.Client{}, fmt.Errorf("failed to get client: %w", err)
	}

	return client, nil
}
