package dal

import (
	"github.com/anurakhan/go-mongo-lb-driver/models"
	"github.com/anurakhan/go-mongo-lb-driver/server"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoInteractor struct {
	Server     *server.Server
	session    *mgo.Session
	collection *mgo.Collection
}

func (interactor *MongoInteractor) StartConn() {
	session, err := mgo.Dial(interactor.Server.Address)

	if err != nil {
		panic(err)
	}

	interactor.session = session
	interactor.collection = session.DB("PhotoInjest").C("PhotoMeta")
}

func (interactor *MongoInteractor) CloseConn() {
	interactor.session.Close()
}

func (interactor *MongoInteractor) InsertFileInfo(fileModel *models.FileModel) {
	col := interactor.collection
	err := col.Insert(fileModel)

	if err != nil {
		panic(err)
	}
}

func (interactor *MongoInteractor) GetFiles() []models.FileModel {
	col := interactor.collection
	var res []models.FileModel
	col.Find(nil).All(&res)
	return res
}

func (interactor *MongoInteractor) GetFileById(id bson.ObjectId) models.FileModel {
	col := interactor.collection
	var res models.FileModel
	col.FindId(id).One(&res)
	return res
}
