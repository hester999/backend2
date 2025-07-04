package supplier

import (
	"backend2/internal/apperr"
	"backend2/internal/entity"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
)

type Supplier interface {
	CreateSupplier(name, phoneNumber, country, city, street string) (entity.Supplier, error)
	GetSupplierById(id string) (entity.Supplier, error)
	GetAllSuppliers() ([]entity.Supplier, error)
	UpdateAddressSupplier(supplierId, country, city, street string) (entity.Supplier, error)
	DeleteSupplierById(id string) error
}

type SupplierHandler struct {
	supplier Supplier
}

func NewSupplierHandler(supplier Supplier) *SupplierHandler {
	return &SupplierHandler{supplier}
}

func (sh *SupplierHandler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var supplier supplierCreateRequest

	err := json.NewDecoder(r.Body).Decode(&supplier)
	if err != nil {
		http.Error(w, `{"status":400,"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	validate := validator.New()
	err = validate.Struct(supplier)
	if err != nil {
		http.Error(w, `{"status":500,"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	res, err := sh.supplier.CreateSupplier(supplier.Name, supplier.PhoneNumber, supplier.Country, supplier.City, supplier.Street)
	if err != nil {
		http.Error(w, `{"status":500,"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (sh *SupplierHandler) GetSupplierById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	if id == "" {
		http.Error(w, `{"status":400,"error":"invalid ID"}`, http.StatusBadRequest)
		return
	}

	supplier, err := sh.supplier.GetSupplierById(id)
	if err != nil {
		if errors.Is(err, apperr.ErrSupplierNotFound) {
			http.Error(w, `{"status":404,"error":"supplier not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"status":500,"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(supplier)
}

func (sh *SupplierHandler) GetAllSuppliers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	suppliers, err := sh.supplier.GetAllSuppliers()
	if err != nil {
		http.Error(w, `{"status":500,"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(suppliers)
}

func (sh *SupplierHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, `{"status":400,"error":"invalid ID"}`, http.StatusBadRequest)
		return
	}
	var supplier supplierUpdateAddressRequest
	err := json.NewDecoder(r.Body).Decode(&supplier)
	if err != nil {
		http.Error(w, `{"status":400,"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	validate := validator.New()
	err = validate.Struct(supplier)
	if err != nil {
		http.Error(w, `{"status":500,"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	res, err := sh.supplier.UpdateAddressSupplier(id, supplier.Street, supplier.Country, supplier.City)
	if err != nil {
		http.Error(w, `{"status":500,"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (sh *SupplierHandler) DeleteSupplierById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, `{"status":400,"error":"invalid ID"}`, http.StatusBadRequest)
		return
	}
	err := sh.supplier.DeleteSupplierById(id)
	if err != nil {
		if errors.Is(err, apperr.ErrSupplierNotFound) {
			http.Error(w, `{"status":404,"error":"supplier not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"status":500,"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
