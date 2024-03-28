package routes

import "github.com/gin-gonic/gin"

func BranchRoutes(r *gin.Engine) {
	branch := r.Group("branch")
	{
		branch.GET("/getall")
		branch.GET("/get/:id")
		branch.POST("/create")
		branch.PUT("/update-all/:id")
		branch.PATCH("/update-branch-status/:id")
		branch.DELETE("/delete/:id")
	}
}
