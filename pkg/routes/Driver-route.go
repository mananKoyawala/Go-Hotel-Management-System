package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
)

func DriverRoutes(r *gin.Engine) {
	driver := r.Group("/driver")
	{
		// * All
		driver.GET("/getall", controllers.GetAllDrivers())
		driver.GET("/get/:id", controllers.GetDriver())
		driver.POST("/login", controllers.DriverLogin())
		// * Admin
		driver.POST("/create", controllers.CreateDriver())
		driver.PUT("/update-all/:id", controllers.UpdateDriverDetails())
		driver.PATCH("/update-status/:id", controllers.UpdateDriverStatus())
		driver.DELETE("/delete/:id", controllers.DeleteDriver())
		// * Driver
		driver.PATCH("/update-password", controllers.ResetDriverPassword())
		driver.PATCH("/update-availability/:id", controllers.UpdateDriverAvailability()) // status is changed based on reservation
	}
}
