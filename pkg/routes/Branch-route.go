package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
	"github.com/mananKoyawala/hotel-management-system/pkg/middleware"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
)

func BranchRoutes(r *gin.Engine) {
	branch := r.Group("/branch")
	{
		// r.Use(middleware.Authentication())
		// * ALL
		branch.GET("/getall", controllers.GetBranches())
		branch.GET("/get/:id", controllers.GetBranch())

		// * Admin
		branch.GET("/get-branch-by-status/:status", middleware.Authentication(models.A_Acc), controllers.GetBranchesByStatus())
		branch.POST("/create", middleware.Authentication(models.A_Acc), controllers.CreateBranch())
		branch.PUT("/update-all/:id", middleware.Authentication(models.A_Acc), controllers.UpdateBranchDetails())
		branch.PATCH("/update-branch-status/:id", middleware.Authentication(models.A_Acc), controllers.UpdateBranchStatus())
		branch.PATCH("/add-image/:id", middleware.Authentication(models.A_Acc), controllers.BranchAddImage())
		branch.DELETE("/delete-image/:id", middleware.Authentication(models.A_Acc), controllers.BranchRemoveImage())
		branch.DELETE("/delete/:id", middleware.Authentication(models.A_Acc), controllers.DeleteBranch())
		branch.POST("/search", middleware.Authentication(models.A_Acc), controllers.SearchBranchData())
		// branch data can be search by branch_name,Address, City, State, Country, status
		branch.POST("/filter", middleware.Authentication(models.A_Acc), controllers.FilterBranch()) // filter by city, state, country, status
	}
}
