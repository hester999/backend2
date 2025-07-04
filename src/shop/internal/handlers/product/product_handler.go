package product

import (
	"backend2/internal/apperr"
	"backend2/internal/dto"
	"backend2/internal/entity"
	"backend2/internal/mapper"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Product interface {
	CreateProduct(product entity.Product) (entity.Product, error)
	GetProductById(id string) (entity.Product, error)
	ReduceProduct(id string, count int) (entity.Product, error)
	GetProducts() ([]entity.Product, error)
	DeleteProduct(id string) error
}
type ProductHandler struct {
	product Product
}

func NewProductHandler(productRepo Product) *ProductHandler {
	return &ProductHandler{product: productRepo}
}

// CreateProduct godoc
// @Summary      Создать товар
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body     dto.ProductCreateRequest  true  "Создаваемый товар"
// @Success      200      {object} dto.ProductResponse
// @Failure      400      {object} dto.Error400
// @Failure      500      {object} dto.Error500
// @Router       /product [post]
func (p *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product dto.ProductCreateRequest

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid JSON",
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.(validator.ValidationErrors).Error(),
		})
		return
	}
	productEntity := mapper.ProductDTOToEntity(product)
	productEntity, err = p.product.CreateProduct(productEntity)
	log.Println(productEntity)
	if err != nil {
		if errors.Is(err, apperr.ErrSupplierNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "supplier not found",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}
	res := mapper.ProductEntityToDTO(productEntity)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// GetProductById godoc
// @Summary      Получить товар по ID
// @Tags         products
// @Produce      json
// @Param        id   path     string  true  "ID товара"
// @Success      200  {object} dto.ProductResponse
// @Failure      400  {object} dto.Error400
// @Failure      404  {object} dto.Error404
// @Failure      500  {object} dto.Error500
// @Router       /product/{id} [get]
func (p *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")
	id := mux.Vars(r)["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid ID",
		})
		return
	}
	product, err := p.product.GetProductById(id)
	if err != nil {
		if errors.Is(err, apperr.ErrProductNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "product not found",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}
	res := mapper.ProductEntityToDTO(product)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// ReduceProduct godoc
// @Summary      Уменьшить количество товара
// @Tags         products
// @Produce      json
// @Param        id     path     string  true  "ID товара"
// @Param        count  query    int     true  "Количество для вычитания"
// @Success      200    {object} dto.ProductResponse
// @Failure      400    {object} dto.Error400
// @Failure      404    {object} dto.Error404
// @Failure      500    {object} dto.Error500
// @Router       /product/{id} [patch]
func (p *ProductHandler) ReduceProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid ID",
		})
		return
	}

	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}
	if count < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid count",
		})
		return
	}
	product, err := p.product.ReduceProduct(id, count)
	if err != nil {
		if errors.Is(err, apperr.ErrProductNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "product not found",
			})
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}
	res := mapper.ProductEntityToDTO(product)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// GetProducts   godoc
// @Summary      Получить список товаров
// @Tags         products
// @Produce      json
// @Success      200  {array}  dto.ProductResponse
// @Success      400  {object}  dto.Error400
// @Success      404  {object}  dto.ProductsResponse
// @Failure      500  {object} dto.Error500
// @Router       /products [get]
func (p *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")
	products, err := p.product.GetProducts()
	if err != nil {
		if errors.Is(err, apperr.ErrProductNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ProductsNotFound{
				Products: make([]entity.Product, 0),
				Message:  "not found",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}
	res := mapper.ProductsEntityToDTOs(products)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// DeleteProduct godoc
// @Summary      Удалить товар
// @Tags         products
// @Param        id   path  string  true  "ID товара"
// @Success      200
// @Failure      400  {object} dto.Error400
// @Failure      404  {object} dto.Error404
// @Failure      500  {object} dto.Error500
// @Router       /product/{id} [delete]
func (p *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid ID",
		})
		return
	}
	err := p.product.DeleteProduct(id)
	if err != nil {
		if errors.Is(err, apperr.ErrProductNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "product not found",
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
	}
	w.WriteHeader(http.StatusOK)
}
