package rest

import (
	"../../models"
	"../../pkg/errors"
	"../../pkg/logger"
	"../../pkg/render"
	"context"
	"github.com/gorilla/mux"
	"net/http"
)

type Rest struct {
	logger.Logger
}

// New initializes the server with routes exposing the given usecases.
func New(logger logger.Logger, pc paymentCallbacks, cc cardCallbacks) http.Handler {
	// setup router with default handlers
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(methodNotAllowedHandler)

	// setup api endpoints
	addPaymentCallback(router, pc, logger)
	addCardCallback(router, cc, logger)

	return router
}

func notFoundHandler(wr http.ResponseWriter, req *http.Request) {
	_ = render.JSON(wr, http.StatusNotFound, errors.ResourceNotFound("path", req.URL.Path))
}

func methodNotAllowedHandler(wr http.ResponseWriter, req *http.Request) {
	_ = render.JSON(wr, http.StatusMethodNotAllowed, errors.ResourceNotFound("path", req.URL.Path))
}

type paymentCallbacks interface {
	InsertPayment(ctx context.Context, payment models.PaymentCallback) error
}

type cardCallbacks interface {
	InsertCard(ctx context.Context, card models.CardCallback) error
}
