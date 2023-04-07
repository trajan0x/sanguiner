// SPDX-License-Identifier: MIT
pragma solidity 0.8.17;

// ══════════════════════════════ LIBRARY IMPORTS ══════════════════════════════
import {Attestation, ExecutionAttestation} from "../libs/Attestation.sol";
import {BaseMessage, BaseMessageLib} from "../libs/BaseMessage.sol";
import {Composite} from "../libs/Composite.sol";
import {SYSTEM_ROUTER, ORIGIN_TREE_HEIGHT, SNAPSHOT_TREE_HEIGHT} from "../libs/Constants.sol";
import {Execution, ExecutionLib} from "../libs/Execution.sol";
import {MerkleLib} from "../libs/Merkle.sol";
import {Header, Message, MessageFlag, MessageLib} from "../libs/Message.sol";
import {MessageStatus} from "../libs/Structures.sol";
import {SystemMessage, SystemMessageLib} from "../libs/SystemMessage.sol";
import {Tips} from "../libs/Tips.sol";
import {TypeCasts} from "../libs/TypeCasts.sol";
import {TypedMemView} from "../libs/TypedMemView.sol";
// ═════════════════════════════ INTERNAL IMPORTS ══════════════════════════════
import {DisputeHub} from "./DisputeHub.sol";
import {ExecutionHubEvents} from "../events/ExecutionHubEvents.sol";
import {IExecutionHub} from "../interfaces/IExecutionHub.sol";
import {IMessageRecipient} from "../interfaces/IMessageRecipient.sol";

/**
 * @notice ExecutionHub is responsible for executing the messages that are
 * proven against the Snapshot Merkle Roots.
 * The Snapshot Merkle Roots themselves are supposed to be dealt with in the child contracts.
 * On the Synapse Chain Notaries are submitting the snapshots that are later used for proving.
 * On the other chains Notaries are submitting the attestations that are later used for proving.
 */
abstract contract ExecutionHub is DisputeHub, ExecutionHubEvents, IExecutionHub {
    using BaseMessageLib for bytes29;
    using MessageLib for bytes;
    using TypeCasts for bytes32;
    using TypedMemView for bytes29;

    /// @notice Struct representing the status of Message in Execution Hub.
    /// @param status       Message execution status
    /// @param attNonce     Nonce of the attestation used for proving the message
    /// @param executor     Executor who successfully executed the message
    struct ExecutionStatus {
        MessageStatus status;
        uint32 attNonce;
        address executor;
    }
    // 16 bits available for tight packing

    // ══════════════════════════════════════════════════ STORAGE ══════════════════════════════════════════════════════

    /// @notice ((origin, nonce) => status)
    /// @dev Messages coming from different origins will always have a different hash
    /// as origin domain is encoded into the formatted message.
    /// Thus we can use hash as a key instead of an (origin, hash) tuple.
    mapping(uint64 => ExecutionStatus) private _executionStatus;

    /// @notice First executor who made a valid attempt of executing a message.
    /// Note: stored only for messages that had Failed status at some point of time
    mapping(uint64 => address) private _firstExecutor;

    /// @dev Tracks all saved attestations
    // (root => attestation)
    mapping(bytes32 => ExecutionAttestation) private _rootAttestations;

    /// @dev gap for upgrade safety
    uint256[47] private __GAP; // solhint-disable-line var-name-mixedcase

    // ═════════════════════════════════════════════ EXECUTE MESSAGES ══════════════════════════════════════════════════

    /// @inheritdoc IExecutionHub
    function execute(
        bytes memory msgPayload,
        bytes32[] calldata originProof,
        bytes32[] calldata snapProof,
        uint256 stateIndex,
        uint64 gasLimit
    ) external {
        // TODO: add reentrancy check
        // This will revert if payload is not a formatted message payload
        Message message = msgPayload.castToMessage();
        Header header = message.header();
        bytes32 msgLeaf = message.leaf();
        // Ensure message was meant for this domain
        require(header.destination() == localDomain, "!destination");
        // Check that message has not been executed before
        uint64 originAndNonce = header.originAndNonce();
        ExecutionStatus memory execStatus = _executionStatus[originAndNonce];
        require(execStatus.status != MessageStatus.Success, "Already executed");
        // Check proofs validity
        ExecutionAttestation memory execAtt =
            _proveAttestation(originAndNonce, msgLeaf, originProof, snapProof, stateIndex);
        // Check if optimistic period has passed
        uint256 proofMaturity = block.timestamp - execAtt.submittedAt;
        require(proofMaturity >= header.optimisticPeriod(), "!optimisticPeriod");
        bool success;
        // Only System/Base message flags exist
        if (message.flag() == MessageFlag.System) {
            // gasLimit is ignored when executing system messages
            success = _executeSystemMessage(originAndNonce, proofMaturity, message.body());
        } else {
            // This will revert if message body is not a formatted BaseMessage payload
            success = _executeBaseMessage(originAndNonce, proofMaturity, gasLimit, message.body().castToBaseMessage());
        }
        if (execStatus.status == MessageStatus.None) {
            // This is the first valid attempt to execute the message, save the attestation nonce
            execStatus.attNonce = execAtt.nonce;
            if (success) {
                // This is the successful attempt to execute the message => save the executor
                execStatus.status = MessageStatus.Success;
                execStatus.executor = msg.sender;
            } else {
                // Save as the "first executor", if execution failed
                execStatus.status = MessageStatus.Failed;
                _firstExecutor[originAndNonce] = msg.sender;
            }
            _executionStatus[originAndNonce] = execStatus;
        } else if (success) {
            // There has been a failed attempt to execute the message before => don't touch attNonce
            // This is the successful attempt to execute the message => save the executor
            execStatus.status = MessageStatus.Success;
            execStatus.executor = msg.sender;
            _executionStatus[originAndNonce] = execStatus;
        }
        emit Executed(header.origin(), msgLeaf);
    }

    // ═══════════════════════════════════════════════════ VIEWS ═══════════════════════════════════════════════════════

    /// @inheritdoc IExecutionHub
    function messageStatus(uint32 origin, uint32 nonce) external view returns (MessageStatus) {
        return _executionStatus[Composite.mergeUint32(origin, nonce)].status;
    }

    /// @inheritdoc IExecutionHub
    function executionData(uint32 origin, uint32 nonce) external view returns (bytes memory data) {
        uint64 originAndNonce = Composite.mergeUint32(origin, nonce);
        ExecutionStatus memory execStatus = _executionStatus[originAndNonce];
        // Return empty payload if there has been no attempt to execute the message
        if (execStatus.status == MessageStatus.None) return "";
        address firstExecutor = _firstExecutor[originAndNonce];
        if (firstExecutor == address(0)) firstExecutor = execStatus.executor;
        // ExecutionHub does not store the tips, the Notary will have to append the tips payload
        return ExecutionLib.formatExecution(
            execStatus.status, origin, nonce, localDomain, execStatus.attNonce, firstExecutor, execStatus.executor, ""
        );
    }

    // ═════════════════════════════════════ INTERNAL LOGIC: MESSAGE EXECUTION ═════════════════════════════════════════

    /// @dev Passes message content to recipient that conforms to IMessageRecipient interface.
    function _executeBaseMessage(uint64 originAndNonce, uint256 proofMaturity, uint64 gasLimit, BaseMessage baseMessage)
        internal
        returns (bool)
    {
        // Check that gas limit covers the one requested by the sender.
        // We let the executor specify gas limit higher than requested to guarantee the execution of
        // messages with gas limit set too low.
        require(gasLimit >= baseMessage.request().gasLimit(), "Gas limit too low");
        (uint32 origin, uint32 nonce) = Composite.splitUint32(originAndNonce);
        // Emit event with the message tips
        emit TipsRecorded(origin, nonce, baseMessage.tips().unwrap().clone());
        // TODO: check that the discarded bits are empty
        address recipient = baseMessage.recipient().bytes32ToAddress();
        // Forward message content to the recipient, and limit the amount of forwarded gas
        require(gasleft() > gasLimit, "Not enough gas supplied");
        try IMessageRecipient(recipient).receiveBaseMessage{gas: gasLimit}(
            origin, nonce, baseMessage.sender(), proofMaturity, baseMessage.content().clone()
        ) {
            return true;
        } catch {
            return false;
        }
    }

    function _executeSystemMessage(uint64 originAndNonce, uint256 proofMaturity, bytes29 body)
        internal
        returns (bool)
    {
        (uint32 origin, uint32 nonce) = Composite.splitUint32(originAndNonce);
        // TODO: introduce incentives for executing System Messages?
        // Forward system message to System Router
        systemRouter.receiveSystemMessage(origin, nonce, proofMaturity, body.clone());
        return true;
    }

    // ══════════════════════════════════════ INTERNAL LOGIC: MESSAGE PROVING ══════════════════════════════════════════

    /// @dev Saves a snapshot root with the attestation data provided by a Notary.
    /// It is assumed that the Notary signature has been checked outside of this contract.
    function _saveAttestation(Attestation att, address notary) internal {
        bytes32 root = att.snapRoot();
        require(_rootAttestations[root].isEmpty(), "Root already exists");
        _rootAttestations[root] = att.toExecutionAttestation(notary);
    }

    /**
     * @notice Attempts to prove the validity of the cross-chain message.
     * First, the origin Merkle Root is reconstructed using the origin proof.
     * Then the origin state's "left leaf" is reconstructed using the origin domain.
     * After that the snapshot Merkle Root is reconstructed using the snapshot proof.
     * The snapshot root needs to have been submitted by an undisputed Notary.
     * @dev Reverts if any of the checks fail.
     * @param msgLeaf       Message Leaf that was inserted in the Origin Merkle Tree
     * @param originProof   Proof of inclusion of Message Leaf in the Origin Merkle Tree
     * @param snapProof     Proof of inclusion of Origin State Left Leaf into Snapshot Merkle Tree
     * @param stateIndex    Index of Origin State in the Snapshot
     * @return execAtt      Attestation data for derived snapshot root
     */
    function _proveAttestation(
        uint64 originAndNonce,
        bytes32 msgLeaf,
        bytes32[] calldata originProof,
        bytes32[] calldata snapProof,
        uint256 stateIndex
    ) internal view returns (ExecutionAttestation memory execAtt) {
        (uint32 origin, uint32 nonce) = Composite.splitUint32(originAndNonce);
        // Reconstruct Origin Merkle Root using the origin proof
        // Message index in the tree is (nonce - 1), as nonce starts from 1
        // This will revert if origin proof length exceeds Origin Tree height
        bytes32 originRoot = MerkleLib.proofRoot(nonce - 1, msgLeaf, originProof, ORIGIN_TREE_HEIGHT);
        // Reconstruct Snapshot Merkle Root using the snapshot proof
        // This will revert if:
        //  - State index is out of range.
        //  - Snapshot Proof length exceeds Snapshot tree Height.
        bytes32 snapshotRoot = _snapshotRoot(originRoot, origin, snapProof, stateIndex);
        // Fetch the attestation data for the snapshot root
        execAtt = _rootAttestations[snapshotRoot];
        // Check if snapshot root has been submitted
        require(!execAtt.isEmpty(), "Invalid snapshot root");
        // Check if Notary who submitted the attestation is still active
        _verifyActive(_agentStatus(execAtt.notary));
        // Check that Notary who submitted the attestation is not in dispute
        require(!_inDispute(execAtt.notary), "Notary is in dispute");
    }

    /// @dev Gets a saved attestation for the given snapshot root.
    /// Will return an empty struct, if the snapshot root hasn't been previously saved.
    function _getRootAttestation(bytes32 root) internal view returns (ExecutionAttestation memory) {
        return _rootAttestations[root];
    }
}
