package service

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/monitorssv/monitorssv/store"
	"math"
	"strconv"
)

type Event struct {
	BlockNumber uint64 `json:"block"`
	TxHash      string `json:"transactionHash"`
	Action      string `json:"action"`
}

func (ms *MonitorSSV) GetEvents(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		monitorLog.Warnw("GetEvents", "page", page)
		ReturnErr(c, badRequestRes)
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		monitorLog.Warnw("GetEvents", "limit", limit)
		ReturnErr(c, badRequestRes)
		return
	}

	search := c.DefaultQuery("search", "")
	monitorLog.Infow("GetEvents", "page", page, "limit", limit, "search", search)

	var eventInfos = make([]store.EventInfo, 0)
	var totalCount int64

	if search == "" {
		monitorLog.Infow("GetEvents", "search is empty")
		totalCount = 0
	}

	if len(search) == ethAddrLength {
		owner := common.HexToAddress(search).String()
		eventInfos, totalCount, err = ms.store.GetEventByAccount(page, limit, owner)
		if err != nil {
			monitorLog.Errorw("GetEvents: GetEventByAccount", "err", err.Error())
			ReturnErr(c, serverErrRes)
			return
		}
	}

	if len(search) == clusterIdLength {
		clusterId := search
		eventInfos, totalCount, err = ms.store.GetEventByClusterId(page, limit, clusterId)
		if err != nil {
			monitorLog.Errorw("GetEvents: GetEventByClusterId", "err", err.Error())
			ReturnErr(c, serverErrRes)
			return
		}
	}

	var events = make([]Event, 0)
	for _, eventInfo := range eventInfos {
		events = append(events, Event{
			BlockNumber: eventInfo.BlockNumber,
			TxHash:      eventInfo.TxHash,
			Action:      eventInfo.Action,
		})
	}

	ReturnOk(c, gin.H{
		"history":     events,
		"totalItems":  totalCount,
		"totalPages":  int(math.Ceil(float64(totalCount) / float64(limit))),
		"currentPage": page,
	})
	return
}
