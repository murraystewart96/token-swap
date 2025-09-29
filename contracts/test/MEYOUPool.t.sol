// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test} from "forge-std/Test.sol";

import {METoken} from "../src/METoken.sol";
import {YOUToken} from "../src/YOUToken.sol";
import {MEYOUPool} from "../src/MEYOUPool.sol";

contract MEYOUPoolTest is Test {
    METoken meToken;
    YOUToken youToken;
    MEYOUPool pool;

    address alice = address(0x1);
    address bob = address(0x2);
    address charlie = address(0x3);
    
    function setUp() public {
        // Deploy tokens
        meToken = new METoken(10000000);
        youToken = new YOUToken(100000);
        
        // Deploy pool
        pool = new MEYOUPool(address(meToken), address(youToken));
        
        // Transfer tokens to test users
        meToken.transfer(alice, 1000);
        youToken.transfer(alice, 500);
        
        meToken.transfer(bob, 2000);
        youToken.transfer(bob, 300);

        meToken.transfer(charlie, 500);
        youToken.transfer(charlie, 100);
    }
    
    function testAddLiquidity() public {
        addLiquidity(address(alice), 1000, 100);

        (uint256 totalMeToken, uint256 totalYouToken) = pool.getReserves();
        
        assertEq(totalMeToken, 1000);
        assertEq(totalYouToken, 100);
    }

    function testSwapMeTokenForYouToken() public {
        addInitialLiquidity();

        vm.startPrank(address(bob));

        uint256 beforeBalance = youToken.balanceOf(address(bob));

        uint256 expectedReturn = pool.getAmountOut(500, MEYOUPool.TokenType.ME_TOKEN);

        meToken.approve(address(pool), 500);
        pool.swapMeTokenForYouToken(500, expectedReturn);

        uint256 afterBalance = youToken.balanceOf(address(bob));

        assertEq(afterBalance, beforeBalance + expectedReturn);

        vm.stopPrank();
    }

    function testSwapYouTokenForMeToken() public {
        addInitialLiquidity();

        vm.startPrank(address(charlie));

        uint256 beforeBalance = meToken.balanceOf(address(charlie));

        uint256 expectedReturn = pool.getAmountOut(50, MEYOUPool.TokenType.YOU_TOKEN);

        youToken.approve(address(pool), 50);
        pool.swapYouTokenForMeToken(50, expectedReturn);

        uint256 afterBalance = meToken.balanceOf(address(charlie));

        assertEq(afterBalance, beforeBalance + expectedReturn);

        vm.stopPrank();
    }

    function testAddLiquidityWithZeroAmount() public {
        vm.startPrank(address(alice));
        
        meToken.approve(address(pool), 1000);
        youToken.approve(address(pool), 100);
        
        vm.expectRevert("ME token amount must be greater than 0");
        pool.addLiquidity(0, 100);
        
        vm.expectRevert("YOU token amount must be greater than 0");
        pool.addLiquidity(1000, 0);
        
        vm.stopPrank();
    }

    function testSwapWithInsufficientOutput() public {
        addInitialLiquidity();

        vm.startPrank(address(bob));

        uint256 expectedReturn = pool.getAmountOut(500, MEYOUPool.TokenType.ME_TOKEN);

        meToken.approve(address(pool), 500);
        
        vm.expectRevert("Insufficient output amount");
        pool.swapMeTokenForYouToken(500, expectedReturn + 1);

        vm.stopPrank();
    }

    function testGetAmountOutWithZeroInput() public {
        addInitialLiquidity();

        vm.expectRevert("Amount in must be greater than 0");
        pool.getAmountOut(0, MEYOUPool.TokenType.ME_TOKEN);
    }

    function testGetAmountOutWithNoLiquidity() public {
        vm.expectRevert("Insufficient liquidity");
        pool.getAmountOut(100, MEYOUPool.TokenType.ME_TOKEN);
    }

    function testConstructorWithInvalidAddresses() public {
        vm.expectRevert("Invalid ME token address");
        new MEYOUPool(address(0), address(youToken));

        vm.expectRevert("Invalid YOU token address");
        new MEYOUPool(address(meToken), address(0));

        vm.expectRevert("Tokens must be different");
        new MEYOUPool(address(meToken), address(meToken));
    }

    // function testEventEmission() public {
    //     // Test Sync event on add liquidity
    //     vm.expectEmit(true, true, true, true);
    //     emit MEYOUPool.Sync(1000, 100);
        
    //     addLiquidity(address(alice), 1000, 100);

    //     // Test Swap event
    //     addInitialLiquidity();
        
    //     vm.startPrank(address(bob));
        
    //     uint256 expectedReturn = pool.getAmountOut(500, MEYOUPool.TokenType.ME_TOKEN);
        
    //     vm.expectEmit(true, true, true, true);
    //     emit MEYOUPool.Swap(address(bob), 500, 0, 0, expectedReturn, address(bob));
        
    //     meToken.approve(address(pool), 500);
    //     pool.swapMeTokenForYouToken(500, expectedReturn);
        
    //     vm.stopPrank();
    // }

    // HELPERS

    function addLiquidity(address spender, uint256 metAmount, uint256 youAmount) private {
        vm.startPrank(spender);

        meToken.approve(address(pool), metAmount);
        youToken.approve(address(pool), youAmount);
        pool.addLiquidity(metAmount, youAmount);

        vm.stopPrank();
    }

    function addInitialLiquidity() private {
        meToken.approve(address(pool), 100000);
        youToken.approve(address(pool), 10000);
        pool.addLiquidity(100000, 10000);
    }
}
