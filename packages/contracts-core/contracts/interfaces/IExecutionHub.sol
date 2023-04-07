// SPDX-License-Identifier: MIT
pragma solidity 0.8.17;

import {MessageStatus} from "../libs/Structures.sol";

interface IExecutionHub {
    /**
     * @notice Attempts to prove inclusion of message into one of Snapshot Merkle Trees,
     * previously submitted to this contract in a form of a signed Attestation.
     * Proven message is immediately executed by passing its contents to the specified recipient.
     * @dev Will revert if any of these is true:
     *  - Message payload is not properly formatted.
     *  - Snapshot root (reconstructed from message hash and proofs) is unknown
     *  - Snapshot root is known, but was submitted by an inactive Notary
     *  - Snapshot root is known, but optimistic period for a message hasn't passed
     *  - Provided gas limit is lower than the one requested in the message
     *  - Recipient doesn't implement a `handle` method (refer to IMessageRecipient.sol)
     *  - Recipient reverted upon receiving a message
     * Note: refer to libs/State.sol for details about Origin State's sub-leafs.
     * @param msgPayload    Raw payload with a formatted message to execute
     * @param originProof   Proof of inclusion of message in the Origin Merkle Tree
     * @param snapProof     Proof of inclusion of Origin State's Left Leaf into Snapshot Merkle Tree
     * @param stateIndex    Index of Origin State in the Snapshot
     * @param gasLimit      Gas limit for message execution
     */
    function execute(
        bytes memory msgPayload,
        bytes32[] calldata originProof,
        bytes32[] calldata snapProof,
        uint256 stateIndex,
        uint64 gasLimit
    ) external;

    // ═══════════════════════════════════════════════════ VIEWS ═══════════════════════════════════════════════════════

    /**
     * @notice Returns message execution status: None/Failed/Success.
     * @param origin            Domain where message originated
     * @param nonce             Message nonce on the origin domain
     * @return flag             Message execution status
     */
    function messageStatus(uint32 origin, uint32 nonce) external view returns (MessageStatus flag);

    /**
     * @notice Returns a formatted payload with message execution data.
     * @dev Notaries could sign this statement and submit it to Summit in order to distribute message tips.
     * @param origin            Domain where message originated
     * @param nonce             Message nonce on the origin domain
     * @return executionPayload Formatted Execution payload for the message.
     */
    function executionStatus(uint32 origin, uint32 nonce) external view returns (bytes memory executionPayload);
}
