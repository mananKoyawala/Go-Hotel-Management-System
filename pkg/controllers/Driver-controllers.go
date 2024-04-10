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

var driverFolder = "driver"

// * DONE
func GetAllDrivers() gin.HandlerFunc {
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
			{Key: "drivers",
				Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
		}}}

		projectStage2 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "total_count", Value: 1},
			{Key: "driver", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$drivers"},
					{Key: "as", Value: "data"},
					{Key: "in", Value: bson.D{
						{Key: "driver_id", Value: "$$data.driver_id"},
						{Key: "first_name", Value: "$$data.first_name"},
						{Key: "last_name", Value: "$$data.last_name"},
						{Key: "age", Value: "$$data.age"},
						{Key: "phone", Value: "$$data.phone"},
						{Key: "email", Value: "$$data.email"},
						{Key: "gender", Value: "$$data.gender"},
						{Key: "car_company", Value: "$$data.car_company"},
						{Key: "car_model", Value: "$$data.car_model"},
						{Key: "car_number_plate", Value: "$$data.car_number_plate"},
						{Key: "availability", Value: "$$data.availability"},
						{Key: "pickup_location", Value: "$$data.pickup_location"},
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

		result, err := database.DriverCollection.Aggregate(ctx, mongo.Pipeline{
			groupStage, projectStage1, projectStage2,
		})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching drivers "+err.Error())
			return
		}

		var alldrivers []bson.M
		if err := result.All(ctx, &alldrivers); err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting the drivers "+err.Error())
			return
		}

		if len(alldrivers) == 0 {
			utils.Response(c, []interface{}{})
			return
		}

		utils.Response(c, alldrivers[0])
	}
}

// * DONE
func GetDriver() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var driver models.Driver
		id := c.Param("id")

		if err := database.DriverCollection.FindOne(ctx, bson.M{"driver_id": id}).Decode(&driver); err != nil {
			utils.Error(c, utils.BadRequest, "Can't find the driver with id.")
			return
		}

		utils.Response(c, driver)
	}
}

// * DONE
func DriverLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var driver models.Driver
		var foundDriver models.Driver

		driver.Email = c.PostForm("email")
		driver.Password = c.PostForm("password")

		// Check Email
		if err := database.DriverCollection.FindOne(ctx, bson.M{"email": driver.Email}).Decode(&foundDriver); err != nil {
			utils.Error(c, utils.NotFound, "Can't find driver with Email id.")
			return
		}

		// Verify Password
		msg, err := helpers.VerifyPassword(driver.Password, foundDriver.Password)
		if err != nil {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// Generate All Tokens
		token, refreshToken, _ := helpers.GenerateAllTokens(foundDriver.Email, foundDriver.First_Name, foundDriver.Last_Name, foundDriver.Driver_id, string(foundDriver.Access_Type))

		// Update Tokens
		if err := helpers.UpdateAllTokens(token, refreshToken, "driver_id", foundDriver.Driver_id); err != nil {
			utils.Error(c, utils.InternalServerError, "Error occured while updating tokens")
			return
		}

		foundDriver.Token = token
		foundDriver.Refresh_Token = refreshToken

		// Return as response
		utils.Response(c, foundDriver)
	}
}

// * DONE
func CreateDriver() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var driver models.Driver

		driver.First_Name = c.PostForm("first_name")
		driver.Last_Name = c.PostForm("last_name")
		driver.Email = c.PostForm("email")
		driver.Password = c.PostForm("password")
		driver.Gender = c.PostForm("gender")
		driver.Age, _ = strconv.Atoi(c.PostForm("age"))
		driver.Car_Company = c.PostForm("car_company")
		driver.Car_Model = c.PostForm("car_model")
		driver.Car_Number_Plate = c.PostForm("car_number_plate")
		driver.Phone, _ = strconv.Atoi(c.PostForm("phone"))
		driver.Salary, _ = strconv.ParseFloat(c.PostForm("salary"), 64)

		// validate details
		msg, val := validateDriverDetails(driver)
		if !val {
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

		// check driver exist with phone, car_number_late
		count, err := database.DriverCollection.CountDocuments(ctx, bson.M{"email": driver.Email, "car_number_plate": driver.Car_Number_Plate})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching the driver details.")
			return
		}

		if count > 0 {
			utils.Error(c, utils.Conflict, "Driver already exists with email and car number plate.")
			return
		}

		// hash password
		password, err := helpers.HashPassword(driver.Password)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Unable to hash password.")
			return
		}
		driver.Password = password

		// create id, timestamps, access_type
		driver.ID = primitive.NewObjectID()
		driver.Driver_id = driver.ID.Hex()
		driver.Status = models.Active
		driver.Availablity = models.Available
		driver.Access_Type = models.D_Acc
		driver.Created_at, _ = helpers.GetTime()
		driver.Updated_at, _ = helpers.GetTime()

		// generate token, refershtoken
		token, refreshToken, _ := helpers.GenerateAllTokens(driver.Email, driver.First_Name, driver.Last_Name, driver.Driver_id, string(driver.Access_Type))
		driver.Token = token
		driver.Refresh_Token = refreshToken

		// upload image
		name := strings.ReplaceAll(handler.Filename, " ", "")
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), name)
		url, err := imageupload.UploadService(file, driverFolder, filename)
		if err != nil {
			log.Println(err.Error())
			url = "https://i.ibb.co/y4BG3Kv/placeholder.jpg"
		}
		driver.Image = url

		// insert driver
		result, err := database.DriverCollection.InsertOne(ctx, driver)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't add driver.")
			return
		}

		// if success return
		utils.Response(c, result)
	}
}

// * DONE
func UpdateDriverDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var driver models.Driver
		id := c.Param("id")

		driver.First_Name = c.PostForm("first_name")
		driver.Last_Name = c.PostForm("last_name")
		driver.Gender = c.PostForm("gender")
		driver.Age, _ = strconv.Atoi(c.PostForm("age"))
		driver.Car_Company = c.PostForm("car_company")
		driver.Car_Model = c.PostForm("car_model")
		driver.Car_Number_Plate = c.PostForm("car_number_plate")
		driver.Phone, _ = strconv.Atoi(c.PostForm("phone"))
		driver.Salary, _ = strconv.ParseFloat(c.PostForm("salary"), 64)

		// Validate data
		msg, isVal := validateDriverUpdateDetails(driver)
		if !isVal {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// check driver exist or not
		count, err := database.DriverCollection.CountDocuments(ctx, bson.M{"driver_id": id})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while checking driver.")
			return
		}

		if !(count > 0) {
			utils.Error(c, utils.BadRequest, "Driver does not exist with id.")
			return
		}

		// Update data here

		var updateObj primitive.D

		updateObj = append(updateObj, bson.E{Key: "first_name", Value: driver.First_Name})

		updateObj = append(updateObj, bson.E{Key: "last_name", Value: driver.Last_Name})

		updateObj = append(updateObj, bson.E{Key: "gender", Value: driver.Gender})

		updateObj = append(updateObj, bson.E{Key: "phone", Value: driver.Phone})

		updateObj = append(updateObj, bson.E{Key: "age", Value: driver.Age})

		updateObj = append(updateObj, bson.E{Key: "car_company", Value: driver.Car_Company})

		updateObj = append(updateObj, bson.E{Key: "car_model", Value: driver.Car_Model})

		updateObj = append(updateObj, bson.E{Key: "car_number_plate", Value: driver.Car_Number_Plate})

		updateObj = append(updateObj, bson.E{Key: "salary", Value: driver.Salary})

		driver.Updated_at, _ = helpers.GetTime()
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: driver.Updated_at})

		filter := bson.M{"driver_id": id}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		// If success send success message
		_, err = database.DriverCollection.UpdateOne(ctx, filter, bson.D{
			{Key: "$set", Value: updateObj},
		}, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update driver details.")
			return
		}

		utils.Message(c, "Driver details updated successfully.")
	}
}

// * DONE
func UpdateDriverStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var driver models.Driver
		id := c.Param("id")

		if err := database.DriverCollection.FindOne(ctx, bson.M{"driver_id": id}).Decode(&driver); err != nil {
			utils.Error(c, utils.NotFound, "Can't find driver with ID.")
			return
		}

		var newStatus models.Status

		if driver.Status == models.Status(models.Active) {
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
		filter := bson.M{"driver_id": id}
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		_, err := database.DriverCollection.UpdateOne(ctx, filter, updateObj, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error updating driver status")
			return
		}

		utils.Message(c, "Driver status updated successfully.")
	}
}

// * DONE
func UpdateDriverAvailability() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var driver models.Driver
		id := c.Param("id")

		if err := database.DriverCollection.FindOne(ctx, bson.M{"driver_id": id}).Decode(&driver); err != nil {
			utils.Error(c, utils.NotFound, "Can't find driver with ID.")
			return
		}

		var newAvailability models.Availablity

		if driver.Availablity == models.Available {
			newAvailability = models.UnAvailable
		} else {
			newAvailability = models.Available
		}

		updated_at, _ := helpers.GetTime()
		updateObj := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "available", Value: newAvailability},
				{Key: "updated_at", Value: updated_at},
			}},
		}

		upsert := true
		filter := bson.M{"driver_id": id}
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		_, err := database.DriverCollection.UpdateOne(ctx, filter, updateObj, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error updating driver status")
			return
		}

		utils.Message(c, "Driver status updated successfully.")
	}
}

// * DONE
func ResetDriverPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var driver models.Driver

		driver.Email = c.PostForm("email")
		driver.Password = c.PostForm("password")

		// validate email
		if !utils.ValidateEmail(driver.Email) {
			utils.Error(c, utils.BadRequest, "Invalid email.")
			return
		}

		// check email exist
		count, err := database.DriverCollection.CountDocuments(ctx, bson.M{"email": driver.Email})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting driver details.")
			return
		}

		if !(count > 0) {
			utils.Error(c, utils.InternalServerError, "Can't find driver with Email id.")
			return
		}

		// validate password
		msg, val := utils.ValidatePassword(driver.Password)
		if !val {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// hash password
		password, err := helpers.HashPassword(driver.Password)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Unable to hash password.")
			return
		}

		// update password,timestamp
		driver.Password = password
		driver.Updated_at, _ = helpers.GetTime()

		// update details

		filter := bson.M{"email": driver.Email}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		updateObj := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "Password", Value: driver.Password},
				{Key: "updated_at", Value: driver.Updated_at},
			}},
		}

		_, err = database.DriverCollection.UpdateOne(ctx, filter, updateObj, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update driver password.")
			return
		}

		// if success return
		utils.Message(c, "Password updated successfully.")
	}
}

// * DONE
func DeleteDriver() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var driver models.Driver
		id := c.Param("id")

		if err := database.DriverCollection.FindOne(ctx, bson.M{"driver_id": id}).Decode(&driver); err != nil {
			utils.Error(c, utils.BadRequest, "Can't find driver with id")
			return
		}

		image := utils.GetTrimedUrl(driver.Image)
		if err := imageupload.DeleteService(image); err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete image."+err.Error())
			return
		}

		_, err := database.DriverCollection.DeleteOne(ctx, bson.M{"driver_id": id})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete driver with id.")
			return
		}

		utils.Message(c, "Driver is successfully deleted.")
	}
}

// * DONE
func UpdateDriverProfilePicture() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var driver models.Driver

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
		if err := database.DriverCollection.FindOne(ctx, bson.M{"driver_id": id}).Decode(&driver); err != nil {
			utils.Error(c, utils.BadRequest, "Can't find driver with id")
			return
		}
		log.Println(driver.Image)
		// delete file
		image := utils.GetTrimedUrl(driver.Image)
		if err := imageupload.DeleteService(image); err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete image."+err.Error())
			return
		}

		// upload new file
		name := strings.ReplaceAll(handler.Filename, " ", "")
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), name)
		url, err := imageupload.UploadService(file, driverFolder, filename)
		if err != nil {
			log.Println(err.Error())
			url = "https://i.ibb.co/y4BG3Kv/placeholder.jpg"
		}

		// update new uploaded file
		updated_at, _ := helpers.GetTime()
		filter := bson.M{"driver_id": id}
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

		_, err = database.DriverCollection.UpdateOne(ctx, filter, updateObj, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update image.")
			return
		}

		// if success return
		utils.Message(c, "Image updated successfully.")
	}
}

func validateDriverDetails(driver models.Driver) (string, bool) {

	if driver.First_Name == "" {
		return "First name required", false
	}

	if driver.Last_Name == "" {
		return "Last name required", false
	}

	if driver.Email == "" {
		return "Email address required", false
	} else if !utils.ValidateEmail(driver.Email) {
		return "Invalid email address", false
	}

	msg, val := utils.ValidatePassword(driver.Password)
	if !val {
		return msg, val
	}

	if driver.Age < 18 || driver.Age > 65 {
		return "Age must be between 18 to 65", false
	}

	if driver.Gender == "" {
		return "Gender is required", false
	}

	if driver.Car_Company == "" {
		return "Car_Compamy is required", false
	}

	if driver.Car_Model == "" {
		return "Car_Model is required", false
	}

	if driver.Car_Number_Plate == "" {
		return "Car_Number_Plate is required", false
	}

	if utils.CheckLength(driver.Phone, 10) {
		return "Phone number must be 10 digits", false
	}

	if !utils.IsNonNegative(int(driver.Salary)) {
		return "Salary must not 0 or negative", false
	}

	return "", true
}

func validateDriverUpdateDetails(driver models.Driver) (string, bool) {

	if driver.First_Name == "" {
		return "First name required", false
	}

	if driver.Last_Name == "" {
		return "Last name required", false
	}

	if driver.Age < 18 || driver.Age > 65 {
		return "Age must be between 18 to 65", false
	}

	if driver.Gender == "" {
		return "Gender is required", false
	}

	if driver.Car_Company == "" {
		return "Car_Compamy is required", false
	}

	if driver.Car_Model == "" {
		return "Car_Model is required", false
	}

	if driver.Car_Number_Plate == "" {
		return "Car_Number_Plate is required", false
	}

	if utils.CheckLength(driver.Phone, 10) {
		return "Phone number must be 10 digits", false
	}

	if !utils.IsNonNegative(int(driver.Salary)) {
		return "Salary must not 0 or negative", false
	}

	return "", true
}

// * DONE
func SearchDriverData() gin.HandlerFunc {
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
				bson.D{{Key: "car_company", Value: bson.D{{Key: "$regex", Value: search_string}, {Key: "$options", Value: "i"}}}},
				bson.D{{Key: "car_model", Value: bson.D{{Key: "$regex", Value: search_string}, {Key: "$options", Value: "i"}}}},
				bson.D{{Key: "car_number_plate", Value: bson.D{{Key: "$regex", Value: search_string}, {Key: "$options", Value: "i"}}}},
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
			{Key: "drivers",
				Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
		}}}

		projectStage2 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "total_count", Value: 1},
			{Key: "driver", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$drivers"},
					{Key: "as", Value: "data"},
					{Key: "in", Value: bson.D{
						{Key: "driver_id", Value: "$$data.driver_id"},
						{Key: "first_name", Value: "$$data.first_name"},
						{Key: "last_name", Value: "$$data.last_name"},
						{Key: "age", Value: "$$data.age"},
						{Key: "phone", Value: "$$data.phone"},
						{Key: "email", Value: "$$data.email"},
						{Key: "gender", Value: "$$data.gender"},
						{Key: "car_company", Value: "$$data.car_company"},
						{Key: "car_model", Value: "$$data.car_model"},
						{Key: "car_number_plate", Value: "$$data.car_number_plate"},
						{Key: "availability", Value: "$$data.availability"},
						{Key: "pickup_location", Value: "$$data.pickup_location"},
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

		result, err := database.DriverCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage1, projectStage2,
		})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching driver  "+err.Error())
			return
		}

		var allDrivers []bson.M
		if err := result.All(ctx, &allDrivers); err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting the drivers "+err.Error())
			return
		}
		if len(allDrivers) == 0 {
			utils.Response(c, []interface{}{})
			return
		}
		utils.Response(c, allDrivers[0])
	}
}

// * DONE
func FilterDriver() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()

		salary, _ := strconv.ParseFloat(c.PostForm("salary"), 64)
		availability := c.PostForm("availability")
		status := models.Status(c.PostForm("status"))
		salaryOperator := c.PostForm("salaryOperator")

		// clean the data
		if salary <= 0.0 {
			salary = 0
		}

		if salaryOperator == "" || (salaryOperator != "$gt" && salaryOperator != "$eq" && salaryOperator != "$lt") {
			salaryOperator = "$gt"
		}

		if availability == "" || (availability != "available" && availability != "unavailable") {
			availability = "available"
		}

		if status == "" || (status != "active" && status != "inactive") {
			status = "active"
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
			{Key: "availability", Value: availability},
			{Key: "status", Value: status},
			{Key: "salary", Value: bson.D{
				{Key: salaryOperator, Value: salary},
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
			{Key: "drivers",
				Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
		}}}

		projectStage2 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "total_count", Value: 1},
			{Key: "driver", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$drivers"},
					{Key: "as", Value: "data"},
					{Key: "in", Value: bson.D{
						{Key: "driver_id", Value: "$$data.driver_id"},
						{Key: "first_name", Value: "$$data.first_name"},
						{Key: "last_name", Value: "$$data.last_name"},
						{Key: "age", Value: "$$data.age"},
						{Key: "phone", Value: "$$data.phone"},
						{Key: "email", Value: "$$data.email"},
						{Key: "gender", Value: "$$data.gender"},
						{Key: "car_company", Value: "$$data.car_company"},
						{Key: "car_model", Value: "$$data.car_model"},
						{Key: "car_number_plate", Value: "$$data.car_number_plate"},
						{Key: "availability", Value: "$$data.availability"},
						{Key: "pickup_location", Value: "$$data.pickup_location"},
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

		result, err := database.DriverCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage1, projectStage2,
		})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching drivers "+err.Error())
			return
		}

		var alldrivers []bson.M
		if err := result.All(ctx, &alldrivers); err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting the drivers "+err.Error())
			return
		}

		if len(alldrivers) == 0 {
			utils.Response(c, []interface{}{})
			return
		}

		utils.Response(c, alldrivers[0])
	}
}
