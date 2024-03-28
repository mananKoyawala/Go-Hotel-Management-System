package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
)

func AdminRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	{
		// only access to admin
		admin.GET("/login", controllers.AdminLogin())
		// admin.POST("/create", controllers.CreateAdmin()) // This route only for testing purposes
	}
}
