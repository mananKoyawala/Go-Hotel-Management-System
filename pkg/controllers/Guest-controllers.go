package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mananKoyawala/hotel-management-system/pkg/database"
	"github.com/mananKoyawala/hotel-management-system/pkg/helpers"
	"github.com/mananKoyawala/hotel-management-system/pkg/models"
	"github.com/mananKoyawala/hotel-management-system/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

		// check json
		if err := c.BindJSON(&guest); err != nil {
			utils.Error(c, utils.BadRequest, "Invalid JSON Format")
			return
		}

		// validate guest details
		msg, isVal := validateGuestDetails(guest)
		if !isVal {
			utils.Error(c, utils.BadRequest, msg)
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
		guest.Created_at, _ = helpers.GetTime()
		guest.Updated_at, _ = helpers.GetTime()

		// generate tokens
		token, refreshToken, _ := helpers.GenerateAllTokens(guest.Email, guest.First_Name, guest.Last_Name, guest.Guest_id, string(guest.Access_Type))

		// update the tokens
		guest.Token = token
		guest.Refresh_Token = refreshToken

		// Insert the details
		result, err := database.GuestCollection.InsertOne(ctx, guest)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't create guest")
			return
		}

		// if success return
		utils.Response(c, result)
	}
}

// * DONE
func GuestLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var guest models.Guest
		var foundGuest models.Guest

		if err := c.BindJSON(&guest); err != nil {
			utils.Error(c, utils.BadRequest, "Invalid JSON Format.")
			return
		}

		// Check Email
		if err := database.GuestCollection.FindOne(ctx, bson.M{"email": guest.Email}).Decode(&foundGuest); err != nil {
			utils.Error(c, utils.InternalServerError, "Can't find guest with Email id.")
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
func UpdateGuestDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var guest models.Guest
		id := c.Param("id")

		// Check json
		if err := c.BindJSON(&guest); err != nil {
			utils.Error(c, utils.BadRequest, "Invalid JSON Format")
			return
		}

		// Validate data
		msg, isVal := validateUpdateGuestDetails(guest)
		if !isVal {
			utils.Error(c, utils.BadRequest, msg)
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

		id := c.Param("id")

		_, err := database.GuestCollection.DeleteOne(ctx, bson.M{"guest_id": id})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete guest.")
			return
		}
		utils.Message(c, "Guest deleted successfully.")
	}
}

// * DONE
func ResetUserPassword() gin.HandlerFunc {
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
			utils.Error(c, utils.InternalServerError, "Can't find guest with Email id.")
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

func checkIdProofType(proof models.ID_Proof_Type) bool {
	if proof == "" {
		return false
	}

	if proof == models.Aadhar_Card {
		return true
	}

	if proof == models.Pan_Card {
		return true
	}

	if proof == models.PassPort {
		return true
	}

	if proof == models.Driving_License {
		return true
	}
	return false
}
