package mapo

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"
)

// ClientType returns MAPO
func (cs ConsensusState) ClientType() string {
	return MAPO
}

// GetRoot returns the commitment Root for the specific
func (cs ConsensusState) GetRoot() common.Hash {
	return common.BytesToHash(cs.CommitmentRoot)
}

// GetTimestamp returns block time in nanoseconds of the header that created consensus state
func (cs ConsensusState) GetTimestamp() uint64 {
	return uint64(cs.Timestamp.UnixNano())
}

// ValidateBasic defines a basic validation for the IBFT consensus state.
func (cs ConsensusState) ValidateBasic() error {
	if cs.Epoch <= 0 {
		return sdkerrors.Wrap(clienttypes.ErrInvalidConsensus, "epoch must be greater than 0")
	}
	if len(cs.Validators.PairKeys) < 1 {
		return sdkerrors.Wrap(clienttypes.ErrInvalidConsensus, "pairKeys length must be greater than 0")
	}
	if len(cs.Validators.PairKeys) != len(cs.Validators.Weights) {
		return sdkerrors.Wrap(clienttypes.ErrInvalidConsensus, "invalid validators")
	}
	if len(cs.CommitmentRoot) == 0 {
		return sdkerrors.Wrap(clienttypes.ErrInvalidConsensus, "commitment root cannot be empty")
	}
	if cs.Timestamp.Unix() <= 0 {
		return sdkerrors.Wrap(clienttypes.ErrInvalidConsensus, "timestamp must be a positive unix time")
	}
	return nil
}
