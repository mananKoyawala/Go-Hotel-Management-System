package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
	"github.com/mananKoyawala/hotel-management-system/pkg/middleware"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
)

func RoomRoutes(r *gin.Engine) {
	room := r.Group("/room")
	{
		// * ALL
		room.GET("/getall", controllers.GetRooms())
		room.GET("/getall/:id", controllers.GetRoomsByBranch()) // By Branch id only
		room.GET("/get/:id", controllers.GetRoom())
		// * Manager
		room.POST("/create", middleware.Authentication(models.M_Acc), controllers.CreateRoom())
		room.PUT("/update-all/:id", middleware.Authentication(models.M_Acc), controllers.UpdateRoomDetails())
		room.PATCH("/add-image/:id", middleware.Authentication(models.M_Acc), controllers.RoomAddImage())
		room.DELETE("/delete-image/:id", middleware.Authentication(models.M_Acc), controllers.RoomRemoveImage())
		room.DELETE("/delete/:id", middleware.Authentication(models.M_Acc), controllers.DeleteRoom())
		room.POST("/filter", middleware.Authentication(models.M_Acc), controllers.FilterRoom())
	}
}
