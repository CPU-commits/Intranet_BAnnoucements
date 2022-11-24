package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const FILES_COLLECTION = "files"

// Type permissions
const (
	PRIVATE   = "private"
	CLASSROOM = "public_classroom"
	PUBLIC    = "public"
)

// File model info
// @Description File attached to annoucement
type File struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id" example:"6376c8283cc695e19d785b08"`
	Type  string             `json:"type" bson:"type" example:"private"`
	Title string             `json:"title" bson:"title" example:"Title"`
}
