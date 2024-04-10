package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
	"github.com/mananKoyawala/hotel-management-system/pkg/middleware"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
)

func PickupServiceRoutes(r *gin.Engine) {
	service := r.Group("/service")
	{
		// * Manager, User, Driver
		service.GET("/getall", middleware.Authentication(models.M_Acc, models.G_Acc, models.D_Acc), controllers.GetAllPickUpServices())
		service.GET("/get/:id", middleware.Authentication(models.M_Acc, models.G_Acc, models.D_Acc), controllers.GetPickUpService())
		service.POST("/create", middleware.Authentication(models.M_Acc, models.G_Acc, models.D_Acc), controllers.CreatePickUpService())
		service.PUT("/update-details/:id", middleware.Authentication(models.M_Acc, models.G_Acc, models.D_Acc), controllers.UpdatePickUpServiceDetails())
		service.PATCH("/update-status/:id", middleware.Authentication(models.M_Acc, models.G_Acc, models.D_Acc), controllers.UpdatePickUpServiceStatus()) // update status of a service completed or not
		service.DELETE("/delete/:id", middleware.Authentication(models.M_Acc, models.G_Acc, models.D_Acc), controllers.DeletePickUpService())             // Cancel the pickup service
	}
}
