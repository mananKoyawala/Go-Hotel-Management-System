package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// * No Need of password beacuse we use Email authentication by Google using Firebase
type Guest struct {
	ID            primitive.ObjectID `bson:"_id"`
	Guest_id      string             `bson:"guest_id,omitempty" json:"guest_id,omitempty"`
	ID_Proof_Type ID_Proof_Type      `bson:"id_proof_type,omitempty" json:"id_proof_type,omitempty"`
	First_Name    string             `bson:"first_name,omitempty" json:"first_name,omitempty"`
	Last_Name     string             `bson:"last_name,omitempty" json:"last_name,omitempty"`
	Phone         int                `bson:"phone,omitempty" json:"phone,omitempty"`
	Email         string             `bson:"email,omitempty" json:"email,omitempty"`
	Password      string             `bson:"Password,omitempty" json:"Password,omitempty"`
	Gender        string             `bson:"gender,omitempty" json:"gender,omitempty"`
	Country       string             `bson:"country,omitempty" json:"country,omitempty"`
	Token         string             `bson:"token,omitempty" json:"token,omitempty"`
	Refresh_Token string             `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
	// Image         string             `bson:"image,omitempty" json:"image,omitempty"
	Access_Type Access_Type `bson:"access_type,omitempty" json:"access_type,omitempty"`
	Created_at  time.Time   `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_at  time.Time   `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// Guest has only one profile photot

// The "omitempty" effect the operation it only for encoding and decoding , so it is easy to send data
