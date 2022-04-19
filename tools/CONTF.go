// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tools

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
)

// CONTFMetaData contains all meta data concerning the CONTF contract.
var CONTFMetaData = &bind.MetaData{
	ABI: "[{\"constant\":false,\"inputs\":[{\"name\":\"s\",\"type\":\"string\"},{\"name\":\"t\",\"type\":\"string\"}],\"name\":\"setWord\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"random\",\"type\":\"uint256\"}],\"name\":\"getRandomWord\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"string\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// CONTFABI is the input ABI used to generate the binding from.
// Deprecated: Use CONTFMetaData.ABI instead.
var CONTFABI = CONTFMetaData.ABI

// CONTF is an auto generated Go binding around an Ethereum contract.
type CONTF struct {
	CONTFCaller     // Read-only binding to the contract
	CONTFTransactor // Write-only binding to the contract
	CONTFFilterer   // Log filterer for contract events
}

// CONTFCaller is an auto generated read-only Go binding around an Ethereum contract.
type CONTFCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CONTFTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CONTFTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CONTFFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CONTFFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CONTFSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CONTFSession struct {
	Contract     *CONTF            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CONTFCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CONTFCallerSession struct {
	Contract *CONTFCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// CONTFTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CONTFTransactorSession struct {
	Contract     *CONTFTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CONTFRaw is an auto generated low-level Go binding around an Ethereum contract.
type CONTFRaw struct {
	Contract *CONTF // Generic contract binding to access the raw methods on
}

// CONTFCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CONTFCallerRaw struct {
	Contract *CONTFCaller // Generic read-only contract binding to access the raw methods on
}

// CONTFTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CONTFTransactorRaw struct {
	Contract *CONTFTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCONTF creates a new instance of CONTF, bound to a specific deployed contract.
func NewCONTF(address common.Address, backend bind.ContractBackend) (*CONTF, error) {
	contract, err := bindCONTF(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CONTF{CONTFCaller: CONTFCaller{contract: contract}, CONTFTransactor: CONTFTransactor{contract: contract}, CONTFFilterer: CONTFFilterer{contract: contract}}, nil
}

// NewCONTFCaller creates a new read-only instance of CONTF, bound to a specific deployed contract.
func NewCONTFCaller(address common.Address, caller bind.ContractCaller) (*CONTFCaller, error) {
	contract, err := bindCONTF(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CONTFCaller{contract: contract}, nil
}

// NewCONTFTransactor creates a new write-only instance of CONTF, bound to a specific deployed contract.
func NewCONTFTransactor(address common.Address, transactor bind.ContractTransactor) (*CONTFTransactor, error) {
	contract, err := bindCONTF(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CONTFTransactor{contract: contract}, nil
}

// NewCONTFFilterer creates a new log filterer instance of CONTF, bound to a specific deployed contract.
func NewCONTFFilterer(address common.Address, filterer bind.ContractFilterer) (*CONTFFilterer, error) {
	contract, err := bindCONTF(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CONTFFilterer{contract: contract}, nil
}

// bindCONTF binds a generic wrapper to an already deployed contract.
func bindCONTF(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CONTFABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CONTF *CONTFRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CONTF.Contract.CONTFCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CONTF *CONTFRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CONTF.Contract.CONTFTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CONTF *CONTFRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CONTF.Contract.CONTFTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CONTF *CONTFCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CONTF.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CONTF *CONTFTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CONTF.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CONTF *CONTFTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CONTF.Contract.contract.Transact(opts, method, params...)
}

// GetRandomWord is a free data retrieval call binding the contract method 0xcc76b911.
//
// Solidity: function getRandomWord(uint256 random) view returns(uint256, string, address, string)
func (_CONTF *CONTFCaller) GetRandomWord(opts *bind.CallOpts, random *big.Int) (*big.Int, string, common.Address, string, error) {
	var out []interface{}
	err := _CONTF.contract.Call(opts, &out, "getRandomWord", random)

	if err != nil {
		return *new(*big.Int), *new(string), *new(common.Address), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)
	out2 := *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	out3 := *abi.ConvertType(out[3], new(string)).(*string)

	return out0, out1, out2, out3, err

}

// GetRandomWord is a free data retrieval call binding the contract method 0xcc76b911.
//
// Solidity: function getRandomWord(uint256 random) view returns(uint256, string, address, string)
func (_CONTF *CONTFSession) GetRandomWord(random *big.Int) (*big.Int, string, common.Address, string, error) {
	return _CONTF.Contract.GetRandomWord(&_CONTF.CallOpts, random)
}

// GetRandomWord is a free data retrieval call binding the contract method 0xcc76b911.
//
// Solidity: function getRandomWord(uint256 random) view returns(uint256, string, address, string)
func (_CONTF *CONTFCallerSession) GetRandomWord(random *big.Int) (*big.Int, string, common.Address, string, error) {
	return _CONTF.Contract.GetRandomWord(&_CONTF.CallOpts, random)
}

// SetWord is a paid mutator transaction binding the contract method 0x8c595a95.
//
// Solidity: function setWord(string s, string t) returns()
func (_CONTF *CONTFTransactor) SetWord(opts *bind.TransactOpts, s string, t string) (*types.Transaction, error) {
	return _CONTF.contract.Transact(opts, "setWord", s, t)
}

// SetWord is a paid mutator transaction binding the contract method 0x8c595a95.
//
// Solidity: function setWord(string s, string t) returns()
func (_CONTF *CONTFSession) SetWord(s string, t string) (*types.Transaction, error) {
	return _CONTF.Contract.SetWord(&_CONTF.TransactOpts, s, t)
}

// SetWord is a paid mutator transaction binding the contract method 0x8c595a95.
//
// Solidity: function setWord(string s, string t) returns()
func (_CONTF *CONTFTransactorSession) SetWord(s string, t string) (*types.Transaction, error) {
	return _CONTF.Contract.SetWord(&_CONTF.TransactOpts, s, t)
}
