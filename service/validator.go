package service

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/monitorssv/monitorssv/store"
	"math"
	"strconv"
	"strings"
)

type Validator struct {
	PublicKey string          `json:"publicKey"`
	Owner     string          `json:"owner"`
	Operators []OperatorIntro `json:"operators"`
	ClusterId string          `json:"clusterId"`
	Status    string          `json:"status"`
	Online    bool            `json:"online"`
}

func (ms *MonitorSSV) GetValidators(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		monitorLog.Warnw("GetValidators", "page", page)
		ReturnErr(c, badRequestRes)
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		monitorLog.Warnw("GetValidators", "limit", limit)
		ReturnErr(c, badRequestRes)
		return
	}

	search := c.DefaultQuery("search", "")
	monitorLog.Infow("GetValidators", "page", page, "limit", limit, "search", search)

	var validatorInfos []store.ValidatorInfo
	var totalCount int64
	if search == "" {
		validatorInfos, totalCount, err = ms.store.GetValidators(page, limit)
		if err != nil {
			monitorLog.Errorw("GetValidators: GetValidators", "err", err.Error())
			ReturnErr(c, serverErrRes)
			return
		}
	}

	// by owner
	if len(search) == ethAddrLength {
		owner := common.HexToAddress(search).String()
		validatorInfos, totalCount, err = ms.store.GetValidatorByOwner(page, limit, owner)
		if err != nil {
			monitorLog.Errorw("GetValidators: GetValidatorByOwner", "err", err.Error())
			ReturnErr(c, serverErrRes)
			return
		}
	}

	// by clusterId
	if len(search) == clusterIdLength {
		clusterId := search
		validatorInfos, totalCount, err = ms.store.GetValidatorByClusterId(page, limit, clusterId)
		if err != nil {
			monitorLog.Errorw("GetValidators: GetValidatorByOwner", "err", err.Error())
			ReturnErr(c, serverErrRes)
			return
		}
	}

	// by public key
	if len(search) == pubKeyLength || len(search) == pubKeyLength+2 {
		pubKey := search
		if search[:2] == "0x" {
			pubKey = pubKey[2:]
		}
		validatorInfo, err := ms.store.GetValidatorByPublicKey(pubKey)
		if err != nil {
			monitorLog.Errorw("GetValidators: GetValidatorByPublicKey", "err", err.Error())
			ReturnErr(c, serverErrRes)
			return
		}

		if validatorInfo != nil {
			validatorInfos = append(validatorInfos, *validatorInfo)
		}
	}

	var validators = make([]Validator, 0)
	for _, info := range validatorInfos {
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

	ReturnOk(c, gin.H{
		"validators":  validators,
		"totalItems":  totalCount,
		"totalPages":  int(math.Ceil(float64(totalCount) / float64(limit))),
		"currentPage": page,
	})
	return
}
