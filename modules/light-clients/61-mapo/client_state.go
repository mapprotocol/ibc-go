package mapo

import (
	"bytes"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"

	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

var _ exported.ClientState = (*ClientState)(nil)

const (
	MAPO string = "61-mapo"
)

type TxProof struct {
	Receipt     *ethtypes.Receipt
	Proof       light.NodeList
	BlockNumber uint64
	TxIndex     uint
}

// ClientType is mapo.
func (cs ClientState) ClientType() string {
	return MAPO
}

// GetLatestHeight returns latest block height.
func (cs ClientState) GetLatestHeight() exported.Height {
	return clienttypes.NewHeight(revisionNumber, cs.LatestHeight)
}

// GetTimestampAtHeight returns the timestamp in nanoseconds of the consensus state at the given height.
func (cs ClientState) GetTimestampAtHeight(
	ctx sdk.Context,
	clientStore sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
) (uint64, error) {
	// get consensus state at height from clientStore to check for expiry
	consState, found := GetConsensusState(clientStore, cdc, height)
	if !found {
		return 0, sdkerrors.Wrapf(clienttypes.ErrConsensusStateNotFound, "height (%s)", height)
	}
	return consState.GetTimestamp(), nil
}

// Status returns the status of the mapo client.
// The client may be:
// - Active: FrozenHeight is zero and client is not expired
// - Frozen: Frozen Height is not zero
// - Expired: the latest consensus state timestamp + trusting period <= current time
//
// A frozen client will become expired, so the Frozen status
// has higher precedence.
func (cs ClientState) Status(
	ctx sdk.Context,
	clientStore sdk.KVStore,
	cdc codec.BinaryCodec,
) exported.Status {
	if cs.Frozen {
		return exported.Frozen
	}

	return exported.Active
}

// Validate performs a basic validation of the client state fields.
func (cs ClientState) Validate() error {
	if strings.TrimSpace(cs.ClientIdentifier) == "" {
		return sdkerrors.Wrap(ErrInvalidClientIdentifier, "client identifier cannot be empty string")
	}
	if cs.LatestEpoch == 0 {
		return sdkerrors.Wrap(ErrInvalidLatestEpoch, "latest epoch must be greater than zero")
	}
	if cs.EpochSize == 0 {
		return sdkerrors.Wrap(ErrInvalidEpochSize, "epoch size must be greater than zero")
	}
	if cs.LatestHeight == 0 {
		return sdkerrors.Wrap(ErrInvalidLatestHeight, "latest height must be greater than zero")
	}
	return nil
}

// ZeroCustomFields returns a ClientState that is a copy of the current ClientState
// with all client customizable fields zeroed out
func (cs ClientState) ZeroCustomFields() exported.ClientState {
	// copy over all chain-specified fields
	// and leave custom fields empty
	return &ClientState{
		Frozen:       cs.Frozen,
		LatestEpoch:  cs.LatestEpoch,
		EpochSize:    cs.EpochSize,
		LatestHeight: cs.LatestHeight,
	}
}

// Initialize checks that the initial consensus state is an 61-mapo consensus state and
// sets the client state, consensus state and associated metadata in the provided client store.
func (cs ClientState) Initialize(ctx sdk.Context, cdc codec.BinaryCodec, clientStore sdk.KVStore, consState exported.ConsensusState) error {
	consensusState, ok := consState.(*ConsensusState)
	if !ok {
		return sdkerrors.Wrapf(clienttypes.ErrInvalidConsensus, "invalid initial consensus state. expected type: %T, got: %T",
			&ConsensusState{}, consState)
	}

	setClientState(clientStore, cdc, &cs)
	setConsensusState(clientStore, cdc, consensusState, cs.GetLatestHeight())
	setConsensusMetadata(ctx, clientStore, cs.GetLatestHeight())
	return nil
}

// VerifyMembership is a generic proof verification method which verifies a proof of the existence of a value at a given CommitmentPath at the specified height.
// The caller is expected to construct the full CommitmentPath from a CommitmentPrefix and a standardized path (as defined in ICS 24).
// If a zero proof height is passed in, it will fail to retrieve the associated consensus state.
func (cs ClientState) VerifyMembership(
	ctx sdk.Context,
	clientStore sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	delayTimePeriod uint64,
	delayBlockPeriod uint64,
	proof []byte,
	path exported.Path,
	value []byte,
) error {
	if cs.GetLatestHeight().LT(height) {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"client state height < proof height (%d < %d), please ensure the client has been updated", cs.GetLatestHeight(), height,
		)
	}

	consensusState, found := GetConsensusState(clientStore, cdc, height)
	if !found {
		return sdkerrors.Wrap(clienttypes.ErrConsensusStateNotFound, "please ensure the proof was constructed against a height that exists on the client")
	}

	txProve, err := decodeTxProve(proof)
	if err != nil {
		return err
	}

	return verifyProof(consensusState.GetRoot(), txProve)
}

// VerifyNonMembership is a generic proof verification method which verifies the absence of a given CommitmentPath at a specified height.
// The caller is expected to construct the full CommitmentPath from a CommitmentPrefix and a standardized path (as defined in ICS 24).
// If a zero proof height is passed in, it will fail to retrieve the associated consensus state.
func (cs ClientState) VerifyNonMembership(
	ctx sdk.Context,
	clientStore sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	delayTimePeriod uint64,
	delayBlockPeriod uint64,
	proof []byte,
	path exported.Path,
) error {
	if cs.GetLatestHeight().LT(height) {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"client state height < proof height (%d < %d), please ensure the client has been updated", cs.GetLatestHeight(), height,
		)
	}
	return nil
}

func decodeTxProve(txProveBytes []byte) (*TxProof, error) {
	var txProof TxProof
	if err := rlp.DecodeBytes(txProveBytes, &txProof); err != nil {
		return nil, err
	}
	return &txProof, nil
}

func verifyProof(receiptsRoot common.Hash, txProof *TxProof) error {
	var buf bytes.Buffer
	rs := ethtypes.Receipts{txProof.Receipt}
	rs.EncodeIndex(0, &buf)
	giveReceipt := buf.Bytes()

	var key []byte
	key = rlp.AppendUint64(key[:0], uint64(txProof.TxIndex))

	getReceipt, err := trie.VerifyProof(receiptsRoot, key, txProof.Proof.NodeSet())
	if err != nil {
		return err
	}
	if !bytes.Equal(giveReceipt, getReceipt) {
		return sdkerrors.Wrap(ErrInvalidReceipt, "receipt mismatch")
	}
	return nil
}
