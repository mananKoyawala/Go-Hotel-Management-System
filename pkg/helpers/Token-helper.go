package helpers

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mananKoyawala/hotel-management-system/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email       string
	First_name  string
	Last_name   string
	Id          string
	Access_Type string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email, firstName, lastName, id, access_type string) (signedToken, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:       email,
		First_name:  firstName,
		Last_name:   lastName,
		Id:          id,
		Access_Type: access_type,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refershClaims := &SignedDetails{
		Access_Type: access_type,
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

// Here idname is provided beacuse there are multiple access ids like admin_id, user_id, etc
func UpdateAllTokens(signedToken, signedRefreshToken, idname, id string) error {
	ctx, cancel := GetContext()
	defer cancel()

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
	Updated_at, _ := GetTime()
	updateObj = append(updateObj, bson.E{Key: "updated_at", Value: Updated_at})
	upsert := true
	filter := bson.M{"admin_id": id}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := database.AdminCollection.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updateObj},
	}, &opt)

	if err != nil {
		return err
	}
	return nil
}
