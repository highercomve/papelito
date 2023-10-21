package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/highercomve/papelito/modules/helper/helperrepo"
	"github.com/highercomve/papelito/server"
	"github.com/highercomve/papelito/utils"
	"github.com/highercomve/papelito/utils/tracer"
)

func init() {
	utils.SetupEnvDefaults()

	buildEnvironmentJS("assets")

	// Start connection to database
	_, err := helperrepo.NewStorage(utils.Env.CollectionPrefix)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// @title Papelito the game
// @version 1.0
// @description Papelito the game
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /

func main() {
	if utils.GetEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "") != "" {
		tp := tracer.Init(utils.GetEnv("OTEL_SERVICE_NAME", "papelito"))
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}()
	}

	if utils.GetEnv("DISABLE_API", "") != "true" {
		serverAddres := fmt.Sprintf("%s:%s", utils.Env.ServerHost, utils.Env.ServerPort)
		server.Start(serverAddres)
	}
}

func buildEnvironmentJS(folder string) error {
	var environment = "if (!window.env) { window.env = {} } \n"
	for _, e := range os.Environ() {
		if strings.Contains(e, "REACT_APP_") {
			pair := strings.SplitN(e, "=", 2)
			environment = environment + "window.env." + pair[0] + "='" + pair[1] + "';"
		}
	}

	return os.WriteFile(folder+"/environment.js", []byte(environment), 0644)
}
