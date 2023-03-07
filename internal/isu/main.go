package main

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/uroborosq/isu/internal/isu/adapters"
	"github.com/uroborosq/isu/internal/isu/app"
	"github.com/uroborosq/isu/internal/isu/ports"
	"golang.org/x/oauth2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

var (
	connString   = os.Getenv("DATABASE_URL")
	clientID     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	redirectURL  = os.Getenv("REDIRECT_URL")
	issuerURL    = os.Getenv("ISSUER_URL")
	port         = os.Getenv("PORT")
)

func main() {
	ctx := context.Background()
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	userRepo, err := adapters.NewUserRepo(db)
	if err != nil {
		log.Fatal(err.Error())
	}

	provider, err := oidc.NewProvider(ctx, issuerURL)
	if err != nil {
		log.Fatal(err.Error())
	}

	oauth2Config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID},
	}

	service := app.NewIsuService(userRepo)
	httpServer := ports.NewHttpServer(service, oauth2Config, ctx, provider)
	r := ports.Init(httpServer)
	log.Println("Server started on port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
