package alert

import (
	"encoding/hex"
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"github.com/monitorssv/monitorssv/crypto"
	"github.com/monitorssv/monitorssv/eth1/client"
	"github.com/monitorssv/monitorssv/eth1/utils"
	"github.com/monitorssv/monitorssv/store"
	"github.com/robfig/cron/v3"
	"math/big"
	"strings"
	"time"
)

var log = logging.Logger("alarm")

type NetworkFeeChangeNotify struct {
	Block         uint64
	OldNetworkFee *big.Int
	NewNetworkFee *big.Int
}
type OperatorFeeChangeNotify struct {
	Block       uint64
	OperatorId  uint64
	OperatorFee *big.Int
}
type ValidatorProposeBlockNotify struct {
	Epoch     uint64
	Slot      uint64
	ClusterId string
	Index     uint64
}
type ValidatorMissedBlockNotify struct {
	Epoch     uint64
	Slot      uint64
	ClusterId string
	Index     uint64
}
type ValidatorBalanceDeltaNotify struct {
	Epoch     uint64
	ClusterId string
	Index     []uint64
}
type ValidatorSlashNotify struct {
	Epoch     uint64
	ClusterId string
	Index     []uint64
}

type AlarmDaemon struct {
	cron *cron.Cron

	client   *client.Eth1Client
	store    *store.Store
	password string

	networkFeeChangeChan      chan NetworkFeeChangeNotify
	operatorFeeChangeChan     chan OperatorFeeChangeNotify
	validatorProposeBlockChan chan ValidatorProposeBlockNotify
	validatorMissedBlockChan  chan ValidatorMissedBlockNotify
	validatorBalanceDeltaChan chan ValidatorBalanceDeltaNotify
	validatorSlashNotifyChan  chan ValidatorSlashNotify

	close chan struct{}
}

// todo avoid a large number of alarm messages in a short period of time
func NewAlarmDaemon(store *store.Store, client *client.Eth1Client, password string) (*AlarmDaemon, error) {
	c := cron.New()
	alarm := &AlarmDaemon{
		cron:                  c,
		client:                client,
		store:                 store,
		password:              password,
		networkFeeChangeChan:  make(chan NetworkFeeChangeNotify, 1),
		operatorFeeChangeChan: make(chan OperatorFeeChangeNotify, 10),

		validatorProposeBlockChan: make(chan ValidatorProposeBlockNotify, 1),
		validatorMissedBlockChan:  make(chan ValidatorMissedBlockNotify, 1),
		validatorBalanceDeltaChan: make(chan ValidatorBalanceDeltaNotify, 100),
		validatorSlashNotifyChan:  make(chan ValidatorSlashNotify, 100),

		close: make(chan struct{}),
	}

	// check password
	_, err := alarm.getAllAlarmInfos()
	if err != nil {
		return nil, err
	}

	return alarm, nil
}

func (d *AlarmDaemon) NetworkFeeChangeChan() chan<- NetworkFeeChangeNotify {
	return d.networkFeeChangeChan
}

func (d *AlarmDaemon) OperatorFeeChangeChan() chan<- OperatorFeeChangeNotify {
	return d.operatorFeeChangeChan
}

func (d *AlarmDaemon) ValidatorProposeBlockChan() chan<- ValidatorProposeBlockNotify {
	return d.validatorProposeBlockChan
}

func (d *AlarmDaemon) ValidatorMissedBlockChan() chan<- ValidatorMissedBlockNotify {
	return d.validatorMissedBlockChan
}

func (d *AlarmDaemon) ValidatorBalanceDeltaChan() chan<- ValidatorBalanceDeltaNotify {
	return d.validatorBalanceDeltaChan
}

func (d *AlarmDaemon) ValidatorSlashNotifyChan() chan<- ValidatorSlashNotify {
	return d.validatorSlashNotifyChan
}

func (d *AlarmDaemon) Start() {
	_, err := d.cron.AddFunc("0 0 * * *", d.liquidationAlarm)
	if err != nil {
		panic(err)
	}
	_, err = d.cron.AddFunc("0 0 * * *", d.validatorExitedButNotRemovedAlarm)
	if err != nil {
		panic(err)
	}
	_, err = d.cron.AddFunc("0 0 * * 1", d.weeklyReport)
	if err != nil {
		panic(err)
	}

	d.cron.Start()
	go d.alarmDaemonLoop()
}

func (d *AlarmDaemon) Stop() {
	ctx := d.cron.Stop()
	select {
	case <-ctx.Done():
	case <-time.After(time.Minute):
		log.Warn("cron task stop too long")
	}
	close(d.close)
}

// day 0 0 * * *
func (d *AlarmDaemon) liquidationAlarm() {
	curBlock, err := d.client.BlockNumber()
	if err != nil {
		log.Warnw("liquidationAlarm: BlockNumber", "err", err)
		return
	}

	alarmConfigs, err := d.getAllAlarmInfos()
	if err != nil {
		log.Errorw("liquidationAlarm: getAllAlarmInfos", "err", err)
		return
	}

	for _, ac := range alarmConfigs {
		clusterInfos, err := d.store.GetAllClusterByEoaOwner(ac.EoaOwner)
		if err != nil {
			log.Errorw("liquidationAlarm: GetAllClusterByEoaOwner", "err", err)
			continue
		}

		alarm, err := NewAlarm(ac.AlarmType, ac.AlarmChannel)
		if err != nil {
			log.Warnw("liquidationAlarm: NewAlarm", "owner", ac.EoaOwner, "err", err)
			continue
		}

		for _, clusterInfo := range clusterInfos {
			if clusterInfo.ValidatorCount == 0 {
				log.Infow("liquidationAlarm: cluster has no validators, skip", "cluster", clusterInfo.ClusterID)
				continue
			}

			log.Infow("liquidationAlarm", "cluster", clusterInfo.ClusterID, "curBlock", curBlock, "LiquidationBlock", clusterInfo.LiquidationBlock, "ReportLiquidationThreshold", ac.ReportLiquidationThreshold)

			if curBlock+ac.ReportLiquidationThreshold >= clusterInfo.LiquidationBlock {
				onChainBalanceStr := store.CalcClusterOnChainBalance(curBlock, &clusterInfo)
				liquidationMsgFormat := "MonitorSSV: Liquidation Warning!\n  Cluster: %s\n  Cluster Balance: %s ssv\n  Liquidation Block: %d\n  Operational Runway: %d days"

				msg := fmt.Sprintf(liquidationMsgFormat, clusterInfo.ClusterID, onChainBalanceStr, clusterInfo.LiquidationBlock, calcRunway(clusterInfo.LiquidationBlock, curBlock))
				log.Infow("liquidationAlarm", "msg", msg)
				err = alarm.Send(msg)
				if err != nil {
					log.Warnw("liquidationAlarm: Send", "msg", msg, "err", err)
				}
			}
		}
	}
}

// day 0 0 * * *
func (d *AlarmDaemon) validatorExitedButNotRemovedAlarm() {
	alarmConfigs, err := d.getAllAlarmInfos()
	if err != nil {
		log.Errorw("validatorExitedButNotRemovedAlarm: getAllAlarmInfos", "err", err)
		return
	}

	for _, ac := range alarmConfigs {
		if !ac.ReportExitedButNotRemoved {
			continue
		}

		clusterInfos, err := d.store.GetAllClusterByEoaOwner(ac.EoaOwner)
		if err != nil {
			log.Errorw("validatorExitedButNotRemovedAlarm: GetAllClusterByEoaOwner", "err", err)
			continue
		}

		alarm, err := NewAlarm(ac.AlarmType, ac.AlarmChannel)
		if err != nil {
			log.Warnw("validatorExitedButNotRemovedAlarm: NewAlarm", "owner", ac.EoaOwner, "err", err)
			continue
		}

		for _, clusterInfo := range clusterInfos {
			if clusterInfo.ValidatorCount == 0 {
				continue
			}

			validators, err := d.store.GetActiveButExitedValidatorsByClusterId(clusterInfo.ClusterID)
			if err != nil {
				log.Errorw("validatorExitedButNotRemovedAlarm: GetActiveButExitedValidatorsByClusterId", "clusterId", clusterInfo.ClusterID, "err", err)
				continue
			}
			if len(validators) == 0 {
				continue
			}
			indexs := ""
			for _, validator := range validators {
				if indexs == "" {
					indexs = fmt.Sprintf("%d", validator.ValidatorIndex)
					continue
				}
				indexs = fmt.Sprintf("%s, %d", indexs, validator.ValidatorIndex)
			}
			validatorNotRemovedMsgFormat := "MonitorSSV: Validator NotRemoved Warning!\n  Cluster: %s\n  Validators: %s"
			msg := fmt.Sprintf(validatorNotRemovedMsgFormat, clusterInfo.ClusterID, indexs)
			log.Infow("validatorExitedButNotRemovedAlarm", "msg", msg)
			err = alarm.Send(msg)
			if err != nil {
				log.Warnw("validatorExitedButNotRemovedAlarm", "msg", msg, "err", err)
			}
		}
	}
}

// monday 0 0 * * 1
func (d *AlarmDaemon) weeklyReport() {
	curBlock, err := d.client.BlockNumber()
	if err != nil {
		log.Warnw("weeklyReport: BlockNumber", "err", err)
		return
	}

	alarmConfigs, err := d.getAllAlarmInfos()
	if err != nil {
		log.Errorw("weeklyReport: getAllAlarmInfos", "err", err)
		return
	}

	for _, ac := range alarmConfigs {
		if !ac.ReportWeekly {
			continue
		}

		clusterInfos, err := d.store.GetAllClusterByEoaOwner(ac.EoaOwner)
		if err != nil {
			log.Errorw("weeklyReport: GetAllClusterByEoaOwner", "err", err)
			continue
		}

		alarm, err := NewAlarm(ac.AlarmType, ac.AlarmChannel)
		if err != nil {
			log.Warnw("weeklyReport: NewAlarm", "owner", ac.EoaOwner, "err", err)
			continue
		}

		for _, clusterInfo := range clusterInfos {
			if clusterInfo.ValidatorCount == 0 {
				log.Infow("weeklyReport: cluster has no validators, skip", "cluster", clusterInfo.ClusterID)
				continue
			}

			log.Infow("weeklyReport: clusterInfo", "cluster", clusterInfo.ClusterID, "LiquidationBlock", clusterInfo.LiquidationBlock)

			onChainBalanceStr := store.CalcClusterOnChainBalance(curBlock, &clusterInfo)
			weeklyReportMsgFormat := "MonitorSSV: Weekly Report!\n  Cluster: %s\n  Validator Count: %d\n  Cluster Balance: %s ssv\n  Liquidation Block: %d\n  Operational Runway: %d days"
			msg := fmt.Sprintf(weeklyReportMsgFormat, clusterInfo.ClusterID, clusterInfo.ValidatorCount, onChainBalanceStr, clusterInfo.LiquidationBlock, calcRunway(clusterInfo.LiquidationBlock, curBlock))
			log.Infow("weeklyReport", "msg", msg)
			err = alarm.Send(msg)
			if err != nil {
				log.Warnw("weeklyReport: Send", "msg", msg, "err", err)
			}
		}
	}
}

func (d *AlarmDaemon) alarmDaemonLoop() {
	for {
		select {
		case <-d.close:
			return
		case networkFeeChange := <-d.networkFeeChangeChan:
			log.Infow("alarmDaemonLoop", "networkFeeChange", networkFeeChange)
			go func() {
				<-time.After(10 * time.Minute)
				d.networkFeeChangeAlarm(networkFeeChange)
			}()
		case operatorFeeChange := <-d.operatorFeeChangeChan:
			log.Infow("alarmDaemonLoop", "operatorFeeChange", operatorFeeChange)
			go func() {
				<-time.After(10 * time.Minute)
				d.operatorFeeChangeAlarm(operatorFeeChange)
			}()
		case validatorProposeBlock := <-d.validatorProposeBlockChan:
			log.Infow("alarmDaemonLoop", "validatorProposeBlockChan", validatorProposeBlock)
			d.proposeBlockAlarm(validatorProposeBlock)
		case validatorMissedBlock := <-d.validatorMissedBlockChan:
			log.Infow("alarmDaemonLoop", "validatorMissedBlockChan", validatorMissedBlock)
			d.missedBlockAlarm(validatorMissedBlock)
		case validatorBalanceDelta := <-d.validatorBalanceDeltaChan:
			log.Infow("alarmDaemonLoop", "validatorBalanceDelta", validatorBalanceDelta)
			d.validatorBalanceDeltaAlarm(validatorBalanceDelta)
		case validatorSlash := <-d.validatorSlashNotifyChan:
			log.Infow("alarmDaemonLoop", "validatorSlashNotifyChan", validatorSlash)
			d.validatorSlashAlarm(validatorSlash)
		}
	}
}

func (d *AlarmDaemon) operatorFeeChangeAlarm(operatorFeeChange OperatorFeeChangeNotify) {
	curBlock, err := d.client.BlockNumber()
	if err != nil {
		log.Warnw("operatorFeeChangeAlarm: BlockNumber", "err", err)
		return
	}

	alarmConfigs, err := d.getAllAlarmInfos()
	if err != nil {
		log.Errorw("operatorFeeChangeAlarm: getAllAlarmInfos", "err", err)
		return
	}

	operatorInfo, err := d.store.GetOperatorByOperatorId(operatorFeeChange.OperatorId)
	if err != nil {
		log.Errorw("operatorFeeChangeAlarm: GetOperatorByOperatorId", "err", err)
		return
	}

	clusterIds := strings.Split(operatorInfo.ClusterIds, ",")

	for _, clusterId := range clusterIds {
		clusterInfo, err := d.store.GetClusterByClusterId(clusterId)
		if err != nil {
			log.Errorw("operatorFeeChangeAlarm: GetClusterByClusterId", "err", err)
			continue
		}
		if clusterInfo.ValidatorCount == 0 {
			log.Infow("operatorFeeChangeAlarm: cluster has no validators, skip", "cluster", clusterInfo.ClusterID)
			continue
		}

		if clusterInfo.EoaOwner == "0x" {
			log.Warnw("operatorFeeChangeAlarm: GetClusterByClusterId", "clusterId", clusterId, "eoaOwner", clusterInfo.EoaOwner)
			continue
		}

		if ac, ok := alarmConfigs[clusterInfo.EoaOwner]; ok {
			if !ac.ReportOperatorFeeChange {
				log.Infow("operatorFeeChangeAlarm: ReportOperatorFeeChange not set", "eoaOwner", clusterInfo.EoaOwner)
				continue
			}

			alarm, err := NewAlarm(ac.AlarmType, ac.AlarmChannel)
			if err != nil {
				log.Warnw("operatorFeeChangeAlarm: NewAlarm", "owner", ac.EoaOwner, "err", err)
				continue
			}

			operatorFee := "0"
			if operatorFeeChange.OperatorFee.Uint64() != 0 {
				fee := big.NewInt(0).Mul(operatorFeeChange.OperatorFee, big.NewInt(2613400))
				operatorFee = utils.ToSSV(fee, "%.2f")
			}

			onChainBalanceStr := store.CalcClusterOnChainBalance(curBlock, clusterInfo)
			weeklyReportMsgFormat := "MonitorSSV: OperatorFee Change Notice!\n  Operator ID: %d\n  Operator Fee: %s\n  Cluster: %s\n  Validator Count: %d\n  Cluster Balance: %s ssv\n  Liquidation Block: %d\n  Operational Runway: %d days"
			msg := fmt.Sprintf(weeklyReportMsgFormat, operatorFeeChange.OperatorId, operatorFee, clusterInfo.ClusterID, clusterInfo.ValidatorCount, onChainBalanceStr, clusterInfo.LiquidationBlock, calcRunway(clusterInfo.LiquidationBlock, curBlock))
			log.Infow("operatorFeeChangeAlarm", "msg", msg)
			err = alarm.Send(msg)
			if err != nil {
				log.Warnw("operatorFeeChangeAlarm: Send", "msg", msg, "err", err)
			}
		}
	}
}

func (d *AlarmDaemon) networkFeeChangeAlarm(networkFeeChange NetworkFeeChangeNotify) {
	curBlock, err := d.client.BlockNumber()
	if err != nil {
		log.Warnw("networkFeeChangeAlarm: BlockNumber", "err", err)
		return
	}

	alarmConfigs, err := d.getAllAlarmInfos()
	if err != nil {
		log.Errorw("networkFeeChangeAlarm: getAllAlarmInfos", "err", err)
		return
	}

	for eoaOwner, ac := range alarmConfigs {
		if !ac.ReportNetworkFeeChange {
			log.Infow("networkFeeChangeAlarm: ReportNetworkFeeChange not set", "eoaOwner", eoaOwner)
			continue
		}
		clusterInfos, err := d.store.GetAllClusterByEoaOwner(eoaOwner)
		if err != nil {
			log.Errorw("networkFeeChangeAlarm: GetAllClusterByEoaOwner", "err", err)
			continue
		}

		alarm, err := NewAlarm(ac.AlarmType, ac.AlarmChannel)
		if err != nil {
			log.Warnw("networkFeeChangeAlarm: NewAlarm", "owner", ac.EoaOwner, "err", err)
			continue
		}

		for _, clusterInfo := range clusterInfos {
			newNetworkFee := "0"
			if networkFeeChange.NewNetworkFee.Uint64() != 0 {
				fee := big.NewInt(0).Mul(networkFeeChange.NewNetworkFee, big.NewInt(2613400))
				newNetworkFee = utils.ToSSV(fee, "%.2f")
			}
			oldNetworkFee := "0"
			if networkFeeChange.OldNetworkFee.Uint64() != 0 {
				fee := big.NewInt(0).Mul(networkFeeChange.OldNetworkFee, big.NewInt(2613400))
				oldNetworkFee = utils.ToSSV(fee, "%.2f")
			}

			onChainBalanceStr := store.CalcClusterOnChainBalance(curBlock, &clusterInfo)
			weeklyReportMsgFormat := "MonitorSSV: NetworkFee Change Notice!\n  Old Network Fee: %s\n  New Network Fee: %s\n  Cluster: %s\n  Validator Count: %d\n  Cluster Balance: %s ssv\n  Liquidation Block: %d\n  Operational Runway: %d days"
			msg := fmt.Sprintf(weeklyReportMsgFormat, oldNetworkFee, newNetworkFee, clusterInfo.ClusterID, clusterInfo.ValidatorCount, onChainBalanceStr, clusterInfo.LiquidationBlock, calcRunway(clusterInfo.LiquidationBlock, curBlock))
			log.Infow("networkFeeChangeAlarm", "msg", msg)
			err = alarm.Send(msg)
			if err != nil {
				log.Warnw("networkFeeChangeAlarm: Send", "msg", msg, "err", err)
			}
		}
	}
}

func (d *AlarmDaemon) proposeBlockAlarm(validatorProposeBlock ValidatorProposeBlockNotify) {
	ac, err := d.getClusterAlarmInfo(validatorProposeBlock.ClusterId)
	if err != nil {
		log.Errorw("proposeBlockAlarm: getClusterAlarmInfo", "err", err)
		return
	}

	if ac != nil {
		if !ac.ReportProposeBlock {
			log.Infow("proposeBlockAlarm: ReportProposeBlock not set", "eoaOwner", ac.EoaOwner)
			return
		}

		alarm, err := NewAlarm(ac.AlarmType, ac.AlarmChannel)
		if err != nil {
			log.Warnw("proposeBlockAlarm: NewAlarm", "owner", ac.EoaOwner, "err", err)
			return
		}

		reportProposeBlockMsgFormat := "MonitorSSV: Validator propose block!\n  Cluster ID: %s\n  Validator Index: %d\n  Epoch: %d\n  Slot: %d\n"
		msg := fmt.Sprintf(reportProposeBlockMsgFormat, validatorProposeBlock.ClusterId, validatorProposeBlock.Index, validatorProposeBlock.Epoch, validatorProposeBlock.Slot)
		log.Infow("proposeBlockAlarm", "msg", msg)
		err = alarm.Send(msg)
		if err != nil {
			log.Warnw("proposeBlockAlarm: Send", "msg", msg, "err", err)
		}
	}
}

func (d *AlarmDaemon) missedBlockAlarm(validatorMissedBlock ValidatorMissedBlockNotify) {
	ac, err := d.getClusterAlarmInfo(validatorMissedBlock.ClusterId)
	if err != nil {
		log.Errorw("missedBlockAlarm: getClusterAlarmInfo", "err", err)
		return
	}

	if ac != nil {
		if !ac.ReportMissedBlock {
			log.Infow("missedBlockAlarm: ReportProposeBlock not set", "eoaOwner", ac.EoaOwner)
			return
		}

		alarm, err := NewAlarm(ac.AlarmType, ac.AlarmChannel)
		if err != nil {
			log.Warnw("missedBlockAlarm: NewAlarm", "owner", ac.EoaOwner, "err", err)
			return
		}

		reportMissedBlockMsgFormat := "MonitorSSV: Validator missed block!\n  Cluster ID: %s\n  Validator Index: %d\n  Epoch: %d\n  Slot: %d\n"
		msg := fmt.Sprintf(reportMissedBlockMsgFormat, validatorMissedBlock.ClusterId, validatorMissedBlock.Index, validatorMissedBlock.Epoch, validatorMissedBlock.Slot)
		log.Infow("missedBlockAlarm", "msg", msg)
		err = alarm.Send(msg)
		if err != nil {
			log.Warnw("missedBlockAlarm: Send", "msg", msg, "err", err)
		}
	}
}

func (d *AlarmDaemon) validatorBalanceDeltaAlarm(validatorBalanceDelta ValidatorBalanceDeltaNotify) {
	ac, err := d.getClusterAlarmInfo(validatorBalanceDelta.ClusterId)
	if err != nil {
		log.Errorw("validatorBalanceDeltaAlarm: getClusterAlarmInfo", "err", err)
		return
	}

	if ac != nil {
		if !ac.ReportBalanceDecrease {
			log.Infow("validatorBalanceDeltaAlarm: ReportProposeBlock not set", "eoaOwner", ac.EoaOwner)
			return
		}

		alarm, err := NewAlarm(ac.AlarmType, ac.AlarmChannel)
		if err != nil {
			log.Warnw("validatorBalanceDeltaAlarm: NewAlarm", "owner", ac.EoaOwner, "err", err)
			return
		}

		reportBalanceDecreaseMsgFormat := "MonitorSSV: Validator balance decreases!\n  Cluster ID: %s\n  Epoch: %d\n  Validator Index: %v\n"
		for i, batch := range chunkSlice(validatorBalanceDelta.Index, 100) {
			msg := fmt.Sprintf(reportBalanceDecreaseMsgFormat, validatorBalanceDelta.ClusterId, validatorBalanceDelta.Epoch, batch)
			log.Infow("validatorBalanceDeltaAlarm", "batch", i, "msg", msg)
			err = alarm.Send(msg)
			if err != nil {
				log.Warnw("validatorBalanceDeltaAlarm: Send", "msg", msg, "err", err)
			}
		}
	}
}

func (d *AlarmDaemon) validatorSlashAlarm(validatorSlashNotify ValidatorSlashNotify) {
	ac, err := d.getClusterAlarmInfo(validatorSlashNotify.ClusterId)
	if err != nil {
		log.Errorw("validatorSlashAlarm: getClusterAlarmInfo", "err", err)
		return
	}

	if ac != nil {
		alarm, err := NewAlarm(ac.AlarmType, ac.AlarmChannel)
		if err != nil {
			log.Warnw("validatorSlashAlarm: NewAlarm", "owner", ac.EoaOwner, "err", err)
			return
		}

		reportValidatorSlashMsgFormat := "MonitorSSV: Validator slashed!\n  Cluster ID: %s\n  Epoch: %d\n  Validator Index: %v\n"
		for i, batch := range chunkSlice(validatorSlashNotify.Index, 100) {
			msg := fmt.Sprintf(reportValidatorSlashMsgFormat, validatorSlashNotify.ClusterId, validatorSlashNotify.Epoch, batch)
			log.Infow("validatorSlashAlarm", "batch", i, "msg", msg)
			err = alarm.Send(msg)
			if err != nil {
				log.Warnw("validatorSlashAlarm: Send", "msg", msg, "err", err)
			}
		}
	}
}

type alarmConfig struct {
	EoaOwner                   string `json:"eoa_owner"`
	AlarmType                  int    `json:"alarm_type"`
	AlarmChannel               string `json:"alarm_channel"`
	ReportLiquidationThreshold uint64 `json:"report_liquidation_threshold"`
	ReportOperatorFeeChange    bool   `json:"report_operator_fee_change"`
	ReportNetworkFeeChange     bool   `json:"report_network_fee_change"`
	ReportProposeBlock         bool   `json:"report_propose_block"`
	ReportMissedBlock          bool   `json:"report_missed_block"`
	ReportBalanceDecrease      bool   `json:"report_balance_decrease"`
	ReportExitedButNotRemoved  bool   `json:"report_exited_but_not_removed"`
	ReportWeekly               bool   `json:"report_weekly"`
}

func (d *AlarmDaemon) getAllAlarmInfos() (map[string]alarmConfig, error) {
	var alarmConfigMap = make(map[string]alarmConfig)
	key := crypto.GenerateEncryptKey([]byte(d.password))

	alarmInfos, err := d.store.GetAllAlarmInfos()
	if err != nil {
		log.Errorw("GetAllAlarmInfos", "err", err)
		return nil, err
	}

	for _, alarmInfo := range alarmInfos {
		ac, err := decryptAlarmInfo(key, &alarmInfo)
		if err != nil {
			log.Errorw("decryptAlarmInfo", "err", err)
			return nil, err
		}

		alarmConfigMap[alarmInfo.EoaOwner] = *ac
	}

	return alarmConfigMap, nil
}

func (d *AlarmDaemon) getClusterAlarmInfo(cluster string) (*alarmConfig, error) {
	clusterInfo, err := d.store.GetClusterByClusterId(cluster)
	if err != nil {
		log.Errorw("GetClusterByClusterId", "err", err)
		return nil, err
	}

	alarmInfo, err := d.store.GetAlarmByEoaOwner(clusterInfo.EoaOwner)
	if err != nil {
		log.Errorw("GetAlarmByEoaOwner", "err", err)
		return nil, err
	}
	if alarmInfo == nil {
		return nil, nil
	}

	key := crypto.GenerateEncryptKey([]byte(d.password))

	return decryptAlarmInfo(key, alarmInfo)
}

func decryptAlarmInfo(key []byte, alarmInfo *store.AlarmInfo) (*alarmConfig, error) {
	encryptedData, err := hex.DecodeString(alarmInfo.AlarmChannel)
	if err != nil {
		log.Warnw("decryptAlarmInfo: DecodeString", "err", err)
		return nil, err
	}

	alarmChannel, err := crypto.DecryptData(encryptedData, key)
	if err != nil {
		log.Warnw("decryptAlarmInfo: DecryptData", "err", err)
		return nil, err
	}
	alarmChannelHash := crypto.Hash256(alarmChannel)
	if hex.EncodeToString(alarmChannelHash) != alarmInfo.AlarmChannelHash {
		log.Warnw("decryptAlarmInfo: alarmChannelHash does not match", "want", alarmInfo.AlarmChannelHash, "get", hex.EncodeToString(alarmChannelHash))
		return nil, err
	}

	var ac alarmConfig
	ac.EoaOwner = alarmInfo.EoaOwner
	ac.AlarmType = alarmInfo.AlarmType
	ac.AlarmChannel = string(alarmChannel)
	ac.ReportLiquidationThreshold = alarmInfo.ReportLiquidationThreshold
	ac.ReportOperatorFeeChange = alarmInfo.ReportOperatorFeeChange
	ac.ReportNetworkFeeChange = alarmInfo.ReportNetworkFeeChange
	ac.ReportProposeBlock = alarmInfo.ReportProposeBlock
	ac.ReportMissedBlock = alarmInfo.ReportMissedBlock
	ac.ReportBalanceDecrease = alarmInfo.ReportBalanceDecrease
	ac.ReportExitedButNotRemoved = alarmInfo.ReportExitedButNotRemoved
	ac.ReportWeekly = alarmInfo.ReportWeekly

	return &ac, nil
}

func chunkSlice(slice []uint64, chunkSize int) [][]uint64 {
	var chunks [][]uint64
	if chunkSize <= 0 {
		return [][]uint64{slice}
	}

	length := len(slice)
	numChunks := (length + chunkSize - 1) / chunkSize

	chunks = make([][]uint64, 0, numChunks)
	for i := 0; i < length; i += chunkSize {
		end := i + chunkSize
		if end > length {
			end = length
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func calcRunway(liquidationBlock, curBlock uint64) uint64 {
	runway := uint64(0)
	if liquidationBlock > curBlock {
		runway = (liquidationBlock - curBlock) / 7200
	}
	return runway
}
