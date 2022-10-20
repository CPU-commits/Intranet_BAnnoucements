package forms

type AnnoucementForm struct {
	Annoucement string   `json:"annoucement" binding:"min=1,required"`
	Files       []string `json:"files" binding:"dive"`
}
