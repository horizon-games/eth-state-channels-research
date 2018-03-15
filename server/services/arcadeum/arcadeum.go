// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package arcadeum

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

// ArcadeumABI is the input ABI used to generate the binding from.
const ArcadeumABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"game\",\"type\":\"address\"},{\"name\":\"matchID\",\"type\":\"uint32\"},{\"name\":\"timestamp\",\"type\":\"uint256\"},{\"name\":\"timestampV\",\"type\":\"uint8\"},{\"name\":\"timestampR\",\"type\":\"bytes32\"},{\"name\":\"timestampS\",\"type\":\"bytes32\"},{\"name\":\"subkeyV\",\"type\":\"uint8\"},{\"name\":\"subkeyR\",\"type\":\"bytes32\"},{\"name\":\"subkeyS\",\"type\":\"bytes32\"}],\"name\":\"canStopWithdrawalXXX\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"subkey\",\"type\":\"address\"}],\"name\":\"subkeyMessage\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"game\",\"type\":\"address\"},{\"name\":\"matchID\",\"type\":\"uint32\"},{\"name\":\"timestamp\",\"type\":\"uint256\"},{\"name\":\"timestampV\",\"type\":\"uint8\"},{\"name\":\"timestampR\",\"type\":\"bytes32\"},{\"name\":\"timestampS\",\"type\":\"bytes32\"},{\"name\":\"subkeyV\",\"type\":\"uint8\"},{\"name\":\"subkeyR\",\"type\":\"bytes32\"},{\"name\":\"subkeyS\",\"type\":\"bytes32\"}],\"name\":\"couldStopWithdrawalXXX\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\"},{\"name\":\"timestampV\",\"type\":\"uint8\"},{\"name\":\"timestampR\",\"type\":\"bytes32\"},{\"name\":\"timestampS\",\"type\":\"bytes32\"}],\"name\":\"timestampSubkeyXXX\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"subkey\",\"type\":\"address\"},{\"name\":\"subkeyV\",\"type\":\"uint8\"},{\"name\":\"subkeyR\",\"type\":\"bytes32\"},{\"name\":\"subkeyS\",\"type\":\"bytes32\"}],\"name\":\"subkeyParentXXX\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"game\",\"type\":\"address\"},{\"name\":\"matchID\",\"type\":\"uint32\"},{\"name\":\"timestamp\",\"type\":\"uint256\"},{\"name\":\"timestampV\",\"type\":\"uint8\"},{\"name\":\"timestampR\",\"type\":\"bytes32\"},{\"name\":\"timestampS\",\"type\":\"bytes32\"},{\"name\":\"subkeyV\",\"type\":\"uint8\"},{\"name\":\"subkeyR\",\"type\":\"bytes32\"},{\"name\":\"subkeyS\",\"type\":\"bytes32\"}],\"name\":\"stopWithdrawalXXX\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"canFinishWithdrawal\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"game\",\"type\":\"address\"},{\"name\":\"matchID\",\"type\":\"uint32\"},{\"name\":\"timestamp\",\"type\":\"uint256\"},{\"name\":\"timestampV\",\"type\":\"uint8\"},{\"name\":\"timestampR\",\"type\":\"bytes32\"},{\"name\":\"timestampS\",\"type\":\"bytes32\"},{\"name\":\"subkeyV\",\"type\":\"uint8\"},{\"name\":\"subkeyR\",\"type\":\"bytes32\"},{\"name\":\"subkeyS\",\"type\":\"bytes32\"}],\"name\":\"playerAccountXXX\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"startWithdrawal\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finishWithdrawal\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"game\",\"type\":\"address\"},{\"name\":\"matchID\",\"type\":\"uint32\"},{\"name\":\"timestamp\",\"type\":\"uint256\"},{\"name\":\"accounts\",\"type\":\"address[2]\"},{\"name\":\"seedRatings\",\"type\":\"uint32[2]\"},{\"name\":\"publicSeeds\",\"type\":\"bytes32[1][2]\"}],\"name\":\"matchHash\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isWithdrawing\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"withdrawalStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"withdrawalStopped\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"game\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"matchID\",\"type\":\"uint32\"},{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"rewardClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"game\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"matchID\",\"type\":\"uint32\"},{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"cheaterReported\",\"type\":\"event\"}]"

// ArcadeumBin is the compiled bytecode used for deploying new contracts.
const ArcadeumBin = `0x`

// DeployArcadeum deploys a new Ethereum contract, binding an instance of Arcadeum to it.
func DeployArcadeum(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Arcadeum, error) {
	parsed, err := abi.JSON(strings.NewReader(ArcadeumABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ArcadeumBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Arcadeum{ArcadeumCaller: ArcadeumCaller{contract: contract}, ArcadeumTransactor: ArcadeumTransactor{contract: contract}, ArcadeumFilterer: ArcadeumFilterer{contract: contract}}, nil
}

// Arcadeum is an auto generated Go binding around an Ethereum contract.
type Arcadeum struct {
	ArcadeumCaller     // Read-only binding to the contract
	ArcadeumTransactor // Write-only binding to the contract
	ArcadeumFilterer   // Log filterer for contract events
}

// ArcadeumCaller is an auto generated read-only Go binding around an Ethereum contract.
type ArcadeumCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArcadeumTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ArcadeumTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArcadeumFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ArcadeumFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ArcadeumSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ArcadeumSession struct {
	Contract     *Arcadeum         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ArcadeumCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ArcadeumCallerSession struct {
	Contract *ArcadeumCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ArcadeumTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ArcadeumTransactorSession struct {
	Contract     *ArcadeumTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ArcadeumRaw is an auto generated low-level Go binding around an Ethereum contract.
type ArcadeumRaw struct {
	Contract *Arcadeum // Generic contract binding to access the raw methods on
}

// ArcadeumCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ArcadeumCallerRaw struct {
	Contract *ArcadeumCaller // Generic read-only contract binding to access the raw methods on
}

// ArcadeumTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ArcadeumTransactorRaw struct {
	Contract *ArcadeumTransactor // Generic write-only contract binding to access the raw methods on
}

// NewArcadeum creates a new instance of Arcadeum, bound to a specific deployed contract.
func NewArcadeum(address common.Address, backend bind.ContractBackend) (*Arcadeum, error) {
	contract, err := bindArcadeum(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Arcadeum{ArcadeumCaller: ArcadeumCaller{contract: contract}, ArcadeumTransactor: ArcadeumTransactor{contract: contract}, ArcadeumFilterer: ArcadeumFilterer{contract: contract}}, nil
}

// NewArcadeumCaller creates a new read-only instance of Arcadeum, bound to a specific deployed contract.
func NewArcadeumCaller(address common.Address, caller bind.ContractCaller) (*ArcadeumCaller, error) {
	contract, err := bindArcadeum(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ArcadeumCaller{contract: contract}, nil
}

// NewArcadeumTransactor creates a new write-only instance of Arcadeum, bound to a specific deployed contract.
func NewArcadeumTransactor(address common.Address, transactor bind.ContractTransactor) (*ArcadeumTransactor, error) {
	contract, err := bindArcadeum(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ArcadeumTransactor{contract: contract}, nil
}

// NewArcadeumFilterer creates a new log filterer instance of Arcadeum, bound to a specific deployed contract.
func NewArcadeumFilterer(address common.Address, filterer bind.ContractFilterer) (*ArcadeumFilterer, error) {
	contract, err := bindArcadeum(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ArcadeumFilterer{contract: contract}, nil
}

// bindArcadeum binds a generic wrapper to an already deployed contract.
func bindArcadeum(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ArcadeumABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Arcadeum *ArcadeumRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Arcadeum.Contract.ArcadeumCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Arcadeum *ArcadeumRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Arcadeum.Contract.ArcadeumTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Arcadeum *ArcadeumRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Arcadeum.Contract.ArcadeumTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Arcadeum *ArcadeumCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Arcadeum.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Arcadeum *ArcadeumTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Arcadeum.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Arcadeum *ArcadeumTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Arcadeum.Contract.contract.Transact(opts, method, params...)
}

// Balance is a free data retrieval call binding the contract method 0xe3d670d7.
//
// Solidity: function balance( address) constant returns(uint256)
func (_Arcadeum *ArcadeumCaller) Balance(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Arcadeum.contract.Call(opts, out, "balance", arg0)
	return *ret0, err
}

// Balance is a free data retrieval call binding the contract method 0xe3d670d7.
//
// Solidity: function balance( address) constant returns(uint256)
func (_Arcadeum *ArcadeumSession) Balance(arg0 common.Address) (*big.Int, error) {
	return _Arcadeum.Contract.Balance(&_Arcadeum.CallOpts, arg0)
}

// Balance is a free data retrieval call binding the contract method 0xe3d670d7.
//
// Solidity: function balance( address) constant returns(uint256)
func (_Arcadeum *ArcadeumCallerSession) Balance(arg0 common.Address) (*big.Int, error) {
	return _Arcadeum.Contract.Balance(&_Arcadeum.CallOpts, arg0)
}

// CanFinishWithdrawal is a free data retrieval call binding the contract method 0x98b6a663.
//
// Solidity: function canFinishWithdrawal(account address) constant returns(bool)
func (_Arcadeum *ArcadeumCaller) CanFinishWithdrawal(opts *bind.CallOpts, account common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Arcadeum.contract.Call(opts, out, "canFinishWithdrawal", account)
	return *ret0, err
}

// CanFinishWithdrawal is a free data retrieval call binding the contract method 0x98b6a663.
//
// Solidity: function canFinishWithdrawal(account address) constant returns(bool)
func (_Arcadeum *ArcadeumSession) CanFinishWithdrawal(account common.Address) (bool, error) {
	return _Arcadeum.Contract.CanFinishWithdrawal(&_Arcadeum.CallOpts, account)
}

// CanFinishWithdrawal is a free data retrieval call binding the contract method 0x98b6a663.
//
// Solidity: function canFinishWithdrawal(account address) constant returns(bool)
func (_Arcadeum *ArcadeumCallerSession) CanFinishWithdrawal(account common.Address) (bool, error) {
	return _Arcadeum.Contract.CanFinishWithdrawal(&_Arcadeum.CallOpts, account)
}

// CanStopWithdrawalXXX is a free data retrieval call binding the contract method 0x0fad2462.
//
// Solidity: function canStopWithdrawalXXX(game address, matchID uint32, timestamp uint256, timestampV uint8, timestampR bytes32, timestampS bytes32, subkeyV uint8, subkeyR bytes32, subkeyS bytes32) constant returns(bool)
func (_Arcadeum *ArcadeumCaller) CanStopWithdrawalXXX(opts *bind.CallOpts, game common.Address, matchID uint32, timestamp *big.Int, timestampV uint8, timestampR [32]byte, timestampS [32]byte, subkeyV uint8, subkeyR [32]byte, subkeyS [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Arcadeum.contract.Call(opts, out, "canStopWithdrawalXXX", game, matchID, timestamp, timestampV, timestampR, timestampS, subkeyV, subkeyR, subkeyS)
	return *ret0, err
}

// CanStopWithdrawalXXX is a free data retrieval call binding the contract method 0x0fad2462.
//
// Solidity: function canStopWithdrawalXXX(game address, matchID uint32, timestamp uint256, timestampV uint8, timestampR bytes32, timestampS bytes32, subkeyV uint8, subkeyR bytes32, subkeyS bytes32) constant returns(bool)
func (_Arcadeum *ArcadeumSession) CanStopWithdrawalXXX(game common.Address, matchID uint32, timestamp *big.Int, timestampV uint8, timestampR [32]byte, timestampS [32]byte, subkeyV uint8, subkeyR [32]byte, subkeyS [32]byte) (bool, error) {
	return _Arcadeum.Contract.CanStopWithdrawalXXX(&_Arcadeum.CallOpts, game, matchID, timestamp, timestampV, timestampR, timestampS, subkeyV, subkeyR, subkeyS)
}

// CanStopWithdrawalXXX is a free data retrieval call binding the contract method 0x0fad2462.
//
// Solidity: function canStopWithdrawalXXX(game address, matchID uint32, timestamp uint256, timestampV uint8, timestampR bytes32, timestampS bytes32, subkeyV uint8, subkeyR bytes32, subkeyS bytes32) constant returns(bool)
func (_Arcadeum *ArcadeumCallerSession) CanStopWithdrawalXXX(game common.Address, matchID uint32, timestamp *big.Int, timestampV uint8, timestampR [32]byte, timestampS [32]byte, subkeyV uint8, subkeyR [32]byte, subkeyS [32]byte) (bool, error) {
	return _Arcadeum.Contract.CanStopWithdrawalXXX(&_Arcadeum.CallOpts, game, matchID, timestamp, timestampV, timestampR, timestampS, subkeyV, subkeyR, subkeyS)
}

// CouldStopWithdrawalXXX is a free data retrieval call binding the contract method 0x4898b506.
//
// Solidity: function couldStopWithdrawalXXX(game address, matchID uint32, timestamp uint256, timestampV uint8, timestampR bytes32, timestampS bytes32, subkeyV uint8, subkeyR bytes32, subkeyS bytes32) constant returns(bool)
func (_Arcadeum *ArcadeumCaller) CouldStopWithdrawalXXX(opts *bind.CallOpts, game common.Address, matchID uint32, timestamp *big.Int, timestampV uint8, timestampR [32]byte, timestampS [32]byte, subkeyV uint8, subkeyR [32]byte, subkeyS [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Arcadeum.contract.Call(opts, out, "couldStopWithdrawalXXX", game, matchID, timestamp, timestampV, timestampR, timestampS, subkeyV, subkeyR, subkeyS)
	return *ret0, err
}

// CouldStopWithdrawalXXX is a free data retrieval call binding the contract method 0x4898b506.
//
// Solidity: function couldStopWithdrawalXXX(game address, matchID uint32, timestamp uint256, timestampV uint8, timestampR bytes32, timestampS bytes32, subkeyV uint8, subkeyR bytes32, subkeyS bytes32) constant returns(bool)
func (_Arcadeum *ArcadeumSession) CouldStopWithdrawalXXX(game common.Address, matchID uint32, timestamp *big.Int, timestampV uint8, timestampR [32]byte, timestampS [32]byte, subkeyV uint8, subkeyR [32]byte, subkeyS [32]byte) (bool, error) {
	return _Arcadeum.Contract.CouldStopWithdrawalXXX(&_Arcadeum.CallOpts, game, matchID, timestamp, timestampV, timestampR, timestampS, subkeyV, subkeyR, subkeyS)
}

// CouldStopWithdrawalXXX is a free data retrieval call binding the contract method 0x4898b506.
//
// Solidity: function couldStopWithdrawalXXX(game address, matchID uint32, timestamp uint256, timestampV uint8, timestampR bytes32, timestampS bytes32, subkeyV uint8, subkeyR bytes32, subkeyS bytes32) constant returns(bool)
func (_Arcadeum *ArcadeumCallerSession) CouldStopWithdrawalXXX(game common.Address, matchID uint32, timestamp *big.Int, timestampV uint8, timestampR [32]byte, timestampS [32]byte, subkeyV uint8, subkeyR [32]byte, subkeyS [32]byte) (bool, error) {
	return _Arcadeum.Contract.CouldStopWithdrawalXXX(&_Arcadeum.CallOpts, game, matchID, timestamp, timestampV, timestampR, timestampS, subkeyV, subkeyR, subkeyS)
}

// IsWithdrawing is a free data retrieval call binding the contract method 0xed095b84.
//
// Solidity: function isWithdrawing(account address) constant returns(bool)
func (_Arcadeum *ArcadeumCaller) IsWithdrawing(opts *bind.CallOpts, account common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Arcadeum.contract.Call(opts, out, "isWithdrawing", account)
	return *ret0, err
}

// IsWithdrawing is a free data retrieval call binding the contract method 0xed095b84.
//
// Solidity: function isWithdrawing(account address) constant returns(bool)
func (_Arcadeum *ArcadeumSession) IsWithdrawing(account common.Address) (bool, error) {
	return _Arcadeum.Contract.IsWithdrawing(&_Arcadeum.CallOpts, account)
}

// IsWithdrawing is a free data retrieval call binding the contract method 0xed095b84.
//
// Solidity: function isWithdrawing(account address) constant returns(bool)
func (_Arcadeum *ArcadeumCallerSession) IsWithdrawing(account common.Address) (bool, error) {
	return _Arcadeum.Contract.IsWithdrawing(&_Arcadeum.CallOpts, account)
}

// MatchHash is a free data retrieval call binding the contract method 0xe13e62d0.
//
// Solidity: function matchHash(game address, matchID uint32, timestamp uint256, accounts address[2], seedRatings uint32[2], publicSeeds bytes32[1][2]) constant returns(bytes32)
func (_Arcadeum *ArcadeumCaller) MatchHash(opts *bind.CallOpts, game common.Address, matchID uint32, timestamp *big.Int, accounts [2]common.Address, seedRatings [2]uint32, publicSeeds [2][1][32]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Arcadeum.contract.Call(opts, out, "matchHash", game, matchID, timestamp, accounts, seedRatings, publicSeeds)
	return *ret0, err
}

// MatchHash is a free data retrieval call binding the contract method 0xe13e62d0.
//
// Solidity: function matchHash(game address, matchID uint32, timestamp uint256, accounts address[2], seedRatings uint32[2], publicSeeds bytes32[1][2]) constant returns(bytes32)
func (_Arcadeum *ArcadeumSession) MatchHash(game common.Address, matchID uint32, timestamp *big.Int, accounts [2]common.Address, seedRatings [2]uint32, publicSeeds [2][1][32]byte) ([32]byte, error) {
	return _Arcadeum.Contract.MatchHash(&_Arcadeum.CallOpts, game, matchID, timestamp, accounts, seedRatings, publicSeeds)
}

// MatchHash is a free data retrieval call binding the contract method 0xe13e62d0.
//
// Solidity: function matchHash(game address, matchID uint32, timestamp uint256, accounts address[2], seedRatings uint32[2], publicSeeds bytes32[1][2]) constant returns(bytes32)
func (_Arcadeum *ArcadeumCallerSession) MatchHash(game common.Address, matchID uint32, timestamp *big.Int, accounts [2]common.Address, seedRatings [2]uint32, publicSeeds [2][1][32]byte) ([32]byte, error) {
	return _Arcadeum.Contract.MatchHash(&_Arcadeum.CallOpts, game, matchID, timestamp, accounts, seedRatings, publicSeeds)
}

// PlayerAccountXXX is a free data retrieval call binding the contract method 0xa8f2de36.
//
// Solidity: function playerAccountXXX(game address, matchID uint32, timestamp uint256, timestampV uint8, timestampR bytes32, timestampS bytes32, subkeyV uint8, subkeyR bytes32, subkeyS bytes32) constant returns(address)
func (_Arcadeum *ArcadeumCaller) PlayerAccountXXX(opts *bind.CallOpts, game common.Address, matchID uint32, timestamp *big.Int, timestampV uint8, timestampR [32]byte, timestampS [32]byte, subkeyV uint8, subkeyR [32]byte, subkeyS [32]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Arcadeum.contract.Call(opts, out, "playerAccountXXX", game, matchID, timestamp, timestampV, timestampR, timestampS, subkeyV, subkeyR, subkeyS)
	return *ret0, err
}

// PlayerAccountXXX is a free data retrieval call binding the contract method 0xa8f2de36.
//
// Solidity: function playerAccountXXX(game address, matchID uint32, timestamp uint256, timestampV uint8, timestampR bytes32, timestampS bytes32, subkeyV uint8, subkeyR bytes32, subkeyS bytes32) constant returns(address)
func (_Arcadeum *ArcadeumSession) PlayerAccountXXX(game common.Address, matchID uint32, timestamp *big.Int, timestampV uint8, timestampR [32]byte, timestampS [32]byte, subkeyV uint8, subkeyR [32]byte, subkeyS [32]byte) (common.Address, error) {
	return _Arcadeum.Contract.PlayerAccountXXX(&_Arcadeum.CallOpts, game, matchID, timestamp, timestampV, timestampR, timestampS, subkeyV, subkeyR, subkeyS)
}

// PlayerAccountXXX is a free data retrieval call binding the contract method 0xa8f2de36.
//
// Solidity: function playerAccountXXX(game address, matchID uint32, timestamp uint256, timestampV uint8, timestampR bytes32, timestampS bytes32, subkeyV uint8, subkeyR bytes32, subkeyS bytes32) constant returns(address)
func (_Arcadeum *ArcadeumCallerSession) PlayerAccountXXX(game common.Address, matchID uint32, timestamp *big.Int, timestampV uint8, timestampR [32]byte, timestampS [32]byte, subkeyV uint8, subkeyR [32]byte, subkeyS [32]byte) (common.Address, error) {
	return _Arcadeum.Contract.PlayerAccountXXX(&_Arcadeum.CallOpts, game, matchID, timestamp, timestampV, timestampR, timestampS, subkeyV, subkeyR, subkeyS)
}

// SubkeyMessage is a free data retrieval call binding the contract method 0x41b677db.
//
// Solidity: function subkeyMessage(subkey address) constant returns(string)
func (_Arcadeum *ArcadeumCaller) SubkeyMessage(opts *bind.CallOpts, subkey common.Address) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Arcadeum.contract.Call(opts, out, "subkeyMessage", subkey)
	return *ret0, err
}

// SubkeyMessage is a free data retrieval call binding the contract method 0x41b677db.
//
// Solidity: function subkeyMessage(subkey address) constant returns(string)
func (_Arcadeum *ArcadeumSession) SubkeyMessage(subkey common.Address) (string, error) {
	return _Arcadeum.Contract.SubkeyMessage(&_Arcadeum.CallOpts, subkey)
}

// SubkeyMessage is a free data retrieval call binding the contract method 0x41b677db.
//
// Solidity: function subkeyMessage(subkey address) constant returns(string)
func (_Arcadeum *ArcadeumCallerSession) SubkeyMessage(subkey common.Address) (string, error) {
	return _Arcadeum.Contract.SubkeyMessage(&_Arcadeum.CallOpts, subkey)
}

// SubkeyParentXXX is a free data retrieval call binding the contract method 0x5e15d8a8.
//
// Solidity: function subkeyParentXXX(subkey address, subkeyV uint8, subkeyR bytes32, subkeyS bytes32) constant returns(address)
func (_Arcadeum *ArcadeumCaller) SubkeyParentXXX(opts *bind.CallOpts, subkey common.Address, subkeyV uint8, subkeyR [32]byte, subkeyS [32]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Arcadeum.contract.Call(opts, out, "subkeyParentXXX", subkey, subkeyV, subkeyR, subkeyS)
	return *ret0, err
}

// SubkeyParentXXX is a free data retrieval call binding the contract method 0x5e15d8a8.
//
// Solidity: function subkeyParentXXX(subkey address, subkeyV uint8, subkeyR bytes32, subkeyS bytes32) constant returns(address)
func (_Arcadeum *ArcadeumSession) SubkeyParentXXX(subkey common.Address, subkeyV uint8, subkeyR [32]byte, subkeyS [32]byte) (common.Address, error) {
	return _Arcadeum.Contract.SubkeyParentXXX(&_Arcadeum.CallOpts, subkey, subkeyV, subkeyR, subkeyS)
}

// SubkeyParentXXX is a free data retrieval call binding the contract method 0x5e15d8a8.
//
// Solidity: function subkeyParentXXX(subkey address, subkeyV uint8, subkeyR bytes32, subkeyS bytes32) constant returns(address)
func (_Arcadeum *ArcadeumCallerSession) SubkeyParentXXX(subkey common.Address, subkeyV uint8, subkeyR [32]byte, subkeyS [32]byte) (common.Address, error) {
	return _Arcadeum.Contract.SubkeyParentXXX(&_Arcadeum.CallOpts, subkey, subkeyV, subkeyR, subkeyS)
}

// TimestampSubkeyXXX is a free data retrieval call binding the contract method 0x529fec07.
//
// Solidity: function timestampSubkeyXXX(timestamp uint256, timestampV uint8, timestampR bytes32, timestampS bytes32) constant returns(address)
func (_Arcadeum *ArcadeumCaller) TimestampSubkeyXXX(opts *bind.CallOpts, timestamp *big.Int, timestampV uint8, timestampR [32]byte, timestampS [32]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Arcadeum.contract.Call(opts, out, "timestampSubkeyXXX", timestamp, timestampV, timestampR, timestampS)
	return *ret0, err
}

// TimestampSubkeyXXX is a free data retrieval call binding the contract method 0x529fec07.
//
// Solidity: function timestampSubkeyXXX(timestamp uint256, timestampV uint8, timestampR bytes32, timestampS bytes32) constant returns(address)
func (_Arcadeum *ArcadeumSession) TimestampSubkeyXXX(timestamp *big.Int, timestampV uint8, timestampR [32]byte, timestampS [32]byte) (common.Address, error) {
	return _Arcadeum.Contract.TimestampSubkeyXXX(&_Arcadeum.CallOpts, timestamp, timestampV, timestampR, timestampS)
}

// TimestampSubkeyXXX is a free data retrieval call binding the contract method 0x529fec07.
//
// Solidity: function timestampSubkeyXXX(timestamp uint256, timestampV uint8, timestampR bytes32, timestampS bytes32) constant returns(address)
func (_Arcadeum *ArcadeumCallerSession) TimestampSubkeyXXX(timestamp *big.Int, timestampV uint8, timestampR [32]byte, timestampS [32]byte) (common.Address, error) {
	return _Arcadeum.Contract.TimestampSubkeyXXX(&_Arcadeum.CallOpts, timestamp, timestampV, timestampR, timestampS)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_Arcadeum *ArcadeumTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Arcadeum.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_Arcadeum *ArcadeumSession) Deposit() (*types.Transaction, error) {
	return _Arcadeum.Contract.Deposit(&_Arcadeum.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_Arcadeum *ArcadeumTransactorSession) Deposit() (*types.Transaction, error) {
	return _Arcadeum.Contract.Deposit(&_Arcadeum.TransactOpts)
}

// FinishWithdrawal is a paid mutator transaction binding the contract method 0xbde6cf64.
//
// Solidity: function finishWithdrawal() returns()
func (_Arcadeum *ArcadeumTransactor) FinishWithdrawal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Arcadeum.contract.Transact(opts, "finishWithdrawal")
}

// FinishWithdrawal is a paid mutator transaction binding the contract method 0xbde6cf64.
//
// Solidity: function finishWithdrawal() returns()
func (_Arcadeum *ArcadeumSession) FinishWithdrawal() (*types.Transaction, error) {
	return _Arcadeum.Contract.FinishWithdrawal(&_Arcadeum.TransactOpts)
}

// FinishWithdrawal is a paid mutator transaction binding the contract method 0xbde6cf64.
//
// Solidity: function finishWithdrawal() returns()
func (_Arcadeum *ArcadeumTransactorSession) FinishWithdrawal() (*types.Transaction, error) {
	return _Arcadeum.Contract.FinishWithdrawal(&_Arcadeum.TransactOpts)
}

// StartWithdrawal is a paid mutator transaction binding the contract method 0xbc2f8dd8.
//
// Solidity: function startWithdrawal() returns()
func (_Arcadeum *ArcadeumTransactor) StartWithdrawal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Arcadeum.contract.Transact(opts, "startWithdrawal")
}

// StartWithdrawal is a paid mutator transaction binding the contract method 0xbc2f8dd8.
//
// Solidity: function startWithdrawal() returns()
func (_Arcadeum *ArcadeumSession) StartWithdrawal() (*types.Transaction, error) {
	return _Arcadeum.Contract.StartWithdrawal(&_Arcadeum.TransactOpts)
}

// StartWithdrawal is a paid mutator transaction binding the contract method 0xbc2f8dd8.
//
// Solidity: function startWithdrawal() returns()
func (_Arcadeum *ArcadeumTransactorSession) StartWithdrawal() (*types.Transaction, error) {
	return _Arcadeum.Contract.StartWithdrawal(&_Arcadeum.TransactOpts)
}

// StopWithdrawalXXX is a paid mutator transaction binding the contract method 0x927e0d56.
//
// Solidity: function stopWithdrawalXXX(game address, matchID uint32, timestamp uint256, timestampV uint8, timestampR bytes32, timestampS bytes32, subkeyV uint8, subkeyR bytes32, subkeyS bytes32) returns()
func (_Arcadeum *ArcadeumTransactor) StopWithdrawalXXX(opts *bind.TransactOpts, game common.Address, matchID uint32, timestamp *big.Int, timestampV uint8, timestampR [32]byte, timestampS [32]byte, subkeyV uint8, subkeyR [32]byte, subkeyS [32]byte) (*types.Transaction, error) {
	return _Arcadeum.contract.Transact(opts, "stopWithdrawalXXX", game, matchID, timestamp, timestampV, timestampR, timestampS, subkeyV, subkeyR, subkeyS)
}

// StopWithdrawalXXX is a paid mutator transaction binding the contract method 0x927e0d56.
//
// Solidity: function stopWithdrawalXXX(game address, matchID uint32, timestamp uint256, timestampV uint8, timestampR bytes32, timestampS bytes32, subkeyV uint8, subkeyR bytes32, subkeyS bytes32) returns()
func (_Arcadeum *ArcadeumSession) StopWithdrawalXXX(game common.Address, matchID uint32, timestamp *big.Int, timestampV uint8, timestampR [32]byte, timestampS [32]byte, subkeyV uint8, subkeyR [32]byte, subkeyS [32]byte) (*types.Transaction, error) {
	return _Arcadeum.Contract.StopWithdrawalXXX(&_Arcadeum.TransactOpts, game, matchID, timestamp, timestampV, timestampR, timestampS, subkeyV, subkeyR, subkeyS)
}

// StopWithdrawalXXX is a paid mutator transaction binding the contract method 0x927e0d56.
//
// Solidity: function stopWithdrawalXXX(game address, matchID uint32, timestamp uint256, timestampV uint8, timestampR bytes32, timestampS bytes32, subkeyV uint8, subkeyR bytes32, subkeyS bytes32) returns()
func (_Arcadeum *ArcadeumTransactorSession) StopWithdrawalXXX(game common.Address, matchID uint32, timestamp *big.Int, timestampV uint8, timestampR [32]byte, timestampS [32]byte, subkeyV uint8, subkeyR [32]byte, subkeyS [32]byte) (*types.Transaction, error) {
	return _Arcadeum.Contract.StopWithdrawalXXX(&_Arcadeum.TransactOpts, game, matchID, timestamp, timestampV, timestampR, timestampS, subkeyV, subkeyR, subkeyS)
}

// ArcadeumBalanceChangedIterator is returned from FilterBalanceChanged and is used to iterate over the raw logs and unpacked data for BalanceChanged events raised by the Arcadeum contract.
type ArcadeumBalanceChangedIterator struct {
	Event *ArcadeumBalanceChanged // Event containing the contract specifics and raw log

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
func (it *ArcadeumBalanceChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArcadeumBalanceChanged)
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
		it.Event = new(ArcadeumBalanceChanged)
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
func (it *ArcadeumBalanceChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ArcadeumBalanceChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ArcadeumBalanceChanged represents a BalanceChanged event raised by the Arcadeum contract.
type ArcadeumBalanceChanged struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBalanceChanged is a free log retrieval operation binding the contract event 0xeadc109e76de61f30acdeb8f317d29d2bfc33fb375d996835e421c72aabb7170.
//
// Solidity: event balanceChanged(account indexed address)
func (_Arcadeum *ArcadeumFilterer) FilterBalanceChanged(opts *bind.FilterOpts, account []common.Address) (*ArcadeumBalanceChangedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Arcadeum.contract.FilterLogs(opts, "balanceChanged", accountRule)
	if err != nil {
		return nil, err
	}
	return &ArcadeumBalanceChangedIterator{contract: _Arcadeum.contract, event: "balanceChanged", logs: logs, sub: sub}, nil
}

// WatchBalanceChanged is a free log subscription operation binding the contract event 0xeadc109e76de61f30acdeb8f317d29d2bfc33fb375d996835e421c72aabb7170.
//
// Solidity: event balanceChanged(account indexed address)
func (_Arcadeum *ArcadeumFilterer) WatchBalanceChanged(opts *bind.WatchOpts, sink chan<- *ArcadeumBalanceChanged, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Arcadeum.contract.WatchLogs(opts, "balanceChanged", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ArcadeumBalanceChanged)
				if err := _Arcadeum.contract.UnpackLog(event, "balanceChanged", log); err != nil {
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

// ArcadeumCheaterReportedIterator is returned from FilterCheaterReported and is used to iterate over the raw logs and unpacked data for CheaterReported events raised by the Arcadeum contract.
type ArcadeumCheaterReportedIterator struct {
	Event *ArcadeumCheaterReported // Event containing the contract specifics and raw log

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
func (it *ArcadeumCheaterReportedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArcadeumCheaterReported)
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
		it.Event = new(ArcadeumCheaterReported)
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
func (it *ArcadeumCheaterReportedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ArcadeumCheaterReportedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ArcadeumCheaterReported represents a CheaterReported event raised by the Arcadeum contract.
type ArcadeumCheaterReported struct {
	Game    common.Address
	MatchID uint32
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterCheaterReported is a free log retrieval operation binding the contract event 0xdfa829d13be375f4a0aefc98a25c8a11d716616e5210e6a69d6a8f342d4dbd38.
//
// Solidity: event cheaterReported(game indexed address, matchID indexed uint32, account indexed address)
func (_Arcadeum *ArcadeumFilterer) FilterCheaterReported(opts *bind.FilterOpts, game []common.Address, matchID []uint32, account []common.Address) (*ArcadeumCheaterReportedIterator, error) {

	var gameRule []interface{}
	for _, gameItem := range game {
		gameRule = append(gameRule, gameItem)
	}
	var matchIDRule []interface{}
	for _, matchIDItem := range matchID {
		matchIDRule = append(matchIDRule, matchIDItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Arcadeum.contract.FilterLogs(opts, "cheaterReported", gameRule, matchIDRule, accountRule)
	if err != nil {
		return nil, err
	}
	return &ArcadeumCheaterReportedIterator{contract: _Arcadeum.contract, event: "cheaterReported", logs: logs, sub: sub}, nil
}

// WatchCheaterReported is a free log subscription operation binding the contract event 0xdfa829d13be375f4a0aefc98a25c8a11d716616e5210e6a69d6a8f342d4dbd38.
//
// Solidity: event cheaterReported(game indexed address, matchID indexed uint32, account indexed address)
func (_Arcadeum *ArcadeumFilterer) WatchCheaterReported(opts *bind.WatchOpts, sink chan<- *ArcadeumCheaterReported, game []common.Address, matchID []uint32, account []common.Address) (event.Subscription, error) {

	var gameRule []interface{}
	for _, gameItem := range game {
		gameRule = append(gameRule, gameItem)
	}
	var matchIDRule []interface{}
	for _, matchIDItem := range matchID {
		matchIDRule = append(matchIDRule, matchIDItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Arcadeum.contract.WatchLogs(opts, "cheaterReported", gameRule, matchIDRule, accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ArcadeumCheaterReported)
				if err := _Arcadeum.contract.UnpackLog(event, "cheaterReported", log); err != nil {
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

// ArcadeumRewardClaimedIterator is returned from FilterRewardClaimed and is used to iterate over the raw logs and unpacked data for RewardClaimed events raised by the Arcadeum contract.
type ArcadeumRewardClaimedIterator struct {
	Event *ArcadeumRewardClaimed // Event containing the contract specifics and raw log

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
func (it *ArcadeumRewardClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArcadeumRewardClaimed)
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
		it.Event = new(ArcadeumRewardClaimed)
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
func (it *ArcadeumRewardClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ArcadeumRewardClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ArcadeumRewardClaimed represents a RewardClaimed event raised by the Arcadeum contract.
type ArcadeumRewardClaimed struct {
	Game    common.Address
	MatchID uint32
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRewardClaimed is a free log retrieval operation binding the contract event 0x10dfa66ec4ab3387564a3de672c82004e490e649ee4b23ae31263fed3bb87e65.
//
// Solidity: event rewardClaimed(game indexed address, matchID indexed uint32, account indexed address)
func (_Arcadeum *ArcadeumFilterer) FilterRewardClaimed(opts *bind.FilterOpts, game []common.Address, matchID []uint32, account []common.Address) (*ArcadeumRewardClaimedIterator, error) {

	var gameRule []interface{}
	for _, gameItem := range game {
		gameRule = append(gameRule, gameItem)
	}
	var matchIDRule []interface{}
	for _, matchIDItem := range matchID {
		matchIDRule = append(matchIDRule, matchIDItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Arcadeum.contract.FilterLogs(opts, "rewardClaimed", gameRule, matchIDRule, accountRule)
	if err != nil {
		return nil, err
	}
	return &ArcadeumRewardClaimedIterator{contract: _Arcadeum.contract, event: "rewardClaimed", logs: logs, sub: sub}, nil
}

// WatchRewardClaimed is a free log subscription operation binding the contract event 0x10dfa66ec4ab3387564a3de672c82004e490e649ee4b23ae31263fed3bb87e65.
//
// Solidity: event rewardClaimed(game indexed address, matchID indexed uint32, account indexed address)
func (_Arcadeum *ArcadeumFilterer) WatchRewardClaimed(opts *bind.WatchOpts, sink chan<- *ArcadeumRewardClaimed, game []common.Address, matchID []uint32, account []common.Address) (event.Subscription, error) {

	var gameRule []interface{}
	for _, gameItem := range game {
		gameRule = append(gameRule, gameItem)
	}
	var matchIDRule []interface{}
	for _, matchIDItem := range matchID {
		matchIDRule = append(matchIDRule, matchIDItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Arcadeum.contract.WatchLogs(opts, "rewardClaimed", gameRule, matchIDRule, accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ArcadeumRewardClaimed)
				if err := _Arcadeum.contract.UnpackLog(event, "rewardClaimed", log); err != nil {
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

// ArcadeumWithdrawalStartedIterator is returned from FilterWithdrawalStarted and is used to iterate over the raw logs and unpacked data for WithdrawalStarted events raised by the Arcadeum contract.
type ArcadeumWithdrawalStartedIterator struct {
	Event *ArcadeumWithdrawalStarted // Event containing the contract specifics and raw log

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
func (it *ArcadeumWithdrawalStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArcadeumWithdrawalStarted)
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
		it.Event = new(ArcadeumWithdrawalStarted)
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
func (it *ArcadeumWithdrawalStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ArcadeumWithdrawalStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ArcadeumWithdrawalStarted represents a WithdrawalStarted event raised by the Arcadeum contract.
type ArcadeumWithdrawalStarted struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalStarted is a free log retrieval operation binding the contract event 0x3ab48c8f15bc703088e6dd5631e3d30792f37e03cbdf5011637a8915b0d0bb47.
//
// Solidity: event withdrawalStarted(account indexed address)
func (_Arcadeum *ArcadeumFilterer) FilterWithdrawalStarted(opts *bind.FilterOpts, account []common.Address) (*ArcadeumWithdrawalStartedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Arcadeum.contract.FilterLogs(opts, "withdrawalStarted", accountRule)
	if err != nil {
		return nil, err
	}
	return &ArcadeumWithdrawalStartedIterator{contract: _Arcadeum.contract, event: "withdrawalStarted", logs: logs, sub: sub}, nil
}

// WatchWithdrawalStarted is a free log subscription operation binding the contract event 0x3ab48c8f15bc703088e6dd5631e3d30792f37e03cbdf5011637a8915b0d0bb47.
//
// Solidity: event withdrawalStarted(account indexed address)
func (_Arcadeum *ArcadeumFilterer) WatchWithdrawalStarted(opts *bind.WatchOpts, sink chan<- *ArcadeumWithdrawalStarted, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Arcadeum.contract.WatchLogs(opts, "withdrawalStarted", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ArcadeumWithdrawalStarted)
				if err := _Arcadeum.contract.UnpackLog(event, "withdrawalStarted", log); err != nil {
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

// ArcadeumWithdrawalStoppedIterator is returned from FilterWithdrawalStopped and is used to iterate over the raw logs and unpacked data for WithdrawalStopped events raised by the Arcadeum contract.
type ArcadeumWithdrawalStoppedIterator struct {
	Event *ArcadeumWithdrawalStopped // Event containing the contract specifics and raw log

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
func (it *ArcadeumWithdrawalStoppedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ArcadeumWithdrawalStopped)
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
		it.Event = new(ArcadeumWithdrawalStopped)
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
func (it *ArcadeumWithdrawalStoppedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ArcadeumWithdrawalStoppedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ArcadeumWithdrawalStopped represents a WithdrawalStopped event raised by the Arcadeum contract.
type ArcadeumWithdrawalStopped struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWithdrawalStopped is a free log retrieval operation binding the contract event 0xd002c0e157290359441a1dfa197f28329220c0e1d378d1974e9a87d3a6480633.
//
// Solidity: event withdrawalStopped(account indexed address)
func (_Arcadeum *ArcadeumFilterer) FilterWithdrawalStopped(opts *bind.FilterOpts, account []common.Address) (*ArcadeumWithdrawalStoppedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Arcadeum.contract.FilterLogs(opts, "withdrawalStopped", accountRule)
	if err != nil {
		return nil, err
	}
	return &ArcadeumWithdrawalStoppedIterator{contract: _Arcadeum.contract, event: "withdrawalStopped", logs: logs, sub: sub}, nil
}

// WatchWithdrawalStopped is a free log subscription operation binding the contract event 0xd002c0e157290359441a1dfa197f28329220c0e1d378d1974e9a87d3a6480633.
//
// Solidity: event withdrawalStopped(account indexed address)
func (_Arcadeum *ArcadeumFilterer) WatchWithdrawalStopped(opts *bind.WatchOpts, sink chan<- *ArcadeumWithdrawalStopped, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Arcadeum.contract.WatchLogs(opts, "withdrawalStopped", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ArcadeumWithdrawalStopped)
				if err := _Arcadeum.contract.UnpackLog(event, "withdrawalStopped", log); err != nil {
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

// DGameABI is the input ABI used to generate the binding from.
const DGameABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"matchDuration\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"},{\"name\":\"secretSeed\",\"type\":\"bytes\"}],\"name\":\"isSecretSeedValid\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"secretSeed\",\"type\":\"bytes\"}],\"name\":\"publicSeed\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[1]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"secretSeed\",\"type\":\"bytes\"}],\"name\":\"secretSeedRating\",\"outputs\":[{\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// DGameBin is the compiled bytecode used for deploying new contracts.
const DGameBin = `0x`

// DeployDGame deploys a new Ethereum contract, binding an instance of DGame to it.
func DeployDGame(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DGame, error) {
	parsed, err := abi.JSON(strings.NewReader(DGameABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(DGameBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DGame{DGameCaller: DGameCaller{contract: contract}, DGameTransactor: DGameTransactor{contract: contract}, DGameFilterer: DGameFilterer{contract: contract}}, nil
}

// DGame is an auto generated Go binding around an Ethereum contract.
type DGame struct {
	DGameCaller     // Read-only binding to the contract
	DGameTransactor // Write-only binding to the contract
	DGameFilterer   // Log filterer for contract events
}

// DGameCaller is an auto generated read-only Go binding around an Ethereum contract.
type DGameCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DGameTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DGameTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DGameFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DGameFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DGameSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DGameSession struct {
	Contract     *DGame            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DGameCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DGameCallerSession struct {
	Contract *DGameCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// DGameTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DGameTransactorSession struct {
	Contract     *DGameTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DGameRaw is an auto generated low-level Go binding around an Ethereum contract.
type DGameRaw struct {
	Contract *DGame // Generic contract binding to access the raw methods on
}

// DGameCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DGameCallerRaw struct {
	Contract *DGameCaller // Generic read-only contract binding to access the raw methods on
}

// DGameTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DGameTransactorRaw struct {
	Contract *DGameTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDGame creates a new instance of DGame, bound to a specific deployed contract.
func NewDGame(address common.Address, backend bind.ContractBackend) (*DGame, error) {
	contract, err := bindDGame(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DGame{DGameCaller: DGameCaller{contract: contract}, DGameTransactor: DGameTransactor{contract: contract}, DGameFilterer: DGameFilterer{contract: contract}}, nil
}

// NewDGameCaller creates a new read-only instance of DGame, bound to a specific deployed contract.
func NewDGameCaller(address common.Address, caller bind.ContractCaller) (*DGameCaller, error) {
	contract, err := bindDGame(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DGameCaller{contract: contract}, nil
}

// NewDGameTransactor creates a new write-only instance of DGame, bound to a specific deployed contract.
func NewDGameTransactor(address common.Address, transactor bind.ContractTransactor) (*DGameTransactor, error) {
	contract, err := bindDGame(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DGameTransactor{contract: contract}, nil
}

// NewDGameFilterer creates a new log filterer instance of DGame, bound to a specific deployed contract.
func NewDGameFilterer(address common.Address, filterer bind.ContractFilterer) (*DGameFilterer, error) {
	contract, err := bindDGame(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DGameFilterer{contract: contract}, nil
}

// bindDGame binds a generic wrapper to an already deployed contract.
func bindDGame(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DGameABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DGame *DGameRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DGame.Contract.DGameCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DGame *DGameRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DGame.Contract.DGameTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DGame *DGameRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DGame.Contract.DGameTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DGame *DGameCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DGame.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DGame *DGameTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DGame.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DGame *DGameTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DGame.Contract.contract.Transact(opts, method, params...)
}

// IsSecretSeedValid is a free data retrieval call binding the contract method 0x67745882.
//
// Solidity: function isSecretSeedValid(account address, secretSeed bytes) constant returns(bool)
func (_DGame *DGameCaller) IsSecretSeedValid(opts *bind.CallOpts, account common.Address, secretSeed []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DGame.contract.Call(opts, out, "isSecretSeedValid", account, secretSeed)
	return *ret0, err
}

// IsSecretSeedValid is a free data retrieval call binding the contract method 0x67745882.
//
// Solidity: function isSecretSeedValid(account address, secretSeed bytes) constant returns(bool)
func (_DGame *DGameSession) IsSecretSeedValid(account common.Address, secretSeed []byte) (bool, error) {
	return _DGame.Contract.IsSecretSeedValid(&_DGame.CallOpts, account, secretSeed)
}

// IsSecretSeedValid is a free data retrieval call binding the contract method 0x67745882.
//
// Solidity: function isSecretSeedValid(account address, secretSeed bytes) constant returns(bool)
func (_DGame *DGameCallerSession) IsSecretSeedValid(account common.Address, secretSeed []byte) (bool, error) {
	return _DGame.Contract.IsSecretSeedValid(&_DGame.CallOpts, account, secretSeed)
}

// MatchDuration is a free data retrieval call binding the contract method 0x0e649a0b.
//
// Solidity: function matchDuration() constant returns(uint256)
func (_DGame *DGameCaller) MatchDuration(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DGame.contract.Call(opts, out, "matchDuration")
	return *ret0, err
}

// MatchDuration is a free data retrieval call binding the contract method 0x0e649a0b.
//
// Solidity: function matchDuration() constant returns(uint256)
func (_DGame *DGameSession) MatchDuration() (*big.Int, error) {
	return _DGame.Contract.MatchDuration(&_DGame.CallOpts)
}

// MatchDuration is a free data retrieval call binding the contract method 0x0e649a0b.
//
// Solidity: function matchDuration() constant returns(uint256)
func (_DGame *DGameCallerSession) MatchDuration() (*big.Int, error) {
	return _DGame.Contract.MatchDuration(&_DGame.CallOpts)
}

// PublicSeed is a free data retrieval call binding the contract method 0x9bf391ec.
//
// Solidity: function publicSeed(secretSeed bytes) constant returns(bytes32[1])
func (_DGame *DGameCaller) PublicSeed(opts *bind.CallOpts, secretSeed []byte) ([1][32]byte, error) {
	var (
		ret0 = new([1][32]byte)
	)
	out := ret0
	err := _DGame.contract.Call(opts, out, "publicSeed", secretSeed)
	return *ret0, err
}

// PublicSeed is a free data retrieval call binding the contract method 0x9bf391ec.
//
// Solidity: function publicSeed(secretSeed bytes) constant returns(bytes32[1])
func (_DGame *DGameSession) PublicSeed(secretSeed []byte) ([1][32]byte, error) {
	return _DGame.Contract.PublicSeed(&_DGame.CallOpts, secretSeed)
}

// PublicSeed is a free data retrieval call binding the contract method 0x9bf391ec.
//
// Solidity: function publicSeed(secretSeed bytes) constant returns(bytes32[1])
func (_DGame *DGameCallerSession) PublicSeed(secretSeed []byte) ([1][32]byte, error) {
	return _DGame.Contract.PublicSeed(&_DGame.CallOpts, secretSeed)
}

// SecretSeedRating is a free data retrieval call binding the contract method 0x9ef7f8d4.
//
// Solidity: function secretSeedRating(secretSeed bytes) constant returns(uint32)
func (_DGame *DGameCaller) SecretSeedRating(opts *bind.CallOpts, secretSeed []byte) (uint32, error) {
	var (
		ret0 = new(uint32)
	)
	out := ret0
	err := _DGame.contract.Call(opts, out, "secretSeedRating", secretSeed)
	return *ret0, err
}

// SecretSeedRating is a free data retrieval call binding the contract method 0x9ef7f8d4.
//
// Solidity: function secretSeedRating(secretSeed bytes) constant returns(uint32)
func (_DGame *DGameSession) SecretSeedRating(secretSeed []byte) (uint32, error) {
	return _DGame.Contract.SecretSeedRating(&_DGame.CallOpts, secretSeed)
}

// SecretSeedRating is a free data retrieval call binding the contract method 0x9ef7f8d4.
//
// Solidity: function secretSeedRating(secretSeed bytes) constant returns(uint32)
func (_DGame *DGameCallerSession) SecretSeedRating(secretSeed []byte) (uint32, error) {
	return _DGame.Contract.SecretSeedRating(&_DGame.CallOpts, secretSeed)
}
