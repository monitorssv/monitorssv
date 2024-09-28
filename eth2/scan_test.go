package eth2

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/monitorssv/monitorssv/alert"
	"github.com/monitorssv/monitorssv/config"
	eth1client "github.com/monitorssv/monitorssv/eth1/client"
	"github.com/monitorssv/monitorssv/eth2/client"
	"github.com/monitorssv/monitorssv/store"
	"sync"
	"testing"
)

func initBeaconMonitor(t *testing.T) *BeaconMonitor {
	_ = logging.SetLogLevel("*", "INFO")
	cfg, err := config.InitConfig("../config/config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	eth1Client, err := eth1client.NewEth1Client(cfg)
	if err != nil {
		t.Fatal(err)
	}

	beaconClient := client.NewClient(cfg.Eth2Rpc)

	db, err := store.NewStore(cfg)
	if err != nil {
		t.Fatal(err)
	}
	password := "test20240908"
	alarmDaemon, err := alert.NewAlarmDaemon(db, eth1Client, password)
	if err != nil {
		t.Fatal(err)
	}
	bm, err := NewBeaconMonitor(cfg, beaconClient, db, alarmDaemon)
	if err != nil {
		t.Fatal(err)
	}

	return bm
}

func TestScan(t *testing.T) {
	bm := initBeaconMonitor(t)

	wg := sync.WaitGroup{}
	wg.Add(1)

	bm.Start()

	wg.Wait()
}
