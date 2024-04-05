package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
)

func BranchRoutes(r *gin.Engine) {
	branch := r.Group("/branch")
	{
		// r.Use(middleware.Authentication())
		// * ALL
		branch.GET("/getall", controllers.GetBranches())
		branch.GET("/get/:id", controllers.GetBranch())
		// * Admin
		branch.GET("/get-branch-by-status/:status", controllers.GetBranchesByStatus())
		branch.POST("/create", controllers.CreateBranch())
		branch.PUT("/update-all/:id", controllers.UpdateBranchDetails())
		branch.PATCH("/update-branch-status/:id", controllers.UpdateBranchStatus())
		branch.PATCH("/add-image/:id", controllers.BranchAddImage())
		branch.DELETE("/delete-image/:id", controllers.BranchRemoveImage())
		branch.DELETE("/delete/:id", controllers.DeleteBranch())
	}
}
