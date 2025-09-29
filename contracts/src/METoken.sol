// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {ERC20} from "openzeppelin-contracts/contracts/token/ERC20/ERC20.sol";

contract METoken is ERC20 {
    constructor(uint256 initialSupply) ERC20("METoken", "MET") {
        _mint(msg.sender, initialSupply);
    }
}
