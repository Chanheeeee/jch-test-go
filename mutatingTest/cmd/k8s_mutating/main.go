package main

import (
    //"internal/pkg/"
	//"github.com/gin-gonic/gin"
)


/*
 Webhook needs to serve TLS
*/

cmd/main.go
func main() {
    // flags
    // ...
    server := http.NewServer(port)
    if err := server.ListenAndServeTLS(tlscert, tlskey); err != nil {
        log.Errorf("Failed to listen and serve: %v", err)
    }
}