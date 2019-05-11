package models

import (
	"gopkg.in/mgo.v2/bson"
)

type FileModel struct {
	Id       bson.ObjectId `bson:"_id,omitempty"`
	FileName string
	FileExt  string
}
