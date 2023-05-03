// autogenerated file

package config

import (
	"math/big"
	"time"
)

// IConfig ...
type IConfig interface {
	// GetMaxBatchSize returns the maximum number of transactions to send in a batch.
	GetMaxBatchSize(chainID int) int
	// GetBatch returns whether or not to batch transactions at the rpc level.
	GetBatch(chainID int) bool
	// GetMaxGasPrice returns the maximum gas price to use for transactions.
	GetMaxGasPrice(chainID int) (maxPrice *big.Int)
	// GetBumpInterval returns the number of seconds to wait before bumping a transaction
	// TODO: test this method.
	GetBumpInterval(chainID int) time.Duration
	// IsL2 returns whether or not this chain is an L2 chain.
	IsL2(chainID int) bool
	// GetGasBumpPercentage returns the percentage to bump the gas price by
	// TODO: test this method.
	GetGasBumpPercentage(chainID int) (gasBumpPercentage int)
	// GetGasEstimate returns the gas estimate to use for transactions
	// TODO: test this method.
	GetGasEstimate(chainID int) (gasEstimate uint64)
	// GetDynamicGasEstimate returns whether or not to use dynamic gas estimation
	// TODO: test this method.
	GetDynamicGasEstimate(chainID int) bool
	// SupportsEIP1559 returns whether or not this chain supports EIP1559.
	SupportsEIP1559(chainID int) bool
	// SetGlobalMaxGasPrice is a helper function that sets the global gas price.
	SetGlobalMaxGasPrice(maxPrice *big.Int)
	// SetGlobalEIP1559Support is a helper function that sets the global EIP1559 support.
	SetGlobalEIP1559Support(supportsEIP1559 bool)
}
