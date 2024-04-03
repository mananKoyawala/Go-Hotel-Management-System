package controllers

import (
	"context"
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
func GetRoom() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := helpers.GetContext()
		defer cancel()
		var room models.Room
		id := c.Param("id")

		if err := database.RoomCollection.FindOne(ctx, bson.M{"room_id": id}).Decode(&room); err != nil {
			utils.Error(c, utils.InternalServerError, "Can't find the room with id.")
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

		// Check json
		if err := c.BindJSON(&room); err != nil {
			utils.Error(c, utils.BadRequest, "Invalid JSON Format")
			return
		}

		// Validate details
		msg, isVal := validateRoomDetails(room)
		if !isVal {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// Check branch exist or not
		count, bErr := database.BranchCollection.CountDocuments(ctx, bson.M{"branch_id": room.Branch_id})
		if bErr != nil || count <= 0 {
			utils.Error(c, utils.InternalServerError, "Branch doesn't exist")
			return
		}

		// Check the room number exists
		count, rErr := database.RoomCollection.CountDocuments(ctx, bson.M{"room_number": room.Room_Number, "branch_id": room.Branch_id})
		if rErr != nil {
			utils.Error(c, utils.InternalServerError, "Error while checking Room Number")
			return
		}

		if count > 0 {
			utils.Error(c, utils.BadRequest, "Room Number already Exist with Branch id.")
			return
		}

		// Generate ID and Timestamp
		room.ID = primitive.NewObjectID()
		room.Room_id = room.ID.Hex()
		room.Room_Availability = models.Room_Available
		room.Cleaning_Status = models.Cleaned
		room.Created_at, _ = helpers.GetTime()
		room.Updated_at, _ = helpers.GetTime()

		// Insert Room Details
		result, err := database.RoomCollection.InsertOne(ctx, room)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't insert room.")
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

		// Check json
		if err := c.BindJSON(&room); err != nil {
			utils.Error(c, utils.BadRequest, "Invalid JSON Format")
			return
		}

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

		id := c.Param("id")

		_, err := database.RoomCollection.DeleteOne(ctx, bson.M{"room_id": id})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete room.")
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

	if !utils.IsNonNegative(room.Capacity) {
		return "Room capacity must be integer.", false
	}

	return "", true
}

func checkRoomType(roomType models.Room_Type) bool {

	if roomType == "" {
		return false
	}

	if roomType == models.Single_Bad {
		return true
	}

	if roomType == models.Double_Bad {
		return true
	}

	if roomType == models.Suite {
		return true
	}

	return false
}

func checkRoomStatus(status models.Cleaning_Status) bool {

	if status == "" {
		return false
	}

	if status == models.Cleaned {
		return true
	}

	if status == models.Dirty {
		return true
	}

	if status == models.InProgress {
		return true
	}

	return true
}

func checkRoomAvailability(avail models.Room_Availability) bool {

	if avail == "" {
		return false
	}

	if avail == models.Room_Available {
		return true
	}

	if avail == models.Room_Unavailable {
		return true
	}

	return true
}

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
