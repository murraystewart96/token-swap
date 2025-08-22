// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// PoolMetaData contains all meta data concerning the Pool contract.
var PoolMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_meTokenAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_youTokenAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addLiquidity\",\"inputs\":[{\"name\":\"_amountMeToken\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_amountYouToken\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAmountOut\",\"inputs\":[{\"name\":\"_amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_tokenType\",\"type\":\"uint8\",\"internalType\":\"enumMEYOUPool.TokenType\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReserves\",\"inputs\":[],\"outputs\":[{\"name\":\"meTokenReserve\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"youTokenReserve\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"meToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeLiquidity\",\"inputs\":[{\"name\":\"_lpTokenAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"swapMeTokenForYouToken\",\"inputs\":[{\"name\":\"_amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_minAmountOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"swapYouTokenForMeToken\",\"inputs\":[{\"name\":\"_amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_minAmountOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"youToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Swap\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"meTokenIn\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"youTokenIn\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"meTokenOut\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"youTokenOut\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Sync\",\"inputs\":[{\"name\":\"meTokenAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"youTokenAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
}

// PoolABI is the input ABI used to generate the binding from.
// Deprecated: Use PoolMetaData.ABI instead.
var PoolABI = PoolMetaData.ABI

// Pool is an auto generated Go binding around an Ethereum contract.
type Pool struct {
	PoolCaller     // Read-only binding to the contract
	PoolTransactor // Write-only binding to the contract
	PoolFilterer   // Log filterer for contract events
}

// PoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type PoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PoolSession struct {
	Contract     *Pool             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PoolCallerSession struct {
	Contract *PoolCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// PoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PoolTransactorSession struct {
	Contract     *PoolTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type PoolRaw struct {
	Contract *Pool // Generic contract binding to access the raw methods on
}

// PoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PoolCallerRaw struct {
	Contract *PoolCaller // Generic read-only contract binding to access the raw methods on
}

// PoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PoolTransactorRaw struct {
	Contract *PoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPool creates a new instance of Pool, bound to a specific deployed contract.
func NewPool(address common.Address, backend bind.ContractBackend) (*Pool, error) {
	contract, err := bindPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Pool{PoolCaller: PoolCaller{contract: contract}, PoolTransactor: PoolTransactor{contract: contract}, PoolFilterer: PoolFilterer{contract: contract}}, nil
}

// NewPoolCaller creates a new read-only instance of Pool, bound to a specific deployed contract.
func NewPoolCaller(address common.Address, caller bind.ContractCaller) (*PoolCaller, error) {
	contract, err := bindPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PoolCaller{contract: contract}, nil
}

// NewPoolTransactor creates a new write-only instance of Pool, bound to a specific deployed contract.
func NewPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*PoolTransactor, error) {
	contract, err := bindPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PoolTransactor{contract: contract}, nil
}

// NewPoolFilterer creates a new log filterer instance of Pool, bound to a specific deployed contract.
func NewPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*PoolFilterer, error) {
	contract, err := bindPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PoolFilterer{contract: contract}, nil
}

// bindPool binds a generic wrapper to an already deployed contract.
func bindPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pool *PoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Pool.Contract.PoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pool *PoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pool.Contract.PoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pool *PoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pool.Contract.PoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pool *PoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Pool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pool *PoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pool *PoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pool.Contract.contract.Transact(opts, method, params...)
}

// GetAmountOut is a free data retrieval call binding the contract method 0x116bdb44.
//
// Solidity: function getAmountOut(uint256 _amountIn, uint8 _tokenType) view returns(uint256)
func (_Pool *PoolCaller) GetAmountOut(opts *bind.CallOpts, _amountIn *big.Int, _tokenType uint8) (*big.Int, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "getAmountOut", _amountIn, _tokenType)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAmountOut is a free data retrieval call binding the contract method 0x116bdb44.
//
// Solidity: function getAmountOut(uint256 _amountIn, uint8 _tokenType) view returns(uint256)
func (_Pool *PoolSession) GetAmountOut(_amountIn *big.Int, _tokenType uint8) (*big.Int, error) {
	return _Pool.Contract.GetAmountOut(&_Pool.CallOpts, _amountIn, _tokenType)
}

// GetAmountOut is a free data retrieval call binding the contract method 0x116bdb44.
//
// Solidity: function getAmountOut(uint256 _amountIn, uint8 _tokenType) view returns(uint256)
func (_Pool *PoolCallerSession) GetAmountOut(_amountIn *big.Int, _tokenType uint8) (*big.Int, error) {
	return _Pool.Contract.GetAmountOut(&_Pool.CallOpts, _amountIn, _tokenType)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint256 meTokenReserve, uint256 youTokenReserve)
func (_Pool *PoolCaller) GetReserves(opts *bind.CallOpts) (struct {
	MeTokenReserve  *big.Int
	YouTokenReserve *big.Int
}, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "getReserves")

	outstruct := new(struct {
		MeTokenReserve  *big.Int
		YouTokenReserve *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MeTokenReserve = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.YouTokenReserve = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint256 meTokenReserve, uint256 youTokenReserve)
func (_Pool *PoolSession) GetReserves() (struct {
	MeTokenReserve  *big.Int
	YouTokenReserve *big.Int
}, error) {
	return _Pool.Contract.GetReserves(&_Pool.CallOpts)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint256 meTokenReserve, uint256 youTokenReserve)
func (_Pool *PoolCallerSession) GetReserves() (struct {
	MeTokenReserve  *big.Int
	YouTokenReserve *big.Int
}, error) {
	return _Pool.Contract.GetReserves(&_Pool.CallOpts)
}

// MeToken is a free data retrieval call binding the contract method 0x08aa6b0f.
//
// Solidity: function meToken() view returns(address)
func (_Pool *PoolCaller) MeToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "meToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MeToken is a free data retrieval call binding the contract method 0x08aa6b0f.
//
// Solidity: function meToken() view returns(address)
func (_Pool *PoolSession) MeToken() (common.Address, error) {
	return _Pool.Contract.MeToken(&_Pool.CallOpts)
}

// MeToken is a free data retrieval call binding the contract method 0x08aa6b0f.
//
// Solidity: function meToken() view returns(address)
func (_Pool *PoolCallerSession) MeToken() (common.Address, error) {
	return _Pool.Contract.MeToken(&_Pool.CallOpts)
}

// YouToken is a free data retrieval call binding the contract method 0x6b11ea20.
//
// Solidity: function youToken() view returns(address)
func (_Pool *PoolCaller) YouToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Pool.contract.Call(opts, &out, "youToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// YouToken is a free data retrieval call binding the contract method 0x6b11ea20.
//
// Solidity: function youToken() view returns(address)
func (_Pool *PoolSession) YouToken() (common.Address, error) {
	return _Pool.Contract.YouToken(&_Pool.CallOpts)
}

// YouToken is a free data retrieval call binding the contract method 0x6b11ea20.
//
// Solidity: function youToken() view returns(address)
func (_Pool *PoolCallerSession) YouToken() (common.Address, error) {
	return _Pool.Contract.YouToken(&_Pool.CallOpts)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0x9cd441da.
//
// Solidity: function addLiquidity(uint256 _amountMeToken, uint256 _amountYouToken) returns()
func (_Pool *PoolTransactor) AddLiquidity(opts *bind.TransactOpts, _amountMeToken *big.Int, _amountYouToken *big.Int) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "addLiquidity", _amountMeToken, _amountYouToken)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0x9cd441da.
//
// Solidity: function addLiquidity(uint256 _amountMeToken, uint256 _amountYouToken) returns()
func (_Pool *PoolSession) AddLiquidity(_amountMeToken *big.Int, _amountYouToken *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.AddLiquidity(&_Pool.TransactOpts, _amountMeToken, _amountYouToken)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0x9cd441da.
//
// Solidity: function addLiquidity(uint256 _amountMeToken, uint256 _amountYouToken) returns()
func (_Pool *PoolTransactorSession) AddLiquidity(_amountMeToken *big.Int, _amountYouToken *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.AddLiquidity(&_Pool.TransactOpts, _amountMeToken, _amountYouToken)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0x9c8f9f23.
//
// Solidity: function removeLiquidity(uint256 _lpTokenAmount) returns()
func (_Pool *PoolTransactor) RemoveLiquidity(opts *bind.TransactOpts, _lpTokenAmount *big.Int) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "removeLiquidity", _lpTokenAmount)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0x9c8f9f23.
//
// Solidity: function removeLiquidity(uint256 _lpTokenAmount) returns()
func (_Pool *PoolSession) RemoveLiquidity(_lpTokenAmount *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.RemoveLiquidity(&_Pool.TransactOpts, _lpTokenAmount)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0x9c8f9f23.
//
// Solidity: function removeLiquidity(uint256 _lpTokenAmount) returns()
func (_Pool *PoolTransactorSession) RemoveLiquidity(_lpTokenAmount *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.RemoveLiquidity(&_Pool.TransactOpts, _lpTokenAmount)
}

// SwapMeTokenForYouToken is a paid mutator transaction binding the contract method 0x43e47db3.
//
// Solidity: function swapMeTokenForYouToken(uint256 _amountIn, uint256 _minAmountOut) returns()
func (_Pool *PoolTransactor) SwapMeTokenForYouToken(opts *bind.TransactOpts, _amountIn *big.Int, _minAmountOut *big.Int) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "swapMeTokenForYouToken", _amountIn, _minAmountOut)
}

// SwapMeTokenForYouToken is a paid mutator transaction binding the contract method 0x43e47db3.
//
// Solidity: function swapMeTokenForYouToken(uint256 _amountIn, uint256 _minAmountOut) returns()
func (_Pool *PoolSession) SwapMeTokenForYouToken(_amountIn *big.Int, _minAmountOut *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.SwapMeTokenForYouToken(&_Pool.TransactOpts, _amountIn, _minAmountOut)
}

// SwapMeTokenForYouToken is a paid mutator transaction binding the contract method 0x43e47db3.
//
// Solidity: function swapMeTokenForYouToken(uint256 _amountIn, uint256 _minAmountOut) returns()
func (_Pool *PoolTransactorSession) SwapMeTokenForYouToken(_amountIn *big.Int, _minAmountOut *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.SwapMeTokenForYouToken(&_Pool.TransactOpts, _amountIn, _minAmountOut)
}

// SwapYouTokenForMeToken is a paid mutator transaction binding the contract method 0xf28e8088.
//
// Solidity: function swapYouTokenForMeToken(uint256 _amountIn, uint256 _minAmountOut) returns()
func (_Pool *PoolTransactor) SwapYouTokenForMeToken(opts *bind.TransactOpts, _amountIn *big.Int, _minAmountOut *big.Int) (*types.Transaction, error) {
	return _Pool.contract.Transact(opts, "swapYouTokenForMeToken", _amountIn, _minAmountOut)
}

// SwapYouTokenForMeToken is a paid mutator transaction binding the contract method 0xf28e8088.
//
// Solidity: function swapYouTokenForMeToken(uint256 _amountIn, uint256 _minAmountOut) returns()
func (_Pool *PoolSession) SwapYouTokenForMeToken(_amountIn *big.Int, _minAmountOut *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.SwapYouTokenForMeToken(&_Pool.TransactOpts, _amountIn, _minAmountOut)
}

// SwapYouTokenForMeToken is a paid mutator transaction binding the contract method 0xf28e8088.
//
// Solidity: function swapYouTokenForMeToken(uint256 _amountIn, uint256 _minAmountOut) returns()
func (_Pool *PoolTransactorSession) SwapYouTokenForMeToken(_amountIn *big.Int, _minAmountOut *big.Int) (*types.Transaction, error) {
	return _Pool.Contract.SwapYouTokenForMeToken(&_Pool.TransactOpts, _amountIn, _minAmountOut)
}

// PoolSwapIterator is returned from FilterSwap and is used to iterate over the raw logs and unpacked data for Swap events raised by the Pool contract.
type PoolSwapIterator struct {
	Event *PoolSwap // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolSwapIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolSwap)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolSwap)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolSwapIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolSwapIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolSwap represents a Swap event raised by the Pool contract.
type PoolSwap struct {
	Sender      common.Address
	MeTokenIn   *big.Int
	YouTokenIn  *big.Int
	MeTokenOut  *big.Int
	YouTokenOut *big.Int
	To          common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterSwap is a free log retrieval operation binding the contract event 0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822.
//
// Solidity: event Swap(address indexed sender, uint256 meTokenIn, uint256 youTokenIn, uint256 meTokenOut, uint256 youTokenOut, address indexed to)
func (_Pool *PoolFilterer) FilterSwap(opts *bind.FilterOpts, sender []common.Address, to []common.Address) (*PoolSwapIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Pool.contract.FilterLogs(opts, "Swap", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PoolSwapIterator{contract: _Pool.contract, event: "Swap", logs: logs, sub: sub}, nil
}

// WatchSwap is a free log subscription operation binding the contract event 0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822.
//
// Solidity: event Swap(address indexed sender, uint256 meTokenIn, uint256 youTokenIn, uint256 meTokenOut, uint256 youTokenOut, address indexed to)
func (_Pool *PoolFilterer) WatchSwap(opts *bind.WatchOpts, sink chan<- *PoolSwap, sender []common.Address, to []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Pool.contract.WatchLogs(opts, "Swap", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolSwap)
				if err := _Pool.contract.UnpackLog(event, "Swap", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSwap is a log parse operation binding the contract event 0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822.
//
// Solidity: event Swap(address indexed sender, uint256 meTokenIn, uint256 youTokenIn, uint256 meTokenOut, uint256 youTokenOut, address indexed to)
func (_Pool *PoolFilterer) ParseSwap(log types.Log) (*PoolSwap, error) {
	event := new(PoolSwap)
	if err := _Pool.contract.UnpackLog(event, "Swap", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolSyncIterator is returned from FilterSync and is used to iterate over the raw logs and unpacked data for Sync events raised by the Pool contract.
type PoolSyncIterator struct {
	Event *PoolSync // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PoolSyncIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolSync)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PoolSync)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PoolSyncIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolSyncIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolSync represents a Sync event raised by the Pool contract.
type PoolSync struct {
	MeTokenAmount  *big.Int
	YouTokenAmount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterSync is a free log retrieval operation binding the contract event 0xcf2aa50876cdfbb541206f89af0ee78d44a2abf8d328e37fa4917f982149848a.
//
// Solidity: event Sync(uint256 meTokenAmount, uint256 youTokenAmount)
func (_Pool *PoolFilterer) FilterSync(opts *bind.FilterOpts) (*PoolSyncIterator, error) {

	logs, sub, err := _Pool.contract.FilterLogs(opts, "Sync")
	if err != nil {
		return nil, err
	}
	return &PoolSyncIterator{contract: _Pool.contract, event: "Sync", logs: logs, sub: sub}, nil
}

// WatchSync is a free log subscription operation binding the contract event 0xcf2aa50876cdfbb541206f89af0ee78d44a2abf8d328e37fa4917f982149848a.
//
// Solidity: event Sync(uint256 meTokenAmount, uint256 youTokenAmount)
func (_Pool *PoolFilterer) WatchSync(opts *bind.WatchOpts, sink chan<- *PoolSync) (event.Subscription, error) {

	logs, sub, err := _Pool.contract.WatchLogs(opts, "Sync")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolSync)
				if err := _Pool.contract.UnpackLog(event, "Sync", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSync is a log parse operation binding the contract event 0xcf2aa50876cdfbb541206f89af0ee78d44a2abf8d328e37fa4917f982149848a.
//
// Solidity: event Sync(uint256 meTokenAmount, uint256 youTokenAmount)
func (_Pool *PoolFilterer) ParseSync(log types.Log) (*PoolSync, error) {
	event := new(PoolSync)
	if err := _Pool.contract.UnpackLog(event, "Sync", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
