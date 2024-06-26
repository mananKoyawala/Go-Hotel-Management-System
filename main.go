package main

import (
	"context"
	"os"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/database"
	"github.com/mananKoyawala/hotel-management-system/pkg/routes"
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

	// server.GET("/email", emailverification.SendMail())

	routes.AdminRoutes(server)
	routes.ManagerRoutes(server)
	routes.BranchRoutes(server)
	routes.ReservationRoutes(server)
	routes.PickupServiceRoutes(server)
	routes.DriverRoutes(server)
	routes.StaffRoutes(server)
	routes.FeedbackRoutes(server)
	routes.GuestRoutes(server)
	routes.RoomRoutes(server)
	color.Red("Please not that when you create the guest please make sure that u enter the correct email address otherwise problem occurs beacuse email can false.")
	server.Run(":" + PORT)
}
