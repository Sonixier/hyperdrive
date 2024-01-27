package wallet

import (
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive/hyperdrive-daemon/api/server"
	"github.com/nodeset-org/hyperdrive/shared/types/api"
)

// ===============
// === Factory ===
// ===============

type walletExportContextFactory struct {
	handler *WalletHandler
}

func (f *walletExportContextFactory) Create(args url.Values) (*walletExportContext, error) {
	c := &walletExportContext{
		handler: f.handler,
	}
	return c, nil
}

func (f *walletExportContextFactory) RegisterRoute(router *mux.Router) {
	server.RegisterQuerylessGet[*walletExportContext, api.WalletExportData](
		router, "export", f, f.handler.serviceProvider,
	)
}

// ===============
// === Context ===
// ===============

type walletExportContext struct {
	handler *WalletHandler
}

func (c *walletExportContext) PrepareData(data *api.WalletExportData, opts *bind.TransactOpts) error {
	sp := c.handler.serviceProvider
	w := sp.GetWallet()

	// Requirements
	err := sp.RequireWalletReady()
	if err != nil {
		return err
	}

	// Get password
	pw, isSet := w.GetPassword()
	if !isSet {
		return fmt.Errorf("password has not been set; cannot decrypt wallet keystore without it")
	}
	data.Password = pw

	// Serialize wallet
	walletString, err := w.String()
	if err != nil {
		return fmt.Errorf("error serializing wallet keystore: %w", err)
	}
	data.Wallet = walletString

	// Get account private key
	data.AccountPrivateKey = w.GetNodePrivateKeyBytes()

	return nil
}
