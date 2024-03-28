package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/routes"
)

func main() {

	// GET PORT
	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8000"
	}

	server := gin.New()
	server.Use(gin.Logger())

	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	routes.AdminRoutes(server)

	server.Run(":" + PORT)
}
