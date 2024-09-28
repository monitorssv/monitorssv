package alert

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/monitorssv/monitorssv/config"
	"github.com/monitorssv/monitorssv/eth1/client"
	"github.com/monitorssv/monitorssv/store"
	"math/big"
	"testing"
)

func initAlarm(t *testing.T) *AlarmDaemon {
	_ = logging.SetLogLevel("*", "INFO")
	cfg, err := config.InitConfig("../config/config.yaml")
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
	alarmDaemon, err := NewAlarmDaemon(db, eth1Client, password)
	if err != nil {
		t.Fatal(err)
	}

	return alarmDaemon
}

func TestLiquidationAlarm(t *testing.T) {
	alarmDaemon := initAlarm(t)
	alarmDaemon.liquidationAlarm()
}
func TestValidatorExitedButNotRemovedAlarm(t *testing.T) {
	alarmDaemon := initAlarm(t)
	alarmDaemon.validatorExitedButNotRemovedAlarm()
}
func TestWeeklyReport(t *testing.T) {
	alarmDaemon := initAlarm(t)
	alarmDaemon.weeklyReport()
}
func TestNetworkFeeChangeAlarm(t *testing.T) {
	alarmDaemon := initAlarm(t)
	alarmDaemon.networkFeeChangeAlarm(NetworkFeeChangeNotify{
		Block:         20814505,
		OldNetworkFee: big.NewInt(264650000000),
		NewNetworkFee: big.NewInt(364650000000),
	})
}
func TestOperatorFeeChangeAlarm(t *testing.T) {
	alarmDaemon := initAlarm(t)
	alarmDaemon.operatorFeeChangeAlarm(OperatorFeeChangeNotify{
		Block:       20814505,
		OperatorFee: big.NewInt(264650000000),
		OperatorId:  407,
	})
}
func TestProposeBlockAlarm(t *testing.T) {
	alarmDaemon := initAlarm(t)
	alarmDaemon.proposeBlockAlarm(ValidatorProposeBlockNotify{
		Epoch:     313246,
		Slot:      10023888,
		ClusterId: "df4e5f2a04ba6ea16eb721b577eac0edfb209389424d8d8217f81151ac80ac20",
		Index:     70812,
	})
}
func TestMissedBlockAlarm(t *testing.T) {
	alarmDaemon := initAlarm(t)
	alarmDaemon.missedBlockAlarm(ValidatorMissedBlockNotify{
		Epoch:     313246,
		Slot:      10023888,
		ClusterId: "df4e5f2a04ba6ea16eb721b577eac0edfb209389424d8d8217f81151ac80ac20",
		Index:     70812,
	})
}
func TestValidatorBalanceDeltaAlarm(t *testing.T) {
	alarmDaemon := initAlarm(t)
	alarmDaemon.validatorBalanceDeltaAlarm(ValidatorBalanceDeltaNotify{
		Epoch:     313246,
		ClusterId: "df4e5f2a04ba6ea16eb721b577eac0edfb209389424d8d8217f81151ac80ac20",
		Index:     []uint64{70812, 70813},
	})
}

func TestValidatorSlashAlarm(t *testing.T) {
	alarmDaemon := initAlarm(t)
	alarmDaemon.validatorSlashAlarm(ValidatorSlashNotify{
		Epoch:     313246,
		ClusterId: "df4e5f2a04ba6ea16eb721b577eac0edfb209389424d8d8217f81151ac80ac20",
		Index:     []uint64{70812, 70813},
	})
}
