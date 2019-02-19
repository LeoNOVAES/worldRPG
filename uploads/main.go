package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

func up(w http.ResponseWriter, r *http.Request) {
	file, handle, err := r.FormFile("file")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	mimeType := handle.Header.Get("Content-Type")
	switch mimeType {
	case "image/jpeg":
		saveFile(w, file, handle)
	default:
		jsonRes(w, http.StatusBadRequest, "Arquivo nao suportado")
	}

}

func saveFile(w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader) {
	data, err := ioutil.ReadAll(file)

	if err != nil {
		log.Fatal(err)
		return
	}

	err = ioutil.WriteFile("./public/"+handle.Filename, data, 0666)

	if err != nil {
		log.Fatal(err)
		return
	}

	jsonRes(w, http.StatusCreated, "File uploaded successfully!.")
}

func jsonRes(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}

func main() {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 //8mb
	router.Static("/", "./public")
	http.HandleFunc("/uploads", up)

	http.ListenAndServe(":9000", nil)
}
