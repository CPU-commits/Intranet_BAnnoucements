package smaps

import "github.com/CPU-commits/Intranet_BAnnoucements/models"

type AnnoucementsMap struct {
	Annoucements []models.AnnoucementWLookUp `json:"annoucements"`
	Total        int64                       `json:"total"`
}
