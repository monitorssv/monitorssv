package utils

import (
	"errors"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"math/big"
	"strings"
)

var (
	// https://etherscan.io/address/0xeefba1e63905ef1d7acba5a8513c70307c1ce441#contracts
	mainnetMultiCallAddr = common.HexToAddress("0xeefba1e63905ef1d7acba5a8513c70307c1ce441")
	// https://holesky.etherscan.io/address/0xa3e09ba95fa2887b351043841fe5d0cc33ff1052#contracts
	holeskyMultiCallAddr = common.HexToAddress("0xA3e09ba95Fa2887b351043841Fe5d0CC33fF1052")
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

// Struct0 is an auto generated low-level Go binding around an user-defined struct.
type Struct0 struct {
	Target   common.Address
	CallData []byte
}

// MulticallMetaData contains all meta data concerning the Multicall contract.
var MulticallMetaData = &bind.MetaData{
	ABI: "[{\"constant\":false,\"inputs\":[{\"components\":[{\"name\":\"target\",\"type\":\"address\"},{\"name\":\"callData\",\"type\":\"bytes\"}],\"name\":\"calls\",\"type\":\"tuple[]\"}],\"name\":\"aggregate\",\"outputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"name\":\"returnData\",\"type\":\"bytes[]\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// MulticallABI is the input ABI used to generate the binding from.
// Deprecated: Use MulticallMetaData.ABI instead.
var MulticallABI = MulticallMetaData.ABI

// Multicall is an auto generated Go binding around an Ethereum contract.
type Multicall struct {
	MulticallCaller // Read-only binding to the contract
}

// MulticallCaller is an auto generated read-only Go binding around an Ethereum contract.
type MulticallCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NewMulticall creates a new instance of Multicall, bound to a specific deployed contract.
func NewMulticall(network string, backend bind.ContractBackend) (*Multicall, error) {
	address := mainnetMultiCallAddr
	if strings.ToLower(network) == "holesky" {
		address = holeskyMultiCallAddr
	}

	contract, err := bindMulticall(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Multicall{MulticallCaller: MulticallCaller{contract: contract}}, nil
}

// bindMulticall binds a generic wrapper to an already deployed contract.
func bindMulticall(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MulticallABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Aggregate : Read-only
// Solidity: function aggregate((address,bytes)[] calls) returns(uint256 blockNumber, bytes[] returnData)
func (mc *MulticallCaller) Aggregate(opts *bind.CallOpts, calls []Struct0) ([]interface{}, error) {
	var out []interface{}
	err := mc.contract.Call(opts, &out, "aggregate", calls)

	if err != nil {
		return nil, err
	}

	return out, err
}
