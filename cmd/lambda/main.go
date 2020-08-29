package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/strfmt"
	log "github.com/sirupsen/logrus"

	"github.com/mintak21/qiitaWrapper/api/handler"
	"github.com/mintak21/qiitaWrapper/gen/restapi"
	qws "github.com/mintak21/qiitaWrapper/gen/restapi/qiitawrapper"
)

var httpAdapter *httpadapter.HandlerAdapter

// handler handles API requests
func handleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return httpAdapter.Proxy(req)
}

func main() {
	// not use init() cause linter checks gochecknoinits
	initialSetup()
	lambda.Start(handleRequest)
}

func initialSetup() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	strfmt.MarshalFormat = strfmt.RFC3339Millis

	// load embedded swagger file
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "2.0")
	if err != nil {
		log.Fatalln(err)
	}

	// create new service API
	api := qws.NewQiitawrapperAPI(swaggerSpec)
	server := restapi.NewServer(api)

	// Set Handler
	api.ItemsGetTagItemsHandler = handler.NewGetTagItemsHandler()
	api.ItemsSyncTagItemsHandler = handler.NewSyncTagItemsHandler()
	api.ItemsGetMonthlyTrendItemsHandler = handler.NewMonthlyTrendItemsHandler()
	server.ConfigureAPI()

	// see https://github.com/go-swagger/go-swagger/issues/962#issuecomment-478382896
	httpAdapter = httpadapter.New(server.GetHandler())
}
