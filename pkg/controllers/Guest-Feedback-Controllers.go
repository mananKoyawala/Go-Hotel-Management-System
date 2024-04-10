package controllers

import (
	"fmt"
	"path/filepath"
	"regexp"
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

var feedbackFolder = "feedback"

// * DONE
func GetAllFeedbacks() gin.HandlerFunc {
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

		matchStage := bson.D{{Key: "$match", Value: bson.D{
			{Key: "branch_id", Value: branch_id},
		}}}

		groupStage := bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "null"},
			{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}}}

		projectStage1 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "feedbacks",
				Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
		}}}

		projectStage2 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "total_count", Value: 1},
			{Key: "feedback", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$feedbacks"},
					{Key: "as", Value: "data"},
					{Key: "in", Value: bson.D{
						{Key: "guest_feedback_id", Value: "$$data.guest_feedback_id"},
						{Key: "branch_id", Value: "$$data.branch_id"},
						{Key: "room_id", Value: "$$data.room_id"},
						{Key: "guest_id", Value: "$$data.guest_id"},
						{Key: "description", Value: "$$data.description"},
						{Key: "feedback_type", Value: "$$data.feedback_type"},
						{Key: "resolution_details", Value: "$$data.resolution_details"},
						{Key: "rating", Value: "$$data.rating"},
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

		result, err := database.GuestFeedbackCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage1, projectStage2,
		})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching feedbacks "+err.Error())
			return
		}

		var allFeedbacks []bson.M
		if err := result.All(ctx, &allFeedbacks); err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting the feedbacks "+err.Error())
			return
		}

		if len(allFeedbacks) == 0 {
			utils.Response(c, []interface{}{})
			return
		}

		utils.Response(c, allFeedbacks[0])
	}
}

// * DONE
func GetFeedback() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var feedback models.GuestFeedback

		id := c.Param("id")

		if err := database.GuestFeedbackCollection.FindOne(ctx, bson.M{"guest_feedback_id": id}).Decode(&feedback); err != nil {
			utils.Error(c, utils.NotFound, "Can't find feedback with id")
			return
		}

		utils.Response(c, feedback)
	}
}

// * DONE
func CreateFeedback() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var feedback models.GuestFeedback

		feedback.Branch_id = c.PostForm("branch_id")
		feedback.Room_id = c.PostForm("room_id")
		feedback.Guest_id = c.PostForm("guest_id")
		feedback.Description = c.PostForm("description")
		feedback.Feedback_Type = c.PostForm("feedback_type")
		feedback.Rating = c.PostForm("rating")

		// validate guest details
		msg, isVal := validateFeedbackDetails(feedback)
		if !isVal {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		count, err := database.GuestFeedbackCollection.CountDocuments(ctx, bson.M{"branch_id": feedback.Branch_id, "room_id": feedback.Room_id, "guest_id": feedback.Guest_id})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting details.")
			return
		}

		if count > 0 {
			utils.Error(c, utils.Conflict, "Feedback already exist with branch, room and guest id.")
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

		// generate id, timestamps
		feedback.ID = primitive.NewObjectID()
		feedback.Guest_Feedback_id = feedback.ID.Hex()
		feedback.Status = models.Pending
		feedback.Created_at, _ = helpers.GetTime()
		feedback.Updated_at, _ = helpers.GetTime()

		// upload image
		name := strings.ReplaceAll(handler.Filename, " ", "")
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), name)
		url, err := imageupload.UploadService(file, feedbackFolder, filename)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't uplaod the image.")
			return
		}
		feedback.Image = url

		// Insert the details
		result, err := database.GuestFeedbackCollection.InsertOne(ctx, feedback)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't create feedback")
			return
		}

		// if success return
		utils.Response(c, result)
	}
}

// * DONE
func UpdateFeedbackResolution() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var feedback models.GuestFeedback
		id := c.Param("id")

		feedback.Resolution_Details = c.PostForm("resolution_details")

		if feedback.Resolution_Details == "" || len(feedback.Description) > 150 {
			utils.Error(c, utils.BadRequest, "Only 150 characters accepted.")
			return
		}

		// check feedback exist with guest_feedback_id
		count, err := database.GuestFeedbackCollection.CountDocuments(ctx, bson.M{"guest_feedback_id": id})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting guest feedback details.")
			return
		}

		if !(count > 0) {
			utils.Error(c, utils.NotFound, "Can't find feedback with id.")
			return
		}

		// Update data here

		var updateObj primitive.D

		updateObj = append(updateObj, bson.E{Key: "resolution_details", Value: feedback.Resolution_Details})

		feedback.Status = models.Resolved
		updateObj = append(updateObj, bson.E{Key: "status", Value: feedback.Status})

		feedback.Updated_at, _ = helpers.GetTime()
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: feedback.Updated_at})

		filter := bson.M{"guest_feedback_id": id}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		// If success send success message
		_, err = database.GuestFeedbackCollection.UpdateOne(ctx, filter, bson.D{
			{Key: "$set", Value: updateObj},
		}, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update feedback details.")
			return
		}

		utils.Message(c, "Feedback details updated successfully.")
	}
}

// * DONE
func DeleteFeedback() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var feedback models.GuestFeedback
		id := c.Param("id")

		if err := database.GuestFeedbackCollection.FindOne(ctx, bson.M{"guest_feedback_id": id}).Decode(&feedback); err != nil {
			utils.Error(c, utils.NotFound, "Can't find feedback with id")
			return
		}

		image := utils.GetTrimedUrl(feedback.Image)
		if err := imageupload.DeleteService(image); err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete image."+err.Error())
			return
		}

		_, err := database.GuestFeedbackCollection.DeleteOne(ctx, bson.M{"guest_feedback_id": id})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete Feedback.")
			return
		}
		utils.Message(c, "Feddback deleted successfully.")
	}
}

func validateFeedbackDetails(feedback models.GuestFeedback) (string, bool) {

	if feedback.Branch_id == "" {
		return "Branch id required", false
	}

	if feedback.Room_id == "" {
		return "Room id required", false
	}

	if feedback.Guest_id == "" {
		return "Guest id required", false
	}

	if feedback.Description == "" {
		return "Description required", false
	}

	if !validateFeedbackType(feedback.Feedback_Type) {
		return "Feedback type must be complaint or rating", false
	}

	if !validateRating(feedback.Rating) {
		return "Rating must be between 0 and 5", false
	}

	return "", true
}

func validateFeedbackType(fType string) bool {

	if models.Feedback_Type(fType) == models.Complaint {
		return true
	}

	if models.Feedback_Type(fType) == models.Rating {
		return true
	}

	return false
}

func validateRating(ratingStr string) bool {
	validRatingRegex := regexp.MustCompile(`^[0-5](\.\d)?$`)
	return validRatingRegex.MatchString(ratingStr)
}

// * DONE
func FilterFeedback() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()

		status := models.Status(c.PostForm("status"))
		feedbackType := c.PostForm("feedback_type")
		rating := c.PostForm("rating")

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

		// cleaning data

		if feedbackType == "" || (feedbackType != "complaint" && feedbackType != "rating") {
			feedbackType = "rating"
		}

		if status == "" || (status != "resolved" && status != "pending") {
			status = "resolved"
		}

		matchStage := bson.D{{Key: "$match", Value: bson.D{

			{Key: "feedback_type", Value: feedbackType},
			{Key: "status", Value: status},
			{Key: "$or", Value: bson.A{
				bson.D{{Key: "rating", Value: bson.D{{Key: "$regex", Value: rating}, {Key: "$options", Value: "i"}}}},
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
			{Key: "feedbacks",
				Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
		}}}

		projectStage2 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "total_count", Value: 1},
			{Key: "feedback", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$feedbacks"},
					{Key: "as", Value: "data"},
					{Key: "in", Value: bson.D{
						{Key: "guest_feedback_id", Value: "$$data.guest_feedback_id"},
						{Key: "branch_id", Value: "$$data.branch_id"},
						{Key: "room_id", Value: "$$data.room_id"},
						{Key: "guest_id", Value: "$$data.guest_id"},
						{Key: "description", Value: "$$data.description"},
						{Key: "feedback_type", Value: "$$data.feedback_type"},
						{Key: "resolution_details", Value: "$$data.resolution_details"},
						{Key: "rating", Value: "$$data.rating"},
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

		result, err := database.GuestFeedbackCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage1, projectStage2,
		})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching feedbacks "+err.Error())
			return
		}

		var allFeedbacks []bson.M
		if err := result.All(ctx, &allFeedbacks); err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting the feedbacks "+err.Error())
			return
		}

		if len(allFeedbacks) == 0 {
			utils.Response(c, []interface{}{})
			return
		}

		utils.Response(c, allFeedbacks[0])
	}
}
