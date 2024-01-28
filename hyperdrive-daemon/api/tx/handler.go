package tx

import (
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/api/server"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/common/services"
)

type TxHandler struct {
	serviceProvider *services.ServiceProvider
	factories       []server.IContextFactory
}

func NewTxHandler(serviceProvider *services.ServiceProvider) *TxHandler {
	h := &TxHandler{
		serviceProvider: serviceProvider,
	}
	h.factories = []server.IContextFactory{
		&txBatchSignTxsContextFactory{h},
		&txBatchSubmitTxsContextFactory{h},
		&txSignTxContextFactory{h},
		&txSubmitTxContextFactory{h},
		&txWaitContextFactory{h},
	}
	return h
}

func (h *TxHandler) RegisterRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/tx").Subrouter()
	for _, factory := range h.factories {
		factory.RegisterRoute(subrouter)
	}
}
