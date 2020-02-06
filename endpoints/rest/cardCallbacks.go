package rest

import (
	"../../models"
	"../../pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
)

func addCardCallback(router *mux.Router, cCall cardCallbacks, lg logger.Logger) {
	pc := &cardCallbacksController{}
	pc.cc = cCall
	pc.Logger = lg

	router.HandleFunc("/pay/ckassa", pc.post).Methods(http.MethodPost)
}

type cardCallbacksController struct {
	logger.Logger

	cc cardCallbacks
}

func (con cardCallbacksController) post(wr http.ResponseWriter, req *http.Request) {
	card := models.CardCallback{}
	if err := readRequest(req, &card); err != nil {
		con.Warnf("failed to read user request: %s", err)
		respond(wr, http.StatusBadRequest, err)
		return
	}

	err := con.cc.InsertCard(req.Context(), card)
	if err != nil {
		respondErr(wr, err)
		return
	}

	respond(wr, http.StatusCreated, "")
}
