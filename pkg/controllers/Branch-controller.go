package controllers

import (
	"fmt"
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
func GetBranch() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var branch models.Branch

		id := c.Param("id")

		if err := database.BranchCollection.FindOne(ctx, bson.M{"branch_id": id}).Decode(&branch); err != nil {
			utils.Error(c, utils.NotFound, "Can't find branch")
			return
		}
		utils.Response(c, branch)
	}
}

// * DONE
func GetBranches() gin.HandlerFunc {
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
			{Key: "branches",
				Value: bson.D{{Key: "$slice", Value: bson.A{"$data", startIndex, recordPerPage}}}},
		}}}

		projectStage2 := bson.D{{Key: "$project", Value: bson.D{
			{Key: "total_count", Value: 1},
			{Key: "branch", Value: bson.D{
				{Key: "$map", Value: bson.D{
					{Key: "input", Value: "$branches"},
					{Key: "as", Value: "data"},
					{Key: "in", Value: bson.D{
						{Key: "branch_id", Value: "$$data.branch_id"},
						{Key: "manager_id", Value: "$$data.manager_id"},
						{Key: "branch_name", Value: "$$data.branch_name"},
						{Key: "address", Value: "$$data.address"},
						{Key: "phone", Value: "$$data.phone"},
						{Key: "email", Value: "$$data.email"},
						{Key: "city", Value: "$$data.city"},
						{Key: "state", Value: "$$data.state"},
						{Key: "country", Value: "$$data.country"},
						{Key: "pincode", Value: "$$data.pincode"},
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

		result, err := database.BranchCollection.Aggregate(ctx, mongo.Pipeline{
			groupStage, projectStage1, projectStage2,
		})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching branches  "+err.Error())
			return
		}

		var allBranches []bson.M
		if err := result.All(ctx, &allBranches); err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting the branches "+err.Error())
			return
		}
		if len(allBranches) == 0 {
			utils.Response(c, []interface{}{})
			return
		}
		utils.Response(c, allBranches[0])
	}
}

// * DONE
func CreateBranch() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var branch models.Branch

		// Check JSON
		if err := c.BindJSON(&branch); err != nil {
			utils.Error(c, utils.BadRequest, "Invalid JSON Format")
			return
		}

		// Validate Details
		msg, isVal := validateBranchDetails(&branch)
		if !isVal {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		// Check Branch already exist or not by Branch_Name,City,Pincode
		count, err := database.BranchCollection.CountDocuments(ctx, bson.M{"branch_name": branch.Branch_Name, "city": branch.City, "pincode": branch.Pincode})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting branch data")
			return
		}

		if count > 0 {
			msg := fmt.Sprintf("Branch already exist with Branch Name %s, City %s and Pincode %d", branch.Branch_Name, branch.City, branch.Pincode)
			utils.Error(c, utils.Conflict, msg)
			return
		}

		// Generate IDs and Timestamps
		branch.ID = primitive.NewObjectID()
		branch.Branch_id = branch.ID.Hex()
		branch.Status = models.Active
		branch.Total_Rooms = 0 // When branch is created there are no rooms that associated with it
		// If admin add room that it will counted here
		branch.Created_at, _ = helpers.GetTime()
		branch.Updated_at, _ = helpers.GetTime()

		// Insert Details
		result, err := database.BranchCollection.InsertOne(ctx, branch)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't add branch to database.")
			return
		}

		// Return is success
		utils.Response(c, result)
	}
}

// * DONE
func DeleteBranch() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()

		id := c.Param("id")

		_, err := database.BranchCollection.DeleteOne(ctx, bson.M{"branch_id": id})
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete branch")
			return
		}
		utils.Message(c, "Branch deleted successfully")
	}
}

// * DONE
func UpdateBranchDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var branch models.Branch
		id := c.Param("id")

		// Check json
		if err := c.BindJSON(&branch); err != nil {
			utils.Error(c, utils.BadRequest, "Invalid JSON Format")
			return
		}

		// Validate data

		var updateObj primitive.D

		// Update details if not empty
		if branch.Branch_Name != "" {
			updateObj = append(updateObj, bson.E{Key: "branch_name", Value: branch.Branch_Name})
		}

		if branch.Address != "" {
			updateObj = append(updateObj, bson.E{Key: "address", Value: branch.Address})
		}

		if branch.Phone != 0 {
			updateObj = append(updateObj, bson.E{Key: "phone", Value: branch.Phone})
		}

		if branch.Email != "" {
			updateObj = append(updateObj, bson.E{Key: "email", Value: branch.Email})
		}

		if branch.City != "" {
			updateObj = append(updateObj, bson.E{Key: "city", Value: branch.City})
		}

		if branch.State != "" {
			updateObj = append(updateObj, bson.E{Key: "state", Value: branch.State})
		}

		if branch.Country != "" {
			updateObj = append(updateObj, bson.E{Key: "country", Value: branch.Country})
		}

		if branch.Pincode != 0 {
			updateObj = append(updateObj, bson.E{Key: "pincode", Value: branch.Pincode})
		}

		// Update timestamp
		branch.Updated_at, _ = helpers.GetTime()
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: branch.Updated_at})

		// Update in database
		filter := bson.M{"branch_id": id}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}

		_, err := database.BranchCollection.UpdateOne(ctx, filter, bson.D{
			{Key: "$set", Value: updateObj},
		}, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update branch data.")
			return
		}

		// If success then return message
		utils.Message(c, "Branch details are updated.")
	}
}

// * DONE
func UpdateBranchStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var branch models.Branch
		id := c.Param("id")

		if err := database.BranchCollection.FindOne(ctx, bson.M{"branch_id": id}).Decode(&branch); err != nil {
			utils.Error(c, utils.InternalServerError, "Can't find branch with id")
			return
		}

		var newStatus models.Status

		if branch.Status == models.Active {
			newStatus = models.Inactive
		} else {
			newStatus = models.Active
		}

		updated_at, _ := helpers.GetTime()

		updateObj := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "status", Value: newStatus},
				{Key: `updated_at`, Value: updated_at},
			}},
		}

		filter := bson.M{"branch_id": id}
		usert := true
		options := options.UpdateOptions{
			Upsert: &usert,
		}

		_, err := database.BranchCollection.UpdateOne(ctx, filter, updateObj, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update status")
			return
		}
		utils.Message(c, "Branch Status updated")
	}
}

// * DONE
func validateBranchDetails(branch *models.Branch) (string, bool) {

	if branch.Branch_Name == "" {
		return "Branch name is required", false
	}

	if branch.Address == "" {
		return "Address is required", false
	}

	if len(strconv.Itoa(branch.Phone)) != 10 {
		return "Phone number must be 10 digits", false
	}

	if branch.Email == "" {
		return "Email address required", false
	} else if !utils.ValidateEmail(branch.Email) {
		return "Invalid email address", false
	}

	if branch.City == "" {
		return "City is required", false
	}

	if branch.State == "" {
		return "State is required", false
	}

	if branch.Country == "" {
		return "Country is required", false
	}

	if len(strconv.Itoa(branch.Pincode)) != 6 {
		return "Pincode must be 6 digits", false
	}

	return "", true
}

// * DONE
func GetBranchesByStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()

		code, err := strconv.Atoi(c.Param("status"))
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't parse status")
			return
		}
		status := "active"
		if code == 0 {
			status = "inactive"
		}

		// 1 means active, 0 means inactive

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 10 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		// Here the aggregation pipeline started

		stage := bson.A{
			bson.D{
				{Key: "$group",
					Value: bson.D{
						{Key: "_id", Value: primitive.Null{}},
						{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
					},
				},
			},
			bson.D{
				{Key: "$project",
					Value: bson.D{
						{Key: "_id", Value: 0},
						{Key: "branches",
							Value: bson.D{
								{Key: "$slice",
									Value: bson.A{
										bson.D{
											{Key: "$map",
												Value: bson.D{
													{Key: "input",
														Value: bson.D{
															{Key: "$filter",
																Value: bson.D{
																	{Key: "input", Value: "$data"},
																	{Key: "as", Value: "branch"},
																	{Key: "cond",
																		Value: bson.D{
																			{Key: "$eq",
																				Value: bson.A{
																					"$$branch.status",
																					status,
																				},
																			},
																		},
																	},
																},
															},
														},
													},
													{Key: "as", Value: "data"},
													{Key: "in", Value: bson.D{
														{Key: `branch_id`, Value: "$$data.branch_id"},
														{Key: "manager_id", Value: "$$data.manager_id"},
														{Key: "branch_name", Value: "$$data.branch_name"},
														{Key: "address", Value: "$$data.address"},
														{Key: "phone", Value: "$$data.phone"},
														{Key: "email", Value: "$$data.email"},
														{Key: "city", Value: "$$data.city"},
														{Key: "state", Value: "$$data.state"},
														{Key: "country", Value: "$$data.country"},
														{Key: "pincode", Value: "$$data.pincode"},
														{Key: "status", Value: "$$data.status"},
													}},
												},
											},
										},
										startIndex,
										recordPerPage,
									},
								},
							},
						},
					},
				},
			},
		}

		result, err := database.BranchCollection.Aggregate(ctx, stage)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Error while fetching branches  "+err.Error())
			return
		}

		var allBranches []bson.M
		if err := result.All(ctx, &allBranches); err != nil {
			utils.Error(c, utils.InternalServerError, "Error while getting the branches "+err.Error())
			return
		}
		if len(allBranches) == 0 {
			utils.Response(c, []interface{}{})
			return
		}
		utils.Response(c, allBranches[0])
	}
}
