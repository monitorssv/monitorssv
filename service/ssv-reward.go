package service

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"strings"
)

//address account,uint256 cumulativeAmount,bytes32 expectedMerkleRoot,bytes32[] merkleProof

type SSVReward struct {
	Account     string   `json:"account"`
	Amount      string   `json:"cumulativeAmount"`
	MerkleRoot  string   `json:"expectedMerkleRoot"`
	MerkleProof []string `json:"merkleProof"`
}

func (ms *MonitorSSV) GetSSVReward(c *gin.Context) {
	account := c.DefaultQuery("account", "")
	if len(account) != ethAddrLength {
		monitorLog.Warnw("GetSSVReward", "account", account)
		ReturnErr(c, badRequestRes)
		return
	}

	owner := common.HexToAddress(account).String()

	var ssvReward SSVReward
	ssvRewardInfo := ms.store.GetSSVReward(owner)
	ssvReward.Account = owner
	ssvReward.Amount = ssvRewardInfo.Amount.String()
	ssvReward.MerkleRoot = ssvRewardInfo.MerkleRoot
	ssvReward.MerkleProof = strings.Split(ssvRewardInfo.Proofs, ",")

	ReturnOk(c, gin.H{
		"ssvRewardInfo": ssvReward,
	})
	return
}
