package usecases

import (
	"backend2/internal/apperr"
	"backend2/internal/entity"
	"backend2/internal/utils"
	"errors"
	"fmt"
	"log"
	"time"
)

type ClientRepository interface {
	CreateClient(newClient entity.Client) (entity.Client, error)
	UpdateClient(id string, newClient entity.Client) (entity.Client, error)
	DeleteClient(id string) error
	GetAllClients(limit, offset int) ([]entity.Client, error)
	GetClientById(id string) (entity.Client, error)
	GetClients(name, surName string) ([]entity.Client, error)
}

type Client struct {
	repo ClientRepository
	adr  AddressRepo
}

func NewClient(repo ClientRepository, adrRepo AddressRepo) *Client {
	return &Client{
		repo: repo,
		adr:  adrRepo,
	}
}

func (c *Client) CreateClient(client entity.Client) (entity.Client, error) {

	id, err := utils.GenerateUUID()

	if err != nil {
		return entity.Client{}, fmt.Errorf("usecase: generate client id: %w", err)
	}

	addrId, err := utils.GenerateUUID()
	if err != nil {
		return entity.Client{}, fmt.Errorf("usecase: generate address id: %w", err)
	}

	newAdr := entity.Address{
		ID:      addrId,
		Country: client.Country,
		City:    client.City,
		Street:  client.Street,
	}

	_, err = c.adr.Save(newAdr)
	if err != nil {
		return entity.Client{}, fmt.Errorf("usecase: failed to add address: %w", err)
	}

	client.Id = id
	client.AddressId = addrId
	client.RegistrationDate = time.Now().UTC()

	res, err := c.repo.CreateClient(client)
	if err != nil {
		log.Printf("usecase: failed to create client: %w", err)
		return entity.Client{}, fmt.Errorf("usecase: failed to add client: %w", err)
	}
	return res, nil
}

func (c *Client) UpdateClient(id string, client entity.Client) (entity.Client, error) {

	client, err := c.repo.UpdateClient(id, client)
	if err != nil {
		if errors.Is(err, apperr.ErrClientNotFound) {
			return entity.Client{}, apperr.ErrClientNotFound
		}
		return entity.Client{}, fmt.Errorf("usecase: get client: %w", err)
	}

	newClient, err := c.repo.GetClientById(id)
	if err != nil {
		return entity.Client{}, fmt.Errorf("usecase: get updated client: %w", err)
	}
	return newClient, nil
}

func (c *Client) DeleteClient(id string) error {
	err := c.repo.DeleteClient(id)
	if err != nil {
		return fmt.Errorf("usecase: failed to delete client: %w", err)
	}
	return nil
}

func (c *Client) GetAllClients(limit, offset int) ([]entity.Client, error) {

	clients, err := c.repo.GetAllClients(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("usecase: failed to get all clients: %w", err)
	}
	return clients, nil
}

func (c *Client) GetClientsByNameSurname(name, surname string) ([]entity.Client, error) {
	clients, err := c.repo.GetClients(name, surname)
	if err != nil {
		return nil, fmt.Errorf("usecase: failed to get clients by name and surname: %w", err)
	}
	return clients, nil
}
