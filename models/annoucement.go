package models

import (
	"sync"
	"time"

	"github.com/CPU-commits/Intranet_BAnnoucements/db"
	"github.com/CPU-commits/Intranet_BAnnoucements/forms"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	annoucementModel *AnnoucementModel
	lock             = &sync.RWMutex{}
)

const ANNOUCEMENT_COLLECTION = "annoucements"

type Annoucement struct {
	ID          primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	User        primitive.ObjectID   `json:"user" bson:"user"`
	Files       []primitive.ObjectID `json:"files" bson:"files,omitempty"`
	Annoucement string               `json:"annoucement" bson:"annoucement"`
	UploadDate  primitive.DateTime   `json:"upload_date" bson:"upload_date"`
	UpdateDate  primitive.DateTime   `json:"update_date" bson:"update_date"`
}

type AnnoucementWLookUp struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	User        SimpleUser         `json:"user" bson:"user"`
	Files       []File             `json:"files" bson:"files,omitempty"`
	Annoucement string             `json:"annoucement" bson:"annoucement"`
	UploadDate  primitive.DateTime `json:"upload_date" bson:"upload_date"`
	UpdateDate  primitive.DateTime `json:"update_date" bson:"update_date"`
}

type AnnoucementModel struct {
	CollectionName string
}

func init() {
	collections, errC := DbConnect.GetCollections()
	if errC != nil {
		panic(errC)
	}
	for _, collection := range collections {
		if collection == ANNOUCEMENT_COLLECTION {
			return
		}
	}
	var jsonSchema = bson.M{
		"bsonType": "object",
		"required": []string{
			"annoucement",
			"user",
			"upload_date",
			"update_date",
		},
		"properties": bson.M{
			"annoucement": bson.M{"bsonType": "string"},
			"user":        bson.M{"bsonType": "objectId"},
			"upload_date": bson.M{"bsonType": "date"},
			"update_date": bson.M{"bsonType": "date"},
			"files": bson.M{
				"bsonType": bson.A{"array"},
				"items": bson.M{
					"bsonType": "objectId",
				},
			},
		},
	}
	var validators = bson.M{
		"$jsonSchema": jsonSchema,
	}
	opts := &options.CreateCollectionOptions{
		Validator: validators,
	}
	err := DbConnect.CreateCollection(ANNOUCEMENT_COLLECTION, opts)
	if err != nil {
		panic(err)
	}
}

func (a *AnnoucementModel) GetLookupUserNFiles() (bson.D, bson.D) {
	lookUpUser := bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from":         USERS_COLLECTION,
			"localField":   "user",
			"foreignField": "_id",
			"as":           "user",
			"pipeline": bson.A{bson.D{{
				Key: "$project",
				Value: bson.M{
					"name":            1,
					"first_lastname":  1,
					"second_lastname": 1,
				},
			}}},
		},
	}}
	lookUpFiles := bson.D{{
		Key: "$lookup",
		Value: bson.M{
			"from":         FILES_COLLECTION,
			"localField":   "files",
			"foreignField": "_id",
			"as":           "files",
			"pipeline": bson.A{bson.D{{
				Key: "$project",
				Value: bson.M{
					"type":  1,
					"title": 1,
				},
			}}},
		},
	}}
	return lookUpUser, lookUpFiles
}

func (a *AnnoucementModel) GetAnnoucementById(idAnnoucement primitive.ObjectID) (*AnnoucementWLookUp, error) {
	matchP := bson.D{{
		Key: "$match",
		Value: bson.M{
			"_id": idAnnoucement,
		},
	}}
	lookUpUser, lookUpFiles := a.GetLookupUserNFiles()
	setP := bson.D{{
		Key: "$set",
		Value: bson.M{
			"user": bson.M{
				"$first": "$user",
			},
		},
	}}

	var annoucements []AnnoucementWLookUp
	cursor, err := DbConnect.GetCollection(a.CollectionName).Aggregate(db.Ctx, mongo.Pipeline{
		matchP,
		lookUpUser,
		lookUpFiles,
		setP,
	})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(db.Ctx, &annoucements); err != nil {
		return nil, err
	}
	if len(annoucements) == 0 {
		return nil, nil
	}
	return &annoucements[0], nil
}

func (a *AnnoucementModel) GetAnnoucements(
	limit int,
	skip int,
	total bool,
) ([]AnnoucementWLookUp, int64, error) {
	sortP := bson.D{{
		Key: "$sort",
		Value: bson.M{
			"upload_date": -1,
		},
	}}
	limitP := bson.D{{
		Key:   "$limit",
		Value: limit,
	}}
	skipP := bson.D{{
		Key:   "$skip",
		Value: skip,
	}}
	lookUpUser, lookUpFiles := a.GetLookupUserNFiles()
	setP := bson.D{{
		Key: "$set",
		Value: bson.M{
			"user": bson.M{
				"$first": "$user",
			},
		},
	}}

	var annoucements []AnnoucementWLookUp
	cursor, err := DbConnect.GetCollection(a.CollectionName).Aggregate(db.Ctx, mongo.Pipeline{
		sortP,
		skipP,
		limitP,
		lookUpUser,
		lookUpFiles,
		setP,
	})
	if err != nil {
		return nil, 0, err
	}
	if err := cursor.All(db.Ctx, &annoucements); err != nil {
		return nil, 0, err
	}
	// Get total
	var totalData int64
	if total {
		totalData, err = DbConnect.GetCollection(a.CollectionName).CountDocuments(db.Ctx, bson.D{})
		if err != nil {
			return nil, 0, err
		}
	}
	return annoucements, totalData, nil
}

func (a *AnnoucementModel) UploadAnnoucement(
	annoucement forms.AnnoucementForm,
	idUser primitive.ObjectID,
) (*mongo.InsertOneResult, error) {

	now := primitive.NewDateTimeFromTime(time.Now())
	modelAnnoucement := &Annoucement{
		Annoucement: annoucement.Annoucement,
		User:        idUser,
		UploadDate:  now,
		UpdateDate:  now,
	}
	if len(annoucement.Files) > 0 {
		var annoucementFiles []primitive.ObjectID
		for _, file := range annoucement.Files {
			idObjFile, err := primitive.ObjectIDFromHex(file)
			if err != nil {
				return nil, err
			}
			annoucementFiles = append(annoucementFiles, idObjFile)
		}
		modelAnnoucement.Files = annoucementFiles
	}

	return DbConnect.GetCollection(a.CollectionName).InsertOne(
		db.Ctx,
		modelAnnoucement,
		&options.InsertOneOptions{},
	)
}

func (a *AnnoucementModel) DeleteAnnoucement(idAnnoucement primitive.ObjectID) error {
	_, err := DbConnect.GetCollection(a.CollectionName).DeleteOne(
		db.Ctx,
		bson.D{{
			Key:   "_id",
			Value: idAnnoucement,
		}},
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func NewAnnoucementModel() *AnnoucementModel {
	if annoucementModel == nil {
		lock.Lock()
		defer lock.Unlock()

		annoucementModel = &AnnoucementModel{
			CollectionName: ANNOUCEMENT_COLLECTION,
		}
	}
	return annoucementModel
}
