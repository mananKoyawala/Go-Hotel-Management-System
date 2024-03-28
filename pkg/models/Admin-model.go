package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admin struct {
	ID            primitive.ObjectID `bson:"_id"`
	Admin_id      string             `bson:"admin_id,omitempty" json:"admin_id,omitempty"`
	First_Name    string             `bson:"first_name,omitempty" json:"first_name,omitempty"`
	Last_Name     string             `bson:"last_name,omitempty" json:"last_name,omitempty"`
	Email         string             `bson:"email,omitempty" json:"email,omitempty"`
	Password      string             `bson:"Password,omitempty" json:"Password,omitempty"`
	Access_Type   Access_Type        `bson:"access_type,omitempty" json:"access_type,omitempty"`
	Token         string             `bson:"token,omitempty" json:"token,omitempty"`
	Refresh_Token string             `bson:"refresh_token,omitempty" json:"refresh_token,omitempty"`
	Created_at    time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_at    time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}
