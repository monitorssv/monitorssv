package ssv

import (
	"errors"
	"github.com/monitorssv/monitorssv/eth1/utils"
	"math/big"
)

var (
	noValidatorErr       = errors.New("no validators")
	alreadyLiquidatedErr = errors.New("already liquidated")
	canLiquidatedErr     = errors.New("can liquidate")
)

func (s *SSV) CalcLiquidation(cluster Cluster) (uint64, uint64, uint64, string, error) {
	if cluster.ClusterInfo.ValidatorCount == 0 {
		return 0, 0, 0, "", noValidatorErr
	}

	if !cluster.ClusterInfo.Active {
		return 0, 0, 0, "", alreadyLiquidatedErr
	}

	curBlock, err := s.client.BlockNumber()
	if err != nil {
		return 0, 0, 0, "", err
	}

	liquidationInfo, err := s.GetSSVLiquidationInfo(cluster)
	if err != nil {
		return 0, 0, 0, "", err
	}

	if liquidationInfo.ClusterBalance.Cmp(liquidationInfo.MinimumLiquidationCollateral) <= 0 {
		return 0, 0, 0, "", canLiquidatedErr
	}

	minimumBlocksBeforeLiquidation := int64(liquidationInfo.LiquidationThresholdPeriod)

	burnRate := uint64(0)
	for _, opFee := range liquidationInfo.OperatorsFee {
		burnRate += opFee
	}

	fee := liquidationInfo.NetworkFee + burnRate

	perLiquidationThreshold := big.NewInt(0).Mul(big.NewInt(minimumBlocksBeforeLiquidation), big.NewInt(int64(fee)))
	liquidationThreshold := big.NewInt(0).Mul(perLiquidationThreshold, big.NewInt(int64(cluster.ClusterInfo.ValidatorCount)))

	if liquidationInfo.ClusterBalance.Cmp(liquidationThreshold) > 0 {
		// The number of validators and minimum collateral will affect the liquidation runway
		reserve := liquidationInfo.MinimumLiquidationCollateral
		if liquidationThreshold.Cmp(reserve) > 0 {
			reserve = liquidationThreshold
		}

		activeBalance := big.NewInt(0).Sub(liquidationInfo.ClusterBalance, reserve)

		preValidatorBalance := big.NewInt(0).Div(activeBalance, big.NewInt(int64(cluster.ClusterInfo.ValidatorCount)))
		if preValidatorBalance.Uint64() == 0 {
			return 0, 0, 0, "", canLiquidatedErr
		}

		activeBlock := big.NewInt(0).Div(preValidatorBalance, big.NewInt(0).SetUint64(fee)).Uint64()
		if activeBlock == 0 {
			return 0, 0, 0, "", canLiquidatedErr
		}

		return curBlock + activeBlock, curBlock, fee, liquidationInfo.ClusterBalance.String(), nil
	}

	return 0, 0, 0, "", canLiquidatedErr
}

type LiquidationInfo struct {
	NetworkFee                   uint64
	LiquidationThresholdPeriod   uint64
	MinimumLiquidationCollateral *big.Int
	ClusterBalance               *big.Int
	OperatorsFee                 []uint64
}

func (s *SSV) GetSSVLiquidationInfo(cluster Cluster) (*LiquidationInfo, error) {
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

	getBalanceData, err := ssvViewABI.Methods[getBalance].Inputs.Pack(cluster.Owner, cluster.OperatorIds, cluster.ClusterInfo)
	if err != nil {
		return nil, err
	}
	callStructs = append(callStructs, utils.Struct0{
		Target:   s.ssvNetworkViewAdd,
		CallData: append(ssvViewABI.Methods[getBalance].ID, getBalanceData...),
	})

	for _, opId := range cluster.OperatorIds {
		data, err := ssvViewABI.Methods[getOperatorFee].Inputs.Pack(opId)
		if err != nil {
			return nil, err
		}
		callStructs = append(callStructs, utils.Struct0{
			Target:   s.ssvNetworkViewAdd,
			CallData: append(ssvViewABI.Methods[getOperatorFee].ID, data...),
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

	var info LiquidationInfo
	info.NetworkFee = results[0].Uint64()
	info.LiquidationThresholdPeriod = results[1].Uint64()
	info.MinimumLiquidationCollateral = results[2]
	info.ClusterBalance = results[3]
	operatorsFee := make([]uint64, 0)
	for _, r := range results[4:] {
		operatorsFee = append(operatorsFee, r.Uint64())
	}
	info.OperatorsFee = operatorsFee

	return &info, nil
}
