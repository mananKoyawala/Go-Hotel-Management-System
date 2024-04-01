package models

// General enums
type Status string // Means Entity is Working or not ex- Manager, Staff, Driver

const (
	Active   Status = "active"
	Inactive Status = "inactive"
)

type Availablity string // Means Entity available fo work or not ex- Driver

const (
	Available   Availablity = "available"
	UnAvailable Availablity = "unavailable"
)

type Access_Type string

const (
	Admin_Access   Access_Type = "admin"
	Manager_Access Access_Type = "manager"
	Guest_Access   Access_Type = "guest"
	Driver_Access  Access_Type = "driver"
)

// ? GuestFeedback enums

type Feedback_Type string

const (
	Complaint Feedback_Type = "complaint"
	Rating    Feedback_Type = "rating"
)

// ? Guest enums

type ID_Proof_Type string

const (
	Aadhar_Card     ID_Proof_Type = "aadhar_card"
	PassPort        ID_Proof_Type = "passport"
	Pan_Card        ID_Proof_Type = "pan_card"
	Driving_License ID_Proof_Type = "driving_license"
)

// ? ROOM enums
type Room_Type string

const (
	Single_Bad Room_Type = "single_bed"
	Double_Bad Room_Type = "double_bed"
	Suite      Room_Type = "suite"
)

type Cleaning_Status string

const (
	Cleaned    Cleaning_Status = "cleaned"
	Dirty      Cleaning_Status = "dirty"
	InProgress Cleaning_Status = "inprogress"
)

type Room_Availability string

const (
	Room_Available   Room_Availability = "available"
	Room_Unavailable Room_Availability = "occupied"
)

// ? PickupService enums
type PickupStatus string

const (
	Completed    PickupStatus = "completed"
	NotCompleted PickupStatus = "not_completed"
)

// ? Feedback enums
type FeedbackStatus string

const (
	Resolved FeedbackStatus = "resolved"
	Pending  FeedbackStatus = "pending"
)
