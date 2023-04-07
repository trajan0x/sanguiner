// SPDX-License-Identifier: MIT
pragma solidity 0.8.17;

/// @notice A collection of events emitted by the ExecutionHub contract
abstract contract ExecutionHubEvents {
    /**
     * @notice Emitted when message is executed.
     * @param remoteDomain  Remote domain where message originated
     * @param messageHash   The keccak256 hash of the message that was executed
     */
    event Executed(uint32 indexed remoteDomain, bytes32 indexed messageHash);

    /**
     * @notice Emitted when message tips are recorded.
     * @param origin        Domain where message originated
     * @param nonce         Message nonce on origin chain
     * @param tipsPayload   Raw payload with tips paid for the off-chain agents
     */
    event TipsRecorded(uint32 origin, uint32 nonce, bytes tipsPayload);
}
