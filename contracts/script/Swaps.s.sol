// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script, console} from "forge-std/Script.sol";
import {METoken} from "../src/METoken.sol";
import {YOUToken} from "../src/YOUToken.sol";
import {MEYOUPool} from "../src/MEYOUPool.sol";

contract SwapTestScript is Script {
    // Update these addresses with your deployed contract addresses
    address constant METOKEN_ADDRESS = 0x5FbDB2315678afecb367f032d93F642f64180aa3;
    address constant YOUTOKEN_ADDRESS = 0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512;
    address constant POOL_ADDRESS = 0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0;
    
    function run() public {
        vm.startBroadcast();
        
        METoken meToken = METoken(METOKEN_ADDRESS);
        YOUToken youToken = YOUToken(YOUTOKEN_ADDRESS);
        MEYOUPool pool = MEYOUPool(POOL_ADDRESS);
        
        console.log("Starting swap test with multiple transactions...");
        
        // Perform 10 swaps alternating between directions
        for (uint i = 0; i < 10; i++) {
            if (i % 2 == 0) {
                // Swap MET for YOU
                uint swapAmount = 100 + (i * 50); // Varying amounts
                uint expectedOut = pool.getAmountOut(swapAmount, MEYOUPool.TokenType.ME_TOKEN);
                
                meToken.approve(address(pool), swapAmount);
                pool.swapMeTokenForYouToken(swapAmount, expectedOut);
                
                console.log("Swap %d: %d MET -> %d YOU", i + 1, swapAmount, expectedOut);
            } else {
                // Swap YOU for MET
                uint swapAmount = 50 + (i * 25); // Varying amounts
                uint expectedOut = pool.getAmountOut(swapAmount, MEYOUPool.TokenType.YOU_TOKEN);
                
                youToken.approve(address(pool), swapAmount);
                pool.swapYouTokenForMeToken(swapAmount, expectedOut);
                
                console.log("Swap %d: %d YOU -> %d MET", i + 1, swapAmount, expectedOut);
            }
            
            // Small delay between swaps
            vm.sleep(100);
        }
        
        // Print final reserves
        (uint reserve0, uint reserve1) = pool.getReserves();
        console.log("Final reserves - MET:", reserve0, "YOU:", reserve1);

        vm.stopBroadcast();
    }
}
