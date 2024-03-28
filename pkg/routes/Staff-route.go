package routes

import "github.com/gin-gonic/gin"

func StaffRoutes(r *gin.Engine) {
	staff := r.Group("/staff")
	{
		staff.GET("/getall")
		staff.GET("/get/:id")
		staff.POST("/create")
		staff.PUT("/update-all/:id")
		staff.PATCH("/update-job-type/:id")
		staff.PATCH("/update-status/:id")
		staff.DELETE("/delete/:id")
	}
}
