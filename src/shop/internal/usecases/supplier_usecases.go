package usecases

import (
	"backend2/internal/entity"
	"backend2/internal/utils"
	"fmt"
)

type Supplier struct {
	repo SupplierRepository
	adr  AddressRepo
}

type SupplierRepository interface {
	CreateSupplier(supplier entity.Supplier) (entity.Supplier, error)
	GetSupplierById(id string) (entity.Supplier, error)
	UpdateSupplier(id string, supplier entity.Supplier) (entity.Supplier, error)
	DeleteSupplierById(id string) error
	GetAllSuppliers() ([]entity.Supplier, error)
}

func NewSupplier(repo SupplierRepository, addressRepo AddressRepo) *Supplier {
	return &Supplier{repo, addressRepo}
}

func (s *Supplier) CreateSupplier(supplier entity.Supplier) (entity.Supplier, error) {
	id, err := utils.GenerateUUID()
	if err != nil {
		return entity.Supplier{}, fmt.Errorf("failed to generate supplier id: %w", err)
	}

	adrId, err := utils.GenerateUUID()
	if err != nil {
		return entity.Supplier{}, fmt.Errorf("failed to generate address id: %w", err)
	}

	_, err = s.adr.Save(entity.Address{
		ID:      adrId,
		Country: supplier.Address.Country,
		City:    supplier.Address.City,
		Street:  supplier.Address.Street,
	})

	if err != nil {
		return entity.Supplier{}, fmt.Errorf("failed to save address: %w", err)
	}

	supplier.Id = id
	supplier.AddressId = adrId

	_, err = s.repo.CreateSupplier(supplier)
	resp, err := s.repo.GetSupplierById(supplier.Id)
	if err != nil {
		return entity.Supplier{}, fmt.Errorf("failed to save supplier: %w", err)
	}
	return resp, nil
}

func (s *Supplier) GetSupplierById(id string) (entity.Supplier, error) {
	supplier, err := s.repo.GetSupplierById(id)
	if err != nil {
		return entity.Supplier{}, fmt.Errorf("failed to get supplier: %w", err)
	}
	return supplier, nil
}

func (s *Supplier) GetAllSuppliers() ([]entity.Supplier, error) {
	suppliers, err := s.repo.GetAllSuppliers()
	if err != nil {
		return nil, fmt.Errorf("failed to get all suppliers: %w", err)
	}
	return suppliers, nil
}

func (s *Supplier) UpdateAddressSupplier(supplierId string, supplier entity.Supplier) (entity.Supplier, error) {

	supplier, err := s.repo.UpdateSupplier(supplierId, supplier)

	if err != nil {
		return entity.Supplier{}, fmt.Errorf("failed to update address: %w", err)
	}

	return supplier, nil
}

func (s *Supplier) DeleteSupplierById(id string) error {
	err := s.repo.DeleteSupplierById(id)
	if err != nil {
		return fmt.Errorf("failed to delete supplier: %w", err)
	}
	return nil
}
