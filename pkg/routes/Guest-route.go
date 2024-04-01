package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
	"github.com/mananKoyawala/hotel-management-system/pkg/middleware"
)

func GuestRoutes(r *gin.Engine) {
	guest := r.Group("/guest")
	{
		// * User
		guest.POST("/signup", controllers.GuestSignup())
		guest.POST("/login", controllers.GuestLogin())
		r.Use(middleware.Authentication())
		guest.GET("/get/:id", controllers.GetGuest())
		guest.PUT("/update/:id", controllers.UpdateGuestDetails())
		guest.PATCH("/update-password", controllers.ResetUserPassword())
		guest.DELETE("/delete/:id", controllers.DeleteGuest())
		// * Admin
		guest.GET("/getall", controllers.GetAllGuest())
	}
}
