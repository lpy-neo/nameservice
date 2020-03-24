package rest

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/gorilla/mux"
)

const (
	restName = "name"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/names", storeName), namesHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/names/{%s}/buy_name", storeName, restName), buyNameHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/names{%s}/set_name", storeName, restName), setNameHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/names/{%s}", storeName, restName), resolveNameHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/names/{%s}/whois", storeName, restName), whoIsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/names", storeName), deleteNameHandler(cliCtx)).Methods("DELETE")
	r.HandleFunc(fmt.Sprintf("/%s/names/{%s}/set_sale", storeName, restName), setSaleHandler(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/%s/names/{%s}/sale_status", storeName, restName), saleStausHandler(cliCtx, storeName)).Methods("GET")
}
