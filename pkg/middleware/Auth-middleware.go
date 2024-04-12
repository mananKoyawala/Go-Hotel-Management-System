package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/helpers"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
	"github.com/mananKoyawala/hotel-management-system/pkg/utils"
)

// access_type ...models.Access_Type
func Authentication(access_type ...models.Access_Type) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("X-Auth-Token")

		if clientToken == "" {
			utils.Error(c, utils.BadRequest, "Token is missing")
			c.Abort()
			return
		}

		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			utils.Error(c, utils.InternalServerError, err)
			c.Abort()
			return
		}

		accessMap := make(map[models.Access_Type]bool)
		for _, access := range access_type {
			accessMap[access] = true
		}
		// log.Println(claims.Access_Type)
		// log.Println(access_type)
		// Check if the access type from claims matches any of the allowed access types
		if !accessMap[models.Access_Type(claims.Access_Type)] {
			utils.Error(c, utils.Unauthorized, "Unauthorized access to "+claims.Access_Type)
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("id", claims.Id)
		c.Set("access_type", claims.Access_Type)
		c.Next()
	}
}
