package mapo

import (
	"sync"
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"

	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
)

var _ exported.ClientMessage = &Header{}

var hasherPool = sync.Pool{
	New: func() interface{} { return sha3.NewLegacyKeccak256() },
}

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

func (h *Header) Hash() common.Hash {
	//Seal is reserved in extra-data. To prove block is signed by the proposer.
	if len(h.ExtraData) >= IstanbulExtraVanity {
		if istanbulHeader := IstanbulFilteredHeader(h, true); istanbulHeader != nil {
			return rlpHash(istanbulHeader)
		}
	}
	return rlpHash(h)
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

// IstanbulFilteredHeader returns a filtered header which some information (like seal, aggregated signature)
// are clean to fulfill the Istanbul hash rules. It returns nil if the extra-data cannot be
// decoded/encoded by rlp.
func IstanbulFilteredHeader(h *Header, keepSeal bool) *Header {
	newHeader := CopyHeader(h)
	istanbulExtra, err := ExtractIstanbulExtra(newHeader)
	if err != nil {
		return nil
	}

	if !keepSeal {
		istanbulExtra.Seal = []byte{}
	}
	istanbulExtra.AggregatedSeal = IstanbulAggregatedSeal{}

	payload, err := rlp.EncodeToBytes(&istanbulExtra)
	if err != nil {
		return nil
	}

	newHeader.ExtraData = append(newHeader.ExtraData[:IstanbulExtraVanity], payload...)

	return newHeader
}

// CopyHeader creates a deep copy of a block header to prevent side effects from
// modifying a header variable.
func CopyHeader(h *Header) *Header {
	cpy := *h
	if len(h.ExtraData) > 0 {
		cpy.ExtraData = make([]byte, len(h.ExtraData))
		copy(cpy.ExtraData, h.ExtraData)
	}
	return &cpy
}

// rlpHash encodes x and hashes the encoded bytes.
func rlpHash(x interface{}) (h common.Hash) {
	sha := hasherPool.Get().(crypto.KeccakState)
	defer hasherPool.Put(sha)
	sha.Reset()
	rlp.Encode(sha, x)
	sha.Read(h[:])
	return h
}
