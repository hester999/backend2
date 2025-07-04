package mapper

import (
	"backend2/internal/dto"
	"backend2/internal/entity"
)

func ClientCreateRequestToEntity(request dto.ClientCreateRequestDTO) entity.Client {
	address := entity.Address{

		Country: request.Address.Country,
		City:    request.Address.City,
		Street:  request.Address.Street,
	}
	return entity.Client{

		ClientName:    request.ClientName,
		ClientSurname: request.ClientSurname,
		BirthDate:     request.BirthDate,
		Gender:        request.Gender,

		Address: address,
	}
}

func ClientCreateResponse(client entity.Client) dto.ClientResponseDTO {
	return dto.ClientResponseDTO{
		Id:               client.Id,
		ClientName:       client.ClientName,
		ClientSurname:    client.ClientSurname,
		BirthDate:        client.BirthDate,
		Gender:           client.Gender,
		RegistrationDate: client.RegistrationDate,
		AddressId:        client.AddressId,
		Address: dto.AddressDTO{
			ID:      client.Address.ID,
			Country: client.Address.Country,
			City:    client.Address.City,
			Street:  client.Address.Street,
		},
	}
}

func ClientUpdateRequestToEntity(request dto.ClientUpdateRequestDTO) entity.Client {
	address := entity.Address{
		Country: request.Country,
		City:    request.City,
		Street:  request.Street,
	}
	return entity.Client{
		Address: address,
	}
}

func ClientUpdateResponse(client entity.Client) dto.ClientResponseDTO {

	return dto.ClientResponseDTO{
		Id:            client.Id,
		ClientName:    client.ClientName,
		ClientSurname: client.ClientSurname,
		BirthDate:     client.BirthDate,
		Gender:        client.Gender,
		AddressId:     client.AddressId,

		Address: dto.AddressDTO{
			ID:      client.Address.ID,
			Country: client.Address.Country,
			City:    client.Address.City,
			Street:  client.Address.Street,
		},
	}
}

func clientGetResponse(client entity.Client) dto.ClientResponseDTO {
	return dto.ClientResponseDTO{
		Id:            client.Id,
		ClientName:    client.ClientName,
		ClientSurname: client.ClientSurname,
		BirthDate:     client.BirthDate,
		Gender:        client.Gender,
		AddressId:     client.AddressId,
		Address: dto.AddressDTO{
			ID:      client.Address.ID,
			Country: client.Address.Country,
			City:    client.Address.City,
			Street:  client.Address.Street,
		},
	}
}

func GetClientsResponse(clients []entity.Client) dto.ClientsResponseDTO {

	clientsResponse := dto.ClientsResponseDTO{
		Clients: make([]dto.ClientResponseDTO, 0, len(clients)),
	}
	for _, v := range clients {
		tmp := clientGetResponse(v)
		clientsResponse.Clients = append(clientsResponse.Clients, tmp)
	}
	return clientsResponse
}
