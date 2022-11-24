package models

type UserTypes string

const USERS_COLLECTION = "users"

const (
	DIRECTOR          = "f"
	DIRECTIVE         = "e"
	TEACHER           = "d"
	ATTORNEY          = "c"
	STUDENT_DIRECTIVE = "b"
	STUDENT           = "a"
)

// User model info
// @Description User account information
// @Description . Global model for all users
type SimpleUser struct {
	ID             string `json:"_id" example:"6376c8283cc695e19d785b08"`
	Name           string `json:"name" bson:"name" example:"Name"`
	FirstLastname  string `json:"first_lastname" bson:"first_lastname" example:"Firstlastname"`
	SecondLastname string `json:"second_lastname" bson:"second_lastname" example:"Secondlastname"`
	Rut            string `json:"rut,omitempty" bson:"rut" extensions:"x-omitempty" example:"12345678-0"`
}
