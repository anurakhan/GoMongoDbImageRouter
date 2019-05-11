package photoscontroller

import (
	keygen "github.com/anurakhan/go-mongo-lb-driver/key-gen-service"
	"github.com/anurakhan/go-mongo-lb-driver/models"
	hashring "github.com/anurakhan/go-mongo-lb-driver/ring"
	photosservice "github.com/anurakhan/go-mongo-lb-driver/services"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type PhotosController struct {
	Router        *mux.Router
	Ring          hashring.IConsistentHashRing
	KeyGenService keygen.IKeyGenService
}

func (controller *PhotosController) RegisterHandlers() {
	controller.Router.HandleFunc("/photos/{id}", controller.getPhotosByIDEndpoint).Methods("GET")
	controller.Router.HandleFunc("/photos", controller.postPhotoEndpoint).Methods("POST")
}

func (controller *PhotosController) getPhotosByIDEndpoint(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	query := models.FileId{Id: vars["id"]}

	service := &photosservice.PhotosService{Ring: controller.Ring,
		KeyGenService: controller.KeyGenService}

	fileRetModel := service.GetPhoto(query.Id)

	RetImage(w, fileRetModel, http.StatusOK)
}

func (controller *PhotosController) postPhotoEndpoint(w http.ResponseWriter, req *http.Request) {
	var Buf bytes.Buffer
	// in your case file would be fileupload
	file, header, err := req.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Println("YOYOYOYY")

	name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", name[0])
	// Copy the file data to my buffer
	io.Copy(&Buf, file)

	service := &photosservice.PhotosService{Ring: controller.Ring,
		KeyGenService: controller.KeyGenService}

	id := service.PostPhoto(&Buf, name[0], name[1])

	Json(w, http.StatusOK, models.FileId{Id: id})

	Buf.Reset()
}

func Json(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func Error(w http.ResponseWriter, code int, message string) {
	Json(w, code, map[string]string{"error": message})
}

func RetImage(w http.ResponseWriter, retFile *models.FileRetModel, code int) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(retFile.Data)))

	if _, err := w.Write(retFile.Data); err != nil {
		panic(err)
	}
}
