package eth2

import (
	"context"
	"github.com/monitorssv/monitorssv/alert"
)

import (
	logging "github.com/ipfs/go-log/v2"
	"github.com/monitorssv/monitorssv/config"
	"github.com/monitorssv/monitorssv/eth2/client"
	"github.com/monitorssv/monitorssv/store"
	"sync/atomic"
	"time"
)

var log = logging.Logger("beacon")

type BeaconMonitor struct {
	cfg *config.Config

	client *client.Client
	store  *store.Store

	lastProcessedSlot uint64
	isSynced          *atomic.Bool

	lastValidatorMonitorEpoch uint64
	validatorBalanceHistory   map[uint64][3]Balance

	validatorProposeBlockAlarmChan chan<- alert.ValidatorProposeBlockNotify
	validatorMissedBlockAlarmChan  chan<- alert.ValidatorMissedBlockNotify
	validatorBalanceDeltaAlarmChan chan<- alert.ValidatorBalanceDeltaNotify
	validatorSlashAlarmChan        chan<- alert.ValidatorSlashNotify

	close chan struct{}
}

func NewBeaconMonitor(cfg *config.Config, client *client.Client, store *store.Store, alarm *alert.AlarmDaemon) (*BeaconMonitor, error) {
	// ssv deploy block: 17507487
	lastProcessedSlot := uint64(6689770)
	if _, slot, err := store.GetScanPoint(); err == nil && slot != 0 {
		lastProcessedSlot = slot
	}

	bm := BeaconMonitor{
		cfg:    cfg,
		client: client,
		store:  store,

		lastProcessedSlot: lastProcessedSlot,
		isSynced:          new(atomic.Bool),

		lastValidatorMonitorEpoch: 0,
		validatorBalanceHistory:   make(map[uint64][3]Balance),

		validatorProposeBlockAlarmChan: alarm.ValidatorProposeBlockChan(),
		validatorMissedBlockAlarmChan:  alarm.ValidatorMissedBlockChan(),
		validatorBalanceDeltaAlarmChan: alarm.ValidatorBalanceDeltaChan(),
		validatorSlashAlarmChan:        alarm.ValidatorSlashNotifyChan(),

		close: make(chan struct{}),
	}

	bm.isSynced.Store(false)
	return &bm, nil
}

func (bm *BeaconMonitor) Start() {
	if bm.cfg.Network != "mainnet" {
		return
	}

	if bm.cfg.Dev {
		log.Info("Beacon monitor does not run in dev mode")
		return
	}

	log.Info("Beacon monitor is running")

	go bm.ScanBeaconBlockLoop()
	go bm.ValidatorMonitorLoop()
}

func (bm *BeaconMonitor) Stop() {
	close(bm.close)
}

func (bm *BeaconMonitor) GetLastProcessedSlot() uint64 {
	return bm.lastProcessedSlot
}

func (bm *BeaconMonitor) GetLastValidatorMonitorEpoch() uint64 {
	return bm.lastValidatorMonitorEpoch
}

func (bm *BeaconMonitor) ScanBeaconBlockLoop() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-bm.close:
			return
		case <-ticker.C:
			finalizedEpoch, err := bm.client.GetFinalizedEpoch()
			if err != nil {
				log.Warnw("failed to fetch finalized epoch", "err", err)
				continue
			}

			finalizedSlot := finalizedEpoch * 32

			log.Infow("ScanBeaconBlockLoop", "lastProcessedSlot", bm.lastProcessedSlot, "finalizedEpoch", finalizedEpoch)

			lastProcessedSlot, err := bm.ScanProposer(bm.lastProcessedSlot, finalizedEpoch)
			if err != nil {
				log.Warnw("failed to scan beacon block", "err", err)
			}

			if lastProcessedSlot > bm.lastProcessedSlot {
				bm.lastProcessedSlot = lastProcessedSlot
				err = bm.store.UpdateScanEth2Slot(lastProcessedSlot)
				if err != nil {
					log.Errorw("failed to update beacon block", "err", err)
				}
			}

			if !bm.isSynced.Load() {
				if lastProcessedSlot >= finalizedSlot {
					log.Infow("ScanBeaconBlockLoop: Sync completed", "lastProcessedSlot", lastProcessedSlot, "finalizedSlot", finalizedSlot)
					bm.isSynced.Store(true)
					ticker = time.NewTicker(13 * time.Minute)
				}
			}
		}
	}
}

func (bm *BeaconMonitor) ScanProposer(startSlot uint64, finalizedEpoch uint64) (uint64, error) {
	startEpoch := slotToEpoch(startSlot)

	endEpoch := startEpoch + 200
	if endEpoch > finalizedEpoch {
		endEpoch = finalizedEpoch
	}

	if endEpoch <= startEpoch {
		return startSlot, nil
	}

	log.Infow("ScanProposer", "startEpoch", startEpoch, "endEpoch", endEpoch, "finalizedEpoch", finalizedEpoch)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fetchBlocks, fetchErrors := bm.fetchBeaconBlocks(ctx, startEpoch, endEpoch, 50)

	lastProcessedSlot, err := bm.handleBlocks(fetchBlocks)
	if err != nil {
		return lastProcessedSlot, err
	}

	if err := <-fetchErrors; err != nil {
		return lastProcessedSlot, err
	}

	return (endEpoch + 1) * 32, nil
}

func (bm *BeaconMonitor) ValidatorMonitorLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	var lastSlot uint64
	for {
		select {
		case <-bm.close:
			return
		case <-ticker.C:
			if !bm.isSynced.Load() {
				continue
			}

			slot, err := bm.client.GetLatestSlot()
			if err != nil {
				log.Errorw("GetLatestSlot", "err", err)
				continue
			}

			curEpoch := slotToEpoch(slot)
			if curEpoch-1 <= bm.lastValidatorMonitorEpoch {
				continue
			}
			bm.lastValidatorMonitorEpoch = curEpoch - 1

			log.Infow("ValidatorMonitorLoop", "epoch", bm.lastValidatorMonitorEpoch, "curEpoch", curEpoch, "slot", slot)

			if slot != lastSlot {
				nextEpochStartSlot := (curEpoch + 1) * 32
				timeToNextEpoch := time.Duration(nextEpochStartSlot-slot) * 12 * time.Second
				log.Infow("ValidatorMonitorLoop: next ticker", "reset", timeToNextEpoch.String())
				ticker.Reset(timeToNextEpoch)
				lastSlot = slot
			}

			err = bm.validatorMonitor(bm.lastValidatorMonitorEpoch)
			if err != nil {
				log.Errorw("validatorMonitor", "err", err)
				continue
			}
		}
	}
}
