package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgDeleteName - struct for unjailing jailed validator
type MsgDeleteName struct {
	Name  string         `json:"name"`
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgDeleteName creates a new MsgDeleteName instance
func NewMsgDeleteName(name string, owner sdk.AccAddress) MsgDeleteName {
	return MsgDeleteName{
		Name:  name,
		Owner: owner,
	}
}

const DeleteNameConst = "delete_name"

// nolint
func (msg MsgDeleteName) Route() string { return RouterKey }
func (msg MsgDeleteName) Type() string  { return DeleteNameConst }

// ValidateBasic runs stateless checks on the message
func (msg MsgDeleteName) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if len(msg.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Name cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgDeleteName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgDeleteName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
