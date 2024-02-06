package authentication

import "github.com/graph-gophers/graphql-go"

type Authentication struct {
	Token        string     `json:"token"`
	RefreshToken string     `json:"refreshToken"`
	UserId       graphql.ID `json:"userId"`
}

type SignUpData struct {
	Username string `json:"username" validate:"min=8,max=32"`
	Password string `json:"password" validate:"min=8,max=32"`
}

type SignInData struct {
	Username string `json:"username" validate:"min=8,max=32"`
	Password string `json:"password" validate:"min=8,max=32"`
}

type SigninResponse struct {
	UserId       string `json:"userId"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenData struct {
	RefreshToken string `json:"refreshToken"`
}

type HashPasswordParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}
