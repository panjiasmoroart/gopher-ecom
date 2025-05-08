package main

import (
	"log"

	"github.com/panjiasmoroart/gopher-ecom/cmd/api"
)

func main() {
	server := api.NewAPIServer(":9090", nil)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
