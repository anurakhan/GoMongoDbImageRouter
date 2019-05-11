package models

import (
	"gopkg.in/mgo.v2/bson"
)

type FileModel struct {
	Id       bson.ObjectId
	FileName string
	FileExt  string
}
