package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgBuyName - struct for unjailing jailed validator
type MsgBuyName struct {
	Name  string         `json:"name"`
	Bid   sdk.Coins      `json:"bid"`
	Buyer sdk.AccAddress `json:"buyer"`
}

// NewMsgBuyName creates a new MsgBuyName instance
func NewMsgBuyName(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyName {
	return MsgBuyName{
		Name:  name,
		Bid:   bid,
		Buyer: buyer,
	}
}

const BuyNameConst = "buy_name"

// nolint
func (msg MsgBuyName) Route() string { return RouterKey }
func (msg MsgBuyName) Type() string  { return BuyNameConst }

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgBuyName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgBuyName) ValidateBasic() error {
	if msg.Buyer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Buyer.String())
	}
	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name cannot be empty")
	}
	if !msg.Bid.IsAllPositive() {
		return sdkerrors.ErrInsufficientFunds
	}
	return nil
}

func (msg MsgBuyName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}
