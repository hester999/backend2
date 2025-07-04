package usecases

import (
	"backend2/internal/apperr"
	"backend2/internal/entity"
	"backend2/internal/utils"
	"fmt"
	"time"
)

//Для товаров:
//
//Добавление товара (на вход подается json, соответствующей структуре, описанной сверху).
//
//Уменьшение количества товара (на вход запросу подается id товара и на сколько уменьшить)
//
//Получение товара по id
//
//Получение всех доступных товаров
//
//Удаление товара по id

type Product struct {
	repo ProductRepository
	sup  SupplierRepository
	img  ImageRepo
}
type ProductRepository interface {
	// Добавить новый продукт
	CreateProduct(product entity.Product) (entity.Product, error)
	// Получить продукт по ID
	GetProductById(id string) (entity.Product, error)
	// Уменьшить остаток по ID на count единиц
	ReduceProduct(id string, count int) (entity.Product, error)
	// Получить все продукты
	GetProducts() ([]entity.Product, error)
	// Удалить продукт по ID
	DeleteProduct(id string) error
}

func NewProduct(repo ProductRepository, supplier SupplierRepository, img ImageRepo) *Product {
	return &Product{repo, supplier, img}
}

//{
//id
//name
//category
//price
//available_stock // число закупленных экземпляров товара
//last_update_date // число последней закупки
//supplier_id
//image_id: UUID
//}

func (p *Product) CreateProduct(name, category, supplierId string, price float64, availableStock int, lastUpdate time.Time, img []byte) (entity.Product, error) {

	id, err := utils.GenerateUUID()
	if err != nil {
		return entity.Product{}, fmt.Errorf("error generating UUID: %w", err)
	}

	_, err = p.sup.GetSupplierById(supplierId)
	if err != nil {
		return entity.Product{}, apperr.ErrSupplierNotFound
	}
	imgId, err := utils.GenerateUUID()
	if err != nil {
		return entity.Product{}, fmt.Errorf("error generating UUID: %w", err)
	}

	_, err = p.img.AddImage(entity.Image{
		Id:    imgId,
		Image: img,
	})
	if err != nil {
		return entity.Product{}, fmt.Errorf("error adding image to product: %w", err)
	}
	newProduct := entity.Product{
		Id:             id,
		Name:           name,
		Category:       category,
		Price:          price,
		AvailableStock: availableStock,
		LastUpdate:     lastUpdate,
		ImageId:        imgId,
	}

	newProduct, err = p.repo.CreateProduct(newProduct)
	if err != nil {
		return entity.Product{}, fmt.Errorf("error creating product: %w", err)
	}
	return newProduct, nil
}

func (p *Product) GetProductById(id string) (entity.Product, error) {
	product, err := p.repo.GetProductById(id)
	if err != nil {
		return entity.Product{}, err
	}
	return product, nil
}

func (p *Product) ReduceProduct(id string, count int) (entity.Product, error) {
	product, err := p.repo.ReduceProduct(id, count)
	if err != nil {
		return entity.Product{}, err
	}
	return product, nil
}

func (p *Product) GetProducts() ([]entity.Product, error) {
	products, err := p.repo.GetProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *Product) DeleteProduct(id string) error {
	err := p.repo.DeleteProduct(id)
	if err != nil {
		return err
	}
	return nil
}
