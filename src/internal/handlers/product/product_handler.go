package product

import (
	"backend2/internal/apperr"
	"backend2/internal/entity"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type Product interface {
	CreateProduct(name, category, supplierId string, price float64, availableStock int, lastUpdate time.Time, img []byte) (entity.Product, error)
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

func (p *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product productCreateRequest

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, `{"status":400,"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(product)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"status":400,"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	res, err := p.product.CreateProduct(product.Name, product.Category, product.SupplierId, product.Price, product.AvailableStock, product.LastUpdate, product.Image)
	if err != nil {
		http.Error(w, `{"status":500,"error": internal server error"}`, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (p *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]

	if id == "" {
		http.Error(w, `{"status":400,"error":"invalid product id"}`, http.StatusBadRequest)
	}
	product, err := p.product.GetProductById(id)
	if err != nil {
		if errors.Is(err, apperr.ErrProductNotFound) {
			http.Error(w, `{"status":404,"error":"product not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"status":500,"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (p *ProductHandler) ReduceProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, `{"status":400,"error":"invalid product id"}`, http.StatusBadRequest)
		return
	}

	count, err := strconv.Atoi(r.URL.Query().Get("count"))
	if err != nil {
		http.Error(w, `{"status":500,"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	if count < 0 {
		http.Error(w, `{"status":400,"error":"invalid product count"}`, http.StatusBadRequest)
		return
	}
	res, err := p.product.ReduceProduct(id, count)
	if err != nil {
		if errors.Is(err, apperr.ErrProductNotFound) {
			http.Error(w, `{"status":404,"error":"product not found"}`, http.StatusNotFound)
		}
		http.Error(w, `{"status":500,"error":"internal server error"}`, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (p *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	products, err := p.product.GetProducts()
	if err != nil {
		http.Error(w, `{"status":500,"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (p *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, `{"status":400,"error":"invalid product id"}`, http.StatusBadRequest)
		return
	}
	err := p.product.DeleteProduct(id)
	if err != nil {
		if errors.Is(err, apperr.ErrProductNotFound) {
			http.Error(w, `{"status":404,"error":"product not found"}`, http.StatusNotFound)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
