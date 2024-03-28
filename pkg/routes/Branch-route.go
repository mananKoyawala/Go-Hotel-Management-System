package routes

import "github.com/gin-gonic/gin"

func BranchRoutes(r *gin.Engine) {
	branch := r.Group("branch")
	{
		// access to all
		branch.GET("/getall")
		branch.GET("/get/:id")
		// access to only admin
		branch.POST("/create")
		branch.PUT("/update-all/:id")
		branch.PATCH("/update-branch-status/:id")
		branch.DELETE("/delete/:id")
	}
}
