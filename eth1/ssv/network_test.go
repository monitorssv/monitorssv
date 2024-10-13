package ssv

import (
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func TestGetNetworkInfo(t *testing.T) {
	ssv := initSSV(t)
	network, err := ssv.GetNetworkInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(network.NetworkFee)
	t.Log(network.LiquidationThresholdPeriod)
	t.Log(network.MinimumLiquidationCollateral)
	t.Log(network.OperatorValidatorLimit)

}

func TestGetSSVRewardMerkleRootOnChain(t *testing.T) {
	ssv := initSSV(t)
	root, err := GetSSVRewardMerkleRootOnChain(ssv.client.GetClient())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(root)
}

func TestGetSSVRewardCumulativeClaimed(t *testing.T) {
	ssv := initSSV(t)
	addr1 := common.HexToAddress("0x00b09f79228ef82d5925669ab94d6188df24e085")
	addr2 := common.HexToAddress("0x057f66b1e1308fa4259631e33ff202d244c8ad9c")
	amounts, err := GetSSVRewardCumulativeClaimed(ssv.client.GetClient(), []common.Address{addr1, addr2})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(amounts)
}
