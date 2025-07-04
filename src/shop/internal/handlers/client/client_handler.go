package client

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
	"strconv"
)

type ClientUsecases interface {
	CreateClient(client entity.Client) (entity.Client, error)
	UpdateClient(id string, client entity.Client) (entity.Client, error)
	DeleteClient(id string) error
	GetAllClients(limit, offset int) ([]entity.Client, error)
	GetClientsByNameSurname(name, surname string) ([]entity.Client, error)
}

type ClientHandler struct {
	client ClientUsecases
}

func NewClientHandler(client ClientUsecases) *ClientHandler {
	return &ClientHandler{
		client: client,
	}
}

// CreateClient godoc
// @Summary      Создать клиента
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        client  body     dto.ClientCreateRequestDTO  true  "Создаваемый клиент"
// @Success      201     {object} dto.ClientResponseDTO
// @Failure      400     {object} dto.Error400
// @Failure      500     {object} dto.Error500
// @Router       /client [post]
func (c *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var client dto.ClientCreateRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Message: "invalid JSON",
			Code:    http.StatusBadRequest,
		})
		//http.Error(w, `{"status":400,"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(client); err != nil {
		//http.Error(w, fmt.Sprintf(`{"status":400,"error":"%s"}`, err.(validator.ValidationErrors).Error()), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Message: err.(validator.ValidationErrors).Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	entityClient := mapper.ClientCreateRequestToEntity(client)
	entityClient, err := c.client.CreateClient(entityClient)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Message: "internal server error",
			Code:    http.StatusInternalServerError,
		})
		//http.Error(w, `{"status":500,"error":"failed to create client"}`, http.StatusInternalServerError)
		return
	}
	result := mapper.ClientCreateResponse(entityClient)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

// UpdateClient godoc
// @Summary      Обновить клиента
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        id      path     string  true  "ID клиента"
// @Param        client  body     dto.ClientUpdateRequestDTO  true  "Обновляемые поля"
// @Success 200 {object} dto.ClientResponseDTO "OK"
// @Example 200 {json} ClientResponseExample
//
//	{
//	  "id": "f19a3a7-12f5-4332-9582-624519c3eaea",
//	  "client_name": "Harry",
//	  "client_sure_name": "Potter",
//	  "birth_date": "2000-07-31T00:00:00Z",
//	  "gender": "male",
//	  "register_date": "2020-09-01T12:00:00Z",
//	  "address_id": "a123b456-c789-d012-e345-67890abcdef1",
//	  "address": {
//	    "id": "a123b456-c789-d012-e345-67890abcdef1",
//	    "country": "UK",
//	    "city": "London",
//	    "street": "Privet Drive"
//	  }
//	}
//
// @Failure      400     {object} dto.Error400 "Bad request: invalid JSON or validation failed
// @Failure      404     {object} dto.Error404 "client not found"
// @Failure      500     {object} dto.Error500 "internal error"
// @Router       /client/{id} [patch]
func (c *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Message: "invalid id",
			Code:    http.StatusBadRequest,
		})
		//http.Error(w, `{"status":400,"error":"missing id"}`, http.StatusBadRequest)
		return
	}

	var client dto.ClientUpdateRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Message: "invalid JSON",
			Code:    http.StatusBadRequest,
		})
		//http.Error(w, `{"status":400,"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	validate := validator.New()

	if err := validate.Struct(client); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Message: err.(validator.ValidationErrors).Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	entityClient := mapper.ClientUpdateRequestToEntity(client)
	entityClient, err := c.client.UpdateClient(id, entityClient)
	if err != nil {
		if errors.Is(err, apperr.ErrClientNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Message: "client not found",
				Code:    http.StatusNotFound,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Message: "internal server error",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	res := mapper.ClientUpdateResponse(entityClient)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// DeleteClient godoc
// @Summary      Удалить клиента
// @Tags         clients
// @Param        id   path  string  true  "ID клиента"
// @Success      200
// @Failure 400 {object} dto.Error400 "Bad request"
// @Failure 404 {object} dto.Error404 "Client not found"
// @Failure 500 {object} dto.Error500 "Internal error"
// @Router       /client/{id} [delete]
func (c *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	if id == "" {
		//http.Error(w, `{"status":400,"error":"missing id"}`, http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Message: "invalid id",
			Code:    http.StatusBadRequest,
		})
		return
	}

	err := c.client.DeleteClient(id)
	if err != nil {
		if errors.Is(err, apperr.ErrClientNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Message: "client not found",
				Code:    http.StatusNotFound,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Message: "internal server error",
			Code:    http.StatusInternalServerError,
		})
		return
	}
	w.WriteHeader(http.StatusOK)

}

// GetAllClients godoc
// @Summary      Получить всех клиентов
// @Tags         clients
// @Produce      json
// @Success      200  {array}  dto.ClientResponseDTO
// @Success      400  {object}  dto.Error400
// @Failure      404  {object} dto.ClientsNotFound
// @Failure 	 500 {object} dto.Error500 "Internal error"
// @Param 		 limit  query string false "количество отоброжаемых клиентов"
// @Param 		 offset  query string false "Смещение выборки"
// @Router       /clients [get]
func (c *ClientHandler) GetAllClients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")

	var limit, offset int
	var err error
	if r.URL.Query().Get("limit") != "" {
		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Message: "internal server error",
				Code:    http.StatusInternalServerError,
			})
			return
		}
	} else {
		limit = 0
	}

	if limit < 0 {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Message: "limit cannot be negative",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if r.URL.Query().Get("offset") != "" {
		offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Message: "internal server error",
				Code:    http.StatusInternalServerError,
			})
			return
		}
	} else {
		offset = 0
	}

	if offset < 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Message: "offset cannot be negative",
			Code:    http.StatusBadRequest,
		})
		return
	}

	clients, err := c.client.GetAllClients(limit, offset)
	if err != nil {
		if errors.Is(err, apperr.ErrClientNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ClientsNotFound{
				Clients: make([]entity.Client, 0),
				Message: "clients not found"})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Message: "internal server error",
			Code:    http.StatusInternalServerError,
		})

		return
	}
	result := mapper.GetClientsResponse(clients)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)

}

// GetClientsByNameSurname godoc
// @Summary      Поиск клиента по имени и фамилии
// @Tags         clients
// @Produce      json
// @Param        name     query    string  true  "Имя"
// @Param        surname  query    string  true  "Фамилия"
// @Success      200      {array}  dto.ClientResponseDTO
// @Failure 400 {object} dto.Error400 "Bad request"
// @Failure      404      {object} dto.ClientsNotFound "not found"
// @Example      404      {json} NotFoundExample
//
//	{
//
//	  "clients": [],
//	  "message": "clients not found"
//	}
//
// @Failure 500 {object} dto.Error500 "Internal error"
// @Router       /client [get]
func (c *ClientHandler) GetClientsByNameSurname(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=60")
	name := r.URL.Query().Get("name")
	surname := r.URL.Query().Get("surname")
	if name == "" {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Message: "missing name",
			Code:    http.StatusBadRequest,
		})
		return
	}
	if surname == "" {

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Message: "missing surname",
			Code:    http.StatusBadRequest,
		})
		return
	}

	clients, err := c.client.GetClientsByNameSurname(name, surname)
	if err != nil {
		if errors.Is(err, apperr.ErrClientNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ClientsNotFound{
				Clients: make([]entity.Client, 0),
				Message: "clients not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Message: "internal server error",
			Code:    http.StatusInternalServerError,
		})
		return
	}
	result := mapper.GetClientsResponse(clients)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
