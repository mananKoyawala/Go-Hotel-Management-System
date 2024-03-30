package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
	"github.com/mananKoyawala/hotel-management-system/pkg/middleware"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
)

func BranchRoutes(r *gin.Engine) {
	branch := r.Group("branch")
	{
		// access to all
		branch.GET("/getall", controllers.GetBranches())
		branch.GET("/get/:id", controllers.GetBranch())
		branch.GET("/get-branch-by-status/:status", controllers.GetBranchesByStatus())
		// access to only admin
		r.Use(middleware.Authentication(models.Admin_Access))
		branch.POST("/create", controllers.CreateBranch())
		branch.PUT("/update-all/:id", controllers.UpdateBranchDetails())
		branch.PATCH("/update-branch-status/:id", controllers.UpdateBranchStatus())
		branch.DELETE("/delete/:id", controllers.DeleteBranch())
	}
}
