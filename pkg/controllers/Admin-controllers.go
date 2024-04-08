package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/database"
	"github.com/mananKoyawala/hotel-management-system/pkg/helpers"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
	"github.com/mananKoyawala/hotel-management-system/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// * ADMIN only Login when it's token is expired
// * DONE
func AdminLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var admin models.Admin
		var foundAdmin models.Admin

		if err := c.BindJSON(&admin); err != nil {
			utils.Error(c, utils.BadRequest, "Invalid JSON Format")
			return
		}

		msg, validate := validateAdmin(admin)
		if !validate {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		if err := database.AdminCollection.FindOne(ctx, bson.M{"email": admin.Email}).Decode(&foundAdmin); err != nil {
			utils.Error(c, utils.NotFound, "Can't find admin email")
			return
		}

		msg, err := helpers.VerifyPassword(admin.Password, foundAdmin.Password)
		if err != nil {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// Generate the tokens
		token, refreshToken, _ := helpers.GenerateAllTokens(foundAdmin.Email, foundAdmin.First_Name, foundAdmin.Last_Name, foundAdmin.Admin_id, string(foundAdmin.Access_Type))

		// Update tokens
		if err = helpers.UpdateAllTokens(token, refreshToken, "admin_id", foundAdmin.Admin_id); err != nil {
			utils.Error(c, utils.InternalServerError, "Error occured while updating tokens")
			return
		}

		// Include tokens in the response
		foundAdmin.Token = token
		foundAdmin.Refresh_Token = refreshToken
		utils.Response(c, foundAdmin)
	}
}

// * DONE
func CreateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var admin models.Admin

		if err := c.BindJSON(&admin); err != nil {
			utils.Error(c, utils.BadRequest, "Invalid JSON Format")
			return
		}

		password, err := helpers.HashPassword(admin.Password)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't generate hash of password")
			return
		}

		admin.Password = password

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
			utils.Error(c, utils.InternalServerError, err.Error())
			return
		}

		utils.Response(c, result)
	}
}

// * DONE
func validateAdmin(admin models.Admin) (string, bool) {

	if admin.Email == "" {
		return "Email required", false
	}

	if admin.Password == "" {
		return "Password required", false
	}
	return "", true
}
