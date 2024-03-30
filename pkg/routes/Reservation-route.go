package routes

import "github.com/gin-gonic/gin"

func ReservationRoutes(r *gin.Engine) {
	reservation := r.Group("/reservation")
	{
		// * Manager, User
		reservation.GET("/getall")
		reservation.GET("/get")
		reservation.POST("/create")
		reservation.PUT("/update-all/:id") // Give details update otherwise otherthings as it is
		reservation.DELETE("/delete/:id")  // Cancel the reservation
	}
}
