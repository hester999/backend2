package mapper

import (
	"backend2/internal/dto"
	"backend2/internal/entity"
)

func ImgDTOToEntity(dto dto.ImageDTO) entity.Image {
	return entity.Image{
		Id:    dto.Id,
		Image: dto.Image,
	}
}

func ImgEntityToDTO(entity entity.Image) dto.ImageDTO {
	return dto.ImageDTO{
		Id:    entity.Id,
		Image: entity.Image,
	}
}
