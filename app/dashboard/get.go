package dashboard

import (
	"net/http"

	"github.com/highercomve/papelito/modules/auth/authservice"
	"github.com/labstack/echo/v4"
)

// GetDashboard get user dashboard
func GetDashboard(c echo.Context) error {
	_, err := authservice.GetOwnerID(c)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "index.html", nil)
}
