package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	emoji "github.com/tmdvs/Go-Emoji-Utils"
)

// RouterKey is the module name router key
const RouterKey = ModuleName // this was defined in your key.go file

// MsgMint defines a mint message
type MsgMint struct {
	Sender sdk.AccAddress
	Minter sdk.AccAddress
	Denom  string
}

// NewMsgMint is a constructor function for NewMsgMint
func NewMsgMint(sender sdk.AccAddress, minter sdk.AccAddress, denom string) MsgMint {
	return MsgMint{Sender: sender, Minter: minter, Denom: denom}
}

// Route should return the name of the module
func (msg MsgMint) Route() string { return RouterKey }

// Type should return the action
func (msg MsgMint) Type() string { return "mint" }

// ValidateBasic runs stateless checks on the message
func (msg MsgMint) ValidateBasic() error {
	if msg.Minter.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Minter.String())
	}
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}
	if msg.Sender.Equals(msg.Minter) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Can't mint to yourself")
	}
	results := emoji.FindAll(msg.Denom)
	if len(results) != 1 {
		return ErrNoEmoji
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgMint) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgMint) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// MsgMint defines a mint message
type MsgFaucetKey struct {
	Sender sdk.AccAddress
	Armor  string
}

// NewMsgFaucetKey is a constructor function for MsgFaucetKey
func NewMsgFaucetKey(sender sdk.AccAddress, armor string) MsgFaucetKey {
	return MsgFaucetKey{Sender: sender, Armor: armor}
}

// Route should return the name of the module
func (msg MsgFaucetKey) Route() string { return RouterKey }

// Type should return the action
func (msg MsgFaucetKey) Type() string { return "faucet-key" }

// ValidateBasic runs stateless checks on the message
func (msg MsgFaucetKey) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Sender.String())
	}
	if len(msg.Armor) == 0 {
		return ErrFaucetKeyEmpty
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgFaucetKey) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgFaucetKey) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
