package mapper

import (
	"backend2/internal/dto"
	"backend2/internal/entity"
)

func SupplierEntityToDTO(supplier entity.Supplier) dto.SupplierResponseDTO {
	return dto.SupplierResponseDTO{
		Id:          supplier.Id,
		Name:        supplier.Name,
		PhoneNumber: supplier.PhoneNumber,
		AddressId:   supplier.AddressId,
		Address: dto.AddressDTO{
			ID:      supplier.Address.ID,
			Country: supplier.Address.Country,
			City:    supplier.Address.City,
			Street:  supplier.Address.Street,
		}}
}

func SupplierDTOToEntity(request dto.SupplierCreateRequestDTO) entity.Supplier {
	address := entity.Address{
		Country: request.Address.Country,
		City:    request.Address.City,
		Street:  request.Address.Street,
	}
	return entity.Supplier{
		Name:        request.Name,
		PhoneNumber: request.PhoneNumber,
		Address:     address,
	}

}

func SuppliersEntityToDTO(suppliers []entity.Supplier) dto.SuppliersResponse {
	suppliersDTO := dto.SuppliersResponse{
		Suppliers: make([]dto.SupplierResponseDTO, 0, len(suppliers)),
	}
	for _, supplier := range suppliers {
		tmp := SupplierEntityToDTO(supplier)
		suppliersDTO.Suppliers = append(suppliersDTO.Suppliers, tmp)
	}
	return suppliersDTO
}

func SupplierUpdateDTOToEntity(request dto.SupplierUpdateAddressRequestDTO) entity.Supplier {
	address := entity.Address{
		Country: request.Country,
		City:    request.City,
		Street:  request.Street,
	}
	return entity.Supplier{
		Address: address,
	}
}
