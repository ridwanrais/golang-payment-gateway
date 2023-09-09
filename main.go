package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ridwanrais/golang-payment-gateway/internal/route"
)

func main() {
	// Load environment variables from .env file
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file: " + err.Error())
	// }

	// constants.InitConstants()
	// config.ConnectToPostgreSQL()

	r := gin.Default()

	route.SetupRoutes(r)

	port := os.Getenv("PORT")
	r.Run(fmt.Sprintf(":%s", port))
}

//...
