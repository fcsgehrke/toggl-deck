package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v8"
	"github.com/fcsgehrke/toggl-deck/api"
	"github.com/fcsgehrke/toggl-deck/internal/repositories/deck/postgres"
	"github.com/fcsgehrke/toggl-deck/internal/services"
)

//	@title			    Toggl Deck API
//	@version		    1.0
//	@description	  This is a sample deck api for testing purposes.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	  Felipe C. Gehrke
//	@contact.url	  https://github.com/fcsgehrke/toggl-deck
//	@contact.email	fcgehrke@outlook.com

//	@license.name	  Apache 2.0
//	@license.url	  http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		        localhost:3000
//	@BasePath	      /api/v1/decks

const (
	Name string = `
████████╗░█████╗░░██████╗░░██████╗░██╗░░░░░░░░░░░██████╗░███████╗░█████╗░██╗░░██╗
╚══██╔══╝██╔══██╗██╔════╝░██╔════╝░██║░░░░░░░░░░░██╔══██╗██╔════╝██╔══██╗██║░██╔╝
░░░██║░░░██║░░██║██║░░██╗░██║░░██╗░██║░░░░░█████╗██║░░██║█████╗░░██║░░╚═╝█████═╝░
░░░██║░░░██║░░██║██║░░╚██╗██║░░╚██╗██║░░░░░╚════╝██║░░██║██╔══╝░░██║░░██╗██╔═██╗░
░░░██║░░░╚█████╔╝╚██████╔╝╚██████╔╝███████╗░░░░░░██████╔╝███████╗╚█████╔╝██║░╚██╗
░░░╚═╝░░░░╚════╝░░╚═════╝░░╚═════╝░╚══════╝░░░░░░╚═════╝░╚══════╝░╚════╝░╚═╝░░╚═╝`
)

type config struct {
	PostgresDSN string `env:"POSTGRES_DSN"`
	ListenAddr  string `env:"LISTEN_ADDR"`
}

func init() {
	fmt.Println(Name)
  fmt.Println()
}

func main() {
	// Logging and Context
	ctx := context.Background()
	log := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)

	// Handle Environment
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("[ERR] - Couldn't parse env. Err: %s\n", err.Error())
	}

	// Db Conn
	repo, err := postgres.NewDeckPostgresRepo(ctx, cfg.PostgresDSN, log)
	if err != nil {
		log.Printf("[ERR] - Repository creation failed w/ err: %s\n", err.Error())
		return
	}

	if err := repo.Migrate(); err != nil {
		log.Printf("[ERR] - Repository migration failed w/ err: %s\n", err.Error())
		return
	}

	// Service Initialization
	service, err := services.NewDeckService(repo, log)
	if err != nil {
		log.Printf("[ERR] - Deck Service creation failed w/ err: %s\n", err.Error())
		return
	}

	// HTTP Routes
	if err := api.Start(ctx, cfg.ListenAddr, service, log); err != nil {
		log.Printf("[ERR] - App Start failed w/ err: %s\n", err.Error())
	}
}
