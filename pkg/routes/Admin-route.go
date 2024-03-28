package routes

import "github.com/gin-gonic/gin"

func AdminRoutes(r *gin.Engine) {
	admin := r.Group("/admin")
	{
		// only access to admin
		admin.GET("/login") // controllers.AdminController()
	}
}
