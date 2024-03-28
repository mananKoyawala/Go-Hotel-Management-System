package routes

import "github.com/gin-gonic/gin"

func FeedbackRoutes(r *gin.Engine) {
	feedback := r.Group("/feedback")
	{
		feedback.GET("/getall/:id") // Get all the feedback by branch
		feedback.GET("/get/:id")    // Get the feedback by one feedback id
		feedback.POST("/create")
		feedback.PATCH("/update-resolution-details/:id") // Only manager can reply the feedback
		feedback.DELETE("/delete/:id/:uid")              // Only guest can delete this feedback
	}
}
