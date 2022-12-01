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

// InvitationABI is the input ABI used to generate the binding from.
const InvitationABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"inviter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"invitee\",\"type\":\"address\"}],\"name\":\"Bind\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"inviter\",\"type\":\"address\"}],\"name\":\"bind\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getInvitation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"inviter\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"invitees\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Invitation is an auto generated Go binding around an Ethereum contract.
type Invitation struct {
	InvitationCaller     // Read-only binding to the contract
	InvitationTransactor // Write-only binding to the contract
	InvitationFilterer   // Log filterer for contract events
}

// InvitationCaller is an auto generated read-only Go binding around an Ethereum contract.
type InvitationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InvitationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type InvitationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InvitationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type InvitationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InvitationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type InvitationSession struct {
	Contract     *Invitation       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// InvitationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type InvitationCallerSession struct {
	Contract *InvitationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// InvitationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type InvitationTransactorSession struct {
	Contract     *InvitationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// InvitationRaw is an auto generated low-level Go binding around an Ethereum contract.
type InvitationRaw struct {
	Contract *Invitation // Generic contract binding to access the raw methods on
}

// InvitationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type InvitationCallerRaw struct {
	Contract *InvitationCaller // Generic read-only contract binding to access the raw methods on
}

// InvitationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type InvitationTransactorRaw struct {
	Contract *InvitationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewInvitation creates a new instance of Invitation, bound to a specific deployed contract.
func NewInvitation(address common.Address, backend bind.ContractBackend) (*Invitation, error) {
	contract, err := bindInvitation(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Invitation{InvitationCaller: InvitationCaller{contract: contract}, InvitationTransactor: InvitationTransactor{contract: contract}, InvitationFilterer: InvitationFilterer{contract: contract}}, nil
}

// NewInvitationCaller creates a new read-only instance of Invitation, bound to a specific deployed contract.
func NewInvitationCaller(address common.Address, caller bind.ContractCaller) (*InvitationCaller, error) {
	contract, err := bindInvitation(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &InvitationCaller{contract: contract}, nil
}

// NewInvitationTransactor creates a new write-only instance of Invitation, bound to a specific deployed contract.
func NewInvitationTransactor(address common.Address, transactor bind.ContractTransactor) (*InvitationTransactor, error) {
	contract, err := bindInvitation(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &InvitationTransactor{contract: contract}, nil
}

// NewInvitationFilterer creates a new log filterer instance of Invitation, bound to a specific deployed contract.
func NewInvitationFilterer(address common.Address, filterer bind.ContractFilterer) (*InvitationFilterer, error) {
	contract, err := bindInvitation(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &InvitationFilterer{contract: contract}, nil
}

// bindInvitation binds a generic wrapper to an already deployed contract.
func bindInvitation(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(InvitationABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Invitation *InvitationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Invitation.Contract.InvitationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Invitation *InvitationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Invitation.Contract.InvitationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Invitation *InvitationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Invitation.Contract.InvitationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Invitation *InvitationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Invitation.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Invitation *InvitationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Invitation.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Invitation *InvitationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Invitation.Contract.contract.Transact(opts, method, params...)
}

// GetInvitation is a free data retrieval call binding the contract method 0x2b17b8be.
//
// Solidity: function getInvitation(address user) view returns(address inviter, address[] invitees)
func (_Invitation *InvitationCaller) GetInvitation(opts *bind.CallOpts, user common.Address) (struct {
	Inviter  common.Address
	Invitees []common.Address
}, error) {
	var out []interface{}
	err := _Invitation.contract.Call(opts, &out, "getInvitation", user)

	outstruct := new(struct {
		Inviter  common.Address
		Invitees []common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Inviter = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Invitees = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

// GetInvitation is a free data retrieval call binding the contract method 0x2b17b8be.
//
// Solidity: function getInvitation(address user) view returns(address inviter, address[] invitees)
func (_Invitation *InvitationSession) GetInvitation(user common.Address) (struct {
	Inviter  common.Address
	Invitees []common.Address
}, error) {
	return _Invitation.Contract.GetInvitation(&_Invitation.CallOpts, user)
}

// GetInvitation is a free data retrieval call binding the contract method 0x2b17b8be.
//
// Solidity: function getInvitation(address user) view returns(address inviter, address[] invitees)
func (_Invitation *InvitationCallerSession) GetInvitation(user common.Address) (struct {
	Inviter  common.Address
	Invitees []common.Address
}, error) {
	return _Invitation.Contract.GetInvitation(&_Invitation.CallOpts, user)
}

// Bind is a paid mutator transaction binding the contract method 0x81bac14f.
//
// Solidity: function bind(address inviter) returns()
func (_Invitation *InvitationTransactor) Bind(opts *bind.TransactOpts, inviter common.Address) (*types.Transaction, error) {
	return _Invitation.contract.Transact(opts, "bind", inviter)
}

// Bind is a paid mutator transaction binding the contract method 0x81bac14f.
//
// Solidity: function bind(address inviter) returns()
func (_Invitation *InvitationSession) Bind(inviter common.Address) (*types.Transaction, error) {
	return _Invitation.Contract.Bind(&_Invitation.TransactOpts, inviter)
}

// Bind is a paid mutator transaction binding the contract method 0x81bac14f.
//
// Solidity: function bind(address inviter) returns()
func (_Invitation *InvitationTransactorSession) Bind(inviter common.Address) (*types.Transaction, error) {
	return _Invitation.Contract.Bind(&_Invitation.TransactOpts, inviter)
}

// InvitationBindIterator is returned from FilterBind and is used to iterate over the raw logs and unpacked data for Bind events raised by the Invitation contract.
type InvitationBindIterator struct {
	Event *InvitationBind // Event containing the contract specifics and raw log

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
func (it *InvitationBindIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InvitationBind)
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
		it.Event = new(InvitationBind)
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
func (it *InvitationBindIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InvitationBindIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InvitationBind represents a Bind event raised by the Invitation contract.
type InvitationBind struct {
	Inviter common.Address
	Invitee common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBind is a free log retrieval operation binding the contract event 0x3693edbeef6f168d1e6eb95f7f23f3184ba7a1e173a83a11ead1fecf6fcbe034.
//
// Solidity: event Bind(address indexed inviter, address indexed invitee)
func (_Invitation *InvitationFilterer) FilterBind(opts *bind.FilterOpts, inviter []common.Address, invitee []common.Address) (*InvitationBindIterator, error) {

	var inviterRule []interface{}
	for _, inviterItem := range inviter {
		inviterRule = append(inviterRule, inviterItem)
	}
	var inviteeRule []interface{}
	for _, inviteeItem := range invitee {
		inviteeRule = append(inviteeRule, inviteeItem)
	}

	logs, sub, err := _Invitation.contract.FilterLogs(opts, "Bind", inviterRule, inviteeRule)
	if err != nil {
		return nil, err
	}
	return &InvitationBindIterator{contract: _Invitation.contract, event: "Bind", logs: logs, sub: sub}, nil
}

// WatchBind is a free log subscription operation binding the contract event 0x3693edbeef6f168d1e6eb95f7f23f3184ba7a1e173a83a11ead1fecf6fcbe034.
//
// Solidity: event Bind(address indexed inviter, address indexed invitee)
func (_Invitation *InvitationFilterer) WatchBind(opts *bind.WatchOpts, sink chan<- *InvitationBind, inviter []common.Address, invitee []common.Address) (event.Subscription, error) {

	var inviterRule []interface{}
	for _, inviterItem := range inviter {
		inviterRule = append(inviterRule, inviterItem)
	}
	var inviteeRule []interface{}
	for _, inviteeItem := range invitee {
		inviteeRule = append(inviteeRule, inviteeItem)
	}

	logs, sub, err := _Invitation.contract.WatchLogs(opts, "Bind", inviterRule, inviteeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InvitationBind)
				if err := _Invitation.contract.UnpackLog(event, "Bind", log); err != nil {
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

// ParseBind is a log parse operation binding the contract event 0x3693edbeef6f168d1e6eb95f7f23f3184ba7a1e173a83a11ead1fecf6fcbe034.
//
// Solidity: event Bind(address indexed inviter, address indexed invitee)
func (_Invitation *InvitationFilterer) ParseBind(log types.Log) (*InvitationBind, error) {
	event := new(InvitationBind)
	if err := _Invitation.contract.UnpackLog(event, "Bind", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
