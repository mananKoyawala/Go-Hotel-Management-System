package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/database"
	"github.com/mananKoyawala/hotel-management-system/pkg/helpers"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AdminLogin() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func CreateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var admin models.Admin

		if err := c.BindJSON(&admin); err != nil {
			c.JSON(500, gin.H{"error": "Invalid JSON"})
			return
		}

		admin.ID = primitive.NewObjectID()
		admin.Admin_id = admin.ID.Hex()
		admin.Access_Type = models.Admin_Access
		admin.Created_at, _ = helpers.GetTime()
		admin.Updated_at, _ = helpers.GetTime()

		token, refershToken, _ := helpers.GenerateAllTokens(admin.Email, admin.First_Name, admin.Last_Name, admin.Admin_id, string(admin.Access_Type))

		admin.Token = token
		admin.Refresh_Token = refershToken

		result, err := database.AdminCollection.InsertOne(ctx, admin)

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, result)
	}
}
