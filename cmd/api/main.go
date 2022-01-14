package main

import (
	"flag"
	"os"

	loads "github.com/go-openapi/loads"
	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"

	"github.com/mintak21/qiitaWrapper/api/handler"
	"github.com/mintak21/qiitaWrapper/gen/restapi"
	qws "github.com/mintak21/qiitaWrapper/gen/restapi/qiitawrapper"
)

const (
	defaultPort = 8090
)

var port int

func main() {
	// not use init() cause linter checks gochecknoinits(golanglintci)
	initialSetup()

	// load embedded swagger file
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "2.0")
	if err != nil {
		log.Fatalln(err)
	}

	// create new service API
	api := qws.NewQiitawrapperAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer shutdown(server)

	// parse flags
	flag.IntVar(&port, "port", defaultPort, "Port to run this service on")
	flag.Parse()
	// set the port this service will be run on
	server.Port = port

	// Set Handler
	api.HealthHealthHandler = handler.NewHealthHandler()
	api.ItemsGetTagItemsHandler = handler.NewGetTagItemsHandler()
	api.ItemsSyncTagItemsHandler = handler.NewSyncTagItemsHandler()
	api.ItemsGetMonthlyTrendItemsHandler = handler.NewMonthlyTrendItemsHandler()

	// configure server
	server.ConfigureAPI()

	// serve API
	if err := server.Serve(); err != nil {
		log.Error(err)
	}
}

func initialSetup() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	strfmt.MarshalFormat = strfmt.RFC3339Millis
}

func shutdown(server *restapi.Server) {
	err := server.Shutdown()
	if err != nil {
		os.Exit(1)
	}
}
