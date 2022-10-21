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
	ID             string `json:"_id"`
	Name           string `json:"name" bson:"name"`
	FirstLastname  string `json:"first_lastname" bson:"first_lastname"`
	SecondLastname string `json:"second_lastname" bson:"second_lastname"`
	Rut            string `json:"rut,omitempty" bson:"rut" extensions:"x-omitempty"`
}
