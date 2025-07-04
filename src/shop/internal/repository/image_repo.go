package repository

import (
	"backend2/internal/apperr"
	"backend2/internal/entity"
	"database/sql"
	"errors"
	"fmt"
)

type ImageRepo struct {
	db *sql.DB
}

func NewImageRepo(db *sql.DB) *ImageRepo {
	return &ImageRepo{
		db: db,
	}
}

func (i *ImageRepo) AddImage(productID string, image entity.Image) (entity.Image, error) {
	tx, err := i.db.Begin()
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Проверка существования продукта
	var exists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM product WHERE id = $1)`, productID).Scan(&exists)
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to check product existence: %w", err)
	}
	if !exists {
		return entity.Image{}, fmt.Errorf("product  not found: %w", apperr.ErrProductNotFound)
	}

	// 2. Вставка изображения
	_, err = tx.Exec(`INSERT INTO images (id, image) VALUES ($1, $2)`, image.Id, image.Image)
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to insert image: %w", err)
	}

	// 3. Обновление поля image_id у продукта
	_, err = tx.Exec(`UPDATE product SET image_id = $1 WHERE id = $2`, image.Id, productID)
	if err != nil {
		return entity.Image{}, fmt.Errorf("failed to update product with image_id: %w", err)
	}

	// 4. Подтверждение транзакции
	if err = tx.Commit(); err != nil {
		return entity.Image{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return image, nil
}

func (i *ImageRepo) GetImageById(id string) (entity.Image, error) {
	query := `SELECT id, image FROM images WHERE id = $1`
	row := i.db.QueryRow(query, id)

	var img entity.Image
	err := row.Scan(&img.Id, &img.Image)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Image{}, apperr.ErrImageNotFound
	}
	if err != nil {
		return entity.Image{}, fmt.Errorf("error getting image: %w", err)
	}
	return img, nil
}

func (i *ImageRepo) UpdateImage(image entity.Image) (entity.Image, error) {
	query := `UPDATE images SET image = $1 WHERE id = $2`

	res, err := i.db.Exec(query, image.Image, image.Id)
	if err != nil {
		return entity.Image{}, fmt.Errorf("%w: %v", apperr.ErrImageUpdate, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return entity.Image{}, fmt.Errorf("checking update rows: %w", err)
	}
	if rowsAffected == 0 {
		return entity.Image{}, apperr.ErrImageNotFound
	}

	newImg, err := i.GetImageById(image.Id)
	if err != nil {
		return entity.Image{}, fmt.Errorf("error getting updated image: %w", err)
	}
	return newImg, nil
}

func (i *ImageRepo) GetProductImageById(productId string) (entity.Image, error) {
	query := `
		SELECT images.id, images.image
		FROM product
		JOIN images ON product.image_id = images.id
		WHERE product.id = $1
	`

	row := i.db.QueryRow(query, productId)

	var img entity.Image
	err := row.Scan(&img.Id, &img.Image)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Image{}, apperr.ErrImageNotFound
	}
	if err != nil {
		return entity.Image{}, fmt.Errorf("error getting product image: %w", err)
	}

	return img, nil
}

func (i *ImageRepo) DeleteImage(id string) error {
	query := `DELETE FROM images WHERE id = $1`

	res, err := i.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", apperr.ErrImageDelete, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking delete rows: %w", err)
	}
	if rowsAffected == 0 {
		return apperr.ErrImageNotFound
	}

	return nil
}
