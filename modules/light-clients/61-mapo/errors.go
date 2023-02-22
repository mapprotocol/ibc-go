package mapo

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// IBC MAPO client sentinel errors
var (
	ErrInvalidHeader = sdkerrors.Register(ModuleName, 6, "invalid header")
)
