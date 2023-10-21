package authapi

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// LoginPayloadResponse login need data for login.html
type LoginPayloadResponse struct{}

// OAuthService service login URL
type OAuthService struct {
	Fqdn string `json:"fdqn"`
	URL  string `json:"url"`
}

// LoginPage login page
func LoginPage(c echo.Context) error {
	response := &LoginPayloadResponse{}

	return c.Render(http.StatusOK, "auth/login.html", response)
}
