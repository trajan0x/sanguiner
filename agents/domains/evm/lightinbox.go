package evm

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/synapsecns/sanguine/agents/contracts/lightinbox"
	"github.com/synapsecns/sanguine/agents/domains"
	"github.com/synapsecns/sanguine/agents/types"
	"github.com/synapsecns/sanguine/ethergo/chain"
	"github.com/synapsecns/sanguine/ethergo/signer/nonce"
	"github.com/synapsecns/sanguine/ethergo/signer/signer"
)

// NewLightInboxContract returns a bound light inbox contract.
//
//nolint:staticcheck
func NewLightInboxContract(ctx context.Context, client chain.Chain, lightInboxAddress common.Address) (domains.LightInboxContract, error) {
	boundCountract, err := lightinbox.NewLightInboxRef(lightInboxAddress, client)
	if err != nil {
		return nil, fmt.Errorf("could not create %T: %w", &lightinbox.LightInboxRef{}, err)
	}

	nonceManager := nonce.NewNonceManager(ctx, client, client.GetBigChainID())
	return lightInboxContract{
		contract:     boundCountract,
		client:       client,
		nonceManager: nonceManager,
	}, nil
}

type lightInboxContract struct {
	// contract contains the contract handle
	contract *lightinbox.LightInboxRef
	// client contains the evm client
	//nolint: staticcheck
	client chain.Chain
	// nonceManager is the nonce manager used for transacting with the chain
	nonceManager nonce.Manager
}

func (a lightInboxContract) SubmitAttestation(
	transactor *bind.TransactOpts,
	attPayload []byte,
	signature signer.Signature,
	agentRoot [32]byte,
	snapGas []*big.Int,
) (tx *ethTypes.Transaction, err error) {
	rawSig, err := types.EncodeSignature(signature)
	if err != nil {
		return nil, fmt.Errorf("could not encode signature: %w", err)
	}

	tx, err = a.contract.SubmitAttestation(transactor, attPayload, rawSig, agentRoot, snapGas)
	if err != nil {
		return nil, fmt.Errorf("could not submit attestation: %w", err)
	}

	return tx, nil
}

func (a lightInboxContract) SubmitStateReportWithSnapshot(ctx context.Context, signer signer.Signer, stateIndex int64, signature signer.Signature, snapPayload []byte, snapSignature []byte) (tx *ethTypes.Transaction, err error) {
	transactor, err := signer.GetTransactor(ctx, a.client.GetBigChainID())
	if err != nil {
		return nil, fmt.Errorf("could not sign tx: %w", err)
	}

	transactOpts, err := a.nonceManager.NewKeyedTransactor(transactor)
	if err != nil {
		return nil, fmt.Errorf("could not create tx: %w", err)
	}

	transactOpts.Context = ctx

	transactOpts.GasLimit = 5000000

	rawSig, err := types.EncodeSignature(signature)
	if err != nil {
		return nil, fmt.Errorf("could not encode signature: %w", err)
	}

	// TODO: Is there a way to get a return value from a contractTransactor call?
	tx, err = a.contract.SubmitStateReportWithSnapshot(transactOpts, big.NewInt(stateIndex), rawSig, snapPayload, snapSignature)
	if err != nil {
		// TODO: Why is this done? And if it is necessary, we should functionalize it.
		if strings.Contains(err.Error(), "nonce too low") {
			a.nonceManager.ClearNonce(signer.Address())
		}
		return nil, fmt.Errorf("could not submit state report: %w", err)
	}

	return tx, nil
}

func (a lightInboxContract) VerifyStateWithSnapshot(ctx context.Context, signer signer.Signer, stateIndex int64, signature signer.Signature, snapPayload []byte, snapSignature []byte) (tx *ethTypes.Transaction, err error) {
	transactor, err := signer.GetTransactor(ctx, a.client.GetBigChainID())
	if err != nil {
		return nil, fmt.Errorf("could not sign tx: %w", err)
	}

	transactOpts, err := a.nonceManager.NewKeyedTransactor(transactor)
	if err != nil {
		return nil, fmt.Errorf("could not create tx: %w", err)
	}

	transactOpts.Context = ctx
	transactOpts.GasLimit = 5000000

	// TODO: Is there a way to get a return value from a contractTransactor call?
	tx, err = a.contract.VerifyStateWithSnapshot(transactOpts, big.NewInt(stateIndex), snapPayload, snapSignature)
	if err != nil {
		// TODO: Why is this done? And if it is necessary, we should functionalize it.
		if strings.Contains(err.Error(), "nonce too low") {
			a.nonceManager.ClearNonce(signer.Address())
		}
		return nil, fmt.Errorf("could not submit state report: %w", err)
	}

	return tx, nil
}

func (a lightInboxContract) SubmitAttestationReport(ctx context.Context, signer signer.Signer, attestation, arSignature, attSignature []byte) (tx *ethTypes.Transaction, err error) {
	transactor, err := signer.GetTransactor(ctx, a.client.GetBigChainID())
	if err != nil {
		return nil, fmt.Errorf("could not sign tx: %w", err)
	}

	transactOpts, err := a.nonceManager.NewKeyedTransactor(transactor)
	if err != nil {
		return nil, fmt.Errorf("could not create tx: %w", err)
	}

	transactOpts.Context = ctx

	// TODO: This will be removed once we pass in transactOpts from Submitter.
	transactOpts.GasLimit = 5000000

	// TODO: Is there a way to get a return value from a contractTransactor call?
	tx, err = a.contract.SubmitAttestationReport(transactOpts, attestation, arSignature, attSignature)
	if err != nil {
		// TODO: Why is this done? And if it is necessary, we should functionalize it.
		if strings.Contains(err.Error(), "nonce too low") {
			a.nonceManager.ClearNonce(signer.Address())
		}
		return nil, fmt.Errorf("could not submit state report: %w", err)
	}

	return tx, nil
}
