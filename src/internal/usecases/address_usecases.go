package usecases

import "backend2/internal/entity"

type AddressRepo interface {
	Save(address entity.Address) (entity.Address, error)

	Update(address entity.Address) (entity.Address, error)
	Delete(address entity.Address) error
	GetById(address entity.Address) (entity.Address, error)
}

//TODO:
// 1) проверка существования
// 2) автоматическое добавление
