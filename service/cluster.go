package service

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/monitorssv/monitorssv/eth1/utils"
	"github.com/monitorssv/monitorssv/store"
	"math"
	"math/big"
	"strconv"
	"strings"
)

type Cluster struct {
	ID         string          `json:"id"`
	Owner      string          `json:"owner"`
	Operators  []OperatorIntro `json:"operators"`
	Validators uint32          `json:"validators"`
	Status     bool            `json:"status"`
}

type PosData struct {
	TotalProposedBlocks int64 `json:"totalProposedBlocks"`
	TotalMissedBlocks   int64 `json:"totalMissedBlocks"`
	TotalOfflineCount   int64 `json:"totalOfflineCount"`
	PendingRemovalCount int64 `json:"pendingRemovalCount"`
}

func (ms *MonitorSSV) GetPosData(c *gin.Context) {
	clusterId := c.DefaultQuery("clusterId", "")

	if len(clusterId) != clusterIdLength {
		monitorLog.Warnw("GetPosData", "clusterId", clusterId)
		ReturnErr(c, badRequestRes)
		return
	}

	monitorLog.Infow("GetPosData", "clusterId", clusterId)

	totalCount, totalMissedCount, err := ms.store.GetBlockInfoByClusterId(clusterId)
	if err != nil {
		monitorLog.Errorw("GetPosData: GetBlockByClusterId", "err", err.Error())
		ReturnErr(c, serverErrRes)
		return
	}

	totalOfflineCount, err := ms.store.GetClusterOfflineValidatorCount(clusterId)
	if err != nil {
		monitorLog.Errorw("GetPosData: GetClusterOfflineValidatorCount", "err", err.Error())
		ReturnErr(c, serverErrRes)
		return
	}

	exitedButNotRemovedCount, err := ms.store.GetActiveButExitedValidatorCount(clusterId)
	if err != nil {
		monitorLog.Errorw("GetPosData: GetActiveButExitedValidatorCount", "err", err.Error())
		ReturnErr(c, serverErrRes)
		return
	}

	var data = PosData{
		TotalProposedBlocks: totalCount,
		TotalMissedBlocks:   totalMissedCount,
		TotalOfflineCount:   totalOfflineCount,
		PendingRemovalCount: exitedButNotRemovedCount,
	}

	ReturnOk(c, gin.H{
		"posData": data,
	})
	return
}

func (ms *MonitorSSV) Get30DayLiquidationRankingClusters(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		monitorLog.Warnw("GetClusters", "page", page)
		ReturnErr(c, badRequestRes)
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		monitorLog.Warnw("GetClusters", "limit", limit)
		ReturnErr(c, badRequestRes)
		return
	}

	var clusterInfos []store.ClusterInfo
	var totalCount int64
	clusterInfos, totalCount, err = ms.store.Get30DayLiquidationRankingClusters(page, limit, ms.ssv.GetLastProcessedBlock())
	if err != nil {
		monitorLog.Errorw("GetClusters: GetClusters", "err", err.Error())
		ReturnErr(c, serverErrRes)
		return
	}

	var clusterDetails []ClusterDetails
	for _, clusterInfo := range clusterInfos {
		var clusterDetail ClusterDetails
		clusterDetail.ID = clusterInfo.ClusterID
		clusterDetail.Owner = clusterInfo.Owner
		clusterDetail.Active = clusterInfo.Active

		burnFeeInt := big.NewInt(0).SetUint64(clusterInfo.BurnFee)
		fee := big.NewInt(0).Mul(burnFeeInt, big.NewInt(2613400))
		burnFeeStr := utils.ToSSV(fee, "%.2f")
		clusterDetail.BurnFee = burnFeeStr

		curBlock := ms.ssv.GetLastProcessedBlock()

		onChainBalanceStr := store.CalcClusterOnChainBalance(curBlock, &clusterInfo)

		monitorLog.Infow("CalcOnChainBalance", "onChainBalance", onChainBalanceStr)

		clusterDetail.OnChainBalance = onChainBalanceStr

		operationalRunaway := uint64(0)
		if clusterInfo.Active && clusterInfo.LiquidationBlock > curBlock {
			operationalRunaway = clusterInfo.LiquidationBlock - curBlock
		}

		clusterDetail.OperationalRunaway = operationalRunaway
		clusterDetail.ValidatorCount = clusterInfo.ValidatorCount

		var operators = make([]OperatorIntro, 0)
		for _, operatorId := range strings.Split(clusterInfo.OperatorIds, ",") {
			id, _ := strconv.Atoi(operatorId)
			operators = append(operators, ms.getOperatorIntro(uint64(id)))
		}
		clusterDetail.Operators = operators
		clusterDetails = append(clusterDetails, clusterDetail)
	}

	ReturnOk(c, gin.H{
		"clusters":    clusterDetails,
		"totalItems":  totalCount,
		"totalPages":  int(math.Ceil(float64(totalCount) / float64(limit))),
		"currentPage": page,
	})
}

func (ms *MonitorSSV) GetClusters(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		monitorLog.Warnw("GetClusters", "page", page)
		ReturnErr(c, badRequestRes)
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		monitorLog.Warnw("GetClusters", "limit", limit)
		ReturnErr(c, badRequestRes)
		return
	}

	search := c.DefaultQuery("search", "")
	monitorLog.Infow("GetClusters", "page", page, "limit", limit, "search", search)

	var clusterInfos []store.ClusterInfo
	var totalCount int64
	if search == "" {
		clusterInfos, totalCount, err = ms.store.GetClusters(page, limit)
		if err != nil {
			monitorLog.Errorw("GetClusters: GetClusters", "err", err.Error())
			ReturnErr(c, serverErrRes)
			return
		}
	}

	// by owner
	if len(search) == ethAddrLength {
		owner := common.HexToAddress(search).String()
		clusterInfos, totalCount, err = ms.store.GetClusterByOwner(page, limit, owner)
		if err != nil {
			monitorLog.Errorw("GetClusters: GetClusterByOwner", "err", err.Error())
			ReturnErr(c, serverErrRes)
			return
		}
	}

	// by clusterId
	if len(search) == clusterIdLength {
		clusterId := search
		clusterInfo, err := ms.store.GetClusterByClusterId(clusterId)
		if err != nil {
			monitorLog.Errorw("GetClusters: GetClusterByClusterId", "err", err.Error())
			ReturnErr(c, serverErrRes)
			return
		}
		if clusterInfo != nil {
			clusterInfos = append(clusterInfos, *clusterInfo)
			totalCount = 1
		}
	}

	var clusters = make([]Cluster, 0)
	for _, info := range clusterInfos {
		operatorIds := strings.Split(info.OperatorIds, ",")
		var operators = make([]OperatorIntro, 0)
		for _, operatorId := range operatorIds {
			id, _ := strconv.Atoi(operatorId)
			operators = append(operators, ms.getOperatorIntro(uint64(id)))
		}

		clusters = append(clusters, Cluster{
			ID:         info.ClusterID,
			Owner:      info.Owner,
			Operators:  operators,
			Validators: info.ValidatorCount,
			Status:     info.Active,
		})
	}

	ReturnOk(c, gin.H{
		"clusters":    clusters,
		"totalItems":  totalCount,
		"totalPages":  int(math.Ceil(float64(totalCount) / float64(limit))),
		"currentPage": page,
	})
	return
}

type ClusterDetails struct {
	ID                  string          `json:"id"`
	Owner               string          `json:"owner"`
	FeeRecipientAddress string          `json:"feeRecipientAddress"`
	Active              bool            `json:"active"`
	OnChainBalance      string          `json:"onChainBalance"`
	BurnFee             string          `json:"burnFee"`
	OperationalRunaway  uint64          `json:"operationalRunaway"`
	ValidatorCount      uint32          `json:"validatorCount"`
	Operators           []OperatorIntro `json:"operators"`
}

func (ms *MonitorSSV) GetClusterDetails(c *gin.Context) {
	clusterId := c.DefaultQuery("clusterId", "")

	if len(clusterId) != clusterIdLength {
		monitorLog.Warnw("GetClusterDetails", "clusterId", clusterId)
		ReturnErr(c, badRequestRes)
		return
	}

	clusterInfo, err := ms.store.GetClusterByClusterId(clusterId)
	if err != nil {
		monitorLog.Errorw("GetClusterDetails: GetClusterByClusterId", "err", err.Error())
		ReturnErr(c, serverErrRes)
		return
	}

	var feeRecipientAddress string = clusterInfo.Owner
	feeAddress, err := ms.store.GetClusterFeeAddress(clusterInfo.Owner)
	if err != nil {
		monitorLog.Errorw("GetClusterDetails: GetClusterFeeAddress", "err", err.Error())
	} else if feeAddress.FeeAddress != "" {
		feeRecipientAddress = feeAddress.FeeAddress
	}

	var clusterDetails ClusterDetails
	clusterDetails.ID = clusterId
	clusterDetails.Owner = clusterInfo.Owner
	clusterDetails.Active = clusterInfo.Active
	clusterDetails.FeeRecipientAddress = feeRecipientAddress

	burnFeeInt := big.NewInt(0).SetUint64(clusterInfo.BurnFee)
	fee := big.NewInt(0).Mul(burnFeeInt, big.NewInt(2613400))
	burnFeeStr := utils.ToSSV(fee, "%.2f")
	clusterDetails.BurnFee = burnFeeStr

	curBlock := ms.ssv.GetLastProcessedBlock()

	onChainBalanceStr := store.CalcClusterOnChainBalance(curBlock, clusterInfo)

	monitorLog.Infow("CalcOnChainBalance", "onChainBalance", onChainBalanceStr)

	clusterDetails.OnChainBalance = onChainBalanceStr

	operationalRunaway := uint64(0)
	if clusterInfo.Active && clusterInfo.LiquidationBlock > curBlock {
		operationalRunaway = clusterInfo.LiquidationBlock - curBlock
	}

	clusterDetails.OperationalRunaway = operationalRunaway
	clusterDetails.ValidatorCount = clusterInfo.ValidatorCount

	var operators = make([]OperatorIntro, 0)
	for _, operatorId := range strings.Split(clusterInfo.OperatorIds, ",") {
		id, _ := strconv.Atoi(operatorId)
		operators = append(operators, ms.getOperatorIntro(uint64(id)))
	}
	clusterDetails.Operators = operators

	ReturnOk(c, gin.H{
		"clusterDetails": clusterDetails,
	})
	return
}

type ClusterReward struct {
	ID                  string `json:"id"`
	TotalProposedBlocks uint64 `json:"totalProposedBlocks"`
	TotalMissedBlocks   uint64 `json:"totalMissedBlocks"`
	TotalRewards        string `json:"totalRewards"`
	TotalPenalties      string `json:"totalPenalties"`
}
