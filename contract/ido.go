// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// IDOABI is the input ABI used to generate the binding from.
const IDOABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Claim\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Invest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"_invitation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_nft\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_partnerCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_sat\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_stakeCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_stop\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_totalAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_totalCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_totalCount2\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_usdt\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"invest1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"invest2\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sat\",\"type\":\"address\"}],\"name\":\"setSatAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"s\",\"type\":\"bool\"}],\"name\":\"setStop\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stakeUsdt\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"userInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"invest\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"invitees\",\"type\":\"address[]\"},{\"internalType\":\"bool\",\"name\":\"isPartner\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isClaimed\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"totalAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalCount2\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stakeCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"partnerCount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// IDO is an auto generated Go binding around an Ethereum contract.
type IDO struct {
	IDOCaller     // Read-only binding to the contract
	IDOTransactor // Write-only binding to the contract
	IDOFilterer   // Log filterer for contract events
}

// IDOCaller is an auto generated read-only Go binding around an Ethereum contract.
type IDOCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IDOTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IDOTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IDOFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IDOFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IDOSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IDOSession struct {
	Contract     *IDO              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IDOCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IDOCallerSession struct {
	Contract *IDOCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// IDOTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IDOTransactorSession struct {
	Contract     *IDOTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IDORaw is an auto generated low-level Go binding around an Ethereum contract.
type IDORaw struct {
	Contract *IDO // Generic contract binding to access the raw methods on
}

// IDOCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IDOCallerRaw struct {
	Contract *IDOCaller // Generic read-only contract binding to access the raw methods on
}

// IDOTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IDOTransactorRaw struct {
	Contract *IDOTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIDO creates a new instance of IDO, bound to a specific deployed contract.
func NewIDO(address common.Address, backend bind.ContractBackend) (*IDO, error) {
	contract, err := bindIDO(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IDO{IDOCaller: IDOCaller{contract: contract}, IDOTransactor: IDOTransactor{contract: contract}, IDOFilterer: IDOFilterer{contract: contract}}, nil
}

// NewIDOCaller creates a new read-only instance of IDO, bound to a specific deployed contract.
func NewIDOCaller(address common.Address, caller bind.ContractCaller) (*IDOCaller, error) {
	contract, err := bindIDO(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IDOCaller{contract: contract}, nil
}

// NewIDOTransactor creates a new write-only instance of IDO, bound to a specific deployed contract.
func NewIDOTransactor(address common.Address, transactor bind.ContractTransactor) (*IDOTransactor, error) {
	contract, err := bindIDO(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IDOTransactor{contract: contract}, nil
}

// NewIDOFilterer creates a new log filterer instance of IDO, bound to a specific deployed contract.
func NewIDOFilterer(address common.Address, filterer bind.ContractFilterer) (*IDOFilterer, error) {
	contract, err := bindIDO(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IDOFilterer{contract: contract}, nil
}

// bindIDO binds a generic wrapper to an already deployed contract.
func bindIDO(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(IDOABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IDO *IDORaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IDO.Contract.IDOCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IDO *IDORaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IDO.Contract.IDOTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IDO *IDORaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IDO.Contract.IDOTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IDO *IDOCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IDO.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IDO *IDOTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IDO.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IDO *IDOTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IDO.Contract.contract.Transact(opts, method, params...)
}

// Invitation is a free data retrieval call binding the contract method 0xbe086508.
//
// Solidity: function _invitation() view returns(address)
func (_IDO *IDOCaller) Invitation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IDO.contract.Call(opts, &out, "_invitation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Invitation is a free data retrieval call binding the contract method 0xbe086508.
//
// Solidity: function _invitation() view returns(address)
func (_IDO *IDOSession) Invitation() (common.Address, error) {
	return _IDO.Contract.Invitation(&_IDO.CallOpts)
}

// Invitation is a free data retrieval call binding the contract method 0xbe086508.
//
// Solidity: function _invitation() view returns(address)
func (_IDO *IDOCallerSession) Invitation() (common.Address, error) {
	return _IDO.Contract.Invitation(&_IDO.CallOpts)
}

// Nft is a free data retrieval call binding the contract method 0x98300e18.
//
// Solidity: function _nft() view returns(address)
func (_IDO *IDOCaller) Nft(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IDO.contract.Call(opts, &out, "_nft")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Nft is a free data retrieval call binding the contract method 0x98300e18.
//
// Solidity: function _nft() view returns(address)
func (_IDO *IDOSession) Nft() (common.Address, error) {
	return _IDO.Contract.Nft(&_IDO.CallOpts)
}

// Nft is a free data retrieval call binding the contract method 0x98300e18.
//
// Solidity: function _nft() view returns(address)
func (_IDO *IDOCallerSession) Nft() (common.Address, error) {
	return _IDO.Contract.Nft(&_IDO.CallOpts)
}

// PartnerCount is a free data retrieval call binding the contract method 0x3ced72ac.
//
// Solidity: function _partnerCount() view returns(uint256)
func (_IDO *IDOCaller) PartnerCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IDO.contract.Call(opts, &out, "_partnerCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PartnerCount is a free data retrieval call binding the contract method 0x3ced72ac.
//
// Solidity: function _partnerCount() view returns(uint256)
func (_IDO *IDOSession) PartnerCount() (*big.Int, error) {
	return _IDO.Contract.PartnerCount(&_IDO.CallOpts)
}

// PartnerCount is a free data retrieval call binding the contract method 0x3ced72ac.
//
// Solidity: function _partnerCount() view returns(uint256)
func (_IDO *IDOCallerSession) PartnerCount() (*big.Int, error) {
	return _IDO.Contract.PartnerCount(&_IDO.CallOpts)
}

// Sat is a free data retrieval call binding the contract method 0x8983e08b.
//
// Solidity: function _sat() view returns(address)
func (_IDO *IDOCaller) Sat(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IDO.contract.Call(opts, &out, "_sat")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Sat is a free data retrieval call binding the contract method 0x8983e08b.
//
// Solidity: function _sat() view returns(address)
func (_IDO *IDOSession) Sat() (common.Address, error) {
	return _IDO.Contract.Sat(&_IDO.CallOpts)
}

// Sat is a free data retrieval call binding the contract method 0x8983e08b.
//
// Solidity: function _sat() view returns(address)
func (_IDO *IDOCallerSession) Sat() (common.Address, error) {
	return _IDO.Contract.Sat(&_IDO.CallOpts)
}

// StakeCount is a free data retrieval call binding the contract method 0xfb12a6ae.
//
// Solidity: function _stakeCount() view returns(uint256)
func (_IDO *IDOCaller) StakeCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IDO.contract.Call(opts, &out, "_stakeCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// StakeCount is a free data retrieval call binding the contract method 0xfb12a6ae.
//
// Solidity: function _stakeCount() view returns(uint256)
func (_IDO *IDOSession) StakeCount() (*big.Int, error) {
	return _IDO.Contract.StakeCount(&_IDO.CallOpts)
}

// StakeCount is a free data retrieval call binding the contract method 0xfb12a6ae.
//
// Solidity: function _stakeCount() view returns(uint256)
func (_IDO *IDOCallerSession) StakeCount() (*big.Int, error) {
	return _IDO.Contract.StakeCount(&_IDO.CallOpts)
}

// Stop is a free data retrieval call binding the contract method 0xbe8343ab.
//
// Solidity: function _stop() view returns(bool)
func (_IDO *IDOCaller) Stop(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _IDO.contract.Call(opts, &out, "_stop")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Stop is a free data retrieval call binding the contract method 0xbe8343ab.
//
// Solidity: function _stop() view returns(bool)
func (_IDO *IDOSession) Stop() (bool, error) {
	return _IDO.Contract.Stop(&_IDO.CallOpts)
}

// Stop is a free data retrieval call binding the contract method 0xbe8343ab.
//
// Solidity: function _stop() view returns(bool)
func (_IDO *IDOCallerSession) Stop() (bool, error) {
	return _IDO.Contract.Stop(&_IDO.CallOpts)
}

// TotalAmount is a free data retrieval call binding the contract method 0x3bbeaab5.
//
// Solidity: function _totalAmount() view returns(uint256)
func (_IDO *IDOCaller) TotalAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IDO.contract.Call(opts, &out, "_totalAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalAmount is a free data retrieval call binding the contract method 0x3bbeaab5.
//
// Solidity: function _totalAmount() view returns(uint256)
func (_IDO *IDOSession) TotalAmount() (*big.Int, error) {
	return _IDO.Contract.TotalAmount(&_IDO.CallOpts)
}

// TotalAmount is a free data retrieval call binding the contract method 0x3bbeaab5.
//
// Solidity: function _totalAmount() view returns(uint256)
func (_IDO *IDOCallerSession) TotalAmount() (*big.Int, error) {
	return _IDO.Contract.TotalAmount(&_IDO.CallOpts)
}

// TotalCount is a free data retrieval call binding the contract method 0xeab26b91.
//
// Solidity: function _totalCount() view returns(uint256)
func (_IDO *IDOCaller) TotalCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IDO.contract.Call(opts, &out, "_totalCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalCount is a free data retrieval call binding the contract method 0xeab26b91.
//
// Solidity: function _totalCount() view returns(uint256)
func (_IDO *IDOSession) TotalCount() (*big.Int, error) {
	return _IDO.Contract.TotalCount(&_IDO.CallOpts)
}

// TotalCount is a free data retrieval call binding the contract method 0xeab26b91.
//
// Solidity: function _totalCount() view returns(uint256)
func (_IDO *IDOCallerSession) TotalCount() (*big.Int, error) {
	return _IDO.Contract.TotalCount(&_IDO.CallOpts)
}

// TotalCount2 is a free data retrieval call binding the contract method 0xde90eed0.
//
// Solidity: function _totalCount2() view returns(uint256)
func (_IDO *IDOCaller) TotalCount2(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IDO.contract.Call(opts, &out, "_totalCount2")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalCount2 is a free data retrieval call binding the contract method 0xde90eed0.
//
// Solidity: function _totalCount2() view returns(uint256)
func (_IDO *IDOSession) TotalCount2() (*big.Int, error) {
	return _IDO.Contract.TotalCount2(&_IDO.CallOpts)
}

// TotalCount2 is a free data retrieval call binding the contract method 0xde90eed0.
//
// Solidity: function _totalCount2() view returns(uint256)
func (_IDO *IDOCallerSession) TotalCount2() (*big.Int, error) {
	return _IDO.Contract.TotalCount2(&_IDO.CallOpts)
}

// Usdt is a free data retrieval call binding the contract method 0xbe3601f8.
//
// Solidity: function _usdt() view returns(address)
func (_IDO *IDOCaller) Usdt(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IDO.contract.Call(opts, &out, "_usdt")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Usdt is a free data retrieval call binding the contract method 0xbe3601f8.
//
// Solidity: function _usdt() view returns(address)
func (_IDO *IDOSession) Usdt() (common.Address, error) {
	return _IDO.Contract.Usdt(&_IDO.CallOpts)
}

// Usdt is a free data retrieval call binding the contract method 0xbe3601f8.
//
// Solidity: function _usdt() view returns(address)
func (_IDO *IDOCallerSession) Usdt() (common.Address, error) {
	return _IDO.Contract.Usdt(&_IDO.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IDO *IDOCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IDO.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IDO *IDOSession) Owner() (common.Address, error) {
	return _IDO.Contract.Owner(&_IDO.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_IDO *IDOCallerSession) Owner() (common.Address, error) {
	return _IDO.Contract.Owner(&_IDO.CallOpts)
}

// UserInfo is a free data retrieval call binding the contract method 0x1959a002.
//
// Solidity: function userInfo(address account) view returns(uint256 invest, address[] invitees, bool isPartner, bool isClaimed, uint256 totalAmount, uint256 totalCount, uint256 totalCount2, uint256 stakeCount, uint256 partnerCount)
func (_IDO *IDOCaller) UserInfo(opts *bind.CallOpts, account common.Address) (struct {
	Invest       *big.Int
	Invitees     []common.Address
	IsPartner    bool
	IsClaimed    bool
	TotalAmount  *big.Int
	TotalCount   *big.Int
	TotalCount2  *big.Int
	StakeCount   *big.Int
	PartnerCount *big.Int
}, error) {
	var out []interface{}
	err := _IDO.contract.Call(opts, &out, "userInfo", account)

	outstruct := new(struct {
		Invest       *big.Int
		Invitees     []common.Address
		IsPartner    bool
		IsClaimed    bool
		TotalAmount  *big.Int
		TotalCount   *big.Int
		TotalCount2  *big.Int
		StakeCount   *big.Int
		PartnerCount *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Invest = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Invitees = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)
	outstruct.IsPartner = *abi.ConvertType(out[2], new(bool)).(*bool)
	outstruct.IsClaimed = *abi.ConvertType(out[3], new(bool)).(*bool)
	outstruct.TotalAmount = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.TotalCount = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.TotalCount2 = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.StakeCount = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.PartnerCount = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// UserInfo is a free data retrieval call binding the contract method 0x1959a002.
//
// Solidity: function userInfo(address account) view returns(uint256 invest, address[] invitees, bool isPartner, bool isClaimed, uint256 totalAmount, uint256 totalCount, uint256 totalCount2, uint256 stakeCount, uint256 partnerCount)
func (_IDO *IDOSession) UserInfo(account common.Address) (struct {
	Invest       *big.Int
	Invitees     []common.Address
	IsPartner    bool
	IsClaimed    bool
	TotalAmount  *big.Int
	TotalCount   *big.Int
	TotalCount2  *big.Int
	StakeCount   *big.Int
	PartnerCount *big.Int
}, error) {
	return _IDO.Contract.UserInfo(&_IDO.CallOpts, account)
}

// UserInfo is a free data retrieval call binding the contract method 0x1959a002.
//
// Solidity: function userInfo(address account) view returns(uint256 invest, address[] invitees, bool isPartner, bool isClaimed, uint256 totalAmount, uint256 totalCount, uint256 totalCount2, uint256 stakeCount, uint256 partnerCount)
func (_IDO *IDOCallerSession) UserInfo(account common.Address) (struct {
	Invest       *big.Int
	Invitees     []common.Address
	IsPartner    bool
	IsClaimed    bool
	TotalAmount  *big.Int
	TotalCount   *big.Int
	TotalCount2  *big.Int
	StakeCount   *big.Int
	PartnerCount *big.Int
}, error) {
	return _IDO.Contract.UserInfo(&_IDO.CallOpts, account)
}

// Claim is a paid mutator transaction binding the contract method 0x1e83409a.
//
// Solidity: function claim(address recipient) returns()
func (_IDO *IDOTransactor) Claim(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _IDO.contract.Transact(opts, "claim", recipient)
}

// Claim is a paid mutator transaction binding the contract method 0x1e83409a.
//
// Solidity: function claim(address recipient) returns()
func (_IDO *IDOSession) Claim(recipient common.Address) (*types.Transaction, error) {
	return _IDO.Contract.Claim(&_IDO.TransactOpts, recipient)
}

// Claim is a paid mutator transaction binding the contract method 0x1e83409a.
//
// Solidity: function claim(address recipient) returns()
func (_IDO *IDOTransactorSession) Claim(recipient common.Address) (*types.Transaction, error) {
	return _IDO.Contract.Claim(&_IDO.TransactOpts, recipient)
}

// Invest1 is a paid mutator transaction binding the contract method 0x9641119e.
//
// Solidity: function invest1() returns()
func (_IDO *IDOTransactor) Invest1(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IDO.contract.Transact(opts, "invest1")
}

// Invest1 is a paid mutator transaction binding the contract method 0x9641119e.
//
// Solidity: function invest1() returns()
func (_IDO *IDOSession) Invest1() (*types.Transaction, error) {
	return _IDO.Contract.Invest1(&_IDO.TransactOpts)
}

// Invest1 is a paid mutator transaction binding the contract method 0x9641119e.
//
// Solidity: function invest1() returns()
func (_IDO *IDOTransactorSession) Invest1() (*types.Transaction, error) {
	return _IDO.Contract.Invest1(&_IDO.TransactOpts)
}

// Invest2 is a paid mutator transaction binding the contract method 0x6ad72abe.
//
// Solidity: function invest2() returns()
func (_IDO *IDOTransactor) Invest2(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IDO.contract.Transact(opts, "invest2")
}

// Invest2 is a paid mutator transaction binding the contract method 0x6ad72abe.
//
// Solidity: function invest2() returns()
func (_IDO *IDOSession) Invest2() (*types.Transaction, error) {
	return _IDO.Contract.Invest2(&_IDO.TransactOpts)
}

// Invest2 is a paid mutator transaction binding the contract method 0x6ad72abe.
//
// Solidity: function invest2() returns()
func (_IDO *IDOTransactorSession) Invest2() (*types.Transaction, error) {
	return _IDO.Contract.Invest2(&_IDO.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IDO *IDOTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IDO.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IDO *IDOSession) RenounceOwnership() (*types.Transaction, error) {
	return _IDO.Contract.RenounceOwnership(&_IDO.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_IDO *IDOTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _IDO.Contract.RenounceOwnership(&_IDO.TransactOpts)
}

// SetSatAddress is a paid mutator transaction binding the contract method 0x8aaf4bfc.
//
// Solidity: function setSatAddress(address sat) returns()
func (_IDO *IDOTransactor) SetSatAddress(opts *bind.TransactOpts, sat common.Address) (*types.Transaction, error) {
	return _IDO.contract.Transact(opts, "setSatAddress", sat)
}

// SetSatAddress is a paid mutator transaction binding the contract method 0x8aaf4bfc.
//
// Solidity: function setSatAddress(address sat) returns()
func (_IDO *IDOSession) SetSatAddress(sat common.Address) (*types.Transaction, error) {
	return _IDO.Contract.SetSatAddress(&_IDO.TransactOpts, sat)
}

// SetSatAddress is a paid mutator transaction binding the contract method 0x8aaf4bfc.
//
// Solidity: function setSatAddress(address sat) returns()
func (_IDO *IDOTransactorSession) SetSatAddress(sat common.Address) (*types.Transaction, error) {
	return _IDO.Contract.SetSatAddress(&_IDO.TransactOpts, sat)
}

// SetStop is a paid mutator transaction binding the contract method 0x641657cb.
//
// Solidity: function setStop(bool s) returns()
func (_IDO *IDOTransactor) SetStop(opts *bind.TransactOpts, s bool) (*types.Transaction, error) {
	return _IDO.contract.Transact(opts, "setStop", s)
}

// SetStop is a paid mutator transaction binding the contract method 0x641657cb.
//
// Solidity: function setStop(bool s) returns()
func (_IDO *IDOSession) SetStop(s bool) (*types.Transaction, error) {
	return _IDO.Contract.SetStop(&_IDO.TransactOpts, s)
}

// SetStop is a paid mutator transaction binding the contract method 0x641657cb.
//
// Solidity: function setStop(bool s) returns()
func (_IDO *IDOTransactorSession) SetStop(s bool) (*types.Transaction, error) {
	return _IDO.Contract.SetStop(&_IDO.TransactOpts, s)
}

// StakeUsdt is a paid mutator transaction binding the contract method 0x85b9a448.
//
// Solidity: function stakeUsdt() returns()
func (_IDO *IDOTransactor) StakeUsdt(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IDO.contract.Transact(opts, "stakeUsdt")
}

// StakeUsdt is a paid mutator transaction binding the contract method 0x85b9a448.
//
// Solidity: function stakeUsdt() returns()
func (_IDO *IDOSession) StakeUsdt() (*types.Transaction, error) {
	return _IDO.Contract.StakeUsdt(&_IDO.TransactOpts)
}

// StakeUsdt is a paid mutator transaction binding the contract method 0x85b9a448.
//
// Solidity: function stakeUsdt() returns()
func (_IDO *IDOTransactorSession) StakeUsdt() (*types.Transaction, error) {
	return _IDO.Contract.StakeUsdt(&_IDO.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IDO *IDOTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _IDO.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IDO *IDOSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IDO.Contract.TransferOwnership(&_IDO.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_IDO *IDOTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _IDO.Contract.TransferOwnership(&_IDO.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address token, address recipient, uint256 amount) returns()
func (_IDO *IDOTransactor) Withdraw(opts *bind.TransactOpts, token common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IDO.contract.Transact(opts, "withdraw", token, recipient, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address token, address recipient, uint256 amount) returns()
func (_IDO *IDOSession) Withdraw(token common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IDO.Contract.Withdraw(&_IDO.TransactOpts, token, recipient, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address token, address recipient, uint256 amount) returns()
func (_IDO *IDOTransactorSession) Withdraw(token common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IDO.Contract.Withdraw(&_IDO.TransactOpts, token, recipient, amount)
}

// IDOClaimIterator is returned from FilterClaim and is used to iterate over the raw logs and unpacked data for Claim events raised by the IDO contract.
type IDOClaimIterator struct {
	Event *IDOClaim // Event containing the contract specifics and raw log

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
func (it *IDOClaimIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IDOClaim)
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
		it.Event = new(IDOClaim)
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
func (it *IDOClaimIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IDOClaimIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IDOClaim represents a Claim event raised by the IDO contract.
type IDOClaim struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterClaim is a free log retrieval operation binding the contract event 0x70eb43c4a8ae8c40502dcf22436c509c28d6ff421cf07c491be56984bd987068.
//
// Solidity: event Claim(address indexed sender, address indexed recipient, uint256 indexed amount)
func (_IDO *IDOFilterer) FilterClaim(opts *bind.FilterOpts, sender []common.Address, recipient []common.Address, amount []*big.Int) (*IDOClaimIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IDO.contract.FilterLogs(opts, "Claim", senderRule, recipientRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &IDOClaimIterator{contract: _IDO.contract, event: "Claim", logs: logs, sub: sub}, nil
}

// WatchClaim is a free log subscription operation binding the contract event 0x70eb43c4a8ae8c40502dcf22436c509c28d6ff421cf07c491be56984bd987068.
//
// Solidity: event Claim(address indexed sender, address indexed recipient, uint256 indexed amount)
func (_IDO *IDOFilterer) WatchClaim(opts *bind.WatchOpts, sink chan<- *IDOClaim, sender []common.Address, recipient []common.Address, amount []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IDO.contract.WatchLogs(opts, "Claim", senderRule, recipientRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IDOClaim)
				if err := _IDO.contract.UnpackLog(event, "Claim", log); err != nil {
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

// ParseClaim is a log parse operation binding the contract event 0x70eb43c4a8ae8c40502dcf22436c509c28d6ff421cf07c491be56984bd987068.
//
// Solidity: event Claim(address indexed sender, address indexed recipient, uint256 indexed amount)
func (_IDO *IDOFilterer) ParseClaim(log types.Log) (*IDOClaim, error) {
	event := new(IDOClaim)
	if err := _IDO.contract.UnpackLog(event, "Claim", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IDOInvestIterator is returned from FilterInvest and is used to iterate over the raw logs and unpacked data for Invest events raised by the IDO contract.
type IDOInvestIterator struct {
	Event *IDOInvest // Event containing the contract specifics and raw log

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
func (it *IDOInvestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IDOInvest)
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
		it.Event = new(IDOInvest)
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
func (it *IDOInvestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IDOInvestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IDOInvest represents a Invest event raised by the IDO contract.
type IDOInvest struct {
	Sender common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterInvest is a free log retrieval operation binding the contract event 0xd90d253a9de34d2fdd5a75ae49ea17fcb43af32fc8ea08cc6d2341991dd3872e.
//
// Solidity: event Invest(address indexed sender, uint256 indexed amount)
func (_IDO *IDOFilterer) FilterInvest(opts *bind.FilterOpts, sender []common.Address, amount []*big.Int) (*IDOInvestIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IDO.contract.FilterLogs(opts, "Invest", senderRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &IDOInvestIterator{contract: _IDO.contract, event: "Invest", logs: logs, sub: sub}, nil
}

// WatchInvest is a free log subscription operation binding the contract event 0xd90d253a9de34d2fdd5a75ae49ea17fcb43af32fc8ea08cc6d2341991dd3872e.
//
// Solidity: event Invest(address indexed sender, uint256 indexed amount)
func (_IDO *IDOFilterer) WatchInvest(opts *bind.WatchOpts, sink chan<- *IDOInvest, sender []common.Address, amount []*big.Int) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _IDO.contract.WatchLogs(opts, "Invest", senderRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IDOInvest)
				if err := _IDO.contract.UnpackLog(event, "Invest", log); err != nil {
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

// ParseInvest is a log parse operation binding the contract event 0xd90d253a9de34d2fdd5a75ae49ea17fcb43af32fc8ea08cc6d2341991dd3872e.
//
// Solidity: event Invest(address indexed sender, uint256 indexed amount)
func (_IDO *IDOFilterer) ParseInvest(log types.Log) (*IDOInvest, error) {
	event := new(IDOInvest)
	if err := _IDO.contract.UnpackLog(event, "Invest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IDOOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the IDO contract.
type IDOOwnershipTransferredIterator struct {
	Event *IDOOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *IDOOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IDOOwnershipTransferred)
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
		it.Event = new(IDOOwnershipTransferred)
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
func (it *IDOOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IDOOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IDOOwnershipTransferred represents a OwnershipTransferred event raised by the IDO contract.
type IDOOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IDO *IDOFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*IDOOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IDO.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &IDOOwnershipTransferredIterator{contract: _IDO.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IDO *IDOFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *IDOOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _IDO.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IDOOwnershipTransferred)
				if err := _IDO.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_IDO *IDOFilterer) ParseOwnershipTransferred(log types.Log) (*IDOOwnershipTransferred, error) {
	event := new(IDOOwnershipTransferred)
	if err := _IDO.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
