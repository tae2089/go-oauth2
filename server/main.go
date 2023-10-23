package main

import (
	"log"
	"net/http"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt"
)

func main() {
	manager := manage.NewDefaultManager()
	// token memory store

	manager.MustTokenStorage(store.NewMemoryTokenStore())
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))

	manager.SetAuthorizeCodeTokenCfg(&manage.Config{
		IsGenerateRefresh: true,
	})

	manager.SetClientTokenCfg(&manage.Config{
		IsGenerateRefresh: true,
	})

	// client memory store
	clientStore := store.NewClientStore()
	clientStore.Set("000000", &models.Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "http://localhost",
		UserID: "123123",
	})
	manager.MapClientStorage(clientStore)

	// srv := server.NewDefaultServer(manager)
	srv := server.NewServer(&server.Config{
		TokenType:             "JWT",
		AllowedGrantTypes:     []oauth2.GrantType{oauth2.AuthorizationCode, oauth2.ClientCredentials},
		AllowGetAccessRequest: true,
		AllowedResponseTypes:  []oauth2.ResponseType{oauth2.Token},
		// AllowedCodeChallengeMethods: []oauth2.CodeChallengeMethod{
		// 	oauth2.CodeChallengePlain,
		// 	oauth2.CodeChallengeS256,
		// },
	}, manager)

	srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	http.HandleFunc("/authorize", func(w http.ResponseWriter, r *http.Request) {
		err := srv.HandleAuthorizeRequest(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		srv.HandleTokenRequest(w, r)

	})

	log.Fatal(http.ListenAndServe(":9096", nil))
}
