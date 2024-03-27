package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Driver struct {
	ID               primitive.ObjectID `bson:"_id"`
	Driver_id        string             `bson:"driver_id,omitempty" json:"driver_id,omitempty"`
	First_Name       string             `bson:"first_name,omitempty" json:"first_name,omitempty"`
	Last_Name        string             `bson:"last_name,omitempty" json:"last_name,omitempty"`
	Age              int                `bson:"age,omitempty" json:"age,omitempty"`
	Gender           string             `bson:"gender,omitempty" json:"gender,omitempty"`
	Car_Compamy      string             `bson:"car_compamy,omitempty" json:"car_compamy,omitempty"`
	Car_Model        string             `bson:"car_model,omitempty" json:"car_model,omitempty"`
	Car_Number_Plate string             `bson:"car_number_plate,omitempty" json:"car_number_plate,omitempty"`
	Status           Status             `bson:"status,omitempty" json:"status,omitempty"`       // Active or Inactive
	Availablity      Availablity        `bson:"available,omitempty" json:"available,omitempty"` // available or not available (on the drive)
	Pickup_Location  string             `bson:"pickup_location,omitempty" json:"pickup_location,omitempty"`
	Salary           float64            `bson:"salary,omitempty" json:"salary,omitempty"`
	Phone            int                `bson:"phone,omitempty" json:"phone,omitempty"`
	Token            string             `bson:"token,omitempty" json:"token,omitempty"`
	Refresh_Token    string             `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
	// Image            string             `bson:"image,omitempty" json:"image,omitempty"`
	Created_at time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_at time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// Driver has only one profile photo
