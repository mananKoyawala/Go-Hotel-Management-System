package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
)

func FeedbackRoutes(r *gin.Engine) {
	feedback := r.Group("/feedback")
	{
		feedback.GET("/getall/:id", controllers.GetAllFeedbacks())                               // Get all the feedback by branch
		feedback.GET("/get/:id", controllers.GetFeedback())                                      // Get the feedback by one feedback id
		feedback.POST("/create", controllers.CreateFeedback())                                   // access to manager, guest
		feedback.PATCH("/update-resolution-details/:id", controllers.UpdateFeedbackResolution()) // Only manager can reply the feedback
		feedback.DELETE("/delete/:id", controllers.DeleteFeedback())                             // Only guest can delete this feedback
	}
}
