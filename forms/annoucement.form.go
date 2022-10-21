package forms

type AnnoucementForm struct {
	Annoucement string   `json:"annoucement" binding:"min=1,required" validate:"required" minLength:"1"`
	Files       []string `json:"files" binding:"dive"`
}
