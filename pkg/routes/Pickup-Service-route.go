package routes

import "github.com/gin-gonic/gin"

func PickupServiceRoutes(r *gin.Engine) {
	service := r.Group("/service")
	{
		service.GET("/getall")
		service.GET("/get/:id")
		service.POST("/create")
		service.PUT("/update-details/:id")
		service.PATCH("/update-status/:id") // update status of a service completed or not
		service.DELETE("/delete/:id")       // Cancel the pickup service
	}
}
