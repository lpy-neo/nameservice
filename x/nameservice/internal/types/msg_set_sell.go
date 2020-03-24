package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgSetSale - struct for unjailing jailed validator
type MsgSetSale struct {
	Owner    sdk.AccAddress `json:"owner"`
	Name     string         `json:"name"`
	SaleType SaleType       `json:"saleType"`
	Price    sdk.Coins      `json:"price"`
}

// NewMsgSetSale creates a new MsgSetSale instance
func NewMsgSetSale(owner sdk.AccAddress, name string, saleType SaleType, price sdk.Coins) MsgSetSale {
	return MsgSetSale{
		Owner:    owner,
		Name:     name,
		SaleType: saleType,
		Price:    price,
	}
}

const SetSaleConst = "set_sell"

// nolint
func (msg MsgSetSale) Route() string { return RouterKey }
func (msg MsgSetSale) Type() string  { return SetSaleConst }

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgSetSale) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgSetSale) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if len(msg.Name) == 0 ||
		!IsSaleTypeValid(msg.SaleType) ||
		msg.Price.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Illegal parameters")
	}
	return nil
}

func (msg MsgSetSale) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
