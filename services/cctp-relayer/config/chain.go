package config

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/richardwilkes/toolbox/collection"
)

// ChainConfig defines the config for a specific chain.
type ChainConfig struct {
	// TempRPC is the temporary RPC endpoint for the chain.
	TempRPC string `yaml:"temp_rpc"`
	// ChainID is the ID of the chain.
	ChainID uint32 `yaml:"chain_id"`
	// OriginAddress is the address of the origin contract.
	OriginAddress string `yaml:"origin_address"`
	// DestinationAddress is the address of the destination contract.
	DestinationAddress string `yaml:"destination_address"`
}

// GetOriginAddress returns the origin address.
func (c ChainConfig) GetOriginAddress() common.Address {
	return common.HexToAddress(c.OriginAddress)
}

// GetDestinationAddress returns the destination address.
func (c ChainConfig) GetDestinationAddress() common.Address {
	return common.HexToAddress(c.DestinationAddress)
}

// ChainConfigs contains an array of ChainConfigs.
type ChainConfigs []ChainConfig

// IsValid validates the chain config by asserting no two chains appear twice.
func (c ChainConfigs) IsValid(ctx context.Context) (ok bool, err error) {
	intSet := collection.Set[uint32]{}

	for _, cfg := range c {
		if intSet.Contains(cfg.ChainID) {
			return false, fmt.Errorf("chain id %d appears twice: %s", cfg.ChainID, "duplicate chain id")
		}
		intSet.Add(cfg.ChainID)
	}

	return true, nil
}

// IsValid validates the chain config.
func (c ChainConfig) IsValid(ctx context.Context) (ok bool, err error) {
	if c.ChainID == 0 {
		return false, fmt.Errorf("%s: chain ID cannot be 0", "invalid chain id")
	}

	if c.OriginAddress == "" {
		return false, fmt.Errorf("field OriginAddress is required")
	}

	if c.DestinationAddress == "" {
		return false, fmt.Errorf("field DestinationAddress is required")
	}

	return true, nil
}
