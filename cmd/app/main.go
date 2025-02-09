package main

import (
	"github.com/keremeti/iq-progers/config"
	_ "github.com/keremeti/iq-progers/docs"
	"github.com/keremeti/iq-progers/internal/app"
)

// @title           Swagger API
// @version         1.0
// @description     This is a sample server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host     localhost:8080
// @BasePath  /v1/

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// Configuration
	conf := config.New()

	app.Run(conf)
}
