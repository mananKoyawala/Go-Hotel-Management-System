package routes

import "github.com/gin-gonic/gin"

func GuestRoutes(r *gin.Engine) {
	guest := r.Group("/guest")
	{
		guest.GET("/getall")
		guest.GET("/get/:id")
		guest.POST("/create")
		guest.PUT("/update/:id")
		guest.DELETE("/delete/:id")
	}
}
