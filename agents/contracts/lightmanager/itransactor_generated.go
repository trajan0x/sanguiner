// autogenerated file

package lightmanager

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ILightManagerTransactor ...
type ILightManagerTransactor interface {
	// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
	//
	// Solidity: function initialize(address origin_, address destination_, address inbox_) returns()
	Initialize(opts *bind.TransactOpts, origin_ common.Address, destination_ common.Address, inbox_ common.Address) (*types.Transaction, error)
	// Multicall is a paid mutator transaction binding the contract method 0x60fc8466.
	//
	// Solidity: function multicall((bool,bytes)[] calls) returns((bool,bytes)[] callResults)
	Multicall(opts *bind.TransactOpts, calls []MultiCallableCall) (*types.Transaction, error)
	// OpenDispute is a paid mutator transaction binding the contract method 0xa2155c34.
	//
	// Solidity: function openDispute(uint32 guardIndex, uint32 notaryIndex) returns()
	OpenDispute(opts *bind.TransactOpts, guardIndex uint32, notaryIndex uint32) (*types.Transaction, error)
	// RemoteWithdrawTips is a paid mutator transaction binding the contract method 0x1fa07138.
	//
	// Solidity: function remoteWithdrawTips(uint32 msgOrigin, uint256 proofMaturity, address recipient, uint256 amount) returns(bytes4 magicValue)
	RemoteWithdrawTips(opts *bind.TransactOpts, msgOrigin uint32, proofMaturity *big.Int, recipient common.Address, amount *big.Int) (*types.Transaction, error)
	// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
	//
	// Solidity: function renounceOwnership() returns()
	RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error)
	// ResolveStuckDispute is a paid mutator transaction binding the contract method 0x89791e17.
	//
	// Solidity: function resolveStuckDispute(uint32 domain, address slashedAgent) returns()
	ResolveStuckDispute(opts *bind.TransactOpts, domain uint32, slashedAgent common.Address) (*types.Transaction, error)
	// SetAgentRoot is a paid mutator transaction binding the contract method 0x58668176.
	//
	// Solidity: function setAgentRoot(bytes32 agentRoot_) returns()
	SetAgentRoot(opts *bind.TransactOpts, agentRoot_ [32]byte) (*types.Transaction, error)
	// SlashAgent is a paid mutator transaction binding the contract method 0x2853a0e6.
	//
	// Solidity: function slashAgent(uint32 domain, address agent, address prover) returns()
	SlashAgent(opts *bind.TransactOpts, domain uint32, agent common.Address, prover common.Address) (*types.Transaction, error)
	// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
	//
	// Solidity: function transferOwnership(address newOwner) returns()
	TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error)
	// UpdateAgentStatus is a paid mutator transaction binding the contract method 0xcbd05965.
	//
	// Solidity: function updateAgentStatus(address agent, (uint8,uint32,uint32) status, bytes32[] proof) returns()
	UpdateAgentStatus(opts *bind.TransactOpts, agent common.Address, status AgentStatus, proof [][32]byte) (*types.Transaction, error)
}