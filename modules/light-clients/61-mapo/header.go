package mapo

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

var _ exported.ClientMessage = &Header{}

const revisionNumber = 0

// ConsensusState returns the updated consensus state associated with the header
func (h Header) ConsensusState() *ConsensusState {
	// todo
	return &ConsensusState{}
}

// ClientType defines that the Header is a MAPO consensus algorithm
func (h Header) ClientType() string {
	return MAPO
}

// GetHeight returns the current height. It returns 0 if the MAPO
// header is nil.
// NOTE: the header.Header is checked to be non nil in ValidateBasic.
func (h Header) GetHeight() exported.Height {
	return clienttypes.NewHeight(revisionNumber, h.Number)
}

// GetTime returns the current block timestamp. It returns a zero time if
// the MAPO header is nil.
// NOTE: the header.Header is checked to be non nil in ValidateBasic.
func (h Header) GetTime() time.Time {
	if h.SignedHeader == nil {
		return time.Time{}
	}
	return time.Unix(int64(h.SignedHeader.Timestamp), 0)
}

// ValidateBasic calls
func (h Header) ValidateBasic() error {
	if h.SignedHeader == nil {
		return sdkerrors.Wrap(clienttypes.ErrInvalidHeader, "map signed header cannot be nil")
	}
	if len(h.CommitmentRoot) == 0 {
		return sdkerrors.Wrap(clienttypes.ErrInvalidHeader, "map commitment root cannot be empty")
	}
	// todo
	return nil
}
