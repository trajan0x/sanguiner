package types

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/synapsecns/sanguine/core"
	"github.com/synapsecns/sanguine/ethergo/signer/signer"
)

const (
	stateOffsetRoot        = 0
	stateOffsetOrigin      = 32
	stateOffsetNonce       = 36
	stateOffsetBlockNumber = 40
	stateOffsetTimestamp   = 45
	stateOffsetGasData     = 50
	stateSize              = 62
)

// State is the state interface.
type State interface {
	// Root is the root of the Origin Merkle Tree.
	Root() [32]byte
	// Origin is the domain where Origin is located.
	Origin() uint32
	// Nonce is the amount of sent messages.
	Nonce() uint32
	// BlockNumber is the block of the last sent message.
	BlockNumber() *big.Int
	// Timestamp is the unix time of the last sent message.
	Timestamp() *big.Int
	// GasData is the gas data from the chain's gas oracle.
	GasData() GasData

	// Hash returns the hash of the state.
	Hash() ([32]byte, error)
	// SubLeaves returns the left and right sub-leaves of the state.
	SubLeaves() (leftLeaf, rightLeaf [32]byte, err error)
	// SignState returns the signature of the state payload signed by the signer.
	SignState(ctx context.Context, signer signer.Signer) (signer.Signature, []byte, common.Hash, error)
}

type state struct {
	root        [32]byte
	origin      uint32
	nonce       uint32
	blockNumber *big.Int
	timestamp   *big.Int
	gasData     GasData
}

// NewState creates a new state.
func NewState(root [32]byte, origin uint32, nonce uint32, blockNumber *big.Int, timestamp *big.Int, gasData GasData) State {
	return &state{
		root:        root,
		origin:      origin,
		nonce:       nonce,
		blockNumber: blockNumber,
		timestamp:   timestamp,
		gasData:     gasData,
	}
}

func (s state) Root() [32]byte {
	return s.root
}

func (s state) Origin() uint32 {
	return s.origin
}

func (s state) Nonce() uint32 {
	return s.nonce
}

func (s state) BlockNumber() *big.Int {
	return s.blockNumber
}

func (s state) Timestamp() *big.Int {
	return s.timestamp
}

func (s state) GasData() GasData {
	return s.gasData
}

func (s state) Hash() ([32]byte, error) {
	leftLeaf, rightLeaf, err := s.SubLeaves()
	if err != nil {
		return [32]byte{}, err
	}

	concatLeafs := append(leftLeaf[:], rightLeaf[:]...)

	return crypto.Keccak256Hash(concatLeafs), nil
}

func (s state) SubLeaves() (leftLeaf, rightLeaf [32]byte, err error) {
	encodedState, err := EncodeState(s)
	if err != nil {
		return
	}

	leftLeaf = crypto.Keccak256Hash(encodedState[stateOffsetRoot:stateOffsetNonce])
	rightLeaf = crypto.Keccak256Hash(encodedState[stateOffsetNonce:stateSize])
	return
}

//nolint:dupl
func (s state) SignState(ctx context.Context, signer signer.Signer) (signer.Signature, []byte, common.Hash, error) {
	encodedState, err := EncodeState(s)
	if err != nil {
		return nil, nil, common.Hash{}, fmt.Errorf("failed to encode state: %w", err)
	}

	stateSalt := crypto.Keccak256Hash([]byte("STATE_INVALID_SALT"))

	hashedEncodedState := crypto.Keccak256Hash(encodedState).Bytes()
	toSign := append(stateSalt.Bytes(), hashedEncodedState...)

	hashedState, err := HashRawBytes(toSign)
	if err != nil {
		return nil, nil, common.Hash{}, fmt.Errorf("failed to hash state: %w", err)
	}

	signature, err := signer.SignMessage(ctx, core.BytesToSlice(hashedState), false)
	if err != nil {
		return nil, nil, common.Hash{}, fmt.Errorf("failed to sign state: %w", err)
	}

	return signature, encodedState, hashedState, nil
}

var _ State = state{}
