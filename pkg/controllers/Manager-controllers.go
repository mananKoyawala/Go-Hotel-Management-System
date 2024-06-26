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
	imageupload "github.com/mananKoyawala/hotel-management-system/pkg/service/image-upload"
	"github.com/mananKoyawala/hotel-management-system/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var managerFolder = "manager"

// * DONE
func GetManagers() gin.HandlerFunc {
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
			{Key: "managers",
				Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
		}}}

		projectStage2 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "total_count", Value: 1},
			{Key: "manager", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$managers"},
					{Key: "as", Value: "data"},
					{Key: "in", Value: bson.D{
						{Key: "manager_id", Value: "$$data.manager_id"},
						{Key: "first_name", Value: "$$data.first_name"},
						{Key: "last_name", Value: "$$data.last_name"},
						{Key: "age", Value: "$$data.age"},
						{Key: "phone", Value: "$$data.phone"},
						{Key: "email", Value: "$$data.email"},
						{Key: "gender", Value: "$$data.gender"},
						{Key: "salary", Value: "$$data.salary"},
						{Key: "status", Value: "$$data.status"},
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

		result, err := database.ManagerCollection.Aggregate(ctx, mongo.Pipeline{
			groupStage, projectStage1, projectStage2,
		})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching managers "+err.Error())
			return
		}

		var allManagers []bson.M
		if err := result.All(ctx, &allManagers); err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting the managers "+err.Error())
			return
		}

		if len(allManagers) == 0 {
			utils.Response(c, []interface{}{})
			return
		}

		utils.Response(c, allManagers[0])
	}
}

// * DONE
func GetManager() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := helpers.GetContext()
		defer cancel()
		var manager models.Manager
		id := c.Param("id")

		if err := database.ManagerCollection.FindOne(ctx, bson.M{"manager_id": id}).Decode(&manager); err != nil {
			utils.Error(c, utils.NotFound, "Can't find manager with id")
			return
		}

		utils.Response(c, manager)
	}
}

// * DONE
func ManagerLoign() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var manager models.Manager
		var foundManager models.Manager

		manager.Email = c.PostForm("email")
		manager.Password = c.PostForm("password")

		// Check Email
		if err := database.ManagerCollection.FindOne(ctx, bson.M{"email": manager.Email}).Decode(&foundManager); err != nil {
			utils.Error(c, utils.NotFound, "Can't find manager with Email id.")
			return
		}

		// Verify Password
		msg, err := helpers.VerifyPassword(manager.Password, foundManager.Password)
		if err != nil {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// Generate All Tokens
		token, refreshToken, _ := helpers.GenerateAllTokens(foundManager.Email, foundManager.First_Name, foundManager.Last_Name, foundManager.Manager_id, string(foundManager.Access_Type))

		// Update Tokens
		if err := helpers.UpdateAllTokens(token, refreshToken, "manager_id", foundManager.Manager_id); err != nil {
			utils.Error(c, utils.InternalServerError, "Error occured while updating tokens")
			return
		}

		foundManager.Token = token
		foundManager.Refresh_Token = refreshToken

		// Return as response
		utils.Response(c, foundManager)
	}
}

// * DONE
func CreateManager() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var manager models.Manager

		manager.First_Name = c.PostForm("first_name")
		manager.Last_Name = c.PostForm("last_name")
		manager.Age, _ = strconv.Atoi(c.PostForm("age"))
		manager.Phone, _ = strconv.Atoi(c.PostForm("phone"))
		manager.Email = c.PostForm("email")
		manager.Password = c.PostForm("password")
		manager.Gender = c.PostForm("gender")
		manager.Salary, _ = strconv.ParseFloat(c.PostForm("salary"), 64)
		manager.Aadhar_Number, _ = strconv.Atoi(c.PostForm("aadhar_number"))

		// Validate details
		msg, isVal := validateManager(manager)
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

		// Check the email and phone already exists
		count1, emailerr := database.ManagerCollection.CountDocuments(ctx, bson.M{"email": manager.Email})
		if emailerr != nil {
			utils.Error(c, utils.InternalServerError, "Error while checking Email")
			return
		}

		count2, phoneerr := database.ManagerCollection.CountDocuments(ctx, bson.M{"phone": manager.Phone})
		if phoneerr != nil {
			utils.Error(c, utils.InternalServerError, "Error while checking Phone")
			return
		}

		if count1 > 0 || count2 > 0 {
			utils.Error(c, utils.Conflict, "User email or phone already exist.")
			return
		}

		// Hash Password
		password, err := helpers.HashPassword(manager.Password)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't generate password hash.")
			return
		}
		manager.Password = password

		// Generate ID and Timestamp
		manager.ID = primitive.NewObjectID()
		manager.Manager_id = manager.ID.Hex()
		manager.Created_at, _ = helpers.GetTime()
		manager.Updated_at, _ = helpers.GetTime()
		manager.Access_Type = models.M_Acc
		manager.Status = models.Active

		// Generate Tokens
		token, refreshToken, _ := helpers.GenerateAllTokens(manager.Email, manager.First_Name, manager.Last_Name, manager.Manager_id, string(manager.Access_Type))
		manager.Token = token
		manager.Refresh_Token = refreshToken

		// upload image
		name := strings.ReplaceAll(handler.Filename, " ", "")
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), name)
		url, err := imageupload.UploadService(file, managerFolder, filename)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't uplaod the image.")
			return
		}
		manager.Image = url

		// Insert Manager
		result, err := database.ManagerCollection.InsertOne(ctx, manager)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't insert manager.")
			return
		}

		// If of reture response
		utils.Response(c, result)
	}
}

// * DONE
func UpdateManagerDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var manager models.Manager
		id := c.Param("id")

		manager.First_Name = c.PostForm("first_name")
		manager.Last_Name = c.PostForm("last_name")
		manager.Age, _ = strconv.Atoi(c.PostForm("age"))
		manager.Gender = c.PostForm("gender")
		manager.Salary, _ = strconv.ParseFloat(c.PostForm("salary"), 64)
		manager.Phone, _ = strconv.Atoi(c.PostForm("phone"))

		// Validate data
		msg, val := validateUpdateManager(manager)
		if !val {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// check manager if exist or not
		count, err := database.ManagerCollection.CountDocuments(ctx, bson.M{"manager_id": id})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching manager.")
			return
		}

		if !(count > 0) {
			utils.Error(c, utils.NotFound, "Can't find manager with id.")
			return
		}

		// Update data here

		var updateObj primitive.D

		if manager.First_Name != "" {
			updateObj = append(updateObj, bson.E{Key: "first_name", Value: manager.First_Name})
		}

		if manager.Last_Name != "" {
			updateObj = append(updateObj, bson.E{Key: "last_name", Value: manager.Last_Name})
		}

		if manager.Age != 0 {
			updateObj = append(updateObj, bson.E{Key: "age", Value: manager.Age})
		}

		if manager.Gender != "" {
			updateObj = append(updateObj, bson.E{Key: "gender", Value: manager.Gender})
		}

		if manager.Salary != 0.0 {
			updateObj = append(updateObj, bson.E{Key: "salary", Value: manager.Salary})
		}

		if manager.Phone != 0.0 {
			updateObj = append(updateObj, bson.E{Key: "phone", Value: manager.Phone})
		}

		manager.Updated_at, _ = helpers.GetTime()
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: manager.Updated_at})

		filter := bson.M{"manager_id": id}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		// If success send success message

		_, err = database.ManagerCollection.UpdateOne(ctx, filter, bson.D{
			{Key: "$set", Value: updateObj},
		}, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update manager details")
			return
		}
		utils.Message(c, "Manager updated successfully.")
	}
}

// * DONE
func UpdateManagerStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var manager models.Manager
		id := c.Param("id")

		if err := database.ManagerCollection.FindOne(ctx, bson.M{"manager_id": id}).Decode(&manager); err != nil {
			utils.Error(c, utils.NotFound, "Can't find manager with ID.")
			return
		}

		var newStatus models.Status

		if manager.Status == models.Status(models.Active) {
			newStatus = models.Inactive
		} else {
			newStatus = models.Active
		}

		updated_at, _ := helpers.GetTime()
		updateObj := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "status", Value: newStatus},
				{Key: "updated_at", Value: updated_at},
			}},
		}

		upsert := true
		filter := bson.M{"manager_id": id}
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		_, err := database.ManagerCollection.UpdateOne(ctx, filter, updateObj, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error updating manager status")
			return
		}

		utils.Message(c, "Manager status updated successfully.")
	}
}

// * DONE
func DeleteManager() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var manager models.Manager
		id := c.Param("id")

		if err := database.ManagerCollection.FindOne(ctx, bson.M{"manager_id": id}).Decode(&manager); err != nil {
			utils.Error(c, utils.NotFound, "Can't find manager with id")
			return
		}

		image := utils.GetTrimedUrl(manager.Image)
		if err := imageupload.DeleteService(image); err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete image."+err.Error())
			return
		}

		_, err := database.ManagerCollection.DeleteOne(ctx, bson.M{"manager_id": id})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete manager.")
			return
		}
		utils.Message(c, "Manager deleted successfully.")
	}
}

// * DONE
func ResetManagerPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var manager models.Manager

		manager.Email = c.PostForm("email")
		manager.Password = c.PostForm("password")

		// validate email & password
		if !utils.ValidateEmail(manager.Email) {
			utils.Error(c, utils.BadRequest, "Invalid Email Address")
			return
		}

		msg, val := utils.ValidatePassword(manager.Password)
		if !val {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// check manager with email exist or not
		count, err := database.ManagerCollection.CountDocuments(ctx, bson.M{"email": manager.Email})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting manager details.")
			return
		}

		if !(count > 0) {
			utils.Error(c, utils.NotFound, "Can't find manager with id.")
			return
		}

		// hash password and update timestamp
		password, err := helpers.HashPassword(manager.Password)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't generate hash password")
			return
		}
		manager.Password = password
		manager.Updated_at, _ = helpers.GetTime()

		// update details
		filter := bson.M{"email": manager.Email}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		updateObj := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "Password", Value: manager.Password},
				{Key: "updated_at", Value: manager.Updated_at},
			}},
		}

		_, err = database.ManagerCollection.UpdateOne(ctx, filter, updateObj, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update password")
			return
		}

		// return if success
		utils.Message(c, "Password was updated")
	}
}

// * DONE
func UpdateManagerProfilePicture() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var manager models.Manager

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
		if err := database.ManagerCollection.FindOne(ctx, bson.M{"manager_id": id}).Decode(&manager); err != nil {
			utils.Error(c, utils.NotFound, "Can't find manager with id")
			return
		}
		log.Println(manager.Image)
		// delete file
		image := utils.GetTrimedUrl(manager.Image)
		if err := imageupload.DeleteService(image); err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete image."+err.Error())
			return
		}

		// upload new file
		name := strings.ReplaceAll(handler.Filename, " ", "")
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), name)
		url, err := imageupload.UploadService(file, managerFolder, filename)
		if err != nil {
			log.Println(err.Error())
			url = "https://i.ibb.co/y4BG3Kv/placeholder.jpg"
		}

		// update new uploaded file
		updated_at, _ := helpers.GetTime()
		filter := bson.M{"manager_id": id}
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

		_, err = database.ManagerCollection.UpdateOne(ctx, filter, updateObj, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update image.")
			return
		}

		// if success return
		utils.Message(c, "Image updated successfully.")
	}
}

// * DONE
func validateManager(manager models.Manager) (string, bool) {
	if manager.First_Name == "" {
		return "First name is required", false
	}

	if manager.Last_Name == "" {
		return "Last name is required", false
	}

	if manager.Age < 18 || manager.Age > 65 {
		return "Age must be between 18 to 65", false
	}

	if utils.CheckLength(manager.Phone, 10) {
		return "Phone number must be 10 digits", false
	}

	if manager.Email == "" {
		return "Email address required", false
	} else if !utils.ValidateEmail(manager.Email) {
		return "Invalid email address", false
	}

	msg, val := utils.ValidatePassword(manager.Password)
	if !val {
		return msg, false
	}

	if manager.Gender == "" {
		return "Gender is required", false
	}

	if !utils.IsNonNegative(int(manager.Salary)) {
		return "Salary must not 0 or negative", false
	}

	if utils.CheckLength(manager.Aadhar_Number, 12) {
		return "Aadhar Number must be of 12 digits", false
	}

	return "", true
}

// * DONE
func validateUpdateManager(manager models.Manager) (string, bool) {
	if manager.First_Name == "" {
		return "First name is required", false
	}

	if manager.Last_Name == "" {
		return "Last name is required", false
	}

	if manager.Age < 18 || manager.Age > 65 {
		return "Age must be between 18 to 65", false
	}

	if utils.CheckLength(manager.Phone, 10) {
		return "Phone number must be 10 digits", false
	}

	if manager.Gender == "" {
		return "Gender is required", false
	}

	if !utils.IsNonNegative(int(manager.Salary)) {
		return "Salary must not 0 or negative", false
	}

	return "", true
}

// * DONE
func SearchManagerData() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()

		search_string := c.PostForm("search")

		if search_string == "" {
			utils.Error(c, utils.BadRequest, "Please provide a search string.")
			return
		}

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

		matchStage := bson.D{{Key: "$match", Value: bson.D{
			{Key: "$or", Value: bson.A{
				bson.D{{Key: "first_name", Value: bson.D{{Key: "$regex", Value: search_string}, {Key: "$options", Value: "i"}}}},
				bson.D{{Key: "last_name", Value: bson.D{{Key: "$regex", Value: search_string}, {Key: "$options", Value: "i"}}}},
				bson.D{{Key: "gender", Value: bson.D{{Key: "$regex", Value: search_string}, {Key: "$options", Value: "i"}}}},
			}},
		}}}

		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "null"},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}}}

		projectStage1 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "managers",
				Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
		}}}

		projectStage2 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "total_count", Value: 1},
			{Key: "manager", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$managers"},
					{Key: "as", Value: "data"},
					{Key: "in", Value: bson.D{
						{Key: "manager_id", Value: "$$data.manager_id"},
						{Key: "first_name", Value: "$$data.first_name"},
						{Key: "last_name", Value: "$$data.last_name"},
						{Key: "age", Value: "$$data.age"},
						{Key: "phone", Value: "$$data.phone"},
						{Key: "email", Value: "$$data.email"},
						{Key: "gender", Value: "$$data.gender"},
						{Key: "salary", Value: "$$data.salary"},
						{Key: "status", Value: "$$data.status"},
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

		result, err := database.ManagerCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage1, projectStage2,
		})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching manager  "+err.Error())
			return
		}

		var allManagers []bson.M
		if err := result.All(ctx, &allManagers); err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting the managers "+err.Error())
			return
		}
		if len(allManagers) == 0 {
			utils.Response(c, []interface{}{})
			return
		}
		utils.Response(c, allManagers[0])
	}
}

// * DONE
func FilterManager() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var manager models.Manager

		manager.Age, _ = strconv.Atoi(c.PostForm("age"))
		manager.Salary, _ = strconv.ParseFloat(c.PostForm("salary"), 64)
		manager.Status = models.Status(c.PostForm("status"))
		ageOperator := c.PostForm("ageOperator")
		salaryOperator := c.PostForm("salaryOperator")

		// Cleaning the input
		if manager.Age <= 0 {
			manager.Age = 0
		}
		if manager.Salary <= 0.0 {
			manager.Salary = 0
		}
		if ageOperator == "" || (ageOperator != "$gt" && ageOperator != "$eq" && ageOperator != "$lt") {
			ageOperator = "$gt"
		}

		if salaryOperator == "" || (salaryOperator != "$gt" && salaryOperator != "$eq" && salaryOperator != "$lt") {
			salaryOperator = "$gt"
		}

		if manager.Status == "" || (manager.Status != "active" && manager.Status != "inactive") {
			manager.Status = "active"
		}

		// log.Println(manager.Age)
		// log.Println(manager.Salary)
		// log.Println(manager.Status)
		// log.Println(ageOperator)
		// log.Println(salaryOperator)

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

		matchStage := bson.D{{Key: "$match", Value: bson.D{
			{Key: "status", Value: manager.Status},
			{Key: "salary", Value: bson.D{
				{Key: salaryOperator, Value: manager.Salary},
			}},
			{Key: "age", Value: bson.D{
				{Key: ageOperator, Value: manager.Age},
			}},
		}}}

		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "null"},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}}}

		projectStage1 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "managers",
				Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
		}}}

		projectStage2 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "total_count", Value: 1},
			{Key: "manager", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$managers"},
					{Key: "as", Value: "data"},
					{Key: "in", Value: bson.D{
						{Key: "manager_id", Value: "$$data.manager_id"},
						{Key: "first_name", Value: "$$data.first_name"},
						{Key: "last_name", Value: "$$data.last_name"},
						{Key: "age", Value: "$$data.age"},
						{Key: "phone", Value: "$$data.phone"},
						{Key: "email", Value: "$$data.email"},
						{Key: "gender", Value: "$$data.gender"},
						{Key: "salary", Value: "$$data.salary"},
						{Key: "status", Value: "$$data.status"},
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

		result, err := database.ManagerCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage1, projectStage2,
		})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching manager  "+err.Error())
			return
		}

		var allManagers []bson.M
		if err := result.All(ctx, &allManagers); err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting the managers "+err.Error())
			return
		}
		if len(allManagers) == 0 {
			utils.Response(c, []interface{}{})
			return
		}
		utils.Response(c, allManagers[0])
	}
}
