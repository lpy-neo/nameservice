package nameservice

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/lpy-neo/nameservice/x/nameservice/internal/types"
)

// NewHandler creates an sdk.Handler for all the nameservice type messages
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgSetName:
			return handleMsgSetName(ctx, keeper, msg)
		case types.MsgBuyName:
			return handleMsgBuyName(ctx, keeper, msg)
		case types.MsgDeleteName:
			return handleMsgDeleteName(ctx, keeper, msg)
		case types.MsgSetSale:
			return handleMsgSetSale(ctx, keeper, msg)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Unrecognized nameservice Msg type: %v", msg.Type()))
		}
	}
}

// handleSetName does x
func handleMsgSetName(ctx sdk.Context, keeper Keeper, msg types.MsgSetName) (*sdk.Result, error) {
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) { // Checks if the the msg sender is the same as the current owner
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner") // If not, throw an error
	}
	keeper.SetName(ctx, msg.Name, msg.Value) // If so, set the name to the value specified in the msg.
	return &sdk.Result{}, nil                // return
}

// Handle a message to buy name
func handleMsgBuyName(ctx sdk.Context, keeper Keeper, msg types.MsgBuyName) (*sdk.Result, error) {
	saleStaus := keeper.GetSaleStaus(ctx, msg.Name)
	switch saleStaus.SaleType {
	case types.SaleTypeNormal:
		return handleNormalBuy(ctx, keeper, msg)
	case types.SaleTypeAuction:
		return handleAuctionBuy(ctx, keeper, msg)
	}
	return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Name is not on sale")
}

func handleNormalBuy(ctx sdk.Context, keeper Keeper, msg types.MsgBuyName) (*sdk.Result, error) {
	// Checks if the the bid price is greater than the price paid by the current owner
	if keeper.GetPrice(ctx, msg.Name).IsAllGT(msg.Bid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid not high enough") // If not, throw an error
	}
	if keeper.HasOwner(ctx, msg.Name) {
		err := keeper.CoinKeeper.SendCoins(ctx, msg.Buyer, keeper.GetOwner(ctx, msg.Name), msg.Bid)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := keeper.CoinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Bid) // If so, deduct the Bid amount from the sender
		if err != nil {
			return nil, err
		}
	}
	keeper.SetOwner(ctx, msg.Name, msg.Buyer)
	keeper.SetPrice(ctx, msg.Name, msg.Bid)
	return &sdk.Result{}, nil
}

func handleAuctionBuy(ctx sdk.Context, keeper Keeper, msg types.MsgBuyName) (*sdk.Result, error) {
	// Checks if the the bid price is greater than the price paid by the current owner
	whois := keeper.GetWhois(ctx, msg.Name)
	if len(whois.SaleStatus.Bids) == 0 && whois.SaleStatus.Price.IsAllGT(msg.Bid) ||
		len(whois.SaleStatus.Bids) > 0 && whois.SaleStatus.Bids[len(whois.SaleStatus.Bids)-1].Price.IsAllGTE(msg.Bid) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "Bid not high enough") // If not, throw an error
	}

	if len(whois.SaleStatus.Bids) > 0 && whois.SaleStatus.Bids[len(whois.SaleStatus.Bids)-1].Buyer.Equals(msg.Buyer) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "The last bid is asking from youself")
	}

	keeper.AddBid(ctx, msg.Name, msg.Buyer, msg.Bid)

	return &sdk.Result{}, nil
}

// Handle a message to delete name
func handleMsgDeleteName(ctx sdk.Context, keeper Keeper, msg types.MsgDeleteName) (*sdk.Result, error) {
	if !keeper.IsNamePresent(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(types.ErrNameDoesNotExist, msg.Name)
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	keeper.DeleteWhois(ctx, msg.Name)
	return &sdk.Result{}, nil
}

// Handle a message to set sale
func handleMsgSetSale(ctx sdk.Context, keeper Keeper, msg types.MsgSetSale) (*sdk.Result, error) {
	if !keeper.IsNamePresent(ctx, msg.Name) {
		return nil, sdkerrors.Wrap(types.ErrNameDoesNotExist, msg.Name)
	}
	if !msg.Owner.Equals(keeper.GetOwner(ctx, msg.Name)) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Incorrect Owner")
	}

	keeper.SetSale(ctx, msg.Name, msg.SaleType, msg.Price)
	return &sdk.Result{}, nil
}
