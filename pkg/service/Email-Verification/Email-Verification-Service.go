package emailverification

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/database"
	"github.com/mananKoyawala/hotel-management-system/pkg/helpers"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
	"github.com/resend/resend-go/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
https://dashboard.pawns.app/verify-email/confirm?expires=1712116416&hash=168bd071e7ff5933bb1583aa10e92563b30f209f&id=4318333&signature=2154d030cb8a95646693859c26eaa2b7d77953d10028a480d630279eab523243
*/

var SECRET_KEY string = os.Getenv("SECRET_KEY")
var PORT = os.Getenv("PORT")

type EmailVerification struct {
	ID              primitive.ObjectID `bson:"_id"`
	Verification_id string             `bson:"verification_id,omitempty" json:"verification_id,omitempty"`
	Token           string             `bson:"token,omitempty" json:"token,omitempty"`
	Guest_id        string             `bson:"guest_id,omitempty" json:"guest_id,omitempty"`
	ExpiresAt       int64              `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	Created_at      time.Time          `bson:"created_at,omitempty" json:"create_at,omitempty"`
	Updated_at      time.Time          `bson:"updated_at,omitempty" json:"update_at,omitempty"`
}

// * DONE
func GenerateEmailVerificationLink(id, userEmail string) error {
	ctx, cancel := helpers.GetContext()
	defer cancel()
	var verification EmailVerification

	// expires_at, token, guest id , signature
	verification.ExpiresAt = time.Now().Local().Add(time.Hour * time.Duration(24)).Unix()
	token, err := generateToken()
	if err != nil {
		return err
	}
	verification.Token = token
	verification.Guest_id = id

	tokenHash := generateTokenHash(verification.Token)

	// signature using expires_at, hash_token, user_id
	signature := generateSignature(verification.ExpiresAt, tokenHash, verification.Guest_id, SECRET_KEY)

	// add the token, and user id in verification collection
	verification.ID = primitive.NewObjectID()
	verification.Verification_id = verification.ID.Hex()
	verification.Created_at, _ = helpers.GetTime()
	verification.Updated_at, _ = helpers.GetTime()
	_, err = database.VerificationEmailCollection.InsertOne(ctx, verification)
	if err != nil {
		return err
	}

	// generate the link for sending
	link := fmt.Sprintf("http://localhost:%s/guest/verify-email/confirm?expires=%d&hash=%s&id=%s&signature=%s", PORT, verification.ExpiresAt, tokenHash, verification.Guest_id, signature)

	// send the email if success
	err = sendMail(link, userEmail)
	if err != nil {
		return err
	}
	// log.Println(verification.ExpiresAt)
	// log.Println(tokenHash)
	// log.Println(verification.Guest_id)
	// log.Println(signature)
	// log.Println(link)

	// commit the transaction

	return nil
}

// * DONE
func VerifyEmail(c *gin.Context, ctx context.Context) error {
	var verification EmailVerification
	var guest models.Guest

	expires := c.Query("expires")
	hash := c.Query("hash")
	id := c.Query("id")
	signature := c.Query("signature")

	// log.Println(expires)
	// log.Println(hash)
	// log.Println(id)
	// log.Println(signature)

	// check above things are ni
	if expires == "" || hash == "" || signature == "" || id == "" {
		return errors.New("invalid verification link")
	}

	expireTime, _ := strconv.ParseInt(expires, 10, 64)

	// check link expires
	if expireTime < time.Now().Local().Unix() {
		// log.Println(expires)
		// log.Println(string(rune(time.Now().Local().Unix())))
		return errors.New("link is expired")
	}

	// generate the signature by expires, hash_token,
	expires_at, _ := strconv.ParseInt(expires, 10, 64)
	generatedSignature := generateSignature(expires_at, hash, id, SECRET_KEY)

	// compare it with sended signature
	if signature != generatedSignature {
		// log.Println(signature)
		// log.Println(generatedSignature)
		return errors.New("signature is not valid")
	}

	if err := database.GuestCollection.FindOne(ctx, bson.M{"guest_id": id}).Decode(&guest); err != nil {
		return errors.New("error while getting data, can't verify user")
	}

	if guest.IsVerified == "true" {
		return errors.New("user already verified")
	}

	// get user from verification collection
	if err := database.VerificationEmailCollection.FindOne(ctx, bson.M{"guest_id": id}).Decode(&verification); err != nil {
		return errors.New("can't find verification details with id")
	}

	// generate hash of found user token
	newHash := generateTokenHash(verification.Token)

	// compare it
	if newHash != hash {
		return errors.New("invalid hash")
	}

	// if ok update user verifed in database
	updated_at, _ := helpers.GetTime()
	filter := bson.M{"guest_id": id}
	upsert := true
	options := options.UpdateOptions{
		Upsert: &upsert,
	}

	updateObj := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "is_verified", Value: "true"},
			{Key: "updated_at", Value: updated_at},
		}},
	}
	_, err := database.GuestCollection.UpdateOne(ctx, filter, updateObj, &options)
	if err != nil {
		return errors.New("error while update guest verifed status")
	}

	_, err = database.VerificationEmailCollection.DeleteOne(ctx, bson.M{"guest_id": id})
	if err != nil {
		return errors.New("error while deleting verification token")
	}

	// if success return user verifed
	return nil
}

func sendMail(link, userEmail string) error {
	var RESEND_API_KEY = os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(RESEND_API_KEY)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      []string{userEmail},
		Subject: "Email verification",
		Html: `<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Email Verification</title>
				<style>
					body {
						font-family: Arial, sans-serif;
						background-color: #f5f5f5;
						margin: 0;
						padding: 0;
						line-height: 1.6;
						color: #333;
					}
					.container {
						width: 400px; /* Fixed width */
						margin: 20px auto;
						padding: 20px;
						background-color: #fff;
						border-radius: 8px;
						box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
					}
					h1 {
						font-size: 24px;
						text-align: center;
						margin-bottom: 20px;
					}
					.btn {
						display: inline-block;
						padding: 10px 20px;
						background-color: #007bff;
						color: white;
						text-decoration: none;
						border-radius: 4px;
						margin-top: 20px;
					}
					.btn:hover {
					  color: white;
						background-color: #0056b3;
					}
					.instructions {
						margin-top: 20px;
						font-size: 14px;
					}
			
				</style>
			</head>
			<body>
				<table style="width: 100%; height: 100%;">
					<tr>
					<td align="center" valign="top">
					<div class="container" style="background-color: #f9f9f9;">
								<h1>Verify Email Address</h1>
								<p>Thank you for becoming a part of My Hotel!</p>
								<p>To complete your registration, please verify your email address by clicking on the link below:</p>
								<a href="` + link + `" class="btn" style="color:white;">Verify Email Address</a>
								<p class="instructions">If you did not sign up with us, please ignore this email.</p>
								<hr>
								<p class="instructions">If you are having trouble clicking the "Verify Email Address" button, copy and paste the URL below into your web browser:</p>
								<p class="instructions">` + link + `</p>
								<hr>
								<p>You have received this email as a registered user of My Hotel</p>
								<p>&copy; 2024 My Hotel. All rights reserved.</p>
							</div>
						</td>
					</tr>
				</table>
			</body>
			</html>
			`,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return err
	}

	return nil
}

func generateToken() (string, error) {
	var length = 20

	// Define the characters you want to include in the token
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Generate random bytes
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Create a token by selecting characters from the defined character set
	for i := 0; i < length; i++ {
		bytes[i] = characters[int(bytes[i])%len(characters)]
	}

	return string(bytes), nil
}

func generateTokenHash(token string) string {
	// Convert the token string to bytes
	tokenBytes := []byte(token)

	// Create a new SHA-256 hasher
	hasher := sha256.New()

	// Write the token bytes to the hasher
	hasher.Write(tokenBytes)

	// Calculate the hash
	hash := hasher.Sum(nil)

	// Convert the hash bytes to a hexadecimal string
	hashString := hex.EncodeToString(hash)

	return hashString
}

func generateSignature(expiresAt int64, token, guestID, secretKey string) string {
	// Construct the message by concatenating parameters
	message := fmt.Sprintf("%s%s%s", expiresAt, token, guestID)

	// Create a new HMAC hasher with SHA-256
	h := hmac.New(sha256.New, []byte(secretKey))

	// Write the sorted message to the hasher
	h.Write([]byte(message))

	// Calculate the HMAC hash
	signature := h.Sum(nil)

	// Encode the hash to base64
	encodedSignature := base64.StdEncoding.EncodeToString(signature)

	return encodedSignature
}
