package keygen

import (
	"encoding/hex"

	"github.com/segmentio/ksuid"
	"gopkg.in/mgo.v2/bson"
)

type IKeyGenService interface {
	GenKey() []byte
}

type KSuidKeyGenService struct {
}

func (service *KSuidKeyGenService) GenKey() []byte {
	id := ksuid.New()
	return id[:]
}

type BsonKeyGenService struct {
}

func (service *BsonKeyGenService) GenKey() []byte {
	id := bson.NewObjectId()
	data, err := hex.DecodeString(id.Hex())
	if err != nil {
		panic(err)
	}
	return data
}
