package controllers

import (
	"context"
	"fmt"
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

var roomFolder = "room"

// * DONE
func GetRoom() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := helpers.GetContext()
		defer cancel()
		var room models.Room
		id := c.Param("id")

		if err := database.RoomCollection.FindOne(ctx, bson.M{"room_id": id}).Decode(&room); err != nil {
			utils.Error(c, utils.NotFound, "Can't find the room with id.")
			return
		}

		utils.Response(c, room)
	}
}

// * DONE
func GetRooms() gin.HandlerFunc {
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
			{Key: "rooms",
				Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
		}}}

		projectStage2 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "total_count", Value: 1},
			{Key: "room", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$rooms"},
					{Key: "as", Value: "data"},
					{Key: "in", Value: bson.D{
						{Key: "room_id", Value: "$$data.room_id"},
						{Key: "branch_id", Value: "$$data.branch_id"},
						{Key: "room_number", Value: "$$data.room_number"},
						{Key: "room_type", Value: "$$data.room_type"},
						{Key: "room_availability", Value: "$$data.room_availability"},
						{Key: "cleaning_status", Value: "$$data.cleaning_status"},
						{Key: "price", Value: "$$data.price"},
						{Key: "capacity", Value: "$$data.capacity"},
						{Key: "images", Value: "$$data.images"},
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

		result, err := database.RoomCollection.Aggregate(ctx, mongo.Pipeline{
			groupStage, projectStage1, projectStage2,
		})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching rooms "+err.Error())
			return
		}

		var allRooms []bson.M
		if err := result.All(ctx, &allRooms); err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting the room "+err.Error())
			return
		}

		if len(allRooms) == 0 {
			utils.Response(c, []interface{}{})
			return
		}

		utils.Response(c, allRooms[0])
	}
}

// * DONE
func GetRoomsByBranch() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()

		branch_id := c.Param("id")

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

		matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "branch_id", Value: branch_id}}}}

		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "null"},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}}}

		projectStage1 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "rooms",
				Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
		}}}

		projectStage2 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "total_count", Value: 1},
			{Key: "room", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$rooms"},
					{Key: "as", Value: "data"},
					{Key: "in", Value: bson.D{
						{Key: "room_id", Value: "$$data.room_id"},
						{Key: "branch_id", Value: "$$data.branch_id"},
						{Key: "room_number", Value: "$$data.room_number"},
						{Key: "room_type", Value: "$$data.room_type"},
						{Key: "room_availability", Value: "$$data.room_availability"},
						{Key: "cleaning_status", Value: "$$data.cleaning_status"},
						{Key: "price", Value: "$$data.price"},
						{Key: "images", Value: "$$data.images"},
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

		result, err := database.RoomCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage1, projectStage2,
		})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching rooms "+err.Error())
			return
		}

		var allRooms []bson.M
		if err := result.All(ctx, &allRooms); err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting the room "+err.Error())
			return
		}

		if len(allRooms) == 0 {
			utils.Response(c, []interface{}{})
			return
		}

		utils.Response(c, allRooms[0])
	}
}

// * DONE
func CreateRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var room models.Room
		var branch models.Branch

		room.Branch_id = c.PostForm("branch_id")
		room.Room_Number, _ = strconv.Atoi(c.PostForm("room_number"))
		room.Room_Type = c.PostForm("room_type")
		room.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)
		room.Capacity, _ = strconv.Atoi(c.PostForm("capacity"))

		// Validate details
		msg, isVal := validateRoomDetails(room)
		if !isVal {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// Check branch exist or not
		if err := database.BranchCollection.FindOne(ctx, bson.M{"branch_id": room.Branch_id}).Decode(&branch); err != nil {
			utils.Error(c, utils.NotFound, "Branch doesn't exist")
			return
		}

		// Check the room number exists
		count, rErr := database.RoomCollection.CountDocuments(ctx, bson.M{"room_number": room.Room_Number, "branch_id": room.Branch_id})
		if rErr != nil {
			utils.Error(c, utils.InternalServerError, "Error while checking Room Number")
			return
		}

		if count > 0 {
			utils.Error(c, utils.Conflict, "Room Number already Exist with Branch id.")
			return
		}

		// Check if no files are provided
		if len(c.Request.MultipartForm.File["file"]) <= 0 {
			utils.Error(c, utils.BadRequest, "No files provided.")
			return
		}

		// Check if more than 5 files are uploaded
		if len(c.Request.MultipartForm.File["file"]) > maxFiles {
			message := fmt.Sprintf("Only %d files can be uploaded.", maxFiles)
			utils.Error(c, utils.BadRequest, message)
			return
		}

		// check that all the provided files are .png .jpeg or .jpg
		for _, fileHeader := range c.Request.MultipartForm.File["file"] {
			file, err := fileHeader.Open()
			if err != nil {
				utils.Error(c, utils.BadRequest, "File was not provided or Invalid file.")
				return
			}
			defer file.Close()

			// Check the image file is .png , .jpg , .jpeg
			ext := filepath.Ext(fileHeader.Filename)
			if ext != ".jpeg" && ext != ".jpg" && ext != ".png" {
				utils.Error(c, utils.BadRequest, "Invalid Image file format. Only JPEG, JPG, or PNG files are allowed.")
				return
			}
		}

		// Generate ID and Timestamp
		room.ID = primitive.NewObjectID()
		room.Room_id = room.ID.Hex()
		room.Room_Availability = string(models.Room_Available)
		room.Cleaning_Status = string(models.Cleaned)
		room.Created_at, _ = helpers.GetTime()
		room.Updated_at, _ = helpers.GetTime()

		// upload files here
		for _, fileHeader := range c.Request.MultipartForm.File["file"] {
			name := strings.ReplaceAll(fileHeader.Filename, " ", "")
			file, _ := fileHeader.Open()
			filename := fmt.Sprintf("%d_%s", time.Now().Unix(), name)
			url, err := imageupload.UploadService(file, roomFolder, filename)
			if err != nil {
				utils.Error(c, utils.InternalServerError, "Can't uplaod the image.")
				return
			}
			room.Images = append(room.Images, url)
		}

		// Insert Room Details
		result, err := database.RoomCollection.InsertOne(ctx, room)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't insert room.")
			return
		}

		branch.Total_Rooms++
		branch.Updated_at, _ = helpers.GetTime()
		updateObj := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "total_rooms", Value: branch.Total_Rooms},
				{Key: "updated_at", Value: branch.Updated_at},
			}},
		}
		_, err = database.BranchCollection.UpdateOne(ctx, bson.M{"branch_id": branch.Branch_id}, updateObj)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update total_rooms.")
			return
		}
		// If of reture response
		utils.Response(c, result)
	}
}

// * DONE
func UpdateRoomDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Update all things
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var room models.Room
		id := c.Param("id")

		room.Room_Number, _ = strconv.Atoi(c.PostForm("room_number"))
		room.Room_Type = c.PostForm("room_type")
		room.Cleaning_Status = c.PostForm("cleaning_status")
		room.Room_Availability = c.PostForm("room_availability")
		room.Price, _ = strconv.ParseFloat(c.PostForm("price"), 64)

		// Validate data
		msg, isval := validateUpdateRoomDetails(room)
		if !isval {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// Update data here

		var updateObj primitive.D

		updateObj = append(updateObj, bson.E{Key: "room_number", Value: room.Room_Number})
		updateObj = append(updateObj, bson.E{Key: "room_type", Value: room.Room_Type})
		updateObj = append(updateObj, bson.E{Key: "room_availability", Value: room.Room_Availability})
		updateObj = append(updateObj, bson.E{Key: "cleaning_status", Value: room.Cleaning_Status})
		updateObj = append(updateObj, bson.E{Key: "price", Value: room.Price})

		room.Updated_at, _ = helpers.GetTime()
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: room.Updated_at})

		filter := bson.M{"room_id": id}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		// If success send success message

		_, err := database.RoomCollection.UpdateOne(ctx, filter, bson.D{
			{Key: "$set", Value: updateObj},
		}, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update room details")
			return
		}
		utils.Message(c, "Room is updated.")
	}
}

// * DONE
func DeleteRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var room models.Room
		var branch models.Branch
		id := c.Param("id")

		if err := database.RoomCollection.FindOne(ctx, bson.M{"room_id": id}).Decode(&room); err != nil {
			utils.Error(c, utils.NotFound, "Can't find room with id")
			return
		}

		for _, i := range room.Images {
			image := utils.GetTrimedUrl(i)
			if err := imageupload.DeleteService(image); err != nil {
				utils.Error(c, utils.InternalServerError, "Can't delete image."+err.Error())
				return
			}
		}

		_, err := database.RoomCollection.DeleteOne(ctx, bson.M{"room_id": id})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete room.")
			return
		}

		if err := database.BranchCollection.FindOne(ctx, bson.M{"branch_id": room.Branch_id}).Decode(&branch); err != nil {
			utils.Error(c, utils.NotFound, "Can't find branch with id")
			return
		}

		branch.Total_Rooms--
		branch.Updated_at, _ = helpers.GetTime()
		updateObj := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "total_rooms", Value: branch.Total_Rooms},
				{Key: "updated_at", Value: branch.Updated_at},
			}},
		}
		_, err = database.BranchCollection.UpdateOne(ctx, bson.M{"branch_id": branch.Branch_id}, updateObj)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update total_rooms.")
			return
		}

		utils.Message(c, "Room deleted successfully.")
	}
}

// * DONE
func validateRoomDetails(room models.Room) (string, bool) {

	if room.Branch_id == "" {
		return "Branch id needed", false
	}

	if !utils.IsNonNegative(room.Room_Number) {
		return "Room number must be integer.", false
	}

	if !checkRoomType(room.Room_Type) {
		return "Room type must be single_bed, double_bed or suite.", false
	}

	if !utils.IsNonNegative(int(room.Price)) {
		return "Room price must be integer.", false
	}

	if !utils.IsNonNegative(room.Capacity) {
		return "Room capacity must be integer.", false
	}

	return "", true
}

// * DONE
func validateUpdateRoomDetails(room models.Room) (string, bool) {

	if !utils.IsNonNegative(room.Room_Number) {
		return "Room number must be integer.", false
	}

	if !checkRoomType(room.Room_Type) {
		return "Room type must be single_bed, double_bed or suite.", false
	}

	if !checkRoomStatus(room.Cleaning_Status) {
		return "Room cleaning status must be cleaned, dirty or inprogress", false
	}

	if !checkRoomAvailability(room.Room_Availability) {
		return "Room availability must be available or occupied.", false
	}

	if !utils.IsNonNegative(int(room.Price)) {
		return "Room price must be integer.", false
	}

	return "", true
}

// * DONE
func RoomAddImage() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := helpers.GetContext()
		defer cancel()
		var room models.Room
		id := c.Param("id")

		if err := database.RoomCollection.FindOne(ctx, bson.M{"room_id": id}).Decode(&room); err != nil {
			utils.Error(c, utils.NotFound, "Can't find room with id.")
			return
		}

		if len(room.Images) >= 5 {
			utils.Error(c, utils.BadRequest, "Maximum 5 images already uploaded.")
			return
		}

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

		// upload image
		name := strings.ReplaceAll(handler.Filename, " ", "")
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), name)
		url, err := imageupload.UploadService(file, roomFolder, filename)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't uplaod the image.")
			return
		}
		room.Images = append(room.Images, url)
		updated_at, _ := helpers.GetTime()

		filter := bson.M{"room_id": id}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}
		updateObj := bson.D{{Key: "$set",
			Value: bson.D{
				{Key: "images", Value: room.Images},
				{Key: "updated_at", Value: updated_at},
			}}}

		_, err = database.RoomCollection.UpdateOne(ctx, filter, updateObj, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update image")
			return
		}
		utils.Message(c, "Image updated successfully.")
	}
}

// * DONE
func RoomRemoveImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var room models.Room
		id := c.Param("id")
		imageUrl := c.PostForm("image")

		if err := database.RoomCollection.FindOne(ctx, bson.M{"room_id": id}).Decode(&room); err != nil {
			utils.Error(c, utils.NotFound, "Can't find room with id.")
			return
		}

		if len(room.Images) <= 0 {
			utils.Error(c, utils.BadRequest, "Room doesn't have any images.")
			return
		}

		if imageUrl == "" {
			utils.Error(c, utils.BadRequest, "Please provide image url.")
			return
		}

		found := false
		for i, img := range room.Images {
			if img == imageUrl {
				found = true
				// remove imageUrl from room.Images
				room.Images = append(room.Images[:i], room.Images[i+1:]...)
				break
			}
		}

		if !found {
			utils.Error(c, utils.BadRequest, "Can't find image in room.")
			return
		}

		image := utils.GetTrimedUrl(imageUrl)
		if err := imageupload.DeleteService(image); err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete image."+err.Error())
			return
		}

		updated_at, _ := helpers.GetTime()

		filter := bson.M{"room_id": id}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}
		updateObj := bson.D{{Key: "$set",
			Value: bson.D{
				{Key: "images", Value: room.Images},
				{Key: "updated_at", Value: updated_at},
			}}}

		_, err := database.RoomCollection.UpdateOne(ctx, filter, updateObj, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update image")
			return
		}
		utils.Message(c, "Image deleted successfully.")
	}
}

func checkRoomType(roomType string) bool {

	if roomType == "" {
		return false
	}

	if models.Room_Type(roomType) == models.Single_Bad {
		return true
	}

	if models.Room_Type(roomType) == models.Double_Bad {
		return true
	}

	if models.Room_Type(roomType) == models.Suite {
		return true
	}

	return false
}

func checkRoomStatus(status string) bool {

	if status == "" {
		return false
	}

	if models.Cleaning_Status(status) == models.Cleaned {
		return true
	}

	if models.Cleaning_Status(status) == models.Dirty {
		return true
	}

	if models.Cleaning_Status(status) == models.InProgress {
		return true
	}

	return true
}

func checkRoomAvailability(avail string) bool {

	if avail == "" {
		return false
	}

	if models.Room_Availability(avail) == models.Room_Available {
		return true
	}

	if models.Room_Availability(avail) == models.Room_Unavailable {
		return true
	}

	return true
}

// NOT USED FOR THIS API
func IncreaseRoomCapacity() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()

		updated_at, _ := helpers.GetTime()
		filter1 := bson.M{"room_type": models.Single_Bad}
		filter2 := bson.M{"room_type": models.Double_Bad}
		filter3 := bson.M{"room_type": models.Suite}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		updateObj1 := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "capacity", Value: 1},
				{Key: `updated_at`, Value: updated_at},
			}},
		}

		updateObj2 := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "capacity", Value: 2},
				{Key: `updated_at`, Value: updated_at},
			}},
		}

		updateObj3 := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "capacity", Value: 5},
				{Key: `updated_at`, Value: updated_at},
			}},
		}
		_, err := database.RoomCollection.UpdateMany(ctx, filter1, updateObj1, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update the status for single bed.")
			return
		}

		_, err = database.RoomCollection.UpdateMany(ctx, filter2, updateObj2, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update the status for double bed.")
			return
		}

		_, err = database.RoomCollection.UpdateMany(ctx, filter3, updateObj3, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update the status for suite.")
			return
		}

		utils.Message(c, "Updated")
	}
}

func UpdateRoomAvailability(id string, avail models.Room_Availability) error {

	updated_at, _ := helpers.GetTime()
	updateObj := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "room_availability", Value: avail},
			{Key: "updated_at", Value: updated_at},
		}},
	}

	upsert := true
	filter := bson.M{"room_id": id}
	options := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := database.RoomCollection.UpdateOne(context.TODO(), filter, updateObj, &options)
	if err != nil {
		return err
	}

	return nil
}

func FilterRoom() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := helpers.GetContext()
		defer cancel()

		room_type := c.PostForm("room_type")
		room_availability := c.PostForm("room_availability")
		cleaning_status := c.PostForm("cleaning_status")
		priceOperator := c.PostForm("priceOperator")
		price, _ := strconv.ParseFloat(c.PostForm("price"), 64)

		if price <= 0.0 {
			price = 0
		}

		if priceOperator == "" || (priceOperator != "$gt" && priceOperator != "$eq" && priceOperator != "$lt") {
			priceOperator = "$gt"
		}

		if room_type == "" || (room_type != "single_bed" && room_type != "double_bed" && room_type != "suite") {
			room_type = "single_bed"
		}

		if room_availability == "" || (room_availability != "available" && room_availability != "occupied") {
			room_availability = "available"
		}

		if cleaning_status == "" || (cleaning_status != "cleaned" && cleaning_status != "dirty" && cleaning_status != "inprogress") {
			cleaning_status = "cleaned"
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
			{Key: "room_type", Value: room_type},
			{Key: "room_availability", Value: room_availability},
			{Key: "cleaning_status", Value: cleaning_status},
			{Key: "price", Value: bson.D{
				{Key: priceOperator, Value: price},
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
			{Key: "rooms",
				Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
		}}}

		projectStage2 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "total_count", Value: 1},
			{Key: "room", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$rooms"},
					{Key: "as", Value: "data"},
					{Key: "in", Value: bson.D{
						{Key: "room_id", Value: "$$data.room_id"},
						{Key: "branch_id", Value: "$$data.branch_id"},
						{Key: "room_number", Value: "$$data.room_number"},
						{Key: "room_type", Value: "$$data.room_type"},
						{Key: "room_availability", Value: "$$data.room_availability"},
						{Key: "cleaning_status", Value: "$$data.cleaning_status"},
						{Key: "price", Value: "$$data.price"},
						{Key: "capacity", Value: "$$data.capacity"},
						{Key: "images", Value: "$$data.images"},
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

		result, err := database.RoomCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage1, projectStage2,
		})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching rooms "+err.Error())
			return
		}

		var allRooms []bson.M
		if err := result.All(ctx, &allRooms); err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting the room "+err.Error())
			return
		}

		if len(allRooms) == 0 {
			utils.Response(c, []interface{}{})
			return
		}

		utils.Response(c, allRooms[0])
	}
}
