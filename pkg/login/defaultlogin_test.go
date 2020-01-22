package login

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/matryer/is"

	"github.com/ZupIT/ritchie-cli/pkg/file/fileutil"
)

const (
	okResponse = `{
		"access_token": "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICIzVFZobDgyVEZzZUxWQkhmb0hJU3FyMldienFKMkpaWkt6ZmFnRDNWOGJFIn0.eyJqdGkiOiI0ZjNmMzUzMS05MzBjLTQ1NzYtOTgxNS05MjM0Mjc1ZWY1OGMiLCJleHAiOjE1NzY2OTc5MjEsIm5iZiI6MCwiaWF0IjoxNTc2Njk2NzIxLCJpc3MiOiJodHRwczovL2tleWNsb2FrLWhhLXNhYXMuYXBpcmVhbHdhdmUuaW8vYXV0aC9yZWFsbXMvVHlrIiwiYXVkIjpbInJlYWxtLW1hbmFnZW1lbnQiLCJ0eWstY2xpZW50IiwiYWNjb3VudCJdLCJzdWIiOiJkOWY2NGUzZi0zYTBkLTRiZDMtYWMxMC05MGZjMjQzMTM4ODgiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJ1c2VyLWxvZ2luIiwiYXV0aF90aW1lIjowLCJzZXNzaW9uX3N0YXRlIjoiYzRjNjZiMGUtYmI2NS00ZDk4LTk3ZWUtNjJjNTVmZDJmNmM2IiwiYWNyIjoiMSIsInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJtb292ZV9yZWFkIiwib2ZmbGluZV9hY2Nlc3MiLCJjb25maWdfcmVhZCIsImNvbmZpZ193cml0ZSIsInVtYV9hdXRob3JpemF0aW9uIiwibW9vdmVfd3JpdGUiXX0sInJlc291cmNlX2FjY2VzcyI6eyJyZWFsbS1tYW5hZ2VtZW50Ijp7InJvbGVzIjpbInZpZXctcmVhbG0iLCJ2aWV3LWlkZW50aXR5LXByb3ZpZGVycyIsIm1hbmFnZS1pZGVudGl0eS1wcm92aWRlcnMiLCJpbXBlcnNvbmF0aW9uIiwicmVhbG0tYWRtaW4iLCJjcmVhdGUtY2xpZW50IiwibWFuYWdlLXVzZXJzIiwicXVlcnktcmVhbG1zIiwidmlldy1hdXRob3JpemF0aW9uIiwicXVlcnktY2xpZW50cyIsInF1ZXJ5LXVzZXJzIiwibWFuYWdlLWV2ZW50cyIsIm1hbmFnZS1yZWFsbSIsInZpZXctZXZlbnRzIiwidmlldy11c2VycyIsInZpZXctY2xpZW50cyIsIm1hbmFnZS1hdXRob3JpemF0aW9uIiwibWFuYWdlLWNsaWVudHMiLCJxdWVyeS1ncm91cHMiXX0sInR5ay1jbGllbnQiOnsicm9sZXMiOlsiYWRtaW4iXX0sImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoiZW1haWwgcHJvZmlsZSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiTWFyY29zIEd1aW1hcmFlcyIsInByZWZlcnJlZF91c2VybmFtZSI6Im1hcmNvcy5tZWRlaXJvc0B6dXAuY29tLmJyIiwiZ2l2ZW5fbmFtZSI6Ik1hcmNvcyIsImZhbWlseV9uYW1lIjoiR3VpbWFyYWVzIiwiZW1haWwiOiJtYXJjb3MubWVkZWlyb3NAenVwLmNvbS5iciJ9.qBC6-AKZ6xcuwf-OlOqSLXhFVcPK6s3TB3wbjPNLLBIjjNH09Qvc-SR1tPdqgJWb_jHrci2Lp9cjuqgfqEaELne-cNTpPc3ttbtcOZTujhCmmruKCqi9YJtCpBtOP67SbyvNuy76gEXyjd0xD3A32cldk-ZZPFIjz8hm81zlamYE5-FNM-UVIzqrIf0_ULT4jfZnh6sAfRl-yEpr2HehH6AT9nUHB3O6rx5EmlqXY8Z_Xhqctj5uGUESu1NVTBtm6qQoa6tKoQshRaSfFKQsHk4u4YiBIhhNN1zd9ucCnWMrbhkaFB3nmAcxvXmRYKBOMOkip866VUelAZuRHnbCJg"
	}`
	pwOk            = "ok"
	pwBadCredential = "badcredential"
	pwError         = "error"
	pwUnavailable   = "unavailable"
)

var (
	homePath string
)

func TestMain(m *testing.M) {
	homePath = os.TempDir()
	os.Exit(m.Run())
}

func TestAuthenticate(t *testing.T) {
	fileutil.RemoveFile(homePath + "/.session")

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cred := &Credential{}
		json.NewDecoder(r.Body).Decode(cred)
		switch cred.Password {
		case pwBadCredential:
			w.WriteHeader(http.StatusUnauthorized)
		case pwError:
			w.WriteHeader(http.StatusInternalServerError)
		case pwUnavailable:
			w.WriteHeader(http.StatusServiceUnavailable)
		default:
			w.Write([]byte(okResponse))
		}
	})

	logman, teardown := newTestingManager(h)
	defer teardown()

	is := is.New(t)

	tests := []struct {
		in  *Credential
		out error
	}{
		{&Credential{Username: "marie.curie", Password: pwOk, Organization: "zup"}, nil},
		{&Credential{Username: "marie.curie", Password: pwBadCredential, Organization: "zup"}, ErrBadCredential},
		{&Credential{Username: "marie.curie", Password: pwError, Organization: "zup"}, ErrUnknown},
		{&Credential{Username: "marie.curie", Password: pwUnavailable, Organization: "zup"}, ErrServiceUnavailable},
	}

	for _, test := range tests {
		t.Run(test.in.Password, func(t *testing.T) {
			err := logman.Authenticate(test.in)
			is.Equal(test.out, err)
		})
	}
}

func newTestingManager(handler http.Handler) (Manager, func()) {
	s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	logman := NewDefaultManager(homePath, s.URL, cli)

	return logman, s.Close
}
