package rest

import (
	"ckassa-callback/models"
	"ckassa-callback/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
)

func addPaymentCallback(router *mux.Router, payCall paymentCallbacks, lg logger.Logger) {
	pc := &paymentCallbacksController{}
	pc.pc = payCall
	pc.Logger = lg

	router.HandleFunc("/pay/ckassa/card-reg", pc.post).Methods(http.MethodPost)
}

type paymentCallbacksController struct {
	logger.Logger

	pc paymentCallbacks
}

func (cc paymentCallbacksController) post(wr http.ResponseWriter, req *http.Request) {
	payment := models.PaymentCallback{}
	if err := readRequest(req, &payment); err != nil {
		cc.Warnf("failed to read user request: %s", err)
		respond(wr, http.StatusBadRequest, err)
		return
	}

	err := cc.pc.InsertPayment(req.Context(), payment)
	if err != nil {
		respondErr(wr, err)
		return
	}

	respond(wr, http.StatusCreated, "")
}
