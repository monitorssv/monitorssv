package ssv

import "testing"

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
	root, err := GetSSVRewardMerkleRootOnChain(ssv.cfg.Network, ssv.client.GetClient())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(root)
}
