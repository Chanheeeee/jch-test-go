package main

import (
	"go-docker/internal/handler"
	"net/http"
	"log"
)


func main() {
	
	err := http.ListenAndServeTLS(":8443", "./server.crt", "./server.key", handler.NewHandler() )
	//err := http.ListenAndServeTLS(":8443", "../../cert/server.crt", "../../cert/server.key", handler.NewHandler() )
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}else{
		log.Println("[IN]main.go")
	}
}