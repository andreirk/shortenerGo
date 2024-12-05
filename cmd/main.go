package main

import (
	"fmt"
	"go/adv-demo/config"
	"go/adv-demo/internal/auth"
	"go/adv-demo/internal/hello"
	"go/adv-demo/internal/link"
	"go/adv-demo/internal/stat"
	"go/adv-demo/internal/user"
	"go/adv-demo/pkg/db"
	"go/adv-demo/pkg/midleware"
	"log"
	"net/http"
)

func main() {
	conf := config.LoadConfig()
	router := http.NewServeMux()
	db := db.NewDb(conf)

	//Repositories
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepository := stat.NewStatRepository(db)

	//Services
	authService := auth.NewAuthService(userRepository)

	//Handlers
	hello.NewHelloHandler(router)
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	link.NewLinkHandler(router, link.LinkDeps{
		LinkRepository: linkRepository,
		StatRepository: statRepository,
		Config:         conf,
	})

	//Midlewares
	midlewareStack := midleware.Chain(
		midleware.CORS,
		midleware.Logging,
	)

	server := http.Server{
		Addr:    "localhost:" + conf.Port,
		Handler: midlewareStack(router),
	}
	fmt.Println("Server is listening on port:", conf.Port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
