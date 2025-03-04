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
	cfg, err := config.InitConfig("../deploy/monitorssv/config.yaml")
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
func TestSimulatedLiquidationAlarm(t *testing.T) {
	alarmDaemon := initAlarm(t)
	alarmDaemon.simulatedLiquidationAlarm()
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
		ClusterId: "5ef22a20d6c57cca7eca1f10853fd0d0c1c607ca5862f6fb1f078bcc9d90ba96",
		Index:     generateIndexs(),
	})
}

func TestValidatorSlashAlarm(t *testing.T) {
	alarmDaemon := initAlarm(t)
	alarmDaemon.validatorSlashAlarm(ValidatorSlashNotify{
		Epoch:     313246,
		ClusterId: "5ef22a20d6c57cca7eca1f10853fd0d0c1c607ca5862f6fb1f078bcc9d90ba96",
		Index:     generateIndexs(),
	})
}

func generateIndexs() []uint64 {
	start := uint64(249366)
	count := 210

	result := make([]uint64, count)
	for i := 0; i < count; i++ {
		result[i] = start + uint64(i)
	}
	return result
}

func TestChunkSlice(t *testing.T) {
	testSlice := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21}
	if len(chunkSlice(testSlice, 5)) != 5 {
		t.Fatal("chunk slice error")
	}
	t.Log(chunkSlice(testSlice, 5))

	if len(chunkSlice(testSlice, 2)) != 11 {
		t.Fatal("chunk slice error")
	}
	t.Log(chunkSlice(testSlice, 2))

	if len(chunkSlice(testSlice, 15)) != 2 {
		t.Fatal("chunk slice error")
	}
	t.Log(chunkSlice(testSlice, 15))

	if len(chunkSlice(testSlice, 10)) != 3 {
		t.Fatal("chunk slice error")
	}
	t.Log(chunkSlice(testSlice, 10))

	if len(chunkSlice(testSlice, 11)) != 2 {
		t.Fatal("chunk slice error")
	}
	t.Log(chunkSlice(testSlice, 11))

	if len(chunkSlice(testSlice, 20)) != 2 {
		t.Fatal("chunk slice error")
	}
	t.Log(chunkSlice(testSlice, 20))

	if len(chunkSlice(testSlice, 21)) != 1 {
		t.Fatal("chunk slice error")
	}
	t.Log(chunkSlice(testSlice, 21))

	if len(chunkSlice(testSlice, 0)) != 1 {
		t.Fatal("chunk slice error")
	}
	t.Log(chunkSlice(testSlice, 0))
}
