package mapper

import (
	"backend2/internal/dto"
	"backend2/internal/entity"
)

func ProductDTOToEntity(request dto.ProductCreateRequest) entity.Product {
	return entity.Product{
		Name:           request.Name,
		Category:       request.Category,
		Price:          request.Price,
		AvailableStock: request.AvailableStock,
		SupplierId:     request.SupplierId,
	}
}

func ProductEntityToDTO(product entity.Product) dto.ProductResponse {
	return dto.ProductResponse{
		Id:             product.Id,
		Name:           product.Name,
		Category:       product.Category,
		Price:          product.Price,
		AvailableStock: product.AvailableStock,
		LastUpdate:     product.LastUpdate,
		SupplierId:     product.SupplierId,
		ImageId:        product.ImageId,
	}
}

func ProductsEntityToDTOs(products []entity.Product) dto.ProductsResponse {
	productsResponse := dto.ProductsResponse{
		Products: make([]dto.ProductResponse, 0, len(products)),
	}

	for _, product := range products {
		productsResponse.Products = append(productsResponse.Products, ProductEntityToDTO(product))
	}
	return productsResponse
}
