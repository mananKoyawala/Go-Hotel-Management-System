package controllers

import (
	"context"
	"strconv"
	"time"

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
func GetAllReservations() gin.HandlerFunc {
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
			{Key: "reservation",
				Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
		}}}

		projectStage2 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "total_count", Value: 1},
			{Key: "reservation", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$reservation"},
					{Key: "as", Value: "data"},
					{Key: "in", Value: bson.D{
						{Key: "reservation_id", Value: "$$data.reservation_id"},
						{Key: "room_id", Value: "$$data.room_id"},
						{Key: "guest_id", Value: "$$data.guest_id"},
						{Key: "check_in_time", Value: "$$data.check_in_time"},
						{Key: "check_out_time", Value: "$$data.check_out_time"},
						{Key: "deposit_amount", Value: "$$data.deposit_amount"},
						{Key: "pending_amount", Value: bson.D{
							{Key: "$cond", Value: bson.A{
								// Condition: Check if pending_amount exists
								bson.D{{Key: "$gt", Value: bson.A{"$$data.pending_amount", nil}}},
								// If true: Keep the existing value
								"$$data.pending_amount",
								// If false: Set to null or some default value
								0, // or your default value
							}},
						}},
						{Key: "numbers_of_guests", Value: "$$data.numbers_of_guests"},
						{Key: "is_check_out", Value: "$$data.is_check_out"},
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

		result, err := database.ReservationCollection.Aggregate(ctx, mongo.Pipeline{
			groupStage, projectStage1, projectStage2,
		})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching reservation "+err.Error())
			return
		}

		var allreservation []bson.M
		if err := result.All(ctx, &allreservation); err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting the reservation "+err.Error())
			return
		}

		if len(allreservation) == 0 {
			utils.Response(c, []interface{}{})
			return
		}

		utils.Response(c, allreservation[0])
	}
}

// * DONE
func GetReservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var reservation models.Reservation
		id := c.Param("id")

		if err := database.ReservationCollection.FindOne(ctx, bson.M{"reservation_id": id}).Decode(&reservation); err != nil {
			utils.Error(c, utils.NotFound, "Can't find the reservation with id.")
			return
		}

		utils.Response(c, reservation)
	}
}

// * DONE
func CreateReservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// here context time is 100 sec beacuse reservation may take some time
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var reservation models.Reservation
		var room models.Room

		reservation.Room_id = c.PostForm("room_id")
		reservation.Guest_id = c.PostForm("guest_id")
		reservation.Check_in_time = c.PostForm("check_in_time")
		reservation.Check_out_time = c.PostForm("check_out_time")
		reservation.Deposit_Amount, _ = strconv.ParseFloat(c.PostForm("desposit_amount"), 64)
		reservation.Numbers_of_guests, _ = strconv.Atoi(c.PostForm("numbers_of_guests"))

		// validate details
		msg, val := valiadateReservationDetails(reservation)
		if !val {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// get room details
		if err := database.RoomCollection.FindOne(ctx, bson.M{"room_id": reservation.Room_id}).Decode(&room); err != nil {
			utils.Error(c, utils.NotFound, "Can't find room with id")
			return
		}

		// check room occupied by guest
		if room.Room_Availability == string(models.Room_Unavailable) {
			utils.Error(c, utils.Conflict, "Room already occupied by guest.")
			return
		}

		// check num of guest more than capacity
		roomCapacity := room.Capacity
		numOfGuest := reservation.Numbers_of_guests
		if numOfGuest > roomCapacity {
			utils.Error(c, utils.BadRequest, "Room capacity exceeded.")
			return
		}

		// pending amount
		if room.Price < reservation.Deposit_Amount {
			utils.Error(c, utils.BadRequest, "Deposit ammout exceeded.")
			return
		}
		pendingAmount := room.Price - reservation.Deposit_Amount
		reservation.Pending_amount = pendingAmount

		// generate id, timestamps
		reservation.ID = primitive.NewObjectID()
		reservation.Reservation_id = reservation.ID.Hex()
		reservation.IsCheckOut = string(models.False)
		reservation.Created_at, _ = helpers.GetTime()
		reservation.Updated_at, _ = helpers.GetTime()

		// create reservation
		result, err := database.ReservationCollection.InsertOne(ctx, reservation)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't reserve room.")
			return
		}

		// update room is occupied
		if err := UpdateRoomAvailability(room.Room_id, models.Room_Unavailable); err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update room availability.")
			return
		}

		// if success return
		utils.Response(c, result)
	}
}

// * DONE
func UpdateReservationDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var reservation models.Reservation
		var foundReservation models.Reservation
		id := c.Param("id")

		reservation.Check_in_time = c.PostForm("check_in_time")
		reservation.Check_out_time = c.PostForm("check_out_time")
		reservation.Deposit_Amount, _ = strconv.ParseFloat(c.PostForm("desposit_amount"), 64)
		reservation.Numbers_of_guests, _ = strconv.Atoi(c.PostForm("numbers_of_guests"))
		reservation.IsCheckOut = c.PostForm("is_check_out")

		// Validate data
		msg, isval := valiadateUpdateReservationDetails(reservation)
		if !isval {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// check already checkout or not
		if err := database.ReservationCollection.FindOne(ctx, bson.M{"reservation_id": id}).Decode(&foundReservation); err != nil {
			utils.Error(c, utils.NotFound, "Can't find reservation with id.")
			return
		}

		if foundReservation.IsCheckOut == string(models.True) {
			utils.Error(c, utils.Conflict, "User has been already checked out.")
			return
		}

		// Update data here
		var updateObj primitive.D

		updateObj = append(updateObj, bson.E{Key: "check_in_time", Value: reservation.Check_in_time})

		updateObj = append(updateObj, bson.E{Key: "check_out_time", Value: reservation.Check_out_time})
		updateObj = append(updateObj, bson.E{Key: "deposit_amount", Value: reservation.Deposit_Amount})
		updateObj = append(updateObj, bson.E{Key: "pending_amount", Value: reservation.Pending_amount})
		updateObj = append(updateObj, bson.E{Key: "numbers_of_guests", Value: reservation.Numbers_of_guests})
		updateObj = append(updateObj, bson.E{Key: "is_check_out", Value: reservation.IsCheckOut})

		reservation.Updated_at, _ = helpers.GetTime()
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: reservation.Updated_at})

		filter := bson.M{"reservation_id": id}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		// If success send success message

		_, err := database.ReservationCollection.UpdateOne(ctx, filter, bson.D{
			{Key: "$set", Value: updateObj},
		}, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update reservation details")
			return
		}
		utils.Message(c, "Reservation is updated.")
	}
}

// * DONE
func DeleteReservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()

		id := c.Param("id")
		room_id := c.Param("room_id")

		_, err := database.ReservationCollection.DeleteOne(ctx, bson.M{"reservation_id": id})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete reservation.")
			return
		}

		if err := UpdateRoomAvailability(room_id, models.Room_Available); err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update room availability.")
			return
		}

		utils.Message(c, "Reservation deleted successfully.")
	}
}

func valiadateReservationDetails(reservation models.Reservation) (string, bool) {

	if reservation.Room_id == "" {
		return "Room id is required", false
	}

	if reservation.Guest_id == "" {
		return "Guest id is required", false
	}

	if reservation.Check_in_time == "" {
		return "Check in time is required", false
	}

	if reservation.Check_out_time == "" {
		return "Check out time is required", false
	}

	if !utils.IsNonNegative(int(reservation.Deposit_Amount)) {
		return "Deposit amount is required", false
	}

	if reservation.Numbers_of_guests <= 0 {
		return "Number of guest is not less or equal 0", false
	}

	return "", true
}

func valiadateUpdateReservationDetails(reservation models.Reservation) (string, bool) {

	if reservation.Check_in_time == "" {
		return "Check in time is required", false
	}

	if reservation.Check_out_time == "" {
		return "Check out time is required", false
	}

	if !utils.IsNonNegative(int(reservation.Deposit_Amount)) {
		return "Deposit amount is required", false
	}

	if reservation.Numbers_of_guests <= 0 {
		return "Number of guest is not less or equal 0", false
	}
	// log.Println(reservation.IsCheckOut)
	if !validateIsCheckOut(reservation.IsCheckOut) {
		return "Is check out must be true or false", false
	}

	return "", true
}

func validateIsCheckOut(checkout string) bool {
	return models.BOOL(checkout) == models.True || models.BOOL(checkout) == models.False
}
