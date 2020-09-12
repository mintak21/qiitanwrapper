package handler

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/mintak21/qiitaWrapper/gen/restapi/qiitawrapper/health"
)

// NewHealthHandler creates New HealthCheck Handler
func NewHealthHandler() health.HealthHandler {
	return &healthHandler{}
}

type healthHandler struct{}

func (h *healthHandler) Handle(params health.HealthParams) middleware.Responder {
	dummy := struct{}{}
	return health.NewHealthOK().WithPayload(dummy)
}
