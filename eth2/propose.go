package eth2

import (
	"context"
	"errors"
	"fmt"
	"github.com/monitorssv/monitorssv/alert"
	"github.com/monitorssv/monitorssv/eth1/utils"
	"github.com/monitorssv/monitorssv/eth2/client"
	"github.com/monitorssv/monitorssv/store"
	"sync"
)

type BlockInfo struct {
	BlockNumber uint64
	Epoch       uint64
	Slot        uint64
	PubKey      string
	Index       uint64
	IsMissed    bool
}

func (bm *BeaconMonitor) fetchBeaconBlocks(ctx context.Context, startEpoch, endEpoch uint64, parallel int) (<-chan BlockInfo, <-chan error) {
	if endEpoch-startEpoch < uint64(parallel) {
		parallel = int(endEpoch - startEpoch)
	}

	fetchBlocks := make(chan BlockInfo, 500)
	fetchError := make(chan error, 1)
	jobs := make(chan uint64, parallel)

	log.Infow("fetchBeaconBlocks", "startEpoch", startEpoch, "endEpoch", endEpoch, "parallel", parallel)

	go func() {
		defer close(fetchBlocks)
		defer close(fetchError)

		errorOccurred := make(chan struct{})
		var wg sync.WaitGroup

		for i := 0; i < parallel; i++ {
			wg.Add(1)
			go func(j int) {
				defer wg.Done()
				for fromEpoch := range jobs {
					select {
					case <-ctx.Done():
						return
					default:
						log.Infow("fetchBeaconBlocks", "job", j, "epoch", fromEpoch)
						err := bm.processEpoch(ctx, fromEpoch, fetchBlocks)
						if err != nil {
							select {
							case fetchError <- err:
								close(errorOccurred)
							default:
							}
							return
						}
					}
				}
			}(i)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(jobs)

			for fromEpoch := startEpoch; fromEpoch <= endEpoch; fromEpoch++ {
				select {
				case <-ctx.Done():
					fetchError <- ctx.Err()
					return
				case <-bm.close:
					fetchError <- fmt.Errorf("close")
					return
				case <-errorOccurred:
					log.Warn("abnormal stop")
					return
				case jobs <- fromEpoch:
					log.Infow("fetchBeaconBlocks", "fromEpoch", fromEpoch, "toEpoch", endEpoch, "progress", fmt.Sprintf("%.2f%%", float64(fromEpoch-startEpoch+1)/float64(endEpoch-startEpoch+1)*100))
				}
			}
		}()

		wg.Wait()
	}()

	return fetchBlocks, fetchError
}

func (bm *BeaconMonitor) processEpoch(ctx context.Context, fromEpoch uint64, fetchBlocks chan<- BlockInfo) error {
	proposers, err := bm.client.GetEpochProposer(fromEpoch)
	if err != nil {
		return err
	}

	for _, proposer := range proposers.Data {
		pubKey := removePubKeyPrefix(proposer.Pubkey)
		validatorInfo, err := bm.store.GetValidatorByPublicKey(pubKey)
		if err != nil {
			log.Warnw("GetValidatorByPublicKey", "pubKey", pubKey, "err", err)
		}
		if validatorInfo == nil {
			log.Infow("GetValidatorByPublicKey: validator does not belong to ssv", "pubkey", pubKey, "err", err)
			continue
		}

		blockNumber, err := bm.getEth1ExBlock(uint64(proposer.Slot))
		if err != nil {
			return err
		}

		isMissed := false
		_, err = bm.client.GetSlotHeader(uint64(proposer.Slot))
		if errors.Is(err, client.ErrNotFound) {
			// slot missed
			isMissed = true
		} else if err != nil {
			log.Warnw("GetSlotHeader", "proposer", proposer.Slot, "err", err)
			return err
		}

		log.Infow("GetSlotHeader", "proposer", proposer.Pubkey, "slot", proposer.Slot, "blockNumber", blockNumber, "isMissed", isMissed)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-bm.close:
			return fmt.Errorf("close")
		default:
			fetchBlocks <- BlockInfo{
				BlockNumber: blockNumber,
				Epoch:       fromEpoch,
				Slot:        uint64(proposer.Slot),
				PubKey:      proposer.Pubkey,
				Index:       uint64(proposer.ValidatorIndex),
				IsMissed:    isMissed,
			}
		}
	}
	return nil
}

func (bm *BeaconMonitor) getEth1ExBlock(slot uint64) (uint64, error) {
	return utils.Retry(func() (uint64, error) {
		block, err := bm.client.GetBlockBySlot(slot)
		if errors.Is(err, client.ErrNotFound) {
			return bm.getEth1ExBlock(slot + 1)
		}
		if err != nil {
			return 0, err
		}

		return uint64(block.Data.Message.Body.ExecutionPayload.BlockNumber), nil
	}, utils.DefaultRetryConfig)

}

func slotToEpoch(slot uint64) uint64 {
	slotsPerEpoch := uint64(32)
	return slot / slotsPerEpoch
}

func (bm *BeaconMonitor) handleBlocks(fetchBlocks <-chan BlockInfo) (uint64, error) {
	var lastProcessedSlot uint64
	for block := range fetchBlocks {
		pubKey := removePubKeyPrefix(block.PubKey)
		validatorInfo, err := bm.store.GetValidatorByPubKeyAndBlock(pubKey, block.BlockNumber)
		if err != nil {
			log.Warnw("GetValidatorByPubKeyAndBlock", "err", err)
			continue
		}

		if validatorInfo == nil {
			log.Infow("GetValidatorByPubKeyAndBlock: Not a validator of SSV", "pubKey", pubKey)
			continue
		}

		if validatorInfo.ValidatorIndex == store.DefaultValidatorIndex {
			err = bm.store.UpdateValidatorIndex(pubKey, int64(block.Index))
			if err != nil {
				log.Warnw("UpdateValidatorIndex", "pubKey", pubKey, "Index", block.Index, "err", err)
			}
		}

		blockNumber := uint64(0)
		if !block.IsMissed {
			blockNumber = block.BlockNumber
			if bm.isSynced.Load() {
				// propose block alarm
				bm.validatorProposeBlockAlarmChan <- alert.ValidatorProposeBlockNotify{
					Epoch:     block.Epoch,
					Slot:      block.Slot,
					ClusterId: validatorInfo.ClusterID,
					Index:     block.Index,
				}
			}
		} else {
			if bm.isSynced.Load() {
				// missed block alarm
				bm.validatorMissedBlockAlarmChan <- alert.ValidatorMissedBlockNotify{
					Epoch:     block.Epoch,
					Slot:      block.Slot,
					ClusterId: validatorInfo.ClusterID,
					Index:     block.Index,
				}
			}
		}

		log.Infow("CreateBlock", "slot", block.Slot, "pubKey", pubKey, "clusterId", validatorInfo.ClusterID, "isMiss", block.IsMissed)

		err = bm.store.CreateBlock(&store.BlockInfo{
			ClusterID:   validatorInfo.ClusterID,
			BlockNumber: blockNumber,
			Epoch:       block.Epoch,
			Slot:        block.Slot,
			Proposer:    block.Index,
			PublicKey:   pubKey,
			IsMissed:    block.IsMissed,
		})
		if err != nil {
			log.Warnw("CreateBlock", "err", err)
			continue
		}

		if lastProcessedSlot < block.Slot {
			lastProcessedSlot = block.Slot
		}
	}

	return lastProcessedSlot, nil
}

func removePubKeyPrefix(pubKey string) string {
	if has0xPrefix(pubKey) {
		pubKey = pubKey[2:]
	}
	return pubKey
}

func has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}
