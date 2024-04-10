package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/controllers"
	"github.com/mananKoyawala/hotel-management-system/pkg/middleware"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
)

func FeedbackRoutes(r *gin.Engine) {
	feedback := r.Group("/feedback")
	{
		feedback.GET("/getall/:id", middleware.Authentication(models.M_Acc, models.A_Acc), controllers.GetAllFeedbacks())                               // Get all the feedback by branch
		feedback.GET("/get/:id", middleware.Authentication(models.M_Acc, models.A_Acc), controllers.GetFeedback())                                      // Get the feedback by one feedback id
		feedback.POST("/create", middleware.Authentication(models.M_Acc, models.G_Acc), controllers.CreateFeedback())                                   // access to manager, guest
		feedback.PATCH("/update-resolution-details/:id", middleware.Authentication(models.M_Acc, models.A_Acc), controllers.UpdateFeedbackResolution()) // Only manager,admin can reply the feedback
		feedback.DELETE("/delete/:id", middleware.Authentication(models.G_Acc), controllers.DeleteFeedback())                                           // Only guest can delete this feedback
		feedback.POST("/filter", middleware.Authentication(models.M_Acc, models.A_Acc), controllers.FilterFeedback())                                   // filter by status, feedback_type, rating

	}
}
