package nameservice

import (
	"github.com/lpy-neo/nameservice/x/nameservice/internal/keeper"
	"github.com/lpy-neo/nameservice/x/nameservice/internal/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewKeeper         = keeper.NewKeeper
	NewQuerier        = keeper.NewQuerier
	NewMsgBuyName     = types.NewMsgBuyName
	NewMsgSetName     = types.NewMsgSetName
	NewMsgDeleteName  = types.NewMsgDeleteName
	NewWhois          = types.NewWhois
	NewMsgSetSale     = types.NewMsgSetSale
	ModuleCdc         = types.ModuleCdc
	RegisterCodec     = types.RegisterCodec
	DefaultParamspace = types.DefaultParamspace
)

type (
	Keeper          = keeper.Keeper
	MsgSetName      = types.MsgSetName
	MsgBuyName      = types.MsgBuyName
	MsgDeleteName   = types.MsgDeleteName
	MsgSetSale      = types.MsgSetSale
	QueryResResolve = types.QueryResResolve
	QueryResNames   = types.QueryResNames
	Whois           = types.Whois
	QuerySaleStaus  = types.QuerySaleStatus
)
