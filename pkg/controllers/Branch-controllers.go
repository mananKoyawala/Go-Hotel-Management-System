package controllers

import (
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

var branchFolder = "branch"
var maxFiles = 5

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
						{Key: "images", Value: "$$data.images"},
						{Key: "total_rooms", Value: "$$data.total_rooms"},
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

		branch.Manager_id = c.PostForm("manager_id")
		branch.Branch_Name = c.PostForm("branch_name")
		branch.Address = c.PostForm("address")
		branch.Phone, _ = strconv.Atoi(c.PostForm("phone"))
		branch.Email = c.PostForm("email")
		branch.City = c.PostForm("city")
		branch.State = c.PostForm("state")
		branch.Country = c.PostForm("country")
		branch.Pincode, _ = strconv.Atoi(c.PostForm("pincode"))

		// Validate Details
		msg, isVal := validateBranchDetails(branch)
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

		// Generate IDs and Timestamps
		branch.ID = primitive.NewObjectID()
		branch.Branch_id = branch.ID.Hex()
		branch.Status = models.Active
		branch.Total_Rooms = 0 // When branch is created there are no rooms that associated with it
		// If admin add room that it will counted here
		branch.Created_at, _ = helpers.GetTime()
		branch.Updated_at, _ = helpers.GetTime()

		// upload files here
		for _, fileHeader := range c.Request.MultipartForm.File["file"] {
			name := strings.ReplaceAll(fileHeader.Filename, " ", "")
			file, _ := fileHeader.Open()
			filename := fmt.Sprintf("%d_%s", time.Now().Unix(), name)
			url, err := imageupload.UploadService(file, branchFolder, filename)
			if err != nil {
				utils.Error(c, utils.InternalServerError, "Can't uplaod the image.")
				return
			}
			branch.Images = append(branch.Images, url)
		}

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
		var branch models.Branch
		id := c.Param("id")

		if err := database.BranchCollection.FindOne(ctx, bson.M{"branch_id": id}).Decode(&branch); err != nil {
			utils.Error(c, utils.BadRequest, "Can't find branch with id")
			return
		}

		for _, i := range branch.Images {
			image := utils.GetTrimedUrl(i)
			if err := imageupload.DeleteService(image); err != nil {
				utils.Error(c, utils.InternalServerError, "Can't delete image."+err.Error())
				return
			}

		}
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

		branch.Manager_id = c.PostForm("manager_id")
		branch.Branch_Name = c.PostForm("branch_name")
		branch.Address = c.PostForm("address")
		branch.Phone, _ = strconv.Atoi(c.PostForm("phone"))
		branch.Email = c.PostForm("email")
		branch.City = c.PostForm("city")
		branch.State = c.PostForm("state")
		branch.Country = c.PostForm("country")
		branch.Pincode, _ = strconv.Atoi(c.PostForm("pincode"))

		// Validate data
		msg, val := validateBranchDetails(branch)
		if !val {
			utils.Error(c, utils.BadRequest, msg)
			return
		}

		var updateObj primitive.D

		// Update details if not empty
		if branch.Manager_id != "" {
			updateObj = append(updateObj, bson.E{Key: "manager_id", Value: branch.Manager_id})
		}

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
			utils.Error(c, utils.NotFound, "Can't find branch with id")
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
func validateBranchDetails(branch models.Branch) (string, bool) {

	if branch.Manager_id == "" {
		return "Manager id is required", false
	}

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
														{Key: "images", Value: "$$data.images"},
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

// * DONE
func BranchAddImage() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := helpers.GetContext()
		defer cancel()
		var branch models.Branch
		id := c.Param("id")

		if err := database.BranchCollection.FindOne(ctx, bson.M{"branch_id": id}).Decode(&branch); err != nil {
			utils.Error(c, utils.NotFound, "Can't find branch with id.")
			return
		}

		if len(branch.Images) >= 5 {
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
		url, err := imageupload.UploadService(file, branchFolder, filename)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't uplaod the image.")
			return
		}
		branch.Images = append(branch.Images, url)
		updated_at, _ := helpers.GetTime()

		filter := bson.M{"branch_id": id}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}
		updateObj := bson.D{{Key: "$set",
			Value: bson.D{
				{Key: "images", Value: branch.Images},
				{Key: "updated_at", Value: updated_at},
			}}}

		_, err = database.BranchCollection.UpdateOne(ctx, filter, updateObj, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update image")
			return
		}
		utils.Message(c, "Image updated successfully.")
	}
}

// * DONE
func BranchRemoveImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()
		var branch models.Branch
		id := c.Param("id")
		imageUrl := c.PostForm("image")

		if err := database.BranchCollection.FindOne(ctx, bson.M{"branch_id": id}).Decode(&branch); err != nil {
			utils.Error(c, utils.NotFound, "Can't find branch with id.")
			return
		}

		if len(branch.Images) <= 0 {
			utils.Error(c, utils.BadRequest, "Branch doesn't have any images.")
			return
		}

		if imageUrl == "" {
			utils.Error(c, utils.BadRequest, "Please provide image url.")
			return
		}

		found := false
		for i, img := range branch.Images {
			if img == imageUrl {
				found = true
				// remove imageUrl from branch.Images
				branch.Images = append(branch.Images[:i], branch.Images[i+1:]...)
				break
			}
		}

		if !found {
			utils.Error(c, utils.BadRequest, "Can't find image in branch.")
			return
		}

		image := utils.GetTrimedUrl(imageUrl)
		if err := imageupload.DeleteService(image); err != nil {
			utils.Error(c, utils.InternalServerError, "Can't delete image."+err.Error())
			return
		}

		updated_at, _ := helpers.GetTime()

		filter := bson.M{"branch_id": id}
		upsert := true
		options := options.UpdateOptions{
			Upsert: &upsert,
		}
		updateObj := bson.D{{Key: "$set",
			Value: bson.D{
				{Key: "images", Value: branch.Images},
				{Key: "updated_at", Value: updated_at},
			}}}

		_, err := database.BranchCollection.UpdateOne(ctx, filter, updateObj, &options)
		if err != nil {
			utils.Error(c, utils.InternalServerError, "Can't update image")
			return
		}
		utils.Message(c, "Image deleted successfully.")
	}
}

// * DONE
func SearchBranchData() gin.HandlerFunc {
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
				bson.D{{Key: "branch_name", Value: bson.D{{Key: "$regex", Value: search_string}, {Key: "$options", Value: "i"}}}},
				bson.D{{Key: "address", Value: bson.D{{Key: "$regex", Value: search_string}, {Key: "$options", Value: "i"}}}},
				bson.D{{Key: "status", Value: bson.D{{Key: "$regex", Value: search_string}, {Key: "$options", Value: "i"}}}},
				bson.D{{Key: "city", Value: bson.D{{Key: "$regex", Value: search_string}, {Key: "$options", Value: "i"}}}},
				bson.D{{Key: "state", Value: bson.D{{Key: "$regex", Value: search_string}, {Key: "$options", Value: "i"}}}},
				bson.D{{Key: "country", Value: bson.D{{Key: "$regex", Value: search_string}, {Key: "$options", Value: "i"}}}},
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
						{Key: "images", Value: "$$data.images"},
						{Key: "total_rooms", Value: "$$data.total_rooms"},
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
			matchStage, groupStage, projectStage1, projectStage2,
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
func FilterBranch() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := helpers.GetContext()
		defer cancel()

		city := c.PostForm("city")
		state := c.PostForm("state")
		country := c.PostForm("country")
		status := models.Status(c.PostForm("status"))

		// clean the data
		if status == "" || (status != "active" && status != "inactive") {
			status = "active"
		}

		// log.Println(city)
		// log.Println(state)
		// log.Println(country)
		// log.Println(status)

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
				bson.D{{Key: "city", Value: bson.D{{Key: "$regex", Value: city}, {Key: "$options", Value: "i"}}}},
				bson.D{{Key: "state", Value: bson.D{{Key: "$regex", Value: state}, {Key: "$options", Value: "i"}}}},
				bson.D{{Key: "country", Value: bson.D{{Key: "$regex", Value: country}, {Key: "$options", Value: "i"}}}},
			}},
			{Key: "status", Value: status},
		}}}

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
						{Key: "images", Value: "$$data.images"},
						{Key: "total_rooms", Value: "$$data.total_rooms"},
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
			matchStage, groupStage, projectStage1, projectStage2,
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
