// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import {IERC20} from "openzeppelin-contracts/contracts/token/ERC20/ERC20.sol";

contract MEYOUPool {
    event Swap(
        address indexed sender, 
        uint256 meTokenIn, 
        uint256 youTokenIn, 
        uint256 meTokenOut, 
        uint256 youTokenOut, 
        address indexed to
    );
    event Sync(uint256 meTokenAmount, uint256 youTokenAmount);

    IERC20 public immutable meToken;
    IERC20 public immutable youToken;

    uint256 private _meTokenReserve;
    uint256 private _youTokenReserve;

    constructor(address _meTokenAddr, address _youTokenAddr) {
        require(_meTokenAddr != address(0), "Invalid ME token address");
        require(_youTokenAddr != address(0), "Invalid YOU token address");
        require(_meTokenAddr != _youTokenAddr, "Tokens must be different");
        
        meToken = IERC20(_meTokenAddr);
        youToken = IERC20(_youTokenAddr);
    }

    function addLiquidity(uint256 _amountMeToken, uint256 _amountYouToken) external {
        require(_amountMeToken > 0, "ME token amount must be greater than 0");
        require(_amountYouToken > 0, "YOU token amount must be greater than 0");

        meToken.transferFrom(msg.sender, address(this), _amountMeToken);
        youToken.transferFrom(msg.sender, address(this), _amountYouToken);

        uint256 balanceMeToken = meToken.balanceOf(address(this));
        uint256 balanceYouToken = youToken.balanceOf(address(this));

        require(balanceMeToken == _meTokenReserve + _amountMeToken, "Unexpected ME token reserve balance");
        require(balanceYouToken == _youTokenReserve + _amountYouToken, "Unexpected YOU token reserve balance");

        _sync(balanceMeToken, balanceYouToken);
    }

    function removeLiquidity(uint256 _lpTokenAmount) external {
        // Implementation pending
        require(_lpTokenAmount > 0, "LP token amount must be greater than 0");
        // TODO: Implement liquidity removal logic
    }

    function swapMeTokenForYouToken(uint256 _amountIn, uint256 _minAmountOut) external {
        uint256 amountOut = getAmountOut(_amountIn, TokenType.ME_TOKEN);
    
        require(amountOut >= _minAmountOut, "Insufficient output amount");
        require(amountOut <= _youTokenReserve, "Insufficient liquidity"); 

        // Execute swap
        meToken.transferFrom(msg.sender, address(this), _amountIn);
        youToken.transfer(msg.sender, amountOut);

        // Verify and sync balances
        uint256 balanceMeToken = meToken.balanceOf(address(this));
        uint256 balanceYouToken = youToken.balanceOf(address(this));

        require(balanceMeToken == _meTokenReserve + _amountIn, "Unexpected ME token reserve balance");
        require(balanceYouToken == _youTokenReserve - amountOut, "Unexpected YOU token reserve balance");

        _sync(balanceMeToken, balanceYouToken);

        emit Swap(msg.sender, _amountIn, 0, 0, amountOut, msg.sender);
    }

    function swapYouTokenForMeToken(uint256 _amountIn, uint256 _minAmountOut) external {
        uint256 amountOut = getAmountOut(_amountIn, TokenType.YOU_TOKEN);
    
        require(amountOut >= _minAmountOut, "Insufficient output amount");
        require(amountOut <= _meTokenReserve, "Insufficient liquidity"); 

        // Execute swap
        youToken.transferFrom(msg.sender, address(this), _amountIn);
        meToken.transfer(msg.sender, amountOut);

        // Verify and sync balances
        uint256 balanceYouToken = youToken.balanceOf(address(this));
        uint256 balanceMeToken = meToken.balanceOf(address(this));

        require(balanceYouToken == _youTokenReserve + _amountIn, "Unexpected YOU token reserve balance");
        require(balanceMeToken == _meTokenReserve - amountOut, "Unexpected ME token reserve balance");

        _sync(balanceMeToken, balanceYouToken);

        emit Swap(msg.sender, 0, _amountIn, amountOut, 0, msg.sender);
    }

    function getReserves() external view returns (uint256 meTokenReserve, uint256 youTokenReserve) {
        return (_meTokenReserve, _youTokenReserve);
    }

    enum TokenType { ME_TOKEN, YOU_TOKEN }

    function getAmountOut(uint256 _amountIn, TokenType _tokenType) public view returns (uint256) {
        require(_amountIn > 0, "Amount in must be greater than 0");
        
        uint256 reserveIn;
        uint256 reserveOut;
        
        if (_tokenType == TokenType.ME_TOKEN) {
            reserveIn = _meTokenReserve;
            reserveOut = _youTokenReserve;
        } else {
            reserveIn = _youTokenReserve;
            reserveOut = _meTokenReserve;
        }
        
        require(reserveIn > 0 && reserveOut > 0, "Insufficient liquidity");
        
        // Using the constant product formula: x * y = k
        // amountOut = (reserveOut * amountIn) / (reserveIn + amountIn)
        uint256 amountOut = (reserveOut * _amountIn) / (reserveIn + _amountIn);
        
        require(amountOut < reserveOut, "Insufficient liquidity for swap");
        
        return amountOut;
    }

    function _sync(uint256 _balanceMeToken, uint256 _balanceYouToken) private {
        _meTokenReserve = _balanceMeToken;
        _youTokenReserve = _balanceYouToken;

        emit Sync(_meTokenReserve, _youTokenReserve);
    }
}
