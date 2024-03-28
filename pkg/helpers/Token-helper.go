package helpers

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type SignedDetails struct {
	Email       string
	First_name  string
	Last_name   string
	Uid         string
	Access_Type string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email, firstName, lastName, uid, access_type string) (signedToken, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:       email,
		First_name:  firstName,
		Last_name:   lastName,
		Uid:         uid,
		Access_Type: access_type,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refershClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refershClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	//Token is invalid

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "The token is invalid " + err.Error()
		return claims, msg
	}

	// Token is expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token is expired " + err.Error()
		return claims, msg
	}

	return claims, ""
}
