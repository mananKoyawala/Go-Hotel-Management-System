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
func GetAllStaff() gin.HandlerFunc {
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
			{Key: "staffs",
				Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
		}}}

		projectStage2 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "total_count", Value: 1},
			{Key: "staff", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$staffs"},
					{Key: "as", Value: "data"},
					{Key: "in", Value: bson.D{
						{Key: "staff_id", Value: "$$data.staff_id"},
						{Key: "branch_id", Value: "$$data.branch_id"},
						{Key: "first_name", Value: "$$data.first_name"},
						{Key: "last_name", Value: "$$data.last_name"},
						{Key: "phone", Value: "$$data.phone"},
						{Key: "email", Value: "$$data.email"},
						{Key: "gender", Value: "$$data.gender"},
						{Key: "aadhar_number", Value: "$$data.aadhar_number"},
						{Key: "age", Value: "$$data.age"},
						{Key: "job_type", Value: "$$data.job_type"},
						{Key: "salary", Value: "$$data.salary"},
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

		result, err := database.StaffCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage1, projectStage2,
		})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching staffs "+err.Error())
			return
		}

		var allstaffs []bson.M
		if err := result.All(ctx, &allstaffs); err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting the staffs "+err.Error())
			return
		}

		if len(allstaffs) == 0 {
			utils.Response(c, []interface{}{})
			return
		}

		utils.Response(c, allstaffs[0])
	}
}

// * DONE
func GetStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var staff models.Staff

		id := c.Param("id")

		if err := database.StaffCollection.FindOne(ctx, bson.M{"staff_id": id}).Decode(&staff); err != nil {
			utils.Error(c, utils.NotFound, "Can't find staff with id")
			return
		}

		utils.Response(c, staff)
	}
}

// * DONE
func CreateStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var staff models.Staff

		// check json
		if err := c.BindJSON(&staff); err != nil {
			utils.Error(c, utils.BadRequest, "Invalid JSON Format")
			return
		}

		// validate staff details
		msg, isVal := validateStaffDetails(staff)
		if !isVal {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		count, err := database.StaffCollection.CountDocuments(ctx, bson.M{"email": staff.Email})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting details.")
			return
		}

		if count > 0 {
			utils.Error(c, utils.Conflict, "Email already in use fo staff, try different email addresses.")
			return
		}

		// generate id, timestamps
		staff.ID = primitive.NewObjectID()
		staff.Staff_id = staff.ID.Hex()
		staff.Status = models.Active
		staff.Created_at, _ = helpers.GetTime()
		staff.Updated_at, _ = helpers.GetTime()

		// Insert the details
		result, err := database.StaffCollection.InsertOne(ctx, staff)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't add staff")
			return
		}

		// if success return
		utils.Response(c, result)
	}
}

// * DONE
func UpdateStaffDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var staff models.Staff
		id := c.Param("id")

		// Check json
		if err := c.BindJSON(&staff); err != nil {
			utils.Error(c, utils.BadRequest, "Invalid JSON Format")
			return
		}

		// Validate data
		msg, isVal := validateStaffDetails(staff)
		if !isVal {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// check staff exist or not
		count, err := database.StaffCollection.CountDocuments(ctx, bson.M{"staff_id": id})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while checking staff.")
			return
		}

		if !(count > 0) {
			utils.Error(c, utils.BadRequest, "Staff does not exist with id.")
			return
		}

		// Update data here

		var updateObj primitive.D

		updateObj = append(updateObj, bson.E{Key: "branch_id", Value: staff.Branch_id})

		updateObj = append(updateObj, bson.E{Key: "first_name", Value: staff.First_Name})

		updateObj = append(updateObj, bson.E{Key: "last_name", Value: staff.Last_Name})

		updateObj = append(updateObj, bson.E{Key: "gender", Value: staff.Gender})

		updateObj = append(updateObj, bson.E{Key: "email", Value: staff.Email})

		updateObj = append(updateObj, bson.E{Key: "phone", Value: staff.Phone})

		updateObj = append(updateObj, bson.E{Key: "job_type", Value: staff.Job_Type})

		updateObj = append(updateObj, bson.E{Key: "age", Value: staff.Age})

		updateObj = append(updateObj, bson.E{Key: "salary", Value: staff.Salary})

		updateObj = append(updateObj, bson.E{Key: "aadhar_number", Value: staff.Aadhar_Number})

		staff.Updated_at, _ = helpers.GetTime()
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: staff.Updated_at})

		filter := bson.M{"staff_id": id}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		// If success send success message
		_, err = database.StaffCollection.UpdateOne(ctx, filter, bson.D{
			{Key: "$set", Value: updateObj},
		}, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update staff details.")
			return
		}

		utils.Message(c, "Staff details updated successfully.")
	}
}

// * DONE
func UpdateStaffStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var staff models.Staff
		id := c.Param("id")

		if err := database.StaffCollection.FindOne(ctx, bson.M{"staff_id": id}).Decode(&staff); err != nil {
			utils.Error(c, utils.NotFound, "Can't find staff with ID.")
			return
		}

		var newStatus models.Status

		if staff.Status == models.Status(models.Active) {
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
		filter := bson.M{"staff_id": id}
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		_, err := database.StaffCollection.UpdateOne(ctx, filter, updateObj, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error updating staff status")
			return
		}

		utils.Message(c, "Staff status updated successfully.")
	}
}

// * DONE
func DeleteStaff() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()

		id := c.Param("id")

		_, err := database.StaffCollection.DeleteOne(ctx, bson.M{"staff_id": id})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete staff.")
			return
		}
		utils.Message(c, "Staff deleted successfully.")
	}
}

func validateStaffDetails(staff models.Staff) (string, bool) {
	if staff.Branch_id == "" {
		return "Branch id required", false
	}

	if staff.First_Name == "" {
		return "First name is required", false
	}

	if staff.Last_Name == "" {
		return "Last name is required", false
	}

	if utils.CheckLength(staff.Phone, 10) {
		return "Phone number must be 10 digits", false
	}

	if staff.Email == "" {
		return "Email address required", false
	} else if !utils.ValidateEmail(staff.Email) {
		return "Invalid email address", false
	}

	if staff.Gender == "" {
		return "Gender is required", false
	}

	if staff.Job_Type == "" {
		return "Job Type is required", false
	}

	if staff.Age < 16 || staff.Age > 65 {
		return "Age must be between 16 to 65", false
	}

	if !utils.IsNonNegative(int(staff.Salary)) {
		return "Salary must not 0 or negative", false
	}

	if utils.CheckLength(staff.Aadhar_Number, 12) {
		return "Aadhar Number must be of 12 digits", false
	}

	return "", true
}