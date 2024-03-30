package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/middleware"
)

func GuestRoutes(r *gin.Engine) {
	guest := r.Group("/guest")
	{
		// * User
		r.Use(middleware.Authentication())
		guest.POST("/signup")
		guest.POST("/login")
		guest.GET("/get/:id")
		guest.PUT("/update/:id")
		guest.DELETE("/delete/:id")
		// * Admin
		guest.GET("/getall")
	}
}
