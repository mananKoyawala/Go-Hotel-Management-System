package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
	"github.com/mananKoyawala/hotel-management-system/pkg/middleware"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
)

func ManagerRoutes(r *gin.Engine) {
	r.Use(middleware.Authentication(models.Admin_Access))
	manager := r.Group("/manager")
	{
		// access to admin
		manager.GET("/getall", controllers.GetManagers())
		manager.GET("/get/:id", controllers.GetManager())
		manager.POST("/login", controllers.ManagerLoign())
		manager.POST("/create", controllers.CreateManager())
		manager.PUT("/update-all/:id", controllers.UpdateManagerDetails())
		manager.PUT("/update-status/:id", controllers.UpdateManagerStatus())
		manager.DELETE("/delete/:id", controllers.DeleteManager())
		// TODO : Password update
	}
}

// PUT All the feilds or many fields
// PATCH particaly one or two fields
