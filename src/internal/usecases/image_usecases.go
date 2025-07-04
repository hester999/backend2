package usecases

import (
	"backend2/internal/entity"
	"backend2/internal/utils"
	"fmt"
)

//добавление изображения (на вход подается byte array изображения и id товара).
//
//Изменение изображения (на вход подается id изображения и новая строка для замены)
//
//Удаление изображения по id изображения
//
//Получение изображения конкретного товара (по id товара)
//
//Получение изображения по id изображения

type ImageRepo interface {
	AddImage(image entity.Image) (entity.Image, error)
	GetImageById(id string) (entity.Image, error)
	UpdateImage(image entity.Image) (entity.Image, error)
	GetProductImageById(productId string) (entity.Image, error)
	DeleteImage(id string) error
}

type Image struct {
	img ImageRepo
}

func NewImage(img ImageRepo) *Image {
	return &Image{img: img}
}

func (i *Image) AddImage(img []byte) (entity.Image, error) {

	id, err := utils.GenerateUUID()

	if err != nil {
		return entity.Image{}, fmt.Errorf("error generating id: %w", err)
	}

	newImg := entity.Image{
		Id:    id,
		Image: img,
	}
	newImg, err = i.img.AddImage(newImg)
	if err != nil {
		return entity.Image{}, fmt.Errorf("error adding image: %w", err)
	}
	return newImg, nil
}

func (i *Image) GetImageById(id string) (entity.Image, error) {
	img, err := i.img.GetImageById(id)
	if err != nil {
		return entity.Image{}, fmt.Errorf("error getting image: %w", err)
	}
	return img, nil
}

func (i *Image) UpdateImage(id string, img []byte) (entity.Image, error) {

	newImg := entity.Image{
		Id:    id,
		Image: img,
	}
	newImg, err := i.img.UpdateImage(newImg)
	if err != nil {
		return entity.Image{}, fmt.Errorf("error updating image: %w", err)
	}
	return newImg, nil
}

func (i *Image) GetProductImageById(productId string) (entity.Image, error) {
	img, err := i.img.GetProductImageById(productId)
	if err != nil {
		return entity.Image{}, fmt.Errorf("error getting product image: %w", err)
	}
	return img, nil
}

func (i *Image) DeleteImage(id string) error {
	err := i.img.DeleteImage(id)
	if err != nil {
		return fmt.Errorf("error deleting image: %w", err)
	}
	return nil
}
