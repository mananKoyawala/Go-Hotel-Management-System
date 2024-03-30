package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
	"github.com/mananKoyawala/hotel-management-system/pkg/middleware"
)

func RoomRoutes(r *gin.Engine) {
	room := r.Group("/room")
	{
		// * ALL
		r.Use(middleware.Authentication())
		room.GET("/getall", controllers.GetRooms())
		room.GET("/getall/:id", controllers.GetRoomsByBranch()) // By Branch id only
		room.GET("/get/:id", controllers.GetRoom())
		// * Manager
		room.POST("/create", controllers.CreateRoom())
		room.PUT("/update-all/:id", controllers.UpdateRoomDetails())
		room.DELETE("/delete/:id", controllers.DeleteRoom())
	}
}
