package routes

import "github.com/gin-gonic/gin"

func DriverRoutes(r *gin.Engine) {
	driver := r.Group("/driver")
	{
		// * Manager, User
		driver.GET("/getall")
		driver.GET("/get/:id")
		// * Manager
		driver.POST("/create")
		driver.PUT("/update-all/:id")
		driver.PATCH("/update-status/:id")
		driver.PATCH("/update-availability/:id")
		driver.DELETE("/delete/:id")
	}
}
