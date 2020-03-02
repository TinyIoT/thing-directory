// Copyright 2014-2016 Fraunhofer Institute for Applied Information Technology FIT

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/justinas/alice"
	_ "github.com/linksmart/go-sec/auth/keycloak/obtainer"
	_ "github.com/linksmart/go-sec/auth/keycloak/validator"
	"github.com/linksmart/go-sec/auth/validator"
	"github.com/linksmart/resource-catalog/catalog"
	"github.com/oleksandr/bonjour"
	uuid "github.com/satori/go.uuid"
)

var (
	confPath = flag.String("conf", "conf/resource-catalog.json", "Resource catalog configuration file path")
)

func main() {
	defer log.Println("Stopped.")
	flag.Parse()

	config, err := loadConfig(*confPath)
	if err != nil {
		panic("Error reading config file:" + err.Error())
	}
	if config.ServiceID == "" {
		config.ServiceID = uuid.NewV4().String()
		log.Printf("Service ID not set. Generated new UUID: %s", config.ServiceID)
	}

	// Setup API storage
	var storage catalog.Storage
	switch config.Storage.Type {
	case catalog.BackendLevelDB:
		storage, err = catalog.NewLevelDBStorage(config.Storage.DSN, nil)
		if err != nil {
			panic("Failed to start LevelDB storage:" + err.Error())
		}
		defer storage.Close()
	default:
		panic("Could not create catalog API storage. Unsupported type:" + config.Storage.Type)
	}

	controller, err := catalog.NewController(storage)
	if err != nil {
		panic("Failed to start the controller:" + err.Error())
	}
	defer controller.Stop()

	// Create catalog API object
	api := catalog.NewHTTPAPI(
		controller,
		config.ServiceID,
	)

	nRouter, err := setupHTTPRouter(config, api)
	if err != nil {
		panic(err)
	}
	// Start listener
	addr := fmt.Sprintf("%s:%d", config.BindAddr, config.BindPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Printf("HTTP server listening on %v", addr)
	go func() { log.Fatalln(http.Serve(listener, nRouter)) }()

	// Announce service using DNS-SD
	var bonjourS *bonjour.Server
	if config.DnssdEnabled {
		go func() {
			bonjourS, err = bonjour.Register(config.Description,
				catalog.DNSSDServiceType,
				"",
				config.BindPort,
				[]string{"uri=/td"},
				nil)
			if err != nil {
				log.Printf("Failed to register DNS-SD service: %s", err.Error())
				return
			}
			log.Println("Registered service via DNS-SD using type", catalog.DNSSDServiceType)
		}()
	}

	// Register in the LinkSmart Service Catalog
	if config.ServiceCatalog != nil {
		unregisterService, err := registerInServiceCatalog(config)
		if err != nil {
			panic("Error registering service:" + err.Error())
		}
		// Unregister from the Service Catalog
		defer unregisterService()
	}

	log.Println("Ready!")

	// Ctrl+C / Kill handling
	handler := make(chan os.Signal, 1)
	signal.Notify(handler, os.Interrupt, os.Kill)
	<-handler
	log.Println("Shutting down...")

	// Stop bonjour registration
	if bonjourS != nil {
		bonjourS.Shutdown()
		time.Sleep(1e9)
	}

}

func setupHTTPRouter(config *Config, api *catalog.HTTPAPI) (*negroni.Negroni, error) {

	commonHandlers := alice.New(
		context.ClearHandler,
	)

	// Append auth handler if enabled
	if config.Auth.Enabled {
		// Setup ticket validator
		v, err := validator.Setup(
			config.Auth.Provider,
			config.Auth.ProviderURL,
			config.Auth.ServiceID,
			config.Auth.BasicEnabled,
			config.Auth.Authz)
		if err != nil {
			return nil, err
		}

		commonHandlers = commonHandlers.Append(v.Handler)
	}

	// Configure http api router
	r := newRouter()
	r.post("/td/", commonHandlers.ThenFunc(api.Post))
	r.get("/td/{id:.+}", commonHandlers.ThenFunc(api.Get))
	r.put("/td/{id:.+}", commonHandlers.ThenFunc(api.Put))
	r.delete("/td/{id:.+}", commonHandlers.ThenFunc(api.Delete))
	r.get("/td", commonHandlers.ThenFunc(api.List))
	r.get("/td/filter/{path}/{op}/{value:.*}", commonHandlers.ThenFunc(api.Filter))

	logger := negroni.NewLogger()
	logFlags := log.LstdFlags
	if evalEnv(EnvDisableLogTime) {
		logFlags = 0
	}
	logger.SetFlags(logFlags)
	logger.SetPrefix("")

	// Configure the middleware
	n := negroni.New(
		negroni.NewRecovery(),
		logger,
	)
	// Mount router
	n.UseHandler(r)

	return n, nil
}
