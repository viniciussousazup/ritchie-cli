package login

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ZupIT/ritchie-cli/pkg/validator"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
	"github.com/ZupIT/ritchie-cli/pkg/env"
	"github.com/ZupIT/ritchie-cli/pkg/session"
)

const (
	callbackUrl = "http://localhost:8888/ritchie/callback"
	providerUrl = "%s/oauth"
	htmlClose   = `<!DOCTYPE html>

<html>
  <head>
    <meta charset="utf-8"/>
    <title>Login Ritchie</title>
    <style>
      html, body {
        height: 100%;
        justify-content: center;
        align-items: center; 
      }

      .container {
      	display: flex;
        justify-content: center;
        color: '#111';
        font-size: 30px;
      }
    </style>
  </head>

  <body>
    <div class="container">
      Login Successful
    </div>
    <div class="container">
      <span id="counter">5</span>s to close browser.
    </div>
    <div class="container">
      If not close return to CLI.
    </div>
  </body>

  <script type="text/javascript"> 
    (function startSetInterval() {
      let count = 5;

      const interval = setInterval(function t() {
        const counter = document.getElementById('counter')
        counter.innerText = count;

        if (count === 0) {
          clearInterval(interval)
          window.close()
        }

        count = count - 1;
        return t;
      }(), 1000); 
    }())
</script>

</html>`
)

type defaultManager struct {
	homePath   string
	serverURL  string
	httpClient *http.Client
	session    session.Manager
}

// NewDefaultManager creates a default instance of Manager interface
func NewDefaultManager(homePath, serverURL string, c *http.Client, s session.Manager) *defaultManager {
	return &defaultManager{homePath, serverURL, c, s}
}

func (d *defaultManager) Authenticate(organization,version string) error {
	providerConfig, err := providerConfig(organization)
	if err != nil {
		return err
	}
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, providerConfig.Url)
	if err != nil {
		return err
	}
	oauth2Config := oauth2.Config{
		ClientID:     providerConfig.ClientId,
		ClientSecret: "",
		RedirectURL:  callbackUrl,
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
	http.HandleFunc("/ritchie/callback", d.handler(provider, state, organization, version, oauth2Config, ctx))
	log.Fatal(http.ListenAndServe("localhost:8888", nil))

	return nil
}

func (d *defaultManager) handler(provider *oidc.Provider, state, organization,version string, oauth2Config oauth2.Config, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			Email    string `json:"email"`
			Username string `json:"preferred_username"`
		}{}
		idToken.Claims(&user)
		err = d.session.Create(token, user.Username, organization)
		if err != nil {
			http.Error(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
			go stopServer()
		}
		w.Write([]byte(htmlClose))
		log.Printf("Login ok!")
		validator.IsValidVersion(version, organization)
		go stopServer()
	}
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

func providerConfig(organization string) (ProviderConfig, error) {
	var provideConfig ProviderConfig
	url := fmt.Sprintf(providerUrl, env.ServerUrl)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return provideConfig, fmt.Errorf("Failed to providerConfig for org %s. \n%v", organization, err)
	}
	req.Header.Set("x-org", organization)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return provideConfig, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return provideConfig, fmt.Errorf("Failed to call url. %v for org %s. Status code: %d\n", url, organization, resp.StatusCode)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return provideConfig, fmt.Errorf("Failed parse response to body: %s\n", string(bodyBytes))
	}
	json.Unmarshal(bodyBytes, &provideConfig)
	return provideConfig, nil
}
