package supplier

import (
	"backend2/internal/apperr"
	"backend2/internal/dto"
	"backend2/internal/entity"
	"backend2/internal/mapper"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
)

type Supplier interface {
	CreateSupplier(supplier entity.Supplier) (entity.Supplier, error)
	GetSupplierById(id string) (entity.Supplier, error)
	GetAllSuppliers() ([]entity.Supplier, error)
	UpdateAddressSupplier(supplierId string, supplier entity.Supplier) (entity.Supplier, error)
	DeleteSupplierById(id string) error
}

type SupplierHandler struct {
	supplier Supplier
}

func NewSupplierHandler(supplier Supplier) *SupplierHandler {
	return &SupplierHandler{supplier}
}

// CreateSupplier godoc
// @Summary создать поставщика
// @Tags suppliers
// @Accept       json
// @Produce      json
// @Param supplier body dto.SupplierCreateRequestDTO true "Создаваемый поставщик"
// @Success 200 {object} dto.SupplierResponseDTO
// @Failure 400 {object} dto.Error400 "Bad request"
// @Failure 500 {object} dto.Error500 "Internal error"
// @Router  /supplier [post]
func (sh *SupplierHandler) CreateSupplier(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var supplier dto.SupplierCreateRequestDTO

	err := json.NewDecoder(r.Body).Decode(&supplier)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid JSON",
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(supplier)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.(validator.ValidationErrors).Error(),
		})
		return
	}

	supplierEntity := mapper.SupplierDTOToEntity(supplier)
	sup, err := sh.supplier.CreateSupplier(supplierEntity)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}
	res := mapper.SupplierEntityToDTO(sup)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// GetSupplierById godoc
// @Summary получить поставщика по id
// @Tags suppliers
// @Produce      json
// @Param       id   path  string  true  "ID поставщика"
// @Success 200 {object} dto.SupplierResponseDTO
// @Failure 400 {object} dto.Error400 "Bad request"
// @Failure 404 {object} dto.Error404 "Client not found"
// @Failure 500 {object} dto.Error500 "Internal error"
// @Router  /supplier/{id} [get]
func (sh *SupplierHandler) GetSupplierById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")

	id := mux.Vars(r)["id"]

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		})
		return
	}

	supplier, err := sh.supplier.GetSupplierById(id)
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
	res := mapper.SupplierEntityToDTO(supplier)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// GetAllSuppliers godoc
// @Summary получить всех поставщиков
// @Tags suppliers
// @Accept       json
// @Produce      json
// @Success 200 {object} dto.SuppliersResponse
// @Failure      400     {object} dto.Error400 "Bad request: invalid JSON or validation failed"
// @Failure      404      {object} dto.SuppliersNotFound
// @Failure      500     {object} dto.Error500 "internal error"
// @Router  /suppliers [get]
func (sh *SupplierHandler) GetAllSuppliers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")
	suppliers, err := sh.supplier.GetAllSuppliers()
	if err != nil {
		if errors.Is(err, apperr.ErrSupplierNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.SuppliersNotFound{
				Suppliers: make([]entity.Supplier, 0),
				Message:   "supplier not found",
			})
			return
		}
		return
	}

	res := mapper.SuppliersEntityToDTO(suppliers)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// UpdateAddress godoc
// @Summary обновить адрес поставщика по id
// @Tags suppliers
// @Produce      json
// @Param       id   path  string  true  "ID поставщика"
// @Param supplier body dto.SupplierUpdateAddressRequestDTO   true  "Обновляемые поля"
// @Success 200 {object} dto.SupplierResponseDTO
// @Failure      400     {object} dto.Error400 "Bad request: invalid JSON or validation failed"
// @Failure      404     {object} dto.Error404 "client not found"
// @Failure      500     {object} dto.Error500 "internal error"
// @Router  /supplier/{id} [patch]
func (sh *SupplierHandler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		})
		return
	}
	var supplier dto.SupplierUpdateAddressRequestDTO
	err := json.NewDecoder(r.Body).Decode(&supplier)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid JSON",
		})
		return
	}
	validate := validator.New()
	err = validate.Struct(supplier)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.(validator.ValidationErrors).Error(),
		})
		return
	}
	supplierEntity := mapper.SupplierUpdateDTOToEntity(supplier)

	supplierEntity, err = sh.supplier.UpdateAddressSupplier(id, supplierEntity)
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
	res := mapper.SupplierEntityToDTO(supplierEntity)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// DeleteSupplierById godoc
// @Summary удалить поставщика по id
// @Tags suppliers
// @Produce      json
// @Param       id   path  string  true  "ID поставщика"
// @Success 200
// @Failure      400     {object} dto.Error400 "Bad request: invalid JSON or validation failed"
// @Failure      404     {object} dto.Error404 "client not found"
// @Failure      500     {object} dto.Error500 "internal error"
// @Router  /supplier/{id} [delete]
func (sh *SupplierHandler) DeleteSupplierById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		})
		return
	}

	err := sh.supplier.DeleteSupplierById(id)
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
	w.WriteHeader(http.StatusOK)

}
