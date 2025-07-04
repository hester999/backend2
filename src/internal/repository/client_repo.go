package repository

import (
	"backend2/internal/apperr"
	"backend2/internal/entity"
	"database/sql"
	"errors"
	"fmt"
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
		return entity.Client{}, fmt.Errorf("%w: %v", apperr.ErrInsertFailed, err)
	}

	return newClient, nil
}

func (c *ClientRepo) UpdateClient(newClient entity.Client) (entity.Client, error) {
	query := `
        UPDATE client 
        SET 
            client_name = $1,
            client_surname = $2,
            birthday = $3,
            gender = $4,
            registration_date = $5,
            address_id = $6
        WHERE id = $7
    `
	result, err := c.db.Exec(
		query,
		newClient.ClientName,
		newClient.ClientSurname,
		newClient.BirthDate,
		newClient.Gender,
		newClient.RegistrationDate,
		newClient.AddressId,
		newClient.Id,
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
		SELECT id, client_name, client_surname, birthday, gender, registration_date, address_id
		FROM client 
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
		query = `SELECT id, client_name, client_surname, birthday, gender, registration_date, address_id FROM client`
		rows, err = c.db.Query(query)
	} else if limit > 0 && offset > 0 {
		query = `SELECT id, client_name, client_surname, birthday, gender, registration_date, address_id FROM client LIMIT $1 OFFSET $2`
		rows, err = c.db.Query(query, limit, offset)
	} else if limit > 0 {
		query = `SELECT id, client_name, client_surname, birthday, gender, registration_date, address_id FROM client LIMIT $1`
		rows, err = c.db.Query(query, limit)
	} else {
		query = `SELECT id, client_name, client_surname, birthday, gender, registration_date, address_id FROM client OFFSET $1`
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
		SELECT id, client_name, client_surname, birthday, gender, registration_date, address_id 
		FROM client 
		WHERE id = $1
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
	)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Client{}, apperr.ErrClientNotFound
	}
	if err != nil {
		return entity.Client{}, fmt.Errorf("failed to get client: %w", err)
	}

	return client, nil
}
