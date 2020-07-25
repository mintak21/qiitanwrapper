package main

import (
	"flag"
	loads "github.com/go-openapi/loads"
	"github.com/mintak21/qiitaWrapper/gen/restapi"
	qws "github.com/mintak21/qiitaWrapper/gen/restapi/qiitawrapper"
	"log"
)

var port int

func main() {
	// load embedded swagger file
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "2.0")
	if err != nil {
		log.Fatalln(err)
	}

	// create new service API
	api := qws.NewQiitawrapperAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	// parse flags
	flag.IntVar(&port, "port", 8090, "Port to run this service on")
	flag.Parse()
	// set the port this service will be run on
	server.Port = port

	// TODO: Set Handle

	// serve API
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
