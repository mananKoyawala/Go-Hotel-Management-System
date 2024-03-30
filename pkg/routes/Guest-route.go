package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/middleware"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
)

func GuestRoutes(r *gin.Engine) {
	guest := r.Group("/guest")
	{
		// unprotected beacuse anyone registered as guest
		guest.POST("/signup")
		guest.POST("/login")
		r.Use(middleware.Authentication(models.Admin_Access))
		guest.GET("/getall")
		guest.GET("/get/:id")
		r.Use(middleware.Authentication(models.Guest_Access))
		guest.PUT("/update/:id")
		guest.DELETE("/delete/:id")
	}
}
