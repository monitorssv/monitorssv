package ssv

import (
	"github.com/ethereum/go-ethereum/common"
	logging "github.com/ipfs/go-log/v2"
	"github.com/monitorssv/monitorssv/alert"
	"github.com/monitorssv/monitorssv/config"
	"github.com/monitorssv/monitorssv/eth1/client"
	"github.com/monitorssv/monitorssv/store"
	"math/big"
	"sync"
	"testing"
)

func initSSV(t *testing.T) *SSV {
	_ = logging.SetLogLevel("*", "INFO")
	cfg, err := config.InitConfig("../../deploy/monitorssv/config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	eth1Client, err := client.NewEth1Client(cfg)
	if err != nil {
		t.Fatal(err)
	}

	db, err := store.NewStore(cfg)
	if err != nil {
		t.Fatal(err)
	}
	password := "test20240908"
	alarmDaemon, err := alert.NewAlarmDaemon(db, eth1Client, password)
	if err != nil {
		t.Fatal(err)
	}
	ssv, err := NewSSV(cfg, eth1Client, db, alarmDaemon)
	if err != nil {
		t.Fatal(err)
	}
	return ssv
}

func TestSSV(t *testing.T) {
	ssv := initSSV(t)
	wg := sync.WaitGroup{}
	wg.Add(3)
	ssv.Start()

	wg.Wait()
}

func TestCalcAllClusterLiquidation(t *testing.T) {
	ssv := initSSV(t)
	err := ssv.calcAllClusterLiquidation()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateOperatorName(t *testing.T) {
	ssv := initSSV(t)
	ssv.updateOperatorName()
}

func TestUpdateOperatorEarning(t *testing.T) {
	ssv := initSSV(t)
	ssv.updateOperatorEarning()
}

func TestUpdateClusterEoaOwner(t *testing.T) {
	ssv := initSSV(t)
	ssv.updateClusterEoaOwner()
}

func TestOperatorFeeUpdateCalcClusterLiquidation(t *testing.T) {
	ssv := initSSV(t)
	err := ssv.operatorFeeUpdateCalcClusterLiquidation(1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClusterOnChainBalance(t *testing.T) {
	ssv := initSSV(t)
	clusterInfos, err := ssv.store.GetAllClusters()
	if err != nil {
		t.Fatal(err)
	}
	for _, clusterInfo := range clusterInfos {
		if clusterInfo.ValidatorCount < 20 {
			continue
		}

		operatorIds, err := getOperatorIds(clusterInfo.OperatorIds)
		if err != nil {
			t.Fatal(err)
		}
		balance, isOk := big.NewInt(0).SetString(clusterInfo.Balance, 10)
		if !isOk {
			t.Fatal(err)
		}
		_, curBlock, _, onChainBalance, err := ssv.CalcLiquidation(Cluster{
			ClusterId:   clusterInfo.ClusterID,
			Owner:       common.HexToAddress(clusterInfo.Owner),
			OperatorIds: operatorIds,
			ClusterInfo: ISSVNetworkCoreCluster{
				ValidatorCount:  clusterInfo.ValidatorCount,
				NetworkFeeIndex: clusterInfo.NetworkFeeIndex,
				Index:           clusterInfo.Index,
				Active:          clusterInfo.Active,
				Balance:         balance,
			},
		})
		if err != nil {
			t.Log("====err====", err)
			continue
		}

		t.Log(onChainBalance)

		onChainBalance4 := store.CalcClusterOnChainBalance(curBlock, &clusterInfo)
		t.Log(onChainBalance4)
		t.Log("-----------------------------")
	}
}

func TestCalcLiquidation(t *testing.T) {
	ssv := initSSV(t)
	clusterInfo, err := ssv.store.GetClusterByClusterId("61bb7fdc5b3ccc9ae6573bdaf86fbc26e0681f42b1b513b56325f5f0a63f2b49")
	if err != nil {
		t.Fatal(err)
	}
	operatorIds, err := getOperatorIds(clusterInfo.OperatorIds)
	if err != nil {
		t.Fatal(err)
	}
	balance, isOk := big.NewInt(0).SetString(clusterInfo.Balance, 10)
	if !isOk {
		t.Fatal(err)
	}
	cluster := Cluster{
		ClusterId:   clusterInfo.ClusterID,
		Owner:       common.HexToAddress(clusterInfo.Owner),
		OperatorIds: operatorIds,
		ClusterInfo: ISSVNetworkCoreCluster{
			ValidatorCount:  clusterInfo.ValidatorCount,
			NetworkFeeIndex: clusterInfo.NetworkFeeIndex,
			Index:           clusterInfo.Index,
			Active:          clusterInfo.Active,
			Balance:         balance,
		},
	}
	t.Log(cluster)
	liquidationBlock, curBlock, burnFee, onChainBalance, err := ssv.CalcLiquidation(cluster)
	t.Log(liquidationBlock)
	t.Log(curBlock)
	t.Log(liquidationBlock - curBlock)
	t.Log((liquidationBlock - curBlock) / 7200)
	t.Log(burnFee)
	t.Log(onChainBalance)
}
