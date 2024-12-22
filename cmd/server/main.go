package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/koodeyo/koodnet/pkg/api"
	"github.com/sirupsen/logrus"
)

func main() {
	godotenv.Load()
	l := logrus.New()

	l.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	lport := os.Getenv("KOODNET_LISTEN_PORT")
	laddress := os.Getenv("KOODNET_LISTEN_ADDRESS")
	koodnetEnv := os.Getenv("KOODNET_ENV") // Get the environment (e.g., "development", "production")

	if lport == "" {
		lport = "8001"
	}

	if laddress == "" {
		laddress = ":"
	}

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

	ctx := context.Background()

	r := api.NewRouter(l, &ctx)

	if err := r.Run(laddress + lport); err != nil {
		log.Fatal(err)
	}
}
