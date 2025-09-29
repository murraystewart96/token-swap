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
	Bin: "0x60c060405234801561000f575f5ffd5b50604051611db8380380611db883398181016040528101906100319190610248565b5f73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361009f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610096906102e0565b60405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361010d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161010490610348565b60405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff160361017b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610172906103b0565b60405180910390fd5b8173ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250508073ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff168152505050506103ce565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610217826101ee565b9050919050565b6102278161020d565b8114610231575f5ffd5b50565b5f815190506102428161021e565b92915050565b5f5f6040838503121561025e5761025d6101ea565b5b5f61026b85828601610234565b925050602061027c85828601610234565b9150509250929050565b5f82825260208201905092915050565b7f496e76616c6964204d4520746f6b656e206164647265737300000000000000005f82015250565b5f6102ca601883610286565b91506102d582610296565b602082019050919050565b5f6020820190508181035f8301526102f7816102be565b9050919050565b7f496e76616c696420594f5520746f6b656e2061646472657373000000000000005f82015250565b5f610332601983610286565b915061033d826102fe565b602082019050919050565b5f6020820190508181035f83015261035f81610326565b9050919050565b7f546f6b656e73206d75737420626520646966666572656e7400000000000000005f82015250565b5f61039a601883610286565b91506103a582610366565b602082019050919050565b5f6020820190508181035f8301526103c78161038e565b9050919050565b60805160a0516119756104435f395f81816104350152818161056d01528181610720015281816108ab015281816109e501528181610bc10152610cfc01525f818161018701528181610397015281816104d20152818161080d0152818161094a01528181610c5f0152610d9701526119755ff3fe608060405234801561000f575f5ffd5b5060043610610086575f3560e01c80636b11ea20116100595780636b11ea20146101135780639c8f9f23146101315780639cd441da1461014d578063f28e80881461016957610086565b806308aa6b0f1461008a5780630902f1ac146100a8578063116bdb44146100c757806343e47db3146100f7575b5f5ffd5b610092610185565b60405161009f919061100f565b60405180910390f35b6100b06101a9565b6040516100be929190611040565b60405180910390f35b6100e160048036038101906100dc91906110b8565b6101b8565b6040516100ee91906110f6565b60405180910390f35b610111600480360381019061010c919061110f565b610300565b005b61011b61071e565b604051610128919061100f565b60405180910390f35b61014b6004803603810190610146919061114d565b610742565b005b6101676004803603810190610162919061110f565b610787565b005b610183600480360381019061017e919061110f565b610b2a565b005b7f000000000000000000000000000000000000000000000000000000000000000081565b5f5f5f54600154915091509091565b5f5f83116101fb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101f2906111d2565b60405180910390fd5b5f5f5f60018111156102105761020f6111f0565b5b846001811115610223576102226111f0565b5b03610236575f5491506001549050610240565b60015491505f5490505b5f8211801561024e57505f81115b61028d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161028490611267565b60405180910390fd5b5f858361029a91906112b2565b86836102a691906112e5565b6102b09190611353565b90508181106102f4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102eb906113cd565b60405180910390fd5b80935050505092915050565b5f61030b835f6101b8565b905081811015610350576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161034790611435565b60405180910390fd5b600154811115610395576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161038c90611267565b60405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166323b872dd3330866040518463ffffffff1660e01b81526004016103f293929190611473565b6020604051808303815f875af115801561040e573d5f5f3e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061043291906114dd565b507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb33836040518363ffffffff1660e01b815260040161048e929190611508565b6020604051808303815f875af11580156104aa573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104ce91906114dd565b505f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401610529919061152f565b602060405180830381865afa158015610544573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610568919061155c565b90505f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016105c4919061152f565b602060405180830381865afa1580156105df573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610603919061155c565b9050845f5461061291906112b2565b8214610653576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161064a906115f7565b60405180910390fd5b826001546106619190611615565b81146106a2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610699906116b8565b60405180910390fd5b6106ac8282610f48565b3373ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822875f5f8860405161070f949392919061170f565b60405180910390a35050505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b5f8111610784576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161077b906117c2565b60405180910390fd5b50565b5f82116107c9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107c090611850565b60405180910390fd5b5f811161080b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610802906118de565b60405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166323b872dd3330856040518463ffffffff1660e01b815260040161086893929190611473565b6020604051808303815f875af1158015610884573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906108a891906114dd565b507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166323b872dd3330846040518463ffffffff1660e01b815260040161090693929190611473565b6020604051808303815f875af1158015610922573d5f5f3e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061094691906114dd565b505f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b81526004016109a1919061152f565b602060405180830381865afa1580156109bc573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906109e0919061155c565b90505f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401610a3c919061152f565b602060405180830381865afa158015610a57573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610a7b919061155c565b9050835f54610a8a91906112b2565b8214610acb576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ac2906115f7565b60405180910390fd5b82600154610ad991906112b2565b8114610b1a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b11906116b8565b60405180910390fd5b610b248282610f48565b50505050565b5f610b368360016101b8565b905081811015610b7b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b7290611435565b60405180910390fd5b5f54811115610bbf576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bb690611267565b60405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166323b872dd3330866040518463ffffffff1660e01b8152600401610c1c93929190611473565b6020604051808303815f875af1158015610c38573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610c5c91906114dd565b507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb33836040518363ffffffff1660e01b8152600401610cb8929190611508565b6020604051808303815f875af1158015610cd4573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610cf891906114dd565b505f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401610d53919061152f565b602060405180830381865afa158015610d6e573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610d92919061155c565b90505f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401610dee919061152f565b602060405180830381865afa158015610e09573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610e2d919061155c565b905084600154610e3d91906112b2565b8214610e7e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e75906116b8565b60405180910390fd5b825f54610e8b9190611615565b8114610ecc576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ec3906115f7565b60405180910390fd5b610ed68183610f48565b3373ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d8225f88875f604051610f3994939291906118fc565b60405180910390a35050505050565b815f81905550806001819055507fcf2aa50876cdfbb541206f89af0ee78d44a2abf8d328e37fa4917f982149848a5f54600154604051610f89929190611040565b60405180910390a15050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f819050919050565b5f610fd7610fd2610fcd84610f95565b610fb4565b610f95565b9050919050565b5f610fe882610fbd565b9050919050565b5f610ff982610fde565b9050919050565b61100981610fef565b82525050565b5f6020820190506110225f830184611000565b92915050565b5f819050919050565b61103a81611028565b82525050565b5f6040820190506110535f830185611031565b6110606020830184611031565b9392505050565b5f5ffd5b61107481611028565b811461107e575f5ffd5b50565b5f8135905061108f8161106b565b92915050565b600281106110a1575f5ffd5b50565b5f813590506110b281611095565b92915050565b5f5f604083850312156110ce576110cd611067565b5b5f6110db85828601611081565b92505060206110ec858286016110a4565b9150509250929050565b5f6020820190506111095f830184611031565b92915050565b5f5f6040838503121561112557611124611067565b5b5f61113285828601611081565b925050602061114385828601611081565b9150509250929050565b5f6020828403121561116257611161611067565b5b5f61116f84828501611081565b91505092915050565b5f82825260208201905092915050565b7f416d6f756e7420696e206d7573742062652067726561746572207468616e20305f82015250565b5f6111bc602083611178565b91506111c782611188565b602082019050919050565b5f6020820190508181035f8301526111e9816111b0565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b7f496e73756666696369656e74206c6971756964697479000000000000000000005f82015250565b5f611251601683611178565b915061125c8261121d565b602082019050919050565b5f6020820190508181035f83015261127e81611245565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6112bc82611028565b91506112c783611028565b92508282019050808211156112df576112de611285565b5b92915050565b5f6112ef82611028565b91506112fa83611028565b925082820261130881611028565b9150828204841483151761131f5761131e611285565b5b5092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f61135d82611028565b915061136883611028565b92508261137857611377611326565b5b828204905092915050565b7f496e73756666696369656e74206c697175696469747920666f722073776170005f82015250565b5f6113b7601f83611178565b91506113c282611383565b602082019050919050565b5f6020820190508181035f8301526113e4816113ab565b9050919050565b7f496e73756666696369656e74206f757470757420616d6f756e740000000000005f82015250565b5f61141f601a83611178565b915061142a826113eb565b602082019050919050565b5f6020820190508181035f83015261144c81611413565b9050919050565b5f61145d82610f95565b9050919050565b61146d81611453565b82525050565b5f6060820190506114865f830186611464565b6114936020830185611464565b6114a06040830184611031565b949350505050565b5f8115159050919050565b6114bc816114a8565b81146114c6575f5ffd5b50565b5f815190506114d7816114b3565b92915050565b5f602082840312156114f2576114f1611067565b5b5f6114ff848285016114c9565b91505092915050565b5f60408201905061151b5f830185611464565b6115286020830184611031565b9392505050565b5f6020820190506115425f830184611464565b92915050565b5f815190506115568161106b565b92915050565b5f6020828403121561157157611570611067565b5b5f61157e84828501611548565b91505092915050565b7f556e6578706563746564204d4520746f6b656e20726573657276652062616c615f8201527f6e63650000000000000000000000000000000000000000000000000000000000602082015250565b5f6115e1602383611178565b91506115ec82611587565b604082019050919050565b5f6020820190508181035f83015261160e816115d5565b9050919050565b5f61161f82611028565b915061162a83611028565b925082820390508181111561164257611641611285565b5b92915050565b7f556e657870656374656420594f5520746f6b656e20726573657276652062616c5f8201527f616e636500000000000000000000000000000000000000000000000000000000602082015250565b5f6116a2602483611178565b91506116ad82611648565b604082019050919050565b5f6020820190508181035f8301526116cf81611696565b9050919050565b5f819050919050565b5f6116f96116f46116ef846116d6565b610fb4565b611028565b9050919050565b611709816116df565b82525050565b5f6080820190506117225f830187611031565b61172f6020830186611700565b61173c6040830185611700565b6117496060830184611031565b95945050505050565b7f4c5020746f6b656e20616d6f756e74206d7573742062652067726561746572205f8201527f7468616e20300000000000000000000000000000000000000000000000000000602082015250565b5f6117ac602683611178565b91506117b782611752565b604082019050919050565b5f6020820190508181035f8301526117d9816117a0565b9050919050565b7f4d4520746f6b656e20616d6f756e74206d7573742062652067726561746572205f8201527f7468616e20300000000000000000000000000000000000000000000000000000602082015250565b5f61183a602683611178565b9150611845826117e0565b604082019050919050565b5f6020820190508181035f8301526118678161182e565b9050919050565b7f594f5520746f6b656e20616d6f756e74206d75737420626520677265617465725f8201527f207468616e203000000000000000000000000000000000000000000000000000602082015250565b5f6118c8602783611178565b91506118d38261186e565b604082019050919050565b5f6020820190508181035f8301526118f5816118bc565b9050919050565b5f60808201905061190f5f830187611700565b61191c6020830186611031565b6119296040830185611031565b6119366060830184611700565b9594505050505056fea26469706673582212202092b5370942ed87613aa37c18cc693aec5781eb663417b16994c32b6bb9c89a64736f6c634300081d0033",
}

// PoolABI is the input ABI used to generate the binding from.
// Deprecated: Use PoolMetaData.ABI instead.
var PoolABI = PoolMetaData.ABI

// PoolBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PoolMetaData.Bin instead.
var PoolBin = PoolMetaData.Bin

// DeployPool deploys a new Ethereum contract, binding an instance of Pool to it.
func DeployPool(auth *bind.TransactOpts, backend bind.ContractBackend, _meTokenAddr common.Address, _youTokenAddr common.Address) (common.Address, *types.Transaction, *Pool, error) {
	parsed, err := PoolMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PoolBin), backend, _meTokenAddr, _youTokenAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Pool{PoolCaller: PoolCaller{contract: contract}, PoolTransactor: PoolTransactor{contract: contract}, PoolFilterer: PoolFilterer{contract: contract}}, nil
}

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
