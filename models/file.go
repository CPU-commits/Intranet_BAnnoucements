package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const FILES_COLLECTION = "files"

// Type permissions
const (
	PRIVATE   = "private"
	CLASSROOM = "public_classroom"
	PUBLIC    = "public"
)

type File struct {
	ID    primitive.ObjectID `json:"_id" bson:"_id"`
	Type  string             `json:"type" bson:"type"`
	Title string             `json:"title" bson:"title"`
}
