// SPDX-License-Identifier: MIT
pragma solidity 0.8.17;

import { BondingSecondary } from "../contracts/bonding/BondingSecondary.sol";
// ═════════════════════════════ INTERNAL IMPORTS ══════════════════════════════
import { DeployerUtils } from "./utils/DeployerUtils.sol";
// ═════════════════════════════ EXTERNAL IMPORTS ══════════════════════════════
import { console, stdJson } from "forge-std/Script.sol";
import { Strings } from "@openzeppelin/contracts/utils/Strings.sol";

contract RemoveAgents002Script is DeployerUtils {
    using stdJson for string;
    using Strings for uint256;

    string public constant MESSAGING_002 = "Messaging002Testnet";

    constructor() {
        setupPK("MESSAGING_DEPLOYER_PRIVATE_KEY");
    }

    /// @dev Function to exclude script from coverage report
    function testScript() external {}

    function run(address manager) external {
        startBroadcast(true);
        BondingSecondary bondingManager = BondingSecondary(manager);
        string memory config = loadGlobalDeployConfig(MESSAGING_002);
        console.log("Removing Agents");
        uint256[] memory domains = config.readUintArray(".domains");
        for (uint256 i = 0; i < domains.length; ++i) {
            uint256 domain = domains[i];
            // Key is ".agents.0: for Guards, ".agents.10" for Optimism Notaries, etc
            address[] memory agents = config.readAddressArray(
                string.concat(".agents.", domain.toString())
            );
            for (uint256 j = 0; j < agents.length; ++j) {
                require(bondingManager.isActiveAgent(uint32(domain), agents[j]), "Not an agent");
                bondingManager.removeAgent(uint32(domain), agents[j]);
                require(
                    !bondingManager.isActiveAgent(uint32(domain), agents[j]),
                    "Failed to remove agent"
                );
                console.log("   %s on domain [%s]", agents[j], domain);
            }
        }
        stopBroadcast();
    }
}
