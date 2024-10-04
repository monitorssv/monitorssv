package service

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/monitorssv/monitorssv/eth1/utils"
	"github.com/monitorssv/monitorssv/store"
	"math"
	"math/big"
	"strconv"
	"strings"
)

type OperatorIntro struct {
	Name string `json:"name"`
	ID   uint64 `json:"id"`
}

type Operator struct {
	ID                 uint64   `json:"id"`
	Name               string   `json:"name"`
	Owner              string   `json:"owner"`
	Validators         uint32   `json:"validators"`
	OperatorFee        string   `json:"operatorFee"`
	OperatorEarnings   string   `json:"operatorEarnings"`
	Privacy            bool     `json:"privacy"`
	Removed            bool     `json:"removed"`
	WhitelistedAddress []string `json:"whitelistedAddress"`
}

func (ms *MonitorSSV) getOperatorIntro(id uint64) OperatorIntro {
	var operator = OperatorIntro{
		ID: id,
	}
	operatorInfo, err := ms.store.GetOperatorByOperatorId(id)
	if err != nil {
		monitorLog.Warnw("getOperatorIntro", "err", err.Error())
	}

	if operatorInfo != nil && operatorInfo.OperatorName != "" {
		operator.Name = operatorInfo.OperatorName
	} else {
		operator.Name = fmt.Sprintf("Operator %d", id)
	}

	return operator
}

func (ms *MonitorSSV) GetOperators(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		monitorLog.Warnw("GetOperators", "page", page)
		ReturnErr(c, badRequestRes)
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		monitorLog.Warnw("GetOperators", "limit", limit)
		ReturnErr(c, badRequestRes)
		return
	}

	search := c.DefaultQuery("search", "")
	monitorLog.Infow("GetOperators", "page", page, "limit", limit, "search", search)

	var operatorInfos []store.OperatorInfo
	var totalCount int64
	if search == "" {
		operatorInfos, totalCount, err = ms.store.GetOperators(page, limit)
		if err != nil {
			monitorLog.Errorw("GetOperators: GetOperators", "err", err.Error())
			ReturnErr(c, serverErrRes)
			return
		}
	}

	// by owner
	if len(search) == ethAddrLength {
		monitorLog.Infow("GetOperators", "type", "byOwner", "search", search)
		owner := common.HexToAddress(search).String()
		operatorInfos, totalCount, err = ms.store.GetOperatorByOwner(page, limit, owner)
		if err != nil {
			monitorLog.Errorw("GetOperators: GetOperatorByOwner", "err", err.Error())
			ReturnErr(c, serverErrRes)
			return
		}
	}

	// by id
	if id, err := strconv.Atoi(search); err == nil {
		monitorLog.Infow("GetOperators", "type", "byId", "search", search)
		operatorInfo, err := ms.store.GetOperatorByOperatorId(uint64(id))
		if err != nil {
			monitorLog.Errorw("GetOperators: GetOperatorByOperatorId", "err", err.Error())
			ReturnErr(c, serverErrRes)
			return
		}
		if operatorInfo != nil {
			operatorInfos = append(operatorInfos, *operatorInfo)
			totalCount = 1
		}
	}

	if len(operatorInfos) == 0 {
		monitorLog.Infow("GetOperators", "type", "byName", "search", search)
		// by name
		operatorName := search
		operatorInfos, totalCount, err = ms.store.GetOperatorByOperatorName(page, limit, operatorName)
		if err != nil {
			monitorLog.Errorw("GetOperators: GetOperatorByOperatorName", "err", err.Error())
			ReturnErr(c, serverErrRes)
			return
		}
	}

	var operators = make([]Operator, 0)
	for _, info := range operatorInfos {
		operatorFee := info.OperatorFee
		if operatorFee != "0" {
			fee, isOk := big.NewInt(0).SetString(info.OperatorFee, 10)
			if isOk {
				fee = big.NewInt(0).Mul(fee, big.NewInt(2613400))
				operatorFee = utils.ToSSV(fee, "%.2f")
			} else {
				monitorLog.Warnw("GetOperators: Failed to parse operatorFee", "operatorFee", info.OperatorFee)
			}
		}

		privacy := info.PrivacyStatus
		if !privacy {
			if len(info.WhitelistedAddress) != 0 || len(info.WhitelistingContract) != 0 {
				privacy = true
			}
		}

		operatorName := info.OperatorName
		if operatorName == "" {
			operatorName = fmt.Sprintf("Operator %d", info.OperatorId)
		}

		operators = append(operators, Operator{
			ID:                 info.OperatorId,
			Name:               operatorName,
			Owner:              info.Owner,
			Validators:         info.ValidatorCount,
			OperatorFee:        operatorFee,
			OperatorEarnings:   info.OperatorEarnings,
			Privacy:            privacy,
			Removed:            info.RemoveBlock != 0,
			WhitelistedAddress: strings.Split(info.WhitelistedAddress, ","),
		})
	}

	ReturnOk(c, gin.H{
		"operators":   operators,
		"totalItems":  totalCount,
		"totalPages":  int(math.Ceil(float64(totalCount) / float64(limit))),
		"currentPage": page,
	})
	return
}
