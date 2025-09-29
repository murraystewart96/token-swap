// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {ERC20} from "openzeppelin-contracts/contracts/token/ERC20/ERC20.sol";

contract YOUToken is ERC20 {
    constructor(uint256 initialSupply) ERC20("YOUToken", "YOU") {
        _mint(msg.sender, initialSupply);
    }
}
