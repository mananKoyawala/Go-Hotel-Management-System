package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PickupService struct {
	ID                primitive.ObjectID `bson:"_id"`
	Pickup_Service_id string             `bson:"pickup_service_id,omitempty" json:"pickup_service_id,omitempty"`
	Guest_id          string             `bson:"guest_id,omitempty" json:"guest_id,omitempty"`
	Branch_id         string             `bson:"branch_id,omitempty" json:"branch_id,omitempty"`
	Driver_id         string             `bson:"driver_id,omitempty" json:"driver_id,omitempty"`
	Pickup_Location   string             `bson:"pickup_location,omitempty" json:"pickup_location,omitempty"`
	Pickup_Time       string             `bson:"pickup_time,omitempty" json:"pickup_time,omitempty"`
	Amount            float64            `bson:"amount,omitempty" json:"amount,omitempty"`
	Status            PickupStatus       `bson:"status,omitempty" json:"status,omitempty"`
	Created_at        time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_at        time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
