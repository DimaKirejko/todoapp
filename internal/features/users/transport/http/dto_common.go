package users_transport_http

import "github.com/DimaKirejko/todoapp/internal/core/domain"

type UsersDTOREsponse struct {
	ID          int     `json:"id"           example:"10"`
	Version     int     `json:"version"      example:"3"`
	FullName    string  `json:"full_name"    example:"Name Name"`
	PhoneNumber *string `json:"phone_number" example:"+380999999"`
}

func userDTOFromDomain(user domain.User) UsersDTOREsponse {
	return UsersDTOREsponse{
		ID:      user.ID,
		Version: user.Version,

		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func usersDTOFromDomains(users []domain.User) []UsersDTOREsponse {
	usersDTO := make([]UsersDTOREsponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}

	return usersDTO
}
