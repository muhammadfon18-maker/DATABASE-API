package main

import (
	"fmt"
	"log"
	"net/http"
	"project/internals/handlers"
	"project/internals/service"
	"project/internals/storage"
)

func main() {

	db, err := storage.DB()
	if err != nil {
		log.Fatal(err)
	}

	storage := storage.NewStorage(db)

	i := service.NewService(storage)

	handler := handlers.NewHttpService(i)

	http.HandleFunc("POST /user", handler.HttpCreate)
	http.HandleFunc("GET /user/{id}", handler.HttpGet)
	http.HandleFunc("GET /users", handler.HttpGetAll)
	http.HandleFunc("PUT /user/{id}", handler.HttpUpdate)
	http.HandleFunc("DELETE /user/{id}", handler.HttpDelete)

	fmt.Println("Server running on :8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("connected")

}

//var _ storage.UserStorage = (*storage.Storage)(nil)
