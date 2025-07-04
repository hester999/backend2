package image

import (
	"backend2/internal/apperr"
	"backend2/internal/entity"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type Image interface {
	AddImage(img []byte) (entity.Image, error)
	GetImageById(id string) (entity.Image, error)
	UpdateImage(id string, img []byte) (entity.Image, error)
	GetProductImageById(productId string) (entity.Image, error)
	DeleteImage(id string) error
}

type ImageHandler struct {
	img Image
}

func NewImageHandler(img Image) *ImageHandler {
	return &ImageHandler{
		img: img,
	}
}

func (i *ImageHandler) AddImage(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/octet-stream,multipart/form-data")

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, `{"status":500,"error":"image upload error"}`, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	imgBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, `{"status":500,"error":"failed to read file"}`, http.StatusInternalServerError)
		return
	}

	res, err := i.img.AddImage(imgBytes)
	if err != nil {
		http.Error(w, `{"status":500,"error":"failed to add image"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

}

func (i *ImageHandler) GetProductImageById(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, `{"status":400,"error":"missing id"}`, http.StatusBadRequest)
		return
	}

	productImage, err := i.img.GetProductImageById(id)
	if err != nil {
		http.Error(w, `{"status":404,"error":"product image not found"}`, http.StatusNotFound)
		return
	}

	img := productImage.Image

	w.Header().Set("Content-Type", "application/octet-stream,image/jpeg,image/png")
	//w.Header().Set("Content-Disposition", `attachment; filename="image.jpg"`)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(img)))

	w.WriteHeader(http.StatusOK)
	w.Write(img)
}

func (i *ImageHandler) GetImageById(w http.ResponseWriter, r *http.Request) {

	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, `{"status":400,"error":"missing id"}`, http.StatusBadRequest)
		return
	}

	productImage, err := i.img.GetImageById(id)
	if err != nil {
		http.Error(w, `{"status":404,"error":"image not found"}`, http.StatusNotFound)
		return
	}

	img := productImage.Image

	w.Header().Set("Content-Type", "application/octet-stream,image/jpeg,image/png")
	//w.Header().Set("Content-Disposition", `attachment; filename="image.jpg"`)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(img)))

	w.WriteHeader(http.StatusOK)
	w.Write(img)
}

func (i *ImageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, `{"status":400,"error":"missing id"}`, http.StatusBadRequest)
		return
	}
	err := i.img.DeleteImage(id)
	if err != nil {
		if errors.Is(err, apperr.ErrImageNotFound) {
			http.Error(w, `{"status":404,"error":"image not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"status":500,"error":"failed to delete image"}`, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (i *ImageHandler) UpdateImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream,image/jpeg,image/png")
	id := mux.Vars(r)["id"]
	if id == "" {
		http.Error(w, `{"status":400,"error":"missing id"}`, http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, `{"status":500,"error":"image upload error"}`, http.StatusInternalServerError)
		return
	}
	defer file.Close()
	imgBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, `{"status":500,"error":"failed to read file"}`, http.StatusInternalServerError)
		return
	}
	res, err := i.img.UpdateImage(id, imgBytes)
	if err != nil {
		if errors.Is(err, apperr.ErrImageNotFound) {
			http.Error(w, `{"status":404,"error":"image not found"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"status":500,"error":"failed to update image 1"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res.Image)
}
