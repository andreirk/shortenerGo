package main

import (
	"fmt"
	"go/adv-demo/config"
	"go/adv-demo/internal/auth"
	"go/adv-demo/internal/hello"
	"go/adv-demo/internal/link"
	"net/http"
)

func main() {
	conf := config.LoadConfig()
	router := http.NewServeMux()

	//Handlers
	hello.NewHelloHandler(router)
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
	})
	link.NewLinkHandler(router, link.LinkDeps{})

	server := http.Server{
		Addr:    "localhost:" + conf.Port,
		Handler: router,
	}
	fmt.Println("Server is listening on port:", conf.Port)
	server.ListenAndServe()
}
