package netatmo2mackerel

import (
	"fmt"
	"net/http"
	"os"

	netatmo "github.com/aereal/gae-go-netatmo-ws-mackerel/netatmo"

	// netatmo "github.com/exzz/netatmo-api-go"
	mkr "github.com/mackerelio/mackerel-client-go"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/appengine"
)

const (
	baseURL = "https://api.netatmo.net/"
	authURL = baseURL + "oauth2/token"
)

var (
	netatmoEmail     = os.Getenv("NETATMO_EMAIL")
	netatmoPassword  = os.Getenv("NETATMO_PASSWORD")
	netatmoAppID     = os.Getenv("NETATMO_APP_ID")
	netatmoAppSecret = os.Getenv("NETATMO_APP_SECRET")
)

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/env", handleDumpEnv)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	metricsValues, err := fetchWeatherStationMetrics(ctx)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Failed: (%s) %#v", err, err)
	}
	fmt.Fprintf(w, "%#v", metricsValues)
}

func handleDumpEnv(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "email:%s; app_id:%s", netatmoEmail, netatmoAppID)
}

func fetchWeatherStationMetrics(ctx context.Context) ([]*mkr.MetricValue, error) {
	netatmoClient, err := netatmo.NewClient(netatmo.Config{
		ClientID:     netatmoAppID,
		ClientSecret: netatmoAppSecret,
		Username:     netatmoEmail,
		Password:     netatmoPassword,
	})
	fmt.Printf("client := %#v\n", netatmoClient)
	return nil, err
}

func newNetatmoClient(ctx context.Context) (*netatmo.Client, error) {
	oauthConfig := &oauth2.Config{
		ClientID:     netatmoAppID,
		ClientSecret: netatmoAppSecret,
		Scopes:       []string{"read_station"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  baseURL,
			TokenURL: authURL,
		},
	}
	token, err := oauthConfig.PasswordCredentialsToken(ctx, netatmoEmail, netatmoPassword)
	if err != nil {
		return nil, err
	}

	client := &netatmo.Client{
		oauth:      oauthConfig,
		httpClient: oauthConfig.Client(ctx, token),
		Dc:         &netatmo.DeviceCollection{},
	}
	return client, nil
}
