package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
	"github.com/mananKoyawala/hotel-management-system/pkg/middleware"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
)

func DriverRoutes(r *gin.Engine) {
	driver := r.Group("/driver")
	{
		// * All
		driver.POST("/login", controllers.DriverLogin())
		driver.GET("/getall", middleware.Authentication(models.M_Acc, models.A_Acc, models.G_Acc), controllers.GetAllDrivers())
		driver.GET("/get/:id", middleware.Authentication(models.M_Acc, models.A_Acc, models.G_Acc), controllers.GetDriver())

		// * Admin
		driver.POST("/create", middleware.Authentication(models.A_Acc), controllers.CreateDriver())
		driver.PUT("/update-all/:id", middleware.Authentication(models.A_Acc), controllers.UpdateDriverDetails())
		driver.PATCH("/update-status/:id", middleware.Authentication(models.A_Acc), controllers.UpdateDriverStatus())
		driver.DELETE("/delete/:id", middleware.Authentication(models.A_Acc), controllers.DeleteDriver())

		// * Driver
		driver.PATCH("/update-profile-pic/:id", middleware.Authentication(models.D_Acc), controllers.UpdateDriverProfilePicture())
		driver.PATCH("/update-password", middleware.Authentication(models.D_Acc), controllers.ResetDriverPassword())
		driver.PATCH("/update-availability/:id", middleware.Authentication(models.D_Acc), controllers.UpdateDriverAvailability()) // status is changed based on reservation

		// * Admin , Manager
		driver.POST("/search", middleware.Authentication(models.M_Acc, models.A_Acc), controllers.SearchDriverData())                                      // search by first_name , last_name , gender, car_company, car_model and car_number_plate
		driver.POST("/filter", middleware.Authentication(models.M_Acc, models.A_Acc), middleware.Authentication(models.D_Acc), controllers.FilterDriver()) // filter by availability, state, salary

	}
}
