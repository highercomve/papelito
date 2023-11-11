package server

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"

	docs "github.com/highercomve/papelito/docs"
	"github.com/highercomve/papelito/modules/api"
	"github.com/highercomve/papelito/modules/app"
	"github.com/highercomve/papelito/modules/app/dashboard"
	"github.com/highercomve/papelito/modules/helpers"
	"github.com/highercomve/papelito/modules/middlewares"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

// CustomValidator payload validation
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate payload
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// Start server
func Start(serverAddress string) {
	docs.SwaggerInfo.BasePath = "/api/v1/"
	docs.SwaggerInfo.Schemes = []string{"https"}
	docs.SwaggerInfo.Host = helpers.Env.HostURL

	e := echo.New()

	// Pre request middlewares
	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(middlewares.MethodOverride)
	e.Pre(middlewares.ConvertJaegerToOtel)

	// Serve static assets
	// e.Static("/assets", "assets")

	e.Renderer = CreateTemplateRenderer()
	e.Validator = &CustomValidator{
		Validator: validator.New(),
	}
	cookieSecret := []byte(helpers.Env.SessionSecret)
	cookieStore := sessions.NewCookieStore(cookieSecret)
	cookieStore.Options.HttpOnly = true
	cookieStore.Options.Secure = true

	corsConfig := middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
	}

	// Set Middlewares
	e.Use(otelecho.Middleware(helpers.GetEnv("OTEL_SERVICE_NAME", "papelito")))
	e.Use(session.Middleware(cookieStore))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORSWithConfig(corsConfig))

	if helpers.GetEnv("DEBUG", "") == "true" {
		e.Debug = true
	}

	// Metrics
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	// Load all api services
	api.LoadAPI(e)
	app.LoadApp(e)

	e.GET("/assets/*", GetStatic)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/", dashboard.GetDashboard, middlewares.CookieAuthentication)

	e.Logger.Fatal(e.Start(serverAddress))
}

func GetStatic(c echo.Context) error {
	if len(c.ParamValues()) < 1 {
		return echo.ErrNotFound
	}

	filePath := "assets/" + c.ParamValues()[0]
	acceptEncoding := strings.Split(c.Request().Header.Get("Accept-Encoding"), ",")
	contentType := mimeTypeFromFilename(filePath)

	c.Response().Header().Set("Cache-Control", "private, max-age=86400, stale-while-revalidate=604800")

	for _, encoding := range acceptEncoding {
		var extension = ""
		encoding := strings.Trim(encoding, " ")
		switch encoding {
		case "br":
			extension = ".br"
		case "gzip":
			extension = ".gz"
		default:
			continue
		}
		_, err := os.Stat(filePath + extension)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}

		content, err := os.Open(filePath + extension)
		if err != nil {
			continue
		}

		c.Response().Header().Set("Content-Encoding", encoding)
		return c.Stream(200, contentType, content)
	}

	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return echo.ErrNotFound
	}

	content, err := os.Open(filePath)
	if err != nil {
		return echo.ErrNotFound
	}

	return c.Stream(200, contentType, content)
}

func mimeTypeFromFilename(filename string) string {
	ext := filepath.Ext(filename)

	switch ext {
	case ".css":
		return "text/css"
	case ".jpg":
		return "image/jpeg"
	case ".jpeg":
		return "image/jpeg"
	case ".js":
		return echo.MIMEApplicationJavaScript
	case ".json":
		return echo.MIMEApplicationJSONCharsetUTF8
	case ".png":
		return "image/png"
	case ".svg":
		return "image/svg+xml"
	case ".html":
		return echo.MIMETextHTMLCharsetUTF8
	default:
		return echo.MIMEOctetStream
	}
}
