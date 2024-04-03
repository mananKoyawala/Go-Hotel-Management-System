package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/database"
	"github.com/mananKoyawala/hotel-management-system/pkg/helpers/image"
	"github.com/mananKoyawala/hotel-management-system/pkg/routes"
	_ "github.com/mananKoyawala/hotel-management-system/pkg/utils"
)

func main() {

	client := database.DBInstance()
	defer client.Disconnect(context.Background())

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

	server.POST("/upload", image.ImageUpload())
	server.DELETE("/delete", image.ImageDelete())
	routes.ReservationRoutes(server)
	routes.PickupServiceRoutes(server)
	routes.DriverRoutes(server)
	routes.StaffRoutes(server)
	routes.FeedbackRoutes(server)
	routes.GuestRoutes(server)
	routes.RoomRoutes(server)
	routes.AdminRoutes(server)
	routes.ManagerRoutes(server)
	routes.BranchRoutes(server)

	server.Run(":" + PORT)
}
