// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ssv

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

// ISSVNetworkCoreCluster is an auto generated low-level Go binding around an user-defined struct.
type ISSVNetworkCoreCluster struct {
	ValidatorCount  uint32
	NetworkFeeIndex uint64
	Index           uint64
	Active          bool
	Balance         *big.Int
}

// SsvMetaData contains all meta data concerning the Ssv contract.
var SsvMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"}],\"name\":\"AddressIsWhitelistingContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ApprovalNotWithinTimeframe\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CallerNotOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"CallerNotOwnerWithData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CallerNotWhitelisted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"operatorId\",\"type\":\"uint64\"}],\"name\":\"CallerNotWhitelistedWithData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClusterAlreadyEnabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClusterDoesNotExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClusterIsLiquidated\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ClusterNotLiquidatable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyPublicKeysList\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"operatorId\",\"type\":\"uint64\"}],\"name\":\"ExceedValidatorLimit\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"operatorId\",\"type\":\"uint64\"}],\"name\":\"ExceedValidatorLimitWithData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeExceedsIncreaseLimit\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeIncreaseNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeTooHigh\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FeeTooLow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectClusterState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"IncorrectValidatorState\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"}],\"name\":\"IncorrectValidatorStateWithData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidContractAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidOperatorIdsLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPublicKeyLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidWhitelistAddressesLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"}],\"name\":\"InvalidWhitelistingContract\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxValueExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewBlockPeriodIsBelowMinimum\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NoFeeDeclared\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotAuthorized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OperatorAlreadyExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OperatorDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OperatorsListNotUnique\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PublicKeysSharesLengthMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SameFeeChangeNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TargetModuleDoesNotExist\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"moduleId\",\"type\":\"uint8\"}],\"name\":\"TargetModuleDoesNotExistWithData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TokenTransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnsortedOperatorsList\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValidatorAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"}],\"name\":\"ValidatorAlreadyExistsWithData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ValidatorDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddressNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferStarted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"clusterOwner\",\"type\":\"address\"},{\"internalType\":\"uint64[]\",\"name\":\"operatorIds\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"validatorCount\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"networkFeeIndex\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structISSVNetworkCore.Cluster\",\"name\":\"cluster\",\"type\":\"tuple\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"clusterOwner\",\"type\":\"address\"},{\"internalType\":\"uint64[]\",\"name\":\"operatorIds\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"validatorCount\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"networkFeeIndex\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structISSVNetworkCore.Cluster\",\"name\":\"cluster\",\"type\":\"tuple\"}],\"name\":\"getBurnRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLiquidationThresholdPeriod\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaximumOperatorFee\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMinimumLiquidationCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNetworkEarnings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNetworkFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNetworkValidatorsCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"operatorId\",\"type\":\"uint64\"}],\"name\":\"getOperatorById\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"operatorId\",\"type\":\"uint64\"}],\"name\":\"getOperatorDeclaredFee\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"id\",\"type\":\"uint64\"}],\"name\":\"getOperatorEarnings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"operatorId\",\"type\":\"uint64\"}],\"name\":\"getOperatorFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperatorFeeIncreaseLimit\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOperatorFeePeriods\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"clusterOwner\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"}],\"name\":\"getValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidatorsPerOperatorLimit\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64[]\",\"name\":\"operatorIds\",\"type\":\"uint64[]\"},{\"internalType\":\"address\",\"name\":\"whitelistedAddress\",\"type\":\"address\"}],\"name\":\"getWhitelistedOperators\",\"outputs\":[{\"internalType\":\"uint64[]\",\"name\":\"whitelistedOperatorIds\",\"type\":\"uint64[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractISSVViews\",\"name\":\"ssvNetwork_\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addressToCheck\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"operatorId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"whitelistingContract\",\"type\":\"address\"}],\"name\":\"isAddressWhitelistedInWhitelistingContract\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isWhitelisted\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"clusterOwner\",\"type\":\"address\"},{\"internalType\":\"uint64[]\",\"name\":\"operatorIds\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"validatorCount\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"networkFeeIndex\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structISSVNetworkCore.Cluster\",\"name\":\"cluster\",\"type\":\"tuple\"}],\"name\":\"isLiquidatable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"clusterOwner\",\"type\":\"address\"},{\"internalType\":\"uint64[]\",\"name\":\"operatorIds\",\"type\":\"uint64[]\"},{\"components\":[{\"internalType\":\"uint32\",\"name\":\"validatorCount\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"networkFeeIndex\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"index\",\"type\":\"uint64\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structISSVNetworkCore.Cluster\",\"name\":\"cluster\",\"type\":\"tuple\"}],\"name\":\"isLiquidated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"contractAddress\",\"type\":\"address\"}],\"name\":\"isWhitelistingContract\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ssvNetwork\",\"outputs\":[{\"internalType\":\"contractISSVViews\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// SsvABI is the input ABI used to generate the binding from.
// Deprecated: Use SsvMetaData.ABI instead.
var SsvABI = SsvMetaData.ABI

// Ssv is an auto generated Go binding around an Ethereum contract.
type Ssv struct {
	SsvCaller     // Read-only binding to the contract
	SsvTransactor // Write-only binding to the contract
	SsvFilterer   // Log filterer for contract events
}

// SsvCaller is an auto generated read-only Go binding around an Ethereum contract.
type SsvCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SsvTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SsvTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SsvFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SsvFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SsvSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SsvSession struct {
	Contract     *Ssv              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SsvCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SsvCallerSession struct {
	Contract *SsvCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// SsvTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SsvTransactorSession struct {
	Contract     *SsvTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SsvRaw is an auto generated low-level Go binding around an Ethereum contract.
type SsvRaw struct {
	Contract *Ssv // Generic contract binding to access the raw methods on
}

// SsvCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SsvCallerRaw struct {
	Contract *SsvCaller // Generic read-only contract binding to access the raw methods on
}

// SsvTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SsvTransactorRaw struct {
	Contract *SsvTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSsv creates a new instance of Ssv, bound to a specific deployed contract.
func NewSsv(address common.Address, backend bind.ContractBackend) (*Ssv, error) {
	contract, err := bindSsv(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ssv{SsvCaller: SsvCaller{contract: contract}, SsvTransactor: SsvTransactor{contract: contract}, SsvFilterer: SsvFilterer{contract: contract}}, nil
}

// NewSsvCaller creates a new read-only instance of Ssv, bound to a specific deployed contract.
func NewSsvCaller(address common.Address, caller bind.ContractCaller) (*SsvCaller, error) {
	contract, err := bindSsv(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SsvCaller{contract: contract}, nil
}

// NewSsvTransactor creates a new write-only instance of Ssv, bound to a specific deployed contract.
func NewSsvTransactor(address common.Address, transactor bind.ContractTransactor) (*SsvTransactor, error) {
	contract, err := bindSsv(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SsvTransactor{contract: contract}, nil
}

// NewSsvFilterer creates a new log filterer instance of Ssv, bound to a specific deployed contract.
func NewSsvFilterer(address common.Address, filterer bind.ContractFilterer) (*SsvFilterer, error) {
	contract, err := bindSsv(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SsvFilterer{contract: contract}, nil
}

// bindSsv binds a generic wrapper to an already deployed contract.
func bindSsv(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SsvMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ssv *SsvRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ssv.Contract.SsvCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ssv *SsvRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ssv.Contract.SsvTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ssv *SsvRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ssv.Contract.SsvTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ssv *SsvCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ssv.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ssv *SsvTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ssv.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ssv *SsvTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ssv.Contract.contract.Transact(opts, method, params...)
}

// GetBalance is a free data retrieval call binding the contract method 0xeb8ecfa7.
//
// Solidity: function getBalance(address clusterOwner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) view returns(uint256)
func (_Ssv *SsvCaller) GetBalance(opts *bind.CallOpts, clusterOwner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (*big.Int, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getBalance", clusterOwner, operatorIds, cluster)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalance is a free data retrieval call binding the contract method 0xeb8ecfa7.
//
// Solidity: function getBalance(address clusterOwner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) view returns(uint256)
func (_Ssv *SsvSession) GetBalance(clusterOwner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (*big.Int, error) {
	return _Ssv.Contract.GetBalance(&_Ssv.CallOpts, clusterOwner, operatorIds, cluster)
}

// GetBalance is a free data retrieval call binding the contract method 0xeb8ecfa7.
//
// Solidity: function getBalance(address clusterOwner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) view returns(uint256)
func (_Ssv *SsvCallerSession) GetBalance(clusterOwner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (*big.Int, error) {
	return _Ssv.Contract.GetBalance(&_Ssv.CallOpts, clusterOwner, operatorIds, cluster)
}

// GetBurnRate is a free data retrieval call binding the contract method 0xca162e5e.
//
// Solidity: function getBurnRate(address clusterOwner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) view returns(uint256)
func (_Ssv *SsvCaller) GetBurnRate(opts *bind.CallOpts, clusterOwner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (*big.Int, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getBurnRate", clusterOwner, operatorIds, cluster)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBurnRate is a free data retrieval call binding the contract method 0xca162e5e.
//
// Solidity: function getBurnRate(address clusterOwner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) view returns(uint256)
func (_Ssv *SsvSession) GetBurnRate(clusterOwner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (*big.Int, error) {
	return _Ssv.Contract.GetBurnRate(&_Ssv.CallOpts, clusterOwner, operatorIds, cluster)
}

// GetBurnRate is a free data retrieval call binding the contract method 0xca162e5e.
//
// Solidity: function getBurnRate(address clusterOwner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) view returns(uint256)
func (_Ssv *SsvCallerSession) GetBurnRate(clusterOwner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (*big.Int, error) {
	return _Ssv.Contract.GetBurnRate(&_Ssv.CallOpts, clusterOwner, operatorIds, cluster)
}

// GetLiquidationThresholdPeriod is a free data retrieval call binding the contract method 0x9040f7c3.
//
// Solidity: function getLiquidationThresholdPeriod() view returns(uint64)
func (_Ssv *SsvCaller) GetLiquidationThresholdPeriod(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getLiquidationThresholdPeriod")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetLiquidationThresholdPeriod is a free data retrieval call binding the contract method 0x9040f7c3.
//
// Solidity: function getLiquidationThresholdPeriod() view returns(uint64)
func (_Ssv *SsvSession) GetLiquidationThresholdPeriod() (uint64, error) {
	return _Ssv.Contract.GetLiquidationThresholdPeriod(&_Ssv.CallOpts)
}

// GetLiquidationThresholdPeriod is a free data retrieval call binding the contract method 0x9040f7c3.
//
// Solidity: function getLiquidationThresholdPeriod() view returns(uint64)
func (_Ssv *SsvCallerSession) GetLiquidationThresholdPeriod() (uint64, error) {
	return _Ssv.Contract.GetLiquidationThresholdPeriod(&_Ssv.CallOpts)
}

// GetMaximumOperatorFee is a free data retrieval call binding the contract method 0xdf02ef7f.
//
// Solidity: function getMaximumOperatorFee() view returns(uint64)
func (_Ssv *SsvCaller) GetMaximumOperatorFee(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getMaximumOperatorFee")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetMaximumOperatorFee is a free data retrieval call binding the contract method 0xdf02ef7f.
//
// Solidity: function getMaximumOperatorFee() view returns(uint64)
func (_Ssv *SsvSession) GetMaximumOperatorFee() (uint64, error) {
	return _Ssv.Contract.GetMaximumOperatorFee(&_Ssv.CallOpts)
}

// GetMaximumOperatorFee is a free data retrieval call binding the contract method 0xdf02ef7f.
//
// Solidity: function getMaximumOperatorFee() view returns(uint64)
func (_Ssv *SsvCallerSession) GetMaximumOperatorFee() (uint64, error) {
	return _Ssv.Contract.GetMaximumOperatorFee(&_Ssv.CallOpts)
}

// GetMinimumLiquidationCollateral is a free data retrieval call binding the contract method 0x5ba3d62a.
//
// Solidity: function getMinimumLiquidationCollateral() view returns(uint256)
func (_Ssv *SsvCaller) GetMinimumLiquidationCollateral(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getMinimumLiquidationCollateral")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMinimumLiquidationCollateral is a free data retrieval call binding the contract method 0x5ba3d62a.
//
// Solidity: function getMinimumLiquidationCollateral() view returns(uint256)
func (_Ssv *SsvSession) GetMinimumLiquidationCollateral() (*big.Int, error) {
	return _Ssv.Contract.GetMinimumLiquidationCollateral(&_Ssv.CallOpts)
}

// GetMinimumLiquidationCollateral is a free data retrieval call binding the contract method 0x5ba3d62a.
//
// Solidity: function getMinimumLiquidationCollateral() view returns(uint256)
func (_Ssv *SsvCallerSession) GetMinimumLiquidationCollateral() (*big.Int, error) {
	return _Ssv.Contract.GetMinimumLiquidationCollateral(&_Ssv.CallOpts)
}

// GetNetworkEarnings is a free data retrieval call binding the contract method 0x777915cb.
//
// Solidity: function getNetworkEarnings() view returns(uint256)
func (_Ssv *SsvCaller) GetNetworkEarnings(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getNetworkEarnings")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNetworkEarnings is a free data retrieval call binding the contract method 0x777915cb.
//
// Solidity: function getNetworkEarnings() view returns(uint256)
func (_Ssv *SsvSession) GetNetworkEarnings() (*big.Int, error) {
	return _Ssv.Contract.GetNetworkEarnings(&_Ssv.CallOpts)
}

// GetNetworkEarnings is a free data retrieval call binding the contract method 0x777915cb.
//
// Solidity: function getNetworkEarnings() view returns(uint256)
func (_Ssv *SsvCallerSession) GetNetworkEarnings() (*big.Int, error) {
	return _Ssv.Contract.GetNetworkEarnings(&_Ssv.CallOpts)
}

// GetNetworkFee is a free data retrieval call binding the contract method 0xfc043830.
//
// Solidity: function getNetworkFee() view returns(uint256)
func (_Ssv *SsvCaller) GetNetworkFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getNetworkFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNetworkFee is a free data retrieval call binding the contract method 0xfc043830.
//
// Solidity: function getNetworkFee() view returns(uint256)
func (_Ssv *SsvSession) GetNetworkFee() (*big.Int, error) {
	return _Ssv.Contract.GetNetworkFee(&_Ssv.CallOpts)
}

// GetNetworkFee is a free data retrieval call binding the contract method 0xfc043830.
//
// Solidity: function getNetworkFee() view returns(uint256)
func (_Ssv *SsvCallerSession) GetNetworkFee() (*big.Int, error) {
	return _Ssv.Contract.GetNetworkFee(&_Ssv.CallOpts)
}

// GetNetworkValidatorsCount is a free data retrieval call binding the contract method 0x9568f9d9.
//
// Solidity: function getNetworkValidatorsCount() view returns(uint32)
func (_Ssv *SsvCaller) GetNetworkValidatorsCount(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getNetworkValidatorsCount")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetNetworkValidatorsCount is a free data retrieval call binding the contract method 0x9568f9d9.
//
// Solidity: function getNetworkValidatorsCount() view returns(uint32)
func (_Ssv *SsvSession) GetNetworkValidatorsCount() (uint32, error) {
	return _Ssv.Contract.GetNetworkValidatorsCount(&_Ssv.CallOpts)
}

// GetNetworkValidatorsCount is a free data retrieval call binding the contract method 0x9568f9d9.
//
// Solidity: function getNetworkValidatorsCount() view returns(uint32)
func (_Ssv *SsvCallerSession) GetNetworkValidatorsCount() (uint32, error) {
	return _Ssv.Contract.GetNetworkValidatorsCount(&_Ssv.CallOpts)
}

// GetOperatorById is a free data retrieval call binding the contract method 0xbe3f058e.
//
// Solidity: function getOperatorById(uint64 operatorId) view returns(address, uint256, uint32, address, bool, bool)
func (_Ssv *SsvCaller) GetOperatorById(opts *bind.CallOpts, operatorId uint64) (common.Address, *big.Int, uint32, common.Address, bool, bool, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getOperatorById", operatorId)

	if err != nil {
		return *new(common.Address), *new(*big.Int), *new(uint32), *new(common.Address), *new(bool), *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(uint32)).(*uint32)
	out3 := *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	out4 := *abi.ConvertType(out[4], new(bool)).(*bool)
	out5 := *abi.ConvertType(out[5], new(bool)).(*bool)

	return out0, out1, out2, out3, out4, out5, err

}

// GetOperatorById is a free data retrieval call binding the contract method 0xbe3f058e.
//
// Solidity: function getOperatorById(uint64 operatorId) view returns(address, uint256, uint32, address, bool, bool)
func (_Ssv *SsvSession) GetOperatorById(operatorId uint64) (common.Address, *big.Int, uint32, common.Address, bool, bool, error) {
	return _Ssv.Contract.GetOperatorById(&_Ssv.CallOpts, operatorId)
}

// GetOperatorById is a free data retrieval call binding the contract method 0xbe3f058e.
//
// Solidity: function getOperatorById(uint64 operatorId) view returns(address, uint256, uint32, address, bool, bool)
func (_Ssv *SsvCallerSession) GetOperatorById(operatorId uint64) (common.Address, *big.Int, uint32, common.Address, bool, bool, error) {
	return _Ssv.Contract.GetOperatorById(&_Ssv.CallOpts, operatorId)
}

// GetOperatorDeclaredFee is a free data retrieval call binding the contract method 0x03b3d436.
//
// Solidity: function getOperatorDeclaredFee(uint64 operatorId) view returns(bool, uint256, uint64, uint64)
func (_Ssv *SsvCaller) GetOperatorDeclaredFee(opts *bind.CallOpts, operatorId uint64) (bool, *big.Int, uint64, uint64, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getOperatorDeclaredFee", operatorId)

	if err != nil {
		return *new(bool), *new(*big.Int), *new(uint64), *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(uint64)).(*uint64)
	out3 := *abi.ConvertType(out[3], new(uint64)).(*uint64)

	return out0, out1, out2, out3, err

}

// GetOperatorDeclaredFee is a free data retrieval call binding the contract method 0x03b3d436.
//
// Solidity: function getOperatorDeclaredFee(uint64 operatorId) view returns(bool, uint256, uint64, uint64)
func (_Ssv *SsvSession) GetOperatorDeclaredFee(operatorId uint64) (bool, *big.Int, uint64, uint64, error) {
	return _Ssv.Contract.GetOperatorDeclaredFee(&_Ssv.CallOpts, operatorId)
}

// GetOperatorDeclaredFee is a free data retrieval call binding the contract method 0x03b3d436.
//
// Solidity: function getOperatorDeclaredFee(uint64 operatorId) view returns(bool, uint256, uint64, uint64)
func (_Ssv *SsvCallerSession) GetOperatorDeclaredFee(operatorId uint64) (bool, *big.Int, uint64, uint64, error) {
	return _Ssv.Contract.GetOperatorDeclaredFee(&_Ssv.CallOpts, operatorId)
}

// GetOperatorEarnings is a free data retrieval call binding the contract method 0x6d0db0e4.
//
// Solidity: function getOperatorEarnings(uint64 id) view returns(uint256)
func (_Ssv *SsvCaller) GetOperatorEarnings(opts *bind.CallOpts, id uint64) (*big.Int, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getOperatorEarnings", id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOperatorEarnings is a free data retrieval call binding the contract method 0x6d0db0e4.
//
// Solidity: function getOperatorEarnings(uint64 id) view returns(uint256)
func (_Ssv *SsvSession) GetOperatorEarnings(id uint64) (*big.Int, error) {
	return _Ssv.Contract.GetOperatorEarnings(&_Ssv.CallOpts, id)
}

// GetOperatorEarnings is a free data retrieval call binding the contract method 0x6d0db0e4.
//
// Solidity: function getOperatorEarnings(uint64 id) view returns(uint256)
func (_Ssv *SsvCallerSession) GetOperatorEarnings(id uint64) (*big.Int, error) {
	return _Ssv.Contract.GetOperatorEarnings(&_Ssv.CallOpts, id)
}

// GetOperatorFee is a free data retrieval call binding the contract method 0x9ad3c745.
//
// Solidity: function getOperatorFee(uint64 operatorId) view returns(uint256)
func (_Ssv *SsvCaller) GetOperatorFee(opts *bind.CallOpts, operatorId uint64) (*big.Int, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getOperatorFee", operatorId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetOperatorFee is a free data retrieval call binding the contract method 0x9ad3c745.
//
// Solidity: function getOperatorFee(uint64 operatorId) view returns(uint256)
func (_Ssv *SsvSession) GetOperatorFee(operatorId uint64) (*big.Int, error) {
	return _Ssv.Contract.GetOperatorFee(&_Ssv.CallOpts, operatorId)
}

// GetOperatorFee is a free data retrieval call binding the contract method 0x9ad3c745.
//
// Solidity: function getOperatorFee(uint64 operatorId) view returns(uint256)
func (_Ssv *SsvCallerSession) GetOperatorFee(operatorId uint64) (*big.Int, error) {
	return _Ssv.Contract.GetOperatorFee(&_Ssv.CallOpts, operatorId)
}

// GetOperatorFeeIncreaseLimit is a free data retrieval call binding the contract method 0x68465f7d.
//
// Solidity: function getOperatorFeeIncreaseLimit() view returns(uint64)
func (_Ssv *SsvCaller) GetOperatorFeeIncreaseLimit(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getOperatorFeeIncreaseLimit")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetOperatorFeeIncreaseLimit is a free data retrieval call binding the contract method 0x68465f7d.
//
// Solidity: function getOperatorFeeIncreaseLimit() view returns(uint64)
func (_Ssv *SsvSession) GetOperatorFeeIncreaseLimit() (uint64, error) {
	return _Ssv.Contract.GetOperatorFeeIncreaseLimit(&_Ssv.CallOpts)
}

// GetOperatorFeeIncreaseLimit is a free data retrieval call binding the contract method 0x68465f7d.
//
// Solidity: function getOperatorFeeIncreaseLimit() view returns(uint64)
func (_Ssv *SsvCallerSession) GetOperatorFeeIncreaseLimit() (uint64, error) {
	return _Ssv.Contract.GetOperatorFeeIncreaseLimit(&_Ssv.CallOpts)
}

// GetOperatorFeePeriods is a free data retrieval call binding the contract method 0xe6d2834d.
//
// Solidity: function getOperatorFeePeriods() view returns(uint64, uint64)
func (_Ssv *SsvCaller) GetOperatorFeePeriods(opts *bind.CallOpts) (uint64, uint64, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getOperatorFeePeriods")

	if err != nil {
		return *new(uint64), *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)
	out1 := *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return out0, out1, err

}

// GetOperatorFeePeriods is a free data retrieval call binding the contract method 0xe6d2834d.
//
// Solidity: function getOperatorFeePeriods() view returns(uint64, uint64)
func (_Ssv *SsvSession) GetOperatorFeePeriods() (uint64, uint64, error) {
	return _Ssv.Contract.GetOperatorFeePeriods(&_Ssv.CallOpts)
}

// GetOperatorFeePeriods is a free data retrieval call binding the contract method 0xe6d2834d.
//
// Solidity: function getOperatorFeePeriods() view returns(uint64, uint64)
func (_Ssv *SsvCallerSession) GetOperatorFeePeriods() (uint64, uint64, error) {
	return _Ssv.Contract.GetOperatorFeePeriods(&_Ssv.CallOpts)
}

// GetValidator is a free data retrieval call binding the contract method 0x3e2ec160.
//
// Solidity: function getValidator(address clusterOwner, bytes publicKey) view returns(bool)
func (_Ssv *SsvCaller) GetValidator(opts *bind.CallOpts, clusterOwner common.Address, publicKey []byte) (bool, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getValidator", clusterOwner, publicKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetValidator is a free data retrieval call binding the contract method 0x3e2ec160.
//
// Solidity: function getValidator(address clusterOwner, bytes publicKey) view returns(bool)
func (_Ssv *SsvSession) GetValidator(clusterOwner common.Address, publicKey []byte) (bool, error) {
	return _Ssv.Contract.GetValidator(&_Ssv.CallOpts, clusterOwner, publicKey)
}

// GetValidator is a free data retrieval call binding the contract method 0x3e2ec160.
//
// Solidity: function getValidator(address clusterOwner, bytes publicKey) view returns(bool)
func (_Ssv *SsvCallerSession) GetValidator(clusterOwner common.Address, publicKey []byte) (bool, error) {
	return _Ssv.Contract.GetValidator(&_Ssv.CallOpts, clusterOwner, publicKey)
}

// GetValidatorsPerOperatorLimit is a free data retrieval call binding the contract method 0x14cb9d7b.
//
// Solidity: function getValidatorsPerOperatorLimit() view returns(uint32)
func (_Ssv *SsvCaller) GetValidatorsPerOperatorLimit(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getValidatorsPerOperatorLimit")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GetValidatorsPerOperatorLimit is a free data retrieval call binding the contract method 0x14cb9d7b.
//
// Solidity: function getValidatorsPerOperatorLimit() view returns(uint32)
func (_Ssv *SsvSession) GetValidatorsPerOperatorLimit() (uint32, error) {
	return _Ssv.Contract.GetValidatorsPerOperatorLimit(&_Ssv.CallOpts)
}

// GetValidatorsPerOperatorLimit is a free data retrieval call binding the contract method 0x14cb9d7b.
//
// Solidity: function getValidatorsPerOperatorLimit() view returns(uint32)
func (_Ssv *SsvCallerSession) GetValidatorsPerOperatorLimit() (uint32, error) {
	return _Ssv.Contract.GetValidatorsPerOperatorLimit(&_Ssv.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(string)
func (_Ssv *SsvCaller) GetVersion(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getVersion")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(string)
func (_Ssv *SsvSession) GetVersion() (string, error) {
	return _Ssv.Contract.GetVersion(&_Ssv.CallOpts)
}

// GetVersion is a free data retrieval call binding the contract method 0x0d8e6e2c.
//
// Solidity: function getVersion() view returns(string)
func (_Ssv *SsvCallerSession) GetVersion() (string, error) {
	return _Ssv.Contract.GetVersion(&_Ssv.CallOpts)
}

// GetWhitelistedOperators is a free data retrieval call binding the contract method 0xa9cf9eec.
//
// Solidity: function getWhitelistedOperators(uint64[] operatorIds, address whitelistedAddress) view returns(uint64[] whitelistedOperatorIds)
func (_Ssv *SsvCaller) GetWhitelistedOperators(opts *bind.CallOpts, operatorIds []uint64, whitelistedAddress common.Address) ([]uint64, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "getWhitelistedOperators", operatorIds, whitelistedAddress)

	if err != nil {
		return *new([]uint64), err
	}

	out0 := *abi.ConvertType(out[0], new([]uint64)).(*[]uint64)

	return out0, err

}

// GetWhitelistedOperators is a free data retrieval call binding the contract method 0xa9cf9eec.
//
// Solidity: function getWhitelistedOperators(uint64[] operatorIds, address whitelistedAddress) view returns(uint64[] whitelistedOperatorIds)
func (_Ssv *SsvSession) GetWhitelistedOperators(operatorIds []uint64, whitelistedAddress common.Address) ([]uint64, error) {
	return _Ssv.Contract.GetWhitelistedOperators(&_Ssv.CallOpts, operatorIds, whitelistedAddress)
}

// GetWhitelistedOperators is a free data retrieval call binding the contract method 0xa9cf9eec.
//
// Solidity: function getWhitelistedOperators(uint64[] operatorIds, address whitelistedAddress) view returns(uint64[] whitelistedOperatorIds)
func (_Ssv *SsvCallerSession) GetWhitelistedOperators(operatorIds []uint64, whitelistedAddress common.Address) ([]uint64, error) {
	return _Ssv.Contract.GetWhitelistedOperators(&_Ssv.CallOpts, operatorIds, whitelistedAddress)
}

// IsAddressWhitelistedInWhitelistingContract is a free data retrieval call binding the contract method 0x46e6d917.
//
// Solidity: function isAddressWhitelistedInWhitelistingContract(address addressToCheck, uint256 operatorId, address whitelistingContract) view returns(bool isWhitelisted)
func (_Ssv *SsvCaller) IsAddressWhitelistedInWhitelistingContract(opts *bind.CallOpts, addressToCheck common.Address, operatorId *big.Int, whitelistingContract common.Address) (bool, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "isAddressWhitelistedInWhitelistingContract", addressToCheck, operatorId, whitelistingContract)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAddressWhitelistedInWhitelistingContract is a free data retrieval call binding the contract method 0x46e6d917.
//
// Solidity: function isAddressWhitelistedInWhitelistingContract(address addressToCheck, uint256 operatorId, address whitelistingContract) view returns(bool isWhitelisted)
func (_Ssv *SsvSession) IsAddressWhitelistedInWhitelistingContract(addressToCheck common.Address, operatorId *big.Int, whitelistingContract common.Address) (bool, error) {
	return _Ssv.Contract.IsAddressWhitelistedInWhitelistingContract(&_Ssv.CallOpts, addressToCheck, operatorId, whitelistingContract)
}

// IsAddressWhitelistedInWhitelistingContract is a free data retrieval call binding the contract method 0x46e6d917.
//
// Solidity: function isAddressWhitelistedInWhitelistingContract(address addressToCheck, uint256 operatorId, address whitelistingContract) view returns(bool isWhitelisted)
func (_Ssv *SsvCallerSession) IsAddressWhitelistedInWhitelistingContract(addressToCheck common.Address, operatorId *big.Int, whitelistingContract common.Address) (bool, error) {
	return _Ssv.Contract.IsAddressWhitelistedInWhitelistingContract(&_Ssv.CallOpts, addressToCheck, operatorId, whitelistingContract)
}

// IsLiquidatable is a free data retrieval call binding the contract method 0x16cff008.
//
// Solidity: function isLiquidatable(address clusterOwner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) view returns(bool)
func (_Ssv *SsvCaller) IsLiquidatable(opts *bind.CallOpts, clusterOwner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (bool, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "isLiquidatable", clusterOwner, operatorIds, cluster)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsLiquidatable is a free data retrieval call binding the contract method 0x16cff008.
//
// Solidity: function isLiquidatable(address clusterOwner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) view returns(bool)
func (_Ssv *SsvSession) IsLiquidatable(clusterOwner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (bool, error) {
	return _Ssv.Contract.IsLiquidatable(&_Ssv.CallOpts, clusterOwner, operatorIds, cluster)
}

// IsLiquidatable is a free data retrieval call binding the contract method 0x16cff008.
//
// Solidity: function isLiquidatable(address clusterOwner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) view returns(bool)
func (_Ssv *SsvCallerSession) IsLiquidatable(clusterOwner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (bool, error) {
	return _Ssv.Contract.IsLiquidatable(&_Ssv.CallOpts, clusterOwner, operatorIds, cluster)
}

// IsLiquidated is a free data retrieval call binding the contract method 0xa694695b.
//
// Solidity: function isLiquidated(address clusterOwner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) view returns(bool)
func (_Ssv *SsvCaller) IsLiquidated(opts *bind.CallOpts, clusterOwner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (bool, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "isLiquidated", clusterOwner, operatorIds, cluster)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsLiquidated is a free data retrieval call binding the contract method 0xa694695b.
//
// Solidity: function isLiquidated(address clusterOwner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) view returns(bool)
func (_Ssv *SsvSession) IsLiquidated(clusterOwner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (bool, error) {
	return _Ssv.Contract.IsLiquidated(&_Ssv.CallOpts, clusterOwner, operatorIds, cluster)
}

// IsLiquidated is a free data retrieval call binding the contract method 0xa694695b.
//
// Solidity: function isLiquidated(address clusterOwner, uint64[] operatorIds, (uint32,uint64,uint64,bool,uint256) cluster) view returns(bool)
func (_Ssv *SsvCallerSession) IsLiquidated(clusterOwner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) (bool, error) {
	return _Ssv.Contract.IsLiquidated(&_Ssv.CallOpts, clusterOwner, operatorIds, cluster)
}

// IsWhitelistingContract is a free data retrieval call binding the contract method 0xbac69e6f.
//
// Solidity: function isWhitelistingContract(address contractAddress) view returns(bool)
func (_Ssv *SsvCaller) IsWhitelistingContract(opts *bind.CallOpts, contractAddress common.Address) (bool, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "isWhitelistingContract", contractAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsWhitelistingContract is a free data retrieval call binding the contract method 0xbac69e6f.
//
// Solidity: function isWhitelistingContract(address contractAddress) view returns(bool)
func (_Ssv *SsvSession) IsWhitelistingContract(contractAddress common.Address) (bool, error) {
	return _Ssv.Contract.IsWhitelistingContract(&_Ssv.CallOpts, contractAddress)
}

// IsWhitelistingContract is a free data retrieval call binding the contract method 0xbac69e6f.
//
// Solidity: function isWhitelistingContract(address contractAddress) view returns(bool)
func (_Ssv *SsvCallerSession) IsWhitelistingContract(contractAddress common.Address) (bool, error) {
	return _Ssv.Contract.IsWhitelistingContract(&_Ssv.CallOpts, contractAddress)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ssv *SsvCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ssv *SsvSession) Owner() (common.Address, error) {
	return _Ssv.Contract.Owner(&_Ssv.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ssv *SsvCallerSession) Owner() (common.Address, error) {
	return _Ssv.Contract.Owner(&_Ssv.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Ssv *SsvCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Ssv *SsvSession) PendingOwner() (common.Address, error) {
	return _Ssv.Contract.PendingOwner(&_Ssv.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Ssv *SsvCallerSession) PendingOwner() (common.Address, error) {
	return _Ssv.Contract.PendingOwner(&_Ssv.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Ssv *SsvCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Ssv *SsvSession) ProxiableUUID() ([32]byte, error) {
	return _Ssv.Contract.ProxiableUUID(&_Ssv.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Ssv *SsvCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Ssv.Contract.ProxiableUUID(&_Ssv.CallOpts)
}

// SsvNetwork is a free data retrieval call binding the contract method 0x10d04858.
//
// Solidity: function ssvNetwork() view returns(address)
func (_Ssv *SsvCaller) SsvNetwork(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ssv.contract.Call(opts, &out, "ssvNetwork")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SsvNetwork is a free data retrieval call binding the contract method 0x10d04858.
//
// Solidity: function ssvNetwork() view returns(address)
func (_Ssv *SsvSession) SsvNetwork() (common.Address, error) {
	return _Ssv.Contract.SsvNetwork(&_Ssv.CallOpts)
}

// SsvNetwork is a free data retrieval call binding the contract method 0x10d04858.
//
// Solidity: function ssvNetwork() view returns(address)
func (_Ssv *SsvCallerSession) SsvNetwork() (common.Address, error) {
	return _Ssv.Contract.SsvNetwork(&_Ssv.CallOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Ssv *SsvTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ssv.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Ssv *SsvSession) AcceptOwnership() (*types.Transaction, error) {
	return _Ssv.Contract.AcceptOwnership(&_Ssv.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Ssv *SsvTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Ssv.Contract.AcceptOwnership(&_Ssv.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address ssvNetwork_) returns()
func (_Ssv *SsvTransactor) Initialize(opts *bind.TransactOpts, ssvNetwork_ common.Address) (*types.Transaction, error) {
	return _Ssv.contract.Transact(opts, "initialize", ssvNetwork_)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address ssvNetwork_) returns()
func (_Ssv *SsvSession) Initialize(ssvNetwork_ common.Address) (*types.Transaction, error) {
	return _Ssv.Contract.Initialize(&_Ssv.TransactOpts, ssvNetwork_)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address ssvNetwork_) returns()
func (_Ssv *SsvTransactorSession) Initialize(ssvNetwork_ common.Address) (*types.Transaction, error) {
	return _Ssv.Contract.Initialize(&_Ssv.TransactOpts, ssvNetwork_)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ssv *SsvTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ssv.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ssv *SsvSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ssv.Contract.RenounceOwnership(&_Ssv.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ssv *SsvTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ssv.Contract.RenounceOwnership(&_Ssv.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ssv *SsvTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Ssv.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ssv *SsvSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ssv.Contract.TransferOwnership(&_Ssv.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ssv *SsvTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ssv.Contract.TransferOwnership(&_Ssv.TransactOpts, newOwner)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Ssv *SsvTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _Ssv.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Ssv *SsvSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _Ssv.Contract.UpgradeTo(&_Ssv.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_Ssv *SsvTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _Ssv.Contract.UpgradeTo(&_Ssv.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Ssv *SsvTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Ssv.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Ssv *SsvSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Ssv.Contract.UpgradeToAndCall(&_Ssv.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Ssv *SsvTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Ssv.Contract.UpgradeToAndCall(&_Ssv.TransactOpts, newImplementation, data)
}

// SsvAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the Ssv contract.
type SsvAdminChangedIterator struct {
	Event *SsvAdminChanged // Event containing the contract specifics and raw log

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
func (it *SsvAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SsvAdminChanged)
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
		it.Event = new(SsvAdminChanged)
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
func (it *SsvAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SsvAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SsvAdminChanged represents a AdminChanged event raised by the Ssv contract.
type SsvAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Ssv *SsvFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*SsvAdminChangedIterator, error) {

	logs, sub, err := _Ssv.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &SsvAdminChangedIterator{contract: _Ssv.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Ssv *SsvFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *SsvAdminChanged) (event.Subscription, error) {

	logs, sub, err := _Ssv.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SsvAdminChanged)
				if err := _Ssv.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_Ssv *SsvFilterer) ParseAdminChanged(log types.Log) (*SsvAdminChanged, error) {
	event := new(SsvAdminChanged)
	if err := _Ssv.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SsvBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the Ssv contract.
type SsvBeaconUpgradedIterator struct {
	Event *SsvBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *SsvBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SsvBeaconUpgraded)
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
		it.Event = new(SsvBeaconUpgraded)
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
func (it *SsvBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SsvBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SsvBeaconUpgraded represents a BeaconUpgraded event raised by the Ssv contract.
type SsvBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Ssv *SsvFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*SsvBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _Ssv.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &SsvBeaconUpgradedIterator{contract: _Ssv.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Ssv *SsvFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *SsvBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _Ssv.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SsvBeaconUpgraded)
				if err := _Ssv.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_Ssv *SsvFilterer) ParseBeaconUpgraded(log types.Log) (*SsvBeaconUpgraded, error) {
	event := new(SsvBeaconUpgraded)
	if err := _Ssv.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SsvInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Ssv contract.
type SsvInitializedIterator struct {
	Event *SsvInitialized // Event containing the contract specifics and raw log

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
func (it *SsvInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SsvInitialized)
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
		it.Event = new(SsvInitialized)
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
func (it *SsvInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SsvInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SsvInitialized represents a Initialized event raised by the Ssv contract.
type SsvInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Ssv *SsvFilterer) FilterInitialized(opts *bind.FilterOpts) (*SsvInitializedIterator, error) {

	logs, sub, err := _Ssv.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &SsvInitializedIterator{contract: _Ssv.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Ssv *SsvFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *SsvInitialized) (event.Subscription, error) {

	logs, sub, err := _Ssv.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SsvInitialized)
				if err := _Ssv.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Ssv *SsvFilterer) ParseInitialized(log types.Log) (*SsvInitialized, error) {
	event := new(SsvInitialized)
	if err := _Ssv.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SsvOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Ssv contract.
type SsvOwnershipTransferStartedIterator struct {
	Event *SsvOwnershipTransferStarted // Event containing the contract specifics and raw log

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
func (it *SsvOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SsvOwnershipTransferStarted)
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
		it.Event = new(SsvOwnershipTransferStarted)
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
func (it *SsvOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SsvOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SsvOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Ssv contract.
type SsvOwnershipTransferStarted struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Ssv *SsvFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SsvOwnershipTransferStartedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ssv.contract.FilterLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SsvOwnershipTransferStartedIterator{contract: _Ssv.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Ssv *SsvFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *SsvOwnershipTransferStarted, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ssv.contract.WatchLogs(opts, "OwnershipTransferStarted", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SsvOwnershipTransferStarted)
				if err := _Ssv.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
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

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed previousOwner, address indexed newOwner)
func (_Ssv *SsvFilterer) ParseOwnershipTransferStarted(log types.Log) (*SsvOwnershipTransferStarted, error) {
	event := new(SsvOwnershipTransferStarted)
	if err := _Ssv.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SsvOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Ssv contract.
type SsvOwnershipTransferredIterator struct {
	Event *SsvOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SsvOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SsvOwnershipTransferred)
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
		it.Event = new(SsvOwnershipTransferred)
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
func (it *SsvOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SsvOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SsvOwnershipTransferred represents a OwnershipTransferred event raised by the Ssv contract.
type SsvOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ssv *SsvFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SsvOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ssv.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SsvOwnershipTransferredIterator{contract: _Ssv.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ssv *SsvFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SsvOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ssv.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SsvOwnershipTransferred)
				if err := _Ssv.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Ssv *SsvFilterer) ParseOwnershipTransferred(log types.Log) (*SsvOwnershipTransferred, error) {
	event := new(SsvOwnershipTransferred)
	if err := _Ssv.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SsvUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Ssv contract.
type SsvUpgradedIterator struct {
	Event *SsvUpgraded // Event containing the contract specifics and raw log

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
func (it *SsvUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SsvUpgraded)
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
		it.Event = new(SsvUpgraded)
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
func (it *SsvUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SsvUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SsvUpgraded represents a Upgraded event raised by the Ssv contract.
type SsvUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Ssv *SsvFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*SsvUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Ssv.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &SsvUpgradedIterator{contract: _Ssv.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Ssv *SsvFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *SsvUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Ssv.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SsvUpgraded)
				if err := _Ssv.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Ssv *SsvFilterer) ParseUpgraded(log types.Log) (*SsvUpgraded, error) {
	event := new(SsvUpgraded)
	if err := _Ssv.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
