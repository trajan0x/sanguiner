// Package gcpsigner utilizes the Key Management Service (KMS) from the Google
// Cloud Platform (GCP).
package gcpsigner

import (
	"context"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"github.com/synapsecns/sanguine/ethergo/signer/signer"
	"math/big"

	"cloud.google.com/go/kms/apiv1/kmspb"
	btcecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/googleapis/gax-go/v2"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// managedKey uses the Key Management Service (KMS) for blockchain operation.
type managedKey struct {
	KeyName string         // identifier within Google cloud project
	Addr    common.Address // public key identifier on blockchain

	// AsymmetricSign method from a Google kms.KeyManagementClient.
	asymmetricSignFunc func(context.Context, *kmspb.AsymmetricSignRequest, ...gax.CallOption) (*kmspb.AsymmetricSignResponse, error)
}

// NewManagedKey executes a fail-fast initialization.
// Key names from the Google cloud are slash-separated paths.
func NewManagedKey(ctx context.Context, client KeyClient, keyName string) (signer.Signer, error) {
	addr, err := resolveAddr(ctx, client, keyName)
	if err != nil {
		return nil, err
	}

	return &managedKey{
		KeyName:            keyName,
		Addr:               addr,
		asymmetricSignFunc: client.AsymmetricSign,
	}, nil
}

// GetTransactor returns a transaction signer for the managed key.
func (mk *managedKey) GetTransactor(ctx context.Context, chainID *big.Int) (*bind.TransactOpts, error) {
	latestSigner := types.LatestSignerForChainID(chainID)
	return &bind.TransactOpts{
		From:   mk.Addr,
		Signer: mk.newEthereumSigner(ctx, latestSigner),
	}, nil
}

func (mk *managedKey) Address() common.Address {
	return mk.Addr
}

var _ signer.Signer = &managedKey{}

func resolveAddr(ctx context.Context, client KeyClient, keyName string) (common.Address, error) {
	resp, err := client.GetPublicKey(ctx, &kmspb.GetPublicKeyRequest{Name: keyName})
	if err != nil {
		return common.Address{}, fmt.Errorf("google KMS public key %q lookup: %w", keyName, err)
	}

	block, _ := pem.Decode([]byte(resp.Pem))
	if block == nil {
		return common.Address{}, fmt.Errorf("google KMS public key %q PEM empty: %.130q", keyName, resp.Pem)
	}

	var info struct {
		AlgID pkix.AlgorithmIdentifier
		Key   asn1.BitString
	}
	_, err = asn1.Unmarshal(block.Bytes, &info)
	if err != nil {
		return common.Address{}, fmt.Errorf("google KMS public key %q PEM block %q: %w", keyName, block.Type, err)
	}

	wantAlg := asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}
	if gotAlg := info.AlgID.Algorithm; !gotAlg.Equal(wantAlg) {
		return common.Address{}, fmt.Errorf("google KMS public key %q ASN.1 algorithm %s intead of %s", keyName, gotAlg, wantAlg)
	}

	return pubKeyAddr(info.Key.Bytes), nil
}

// NewEthereumTransactor returns a KMS-backed instance. Ctx applies to the
// entire lifespan of the bind.TransactOpts.
func (mk *managedKey) NewEthereumTransactor(ctx context.Context, txIdentification types.Signer) *bind.TransactOpts {
	return &bind.TransactOpts{
		Context: ctx,
		From:    mk.Addr,
		Signer:  mk.newEthereumSigner(ctx, txIdentification),
	}
}

// SignMessage returns a signature of the message. If hash is true, the message,.
func (mk *managedKey) SignMessage(ctx context.Context, message []byte, hash bool) (signer.Signature, error) {
	if hash {
		message = crypto.Keccak256(message)
	}

	etcSig, err := mk.getSignatureFromKMS(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("could not get kms signature: %w", err)
	}

	return signer.DecodeSignature(etcSig), nil
}

// nolint: cyclop
func (mk *managedKey) getSignatureFromKMS(ctx context.Context, messageBytes []byte) ([]byte, error) {
	// resolve a signature
	req := kmspb.AsymmetricSignRequest{
		Name: mk.KeyName,
		Digest: &kmspb.Digest{
			Digest: &kmspb.Digest_Sha256{
				Sha256: messageBytes,
			},
		},
	}
	resp, err := mk.asymmetricSignFunc(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("google KMS asymmetric sign operation: %w", err)
	}

	// parse signature
	var params struct{ R, S *big.Int }
	_, err = asn1.Unmarshal(resp.Signature, &params)
	if err != nil {
		return nil, fmt.Errorf("google KMS asymmetric signature encoding: %w", err)
	}
	var rLen, sLen int // byte size
	if params.R != nil {
		rLen = (params.R.BitLen() + 7) / 8
	}
	if params.S != nil {
		sLen = (params.S.BitLen() + 7) / 8
	}
	if rLen == 0 || rLen > 32 || sLen == 0 || sLen > 32 {
		return nil, fmt.Errorf("google KMS asymmetric signature with %d-byte r and %d-byte s denied on size", rLen, sLen)
	}

	// Need uncompressed signature with "recovery ID" at end:
	// https://bitcointalk.org/index.php?topic=5249677.0
	// https://ethereum.stackexchange.com/a/53182/39582
	var sig [66]byte // + 1-byte header + 1-byte tailer
	params.R.FillBytes(sig[33-rLen : 33])
	params.S.FillBytes(sig[65-sLen : 65])

	// brute force try includes KMS verification
	var recoverErr error
	for recoveryID := byte(0); recoveryID < 2; recoveryID++ {
		sig[0] = recoveryID + 27 // BitCoin header
		btcsig := sig[:65]       // exclude Ethereum 'v' parameter
		pubKey, _, err := btcecdsa.RecoverCompact(btcsig, messageBytes)
		if err != nil {
			recoverErr = err
			continue
		}

		if pubKeyAddr(pubKey.SerializeUncompressed()) == mk.Addr {
			// sign the transaction
			sig[65] = recoveryID // Ethereum 'v' parameter
			etcsig := sig[1:]    // exclude BitCoin header
			// nolint: wrapcheck
			return etcsig, nil
		}
	}

	return nil, fmt.Errorf("google KMS asymmetric signature address recovery mis: %w", recoverErr)
}

// newEthereumSigner returns a KMS-backed instance. Ctx applies to the entire
// lifespan of the bind.SignerFn.
// nolint: cyclop
func (mk *managedKey) newEthereumSigner(ctx context.Context, txIdentification types.Signer) bind.SignerFn {
	return func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
		if addr != mk.Addr {
			return nil, bind.ErrNotAuthorized
		}

		// hash the transaction (with Keccak-256 probably)
		txHash := txIdentification.Hash(tx)

		etcSig, err := mk.getSignatureFromKMS(ctx, txHash[:])
		if err != nil {
			return nil, fmt.Errorf("could not get kms signature: %w", err)
		}

		// nolint: wrapcheck
		return tx.WithSignature(txIdentification, etcSig)
	}
}

// PubKeyAddr returns the Ethereum address for the (uncompressed) key bytes.
func pubKeyAddr(bytes []byte) common.Address {
	digest := crypto.Keccak256(bytes[1:])
	var addr common.Address
	copy(addr[:], digest[12:])
	return addr
}
