package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID                primitive.ObjectID `bson:"_id"`
	Room_id           string             `bson:"room_id,omitempty" json:"room_id,omitempty"`
	Branch_id         string             `bson:"branch_id,omitempty" json:"branch_id,omitempty"`
	Room_Type         Room_Type          `bson:"room_type,omitempty" json:"room_type,omitempty"`
	Room_Availability Room_Availability  `bson:"room_availability,omitempty" json:"room_availability,omitempty"`
	Cleaning_Status   Cleaning_Status    `bson:"cleaning_status,omitempty" json:"cleaning_status,omitempty"`
	Price             float64            `bson:"price,omitempty" json:"price,omitempty"`
	// Images          []string           `bson:"images,omitempty" json:"images,omitempty"`
	Created_at time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_at time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}