package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/middleware"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
)

func ManagerRoutes(r *gin.Engine) {
	r.Use(middleware.Authentication(models.Admin_Access))
	manager := r.Group("/manager")
	{
		// access to admin
		manager.GET("/getall") // controllers.LoginController()
		manager.GET("/get/:id")
		manager.POST("/create")
		manager.PUT("/update-all/:id")
		manager.PUT("/update-status/:id")
		manager.PUT("/update-branch/:id")
		manager.DELETE("/delete/:id")
	}
}

// PUT All the feilds or many fields
// PATCH particaly one or two fields
