package image

import (
	"backend2/internal/apperr"
	"backend2/internal/dto"
	"backend2/internal/entity"
	"backend2/internal/mapper"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type Image interface {
	AddImage(productID string, img []byte) (entity.Image, error)
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

// AddImage godoc
// @Summary      Загрузить изображение
// @Tags         images
// @Accept       multipart/form-data
// @Produce      json
// @Param        id     path     string  true  "ID продукта"
// @Param        image  formData file    true  "Файл изображения"
// @Success      200    {object} dto.ImageDTO
// @Failure      400    {object} dto.ErrorResponse
// @Failure      500    {object} dto.ErrorResponse
// @Router       /image/{id} [post]
func (i *ImageHandler) AddImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "id is required"})
		return
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "image upload error",
		})
		//http.Error(w, `{"status":500,"error":"image upload error"}`, http.StatusInternalServerError)
		return
	}
	defer file.Close()

	imgBytes, err := io.ReadAll(file)
	if err != nil {
		//http.Error(w, `{"status":500,"error":"failed to read file"}`, http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	img, err := i.img.AddImage(id, imgBytes)
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
	res := mapper.ImgEntityToDTO(img)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

}

// GetProductImageById godoc
// @Summary      Получить изображение товара
// @Tags         images
// @Produce      octet-stream
// @Param        id   path  string  true  "ID продукта"
// @Success      200  {file} binary
// @Failure      404  {object} dto.ErrorResponse
// @Router       /products/{id}/image [get]
func (i *ImageHandler) GetProductImageById(w http.ResponseWriter, r *http.Request) {
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

	productImage, err := i.img.GetProductImageById(id)
	if err != nil {
		if errors.Is(err, apperr.ErrImageNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "image not found",
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

	img := productImage.Image

	w.Header().Set("Content-Type", "application/octet-stream,image/jpeg,image/png")
	//w.Header().Set("Content-Disposition", `attachment; filename="image.jpg"`)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(img)))

	w.WriteHeader(http.StatusOK)
	w.Write(img)
}

// GetImageById godoc
// @Summary      Получить изображение по ID
// @Tags         images
// @Produce      octet-stream
// @Param        id   path  string  true  "ID изображения"
// @Success      200  {file} binary
// @Failure      404  {object} dto.ErrorResponse
// @Router       /image/{id} [get]
func (i *ImageHandler) GetImageById(w http.ResponseWriter, r *http.Request) {
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

	productImage, err := i.img.GetImageById(id)
	if err != nil {
		if errors.Is(err, apperr.ErrImageNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "image not found",
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

	img := productImage.Image

	w.Header().Set("Content-Type", "application/octet-stream,image/jpeg,image/png")
	//w.Header().Set("Content-Disposition", `attachment; filename="image.jpg"`)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(img)))

	w.WriteHeader(http.StatusOK)
	w.Write(img)
}

// DeleteImage godoc
// @Summary      Удалить изображение
// @Tags         images
// @Param        id   path  string  true  "ID изображения"
// @Success      200
// @Failure      404  {object} dto.ErrorResponse
// @Router       /image/{id} [delete]
func (i *ImageHandler) DeleteImage(w http.ResponseWriter, r *http.Request) {
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
	err := i.img.DeleteImage(id)
	if err != nil {
		if errors.Is(err, apperr.ErrImageNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "image not found",
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

// UpdateImage godoc
// @Summary      Обновить изображение
// @Tags         images
// @Accept       multipart/form-data
// @Produce      octet-stream
// @Param        id     path     string  true  "ID изображения"
// @Param        image  formData file    true  "Новый файл"
// @Success      200    {file}   binary
// @Failure      400    {object} dto.ErrorResponse
// @Failure      500    {object} dto.ErrorResponse
// @Router       /image/{id} [patch]
func (i *ImageHandler) UpdateImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream,image/jpeg,image/png,application/json")
	id := mux.Vars(r)["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid id",
		})
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	defer file.Close()
	imgBytes, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	img, err := i.img.UpdateImage(id, imgBytes)
	if err != nil {
		if errors.Is(err, apperr.ErrImageNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "image not found",
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
	res := mapper.ImgEntityToDTO(img)
	w.WriteHeader(http.StatusOK)
	w.Write(res.Image)
}
