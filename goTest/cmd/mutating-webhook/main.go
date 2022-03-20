package main

import (
	"go-docker/internal/router"
)

func main() {
	r := router.NewRouter()
	r.InitDeploymentRouter()
}