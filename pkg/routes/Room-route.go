package routes

import "github.com/gin-gonic/gin"

func RoomRoutes(r *gin.Engine) {
	room := r.Group("/room")
	{
		room.GET("/getall")
		room.GET("/get/:id")
		room.POST("/create")
		room.PUT("/update-all/:id")
		room.PATCH("/update-room-type/:id")
		room.PATCH("/update-room-availability/:id")
		room.PATCH("/update-cleaning-status/:id")
		room.DELETE("/delete/:id")
	}
}
