package authentication

type Authentication struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	UserID       string `json:"userId"`
}

type SignUpData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type HashPasswordParams struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}
