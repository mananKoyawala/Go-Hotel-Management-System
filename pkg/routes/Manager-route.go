package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
	"github.com/mananKoyawala/hotel-management-system/pkg/middleware"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
)

func ManagerRoutes(r *gin.Engine) {
	manager := r.Group("/manager")
	{
		manager.POST("/login", controllers.ManagerLoign())

		// * Admin
		manager.GET("/getall", middleware.Authentication(models.A_Acc), controllers.GetManagers())
		manager.POST("/create", middleware.Authentication(models.A_Acc), controllers.CreateManager())
		manager.PATCH("/update-status/:id", middleware.Authentication(models.A_Acc), controllers.UpdateManagerStatus())
		manager.PUT("/update-all/:id", middleware.Authentication(models.A_Acc), controllers.UpdateManagerDetails())
		manager.DELETE("/delete/:id", middleware.Authentication(models.A_Acc), controllers.DeleteManager())
		manager.GET("/get/:id", middleware.Authentication(models.A_Acc), controllers.GetManager())
		manager.PATCH("/update-password", middleware.Authentication(models.A_Acc), controllers.ResetManagerPassword())
		manager.PATCH("/update-profile-pic/:id", middleware.Authentication(models.A_Acc), controllers.UpdateManagerProfilePicture())
		manager.POST("/search", middleware.Authentication(models.A_Acc), controllers.SearchManagerData()) // search by first_name , last_name and gender
		manager.POST("/filter", middleware.Authentication(models.A_Acc), controllers.FilterManager())     // filter by age, salary and status

	}
}

// PUT All the feilds or many fields
// PATCH particaly one or two fields
