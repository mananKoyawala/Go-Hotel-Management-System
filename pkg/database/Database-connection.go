package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBInstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongodb := os.Getenv("MONGODB_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodb))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client
}

var client *mongo.Client = DBInstance()

func OpenCollection(collectionName string) *mongo.Collection {
	return client.Database("hotel-managment-system").Collection(collectionName)
}

var (
	AdminCollection             *mongo.Collection = OpenCollection("admin")
	BranchCollection            *mongo.Collection = OpenCollection("branch")
	DriverCollection            *mongo.Collection = OpenCollection("driver")
	GuestFeedbackCollection     *mongo.Collection = OpenCollection("guest-feedback")
	GuestCollection             *mongo.Collection = OpenCollection("guest")
	ManagerCollection           *mongo.Collection = OpenCollection("manager")
	PickupServiceCollection     *mongo.Collection = OpenCollection("pickup-service")
	ReservationCollection       *mongo.Collection = OpenCollection("reservation")
	RoomCollection              *mongo.Collection = OpenCollection("room")
	StaffCollection             *mongo.Collection = OpenCollection("staff")
	VerificationEmailCollection *mongo.Collection = OpenCollection("verification-email")
)
