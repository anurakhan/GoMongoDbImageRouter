package photosservice

import (
	"bytes"

	"github.com/anurakhan/go-mongo-lb-driver/models"

	dal "github.com/anurakhan/go-mongo-lb-driver/data-access-layer"
	keygen "github.com/anurakhan/go-mongo-lb-driver/key-gen-service"
	hashring "github.com/anurakhan/go-mongo-lb-driver/ring"
)

type PhotosService struct {
	Ring          hashring.IConsistentHashRing
	KeyGenService keygen.IKeyGenService
}

func (service *PhotosService) PostPhoto(buf *bytes.Buffer, fileName string, fileExt string) string {
	id := service.KeyGenService.GenKey()
	server := service.Ring.GetServerForKey(string(id))

	repo := &dal.PhotosRepository{Server: server, DbInteractor: &dal.MongoInteractor{Server: server}}
	return repo.PostPhoto(buf, fileName, fileExt, id)
}

func (service *PhotosService) GetPhoto(id string) *models.FileRetModel {
	server := service.Ring.GetServerForKey(id)

	repo := &dal.PhotosRepository{Server: server, DbInteractor: &dal.MongoInteractor{Server: server}}
	fileRetModel := repo.GetPhotoById(id)
	return fileRetModel
}
