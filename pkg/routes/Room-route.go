package routes

import "github.com/gin-gonic/gin"

func RoomRoutes(r *gin.Engine) {
	room := r.Group("/room")
	{
		// access to all
		room.GET("/getall")
		room.GET("/get/:id")
		// access to manager of that branch
		room.POST("/create")
		room.PUT("/update-all/:id")
		room.PATCH("/update-room-type/:id")
		room.PATCH("/update-room-availability/:id")
		room.PATCH("/update-cleaning-status/:id")
		room.DELETE("/delete/:id")
	}
}
