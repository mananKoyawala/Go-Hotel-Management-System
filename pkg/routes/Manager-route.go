package routes

import "github.com/gin-gonic/gin"

func ManagerRoutes(r *gin.Engine) {
	manager := r.Group("/manager")
	{
		manager.GET("/getall") // controllers.LoginController()
		manager.GET("/get/:id")
		manager.POST("/create")
		manager.PUT("/update-all/:id")
		manager.PUT("/update-status/:id")
		manager.PUT("/update-branch/:id")
		manager.DELETE("/delete/:id")
	}
}

// PUT All the feilds or many fields
// PATCH particaly one or two fields
