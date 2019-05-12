package main

import (
	"log"
	"net/http"

	photoscontroller "github.com/anurakhan/go-mongo-lb-driver/controllers"
	keygen "github.com/anurakhan/go-mongo-lb-driver/key-gen-service"
	hashring "github.com/anurakhan/go-mongo-lb-driver/ring"
	"github.com/anurakhan/go-mongo-lb-driver/server"
	"github.com/gorilla/mux"
)

func main() {

	var hashRing hashring.IConsistentHashRing = new(hashring.ConsistentHashRing)
	hashRing.InitRing()
	hashRing.AddServer(&server.Server{Name: "MongoDb1", Address: "localhost:27017",
		FilePath: "FileSystem1", KeyForCh: []byte("\\؁�Н{�y�")})
	hashRing.AddServer(&server.Server{Name: "MongoDb2", Address: "localhost:27018",
		FilePath: "FileSystem2", KeyForCh: []byte("\\؂%Н{$a�y")})
	hashRing.AddServer(&server.Server{Name: "MongoDb3", Address: "localhost:27019",
		FilePath: "FileSystem3", KeyForCh: []byte("\\؂PН{��ga")})

	router := mux.NewRouter()
	photosController := &photoscontroller.PhotosController{Router: router,
		Ring: hashRing, KeyGenService: new(keygen.BsonKeyGenService)}
	photosController.RegisterHandlers()

	log.Fatal(http.ListenAndServe(":80", router))
}
