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
func GetAllPickUpServices() gin.HandlerFunc {
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
			{Key: "services",
				Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
		}}}

		projectStage2 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "total_count", Value: 1},
			{Key: "service", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$services"},
					{Key: "as", Value: "data"},
					{Key: "in", Value: bson.D{
						{Key: "pickup_service_id", Value: "$$data.pickup_service_id"},
						{Key: "branch_id", Value: "$$data.branch_id"},
						{Key: "guest_id", Value: "$$data.guest_id"},
						{Key: "driver", Value: "$$data.driver_id"},
						{Key: "pickup_location", Value: "$$data.pickup_location"},
						{Key: "pickup_time", Value: "$$data.pickup_time"},
						{Key: "amount", Value: "$$data.amount"},
						{Key: "status", Value: "$$data.status"},
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

		result, err := database.PickupServiceCollection.Aggregate(ctx, mongo.Pipeline{
			groupStage, projectStage1, projectStage2,
		})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching services "+err.Error())
			return
		}

		var allservices []bson.M
		if err := result.All(ctx, &allservices); err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting the service "+err.Error())
			return
		}

		if len(allservices) == 0 {
			utils.Response(c, []interface{}{})
			return
		}

		utils.Response(c, allservices[0])
	}
}

// * DONE
func GetPickUpService() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var service models.PickupService
		id := c.Param("id")

		if err := database.PickupServiceCollection.FindOne(ctx, bson.M{"pickup_service_id": id}).Decode(&service); err != nil {
			utils.Error(c, utils.NotFound, "Can't find the service with id.")
			return
		}

		utils.Response(c, service)
	}
}

// * DONE
func CreatePickUpService() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var service models.PickupService

		service.Guest_id = c.PostForm("guest_id")
		service.Branch_id = c.PostForm("branch_id")
		service.Driver_id = c.PostForm("driver_id")
		service.Pickup_Location = c.PostForm("pickup_location")
		service.Pickup_Time = c.PostForm("pickup_time")

		// Validate details
		msg, isVal := validatePickUpServiceDetails(service)
		if !isVal {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// Check branch exist or not
		count, bErr := database.PickupServiceCollection.CountDocuments(ctx, bson.M{"branch_id": service.Branch_id, "guest_id": service.Guest_id, "driver_id": service.Driver_id})
		if bErr != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching pickup service data.")
			return
		}

		if count > 0 {
			utils.Error(c, utils.Conflict, "Service already exists with branch_id, guest_id and driver_id")
			return
		}

		// Generate ID and Timestamp
		service.ID = primitive.NewObjectID()
		service.Pickup_Service_id = service.ID.Hex()
		service.Status = models.NotCompleted
		service.Amount = 500
		service.Created_at, _ = helpers.GetTime()
		service.Updated_at, _ = helpers.GetTime()

		// Insert Service Details
		result, err := database.PickupServiceCollection.InsertOne(ctx, service)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't create service.")
			return
		}

		// If of reture response
		utils.Response(c, result)
	}
}

// * DONE
func UpdatePickUpServiceDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var service models.PickupService
		id := c.Param("id")

		service.Pickup_Location = c.PostForm("pickup_location")
		service.Pickup_Time = c.PostForm("pickup_time")

		// Validate data
		msg, isVal := validatePickUpServiceUpdateDetails(service)
		if !isVal {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// check service exist or not
		count, err := database.PickupServiceCollection.CountDocuments(ctx, bson.M{"pickup_service_id": id})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while checking service.")
			return
		}

		if !(count > 0) {
			utils.Error(c, utils.NotFound, "Service does not exist with id.")
			return
		}

		// Update data here

		var updateObj primitive.D

		updateObj = append(updateObj, bson.E{Key: "pickup_location", Value: service.Pickup_Location})

		updateObj = append(updateObj, bson.E{Key: "pickup_time", Value: service.Pickup_Time})

		service.Updated_at, _ = helpers.GetTime()
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: service.Updated_at})

		filter := bson.M{"pickup_service_id": id}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		// If success send success message
		_, err = database.PickupServiceCollection.UpdateOne(ctx, filter, bson.D{
			{Key: "$set", Value: updateObj},
		}, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update service details.")
			return
		}

		utils.Message(c, "Service details updated successfully.")
	}
}

// * DONE
func DeletePickUpService() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()

		id := c.Param("id")

		_, err := database.PickupServiceCollection.DeleteOne(ctx, bson.M{"pickup_service_id": id})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete service.")
			return
		}
		utils.Message(c, "Service deleted successfully.")
	}
}

// * DONE
func UpdatePickUpServiceStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var service models.PickupService
		id := c.Param("id")

		if err := database.PickupServiceCollection.FindOne(ctx, bson.M{"pickup_service_id": id}).Decode(&service); err != nil {
			utils.Error(c, utils.NotFound, "Can't find service with ID.")
			return
		}

		updated_at, _ := helpers.GetTime()
		updateObj := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "status", Value: models.Completed},
				{Key: "updated_at", Value: updated_at},
			}},
		}

		upsert := true
		filter := bson.M{"pickup_service_id": id}
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		_, err := database.PickupServiceCollection.UpdateOne(ctx, filter, updateObj, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error updating service status")
			return
		}

		utils.Message(c, "Service status updated successfully.")
	}
}

func validatePickUpServiceDetails(service models.PickupService) (string, bool) {

	if service.Guest_id == "" {
		return "Guest id required", false
	}

	if service.Branch_id == "" {
		return "Branch id required", false
	}

	if service.Driver_id == "" {
		return "Driver id required", false
	}

	if service.Pickup_Location == "" {
		return "Pickup location required", false
	}

	if service.Pickup_Time == "" {
		return "Pickup time required", false
	}

	return "", true
}

func validatePickUpServiceUpdateDetails(service models.PickupService) (string, bool) {

	if service.Pickup_Location == "" {
		return "Pickup location required", false
	}

	if service.Pickup_Time == "" {
		return "Pickup time required", false
	}

	return "", true
}
