package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
	"github.com/mananKoyawala/hotel-management-system/pkg/middleware"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
)

func ReservationRoutes(r *gin.Engine) {
	reservation := r.Group("/reservation")
	{
		// * Manager, User
		reservation.GET("/getall", middleware.Authentication(models.M_Acc), controllers.GetAllReservations())
		reservation.GET("/get/:id", middleware.Authentication(models.M_Acc), controllers.GetReservation())
		reservation.POST("/create", middleware.Authentication(models.M_Acc, models.G_Acc), controllers.CreateReservation())
		reservation.PUT("/update-all/:id", middleware.Authentication(models.M_Acc, models.G_Acc), controllers.UpdateReservationDetails())  // Give details update otherwise otherthings as it is
		reservation.DELETE("/delete/:id/:room_id", middleware.Authentication(models.M_Acc, models.G_Acc), controllers.DeleteReservation()) // Cancel the reservation
	}
}
