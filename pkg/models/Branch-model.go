package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Branch struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Branch_id   string             `bson:"branch_id,omitempty" json:"branch_id,omitempty"`
	Manager_id  string             `bson:"manager_id,omitempty" json:"manager_id,omitempty"`
	Branch_Name string             `bson:"branch_name,omitempty" json:"branch_name,omitempty"`
	Address     string             `bson:"address,omitempty" json:"address,omitempty"`
	Phone       int                `bson:"phone,omitempty" json:"phone,omitempty"`
	Email       string             `bson:"email,omitempty" json:"email,omitempty"`
	City        string             `bson:"city,omitempty" json:"city,omitempty"`
	State       string             `bson:"state,omitempty" json:"state,omitempty"`
	Country     string             `bson:"country,omitempty" json:"country,omitempty"`
	Pincode     int                `bson:"pincode,omitempty" json:"pincode,omitempty"`
	Status      Status             `bson:"status,omitempty" json:"status,omitempty"`
	Total_Rooms int                `bosn:"total_rooms,omitempty" json:"total_rooms,omitempty"`
	// Images      []string           `bson:"images,omitempty" json:"images,omitempty"`
	Created_at time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_at time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
