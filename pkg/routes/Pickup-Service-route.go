package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
)

func PickupServiceRoutes(r *gin.Engine) {
	service := r.Group("/service")
	{
		// * Manager, User, Driver
		service.GET("/getall", controllers.GetAllPickUpServices())
		service.GET("/get/:id", controllers.GetPickUpService())
		service.POST("/create", controllers.CreatePickUpService())
		service.PUT("/update-details/:id", controllers.UpdatePickUpServiceDetails())
		service.PATCH("/update-status/:id", controllers.UpdatePickUpServiceStatus()) // update status of a service completed or not
		service.DELETE("/delete/:id", controllers.DeletePickUpService())             // Cancel the pickup service
	}
}
