package ssv

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/monitorssv/monitorssv/alert"
	"github.com/monitorssv/monitorssv/store"
	"math/big"
	"sort"
	"strings"
	"time"
)

func (s *SSV) ScanSSVEvent(startBlock, endBlock uint64) (uint64, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fetchLogs, fetchError := s.fetchEvents(ctx, startBlock, endBlock)

	lastProcessedBlock, err := s.handleEvents(fetchLogs)
	if err != nil {
		return lastProcessedBlock, err
	}

	if err := <-fetchError; err != nil {
		return lastProcessedBlock, err
	}

	return endBlock, nil
}

type BlockLogs struct {
	BlockNumber uint64
	Logs        []ethtypes.Log
}

const (
	logsBatchSize = 5000
	defaultLogBuf = 8192
)

func (s *SSV) fetchEvents(ctx context.Context, startBlock, endBlock uint64) (<-chan BlockLogs, <-chan error) {
	ssvLog.Infow("fetchEvent", "startBlock", startBlock, "endBlock", endBlock)

	fetchLogs := make(chan BlockLogs, defaultLogBuf)
	fetchError := make(chan error, 1)

	go func() {
		defer close(fetchLogs)
		defer close(fetchError)

		if startBlock > endBlock {
			fetchError <- fmt.Errorf("bad input")
			return
		}

		for fromBlock := startBlock; fromBlock <= endBlock; fromBlock += logsBatchSize {
			toBlock := fromBlock + logsBatchSize - 1
			if toBlock > endBlock {
				toBlock = endBlock
			}

			addresses := []common.Address{s.ssvNetworkAddr}
			if s.cfg.Network == "mainnet" { // only mainnet
				addresses = append(addresses, ssvRewardContractAddr)
			}

			start := time.Now()
			results, err := s.client.FilterLogs(ethereum.FilterQuery{
				Addresses: addresses,
				FromBlock: new(big.Int).SetUint64(fromBlock),
				ToBlock:   new(big.Int).SetUint64(toBlock),
			})
			if err != nil {
				fetchError <- err
				return
			}

			ssvLog.Infow("fetch events",
				"fromBlock", fromBlock,
				"toBlock", toBlock,
				"endBlock", endBlock,
				"progress", fmt.Sprintf("%.2f%%", float64(toBlock-startBlock+1)/float64(endBlock-startBlock+1)*100),
				"events", len(results),
				"took", time.Since(start),
			)

			select {
			case <-ctx.Done():
				fetchError <- ctx.Err()
				return

			case <-s.close:
				fetchError <- fmt.Errorf("closed")
				return

			default:
				validLogs := make([]ethtypes.Log, 0, len(results))
				for _, log := range results {
					if log.Removed {
						ssvLog.Warn("log is removed",
							"block_hash", log.BlockHash.Hex(),
							"tx_hash", log.TxHash.Hex(),
							"log_index", log.Index)
						continue
					}
					validLogs = append(validLogs, log)
				}
				if len(validLogs) == 0 {
					// Emit empty block logs to indicate that we have advanced to this block.
					fetchLogs <- BlockLogs{BlockNumber: toBlock}
				} else {
					for _, blockLogs := range packLogs(validLogs) {
						fetchLogs <- blockLogs
					}
				}
			}
		}
	}()

	return fetchLogs, fetchError
}

func packLogs(logs []ethtypes.Log) []BlockLogs {
	// Sort the logs by block number.
	sort.Slice(logs, func(i, j int) bool {
		if logs[i].BlockNumber == logs[j].BlockNumber {
			return logs[i].TxIndex < logs[j].TxIndex
		}
		return logs[i].BlockNumber < logs[j].BlockNumber
	})

	var all []BlockLogs
	for _, log := range logs {
		// Create a BlockLogs if there isn't one for this block yet.
		if len(all) == 0 || all[len(all)-1].BlockNumber != log.BlockNumber {
			all = append(all, BlockLogs{
				BlockNumber: log.BlockNumber,
			})
		}

		// Append the log to the current BlockLogs.
		all[len(all)-1].Logs = append(all[len(all)-1].Logs, log)
	}

	return all
}

func (s *SSV) handleEvents(fetchLogs <-chan BlockLogs) (uint64, error) {
	var lastProcessedBlock uint64
	for event := range fetchLogs {
		lastProcessedBlock = event.BlockNumber
		err := s.processBlockEvents(event.Logs)
		if err != nil {
			return lastProcessedBlock - 1, err
		}
	}

	return lastProcessedBlock, nil
}

func (s *SSV) processBlockEvents(logs []ethtypes.Log) error {
	for _, vLog := range logs {
		// SSV Network’s Incentive Program contract event
		if vLog.Address == ssvRewardContractAddr {
			if vLog.Topics[0] == ssvRewardClaimedTopic {
				event := ssvRewardABI.Events[ssvRewardClaimEvent]
				data, err := event.Inputs.Unpack(vLog.Data)
				if err != nil {
					ssvLog.Errorw("processBlockEvents: SSVReward Unpack", "err", err)
					continue
				}
				account := data[0].(common.Address)
				amount := data[1].(*big.Int)
				err = s.store.UpdateClaimed(account.String(), amount)
				if err != nil {
					ssvLog.Errorw("processBlockEvents: SSVReward UpdateClaimed", "err", err)
				}
			}
			continue
		}

		// ssv network event
		event, ok := s.events[vLog.Topics[0]]
		if !ok {
			ssvLog.Warnw("unknown event topic", "topic", vLog.Topics[0].Hex(), "txHash", vLog.TxHash.Hex())
			continue
		}

		ssvLog.Infow("processing block event", "event", event.Name, "block", vLog.BlockNumber, "txHash", vLog.TxHash.Hex())

		switch event.Name {
		case Initialized, AdminChanged, Upgraded, ModuleUpgraded, BeaconUpgraded, OwnershipTransferred, OwnershipTransferStarted:
			if err := s.recordEvent(vLog, vLog.Address.String(), event.Name, ""); err != nil {
				return err
			}
		case OperatorFeeIncreaseLimitUpdated, DeclareOperatorFeePeriodUpdated, ExecuteOperatorFeePeriodUpdated,
			LiquidationThresholdPeriodUpdated, MinimumLiquidationCollateralUpdated,
			NetworkEarningsWithdrawn, OperatorMaximumFeeUpdated:
			if err := s.recordEvent(vLog, vLog.Address.String(), event.Name, ""); err != nil {
				return err
			}
		case NetworkFeeUpdated:
			if s.isSynced.Load() {
				ssvLog.Info("NetworkFeeUpdated: calculate the liquidation block: allCluster")
				s.calcAllClusterLiquidationChan <- 0 // calc all cluster liquidation block

				// notify cluster owner
				data, err := event.Inputs.Unpack(vLog.Data)
				if err != nil {
					return err
				}

				oldNetworkFee := data[0].(*big.Int)
				newNetworkFee := data[1].(*big.Int)

				s.networkFeeChangeAlarmChan <- alert.NetworkFeeChangeNotify{
					Block:         vLog.BlockNumber,
					OldNetworkFee: oldNetworkFee,
					NewNetworkFee: newNetworkFee,
				}
			}
			if err := s.recordEvent(vLog, vLog.Address.String(), event.Name, ""); err != nil {
				return err
			}
		case OperatorFeeDeclared, OperatorFeeDeclarationCancelled, OperatorWithdrawn:
			var owner common.Address
			copy(owner[:], vLog.Topics[1][12:])
			if err := s.recordEvent(vLog, owner.String(), event.Name, ""); err != nil {
				return err
			}
		case OperatorAdded:
			var operatorId = big.NewInt(0).SetBytes(vLog.Topics[1][:]).Uint64()
			var owner common.Address
			copy(owner[:], vLog.Topics[2][12:])
			data, err := event.Inputs.Unpack(vLog.Data)
			if err != nil {
				return err
			}

			pubKeyBytes := data[0].([]byte)
			pubKeyUnpack, err := pubKeyInputs.Unpack(pubKeyBytes)
			pubKey := ""
			if err != nil {
				ssvLog.Warnw("processBlockEvents: SSVOperator PubKey Unpack", "operatorId", operatorId, "err", err)
			} else {
				pubKey = pubKeyUnpack[0].(string)
				_, err = base64.StdEncoding.DecodeString(pubKey)
				if err != nil {
					ssvLog.Warnw("invalid pubKey", "operatorId", operatorId, "pubKey", pubKey)
					pubKey = hex.EncodeToString(pubKeyBytes)
				}
			}

			operatorFee := data[1].(*big.Int)
			if err = s.store.CreateOperator(&store.OperatorInfo{
				Owner:              owner.String(),
				OperatorId:         operatorId,
				PubKey:             pubKey,
				OperatorName:       "",
				ValidatorCount:     0,
				OperatorFee:        operatorFee.String(),
				WhitelistedAddress: "",
				PrivacyStatus:      false,
				RegistrationBlock:  int64(vLog.BlockNumber),
				RemoveBlock:        0,
			}); err != nil {
				return err
			}
			if err := s.recordEvent(vLog, owner.String(), event.Name, ""); err != nil {
				return err
			}
		case OperatorRemoved:
			var operatorId = big.NewInt(0).SetBytes(vLog.Topics[1][:]).Uint64()
			err := s.store.RemoveOperator(operatorId, int64(vLog.BlockNumber))
			if err != nil {
				return err
			}
			operatorInfo, err := s.store.GetOperatorByOperatorId(operatorId)
			if err != nil {
				return err
			}
			if err := s.recordEvent(vLog, operatorInfo.Owner, event.Name, ""); err != nil {
				return err
			}
		case OperatorPrivacyStatusUpdated:
			data, err := event.Inputs.Unpack(vLog.Data)
			if err != nil {
				return err
			}
			operatorIds := data[0].([]uint64)
			toPrivate := data[1].(bool)
			for _, operatorId := range operatorIds {
				err = s.store.UpdateOperatorPrivacyStatus(operatorId, toPrivate)
				if err != nil {
					return err
				}
			}

			operatorInfo, err := s.store.GetOperatorByOperatorId(operatorIds[0])
			if err != nil {
				return err
			}
			if err := s.recordEvent(vLog, operatorInfo.Owner, event.Name, ""); err != nil {
				return err
			}
		case OperatorFeeExecuted:
			var owner common.Address
			copy(owner[:], vLog.Topics[1][12:])
			var operatorId = big.NewInt(0).SetBytes(vLog.Topics[2][:]).Uint64()

			data, err := event.Inputs.Unpack(vLog.Data)
			if err != nil {
				return err
			}
			operatorFee := data[1].(*big.Int)
			if err = s.store.UpdateOperatorFee(operatorId, operatorFee.String()); err != nil {
				return err
			}

			if err := s.recordEvent(vLog, owner.String(), event.Name, ""); err != nil {
				return err
			}
			if s.isSynced.Load() {
				ssvLog.Infow("OperatorFeeExecuted: calculate the liquidation block", "operatorId", operatorId)
				// Calculate the liquidation block of the cluster associated with this operator
				s.calcAllClusterLiquidationChan <- operatorId
				s.operatorFeeChangeAlarmChan <- alert.OperatorFeeChangeNotify{
					Block:       vLog.BlockNumber,
					OperatorId:  operatorId,
					OperatorFee: operatorFee,
				}
			}
		case OperatorMultipleWhitelistUpdated:
			data, err := event.Inputs.Unpack(vLog.Data)
			if err != nil {
				return err
			}
			operatorIds := data[0].([]uint64)
			whitelistAddresses := data[1].([]common.Address)
			whitelistUpdatedAddress := ""
			for j, whitelistAddr := range whitelistAddresses {
				if j == 0 {
					whitelistUpdatedAddress = whitelistAddr.String()
					continue
				}

				whitelistUpdatedAddress = fmt.Sprintf("%s,%s", whitelistUpdatedAddress, whitelistAddr.String())
			}

			owner := ""
			for _, operatorId := range operatorIds {
				operatorInfo, err := s.store.GetOperatorByOperatorId(operatorId)
				if err != nil {
					return err
				}

				if operatorInfo.WhitelistedAddress != "" {
					curWhitelistAddresses := strings.Split(operatorInfo.WhitelistedAddress, ",")
					for _, whitelistedAddr := range curWhitelistAddresses {
						if !strings.Contains(whitelistUpdatedAddress, whitelistedAddr) {
							whitelistUpdatedAddress = fmt.Sprintf("%s,%s", whitelistUpdatedAddress, whitelistedAddr)
						}
					}
				}

				if err = s.store.UpdateWhitelistedAddress(operatorId, whitelistUpdatedAddress); err != nil {
					return err
				}
				if owner == "" {
					owner = operatorInfo.Owner
				}
			}
			if err := s.recordEvent(vLog, owner, event.Name, ""); err != nil {
				return err
			}
		case OperatorMultipleWhitelistRemoved:
			data, err := event.Inputs.Unpack(vLog.Data)
			if err != nil {
				return err
			}
			operatorIds := data[0].([]uint64)
			whitelistAddresses := data[1].([]common.Address)
			owner := ""
			for _, operatorId := range operatorIds {
				operatorInfo, err := s.store.GetOperatorByOperatorId(operatorId)
				if err != nil {
					return err
				}

				for _, whitelistAddr := range whitelistAddresses {
					if !strings.Contains(operatorInfo.WhitelistedAddress, whitelistAddr.String()) {
						ssvLog.Warnw("Event: OperatorMultipleWhitelistRemoved", "whitelistedAddress", operatorInfo.WhitelistedAddress, "whitelistAddr", whitelistAddr)
					}
					operatorInfo.WhitelistedAddress = strings.Replace(operatorInfo.WhitelistedAddress, whitelistAddr.String(), "", -1)
					operatorInfo.WhitelistedAddress = strings.Replace(operatorInfo.WhitelistedAddress, ",,", ",", -1)
					operatorInfo.WhitelistedAddress = strings.Trim(operatorInfo.WhitelistedAddress, ",")
				}

				if err = s.store.UpdateWhitelistedAddress(operatorId, operatorInfo.WhitelistedAddress); err != nil {
					return err
				}
				if owner == "" {
					owner = operatorInfo.Owner
				}
			}
			if err := s.recordEvent(vLog, owner, event.Name, ""); err != nil {
				return err
			}
		case OperatorWhitelistingContractUpdated:
			data, err := event.Inputs.Unpack(vLog.Data)
			if err != nil {
				return err
			}
			operatorIds := data[0].([]uint64)
			whitelistingContract := data[1].(common.Address)
			owner := ""
			for _, operatorId := range operatorIds {
				operatorInfo, err := s.store.GetOperatorByOperatorId(operatorId)
				if err != nil {
					return err
				}
				if err = s.store.UpdateWhitelistingContract(operatorId, whitelistingContract.String()); err != nil {
					return err
				}
				if owner == "" {
					owner = operatorInfo.Owner
				}
			}
			if err := s.recordEvent(vLog, owner, event.Name, ""); err != nil {
				return err
			}
		case ValidatorAdded:
			var owner common.Address
			copy(owner[:], vLog.Topics[1][12:])

			data, err := event.Inputs.Unpack(vLog.Data)
			if err != nil {
				return err
			}
			operatorIds := data[0].([]uint64)
			pubKey := data[1].([]byte)
			clusterBytes, err := json.Marshal(data[3])
			if err != nil {
				return err
			}
			cluster := ISSVNetworkCoreCluster{}
			err = json.Unmarshal(clusterBytes, &cluster)
			if err != nil {
				return err
			}
			clusterId := CalcClusterId(owner, operatorIds)

			if err = s.store.CreateOrUpdateCluster(&store.ClusterInfo{
				ClusterID:        clusterId,
				Owner:            owner.String(),
				EoaOwner:         "0x",
				OperatorIds:      toStrOperatorIds(operatorIds),
				ValidatorCount:   cluster.ValidatorCount,
				NetworkFeeIndex:  cluster.NetworkFeeIndex,
				Index:            cluster.Index,
				Active:           cluster.Active,
				Balance:          cluster.Balance.String(),
				LiquidationBlock: 0, // not calculated
			}); err != nil {
				return err
			}

			if err = s.store.CreateValidator(&store.ValidatorInfo{
				ClusterID:         clusterId,
				Owner:             owner.String(),
				OperatorIds:       toStrOperatorIds(operatorIds),
				PublicKey:         hex.EncodeToString(pubKey),
				ValidatorIndex:    store.DefaultValidatorIndex,
				RegistrationBlock: int64(vLog.BlockNumber),
				RemoveBlock:       0,
				ExitedBlock:       0,
				Status:            "Active",
			}); err != nil {
				return err
			}

			if err = s.store.BatchUpdateOperatorValidatorCount(operatorIds, true); err != nil {
				return err
			}

			for _, operatorId := range operatorIds {
				if err = s.store.UpdateOperatorClusterIds(operatorId, clusterId, true); err != nil {
					return err
				}
			}

			if err = s.recordEvent(vLog, owner.String(), event.Name, clusterId); err != nil {
				return err
			}

			s.calcLiquidation(clusterId, owner, operatorIds, cluster)
		case ValidatorRemoved:
			var owner common.Address
			copy(owner[:], vLog.Topics[1][12:])

			data, err := event.Inputs.Unpack(vLog.Data)
			if err != nil {
				return err
			}
			operatorIds := data[0].([]uint64)
			pubKey := data[1].([]byte)
			clusterBytes, err := json.Marshal(data[2])
			if err != nil {
				return err
			}
			cluster := ISSVNetworkCoreCluster{}
			err = json.Unmarshal(clusterBytes, &cluster)
			if err != nil {
				return err
			}
			clusterId := CalcClusterId(owner, operatorIds)

			if err = s.store.CreateOrUpdateCluster(&store.ClusterInfo{
				ClusterID:       clusterId,
				Owner:           owner.String(),
				EoaOwner:        "0x",
				OperatorIds:     toStrOperatorIds(operatorIds),
				ValidatorCount:  cluster.ValidatorCount,
				NetworkFeeIndex: cluster.NetworkFeeIndex,
				Index:           cluster.Index,
				Active:          cluster.Active,
				Balance:         cluster.Balance.String(),
			}); err != nil {
				return err
			}

			if err = s.store.RemoveValidator(hex.EncodeToString(pubKey), clusterId, int64(vLog.BlockNumber)); err != nil {
				return err
			}

			if cluster.Active {
				if err = s.store.BatchUpdateOperatorValidatorCount(operatorIds, false); err != nil {
					return err
				}
			}

			if cluster.ValidatorCount == 0 {
				for _, operatorId := range operatorIds {
					if err = s.store.UpdateOperatorClusterIds(operatorId, clusterId, false); err != nil {
						return err
					}
				}
			}

			if err = s.recordEvent(vLog, owner.String(), event.Name, clusterId); err != nil {
				return err
			}
			s.calcLiquidation(clusterId, owner, operatorIds, cluster)
		case ClusterLiquidated, ClusterReactivated:
			var owner common.Address
			copy(owner[:], vLog.Topics[1][12:])

			data, err := event.Inputs.Unpack(vLog.Data)
			if err != nil {
				return err
			}
			operatorIds := data[0].([]uint64)
			clusterBytes, err := json.Marshal(data[1])
			if err != nil {
				return err
			}
			cluster := ISSVNetworkCoreCluster{}
			err = json.Unmarshal(clusterBytes, &cluster)
			if err != nil {
				return err
			}
			clusterId := CalcClusterId(owner, operatorIds)

			if err = s.store.CreateOrUpdateCluster(&store.ClusterInfo{
				ClusterID:       clusterId,
				Owner:           owner.String(),
				EoaOwner:        "0x",
				OperatorIds:     toStrOperatorIds(operatorIds),
				ValidatorCount:  cluster.ValidatorCount,
				NetworkFeeIndex: cluster.NetworkFeeIndex,
				Index:           cluster.Index,
				Active:          cluster.Active,
				Balance:         cluster.Balance.String(),
			}); err != nil {
				return err
			}

			if event.Name == ClusterLiquidated {
				if err = s.store.ClusterLiquidation(clusterId, vLog.BlockNumber); err != nil {
					return err
				}

				if err = s.store.BatchUpdateOperatorValidatorCounts(operatorIds, cluster.ValidatorCount, false); err != nil {
					return err
				}

				for _, operatorId := range operatorIds {
					if err = s.store.UpdateOperatorClusterIds(operatorId, clusterId, false); err != nil {
						return err
					}
				}
			} else {
				// ClusterReactivated
				if err = s.store.BatchUpdateOperatorValidatorCounts(operatorIds, cluster.ValidatorCount, true); err != nil {
					return err
				}

				for _, operatorId := range operatorIds {
					if err = s.store.UpdateOperatorClusterIds(operatorId, clusterId, true); err != nil {
						return err
					}
				}
			}

			s.calcLiquidation(clusterId, owner, operatorIds, cluster)

			if err = s.recordEvent(vLog, owner.String(), event.Name, clusterId); err != nil {
				return err
			}
		case ClusterWithdrawn, ClusterDeposited:
			var owner common.Address
			copy(owner[:], vLog.Topics[1][12:])

			data, err := event.Inputs.Unpack(vLog.Data)
			if err != nil {
				return err
			}
			operatorIds := data[0].([]uint64)
			clusterBytes, err := json.Marshal(data[2])
			if err != nil {
				return err
			}
			cluster := ISSVNetworkCoreCluster{}
			err = json.Unmarshal(clusterBytes, &cluster)
			if err != nil {
				return err
			}
			clusterId := CalcClusterId(owner, operatorIds)

			if err = s.store.CreateOrUpdateCluster(&store.ClusterInfo{
				ClusterID:       clusterId,
				Owner:           owner.String(),
				EoaOwner:        "0x",
				OperatorIds:     toStrOperatorIds(operatorIds),
				ValidatorCount:  cluster.ValidatorCount,
				NetworkFeeIndex: cluster.NetworkFeeIndex,
				Index:           cluster.Index,
				Active:          cluster.Active,
				Balance:         cluster.Balance.String(),
			}); err != nil {
				return err
			}

			if err = s.recordEvent(vLog, owner.String(), event.Name, clusterId); err != nil {
				return err
			}

			s.calcLiquidation(clusterId, owner, operatorIds, cluster)
		case ValidatorExited:
			var owner common.Address
			copy(owner[:], vLog.Topics[1][12:])

			data, err := event.Inputs.Unpack(vLog.Data)
			if err != nil {
				return err
			}
			operatorIds := data[0].([]uint64)
			pubKey := data[1].([]byte)
			clusterId := CalcClusterId(owner, operatorIds)

			if err = s.store.ExitValidator(hex.EncodeToString(pubKey), int64(vLog.BlockNumber)); err != nil {
				return err
			}

			if err = s.recordEvent(vLog, owner.String(), event.Name, clusterId); err != nil {
				return err
			}
		case FeeRecipientAddressUpdated:
			var owner common.Address
			copy(owner[:], vLog.Topics[1][12:])

			data, err := event.Inputs.Unpack(vLog.Data)
			if err != nil {
				return err
			}
			recipientAddress := data[0].(common.Address)
			if err = s.store.CreateOrUpdateClusterFeeAddress(&store.FeeAddressInfo{
				Owner:      owner.String(),
				FeeAddress: recipientAddress.String(),
			}); err != nil {
				return err
			}

			if err = s.recordEvent(vLog, owner.String(), event.Name, ""); err != nil {
				return err
			}
		case OperatorWhitelistUpdated:
			var operatorId = big.NewInt(0).SetBytes(vLog.Topics[1][:]).Uint64()
			data, err := event.Inputs.Unpack(vLog.Data)
			if err != nil {
				return err
			}

			operatorInfo, err := s.store.GetOperatorByOperatorId(operatorId)
			if err != nil {
				return err
			}

			whitelistAddress := data[0].(common.Address)
			if whitelistAddress.String() == "0x0000000000000000000000000000000000000000" {
				operatorInfo.WhitelistedAddress = ""
			} else {
				operatorInfo.WhitelistedAddress = whitelistAddress.String()
			}

			if err = s.store.UpdateWhitelistedAddress(operatorId, operatorInfo.WhitelistedAddress); err != nil {
				return err
			}

			if err := s.recordEvent(vLog, operatorInfo.Owner, event.Name, ""); err != nil {
				return err
			}
		default:
			ssvLog.Warnw("unknown event:", "name", event.Name, "txHash", vLog.TxHash.Hex())
		}
	}
	return nil
}

func (s *SSV) recordEvent(vLog ethtypes.Log, owner string, name string, clusterId string) error {
	ssvLog.Infow("recordEvent", "txHash", vLog.TxHash.Hex(), "name", name)
	events := []*store.EventInfo{
		&store.EventInfo{
			BlockNumber: vLog.BlockNumber,
			Owner:       owner,
			TxHash:      vLog.TxHash.Hex(),
			LogIndex:    vLog.Index,
			Action:      name,
			ClusterID:   clusterId,
		},
	}
	return s.store.CreateEvent(events)
}

func (s *SSV) calcLiquidation(clusterId string, owner common.Address, operatorIds []uint64, cluster ISSVNetworkCoreCluster) {
	s.calcLiquidationChan <- Cluster{
		ClusterId:   clusterId,
		Owner:       owner,
		OperatorIds: operatorIds,
		ClusterInfo: cluster,
	}
}
