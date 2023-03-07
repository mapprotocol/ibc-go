package mapo

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// IBC MAPO client sentinel errors
var (
	ErrNotLastHeaderOfEpoch       = sdkerrors.Register(ModuleName, 1, "not the last header of epoch")
	ErrInvalidHeader              = sdkerrors.Register(ModuleName, 1, "invalid header")
	ErrInvalidIstanbulHeaderExtra = sdkerrors.Register(ModuleName, 2, "invalid istanbul header extra-data")
	ErrEmptyAggregatedSeal        = sdkerrors.Register(ModuleName, 3, "empty aggregated seal")
	ErrInvalidAggregatedSeal      = sdkerrors.Register(ModuleName, 3, "invalid aggregated seal")
	ErrInsufficientSeals          = sdkerrors.Register(ModuleName, 4, "not enough seals to reach quorum")
	ErrInvalidSignature           = sdkerrors.Register(ModuleName, 5, "invalid signature")
)
