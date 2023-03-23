package mapo

import (
	"fmt"
	"math"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	istanbulCore "github.com/mapprotocol/atlas/consensus/istanbul/core"
	blscrypto "github.com/mapprotocol/atlas/helper/bls"

	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

const (
	EpochSize                 = 50000
	BN256Fork                 = 2350001
	IstanbulExtraVanity       = 32 // Fixed number of extra-data bytes reserved for validator vanity
	IstanbulExtraBlsSignature = 64 // Fixed number of extra-data bytes reserved for validator seal on the current block
	IstanbulExtraSeal         = 65 // Fixed number of extra-data bytes reserved for validator seal
	PublicKeyBytes            = 128
)

// VerifyClientMessage checks if the clientMessage is of type Header or Misbehaviour and verifies the message
func (cs *ClientState) VerifyClientMessage(
	ctx sdk.Context, cdc codec.BinaryCodec, clientStore sdk.KVStore,
	clientMsg exported.ClientMessage,
) error {
	switch msg := clientMsg.(type) {
	case *Header:
		return cs.verifyHeader(ctx, clientStore, cdc, msg)
	case *Misbehaviour:
		return cs.verifyMisbehaviour(ctx, clientStore, cdc, msg)
	default:
		return clienttypes.ErrInvalidClientType
	}
}

func (cs *ClientState) verifyHeader(
	ctx sdk.Context, clientStore sdk.KVStore, cdc codec.BinaryCodec,
	header *Header,
) error {
	if header.Number == 0 {
		return ErrInvalidHeader
	}
	if header.Number%EpochSize != 0 {
		return ErrNotLastHeaderOfEpoch
	}
	extra, err := ExtractIstanbulExtra(header)
	if err != nil {
		return err
	}
	// The length of Committed seals should be larger than 0
	if len(extra.AggregatedSeal.Signature) == 0 {
		return ErrEmptyAggregatedSeal
	}
	// todo replace
	// Retrieve trusted consensus states for each Header in misbehaviour
	consState, found := GetConsensusState(clientStore, cdc, header.GetHeight())
	if !found {
		return sdkerrors.Wrapf(clienttypes.ErrConsensusStateNotFound, "could not get consensus state from clientStore for Header at Number: %s", header.Number)
	}

	pks := consState.Validators.PairKeys
	validators := make([]SerializedPublicKey, 0, len(pks))
	for i, pk := range pks {
		copy(validators[i][:], pk.G2PubKey)
	}

	vs := validatorSet{
		validators: validators,
	}
	fork, cur := big.NewInt(BN256Fork), new(big.Int).SetUint64(header.Number)
	return verifyAggregatedSeal(header.Hash(), vs, extra.AggregatedSeal, fork, cur)
}

func verifyAggregatedSeal(hash common.Hash, validators validatorSet, aggregatedSeal IstanbulAggregatedSeal, fork, cur *big.Int) error {
	if len(aggregatedSeal.Signature) != IstanbulExtraBlsSignature {
		return ErrInvalidAggregatedSeal
	}

	proposalSeal := istanbulCore.PrepareCommittedSeal(hash, aggregatedSeal.Round)
	// Find which public keys signed from the provided validator set
	publicKeys := make([]blscrypto.SerializedPublicKey, 0)
	for i := 0; i < validators.Size(); i++ {
		if aggregatedSeal.Bitmap.Bit(i) == 1 {
			pubKey := validators.BLSPublicKey(i)
			publicKeys = append(publicKeys, blscrypto.SerializedPublicKey(pubKey))
		}
	}
	// The length of a valid seal should be greater than the minimum quorum size
	if len(publicKeys) < validators.MinQuorumSize() {
		return ErrInsufficientSeals
	}
	if err := blscrypto.CryptoType().VerifyAggregatedSignature(publicKeys, proposalSeal, []byte{}, aggregatedSeal.Signature,
		false, false, fork, cur); err != nil {
		return ErrInvalidSignature
	}

	return nil
}

type SerializedPublicKey [PublicKeyBytes]byte

type validatorSet struct {
	validators []SerializedPublicKey
}

func (vs *validatorSet) Size() int {
	return len(vs.validators)
}

func (vs *validatorSet) BLSPublicKey(i int) SerializedPublicKey {
	if i < vs.Size() {
		return vs.validators[i]
	}
	return SerializedPublicKey{}
}

func (vs *validatorSet) MinQuorumSize() int {
	return int(math.Ceil(float64(2*vs.Size()) / 3))
}

type IstanbulAggregatedSeal struct {
	// Bitmap is a bitmap having an active bit for each validator that signed this block
	Bitmap *big.Int
	// Signature is an aggregated BLS signature resulting from signatures by each validator that signed this block
	Signature []byte
	// Round is the round in which the signature was created.
	Round *big.Int
}

type IstanbulExtra struct {
	// AddedValidators are the validators that have been added in the block
	AddedValidators []common.Address
	// AddedValidatorsPublicKeys are the BLS public keys for the validators added in the block
	AddedValidatorsPublicKeys []blscrypto.SerializedPublicKey
	// AddedValidatorsG1PublicKeys are the BLS public keys for the validators added in the block
	AddedValidatorsG1PublicKeys []blscrypto.SerializedG1PublicKey
	// RemovedValidators is a bitmap having an active bit for each removed validator in the block
	RemovedValidators *big.Int
	// Seal is an ECDSA signature by the proposer
	Seal []byte
	// AggregatedSeal contains the aggregated BLS signature created via IBFT consensus.
	AggregatedSeal IstanbulAggregatedSeal
	// ParentAggregatedSeal contains and aggregated BLS signature for the previous block.
	ParentAggregatedSeal IstanbulAggregatedSeal
}

func ExtractIstanbulExtra(h *Header) (*IstanbulExtra, error) {
	if len(h.ExtraData) < IstanbulExtraVanity {
		return nil, ErrInvalidIstanbulHeaderExtra
	}

	var istanbulExtra *IstanbulExtra
	err := rlp.DecodeBytes(h.ExtraData[IstanbulExtraVanity:], &istanbulExtra)
	if err != nil {
		return nil, err
	}
	return istanbulExtra, nil
}

// UpdateState may be used to either create a consensus state for:
// - a future height greater than the latest client state height
// - a past height that was skipped during bisection
// If we are updating to a past height, a consensus state is created for that height to be persisted in client store
// If we are updating to a future height, the consensus state is created and the client state is updated to reflect
// the new latest height
// A list containing the updated consensus height is returned.
// UpdateState must only be used to update within a single revision, thus header revision number and trusted height's revision
// number must be the same. To update to a new revision, use a separate upgrade path
// UpdateState will prune the oldest consensus state if it is expired.
func (cs ClientState) UpdateState(ctx sdk.Context, cdc codec.BinaryCodec, clientStore sdk.KVStore, clientMsg exported.ClientMessage) []exported.Height {
	header, ok := clientMsg.(*Header)
	if !ok {
		panic(fmt.Errorf("expected type %T, got %T", &Header{}, clientMsg))
	}

	cs.pruneOldestConsensusState(ctx, cdc, clientStore)

	// check for duplicate update
	if consensusState, _ := GetConsensusState(clientStore, cdc, header.GetHeight()); consensusState != nil {
		// perform no-op
		return []exported.Height{header.GetHeight()}
	}

	consensusState := &ConsensusState{
		Epoch:          0,
		Validators:     nil, // todo
		CommitmentRoot: header.CommitmentRoot,
		Timestamp:      header.GetTime(),
	}

	// set client state, consensus state and asssociated metadata
	setClientState(clientStore, cdc, &cs)
	setConsensusState(clientStore, cdc, consensusState, header.GetHeight())
	setConsensusMetadata(ctx, clientStore, header.GetHeight())

	return []exported.Height{header.GetHeight()}
}

// pruneOldestConsensusState will retrieve the earliest consensus state for this clientID and check if it is expired. If it is,
// that consensus state will be pruned from store along with all associated metadata. This will prevent the client store from
// becoming bloated with expired consensus states that can no longer be used for updates and packet verification.
func (cs ClientState) pruneOldestConsensusState(ctx sdk.Context, cdc codec.BinaryCodec, clientStore sdk.KVStore) {
	// todo
}

// UpdateStateOnMisbehaviour updates state upon misbehaviour, freezing the ClientState. This method should only be called when misbehaviour is detected
// as it does not perform any misbehaviour checks.
func (cs ClientState) UpdateStateOnMisbehaviour(ctx sdk.Context, cdc codec.BinaryCodec, clientStore sdk.KVStore, _ exported.ClientMessage) {
	// todo
}
