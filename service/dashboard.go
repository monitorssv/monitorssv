package service

import (
	"github.com/gin-gonic/gin"
	"github.com/monitorssv/monitorssv/store"
	"strconv"
	"strings"
)

type DashboardData struct {
	ActiveOperators        int64             `json:"activeOperators"`
	ActiveValidators       int64             `json:"activeValidators"`
	ActiveClusters         int64             `json:"activeClusters"`
	StakedETH              int64             `json:"stakedETH"`
	ProposedBlocks         int64             `json:"proposedBlocks"`
	NetworkFee             string            `json:"networkFee"`
	OperatorValidatorLimit int64             `json:"operatorValidatorLimit"`
	LiquidationThreshold   int64             `json:"liquidationThreshold"`
	MinimumCollateral      string            `json:"minimumCollateral"`
	Events                 []Event           `json:"events"`
	Blocks                 []Block           `json:"blocks"`
	Validators             []Validator       `json:"validators"`
	Charts                 []store.ChartData `json:"charts"`
}

func (ms *MonitorSSV) Dashboard(c *gin.Context) {
	networkInfo, err := ms.ssv.GetNetworkInfo()
	if err != nil {
		monitorLog.Errorw("Dashboard:GetNetworkInfo", "err", err)
		ReturnErr(c, serverErrRes)
		return
	}

	var proposedBlocks int64
	var activeClusterCount int64
	if ms.ssv.GetCfg().Network == "mainnet" {
		proposedBlocks, err = ms.store.GetTotalBlockCount()
		if err != nil {
			monitorLog.Errorw("Dashboard:GetTotalBlockCount", "err", err)
			ReturnErr(c, serverErrRes)
			return
		}
	} else {
		activeClusterCount, err = ms.store.GetActiveClusterCount()
		if err != nil {
			monitorLog.Errorw("Dashboard:GetActiveClusterCount", "err", err)
			ReturnErr(c, serverErrRes)
			return
		}
	}

	activeOperators, err := ms.store.GetActiveOperatorCount()
	if err != nil {
		monitorLog.Errorw("Dashboard:GetActiveOperatorCount", "err", err)
		ReturnErr(c, serverErrRes)
		return
	}

	activeValidators, err := ms.store.GetActiveValidatorCount()
	if err != nil {
		monitorLog.Errorw("Dashboard:GetActiveValidatorCount", "err", err)
		ReturnErr(c, serverErrRes)
		return
	}

	latestEvents, err := ms.store.GetLatestEvents()
	if err != nil {
		monitorLog.Errorw("Dashboard:GetLatestEvents", "err", err)
		ReturnErr(c, serverErrRes)
		return
	}

	var dashboardData DashboardData

	if ms.ssv.GetCfg().Network == "mainnet" {
		latestBlocks, err := ms.store.GetLatestBlocks()
		if err != nil {
			monitorLog.Errorw("Dashboard:GetLatestBlocks", "err", err)
			ReturnErr(c, serverErrRes)
			return
		}

		var blocks = make([]Block, 0)
		for _, block := range latestBlocks {
			blocks = append(blocks, Block{
				Proposer:    block.Proposer,
				Epoch:       block.Epoch,
				Slot:        block.Slot,
				BlockNumber: block.BlockNumber,
			})
		}
		dashboardData.Blocks = blocks
	} else {
		latestValidators, err := ms.store.GetLatestValidators()
		if err != nil {
			monitorLog.Errorw("Dashboard:GetLatestValidators", "err", err)
			ReturnErr(c, serverErrRes)
			return
		}

		var validators = make([]Validator, 0)
		for _, info := range latestValidators {
			operatorIds := strings.Split(info.OperatorIds, ",")
			var operators = make([]OperatorIntro, 0)
			for _, operatorId := range operatorIds {
				id, _ := strconv.Atoi(operatorId)
				operators = append(operators, ms.getOperatorIntro(uint64(id)))
			}
			validators = append(validators, Validator{
				PublicKey: info.PublicKey,
				Owner:     info.Owner,
				Operators: operators,
				ClusterId: info.ClusterID,
				Status:    info.Status,
				Online:    info.IsOnline,
			})
		}
		dashboardData.Validators = validators
	}

	chartData, err := ms.store.CalculateChartData()
	if err != nil {
		monitorLog.Errorw("Dashboard:CalculateChartData", "err", err)
		ReturnErr(c, serverErrRes)
		return
	}

	dashboardData.ActiveOperators = activeOperators
	dashboardData.ActiveValidators = activeValidators
	dashboardData.StakedETH = activeValidators * 32
	dashboardData.ActiveClusters = activeClusterCount
	dashboardData.ProposedBlocks = proposedBlocks
	dashboardData.NetworkFee = networkInfo.NetworkFee
	dashboardData.MinimumCollateral = networkInfo.MinimumLiquidationCollateral
	dashboardData.OperatorValidatorLimit = networkInfo.OperatorValidatorLimit
	dashboardData.LiquidationThreshold = networkInfo.LiquidationThresholdPeriod
	dashboardData.OperatorValidatorLimit = networkInfo.OperatorValidatorLimit
	dashboardData.Charts = chartData

	var events = make([]Event, 0)
	for _, event := range latestEvents {
		events = append(events, Event{
			BlockNumber: event.BlockNumber,
			TxHash:      event.TxHash,
			Action:      event.Action,
		})
	}
	dashboardData.Events = events

	ReturnOk(c, dashboardData)
}
