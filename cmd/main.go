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
	"go/adv-demo/pkg/event"
	"go/adv-demo/pkg/midleware"
	"log"
	"net/http"
)

func App(conf *config.Config) http.Handler {
	router := http.NewServeMux()
	dbInstance := db.NewDb(conf)
	eventBus := event.NewEventBus()

	//Repositories
	linkRepository := link.NewLinkRepository(dbInstance)
	userRepository := user.NewUserRepository(dbInstance)
	statRepository := stat.NewStatRepository(dbInstance)

	//Services
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepository,
	})

	//Handlers
	hello.NewHelloHandler(router)
	auth.NewHandler(router, auth.HandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepository: statRepository,
		Config:         conf,
	})

	link.NewLinkHandler(router, link.LinkDeps{
		LinkRepository: linkRepository,
		EventBus:       eventBus,
		Config:         conf,
	})

	go statService.AddClick()

	//Midlewares
	midlewareStack := midleware.Chain(
		midleware.CORS,
		midleware.Logging,
	)

	return midlewareStack(router)
}

func main() {
	conf := config.LoadConfig("local")
	app := App(conf)
	server := http.Server{
		Addr:    "localhost:" + conf.Port,
		Handler: app,
	}

	fmt.Println("Server is listening on:", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
