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
	UpdateSupplier(supplier entity.Supplier) (entity.Supplier, error)
	DeleteSupplierById(id string) error
	GetAllSuppliers() ([]entity.Supplier, error)
}

func NewSupplier(repo SupplierRepository, addressRepo AddressRepo) *Supplier {
	return &Supplier{repo, addressRepo}
}

func (s *Supplier) CreateSupplier(name, phoneNumber, country, city, street string) (entity.Supplier, error) {
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
		Country: country,
		City:    city,
		Street:  street,
	})
	if err != nil {
		return entity.Supplier{}, fmt.Errorf("failed to save address: %w", err)
	}

	newSupplier := entity.Supplier{
		Id:          id,
		Name:        name,
		AddressId:   adrId,
		PhoneNumber: phoneNumber,
	}

	newSupplier, err = s.repo.CreateSupplier(newSupplier)
	if err != nil {
		return entity.Supplier{}, fmt.Errorf("failed to save supplier: %w", err)
	}
	return newSupplier, nil
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

func (s *Supplier) UpdateAddressSupplier(supplierId, country, city, street string) (entity.Supplier, error) {
	supplier, err := s.repo.GetSupplierById(supplierId)
	if err != nil {
		return entity.Supplier{}, fmt.Errorf("failed to get supplier: %w", err)
	}

	_, err = s.adr.Update(entity.Address{
		ID:      supplier.AddressId,
		Country: country,
		City:    city,
		Street:  street,
	})
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
