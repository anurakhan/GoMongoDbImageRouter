package dal

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/anurakhan/go-mongo-lb-driver/models"
	"github.com/anurakhan/go-mongo-lb-driver/server"
	"gopkg.in/mgo.v2/bson"
	"encoding/hex"
)

type PhotosRepository struct {
	Server       *server.Server
	DbInteractor *MongoInteractor
}

func (repo *PhotosRepository) PostPhoto(buf *bytes.Buffer, fileName string, fileExt string, id []byte) string {
	path := initPath(repo)

	file, err := os.Create(path + "/" + string(id) + "." + fileExt)
	if err != nil {
		panic(err)
	}

	io.Copy(file, buf)

	interactor := repo.DbInteractor

	interactor.StartConn()

	interactor.InsertFileInfo(&models.FileModel{
		Id:       bson.ObjectId(string(id)),
		FileName: fileName,
		FileExt:  fileExt})

	interactor.CloseConn()

	return hex.EncodeToString(id)
}

func (repo *PhotosRepository) GetPhotoById(id string) *models.FileRetModel {
	path := initPath(repo)
	fmt.Println(path)
	interactor := repo.DbInteractor

	interactor.StartConn()

	fileModel := interactor.GetFileById(bson.ObjectId(id))

	interactor.CloseConn()
	fmt.Println(path + "/" + id + "." + fileModel.FileExt)
	file, err := os.Open(path + "/" + id + "." + fileModel.FileExt)

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return &models.FileRetModel{
		Id:       id,
		FileName: fileModel.FileName,
		FileExt:  fileModel.FileExt,
		Data:     data}
}

func initPath(repo *PhotosRepository) string {
	path := repo.Server.FilePath

	path = fromFileSystemDir(path)

	fmt.Println(path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	return path
}

func fromFileSystemDir(path string) string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir + "/go-mongo-lb-driver-file-system/" + path
}
