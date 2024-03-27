package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reservation struct {
	ID                primitive.ObjectID `bson:"_id"`
	Reservation_id    string             `bson:"reservation_id,omitempty" json:"reservation_id,omitempty"`
	Room_id           string             `bson:"room_id,omitempty" json:"room_id,omitempty"`
	Guest_id          string             `bson:"guest_id,omitempty" json:"guest_id,omitempty"`
	Check_in_time     time.Time          `bson:"check_in_time,omitempty" json:"check_in_time,omitempty"`
	Check_out_time    time.Time          `bson:"check_out_time,omitempty" json:"check_out_time,omitempty"`
	Deposit_Amount    float64            `bson:"deposit_amount,omitempty" json:"deposit_amount,omitempty"`
	Pending_amount    float64            `bson:"pending_amount,omitempty" json:"pending_amount,omitempty"`
	Numbers_of_guests int                `bson:"numbers_of_guests,omitempty" json:"numbers_of_guests,omitempty"`
	Created_at        time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_at_at     time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
