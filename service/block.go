package service

import (
	"github.com/gin-gonic/gin"
	"math"
	"strconv"
)

type Block struct {
	Proposer    uint64 `json:"proposer"`
	Epoch       uint64 `json:"epoch"`
	Slot        uint64 `json:"slot"`
	BlockNumber uint64 `json:"blockNumber"`
}

func (ms *MonitorSSV) GetBlocks(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		monitorLog.Warnw("GetBlocks", "page", page)
		ReturnErr(c, badRequestRes)
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		monitorLog.Warnw("GetBlocks", "limit", limit)
		ReturnErr(c, badRequestRes)
		return
	}

	clusterId := c.DefaultQuery("clusterId", "")

	if len(clusterId) != clusterIdLength {
		monitorLog.Warnw("GetBlocks", "clusterId", clusterId)
		ReturnErr(c, badRequestRes)
		return
	}

	monitorLog.Infow("GetBlocks", "page", page, "limit", limit, "clusterId", clusterId)

	blockInfos, totalCount, err := ms.store.GetBlockByClusterId(page, limit, clusterId)
	if err != nil {
		monitorLog.Errorw("GetBlocks: GetBlockByClusterId", "err", err.Error())
		ReturnErr(c, serverErrRes)
		return
	}

	var blocks = make([]Block, 0)
	for _, blockInfo := range blockInfos {
		blocks = append(blocks, Block{
			Proposer:    blockInfo.Proposer,
			Epoch:       blockInfo.Epoch,
			Slot:        blockInfo.Slot,
			BlockNumber: blockInfo.BlockNumber,
		})
	}

	ReturnOk(c, gin.H{
		"blocks":      blocks,
		"totalItems":  totalCount,
		"totalPages":  int(math.Ceil(float64(totalCount) / float64(limit))),
		"currentPage": page,
	})
	return
}
