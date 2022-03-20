package router

import (
	log "github.com/sirupsen/logrus"
	deploymentHandler "go-docker/internal/handler"
	"net/http"
)

type Router struct{
	handler deploymentHandler.Handler
}

func NewRouter() *Router{
	return &Router{
		handler: deploymentHandler.Handler{},
	}
}

func (r Router) InitDeploymentRouter() {
	handler := deploymentHandler.NewHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("/mutate", handler.MutateHandler)

	err := http.ListenAndServeTLS(":8443", "./server.crt", "./server.key", mux)
	if err != nil {
		log.Errorf("ListenAndServe: %v", err)
	}
}