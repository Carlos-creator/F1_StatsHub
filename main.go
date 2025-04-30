package main

import (
	"log"
	"f1statshub/endpoints"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	endpoints.SetupRoutes(router)

	log.Println("🚀 Server running on http://localhost:8080")
	router.Run(":8080")
}
