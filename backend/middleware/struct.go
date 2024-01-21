package middleware

type AuthenticationHeader struct {
	Authorization string `header:"Authentication"`
	ProjectId     string `header:"x-project-id"`
}

const (
	Create = "CREATE"
	Update = "UPDATE"
	Read   = "READ"
	Delete = "DELETE"
)
