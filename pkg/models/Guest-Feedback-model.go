package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GuestFeedback struct {
	ID                 primitive.ObjectID `bson:"_id"`
	Guest_Feedback_id  string             `bson:"guest_feedback_id,omitempty" json:"guest_feedback_id,omitempty"`
	Guest_id           string             `bson:"guest_id,omitempty" json:"guest_id,omitempty"`
	Description        string             `bson:"description,omitempty" json:"description,omitempty"`
	Feedback_Type      Feedback_Type      `bson:"feedback_type,omitempty" json:"feedback_type,omitempty"`
	Resolution_Details string             `bson:"resolution_details,omitempty" json:"resolution_details,omitempty"` //Reply from hotel management
	// Images             []string           `bson:"images,omitempty" json:"images,omitempty"`
	Rating     int       `bson:"rating,omitempty" json:"rating,omitempty"`
	Created_at time.Time `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_at time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
