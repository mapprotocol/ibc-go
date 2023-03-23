package mapo

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// IBC MAPO client sentinel errors
var (
	// todo code replace
	ErrInvalidClientIdentifier    = sdkerrors.Register(ModuleName, 1, "invalid client identifier")
	ErrInvalidLatestEpoch         = sdkerrors.Register(ModuleName, 1, "invalid latest epoch")
	ErrInvalidEpochSize           = sdkerrors.Register(ModuleName, 1, "invalid epoch size")
	ErrInvalidLatestHeight        = sdkerrors.Register(ModuleName, 1, "invalid latest height")
	ErrInvalidReceipt             = sdkerrors.Register(ModuleName, 1, "invalid receipt")
	ErrNotLastHeaderOfEpoch       = sdkerrors.Register(ModuleName, 1, "not the last header of epoch")
	ErrInvalidHeader              = sdkerrors.Register(ModuleName, 1, "invalid header")
	ErrInvalidIstanbulHeaderExtra = sdkerrors.Register(ModuleName, 2, "invalid istanbul header extra-data")
	ErrEmptyAggregatedSeal        = sdkerrors.Register(ModuleName, 3, "empty aggregated seal")
	ErrInvalidAggregatedSeal      = sdkerrors.Register(ModuleName, 3, "invalid aggregated seal")
	ErrInsufficientSeals          = sdkerrors.Register(ModuleName, 4, "not enough seals to reach quorum")
	ErrInvalidSignature           = sdkerrors.Register(ModuleName, 5, "invalid signature")
)
