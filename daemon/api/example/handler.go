package example

import (
	"github.com/gorilla/mux"

)

type AuctionHandler struct {
	serviceProvider *services.ServiceProvider
	factories       []server.IContextFactory
}

func NewAuctionHandler(serviceProvider *services.ServiceProvider) *AuctionHandler {
	h := &AuctionHandler{
		serviceProvider: serviceProvider,
	}
	h.factories = []server.IContextFactory{
		&auctionBidContextFactory{h},
		&auctionClaimContextFactory{h},
		&auctionCreateContextFactory{h},
		&auctionLotContextFactory{h},
		&auctionRecoverContextFactory{h},
		&auctionStatusContextFactory{h},
	}
	return h
}

func (h *AuctionHandler) RegisterRoutes(router *mux.Router) {
	subrouter := router.PathPrefix("/auction").Subrouter()
	for _, factory := range h.factories {
		factory.RegisterRoute(subrouter)
	}
}
