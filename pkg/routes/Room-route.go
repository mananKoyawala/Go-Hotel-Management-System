package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
)

func RoomRoutes(r *gin.Engine) {
	room := r.Group("/room")
	{
		// * ALL
		// r.Use(middleware.Authentication())
		room.GET("/getall", controllers.GetRooms())
		room.GET("/getall/:id", controllers.GetRoomsByBranch()) // By Branch id only
		room.GET("/get/:id", controllers.GetRoom())
		// * Manager
		room.POST("/create", controllers.CreateRoom())
		room.PUT("/update-all/:id", controllers.UpdateRoomDetails())
		room.PATCH("/add-image/:id", controllers.RoomAddImage())
		room.DELETE("/delete-image/:id", controllers.RoomRemoveImage())
		room.DELETE("/delete/:id", controllers.DeleteRoom())
		room.POST("/filter", controllers.FilterRoom())
	}
}
