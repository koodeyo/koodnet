package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/koodeyo/koodnet/pkg/api"
	"github.com/koodeyo/koodnet/pkg/database"
	"github.com/sirupsen/logrus"
)

func init() {
	// Load Env
	godotenv.Load()

	// Connect to database
	database.Connect()

	// Migrate
	database.Migrate()
}

func listenAddress() string {
	lport := os.Getenv("KOODNET_LISTEN_PORT")
	laddress := os.Getenv("KOODNET_LISTEN_ADDRESS")

	if lport == "" {
		lport = "8001"
	}

	if laddress == "" {
		laddress = ":"
	}

	return laddress + lport
}

// @title           Koodnet API
// @version         1.0
// @description     Server API documentation.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8001
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	l := logrus.New()

	l.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	koodnetEnv := os.Getenv("KOODNET_ENV") // Get the environment (e.g., "development", "production")

	// Dynamically set Gin mode based on KOODNET_ENV
	switch koodnetEnv {
	case "production":
		gin.SetMode(gin.ReleaseMode)
		l.SetLevel(logrus.InfoLevel)
		fmt.Println("Running in production mode")
	default:
		gin.SetMode(gin.DebugMode) // Default to debug mode
		l.SetLevel(logrus.DebugLevel)
	}

	r := api.NewRouter(l)

	if err := r.Run(listenAddress()); err != nil {
		log.Fatal(err)
	}
}
