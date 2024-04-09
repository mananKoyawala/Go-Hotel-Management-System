package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
)

func ManagerRoutes(r *gin.Engine) {
	manager := r.Group("/manager")
	{
		// * Admin
		// r.Use(middleware.Authentication())
		manager.GET("/getall", controllers.GetManagers())
		manager.POST("/create", controllers.CreateManager())
		manager.PATCH("/update-status/:id", controllers.UpdateManagerStatus())
		manager.PUT("/update-all/:id", controllers.UpdateManagerDetails())
		manager.POST("/login", controllers.ManagerLoign())
		manager.DELETE("/delete/:id", controllers.DeleteManager())
		manager.GET("/get/:id", controllers.GetManager())
		manager.PATCH("/update-password", controllers.ResetManagerPassword())
		manager.PATCH("/update-profile-pic/:id", controllers.UpdateManagerProfilePicture())
		manager.POST("/search", controllers.SearchManagerData()) // search by first_name , last_name and gender

	}
}

// PUT All the feilds or many fields
// PATCH particaly one or two fields
