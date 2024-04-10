package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
	"github.com/mananKoyawala/hotel-management-system/pkg/middleware"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
)

func GuestRoutes(r *gin.Engine) {
	guest := r.Group("/guest")
	{
		// * All
		guest.POST("/signup", controllers.GuestSignup())
		guest.POST("/login", controllers.GuestLogin())
		guest.GET("/verify-email/confirm", controllers.VerifyGuest())
		// * Only access to user
		guest.GET("/get/:id", middleware.Authentication(models.G_Acc), controllers.GetGuest())
		guest.PUT("/update/:id", middleware.Authentication(models.G_Acc), controllers.UpdateGuestDetails())
		guest.PATCH("/update-password", middleware.Authentication(models.G_Acc), controllers.ResetGuestPassword())
		guest.PATCH("/update-profile-pic/:id", middleware.Authentication(models.G_Acc), controllers.UpdateGuestProfilePicture())
		guest.DELETE("/delete/:id", middleware.Authentication(models.G_Acc), controllers.DeleteGuest())
		// * Admin
		guest.GET("/getall", middleware.Authentication(models.A_Acc), controllers.GetAllGuest())
	}
}
