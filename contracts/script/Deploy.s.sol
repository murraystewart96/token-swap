// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {METoken} from "../src/METoken.sol";
import {YOUToken} from "../src/YOUToken.sol";
import {MEYOUPool} from "../src/MEYOUPool.sol";

contract DeployScript is Script {
    function run() public {
        vm.startBroadcast();
        
        // 1. Deploy tokens
        METoken meToken = new METoken(1000000);
        YOUToken youToken = new YOUToken(100000);
        
        console.log("METoken deployed at: ", address(meToken));
        console.log("YOUToken deployed at: ", address(youToken));
        
        // 2. Deploy pool
        MEYOUPool pool = new MEYOUPool(address(meToken), address(youToken));
        console.log("Pool deployed at:", address(pool));
        
        // 3. Add initial liquidity (500,000 MET + 50,000 YOU)
        meToken.approve(address(pool), 500000);
        youToken.approve(address(pool), 50000);
        pool.addLiquidity(500000, 50000);
        console.log("Initial liquidity added: 500,000 MET + 50,000 YOU");
        
        // 4. Perform test swaps
        console.log("Performing test swaps...");
        
        // Swap 1: 1000 MET for YOU
        uint expectedOut1 = pool.getAmountOut(1000, MEYOUPool.TokenType.ME_TOKEN);
        meToken.approve(address(pool), 1000);
        pool.swapMeTokenForYouToken(1000, expectedOut1);
        console.log("Swapped 1000 MET for", expectedOut1, "YOU");
        
        // Swap 2: 500 YOU for MET  
        uint expectedOut2 = pool.getAmountOut(500, MEYOUPool.TokenType.YOU_TOKEN);
        youToken.approve(address(pool), 500);
        pool.swapYouTokenForMeToken(500, expectedOut2);
        console.log("Swapped 500 YOU for", expectedOut2, "MET");
        
        // Print final reserves
        (uint reserve0, uint reserve1) = pool.getReserves();
        console.log("Final reserves: MET (%d) - YOU (%d)", reserve0, reserve1);
        
        vm.stopBroadcast();
    }
}
