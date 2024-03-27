package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Staff struct {
	ID            primitive.ObjectID `bson:"_id"`
	Staff_id      string             `bson:"staff_id,omitempty" json:"staff_id,omitempty"`
	First_Name    string             `bson:"first_name,omitempty" json:"first_name,omitempty"`
	Last_Name     string             `bson:"last_name,omitempty" json:"last_name,omitempty"`
	Age           int                `bson:"age,omitempty" json:"age,omitempty"`
	Gender        string             `bson:"gender,omitempty" json:"gender,omitempty"`
	Job_Type      string             `bson:"job_type,omitempty" json:"job_type,omitempty"`
	Salary        float64            `bson:"salary,omitempty" json:"salary,omitempty"`
	Aadhar_Number string             `bson:"aadhar_number,omitempty" json:"aadhar_number,omitempty"`
	Email         string             `bson:"email,omitempty" json:"email,omitempty"`
	Phone         int                `bson:"phone,omitempty" json:"phone,omitempty"`
	Status        Status             `bson:"status,omitempty" json:"status,omitempty"`
	// Image         string             `bson:"image,omitempty" json:"image,omitempty"`
	Created_at time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_at time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// Staff has only one profile photo
