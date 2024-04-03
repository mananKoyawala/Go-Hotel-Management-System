package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
)

func ReservationRoutes(r *gin.Engine) {
	reservation := r.Group("/reservation")
	{
		// * Manager, User
		reservation.GET("/getall", controllers.GetAllReservations())
		reservation.GET("/get/:id", controllers.GetReservation())
		reservation.POST("/create", controllers.CreateReservation())
		reservation.PUT("/update-all/:id", controllers.UpdateReservationDetails())  // Give details update otherwise otherthings as it is
		reservation.DELETE("/delete/:id/:room_id", controllers.DeleteReservation()) // Cancel the reservation
	}
}
