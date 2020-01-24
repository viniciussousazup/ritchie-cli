package login

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/crypto/cryptoutil"
	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
	oidc "github.com/coreos/go-oidc"
	"github.com/denisbrodbeck/machineid"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const (
	urlPattern         = "%s/login"
	sessionFilePattern = "%s/.session"

	// AES passphrase
	passphrase = "zYtBIK67fCmhrU0iUbPQ1Cf9"
)

type defaultManager struct {
	homePath   string
	serverURL  string
	httpClient *http.Client
}

// NewDefaultManager creates a default instance of Manager interface
func NewDefaultManager(homePath, serverURL string, httpClient *http.Client) *defaultManager {
	return &defaultManager{homePath, serverURL, httpClient}
}

func (d *defaultManager) Authenticate(organization string) error {
	providerConfig := getProviderConfig(organization)
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, providerConfig.ConfigUrl)
	if err != nil {
		return err
	}
	oauth2Config := oauth2.Config{
		ClientID:     providerConfig.ClientId,
		ClientSecret: providerConfig.ClientSecret,
		RedirectURL:  CallbackUrl,
		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email"},
	}
	state := "somestate"
	err = openBrowser(oauth2Config.AuthCodeURL(state))
	if err != nil {
		return err
	}
	http.HandleFunc("/ritchie/callback", d.handler(provider, state, organization, oauth2Config, ctx))
	log.Fatal(http.ListenAndServe("localhost:8888", nil))

	return nil
}

func (d *defaultManager) handler(provider *oidc.Provider, state, organization string, oauth2Config oauth2.Config, ctx context.Context) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oidcConfig := &oidc.Config{
			ClientID: oauth2Config.ClientID,
		}
		verifier := provider.Verifier(oidcConfig)
		if r.URL.Query().Get("state") != state {
			http.Error(w, "state did not match", http.StatusBadRequest)
			go stopServer()
		}

		oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			go stopServer()
		}
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
			go stopServer()
		}
		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
			go stopServer()
		}
		token := oauth2Token.AccessToken
		user := struct {
			Email string `json:"email"`
			Username string `json:"preferred_username"`
		}{}
		idToken.Claims(&user)
		err = d.createSession(token, user.Username, organization)
		if err != nil {
			http.Error(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
			go stopServer()
		}
		w.Write([]byte("Login ok, return to Rit CLI!"))
		log.Printf("Login ok!")
		go stopServer()
	})
}

func (d *defaultManager) createSession(token, username, organization string) error {
	session := &Session{
		AccessToken:  	token,
		Organization: 	organization,
		Username:    	username,
	}

	b, err := json.Marshal(session)
	if err != nil {
		return err
	}
	id, err := machineid.ID()
	if err != nil {
		return err
	}
	h := md5.New()
	io.WriteString(h, passphrase)
	io.WriteString(h, id)
	cipher := cryptoutil.Encrypt(string(h.Sum(nil)), string(b))
	sessFilePath := fmt.Sprintf(sessionFilePattern, d.homePath)
	err = fileutil.WriteFile(sessFilePath, []byte(cipher))
	if err != nil {
		return err
	}
	err = os.Chmod(sessFilePath, 0600)
	if err != nil {
		return err
	}
	return nil
}

func stopServer() {
	time.Sleep(5 * time.Second)
	os.Exit(0)
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = nil
	}
	return err
}

func getProviderConfig(organization string) ProviderConfig {
	//TODO: put get OAuth config in server
	return ProviderConfig{
		ConfigUrl:    "https://ritchie-keycloak.itiaws.dev/auth/realms/ritchie",
		ClientId:     "user-login",
		ClientSecret: "71da7105-477b-4854-ae6c-1263abf7ea50",
	}

}

func (d *defaultManager) Session() (*Session, error) {
	sessFilePath := fmt.Sprintf(sessionFilePattern, d.homePath)
	if !fileutil.Exists(sessFilePath) {
		return nil, errors.New("Please, you need to login first")
	}
	b, err := fileutil.ReadFile(sessFilePath)
	if err != nil {
		return nil, err
	}
	id, err := machineid.ID()
	if err != nil {
		return nil, err
	}
	h := md5.New()
	io.WriteString(h, passphrase)
	io.WriteString(h, id)
	plain := cryptoutil.Decrypt(string(h.Sum(nil)), string(b))
	session := &Session{}
	if err := json.Unmarshal([]byte(plain), session); err != nil {
		return nil, err
	}
	return session, nil
}
