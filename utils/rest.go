package utils

import "github.com/labstack/echo/v4"

// NotImplemented rest handler not implemented
func NotImplemented(c echo.Context) error {
	return nil
}

// RestHTTPError rest http
type RestHTTPError struct {
	Code     int         `json:"-"`
	Message  interface{} `json:"message"`
	Internal error       `json:"-"` // Stores the error returned by an external dependency
}
