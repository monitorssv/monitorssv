package service

import (
	"github.com/gin-gonic/gin"
	logging "github.com/ipfs/go-log/v2"
	"github.com/monitorssv/monitorssv/alert"
	"github.com/monitorssv/monitorssv/eth1/ssv"
	"github.com/monitorssv/monitorssv/eth2"
	"github.com/monitorssv/monitorssv/store"
	"time"
)

var monitorLog = logging.Logger("monitor-service")

type DashboardCache struct {
	lastTime time.Time
	data     *DashboardData
}

type MonitorSSV struct {
	store         *store.Store
	ssv           *ssv.SSV
	beaconMonitor *eth2.BeaconMonitor
	alarm         *alert.AlarmDaemon
	password      string
	close         chan struct{}
}

func NewMonitorSSV(store *store.Store, ssv *ssv.SSV, beaconMonitor *eth2.BeaconMonitor, alarm *alert.AlarmDaemon, password string) (*MonitorSSV, error) {
	ms := &MonitorSSV{
		store:         store,
		ssv:           ssv,
		beaconMonitor: beaconMonitor,
		alarm:         alarm,
		password:      password,
		close:         make(chan struct{}),
	}

	ms.ssv.Start()
	ms.beaconMonitor.Start()
	ms.alarm.Start()

	return ms, nil
}

func (ms *MonitorSSV) Stop() {
	close(ms.close)
	ms.ssv.Stop()
	ms.beaconMonitor.Stop()
}

type Status struct {
	ELLastMonitoringBlock          uint64 `json:"el_last_monitoring_block"`
	CLLastProposalMonitoringEpoch  uint64 `json:"cl_last_proposal_monitoring_epoch"`
	CLLastValidatorMonitoringEpoch uint64 `json:"cl_last_validator_monitoring_epoch"`
}

func (ms *MonitorSSV) Status(c *gin.Context) {
	var status Status
	status.ELLastMonitoringBlock = ms.ssv.GetLastProcessedBlock()
	status.CLLastProposalMonitoringEpoch = ms.beaconMonitor.GetLastProcessedSlot()/32 - 1
	status.CLLastValidatorMonitoringEpoch = ms.beaconMonitor.GetLastValidatorMonitorEpoch()

	ReturnOk(c, status)
}
