package routes

import "github.com/gin-gonic/gin"

func GuestRoutes(r *gin.Engine) {
	guest := r.Group("/guest")
	{
		// access to admin
		guest.GET("/getall")
		guest.GET("/get/:id")
		// unprotected beacuse anyone registered as guest
		guest.POST("/signup")
		guest.POST("/login")
		// only access to guest
		guest.PUT("/update/:id")
		guest.DELETE("/delete/:id")
	}
}
