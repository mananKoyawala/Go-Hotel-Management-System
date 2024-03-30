package utils

import (
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

var TimeOut = 30 * time.Second

// 500, 404, 400, 200, 409, 401
var InternalServerError = http.StatusInternalServerError // for error while id convertion
var NotFound = http.StatusNotFound                       // for id not found
var BadRequest = http.StatusBadRequest                   // for invalid json or empty data
var OK = http.StatusOK                                   // for ok status
var Conflict = http.StatusConflict                       // for inserting data already exist like email, password
var Unauthorized = http.StatusUnauthorized               // Email or password incorrect

func Error(c *gin.Context, statusCode int, errorMessage string) {
	log.Printf("Error: %s\n", errorMessage)
	c.Header("Content-Type", "application/json")
	c.JSON(statusCode, gin.H{"error": errorMessage})
} // to return error

func Message(c *gin.Context, message string) {
	c.Header("Content-Type", "application/json")
	c.JSON(OK, gin.H{"message": message})
} // to return messages

func Response(c *gin.Context, data interface{}) {
	c.Header("Content-Type", "application/json")
	c.JSON(OK, data)
} // to return response

func ValidateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
