package netatmo2mackerel

import (
	"fmt"
	"net/http"
	"os"
	"time"

	netatmo "github.com/aereal/netatmo-api-go"
	mkr "github.com/mackerelio/mackerel-client-go"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
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

var metricPrefixes = map[string]string{
	"BatteryPercent": "battery",
	"RFStatus":       "rf_status",
	"WifiStatus":     "wifi",
	"Temperature":    "temperature",
	"Humidity":       "humidity",
	"CO2":            "co2",
	"Noise":          "noise",
	"Pressure":       "pressure",
}

func init() {
	http.HandleFunc("/metrics", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	metricsValues, err := fetchWeatherStationMetrics(ctx, time.Duration(1)*time.Minute)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Failed: (%s) %#v", err, err)
	}
	fmt.Fprintf(w, "%#v", metricsValues)
}

func fetchWeatherStationMetrics(ctx context.Context, resolution time.Duration) ([]*mkr.MetricValue, error) {
	resolutionInt := int64(resolution / time.Second)
	netatmoClient, err := netatmo.NewClientWithContext(ctx, netatmo.Config{
		ClientID:     netatmoAppID,
		ClientSecret: netatmoAppSecret,
		Username:     netatmoEmail,
		Password:     netatmoPassword,
	})
	dc, err := netatmoClient.Read()
	if err != nil {
		return nil, err
	}
	metrics := make([]*mkr.MetricValue, 0)
	for _, station := range dc.Stations() {
		for _, module := range station.Modules() {
			infoEpoch, info := module.Info()
			if resolution != time.Duration(0) {
				infoEpoch = infoEpoch / resolutionInt * resolutionInt
			}
			infoTs := time.Unix(infoEpoch, 0)
			for name, value := range info {
				log.Debugf(ctx, "Timestamp:%s Station:%s Module:%s Name:%s Value:%#v\n", infoTs, module.StationName, module.ModuleName, name, value)
				if metricPrefix, ok := metricPrefixes[name]; ok {
					if v, ok := float64of(value); ok {
						metrics = append(metrics, &mkr.MetricValue{
							Name:  metricPrefix + "." + module.ModuleName,
							Value: v,
							Time:  infoTs.Unix(),
						})
					}
				}
			}

			dataEpoch, data := module.Data()
			if resolution != time.Duration(0) {
				dataEpoch = dataEpoch / resolutionInt * resolutionInt
			}
			dataTs := time.Unix(dataEpoch, 0)
			for name, value := range data {
				log.Debugf(ctx, "Timestamp:%s Station:%s Module:%s Name:%s Value:%#v\n", dataTs, module.StationName, module.ModuleName, name, value)
				if metricPrefix, ok := metricPrefixes[name]; ok {
					if v, ok := float64of(value); ok {
						metrics = append(metrics, &mkr.MetricValue{
							Name:  metricPrefix + "." + module.ModuleName,
							Value: v,
							Time:  dataTs.Unix(),
						})
					}
				}
			}
		}
	}
	return metrics, nil
}

func float64of(value interface{}) (float64, bool) {
	switch t := value.(type) {
	case float64:
		return t, true
	case *float64:
		return *t, true
	case float32:
		return float64(t), true
	case *float32:
		return float64(*t), true
	case int32:
		return float64(t), true
	case *int32:
		return float64(*t), true
	case int64:
		return float64(t), true
	case *int64:
		return float64(*t), true
	default:
		return 0, false
	}
}
