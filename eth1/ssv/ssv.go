package ssv

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	logging "github.com/ipfs/go-log/v2"
	"github.com/monitorssv/monitorssv/alert"
	"github.com/monitorssv/monitorssv/config"
	"github.com/monitorssv/monitorssv/eth1/client"
	"github.com/monitorssv/monitorssv/store"
	"math"
	"math/big"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var ssvLog = logging.Logger("ssv-scan")

type SSV struct {
	cfg *config.Config

	client *client.Eth1Client
	store  *store.Store

	ssvNetworkAddr    common.Address
	ssvNetworkViewAdd common.Address

	lastProcessedBlock uint64
	isSynced           *atomic.Bool

	calcLiquidationChan           chan Cluster
	calcAllClusterLiquidationChan chan uint64

	networkFeeChangeAlarmChan  chan<- alert.NetworkFeeChangeNotify
	operatorFeeChangeAlarmChan chan<- alert.OperatorFeeChangeNotify

	events map[common.Hash]abi.Event
	close  chan struct{}
}

func NewSSV(cfg *config.Config, client *client.Eth1Client, store *store.Store, alarm *alert.AlarmDaemon) (*SSV, error) {
	contractInfo, err := GetSSVContractInfo(cfg.Network)
	if err != nil {
		return nil, err
	}

	lastProcessedBlock := contractInfo.DeployBlock
	if block, _, err := store.GetScanPoint(); err == nil && block != 0 {
		lastProcessedBlock = block
	}

	ssv := &SSV{
		cfg:                           cfg,
		client:                        client,
		store:                         store,
		ssvNetworkAddr:                contractInfo.SSVNetwork,
		ssvNetworkViewAdd:             contractInfo.SSVNetworkView,
		lastProcessedBlock:            lastProcessedBlock,
		isSynced:                      new(atomic.Bool),
		calcLiquidationChan:           make(chan Cluster, 100),
		calcAllClusterLiquidationChan: make(chan uint64, 100),
		networkFeeChangeAlarmChan:     alarm.NetworkFeeChangeChan(),
		operatorFeeChangeAlarmChan:    alarm.OperatorFeeChangeChan(),
		events:                        GetAllSSVEvent(),
		close:                         make(chan struct{}),
	}
	ssv.isSynced.Store(false)
	return ssv, nil
}

func (s *SSV) Start() {
	go s.ScanSSVEventLoop()
	go s.UpdateClusterLiquidationLoop()
	go s.UpdateOperatorLoop()
	go s.UpdateClusterEoaOwnerLoop()
}

func (s *SSV) Stop() {
	close(s.close)
}

func (s *SSV) GetLastProcessedBlock() uint64 {
	return s.lastProcessedBlock
}

func (s *SSV) GetCfg() *config.Config {
	return s.cfg
}

const safeBlockRange uint64 = 5

func (s *SSV) ScanSSVEventLoop() {
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-s.close:
			return
		case <-ticker.C:
			curBlock, err := s.client.BlockNumber()
			if err != nil {
				ssvLog.Errorf("error getting block number: %s", err)
				continue
			}
			toBlock := curBlock - safeBlockRange
			if s.lastProcessedBlock+1 > toBlock {
				continue
			}

			ssvLog.Infow("ScanSSVEventLoop", "lastProcessedBlock", s.lastProcessedBlock, "toBlock", toBlock, "curBlock", curBlock)

			lastProcessedBlock, err := s.ScanSSVEvent(s.lastProcessedBlock+1, toBlock)
			if err != nil {
				ssvLog.Errorf("ScanSSVEventLoop: ScanSSVEvent failed: %s", err)
			}

			if lastProcessedBlock > s.lastProcessedBlock {
				s.lastProcessedBlock = lastProcessedBlock
				err = s.store.UpdateScanEth1Block(lastProcessedBlock)
				if err != nil {
					ssvLog.Warnf("ScanSSVEventLoop: UpdateScanEth1Block failed: %s", err)
				}
			}

			if !s.isSynced.Load() {
				if lastProcessedBlock >= toBlock {
					ssvLog.Infow("ScanSSVEventLoop: Sync completed", "lastProcessedBlock", lastProcessedBlock, "toBlock", toBlock)
					s.isSynced.Store(true)
				}
			}
		}
	}
}

func (s *SSV) UpdateClusterLiquidationLoop() {
	var calcLiquidationQueue = map[string]Cluster{}
	ticker := time.NewTicker(24 * time.Second)
	for {
		select {
		case <-s.close:
			return
		case <-ticker.C:
			if !s.isSynced.Load() {
				continue
			}

			for _, cluster := range calcLiquidationQueue {
				ssvLog.Infow("will calc liquidation block", "clusterId", cluster.ClusterId)

				err := s.calcAndUpdateClusterLiquidation(cluster)
				if err != nil {
					ssvLog.Warnf("calcAndUpdateClusterLiquidation failed: %s", err)
					continue
				}

				delete(calcLiquidationQueue, cluster.ClusterId)
			}

		case cluster := <-s.calcLiquidationChan:
			ssvLog.Infow("updateClusterLiquidationLoop: notify", "clusterId", cluster.ClusterId)
			calcLiquidationQueue[cluster.ClusterId] = cluster

		case operatorId := <-s.calcAllClusterLiquidationChan:
			if operatorId == 0 {
				err := s.calcAllClusterLiquidation()
				if err != nil {
					ssvLog.Warnf("CalcAllClusterLiquidation failed: %s", err)
				}
				continue
			}
			err := s.operatorFeeUpdateCalcClusterLiquidation(operatorId)
			if err != nil {
				ssvLog.Warnf("CalcAllClusterLiquidation failed: %s", err)
			}
		}
	}
}

func (s *SSV) calcAllClusterLiquidation() error {
	ssvLog.Info("will calc all cluster liquidation block")

	clusterInfos, err := s.store.GetAllClusters()
	if err != nil {
		return err
	}
	for _, clusterInfo := range clusterInfos {
		if !clusterInfo.Active {
			ssvLog.Infow("cluster liquidated", "clusterId", clusterInfo.ClusterID, "validatorCount", clusterInfo.ValidatorCount)
			continue
		}
		operatorIds, err := getOperatorIds(clusterInfo.OperatorIds)
		if err != nil {
			return err
		}
		balance, isOk := big.NewInt(0).SetString(clusterInfo.Balance, 10)
		if !isOk {
			return errors.New("failed to parse balance")
		}
		err = s.calcAndUpdateClusterLiquidation(Cluster{
			ClusterId:   clusterInfo.ClusterID,
			Owner:       common.HexToAddress(clusterInfo.Owner),
			OperatorIds: operatorIds,
			ClusterInfo: ISSVNetworkCoreCluster{
				ValidatorCount:  clusterInfo.ValidatorCount,
				NetworkFeeIndex: clusterInfo.NetworkFeeIndex,
				Index:           clusterInfo.Index,
				Active:          clusterInfo.Active,
				Balance:         balance,
			},
		})
	}

	return nil
}

func (s *SSV) operatorFeeUpdateCalcClusterLiquidation(operatorId uint64) error {
	operatorInfo, err := s.store.GetOperatorByOperatorId(operatorId)
	if err != nil {
		return err
	}
	if operatorInfo.ClusterIds == "" {
		return nil
	}
	clusterIds := strings.Split(operatorInfo.ClusterIds, ",")
	ssvLog.Infow("will calc cluster liquidation block", "clusterIds", clusterIds)

	for _, clusterId := range clusterIds {
		clusterInfo, err := s.store.GetClusterByClusterId(clusterId)
		if err != nil {
			return err
		}

		operatorIds, err := getOperatorIds(clusterInfo.OperatorIds)
		if err != nil {
			return err
		}

		balance, isOk := big.NewInt(0).SetString(clusterInfo.Balance, 10)
		if !isOk {
			return errors.New("failed to parse balance")
		}
		err = s.calcAndUpdateClusterLiquidation(Cluster{
			ClusterId:   clusterInfo.ClusterID,
			Owner:       common.HexToAddress(clusterInfo.Owner),
			OperatorIds: operatorIds,
			ClusterInfo: ISSVNetworkCoreCluster{
				ValidatorCount:  clusterInfo.ValidatorCount,
				NetworkFeeIndex: clusterInfo.NetworkFeeIndex,
				Index:           clusterInfo.Index,
				Active:          clusterInfo.Active,
				Balance:         balance,
			},
		})
	}

	return nil
}

func getOperatorIds(operatorIdsStr string) ([]uint64, error) {
	operatorIds := make([]uint64, 0)
	for _, operatorIdStr := range strings.Split(operatorIdsStr, ",") {
		operatorId, err := strconv.ParseUint(operatorIdStr, 10, 64)
		if err != nil {
			return nil, err
		}
		operatorIds = append(operatorIds, operatorId)
	}
	return operatorIds, nil
}

func (s *SSV) calcAndUpdateClusterLiquidation(cluster Cluster) error {
	liquidationBlock, curBlock, burnFee, onChainBalance, err := s.CalcLiquidation(cluster)
	if err != nil {
		if errors.Is(err, noValidatorErr) || errors.Is(err, alreadyLiquidatedErr) || errors.Is(err, canLiquidatedErr) {
			return nil
		}

		ssvLog.Warnw("failed to calc liquidation block", "clusterId", cluster.ClusterId, "err", err)
		return err
	}
	err = s.store.UpdateClusterLiquidationInfo(cluster.ClusterId, liquidationBlock, curBlock, burnFee, onChainBalance)
	if err != nil {
		ssvLog.Warnw("failed to update liquidation block", "clusterId", cluster.ClusterId, "err", err)
		return err
	}
	ssvLog.Infow("cluster liquidation block", "clusterId", cluster.ClusterId, "liquidationBlock", liquidationBlock, "curBlock", curBlock, "burnFee", burnFee, "onChainBalance", onChainBalance)

	return nil
}

func (s *SSV) UpdateOperatorLoop() {
	ticker := time.NewTicker(8 * time.Hour)
	for {
		select {
		case <-s.close:
			return
		case <-ticker.C:
			if !s.isSynced.Load() {
				continue
			}

			s.updateOperatorName()
			s.updateOperatorEarning()
		}
	}
}

func (s *SSV) updateOperatorEarning() {
	itemsPerPage := 100
	page := 1

	for {
		operators, totalCount, err := s.store.GetOperators(page, itemsPerPage)
		if err != nil {
			ssvLog.Errorw("failed to get operators", "err", err)
			return
		}

		totalPages := int(math.Ceil(float64(totalCount) / float64(itemsPerPage)))
		ssvLog.Infow("updateOperatorEarning", "page", page, "totalPages", totalPages, "itemsPerPage", itemsPerPage)

		var operatorIds = make([]uint64, 0)
		for _, operator := range operators {
			if operator.OperatorFee == "0" && operator.OperatorEarnings == "0" {
				continue
			}
			operatorIds = append(operatorIds, operator.OperatorId)
		}

		if len(operatorIds) > 0 {
			operatorEarnings, err := s.GetOperatorEarnings(operatorIds)
			if err != nil {
				ssvLog.Warnw("failed to get operator earnings", "err", err)
			} else {
				for i, operatorEarning := range operatorEarnings {
					operatorId := operatorIds[i]
					err = s.store.UpdateOperatorEarning(operatorId, operatorEarning)
					if err != nil {
						ssvLog.Warnw("failed to update operator earning", "err", err)
					}
				}
			}
		}

		if page*itemsPerPage >= int(totalCount) {
			break
		}
		page++
	}
}

func (s *SSV) updateOperatorName() {
	maxOperatorId, err := s.store.GetMaxOperatorId()
	if err != nil {
		ssvLog.Warn("failed to get max operator id: %v", err)
		return
	}
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for i := uint64(1); i <= maxOperatorId; i++ {
		<-ticker.C

		operatorName, err := GetOperatorName(s.cfg.Network, i)
		if err != nil {
			ssvLog.Warn("failed to get operator name: %v (operator id: %d)", err, i)
			continue
		}
		ssvLog.Infow("get operator name", "operatorId", i, "operatorName", operatorName)

		if operatorName == "" {
			operatorName = fmt.Sprintf("Operator-%d", i)
		}
		err = s.store.UpdateOperatorName(i, operatorName)
		if err != nil {
			ssvLog.Warn("failed to update operator name: %v (operator id: %d)", err, i)
			continue
		}
	}
}

func (s *SSV) UpdateClusterEoaOwnerLoop() {
	ticker := time.NewTicker(24 * time.Second)
	for {
		select {
		case <-s.close:
			return
		case <-ticker.C:
			if !s.isSynced.Load() {
				continue
			}

			s.updateClusterEoaOwner()
		}
	}

}

func (s *SSV) updateClusterEoaOwner() {
	owners, err := s.store.GetNoUpdatedClustersOwner()
	if err != nil {
		ssvLog.Warnf("failed to GetNoUpdatedClustersOwner: %v", err)
		return
	}
	if len(owners) == 0 {
		return
	}

	var contractOwners = make([]string, 0)
	for _, owner := range owners {
		code, err := s.client.CodeAt(owner)
		if err != nil {
			ssvLog.Warnf("failed to CodeAt: %v", err)
			continue
		}

		if len(code) == 0 {
			ssvLog.Infow("cluster owner", "owner", owner, "eoaOwner", owner)
			err = s.store.UpdateClusterEoaOwner(owner, owner)
			if err != nil {
				ssvLog.Warnf("failed to UpdateClusterEoaOwner: %v", err)
				continue
			}
			continue
		}

		contractOwners = append(contractOwners, owner)
	}

	if len(contractOwners) == 0 {
		return
	}

	// Etherscan only provides an API to obtain the creator of the mainnet network.
	// Obtaining the creator through the Ethereum full node API will generate a large number of polling requests.
	// Since the monitoring service of the holesky network is not necessary, the contract owner is not obtained here.

	ticker := time.NewTicker(250 * time.Millisecond) // limit 5 calls/s
	defer ticker.Stop()

	for _, owner := range contractOwners {
		<-ticker.C

		contractCreator := "0x"
		if s.cfg.Network == "mainnet" {
			info, err := GetContractCreator(owner, s.cfg.EtherScan.Endpoint, s.cfg.EtherScan.ApiKey)
			if err != nil {
				ssvLog.Warnf("failed to GetContractCreator: %v", err)
				continue
			}

			ssvLog.Infow("cluster owner", "owner", info.ContractAddress, "eoaOwner", info.ContractCreator)
			contractCreator = info.ContractCreator
		} else { // holesky
			contractCreator = "contract"
			ssvLog.Infow("cluster owner: holesky", "owner", owner, "eoaOwner", contractCreator)
		}

		err = s.store.UpdateClusterEoaOwner(owner, contractCreator)
		if err != nil {
			ssvLog.Warnf("failed to UpdateClusterEoaOwner: %v", err)
			continue
		}
	}
}
