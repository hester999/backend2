package usecases

import (
	"backend2/internal/apperr"
	"backend2/internal/entity"
	"backend2/internal/utils"
	"fmt"
	"log"
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

func (p *Product) CreateProduct(product entity.Product) (entity.Product, error) {

	id, err := utils.GenerateUUID()
	if err != nil {
		return entity.Product{}, fmt.Errorf("error generating UUID: %w", err)
	}

	_, err = p.sup.GetSupplierById(product.SupplierId)
	if err != nil {
		return entity.Product{}, apperr.ErrSupplierNotFound
	}

	product.Id = id
	product.LastUpdate = time.Now()
	
	product, err = p.repo.CreateProduct(product)
	log.Println(product)
	if err != nil {
		return entity.Product{}, fmt.Errorf("error creating product: %w", err)
	}
	return product, nil
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
	product.LastUpdate = time.Now()
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
