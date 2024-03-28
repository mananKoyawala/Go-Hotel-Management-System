package routes

import "github.com/gin-gonic/gin"

func DriverRoutes(r *gin.Engine) {
	driver := r.Group("/driver")
	{
		// access to manager, guest
		driver.GET("/getall")
		driver.GET("/get/:id")
		// access to manager
		driver.POST("/create")
		driver.PUT("/update-all/:id")
		driver.PATCH("/update-status/:id")
		driver.PATCH("/update-availability/:id")
		driver.DELETE("/delete/:id")
	}
}
