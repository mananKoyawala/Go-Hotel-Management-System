package controllers

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/database"
	"github.com/mananKoyawala/hotel-management-system/pkg/helpers"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
	emailverification "github.com/mananKoyawala/hotel-management-system/pkg/service/Email-Verification"
	imageupload "github.com/mananKoyawala/hotel-management-system/pkg/service/image-upload"
	"github.com/mananKoyawala/hotel-management-system/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var guestFolder = "guest"

// * DONE
func GetAllGuest() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 10 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		currentIndex := page * recordPerPage
		// Here the aggregation pipeline started

		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "null"},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}}}

		projectStage1 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "guests",
				Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
		}}}

		projectStage2 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "total_count", Value: 1},
			{Key: "guest", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$guests"},
					{Key: "as", Value: "data"},
					{Key: "in", Value: bson.D{
						{Key: "guest_id", Value: "$$data.guest_id"},
						{Key: "id_proof_type", Value: "$$data.id_proof_type"},
						{Key: "first_name", Value: "$$data.first_name"},
						{Key: "last_name", Value: "$$data.last_name"},
						{Key: "phone", Value: "$$data.phone"},
						{Key: "email", Value: "$$data.email"},
						{Key: "gender", Value: "$$data.gender"},
						{Key: "country", Value: "$$data.country"},
						{Key: "image", Value: "$$data.image"},
					}},
				}},
			}},
			{Key: "hashMoreData", Value: bson.D{
				{Key: "$cond", Value: bson.D{
					{Key: "if", Value: bson.D{
						{Key: "$eq", Value: bson.A{"$total_count", currentIndex}},
					}},
					{Key: "then", Value: false},
					{Key: "else", Value: bson.D{
						{Key: "$cond", Value: bson.D{
							{Key: "if", Value: bson.D{
								{Key: "$gt", Value: bson.A{"$total_count", currentIndex}},
							}},
							{Key: "then", Value: true},
							{Key: "else", Value: false},
						}},
					}},
				}},
			}},
		}}}

		result, err := database.GuestCollection.Aggregate(ctx, mongo.Pipeline{
			groupStage, projectStage1, projectStage2,
		})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching guests "+err.Error())
			return
		}

		var allguests []bson.M
		if err := result.All(ctx, &allguests); err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting the guests "+err.Error())
			return
		}

		if len(allguests) == 0 {
			utils.Response(c, []interface{}{})
			return
		}

		utils.Response(c, allguests[0])
	}
}

// * DONE
func GetGuest() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var guest models.Guest

		id := c.Param("id")

		if err := database.GuestCollection.FindOne(ctx, bson.M{"guest_id": id}).Decode(&guest); err != nil {
			utils.Error(c, utils.NotFound, "Can't find guest with id")
			return
		}

		utils.Response(c, guest)
	}
}

// * DONE
func GuestSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var guest models.Guest

		guest.First_Name = c.PostForm("first_name")
		guest.Last_Name = c.PostForm("last_name")
		guest.Phone, _ = strconv.Atoi(c.PostForm("phone"))
		guest.Email = c.PostForm("email")
		guest.Password = c.PostForm("password")
		guest.Country = c.PostForm("country")
		guest.Gender = c.PostForm("gender")
		guest.ID_Proof_Type = c.PostForm("id_proof_type")

		// validate guest details
		msg, isVal := validateGuestDetails(guest)
		if !isVal {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// check file is valid
		file, handler, err := c.Request.FormFile("file")
		if err != nil {
			utils.Error(c, utils.BadRequest, "File was not provided or Invalid file.")
			return
		}
		defer file.Close()

		// check the image file is .png , .jpg , .jpeg
		ext := filepath.Ext(handler.Filename)
		if ext != ".jpeg" && ext != ".jpg" && ext != ".png" {
			utils.Error(c, utils.BadRequest, "Invalid Image file format. Only JPEG, JPG, or PNG files are allowed.")
			return
		}

		count, err := database.GuestCollection.CountDocuments(ctx, bson.M{"email": guest.Email})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting details.")
			return
		}

		if count > 0 {
			utils.Error(c, utils.Conflict, "Email already in use, try different email addresses.")
			return
		}

		// hash password
		password, err := helpers.HashPassword(guest.Password)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Unable to generate hash password")
			return
		}
		guest.Password = password

		// generate id, timestamps
		guest.Access_Type = models.Guest_Access
		guest.ID = primitive.NewObjectID()
		guest.Guest_id = guest.ID.Hex()
		guest.IsVerified = "false"
		guest.Created_at, _ = helpers.GetTime()
		guest.Updated_at, _ = helpers.GetTime()

		// generate tokens
		token, refreshToken, _ := helpers.GenerateAllTokens(guest.Email, guest.First_Name, guest.Last_Name, guest.Guest_id, string(guest.Access_Type))

		// update the tokens
		guest.Token = token
		guest.Refresh_Token = refreshToken

		// upload image
		name := strings.ReplaceAll(handler.Filename, " ", "")
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), name)
		url, err := imageupload.UploadService(file, guestFolder, filename)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't uplaod the image.")
			return
		}
		guest.Image = url

		// Insert the details
		result, err := database.GuestCollection.InsertOne(ctx, guest)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't create guest")
			return
		}

		if err := emailverification.GenerateEmailVerificationLink(guest.Guest_id, guest.Email); err != nil {
			utils.Error(c, utils.InternalServerError, err.Error())
			return
		}

		// if success return
		c.JSON(utils.OK, gin.H{
			"result":  result,
			"message": "Verification link is send to email",
		})
	}
}

// * DONE
func GuestLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var guest models.Guest
		var foundGuest models.Guest

		guest.Email = c.PostForm("email")
		guest.Password = c.PostForm("password")

		// Check Email
		if err := database.GuestCollection.FindOne(ctx, bson.M{"email": guest.Email}).Decode(&foundGuest); err != nil {
			utils.Error(c, utils.NotFound, "Can't find guest with Email id.")
			return
		}

		if foundGuest.IsVerified != "true" {
			utils.Error(c, utils.BadRequest, "Guest is not verified.")
			return
		}

		// Verify Password
		msg, err := helpers.VerifyPassword(guest.Password, foundGuest.Password)
		if err != nil {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// Generate All Tokens
		token, refreshToken, _ := helpers.GenerateAllTokens(foundGuest.Email, foundGuest.First_Name, foundGuest.Last_Name, foundGuest.Guest_id, string(foundGuest.Access_Type))

		// Update Tokens
		if err := helpers.UpdateAllTokens(token, refreshToken, "guest_id", foundGuest.Guest_id); err != nil {
			utils.Error(c, utils.InternalServerError, "Error occured while updating tokens")
			return
		}

		foundGuest.Token = token
		foundGuest.Refresh_Token = refreshToken

		// Return as response
		utils.Response(c, foundGuest)
	}
}

// * DONE
func VerifyGuest() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()

		if err := emailverification.VerifyEmail(c, ctx); err != nil {
			htmlContentError := `<!DOCTYPE html>
		<html>
		  <head>
			<style>
			  body {
				display: flex;
				justify-content: center;
				align-items: center;
				height: 100vh;
				margin: 0;
				background-color: #f0f0f0;
			  }
		
			  .message {
				text-align: center;
				font-size: 34px;
				color: black;
			  }
			</style>
		  </head>
		  <body>
			<div class="message">` + err.Error() + `</div>
		  </body>
		</html>`
			log.Println(err.Error())
			c.Data(utils.OK, "text/html; charset=utf-8", []byte(htmlContentError))
			return
		}

		htmlContent := `<!DOCTYPE html>
		<html>
		  <head>
			<style>
			  body {
				display: flex;
				justify-content: center;
				align-items: center;
				height: 100vh;
				margin: 0;
				background-color: #f0f0f0;
			  }
		
			  .message {
				text-align: center;
				font-size: 34px;
				color: black;
			  }
			</style>
		  </head>
		  <body>
			<div class="message">Guest verification successfully done.</div>
		  </body>
		</html>`

		c.Data(utils.OK, "text/html; charset=utf-8", []byte(htmlContent))
	}
}

// * DONE
func UpdateGuestDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var guest models.Guest
		var foundGuest models.Guest
		id := c.Param("id")

		guest.First_Name = c.PostForm("first_name")
		guest.Last_Name = c.PostForm("last_name")
		guest.Phone, _ = strconv.Atoi(c.PostForm("phone"))
		guest.Country = c.PostForm("country")
		guest.Gender = c.PostForm("gender")

		// Validate data
		msg, isVal := validateUpdateGuestDetails(guest)
		if !isVal {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// check guest exist or not
		if err := database.GuestCollection.FindOne(ctx, bson.M{"guest_id": id}).Decode(&foundGuest); err != nil {
			utils.Error(c, utils.NotFound, "Can't find guest with id")
			return
		}

		// Update data here
		var updateObj primitive.D

		updateObj = append(updateObj, bson.E{Key: "first_name", Value: guest.First_Name})

		updateObj = append(updateObj, bson.E{Key: "last_name", Value: guest.Last_Name})

		updateObj = append(updateObj, bson.E{Key: "gender", Value: guest.Gender})

		updateObj = append(updateObj, bson.E{Key: "country", Value: guest.Country})

		updateObj = append(updateObj, bson.E{Key: "phone", Value: guest.Phone})

		guest.Updated_at, _ = helpers.GetTime()
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: guest.Updated_at})

		filter := bson.M{"guest_id": id}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		// If success send success message
		_, err := database.GuestCollection.UpdateOne(ctx, filter, bson.D{
			{Key: "$set", Value: updateObj},
		}, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update guest details.")
			return
		}

		utils.Message(c, "Guest details updated successfully.")
	}
}

// * DONE
func DeleteGuest() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var guest models.Guest

		id := c.Param("id")

		if err := database.GuestCollection.FindOne(ctx, bson.M{"guest_id": id}).Decode(&guest); err != nil {
			utils.Error(c, utils.NotFound, "Can't find guest with id")
			return
		}

		image := utils.GetTrimedUrl(guest.Image)
		if err := imageupload.DeleteService(image); err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete image."+err.Error())
			return
		}

		_, err := database.GuestCollection.DeleteOne(ctx, bson.M{"guest_id": id})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete guest.")
			return
		}
		utils.Message(c, "Guest deleted successfully.")
	}
}

// * DONE
func ResetGuestPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var guest models.Guest

		// check format
		if err := c.BindJSON(&guest); err != nil {
			utils.Error(c, utils.BadRequest, "Invalid JSON Format.")
			return
		}

		// validate email
		if !utils.ValidateEmail(guest.Email) {
			utils.Error(c, utils.BadRequest, "Invalid email.")
			return
		}

		// check email exist
		count, err := database.GuestCollection.CountDocuments(ctx, bson.M{"email": guest.Email})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting guest details.")
			return
		}

		if !(count > 0) {
			utils.Error(c, utils.NotFound, "Can't find guest with Email id.")
			return
		}

		// validate password
		msg, val := utils.ValidatePassword(guest.Password)
		if !val {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// hash password
		password, err := helpers.HashPassword(guest.Password)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Unable to hash password.")
			return
		}

		// update password,timestamp
		guest.Password = password
		guest.Updated_at, _ = helpers.GetTime()

		// update details

		filter := bson.M{"email": guest.Email}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		updateObj := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "Password", Value: guest.Password},
				{Key: "updated_at", Value: guest.Updated_at},
			}},
		}

		_, err = database.GuestCollection.UpdateOne(ctx, filter, updateObj, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update guest password.")
			return
		}

		// if success return
		utils.Message(c, "Password updated successfully.")
	}
}

// * DONE
func UpdateGuestProfilePicture() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var guest models.Guest

		// get id
		id := c.Param("id")

		// check the file
		file, handler, err := c.Request.FormFile("file")
		if err != nil {
			utils.Error(c, utils.BadRequest, "File was not provided or Invalid file.")
			return
		}
		defer file.Close()

		// check the image file is .png , .jpg , .jpeg
		ext := filepath.Ext(handler.Filename)
		if ext != ".jpeg" && ext != ".jpg" && ext != ".png" {
			utils.Error(c, utils.BadRequest, "Invalid Image file format. Only JPEG, JPG, or PNG files are allowed.")
			return
		}

		// get url details for image url
		if err := database.GuestCollection.FindOne(ctx, bson.M{"guest_id": id}).Decode(&guest); err != nil {
			utils.Error(c, utils.NotFound, "Can't find guest with id")
			return
		}
		log.Println(guest.Image)
		// delete file
		image := utils.GetTrimedUrl(guest.Image)
		if err := imageupload.DeleteService(image); err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete image."+err.Error())
			return
		}

		// upload new file
		name := strings.ReplaceAll(handler.Filename, " ", "")
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), name)
		url, err := imageupload.UploadService(file, guestFolder, filename)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't upload image."+err.Error())
			return
		}

		// update new uploaded file
		updated_at, _ := helpers.GetTime()
		filter := bson.M{"guest_id": id}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}
		updateObj := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "image", Value: url},
				{Key: `updated_at`, Value: updated_at},
			}},
		}

		_, err = database.GuestCollection.UpdateOne(ctx, filter, updateObj, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update image.")
			return
		}

		// if success return
		utils.Message(c, "Image updated successfully.")
	}
}

func validateGuestDetails(guest models.Guest) (string, bool) {

	if guest.First_Name == "" {
		return "First name is required", false
	}

	if guest.Last_Name == "" {
		return "Last name is required", false
	}

	if utils.CheckLength(guest.Phone, 10) {
		return "Phone number must be 10 digits", false
	}

	if guest.Email == "" {
		return "Email address required", false
	} else if !utils.ValidateEmail(guest.Email) {
		return "Invalid email address", false
	}

	msg, val := utils.ValidatePassword(guest.Password)
	if !val {
		return msg, false
	}

	if guest.Gender == "" {
		return "Gender is required", false
	}

	if guest.Country == "" {
		return "Country is required", false
	}

	if !checkIdProofType(guest.ID_Proof_Type) {
		return "Guest Id proof must be aadhar_card,passport,pan_card or driving_license.", false
	}

	return "", true
}

func validateUpdateGuestDetails(guest models.Guest) (string, bool) {

	if guest.First_Name == "" {
		return "First name is required", false
	}

	if guest.Last_Name == "" {
		return "Last name is required", false
	}

	if utils.CheckLength(guest.Phone, 10) {
		return "Phone number must be 10 digits", false
	}

	if guest.Gender == "" {
		return "Gender is required", false
	}

	if guest.Country == "" {
		return "Country is required", false
	}

	return "", true
}

func checkIdProofType(proof string) bool {
	if proof == "" {
		return false
	}

	if models.ID_Proof_Type(proof) == models.Aadhar_Card {
		return true
	}

	if models.ID_Proof_Type(proof) == models.Pan_Card {
		return true
	}

	if models.ID_Proof_Type(proof) == models.PassPort {
		return true
	}

	if models.ID_Proof_Type(proof) == models.Driving_License {
		return true
	}
	return false
}
