package ssv

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/monitorssv/monitorssv/eth1/utils"
	"math/big"
)

type NetworkInfo struct {
	NetworkFee                   string
	LiquidationThresholdPeriod   int64
	MinimumLiquidationCollateral string
	OperatorValidatorLimit       int64
}

func (s *SSV) GetNetworkInfo() (*NetworkInfo, error) {
	var callStructs = make([]utils.Struct0, 0)
	callStructs = append(callStructs, utils.Struct0{
		Target:   s.ssvNetworkViewAdd,
		CallData: ssvViewABI.Methods[getNetworkFee].ID,
	})
	callStructs = append(callStructs, utils.Struct0{
		Target:   s.ssvNetworkViewAdd,
		CallData: ssvViewABI.Methods[getLiquidationThresholdPeriod].ID,
	})
	callStructs = append(callStructs, utils.Struct0{
		Target:   s.ssvNetworkViewAdd,
		CallData: ssvViewABI.Methods[getMinimumLiquidationCollateral].ID,
	})
	callStructs = append(callStructs, utils.Struct0{
		Target:   s.ssvNetworkViewAdd,
		CallData: ssvViewABI.Methods[getValidatorsPerOperatorLimit].ID,
	})

	multiCall, err := utils.NewMulticall(s.cfg.Network, s.client.GetClient())
	if err != nil {
		return nil, err
	}

	outs, err := multiCall.MulticallCaller.Aggregate(nil, callStructs)
	if err != nil {
		return nil, err
	}
	results := make([]*big.Int, 0)
	for _, r := range outs[1].([][]uint8) {
		results = append(results, big.NewInt(0).SetBytes(r))
	}

	var info NetworkInfo
	info.LiquidationThresholdPeriod = results[1].Int64()
	info.OperatorValidatorLimit = results[3].Int64()
	fee := big.NewInt(0).Mul(results[0], big.NewInt(2613400))
	networkFee := utils.ToSSV(fee, "%.2f")
	minimumLiquidationCollateral := utils.ToSSV(results[2], "%.9f")
	info.NetworkFee = networkFee
	info.MinimumLiquidationCollateral = minimumLiquidationCollateral
	return &info, nil
}

type OperatorDeclaredFee struct {
	IsDeclared        bool
	Fee               string
	ApprovalBeginTime uint64
	ApprovalEndTime   uint64
}

func (s *SSV) GetOperatorDeclaredFee(operatorIds []uint64) ([]OperatorDeclaredFee, error) {
	var callStructs = make([]utils.Struct0, 0)
	for _, opId := range operatorIds {
		data, err := ssvViewABI.Methods[getOperatorDeclaredFee].Inputs.Pack(opId)
		if err != nil {
			return nil, err
		}
		callStructs = append(callStructs, utils.Struct0{
			Target:   s.ssvNetworkViewAdd,
			CallData: append(ssvViewABI.Methods[getOperatorDeclaredFee].ID, data...),
		})
	}

	multiCall, err := utils.NewMulticall(s.cfg.Network, s.client.GetClient())
	if err != nil {
		return nil, err
	}

	outs, err := multiCall.MulticallCaller.Aggregate(nil, callStructs)
	if err != nil {
		return nil, err
	}

	var declaredFees = make([]OperatorDeclaredFee, 0)
	for _, out := range outs[1].([][]uint8) {
		data, err := ssvViewABI.Methods[getOperatorDeclaredFee].Outputs.Unpack(out)
		if err != nil {
			return nil, err
		}

		isDeclared := data[0].(bool)
		fee := data[1].(*big.Int).String()
		beginTime := data[2].(uint64)
		endTime := data[3].(uint64)

		declaredFee := OperatorDeclaredFee{
			IsDeclared:        isDeclared,
			Fee:               fee,
			ApprovalBeginTime: beginTime,
			ApprovalEndTime:   endTime,
		}
		declaredFees = append(declaredFees, declaredFee)
	}

	return declaredFees, nil
}

func (s *SSV) GetOperatorEarnings(operatorIds []uint64) ([]string, error) {
	var callStructs = make([]utils.Struct0, 0)
	for _, opId := range operatorIds {
		data, err := ssvViewABI.Methods[getOperatorEarnings].Inputs.Pack(opId)
		if err != nil {
			return nil, err
		}
		callStructs = append(callStructs, utils.Struct0{
			Target:   s.ssvNetworkViewAdd,
			CallData: append(ssvViewABI.Methods[getOperatorEarnings].ID, data...),
		})
	}

	multiCall, err := utils.NewMulticall(s.cfg.Network, s.client.GetClient())
	if err != nil {
		return nil, err
	}

	outs, err := multiCall.MulticallCaller.Aggregate(nil, callStructs)
	if err != nil {
		return nil, err
	}
	results := make([]*big.Int, 0)
	for _, r := range outs[1].([][]uint8) {
		results = append(results, big.NewInt(0).SetBytes(r))
	}

	var operatorEarnings = make([]string, 0)
	for _, result := range results {
		if result.Cmp(big.NewInt(0)) == 0 {
			operatorEarnings = append(operatorEarnings, "0")
			continue
		}
		earnings := utils.ToSSV(result, "%.9f")
		operatorEarnings = append(operatorEarnings, earnings)
	}

	return operatorEarnings, nil
}

func GetSSVRewardCumulativeClaimed(client *ethclient.Client, accounts []common.Address) ([]*big.Int, error) {
	var callStructs = make([]utils.Struct0, 0)
	for _, account := range accounts {
		data, err := ssvRewardABI.Methods[cumulativeClaimedFunc].Inputs.Pack(account)
		if err != nil {
			return nil, err
		}
		callStructs = append(callStructs, utils.Struct0{
			Target:   ssvRewardContractAddr,
			CallData: append(ssvRewardABI.Methods[cumulativeClaimedFunc].ID, data...),
		})
	}

	multiCall, err := utils.NewMulticall("mainnet", client)
	if err != nil {
		return nil, err
	}

	outs, err := multiCall.MulticallCaller.Aggregate(nil, callStructs)
	if err != nil {
		return nil, err
	}

	results := make([]*big.Int, 0)
	for _, r := range outs[1].([][]uint8) {
		results = append(results, big.NewInt(0).SetBytes(r))
	}

	return results, nil
}

func GetSSVRewardMerkleRootOnChain(client *ethclient.Client) (string, error) {
	var callStructs = make([]utils.Struct0, 0)
	callStructs = append(callStructs, utils.Struct0{
		Target:   ssvRewardContractAddr,
		CallData: ssvRewardABI.Methods[getMerkleRootFunc].ID,
	})

	multiCall, err := utils.NewMulticall("mainnet", client)
	if err != nil {
		return "", err
	}

	outs, err := multiCall.MulticallCaller.Aggregate(nil, callStructs)
	if err != nil {
		return "", err
	}
	results := make([]*big.Int, 0)
	for _, r := range outs[1].([][]uint8) {
		results = append(results, big.NewInt(0).SetBytes(r))
	}

	var root = fmt.Sprintf("0x%s", hex.EncodeToString(results[0].Bytes()))
	return root, nil
}
