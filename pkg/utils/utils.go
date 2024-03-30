package utils

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
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

func IsNonNegative(n int) bool {
	return n > 0
}

func ValidatePassword(password string) (string, bool) {
	minLength := 6
	maxLength := 10
	hasUppercase := false
	hasLowercase := false
	hasDigit := false
	hasSpecial := false

	if password == "" {
		return "Password is required", false
	}

	if len(password) < minLength || len(password) > maxLength {
		return "Password must be between 6 to 10 characters", false
	}

	for _, char := range password {
		switch {
		case 'a' <= char && char <= 'z':
			hasLowercase = true
		case 'A' <= char && char <= 'Z':
			hasUppercase = true
		case '0' <= char && char <= '9':
			hasDigit = true
		case char == '@' || char == '$' || char == '!' || char == '*' || char == '%' || char == '?' || char == '&':
			hasSpecial = true
		}
	}

	if !hasLowercase {
		return "Password must contain at least one lowercase letter", false
	}

	if !hasUppercase {
		return "Password must contain at least one uppercase letter", false
	}

	if !hasDigit {
		return "Password must contain at least one digit", false
	}

	if !hasSpecial {
		return "Password must contain at least one special character (@$!*%?&)", false
	}

	return "", true // Password passed all validations
}

func CheckLength(num int, length int) bool {
	return len(strconv.Itoa(num)) != length
}
