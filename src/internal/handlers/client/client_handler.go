package client

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

type ClientUsecases interface {
	CreateClient(name, surname, gender string, regData, birthday time.Time, country, city, street string) (entity.Client, error)
	UpdateClient(id, country, city, street string) (entity.Client, error)
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

func (c *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var client clientCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, `{"status":400,"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(client); err != nil {
		http.Error(w, fmt.Sprintf(`{"status":400,"error":"%s"}`, err.(validator.ValidationErrors).Error()), http.StatusBadRequest)
		return
	}

	client.RegistrationDate = time.Now().UTC()
	result, err := c.client.CreateClient(
		client.ClientName,
		client.ClientSurname,
		client.Gender,
		client.RegistrationDate,
		client.BirthDate,
		client.Address.Country,
		client.Address.City,
		client.Address.Street,
	)

	if err != nil {
		http.Error(w, `{"status":500,"error":"failed to create client"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (c *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, `{"status":400,"error":"missing id"}`, http.StatusBadRequest)
	}

	var client clientUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		http.Error(w, `{"status":400,"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	validate := validator.New()
	if err := validate.Struct(client); err != nil {
		http.Error(w, fmt.Sprintf(`{"status":400,"error":"%s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	res, err := c.client.UpdateClient(id, client.Country, client.City, client.Street)
	if err != nil {
		if errors.Is(err, apperr.ErrClientNotFound) {
			http.Error(w, `{"status":404,"error":"client not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"status":500,"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (c *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	if id == "" {
		http.Error(w, `{"status":400,"error":"missing id"}`, http.StatusBadRequest)
		return
	}

	err := c.client.DeleteClient(id)
	if err != nil {
		http.Error(w, `{"status":500,"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (c *ClientHandler) GetAllClients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var limit, offset int
	var err error
	if r.URL.Query().Get("limit") != "" {
		limit, err = strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, `{"status":500,"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
	} else {
		limit = 0
	}

	if limit < 0 {
		http.Error(w, `{"status":400,"error":"limit cannot be negative"}`, http.StatusBadRequest)
		return
	}

	if r.URL.Query().Get("offset") != "" {
		offset, err = strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, `{"status":500,"error":"internal server error"}`, http.StatusInternalServerError)
			return
		}
	} else {
		offset = 0
	}

	if offset < 0 {
		http.Error(w, `{"status":400,"error":"offset cannot be negative"}`, http.StatusBadRequest)
		return
	}

	clients, err := c.client.GetAllClients(limit, offset)
	if err != nil {
		if errors.Is(err, apperr.ErrClientNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]any{
				"clients": clients,
				"error":   "no clients found",
			})
			return
		}
		http.Error(w, `{"status":500,"internal err"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(clients)

}

func (c *ClientHandler) GetClientByNameSurname(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.URL.Query().Get("name")
	surname := r.URL.Query().Get("surname")
	if name == "" {
		http.Error(w, `{"status":400,"error":"missing name"}`, http.StatusBadRequest)
		return
	}
	if surname == "" {
		http.Error(w, `{"status":400,"error":"missing surname"}`, http.StatusBadRequest)
		return
	}

	res, err := c.client.GetClientsByNameSurname(name, surname)
	if err != nil {
		if errors.Is(err, apperr.ErrClientNotFound) {
			http.Error(w, `{"status":404,"error":"clients not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"status":500,"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
