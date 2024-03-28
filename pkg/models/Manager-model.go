package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Manager struct {
	ID            primitive.ObjectID `bson:"_id"`
	Manager_id    string             `bson:"manager_id,omitempty" json:"manager_id,omitempty"`
	First_Name    string             `bson:"first_name,omitempty" json:"first_name,omitempty"`
	Last_Name     string             `bson:"last_name,omitempty" json:"last_name,omitempty"`
	Age           int                `bson:"age,omitempty" json:"age,omitempty"`
	Phone         int                `bson:"phone,omitempty" json:"phone,omitempty"`
	Email         string             `bson:"email,omitempty" json:"email,omitempty"`
	Password      string             `bson:"Password,omitempty" json:"Password,omitempty"`
	Gender        string             `bson:"gender,omitempty" json:"gender,omitempty"`
	Salary        float64            `bson:"salary,omitempty" json:"salary,omitempty"`
	Aadhar_Number string             `bson:"aadhar_number,omitempty" json:"aadhar_number,omitempty"`
	Status        Status             `bson:"status,omitempty" json:"status,omitempty"`
	Token         string             `bson:"token,omitempty" json:"token,omitempty"`
	Refresh_Token string             `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
	// Image         string             `bson:"image,omitempty" json:"image,omitempty"`
	Access_Type Access_Type `bson:"access_type,omitempty" json:"access_type,omitempty"`
	Created_at  time.Time   `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_at  time.Time   `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// Manager has only one profile photo
