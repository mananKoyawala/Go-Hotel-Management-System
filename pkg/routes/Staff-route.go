package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
)

func StaffRoutes(r *gin.Engine) {
	staff := r.Group("/staff")
	{
		// * Manager
		staff.GET("/getall/:id", controllers.GetAllStaff()) // by branch id
		staff.GET("/get/:id", controllers.GetStaff())
		staff.POST("/create", controllers.CreateStaff())
		staff.PUT("/update-all/:id", controllers.UpdateStaffDetails())
		staff.PATCH("/update-status/:id", controllers.UpdateStaffStatus())
		staff.PATCH("/update-profile-pic/:id", controllers.UpdateStaffProfilePicture())
		staff.DELETE("/delete/:id", controllers.DeleteStaff())
	}
}
