package service

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/monitorssv/monitorssv/alert"
	"github.com/monitorssv/monitorssv/crypto"
	"github.com/monitorssv/monitorssv/store"
	"strconv"
)

func (ms *MonitorSSV) GetClusterMonitorInfo(c *gin.Context) {
	owner := c.DefaultQuery("owner", "")
	monitorLog.Infow("GetClusterMonitor", "owner", owner)
	if owner == "" {
		ReturnErr(c, badRequestRes)
		return
	}

	block := ms.ssv.GetLastProcessedBlock()

	alarmInfo, err := ms.store.GetAlarmByEoaOwner(owner)
	if err != nil {
		monitorLog.Errorw("GetClusterMonitorInfo: GetAllClusterByEoaOwner", "err", err.Error())
		ReturnErr(c, serverErrRes)
		return
	}

	isMonitoring := alarmInfo != nil

	totalClusterCount, totalActiveClusterCount, err := ms.getOwnerClusterInfo(owner)
	if err != nil {
		monitorLog.Errorw("getOwnerClusterInfo", "err", err.Error())
		ReturnErr(c, serverErrRes)
		return
	}

	ReturnOk(c, gin.H{
		"totalClusters":      totalClusterCount,
		"totalActiveCluster": totalActiveClusterCount,
		"isMonitoring":       isMonitoring,
		"block":              block,
	})
}

var getMonitorConfigFormat = "Signature required for cluster ownership. Block: %d"

type MonitorConfig struct {
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

func (ms *MonitorSSV) DeleteMonitorConfig(c *gin.Context) {
	type Request struct {
		Owner     string `json:"owner"`
		Signature string `json:"signature"`
		Block     uint64 `json:"block"`
	}

	param := Request{}
	err := c.ShouldBind(&param)
	if err != nil {
		monitorLog.Warnw("DeleteMonitorConfig", "err", err)
		ReturnErr(c, badRequestRes)
		return
	}

	processedBlock := ms.ssv.GetLastProcessedBlock()
	if param.Block+300 < processedBlock {
		ReturnErr(c, badRequestRes)
		return
	}

	alarmInfo, err := ms.store.GetAlarmByEoaOwner(param.Owner)
	if err != nil {
		monitorLog.Errorw("DeleteMonitorConfig: GetAlarmByEoaOwner", "err", err.Error())
		ReturnErr(c, serverErrRes)
		return
	}

	if alarmInfo == nil {
		ReturnErr(c, badRequestRes)
		return
	}

	sign := common.FromHex(param.Signature)
	msg := fmt.Sprintf(getMonitorConfigFormat, param.Block)
	addr, err := crypto.Ecrecover([]byte(msg), sign)
	if err != nil {
		ReturnErr(c, badRequestRes)
		return
	}
	if addr != param.Owner {
		ReturnErr(c, badRequestRes)
		return
	}

	if err := ms.store.DeleteAlarmByEoaOwner(param.Owner); err != nil {
		monitorLog.Errorw("DeleteMonitorConfig: DeleteAlarmByEoaOwner", "err", err.Error())
		ReturnErr(c, serverErrRes)
		return
	}

	ReturnOk(c, nil)
}

func (ms *MonitorSSV) GetClusterMonitorConfig(c *gin.Context) {
	owner := c.DefaultQuery("owner", "")
	if owner == "" {
		monitorLog.Warnw("GetClusterMonitorConfig", "owner", owner)
		ReturnErr(c, badRequestRes)
		return
	}

	blockStr := c.DefaultQuery("block", "")
	if blockStr == "" {
		monitorLog.Warnw("GetClusterMonitorConfig", "block", blockStr)
		ReturnErr(c, badRequestRes)
		return
	}

	block, err := strconv.ParseUint(blockStr, 10, 64)
	if err != nil {
		ReturnErr(c, badRequestRes)
		return
	}

	signature := c.DefaultQuery("signature", "")
	if signature == "" {
		monitorLog.Warnw("GetClusterMonitorConfig", "signature", signature)
		ReturnErr(c, badRequestRes)
		return
	}

	processedBlock := ms.ssv.GetLastProcessedBlock()
	monitorLog.Infow("GetClusterMonitor", "owner", owner, "processedBlock", processedBlock, "block", block, "signature", signature)

	if block+300 < processedBlock {
		ReturnErr(c, badRequestRes)
		return
	}

	sign := common.FromHex(signature)
	msg := fmt.Sprintf(getMonitorConfigFormat, block)
	addr, err := crypto.Ecrecover([]byte(msg), sign)
	if err != nil {
		ReturnErr(c, badRequestRes)
		return
	}
	if addr != owner {
		ReturnErr(c, badRequestRes)
		return
	}

	_, totalActiveClusterCount, err := ms.getOwnerClusterInfo(owner)
	if err != nil {
		monitorLog.Errorw("getOwnerClusterInfo", "err", err.Error())
		ReturnErr(c, serverErrRes)
		return
	}
	if totalActiveClusterCount == 0 {
		ReturnErr(c, badRequestRes)
		return
	}

	alarmInfo, err := ms.store.GetAlarmByEoaOwner(owner)
	if err != nil {
		monitorLog.Errorw("GetAlarmByEoaOwner", "err", err)
		ReturnErr(c, serverErrRes)
		return
	}

	key := crypto.GenerateEncryptKey([]byte(ms.password))
	encryptedData, err := hex.DecodeString(alarmInfo.AlarmChannel)
	if err != nil {
		monitorLog.Errorw("GetClusterMonitorConfig: DecodeString", "err", err)
		ReturnErr(c, serverErrRes)
		return
	}

	alarmChannel, err := crypto.DecryptData(encryptedData, key)
	if err != nil {
		monitorLog.Errorw("GetClusterMonitorConfig: DecryptData", "err", err)
		ReturnErr(c, serverErrRes)
		return
	}
	alarmChannelHash := crypto.Hash256(alarmChannel)
	if hex.EncodeToString(alarmChannelHash) != alarmInfo.AlarmChannelHash {
		monitorLog.Errorw("GetClusterMonitorConfig: alarmChannelHash does not match", "want", alarmInfo.AlarmChannelHash, "get", hex.EncodeToString(alarmChannelHash))
		ReturnErr(c, serverErrRes)
		return
	}

	var mc MonitorConfig
	mc.AlarmType = alarmInfo.AlarmType
	mc.AlarmChannel = string(alarmChannel)
	mc.ReportLiquidationThreshold = alarmInfo.ReportLiquidationThreshold / 7200
	mc.ReportOperatorFeeChange = alarmInfo.ReportOperatorFeeChange
	mc.ReportNetworkFeeChange = alarmInfo.ReportNetworkFeeChange
	mc.ReportProposeBlock = alarmInfo.ReportProposeBlock
	mc.ReportMissedBlock = alarmInfo.ReportMissedBlock
	mc.ReportBalanceDecrease = alarmInfo.ReportBalanceDecrease
	mc.ReportExitedButNotRemoved = alarmInfo.ReportExitedButNotRemoved
	mc.ReportWeekly = alarmInfo.ReportWeekly

	ReturnOk(c, gin.H{
		"monitorConfig": mc,
		"block":         processedBlock,
	})
}

var saveMonitorConfigFormat = "Signature required for cluster ownership. Block: %d\n%s"

func (ms *MonitorSSV) SaveClusterMonitorConfig(c *gin.Context) {
	type Request struct {
		MonitorConfig string `json:"monitorConfig"`
		Owner         string `json:"owner"`
		Signature     string `json:"signature"`
		Block         uint64 `json:"block"`
	}

	param := Request{}
	err := c.ShouldBind(&param)
	if err != nil {
		monitorLog.Warnw("SaveClusterMonitorConfig", "err", err)
		ReturnErr(c, badRequestRes)
		return
	}

	var monitorConfig MonitorConfig
	err = json.Unmarshal([]byte(param.MonitorConfig), &monitorConfig)
	if err != nil {
		monitorLog.Warnw("SaveClusterMonitorConfig", "err", err)
		ReturnErr(c, badRequestRes)
		return
	}

	processedBlock := ms.ssv.GetLastProcessedBlock()

	if param.Block+300 < processedBlock {
		ReturnErr(c, badRequestRes)
		return
	}

	sign := common.FromHex(param.Signature)
	msg := fmt.Sprintf(saveMonitorConfigFormat, param.Block, param.MonitorConfig)

	addr, err := crypto.Ecrecover([]byte(msg), sign)
	if err != nil {
		ReturnErr(c, badRequestRes)
		return
	}

	if addr != param.Owner {
		ReturnErr(c, badRequestRes)
		return
	}

	_, totalActiveClusterCount, err := ms.getOwnerClusterInfo(addr)
	if err != nil {
		monitorLog.Errorw("getOwnerClusterInfo", "err", err.Error())
		ReturnErr(c, serverErrRes)
		return
	}
	if totalActiveClusterCount == 0 {
		ReturnErr(c, badRequestRes)
		return
	}

	alarmChannelHash := crypto.Hash256([]byte(monitorConfig.AlarmChannel))
	key := crypto.GenerateEncryptKey([]byte(ms.password))
	encryptedData, err := crypto.EncryptData([]byte(monitorConfig.AlarmChannel), key)
	if err != nil {
		monitorLog.Warnw("SaveClusterMonitorConfig: EncryptData", "err", err)
		ReturnErr(c, serverErrRes)
		return
	}

	info := &store.AlarmInfo{
		EoaOwner:                   addr,
		AlarmType:                  monitorConfig.AlarmType,
		AlarmChannel:               hex.EncodeToString(encryptedData),
		AlarmChannelHash:           hex.EncodeToString(alarmChannelHash),
		ReportLiquidationThreshold: monitorConfig.ReportLiquidationThreshold * 7200,
		ReportOperatorFeeChange:    monitorConfig.ReportOperatorFeeChange,
		ReportNetworkFeeChange:     monitorConfig.ReportNetworkFeeChange,
		ReportProposeBlock:         monitorConfig.ReportProposeBlock,
		ReportMissedBlock:          monitorConfig.ReportMissedBlock,
		ReportBalanceDecrease:      monitorConfig.ReportBalanceDecrease,
		ReportExitedButNotRemoved:  monitorConfig.ReportExitedButNotRemoved,
		ReportWeekly:               monitorConfig.ReportWeekly,
	}
	err = ms.store.CreateOrUpdateAlarmInfo(info)

	if err != nil {
		monitorLog.Errorw("SaveClusterMonitorConfig", "err", err)
		ReturnErr(c, serverErrRes)
		return
	}

	monitorLog.Infow("SaveClusterMonitorConfig", "owner", param.Owner, "processedBlock", processedBlock, "block", param.Block, "signature", param.Signature, "monitorConfig", info)

	ReturnOk(c, gin.H{
		"monitorConfig": monitorConfig,
		"block":         processedBlock,
	})
}

func (ms *MonitorSSV) getOwnerClusterInfo(owner string) (uint64, uint64, error) {
	clusters, err := ms.store.GetAllClusterByEoaOwner(owner)
	if err != nil {
		monitorLog.Errorw("GetAllClusterByEoaOwner", "err", err.Error())
		return 0, 0, err
	}
	totalClusterCount := len(clusters)
	totalActiveClusterCount := 0
	for _, cluster := range clusters {
		if cluster.Active {
			totalActiveClusterCount++
		}
	}

	return uint64(totalClusterCount), uint64(totalActiveClusterCount), nil
}

func (ms *MonitorSSV) TestAlarm(c *gin.Context) {
	type Request struct {
		AlarmType    int    `json:"alarm_type"`
		AlarmChannel string `json:"alarm_channel"`
	}

	param := Request{}
	err := c.ShouldBind(&param)
	if err != nil {
		monitorLog.Warnw("TestAlarm", "err", err)
		ReturnErr(c, badRequestRes)
		return
	}

	alarm, err := alert.NewAlarm(param.AlarmType, param.AlarmChannel)
	if err != nil {
		monitorLog.Warnw("TestAlarm", "err", err)
		ReturnErr(c, newResponse(badRequestCode, err.Error()))
		return
	}

	err = alarm.Send(alert.TestAlarmMsg)
	if err != nil {
		monitorLog.Warnw("TestAlarm: Send", "err", err)
		ReturnErr(c, newResponse(badRequestCode, err.Error()))
		return
	}

	ReturnOk(c, nil)
}
