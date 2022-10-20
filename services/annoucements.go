package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sync"

	commonFunc "github.com/CPU-commits/Intranet_BAnnoucements/common/func"
	"github.com/CPU-commits/Intranet_BAnnoucements/forms"
	"github.com/CPU-commits/Intranet_BAnnoucements/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	annoucementService *AnnoucementService
	lock               = &sync.Mutex{}
)

type AnnoucementService struct{}

func (a *AnnoucementService) GetAnnoucements(
	skip int,
	limit int,
	total bool,
) ([]models.AnnoucementWLookUp, int64, *ErrorRes) {
	annoucements, totalData, err := annoucementModel.GetAnnoucements(limit, skip, total)
	if err != nil {
		return nil, 0, &ErrorRes{
			Err:        err,
			StatusCode: http.StatusServiceUnavailable,
		}
	}
	return annoucements, totalData, nil
}

func (a *AnnoucementService) UploadAnnoucement(
	annoucement forms.AnnoucementForm,
	idUser string,
) (primitive.ObjectID, *ErrorRes) {
	idObjUser, err := primitive.ObjectIDFromHex(idUser)
	if err != nil {
		return primitive.NilObjectID, &ErrorRes{
			Err:        err,
			StatusCode: http.StatusNotAcceptable,
		}
	}
	if len(annoucement.Files) > 0 {
		// Make request
		permissionsRequest := make(map[string]interface{})
		permissionsRequest["id_user"] = idUser
		permissionsRequest["files"] = annoucement.Files
		data, err := json.Marshal(permissionsRequest)
		if err != nil {
			return primitive.NilObjectID, &ErrorRes{
				Err:        err,
				StatusCode: http.StatusBadRequest,
			}
		}

		// Request
		var permissions []string
		message, err := nats.Request("get_permissions_files", data)
		if err != nil {
			return primitive.NilObjectID, &ErrorRes{
				Err:        err,
				StatusCode: http.StatusServiceUnavailable,
			}
		}
		err = json.Unmarshal(message.Data, &permissions)
		if err != nil {
			return primitive.NilObjectID, &ErrorRes{
				Err:        err,
				StatusCode: http.StatusInternalServerError,
			}
		}
		someInaccessibleFile, _ := commonFunc.Some(func(permission interface{}) bool {
			v := reflect.ValueOf(permission)
			if v.Kind() == reflect.String {
				if v.Interface() == models.PRIVATE || v.Interface() == models.CLASSROOM {
					return true
				}
			}
			return false
		}, permissions)
		if someInaccessibleFile {
			return primitive.NilObjectID, &ErrorRes{
				Err:        fmt.Errorf("Todos los archivos a publicar deben ser públicos"),
				StatusCode: http.StatusConflict,
			}
		}
	}
	// Upload annoucement
	insertedAnnoucement, err := annoucementModel.UploadAnnoucement(
		annoucement,
		idObjUser,
	)
	if err != nil {
		return primitive.NilObjectID, &ErrorRes{
			Err:        err,
			StatusCode: http.StatusServiceUnavailable,
		}
	}
	return insertedAnnoucement.InsertedID.(primitive.ObjectID), nil
}

func (a *AnnoucementService) DeleteAnnoucement(idAnnoucement string) *ErrorRes {
	idObjAnnoucement, err := primitive.ObjectIDFromHex(idAnnoucement)
	if err != nil {
		return &ErrorRes{
			StatusCode: http.StatusBadRequest,
			Err:        err,
		}
	}

	// Delete
	err = annoucementModel.DeleteAnnoucement(idObjAnnoucement)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return &ErrorRes{
			StatusCode: http.StatusServiceUnavailable,
			Err:        err,
		}
	}
	return nil
}

func NewAnnoucementService() *AnnoucementService {
	if annoucementService == nil {
		lock.Lock()
		defer lock.Unlock()

		annoucementService = &AnnoucementService{}
	}
	return annoucementService
}
