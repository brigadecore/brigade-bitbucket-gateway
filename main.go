package main

import (
	"log"
	"net/http"

	"github.com/brigadecore/brigade-bitbucket-gateway/internal/webhooks"
	libHTTP "github.com/brigadecore/brigade-foundations/http"
	"github.com/brigadecore/brigade-foundations/signals"
	"github.com/brigadecore/brigade-foundations/version"
	"github.com/brigadecore/brigade/sdk/v3"
	"github.com/gorilla/mux"
)

func main() {

	log.Printf(
		"Starting Brigade Bitbucket Gateway -- version %s -- commit %s",
		version.Version(),
		version.Commit(),
	)

	var webhooksService webhooks.Service
	{
		address, token, opts, err := apiClientConfig()
		if err != nil {
			log.Fatal(err)
		}
		webhooksService = webhooks.NewService(
			sdk.NewEventsClient(address, token, &opts),
			webhookServiceConfig(),
		)
	}

	var ipFilter libHTTP.Filter
	{
		config, err := ipFilterConfig()
		if err != nil {
			log.Fatal(err)
		}
		ipFilter = libHTTP.NewIPFilter(config)
	}

	webhooksHandler, err := webhooks.NewHandler(webhooksService)
	if err != nil {
		log.Fatal(err)
	}

	var server libHTTP.Server
	{
		router := mux.NewRouter()
		router.StrictSlash(true)
		router.Handle(
			"/events",
			ipFilter.Decorate(webhooksHandler.ServeHTTP),
		).Methods(http.MethodPost)
		router.HandleFunc("/healthz", libHTTP.Healthz).Methods(http.MethodGet)
		serverConfig, err := serverConfig()
		if err != nil {
			log.Fatal(err)
		}
		server = libHTTP.NewServer(router, &serverConfig)
	}

	log.Println(
		server.ListenAndServe(signals.Context()),
	)
}
