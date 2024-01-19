package middleware

type AuthenticationHeader struct {
	Authorization string `header:"Authentication"`
	ProjectId     string `header:"x-project-id"`
}
