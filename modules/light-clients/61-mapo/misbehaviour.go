package mapo

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

var _ exported.ClientMessage = &Misbehaviour{}

// ClientType is Tendermint light client
func (m Misbehaviour) ClientType() string {
	return MAPO
}

// GetTime returns the timestamp at which misbehaviour occurred. It uses the
// maximum value from both headers to prevent producing an invalid header outside of
// the misbehaviour age range.
func (m Misbehaviour) GetTime() time.Time {
	t1, t2 := m.Header1.GetTime(), m.Header2.GetTime()
	if t1.After(t2) {
		return t1
	}
	return t2
}

// ValidateBasic implements Misbehaviour interface
func (m Misbehaviour) ValidateBasic() error {
	if m.Header1 == nil {
		return sdkerrors.Wrap(ErrInvalidHeader, "misbehaviour Header1 cannot be nil")
	}
	if m.Header2 == nil {
		return sdkerrors.Wrap(ErrInvalidHeader, "misbehaviour Header2 cannot be nil")
	}

	// ValidateBasic on both validators
	if err := m.Header1.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(
			clienttypes.ErrInvalidMisbehaviour,
			sdkerrors.Wrap(err, "header 1 failed validation").Error(),
		)
	}
	if err := m.Header2.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(
			clienttypes.ErrInvalidMisbehaviour,
			sdkerrors.Wrap(err, "header 2 failed validation").Error(),
		)
	}
	// Ensure that Height1 is greater than or equal to Height2
	if m.Header1.GetHeight().LT(m.Header2.GetHeight()) {
		return sdkerrors.Wrapf(clienttypes.ErrInvalidMisbehaviour, "Header1 height is less than Header2 height (%s < %s)", m.Header1.GetHeight(), m.Header2.GetHeight())
	}
	return nil
}
