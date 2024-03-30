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
		// access to admin
		r.Use(middleware.Authentication(models.Admin_Access))
		manager.GET("/getall", controllers.GetManagers())
		manager.POST("/create", controllers.CreateManager())
		manager.PATCH("/update-status/:id", controllers.UpdateManagerStatus())
		r.Use(middleware.Authentication(models.Admin_Access, models.Manager_Access))
		manager.PUT("/update-all/:id", controllers.UpdateManagerDetails())
		manager.POST("/login", controllers.ManagerLoign())
		manager.DELETE("/delete/:id", controllers.DeleteManager())
		manager.GET("/get/:id", controllers.GetManager())
		manager.PATCH("/update-password", controllers.ResetManagerPassword())
	}
}

// PUT All the feilds or many fields
// PATCH particaly one or two fields
