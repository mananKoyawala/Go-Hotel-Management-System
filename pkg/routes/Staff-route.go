package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
	"github.com/mananKoyawala/hotel-management-system/pkg/middleware"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
)

func StaffRoutes(r *gin.Engine) {
	staff := r.Group("/staff")
	{
		// * Manager
		staff.GET("/getall/:id", middleware.Authentication(models.M_Acc), controllers.GetAllStaff()) // by branch id
		staff.GET("/getall", middleware.Authentication(models.M_Acc), controllers.GetStaffs())
		staff.GET("/get/:id", middleware.Authentication(models.M_Acc), controllers.GetStaff())
		staff.POST("/create", middleware.Authentication(models.M_Acc), controllers.CreateStaff())
		staff.PUT("/update-all/:id", middleware.Authentication(models.M_Acc), controllers.UpdateStaffDetails())
		staff.PATCH("/update-status/:id", middleware.Authentication(models.M_Acc), controllers.UpdateStaffStatus())
		staff.PATCH("/update-profile-pic/:id", middleware.Authentication(models.M_Acc), controllers.UpdateStaffProfilePicture())
		staff.DELETE("/delete/:id", middleware.Authentication(models.M_Acc), controllers.DeleteStaff())
		staff.POST("/search", middleware.Authentication(models.M_Acc), controllers.SearchStaffData()) // search by first_name , last_name and gender
		staff.POST("/filter", middleware.Authentication(models.M_Acc), controllers.FilterStaff())
	}
}
