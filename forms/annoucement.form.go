package forms

type AnnoucementForm struct {
	Annoucement string   `json:"annoucement" binding:"min=1,required" validate:"required" minLength:"1" example:"This is a annoucement!"`
	Files       []string `json:"files" binding:"dive" validate:"optional" example:"6376c8283cc695e19d785b08"`
}
