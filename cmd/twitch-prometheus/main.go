package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/nicklaw5/helix"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var cfgfile = flag.String("config", "config.yml", "Config file location")
	flag.Parse()

	config, err := ReadConfig(*cfgfile)
	if err != nil {
		panic(err)
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
	})
	if err != nil {
		panic(err)
	}

	token, err := client.RequestAppAccessToken([]string{})
	if err != nil {
		panic(err)
	}

	client.SetAppAccessToken(token.Data.AccessToken)

	collector, err := NewFollowersCollector(client, config.Followers)
	if err != nil {
		panic(err)
	}
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.ListenPort), nil))
}
