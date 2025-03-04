package service

import (
	"github.com/gin-gonic/gin"
	logging "github.com/ipfs/go-log/v2"
	"github.com/monitorssv/monitorssv/alert"
	"github.com/monitorssv/monitorssv/eth1/ssv"
	"github.com/monitorssv/monitorssv/eth1/utils"
	"github.com/monitorssv/monitorssv/eth2"
	"github.com/monitorssv/monitorssv/store"
	"math/big"
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

type NetworkFee struct {
	CurrentFee  string `json:"current"`
	UpcomingFee string `json:"upcoming"`
}

func (ms *MonitorSSV) GetNetworkFees(c *gin.Context) {
	chainNetworkInfo, err := ms.ssv.GetNetworkInfo()
	if err != nil {
		monitorLog.Errorw("GetNetworkFees: chain.GetNetworkInfo", "err", err.Error())
		ReturnErr(c, serverErrRes)
		return
	}

	storeNetworkInfo, err := ms.store.GetNetworkInfo()
	if err != nil {
		monitorLog.Errorw("GetNetworkFees: store.GetNetworkInfo", "err", err.Error())
		ReturnErr(c, serverErrRes)
		return
	}

	var upcomingNetworkFee = chainNetworkInfo.NetworkFee
	if storeNetworkInfo != nil {
		storeNetworkFee := big.NewInt(0).SetUint64(storeNetworkInfo.UpcomingNetworkFee)
		fee := big.NewInt(0).Mul(storeNetworkFee, big.NewInt(2613400))
		upcomingNetworkFee = utils.ToSSV(fee, "%.2f")
	}

	networkFee := NetworkFee{
		CurrentFee:  chainNetworkInfo.NetworkFee,
		UpcomingFee: upcomingNetworkFee,
	}

	ReturnOk(c, networkFee)
}
